package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// accountCmd represents the account command
var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Parent for commands concerning your loophole account. Always use with one of the following: login, logout",
	Long:  "Parent for commands concerning your loophole account. Always use with one of the following: login, logout",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("No subcommand specified. Exiting.")
	},
}

func init() {
	rootCmd.AddCommand(accountCmd)
}
