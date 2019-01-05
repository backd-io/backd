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

	"github.com/backd-io/backd/backd"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// bootstrapCmd represents the bootstrap command
var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Configures (bootstraps) a new Backd cluster.",
	Long:  `Configures (bootstraps) a new Backd cluster.`,
	Run:   bootstrapFunc,
}

func init() {
	rootCmd.AddCommand(bootstrapCmd)
}

func bootstrapFunc(cmd *cobra.Command, args []string) {

	var (
		client  *backd.Backd
		request backd.BootstrapRequest
		err     error
	)

	isTheCliConfigured()
	client = newBackdClient()

	emptyLines(2)
	title("backd cluster first configuration.")
	title("This operation will make the first user on the cluster, that will become `Domain and Application Admin`.")
	emptyLines(2)

	request.Code = promptText("Bootstrap Code", "", nil)
	request.Name = promptText("Your name", "John Doe", max254)
	request.Username = promptText("Username", "john.doe", min2max254)
	request.Email = promptText("Email", "john.doe@example.com", isEmail)
	request.Password = promptPassword("Password", nil)

	err = client.BootstrapCluster(request.Code,
		request.Name,
		request.Username,
		request.Email,
		request.Password,
	)

	emptyLines(2)
	if err != nil {
		printError(err.Error())
		emptyLines(2)
		os.Exit(1)
	}

	printColor("Server successfully bootstrapped.", true, color.FgWhite)
	emptyLines(2)
}
