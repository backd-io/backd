package lua

import (
	lua "github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
)

// module objects
func (l *Lang) backdModule(L *lua.LState) int {

	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"me": l.getUsersMe,
		// "set":    l.setRBAC,
		// "add":    l.addRBAC,
		// "remove": l.removeRBAC,
	})

	L.SetField(mod, "name", lua.LString("backd"))

	// returns the module
	L.Push(mod)
	return 1

}

func (l *Lang) getUsersMe(L *lua.LState) int {

	var (
		user map[string]interface{}
		err  error
	)

	user, err = l.b.MeMapInterface()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(luajson.DecodeValue(l.env, user))
	return 1

}
