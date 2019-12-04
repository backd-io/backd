package cmd

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"os"

// 	"github.com/backd-io/backd/backd"
// 	"github.com/backd-io/backd/internal/utils"
// 	"github.com/spf13/cobra"
// )

// // variables for applications commands
// var (
// 	flagApplicationID string
// 	flagAPI           string
// 	flagRunAs         string
// )

// // applicationsCmd represents the applications command
// var applicationsCmd = &cobra.Command{
// 	Use:   "applications",
// 	Short: "Helper commands for applications workflow.",
// 	Long:  `Helper commands for applications workflow.`,
// 	Run:   applicationsFunc,
// }

// // subcommands
// var lsCmd = &cobra.Command{
// 	Use:   "ls",
// 	Short: "Returns all applications the user has access to",
// 	Run:   applicationsLSFunc,
// }

// var createCmd = &cobra.Command{
// 	Use:   "create [name]",
// 	Short: "Creates a new application",
// 	Run:   applicationsCreateFunc,
// 	Args:  cobra.MinimumNArgs(2),
// }

// var updateCmd = &cobra.Command{
// 	Use:   "update",
// 	Short: "Updates a function from a source file",
// 	Run:   applicationsUpdateFunc,
// }

// var deleteCmd = &cobra.Command{
// 	Use:   "delete",
// 	Short: "Deletes a function",
// 	Run:   applicationsDeleteFunc,
// }

// func init() {
// 	applicationsCmd.AddCommand(lsCmd)
// 	applicationsCmd.AddCommand(createCmd)
// 	applicationsCmd.AddCommand(updateCmd)
// 	applicationsCmd.AddCommand(deleteCmd)
// 	rootCmd.AddCommand(applicationsCmd)
// 	lsCmd.Flags().StringVarP(&flagApplicationID, "application-id", "a", "", "Application ID to work with")
// 	createCmd.Flags().StringVarP(&flagApplicationID, "application-id", "a", "", "Application ID to work with")
// 	createCmd.Flags().StringVar(&flagAPI, "api", "", "Queryable by applications API?")
// 	createCmd.Flags().StringVar(&flagRunAs, "run-as", "", "User that runs the function on behalf of the logged in user")
// 	updateCmd.Flags().StringVarP(&flagApplicationID, "application-id", "a", "", "Application ID to work with")
// 	updateCmd.Flags().StringVar(&flagAPI, "api", "", "Queryable by applications API?")
// 	updateCmd.Flags().StringVar(&flagRunAs, "run-as", "", "User that runs the function on behalf of the logged in user")
// 	deleteCmd.Flags().StringVarP(&flagApplicationID, "application-id", "a", "", "Application ID to work with")
// }

// func applicationsFunc(cmd *cobra.Command, args []string) {
// 	fmt.Println("applications help")
// }

// func queryApplicationID() {
// 	if flagApplicationID == "" {
// 		flagApplicationID = promptText("application-id", "", min2max32)
// 		if flagApplicationID == "" {
// 			printError("Application ID is required")
// 			os.Exit(2)
// 		}
// 	}
// }

// func queryAPI() bool {
// 	if flagAPI == "" {
// 		flagAPI = promptOptions("Publish as API?", "", []string{"Yes", "No"})
// 	}

// 	if flagAPI == "Yes" {
// 		return true
// 	}

// 	return false
// }

// func queryRunAs() {
// 	if flagRunAs == "" {
// 		flagRunAs = promptText("run as user (blank for none)", "", max254)
// 	}
// }

// func applicationsLSFunc(cmd *cobra.Command, args []string) {

// 	var (
// 		err       error
// 		applications []backd.Function
// 	)

// 	tryLogin()
// 	queryApplicationID()

// 	err = cli.backd.Functions(flagApplicationID).GetMany(backd.QueryOptions{}, &applications)
// 	if err != nil {
// 		printError(err.Error())
// 		os.Exit(2)
// 	}

// 	utils.Prettify(applications)

// 	fmt.Printf("`applications ls` with args: %v\n", args)

// }

// func applicationsCreateFunc(cmd *cobra.Command, args []string) {

// 	var (
// 		err      error
// 		fileByte []byte
// 		function backd.Function
// 	)

// 	tryLogin()
// 	queryApplicationID()
// 	function.API = queryAPI()
// 	queryRunAs()

// 	function.Name = args[0]
// 	function.RunAs = flagRunAs

// 	fileByte, err = ioutil.ReadFile(args[1])
// 	if err != nil {
// 		printError(err.Error())
// 		os.Exit(2)
// 	}

// 	function.Code = string(fileByte)

// 	_, err = cli.backd.Functions(flagApplicationID).Insert(function)
// 	er(err)
// 	utils.Prettify(function)
// 	fmt.Printf("`applications create` with args: %v\n", args)
// }

// func applicationsUpdateFunc(cmd *cobra.Command, args []string) {
// 	fmt.Printf("`applications update` with args: %v\n", args)
// }

// func applicationsDeleteFunc(cmd *cobra.Command, args []string) {
// 	fmt.Printf("`applications delete` with args: %v\n", args)
// }
