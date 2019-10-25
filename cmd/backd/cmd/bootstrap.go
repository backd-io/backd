package cmd

import (
	"os"

	"github.com/fatih/color"
	"github.com/fernandezvara/backd/backd"
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
		request backd.BootstrapRequest
		err     error
	)

	isTheCliConfigured()

	emptyLines(2)
	title("backd cluster first configuration.")
	title("This operation will make the first user on the cluster, that will become `Domain and Application Admin`.")
	emptyLines(2)

	request.Code = promptText("Bootstrap Code", "", nil)
	request.Name = promptText("Your name", "John Doe", max254)
	request.Username = promptText("Username", "john.doe", min2max254)
	request.Email = promptText("Email", "john.doe@example.com", isEmail)
	request.Password = promptPassword("Password", nil)

	err = cli.backd.BootstrapCluster(request.Code,
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
