package lua

import (
	"github.com/fernandezvara/backd/backd"
	lua "github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
)

const (
	version = "0.1"
)

// Lang is the struct that holds all language commands
type Lang struct {
	b            *backd.Backd
	env          *lua.LState
	help         map[string]man
	helpChars    int
	currentAppID string
	source       string
	deep         int
}

type man struct {
	Short   string
	Long    string
	Example string
}

// Clone returns a new instance of Lang copying current configuration of
//    backd and Lua configuration
func (l *Lang) Clone() *Lang {
	var newLang Lang
	newLang.b = l.b
	newLang.env = l.env
	return &newLang
}

// New returns an instance of a Lua interpreter with the backd
//   client inside. This allow us to integrate it in many places.
// Language initialization requires a backd client with endpoints
//   configurated
func New(backd *backd.Backd) *Lang {

	var (
		lang Lang
		err  error
	)

	lang.env = lua.NewState(lua.Options{SkipOpenLibs: true}) // TODO: remember to lang.env.Close()
	lang.help = make(map[string]man)
	lang.b = backd

	// lua package
	if err = lang.LoadLuaPackage(lua.LoadLibName, lua.OpenPackage); err != nil { // it's needed to allow open another packages
		panic(err)
	}

	// lua base language
	if err = lang.LoadLuaPackage(lua.BaseLibName, lua.OpenBase); err != nil {
		panic(err)
	}

	// lua table package
	if err = lang.LoadLuaPackage(lua.TabLibName, lua.OpenTable); err != nil {
		panic(err)
	}

	// lua math package
	if err = lang.LoadLuaPackage(lua.MathLibName, lua.OpenMath); err != nil {
		panic(err)
	}

	// lua string manipullation package
	if err = lang.LoadLuaPackage(lua.StringLibName, lua.OpenString); err != nil {
		panic(err)
	}

	// lua json encode/decode package
	luajson.Preload(lang.env)

	// lang.addCommonCommands()
	// lang.addBackdAdminCommands()
	// lang.addBackdFunctionsCommands()

	return &lang

}

// LoadLuaPackage imports a package to the virtual machine
func (l *Lang) LoadLuaPackage(name string, fn lua.LGFunction) error {
	return l.env.CallByParam(lua.P{
		Fn:      l.env.NewFunction(fn),
		NRet:    0,
		Protect: true,
	}, lua.LString(name))
}
