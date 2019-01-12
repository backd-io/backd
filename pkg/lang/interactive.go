package lang

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/backd-io/anko/parser"
	"github.com/backd-io/anko/vm"
	"github.com/chzyer/readline"
	homedir "github.com/mitchellh/go-homedir"
)

// Interactive starts the anko powered backd cli
func (l *Lang) Interactive() (returnCode int) {

	parser.EnableErrorVerbose()

	fmt.Println("")
	title("backd interactive shell v. %s\n", version)
	fmt.Print(`
help            - Show a brief description of each command. 
help("command") - Show longer help for the 'command'.

	`)

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

	// initial prompt
	l.currentAppID = noAppID
	rl.SetPrompt(l.promptShell(false))

	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		l.source += line
		if l.source == "exit" {
			return 0
		}

		// wait for more code?
		multiline := l.parseSource()
		rl.SetPrompt(l.promptShell(multiline))

		if multiline == false {
			// reset source for the next set
			l.source = ""
		}

	}

	return
}

func (l *Lang) promptShell(multiline bool) string {

	if multiline {
		format := "%" + strconv.Itoa(len(l.currentAppID)) + "s " + spaces(l.deep*2)
		return fmt.Sprintf(format, "...")
	}
	return fmt.Sprintf("\033[1;34m%s»\033[0m ", l.currentAppID)

}

func spaces(n int) (s string) {
	for a := 0; a < n; a++ {
		s += " "
	}
	return
}

func (l *Lang) parseSource() bool {

	stmts, err := parser.ParseSrc(l.source)

	// fmt.Println("source:", l.source)
	// fmt.Println("stmts:", stmts)
	// fmt.Println("err:", err)

	if e, ok := err.(*parser.Error); ok {
		es := e.Error()
		if strings.HasPrefix(es, "syntax error: unexpected") {
			if strings.HasPrefix(es, "syntax error: unexpected $end,") {
				if strings.HasSuffix(l.source, "{") {
					l.deep++
				}
				if strings.HasSuffix(l.source, "}") && !strings.HasSuffix(l.source, "{}") {
					l.deep--
				}
				l.source += "\n"
				return true
			}
		} else {
			if e.Pos.Column == len(l.source) && !e.Fatal {
				fmt.Fprintln(os.Stderr, e)
				// l.deep++
				l.source += "\n"
				return true
			}
			if e.Error() == "unexpected EOF" {
				// l.deep++
				l.source += "\n"
				return true
			}
		}
	}

	var v interface{}

	if err == nil {
		v, err = vm.Run(stmts, l.env)
		// _, err = vm.Run(stmts, l.env)
	}
	if err != nil {
		if e, ok := err.(*vm.Error); ok {
			fmt.Fprintf(os.Stderr, "%d:%d %s\n", e.Pos.Line, e.Pos.Column, err)
		} else if e, ok := err.(*parser.Error); ok {
			fmt.Fprintf(os.Stderr, "%d:%d %s\n", e.Pos.Line, e.Pos.Column, err)
		} else {
			fmt.Fprintln(os.Stderr, err)
		}
	}
	if v != nil {
		if strings.HasPrefix(reflect.TypeOf(v).String(), "func") {
			fmt.Printf("Usage: %s\n", reflect.TypeOf(v).String())
		}
	}

	l.deep = 0
	return false
}

func title(text string, args ...interface{}) {

	t := fmt.Sprintf(text, args...)
	fmt.Printf("\033[1;34m%s\033[0m", t)

}
