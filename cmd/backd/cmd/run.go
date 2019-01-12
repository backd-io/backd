package cmd

import (
	"fmt"
	"os"

	"github.com/backd-io/backd/pkg/lang"
	"github.com/spf13/cobra"
)

func runFunc(cmd *cobra.Command, args []string) {

	filename := args[0]

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Println("File does not exists, exitting...")
		os.Exit(4)
	}

	tryLogin()

	shell := lang.New(cli.backd)
	os.Exit(shell.RunScript(filename))

}
