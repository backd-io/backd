package lua

import (
	"encoding/json"

	"github.com/backd-io/backd/backd"
	lua "github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
)

// module objects
func (l *Lang) backdRBACModule(L *lua.LState) int {

	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"get":    l.getRBAC,
		"set":    l.setRBAC,
		"add":    l.addRBAC,
		"remove": l.removeRBAC,
	})

	L.SetField(mod, "name", lua.LString("rbac"))

	// returns the module
	L.Push(mod)
	return 1

}

func (l *Lang) getRBAC(L *lua.LState) int {

	var (
		from []byte
		rbac backd.RBAC
		err  error
	)

	obj := L.CheckAny(1)

	from, err = luajson.Encode(obj)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	err = json.Unmarshal(from, &rbac)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	err = l.b.RBAC(l.currentAppID).Get(rbac)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(luajson.DecodeValue(l.env, rbac))
	return 1

}

func (l *Lang) setRBAC(L *lua.LState) int {

	var (
		from []byte
		rbac backd.RBAC
		err  error
	)

	obj := L.CheckAny(1)

	from, err = luajson.Encode(obj)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	err = json.Unmarshal(from, &rbac)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	err = l.b.RBAC(l.currentAppID).Set(rbac)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LTrue)
	return 1

}

func (l *Lang) addRBAC(L *lua.LState) int {

	var (
		from []byte
		rbac backd.RBAC
		err  error
	)

	obj := L.CheckAny(1)

	from, err = luajson.Encode(obj)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	err = json.Unmarshal(from, &rbac)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	err = l.b.RBAC(l.currentAppID).Add(rbac)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LTrue)
	return 1

}

func (l *Lang) removeRBAC(L *lua.LState) int {

	var (
		from []byte
		rbac backd.RBAC
		err  error
	)

	obj := L.CheckAny(1)

	from, err = luajson.Encode(obj)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	err = json.Unmarshal(from, &rbac)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	err = l.b.RBAC(l.currentAppID).Remove(rbac)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LTrue)
	return 1

}
