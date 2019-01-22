package lua

import (
	"encoding/json"

	lua "github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
)

// PrepareFunctions executes a script file passed as argument
func (l *Lang) PrepareFunctions() *Lang {

	// preload auth
	l.env.PreloadModule("backd.auth", l.backdAuthModule)

	// preload objects
	l.env.PreloadModule("backd.objects", l.backdObjectsModule)

	// preload relations
	l.env.PreloadModule("backd.relations", l.backdRelationsModule)

	// preload rbac
	l.env.PreloadModule("backd.rbac", l.backdRBACModule)

	return l
}

// SetSession sets the correct session to be used from the backd client
func (l *Lang) SetSession(sessionID string, expiresAt int64) {
	l.b.SetSession(sessionID, expiresAt)
}

// SetAppID allows to set an application ID from Go
func (l *Lang) SetAppID(appID string) {
	l.currentAppID = appID
}

// RunFunction executes the code passed as string.
//   every function looks for a map as input, returns a map and error (if any)
func (l *Lang) RunFunction(code string, input map[string]interface{}) (output map[string]interface{}, err error) {

	// function code(input)
	luaInput := luajson.DecodeValue(l.env, input)

	// load lua code
	err = l.env.DoString(code)
	if err != nil {
		return
	}

	// execute lua: function code(input)
	err = l.env.CallByParam(lua.P{
		Fn:      l.env.GetGlobal("code"),
		NRet:    1,
		Protect: true,
	}, lua.LValue(luaInput))

	if err != nil {
		return
	}

	returned := l.env.Get(1) // returned value
	l.env.Pop(1)             // remove value

	// since we got the response encodeable then we need to reencode for the result
	returnedBytes, err := luajson.Encode(returned)
	if err != nil {
		return
	}

	err = json.Unmarshal(returnedBytes, &output)
	return

}
