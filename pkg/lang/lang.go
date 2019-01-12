package lang

import (
	"fmt"
	"os"

	"github.com/backd-io/anko/vm"
	"github.com/backd-io/backd/backd"
)

const (
	version = "0.0.1"
)

// Lang is the struct that holds all language commands
type Lang struct {
	b            *backd.Backd
	env          *vm.Env
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

// New retuns an instance of the `backd` scripting language
//   backd language is nothing without anko language by mattn
//   github.com/mattn/anko
// Language initialization requires a backd client with endpoints
//   configurated
func New(backd *backd.Backd) *Lang {

	var (
		lang Lang
	)

	lang.env = vm.NewEnv()
	lang.help = make(map[string]man)
	lang.b = backd
	lang.addCommonCommands()
	lang.addBackdAdminCommands()
	lang.addBackdFunctionsCommands()

	return &lang

}

// AddCommand is the helper to add a command to the parser that also creates
//   the help
func (l *Lang) AddCommand(cmd, shortHelp, longHelp, example string, fn interface{}) {

	var err error

	l.help[cmd] = man{
		Short:   shortHelp,
		Long:    longHelp,
		Example: example,
	}

	// commodity to format help
	if len(cmd) > l.helpChars {
		l.helpChars = len(cmd)
	}

	err = l.env.Define(cmd, fn)
	if err != nil {
		fmt.Println(err)
		os.Exit(99)
	}

}
