package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/backd-io/backd/pkg/lua"
	"github.com/spf13/cobra"
)

func runFunc(cmd *cobra.Command, args []string) {

	var (
		filename string
		err      error
	)

	filename, err = filepath.Abs(args[0])
	if err != nil {
		fmt.Println("File appers to not to have a correct path. Detailed error:", err)
		os.Exit(4)
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Println("File does not exists, exitting...")
		os.Exit(4)
	}

	tryLogin()

	shell := lua.New(cli.backd)
	os.Exit(shell.RunScript(filename))

}
