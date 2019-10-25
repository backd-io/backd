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
		err := cli.backd.Logout()

		switch err {
		case nil:
			if !flagQuiet {
				printSuccess("User logged out successfully")
			}
			os.Exit(0)
		default:
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
