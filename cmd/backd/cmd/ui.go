package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

func mustDisableColor() {
	color.NoColor = flagDisableColor
}

func emptyLines(n int) {
	for a := 0; a < n; a++ {
		fmt.Println("")
	}
}

func printError(text string) {
	printlnColor(text, true, color.FgRed)
}

func printSuccess(text string) {
	printlnColor(text, true, color.FgGreen)
}

func printColor(text string, bold bool, aColor color.Attribute, args ...color.Attribute) {

	c := color.New(aColor)
	if bold {
		c.Add(color.Bold)
	}
	c.Add(args...).Print(text)

}

func printfColor(text string, bold bool, aColor color.Attribute, args ...interface{}) {

	c := color.New(aColor)
	if bold {
		c.Add(color.Bold)
	}
	c.Printf(text, args...)

}

func printlnColor(text string, bold bool, aColor color.Attribute, args ...color.Attribute) {

	c := color.New(aColor)
	if bold {
		c.Add(color.Bold)
	}
	c.Add(args...).Println(text)

}

func title(message string) {
	printlnColor(message, true, color.FgGreen, color.Underline)
}

func promptOptionsBool(promptLabel, trueString, falseString string, initialValue bool) bool {

	var (
		items  []string
		result string
		err    error
	)

	if initialValue {
		items = []string{trueString, falseString}
	} else {
		items = []string{falseString, trueString}
	}

	prompt := promptui.Select{
		Label: promptLabel,
		Items: items,
	}

	_, result, err = prompt.Run()
	er(err)

	if result == trueString {
		return true
	}
	return false

}

func promptOptions(promptLabel, cancelLabel string, promptItems []string) string {

	var (
		prompt promptui.Select
		result string
		err    error
	)

	prompt = promptui.Select{
		Label: promptLabel,
		Items: promptItems,
	}

	_, result, err = prompt.Run()
	er(err)

	if cancelLabel != "" {
		if result == cancelLabel {
			os.Exit(2)
		}
	}

	return result
}

func promptText(promptLabel, promptDefault string, validateFunc promptui.ValidateFunc) string {

	var (
		prompt promptui.Prompt
		result string
		err    error
	)

	prompt.Label = promptLabel
	if promptDefault != "" {
		prompt.Default = promptDefault
	}

	if validateFunc != nil {
		prompt.Validate = validateFunc
	}

	result, err = prompt.Run()
	er(err)

	return result

}

func promptPassword(promptLabel string, validateFunc promptui.ValidateFunc) string {

	var (
		prompt promptui.Prompt
		result string
		err    error
	)

	prompt.Label = promptLabel
	prompt.Mask = '*'
	if validateFunc != nil {
		prompt.Validate = validateFunc
	}

	result, err = prompt.Run()
	er(err)

	return result

}

// VALIDATIONS
func validateURL(input string) error {

	if govalidator.IsURL(input) {
		return nil
	}

	return errors.New("Invalid URL")

}

func min2max254(input string) error {
	return minmax(input, 2, 254)
}

func min2max32(input string) error {
	return minmax(input, 2, 32)
}

func max254(input string) error {
	return minmax(input, 0, 254)
}

func minmax(input string, min, max int) error {
	l := len(input)
	if l < min+1 && min != 0 {
		return errors.New("Too short")
	}
	if l > max && max != 0 {
		return errors.New("Too long")
	}
	return nil
}

func isEmail(input string) error {
	if !govalidator.IsEmail(input) {
		return errors.New("Not a valid Email")
	}
	return nil
}

// er helper
func er(err error) {
	if err != nil {
		printError("Unexpected error:")
		printError(err.Error())
		os.Exit(1)
	}
}
