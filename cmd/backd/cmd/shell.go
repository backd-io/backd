package cmd

import (
	"os"

	"github.com/fernandezvara/backd/pkg/lua"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// shellCmd represents the shell command
var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Interactive shell to operate with the APIs",
	Long:  `Interactive shell to operate with the APIs.`,
	Run:   shellFunc,
}

func init() {
	rootCmd.AddCommand(shellCmd)
}

func shellFunc(cmd *cobra.Command, args []string) {

	tryLogin()

	shell := lua.New(cli.backd)
	os.Exit(shell.Interactive())

}

func tryLogin() (username, password, domain string) {

	var (
		err error
	)

	username = viper.GetString(configLoginUsername)
	password = viper.GetString(configLoginPassword)
	domain = viper.GetString(configLoginDomain)

	if username == "" {
		username = promptText("Username", "", min2max254)
	}

	if password == "" {
		password = promptPassword("Password", min2max254)
	}

	if domain == "" {
		domain = promptText("Domain", "", nil)
	}

	err = cli.backd.Login(username, password, domain)
	if err != nil {
		emptyLines(2)
		printError("User/Password/Domain does not match")
		emptyLines(2)
		os.Exit(3)
	}
	return
}
