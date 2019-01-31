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

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes the backd-cli to use a backd instance.",
	Long:  `Initializes the backd-cli to use a backd instance.`,
	Run:   initFunc,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func initFunc(cmd *cobra.Command, args []string) {

	emptyLines(2)
	title("backd cli initialization.")
	emptyLines(2)

	if cliConfigured {
		title("Current Configuration:")
		emptyLines(1)
		for _, i := range viper.AllKeys() {
			printColor(i, true, color.FgWhite)
			printColor(": ", true, color.FgWhite)
			printfColor("%v\n", false, color.FgWhite, viper.Get(i))
		}
		emptyLines(1)
		// if user cancels it will close so value is useless
		_ = promptOptions("Overwrite configuration?", "No", []string{"Yes", "No"})
		emptyLines(2)
	}

	authURL := promptText("Auth URL", "https://auth.backd.io", validateURL)
	objectsURL := promptText("Objects URL", "https://objects.backd.io", validateURL)
	functionsURL := promptText("Functions URL", "https://functions.backd.io", validateURL)
	adminURL := promptText("Admin URL (optional)", "https://auth.backd.io", validateURL)

	viper.Set(configURLAdmin, adminURL)
	viper.Set(configURLAuth, authURL)
	viper.Set(configURLObjects, objectsURL)
	viper.Set(configURLFunctions, functionsURL)

	if cliConfigured {
		emptyLines(2)
		printlnColor(fmt.Sprintf("Saving configuration on file: %s", viper.ConfigFileUsed()), true, color.FgWhite)
		if err := viper.WriteConfig(); err != nil {
			er(err)
		}
		// finish execution since all job is done
		os.Exit(0)
	}

	emptyLines(2)
	printlnColor(fmt.Sprintf("Saving configuration on file: %s/%s.%s", homeDirectory, configFile, configFormat), true, color.FgWhite)
	if err := viper.WriteConfigAs(fmt.Sprintf("%s/%s.%s", homeDirectory, configFile, configFormat)); err != nil {
		er(err)
	}

}
