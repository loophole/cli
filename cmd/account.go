package cmd

import (
	"github.com/spf13/cobra"
)

// accountCmd represents the account command
var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Parent for commands concerning your loophole account. Always use with one of the following: login, logout",
	Long:  "Parent for commands concerning your loophole account. Always use with one of the following: login, logout",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(accountCmd)
}
