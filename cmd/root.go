package cmd

import (
	"fmt"
	"notify-chat/pkgs"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ias-uid",
	Short: "Get user info from Keycloak",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(getFedUserIdCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "not yet",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("1.3.0-dev.2 -- HEAD")
	},
}

var getFedUserIdCmd = &cobra.Command{
	Use:   "getfuid",
	Short: "Get Federal ID from username",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		kc, err := keycloak.InitKeyCloak()
		if err != nil {
			os.Exit(2)
		}
		fmt.Println(kc.GetFUIdFromUId(args[0]))
		kc.Logout()

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
