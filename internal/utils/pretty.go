package utils

import (
	"encoding/json"
	"fmt"

	"github.com/tidwall/pretty"
)

// Prettify writes the item formated a prettified to the stdin, if err it writes the error
func Prettify(item interface{}) {

	by, err := json.Marshal(item)
	if err != nil {
		fmt.Println("Prettify failed. Error:")
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("%s\n", string(pretty.Pretty(by)))

}

// PrettifyString returns the string generated
func PrettifyString(item interface{}) string {

	by, err := json.Marshal(item)
	if err != nil {
		return "{}"
	}
	return string(pretty.Pretty(by))

}
