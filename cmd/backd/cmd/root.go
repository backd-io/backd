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
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "backd",
	Short: "backd-cli allows to manage and use a backd instance",
	Long:  `backd-cli allows to manage and use a backd instance.`,
	Run:   runFunc,
	Args: func(cmd *cobra.Command, args []string) error {

		switch len(args) {
		case 0:
			return errors.Errorf("showing help")
		case 1:
			return nil
		default:
			return errors.Errorf("wrong number of arguments")
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var (
	flagDisableColor bool
	flagQuiet        bool
)

func init() {
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.backd.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&flagDisableColor, "no-color", "n", false, "disable color (default: false)")
	rootCmd.PersistentFlags().BoolVarP(&flagQuiet, "quiet", "q", false, "Dont't ask accessory questions (default: false)")

	// read the config from the filesystem (if exists)
	cobra.OnInitialize(initConfig)

	// initialize client
	cobra.OnInitialize(initCli)

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

func initCli() {
	// if cli is configured then build the cli helper
	if cliConfigured {
		cli = cliStruct{
			backd: backd.NewClient(viper.GetString(configURLAuth),
				viper.GetString(configURLObjects),
				viper.GetString(configURLAdmin),
				viper.GetString(configURLFunctions),
			),
		}

		if viper.GetString(configSessionID) != "" && viper.GetInt64(configSessionExpiresAt) != 0 {
			if time.Now().Unix() < viper.GetInt64(configSessionExpiresAt) {
				cli.backd.SetSession(viper.GetString(configSessionID), viper.GetInt64(configSessionExpiresAt))
				if !flagQuiet {
					fmt.Println("Using session from configuration...")
				}
			}
		}
	}

}

// global variables
var (
	// cliConfigured is a global variable to determine is the cli is already configured
	//   if false, only `init` command can be executed
	cliConfigured bool

	// homeDirectory is the current home path for the user
	homeDirectory string

	// cli is the struct that holds the commands and it's initialized at start
	cli cliStruct
)

const (
	configFile   = ".backd"
	configFormat = "yaml"
)

func isTheCliConfigured() {
	if !cliConfigured {
		emptyLines(1)
		printError("Please configure the client first using `init`")
		emptyLines(1)
		os.Exit(1)
	}
}

type cliStruct struct {
	backd     *backd.Backd
	sessionID string
	expiresAt int64
}
