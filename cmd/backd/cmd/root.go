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

	"github.com/backd-io/backd/backd"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "backd",
	Short: "backd-cli allows to manage and use a backd instance",
	Long:  `backd-cli allows to manage and use a backd instance.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.backd.yaml)")

	rootCmd.PersistentFlags().BoolVar(&disableColor, "no-color", false, "disable color (default: false)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// see if we must disable color on TTY
	cobra.OnInitialize(mustDisableColor)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	var err error

	// Find home directory.
	homeDirectory, err = homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Search config in home directory with name ".backd" (without extension).
	viper.AddConfigPath(homeDirectory)
	viper.SetConfigType(configFormat)
	viper.SetConfigName(configFile)

	// If a config file is found, read it in.
	err = viper.ReadInConfig()
	if err == nil {
		cliConfigured = true
	}

}

// cliConfigured is a global variable to determine is the cli is already configured
//   if false, only `init` command can be executed
var cliConfigured bool

// homeDirectory is the current home path for the user
var homeDirectory string

const (
	configFile   = ".backd"
	configFormat = "yaml"
)

func newBackdClient() *backd.Backd {
	return backd.NewClient(viper.GetString("url.auth"),
		viper.GetString("url.objects"),
		viper.GetString("url.admin"),
	)
}

func isTheCliConfigured() {
	if !cliConfigured {
		emptyLines(1)
		printError("Please configure the client first using `init`")
		emptyLines(1)
		os.Exit(1)
	}
}
