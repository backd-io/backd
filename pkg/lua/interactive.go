package lua

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
	homedir "github.com/mitchellh/go-homedir"
	lua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
)

// Interactive starts the lua powered backd cli
func (l *Lang) Interactive() int {

	// load lua os package
	if err := l.LoadLuaPackage(lua.OsLibName, lua.OpenOs); err != nil {
		panic(err)
	}

	// load lua io package
	if err := l.LoadLuaPackage(lua.IoLibName, lua.OpenIo); err != nil {
		panic(err)
	}

	// load lua debug package
	if err := l.LoadLuaPackage(lua.DebugLibName, lua.OpenDebug); err != nil {
		panic(err)
	}

	// load lua coroutines package
	if err := l.LoadLuaPackage(lua.CoroutineLibName, lua.OpenCoroutine); err != nil {
		panic(err)
	}

	// preload auth
	l.env.PreloadModule("backd.auth", l.backdAuthModule)

	// preload objects
	l.env.PreloadModule("backd.objects", l.backdObjectsModule)

	// preload relations
	l.env.PreloadModule("backd.relations", l.backdRelationsModule)

	// preload rbac
	l.env.PreloadModule("backd.rbac", l.backdRBACModule)

	// Find home directory.
	homeDirectory, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	rl, err := readline.NewEx(&readline.Config{
		DisableAutoSaveHistory: false,
		HistoryFile:            fmt.Sprintf("%s/%s", homeDirectory, ".backd.history"),
		InterruptPrompt:        "^C",
		EOFPrompt:              "exit",
		HistorySearchFold:      true,
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	// set up interactive
	l.currentAppID = noAppID
	rl.SetPrompt(l.promptShell(false))
	// add commands
	l.addCliCommands()

	for {
		if str, err := l.loadline(rl, l.env); err == nil {
			if err := l.env.DoString(str); err != nil {
				fmt.Println(err)
			}
		} else { // error on loadline
			fmt.Println(err)
			return 1
		}
	}

	// unreachable but it's good to respect your own rules
	// return 0

}

func (l *Lang) incomplete(err error) bool {
	if lerr, ok := err.(*lua.ApiError); ok {
		if perr, ok := lerr.Cause.(*parse.Error); ok {
			return perr.Pos.Line == parse.EOF
		}
	}
	return false
}

func (l *Lang) loadline(rl *readline.Instance, L *lua.LState) (string, error) {

	var (
		line string
		err  error
	)
	rl.SetPrompt(l.promptShell(false))
	if line, err = rl.Readline(); err == nil {
		// exit gracefully
		if strings.TrimSpace(line) == "exit" || strings.TrimSpace(line) == "exit()" {
			rl.Close()
			l.env.Close()
			os.Exit(0)
		}
		if _, err = L.LoadString("return " + line); err == nil { // try add return <...> then compile
			return line, nil
		}
		return l.multiline(line, rl, L)
	}
	// else
	return "", err
}

func (l *Lang) multiline(ml string, rl *readline.Instance, L *lua.LState) (string, error) {
	for {
		if _, err := L.LoadString(ml); err == nil { // try compile
			return ml, nil
		} else if !l.incomplete(err) { // syntax error , but not EOF
			return ml, nil
		} else {
			rl.SetPrompt(l.promptShell(true))
			if line, err := rl.Readline(); err == nil {
				ml = ml + "\n" + line
			} else {
				return "", err
			}
		}
	}
}

func (l *Lang) promptShell(multiline bool) string {

	if multiline {
		format := "%" + strconv.Itoa(len(l.currentAppID)) + "s» " //+ spaces(l.deep*2)
		return fmt.Sprintf(format, "")
	}
	return fmt.Sprintf("\033[1;34m%s»\033[0m ", l.currentAppID)

}

func spaces(n int) (s string) {
	for a := 0; a < n; a++ {
		s += " "
	}
	return
}
