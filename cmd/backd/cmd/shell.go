// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"os"

	"github.com/backd-io/backd/pkg/lang"
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

	var (
		username string
		password string
		domain   string
		err      error
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

	shell := lang.New(cli.backd)
	os.Exit(shell.Interactive())

}
