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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logs the user off the API",
	Long:  `Logs the user off the API`,
	Run:   logoutFunc,
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}

func logoutFunc(cmd *cobra.Command, args []string) {

	sessionID := viper.GetString(configSessionID)
	expiresAt := viper.GetInt64(configSessionExpiresAt)

	if sessionID != "" && expiresAt != 0 {
		cli.backd.SetSession(sessionID, expiresAt)
		success, err := cli.backd.Logout()

		switch success {
		case true:
			if !flagQuiet {
				printSuccess("User logged out successfully")
			}
			os.Exit(0)
		case false:
			if !flagQuiet {
				printError("Error logging out user")
				if err != nil {
					printError(err.Error())
				}
			}
			os.Exit(1)
		}
	}

	viper.Set(configSessionID, "")
	viper.Set(configSessionExpiresAt, 0)

	// if no session stored on the cli....
	if !flagQuiet {
		printError("No session found on configuration. Exitting...")
	}
	os.Exit(1)

}
