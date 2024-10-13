package cmd

import (
	"fmt"
	"notify-chat/pkgs/keycloak"
	"os"

	"github.com/spf13/cobra"
	"notify-chat/handler"
)

var rootCmd = &cobra.Command{
	Use:   "ias-uid",
	Short: "Get user info from Keycloak",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		// google.Init_client()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(getFedUserIdCmd)
	rootCmd.AddCommand(chat)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "not yet",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("1.3.1-dev.2 -- HEAD")
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

var chat = &cobra.Command{
	Use:   "chat",
	Short: "botchat",
	Run: func(cmd *cobra.Command, args []string) {
		handler.Handler()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
