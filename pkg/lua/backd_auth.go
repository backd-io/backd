package lua

import (
	"time"

	lua "github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
)

// module objects
func (l *Lang) backdAuthModule(L *lua.LState) int {

	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"me":                     l.getUsersMe,
		"login":                  l.login,
		"logout":                 l.logout,
		"get_session_id":         l.getSessionID,
		"get_session_state":      l.getSessionState,
		"get_session_expiration": l.getSessionExpiresAt,
		"set_session":            l.setSession,
		"set_session_id":         l.setSessionID,
		// "add":    l.addRBAC,
		// "remove": l.removeRBAC,
	})

	L.SetField(mod, "name", lua.LString("auth"))

	// returns the module
	L.Push(mod)
	return 1

}

func (l *Lang) getUsersMe(L *lua.LState) int {

	var (
		userMap map[string]interface{}
		err     error
	)

	userMap, err = l.b.MeMapInterface()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(luajson.DecodeValue(l.env, userMap))
	return 1

}

func (l *Lang) login(L *lua.LState) int {

	var (
		domain   string
		username string
		password string
		err      error
	)

	domain = L.ToString(1)
	username = L.ToString(2)
	password = L.ToString(3)

	err = l.b.Login(username, password, domain)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LTrue)
	return 1

}

func (l *Lang) logout(L *lua.LState) int {

	var (
		err error
	)

	err = l.b.Logout()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LTrue)
	return 1

}

func (l *Lang) getSessionID(L *lua.LState) int {

	var (
		sessionID string
	)

	sessionID, _, _ = l.b.Session()

	L.Push(lua.LString(sessionID))
	return 1

}

func (l *Lang) getSessionState(L *lua.LState) int {

	var (
		sessionState int
	)

	_, sessionState, _ = l.b.Session()

	L.Push(lua.LNumber(float64(sessionState)))
	return 1

}

func (l *Lang) getSessionExpiresAt(L *lua.LState) int {

	var (
		expiresAt time.Time
	)

	_, _, expiresAt = l.b.Session()

	L.Push(lua.LNumber(float64(expiresAt.Unix())))
	return 1

}

func (l *Lang) setSession(L *lua.LState) int {

	var (
		sessionID string
		expiresAt int64
	)

	sessionID = L.ToString(1)
	expiresAt = int64(L.ToNumber(2))

	l.b.SetSession(sessionID, expiresAt)

	L.Push(lua.LTrue)
	return 1

}

func (l *Lang) setSessionID(L *lua.LState) int {

	var (
		sessionID string
	)

	sessionID = L.ToString(1)

	l.b.SetSessionID(sessionID)

	L.Push(lua.LTrue)
	return 1

}

// func (b *Backd) Session() (string, int, time.Time) {
// func (b *Backd) SetSession(sessionID string, expiresAt int64) {
