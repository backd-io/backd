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
	"fmt"
	"os"
	"time"

	"github.com/backd-io/backd/backd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	flagSaveSession bool
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the API",
	Long:  `Before be able to use the API you need to log in with an user.`,
	Run:   loginFunc,
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().BoolVarP(&flagSaveSession, "save-session", "s", false, "saves the session information for later usage")
}

func loginFunc(cmd *cobra.Command, args []string) {

	var (
		username              string
		password              string
		domain                string
		saveSessionIDQuestion string
	)

	username, password, domain = tryLogin()

	sessionID, state, expiresAt := cli.backd.Session()
	if state != backd.StateLoggedIn {
		printError("session state unexpected")
	}

	// if the user wants to be 'quiet' be it!
	if flagQuiet {
		fmt.Printf("%s", sessionID)
		os.Exit(0)
	}

	emptyLines(2)
	printSuccess("Login successful")
	printSuccess(fmt.Sprintf("SessionID: '%s', Expires in %s (%s)", sessionID, time.Until(expiresAt), expiresAt))
	emptyLines(2)

	if !flagQuiet {

		// save session can be passed as flag, as configuration parameter or interactive answer
		if !viper.GetBool(configCliSaveSession) && !flagSaveSession {
			saveSessionIDQuestion = promptOptions("Do you want to save the session on the configuration?", "", []string{answerYes, answerAlways, answerNo, answerNever})
		} else {
			saveSessionIDQuestion = answerYes // no need to rewrite the 'always save'
		}

		if saveSessionIDQuestion == answerAlways || saveSessionIDQuestion == answerYes {

		}

		// stop asking...
		switch saveSessionIDQuestion {
		case answerYes, answerAlways:
			viper.Set(configSessionID, sessionID)
			viper.Set(configSessionExpiresAt, expiresAt.Unix())
			if saveSessionIDQuestion == answerAlways {
				viper.Set(configCliSaveSession, true)
				viper.Set(configCliDontAskSession, false)
			}
			viper.WriteConfig()
		case answerNever:
			viper.Set(configSessionID, "")
			viper.Set(configSessionExpiresAt, 0)
			viper.Set(configCliSaveSession, false)
			viper.Set(configCliDontAskSession, true)
			viper.WriteConfig()
		}

		if viper.Get(configLoginUsername) != username || viper.Get(configLoginDomain) != domain {
			saveUsernameDomain := promptOptions("Do you want to save the username/domain on configuration for later?", "", []string{answerYes, answerNo, answerNever})
			switch saveUsernameDomain {
			case answerYes:
				viper.Set(configLoginUsername, username)
				viper.Set(configLoginDomain, domain)
				viper.WriteConfig()
				savePassword := promptOptions("Do you want to save the password on configuration for later? (DANGEROUS!)", "", []string{answerYes, answerNo, answerNever})
				switch savePassword {
				case answerYes:
					viper.Set(configLoginPassword, password)
					viper.WriteConfig()
				case answerNever:
					viper.Set(configCliDontAskPassword, true)
					viper.WriteConfig()
				}
			case answerNever:
				viper.Set(configCliDontAskUserDomain, true)
				viper.Set(configCliDontAskPassword, true)
				viper.WriteConfig()
			}
		}
	}

}
