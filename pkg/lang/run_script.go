package lang

import (
	"fmt"
	"io/ioutil"
)

// RunScript executes a script file passed as argument
func (l *Lang) RunScript(filename string) int {

	sourceBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Could not read the script. Detailed error:", err)
		return 5
	}

	_, err = l.env.Execute(string(sourceBytes))
	if err != nil {
		fmt.Println("Could not execute the script. Detailed error:", err)
		return 5
	}

	return 0

}
