package lua

import (
	"encoding/json"
	"fmt"

	"github.com/fernandezvara/backd/internal/utils"
	lua "github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
)

// module functions
func (l *Lang) backdFunctionsModule(L *lua.LState) int {

	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"func": l.execFunction,
	})

	L.SetField(mod, "name", lua.LString("functions"))

	// returns the module
	L.Push(mod)
	return 1

}

func (l *Lang) execFunction(L *lua.LState) int {

	var (
		id        string
		obj       lua.LValue
		inputByte []byte
		input     map[string]interface{}
		output    map[string]interface{}
		err       error
	)

	id = L.ToString(1)
	obj = L.CheckAny(2)

	inputByte, err = luajson.Encode(obj)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	err = json.Unmarshal(inputByte, &input)
	if err != nil {
		// get an empty map
		input = map[string]interface{}{}
	}

	fmt.Println("input:")
	utils.Prettify(input)

	output, err = l.b.Functions(l.currentAppID).Run(id, input)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(luajson.DecodeValue(l.env, output))
	L.Push(lua.LString(""))
	return 2

}
