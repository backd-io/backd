package lang

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/tidwall/pretty"
)

func (l *Lang) addCommonCommands() {

	l.AddCommand("help",
		"Shows help of the backd language, help(\"summary\") for a list of commands.",
		`Shows help of the backd language. This is the long help`,
		l.showHelp)

	l.AddCommand("println",
		"Prints any data using its default format. Appends a newline.",
		`Prints any data using its default format. Appends a newline.`,
		printLN)

	l.AddCommand("json",
		"Returns an object as JSON.",
		`Returns an object as JSON.`,
		toJSON)

	l.AddCommand("pretty",
		"Returns an object as JSON but pretty printed.",
		`Returns an object as JSON but pretty printed.`,
		prettifyString)

}

func (l *Lang) showHelp(cmd string) {

	var help string

	switch cmd {
	case "summary", "":
		//format := fmt.Sprintf("%%s%%%ds | %%s\n", l.helpChars)
		fmt.Println()
		var (
			keys []string
			key  string
		)
		// sort keys
		for key = range l.help {
			keys = append(keys, key)
		}

		sort.Strings(keys)
		for _, key = range keys {
			//help = fmt.Sprintf(format, help, key, l.help[key].Short)
			title("%20s", key)
			fmt.Printf(" | %s\n", l.help[key].Short)
		}
	default:
		value, ok := l.help[cmd]
		if !ok {
			help = fmt.Sprintf("Error: Command '%s' unknown.\n", cmd)
			break
		}

		help = fmt.Sprintf("Command '%s'\n\n%s\n", cmd, value.Long)
	}

	fmt.Print(help)
	fmt.Println("")

}

// printLN is a fmt.Println without return to simplify the function
func printLN(a ...interface{}) {
	fmt.Println(a...)
}

func toJSON(i interface{}) string {
	b, err := json.Marshal(i)
	fmt.Println("json.err:", err)
	return string(b)
}

// prettifyString returns the string generated
func prettifyString(item interface{}) string {

	by, err := json.Marshal(item)
	if err != nil {
		return "{}"
	}
	return string(pretty.Pretty(by))

}
