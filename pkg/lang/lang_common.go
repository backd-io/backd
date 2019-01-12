package lang

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/tidwall/pretty"
)

func (l *Lang) addCommonCommands() {

	var verbs = `
Boolean:
%t	the word true or false

Integer:
%c	the character represented by the corresponding Unicode code point
%d	base 10 (%o	base 8) (%b	base 2)
%q	a single-quoted character literal safely escaped with Go syntax.
%x	base 16, with lower-case letters for a-f (%X for upper-case letters)
%U	Unicode format: U+1234; same as "U+%04X"

Floating-point and complex constituents:
%b	decimalless scientific notation with exponent a power of two,
	in the manner of strconv.FormatFloat with the 'b' format,
	e.g. -123456p-78
%e	scientific notation, e.g. -1.234456e+78
%E	scientific notation, e.g. -1.234456E+78
%f	decimal point but no exponent, e.g. 123.456

String and slice of bytes (treated equivalently with these verbs):
%s	the uninterpreted bytes of the string or slice
%q	a double-quoted string safely escaped with Go syntax
%x	base 16, lower-case, two characters per byte
%X	base 16, upper-case, two characters per byte

Modificators / helpers:
%%       it's the escape form to print a single '%'
%+<t>    always add the sign if numeric value 
%nn<t>   write the data type <t> with the width of nn
%-nn<t>  write the data type <t> padded to the right with width of nn
`

	l.AddCommand("help",
		"Shows help of the backd language, help(\"summary\") for a list of commands.",
		`
Shows help of the backd language. This is the long help
		`,
		`
help("println"
		`,
		l.showHelp)

	l.AddCommand("println",
		"Prints any data using its default format. Appends a newline.",
		`
Prints any data evaluating its default format, also appends a newline.

It allows to mix any kind of variables and print them together without care about formatting. Every item in the final string with be separed with a space directly.
		`,
		`
// single item
println("Hello, world!")

// many items 
println("There are", 10, "users")
		`,
		printLN)

	l.AddCommand("printf",
		"Prints any data using the desired format.",
		`
Prints any data using its desired format. Format verbs:

`+verbs+`		
    `,
		`
// string -> Hello, world!
println("Hello, %s!, "world")

// integer
println("There are %d users.", 10)

// longer sample
items = [{"desc": "item with text", "amount": 4}, {"desc": "another item", "amount": 1232}]
for item in items {
  printf("| %-20s | %5d€ |", item.desc, item.amount)
}

// returns:
// | item with text       |     4€ |
// | another item         |  1232€ |

		`,
		printF)

	l.AddCommand("sprintf",
		"Returns a string formatted as desired.",
		`
Returns a string formatted as desired. Format verbs:

`+verbs+`		
		`,
		`
aText = sprintf("%+d items.", 37)
aText2 = sprintf("%+d items.", -37)
println(aText)
println(aText2)

// prints:
// +37 items.
// -37 items.
		`,
		sprintF)

	l.AddCommand("json",
		"Returns an object as JSON.",
		`
Returns an object formatted as JSON. Does not print it to the console.
`,
		`
a = {}
a.text = "this is a text"
a.number = 1234
a.boolean = true 

j = json(a)
println(j)

// returns
// {"boolean":true,"number":1234,"text":"this is a text"}

		`,
		toJSON)

	l.AddCommand("pretty",
		"Prints an object to the console prettifying the output for human.",
		`
Prints an object to the console prettifying the output for human. Useful for debug using shell or for write information as output of scripts.
		`,
		`
a = {}
a.text = "this is a text"
a.number = 1234
a.boolean = true 

pretty(a)

// returns: 
// {
//   "boolean": true,
//   "number": 1234,
//   "text": "this is a text"
// }
		`,
		prettifyString)

}

func (l *Lang) showHelp(cmd string) {

	var help string

	switch cmd {
	case "summary", "":
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
			title("%20s", key)
			fmt.Printf(" | %s\n", l.help[key].Short)
		}
	default:
		value, ok := l.help[cmd]
		if !ok {
			help = fmt.Sprintf("\033[1;31mError\033[0m: Command '%s' unknown.\n\n", cmd)
			break
		}

		help = fmt.Sprintf("\033[1;34mCommand '%s'\033[0m\n\n%s\n\n\033[1;34mExample:\033[0m\n%s\n\n", cmd, value.Long, value.Example)
	}
	fmt.Print(help)

}

// printLN is a fmt.Println without return to simplify the function
func printLN(a ...interface{}) {
	fmt.Println(a...)
}

// printF is a fmt.Printf without return to simplify the function
func printF(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
}

// sprintF is a fmt.Sprintf without return to simplify the function
func sprintF(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}

func toJSON(i interface{}) string {

	b, err := json.Marshal(i)
	if err != nil {
		return "{}"
	}
	return string(b)

}

// prettifyString prints to console the object prettified
func prettifyString(item interface{}) {

	by, err := json.Marshal(item)
	if err != nil {
		fmt.Println("{}")
	}
	fmt.Println(string(pretty.Pretty(by)))

}
