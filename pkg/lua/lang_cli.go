package lua

import (
	"encoding/json"
	"fmt"

	"github.com/tidwall/pretty"
	lua "github.com/yuin/gopher-lua"
)

// TODO: This needs to go to a backd module or something like that

// addCliCommands - add commands useful on the interactive shell (mainly for debugging)
func (l *Lang) addCliCommands() {

	l.env.SetGlobal("pretty", l.env.NewFunction(l.pretty))
	l.env.SetGlobal("appid", l.env.NewFunction(l.appID))

}

// appID set the application id to work with using the client
func (l *Lang) appID(L *lua.LState) int {

	l.currentAppID = L.ToString(1)
	return 0

}

// pretty prints to console the object prettified
func (l *Lang) pretty(L *lua.LState) int {

	item := L.ToUserData(1)

	fmt.Println("item:", item)

	by, err := json.Marshal(item)
	if err != nil {
		fmt.Println("{}")
	}
	fmt.Println(string(pretty.Pretty(by)))

	return 0

}
