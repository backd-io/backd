package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fernandezvara/backd/backd"
	"github.com/fernandezvara/backd/internal/utils"
	"github.com/spf13/cobra"
)

// variables for functions commands
var (
	flagApplicationID string
	flagAPI           string
	flagRunAs         string
)

// functionsCmd represents the functions command
var functionsCmd = &cobra.Command{
	Use:   "functions",
	Short: "Helper commands for functions workflow.",
	Long:  `Helper commands for functions workflow.`,
	Run:   functionsFunc,
}

// subcommands
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Returns all functions defined for the application",
	Run:   functionsLSFunc,
}

var createCmd = &cobra.Command{
	Use:   "create [name] [source code file]",
	Short: "Creates a new function from a source file",
	Run:   functionsCreateFunc,
	Args:  cobra.MinimumNArgs(2),
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates a function from a source file",
	Run:   functionsUpdateFunc,
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes a function",
	Run:   functionsDeleteFunc,
}

func init() {
	functionsCmd.AddCommand(lsCmd)
	functionsCmd.AddCommand(createCmd)
	functionsCmd.AddCommand(updateCmd)
	functionsCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(functionsCmd)
	lsCmd.Flags().StringVarP(&flagApplicationID, "application-id", "a", "", "Application ID to work with")
	createCmd.Flags().StringVarP(&flagApplicationID, "application-id", "a", "", "Application ID to work with")
	createCmd.Flags().StringVar(&flagAPI, "api", "", "Queryable by functions API?")
	createCmd.Flags().StringVar(&flagRunAs, "run-as", "", "User that runs the function on behalf of the logged in user")
	updateCmd.Flags().StringVarP(&flagApplicationID, "application-id", "a", "", "Application ID to work with")
	updateCmd.Flags().StringVar(&flagAPI, "api", "", "Queryable by functions API?")
	updateCmd.Flags().StringVar(&flagRunAs, "run-as", "", "User that runs the function on behalf of the logged in user")
	deleteCmd.Flags().StringVarP(&flagApplicationID, "application-id", "a", "", "Application ID to work with")
}

func functionsFunc(cmd *cobra.Command, args []string) {
	fmt.Println("functions help")
}

func queryApplicationID() {
	if flagApplicationID == "" {
		flagApplicationID = promptText("application-id", "", min2max32)
		if flagApplicationID == "" {
			printError("Application ID is required")
			os.Exit(2)
		}
	}
}

func queryAPI() bool {
	if flagAPI == "" {
		flagAPI = promptOptions("Publish as API?", "", []string{"Yes", "No"})
	}

	if flagAPI == "Yes" {
		return true
	}

	return false
}

func queryRunAs() {
	if flagRunAs == "" {
		flagRunAs = promptText("run as user (blank for none)", "", max254)
	}
}

func functionsLSFunc(cmd *cobra.Command, args []string) {

	var (
		err       error
		functions []backd.Function
	)

	tryLogin()
	queryApplicationID()

	err = cli.backd.Functions(flagApplicationID).GetMany(backd.QueryOptions{}, &functions)
	if err != nil {
		printError(err.Error())
		os.Exit(2)
	}

	utils.Prettify(functions)

	fmt.Printf("`functions ls` with args: %v\n", args)

}

func functionsCreateFunc(cmd *cobra.Command, args []string) {

	var (
		err      error
		fileByte []byte
		function backd.Function
	)

	tryLogin()
	queryApplicationID()
	function.API = queryAPI()
	queryRunAs()

	function.Name = args[0]
	function.RunAs = flagRunAs

	fileByte, err = ioutil.ReadFile(args[1])
	if err != nil {
		printError(err.Error())
		os.Exit(2)
	}

	function.Code = string(fileByte)

	_, err = cli.backd.Functions(flagApplicationID).Insert(function)
	er(err)
	utils.Prettify(function)
	fmt.Printf("`functions create` with args: %v\n", args)
}

func functionsUpdateFunc(cmd *cobra.Command, args []string) {
	fmt.Printf("`functions update` with args: %v\n", args)
}

func functionsDeleteFunc(cmd *cobra.Command, args []string) {
	fmt.Printf("`functions delete` with args: %v\n", args)
}
