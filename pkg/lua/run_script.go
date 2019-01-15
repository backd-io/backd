package lua

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

// RunScript executes a script file passed as argument
func (l *Lang) RunScript(filename string) int {

	var err error

	// load lua os package
	if err = l.LoadLuaPackage(lua.OsLibName, lua.OpenOs); err != nil {
		panic(err)
	}

	// load lua io package
	if err = l.LoadLuaPackage(lua.IoLibName, lua.OpenIo); err != nil {
		panic(err)
	}

	// do not load lua debug package
	if err = l.LoadLuaPackage(lua.DebugLibName, lua.OpenDebug); err != nil {
		panic(err)
	}

	// load lua coroutines package
	if err = l.LoadLuaPackage(lua.CoroutineLibName, lua.OpenCoroutine); err != nil {
		panic(err)
	}

	// preload backd
	l.env.PreloadModule("backd", l.backdModule)

	// preload objects
	l.env.PreloadModule("backd.objects", l.backdObjectsModule)

	// preload relations
	l.env.PreloadModule("backd.relations", l.backdRelationsModule)

	// preload rbac
	l.env.PreloadModule("backd.rbac", l.backdRBACModule)

	// set up backd
	l.currentAppID = noAppID

	// add allowed commands for shell scripts
	l.addCliCommands()

	err = l.env.DoFile(filename)
	if err != nil {
		fmt.Println("Could not execute the script. Detailed error:", err)
		return 5
	}

	return 0

}
