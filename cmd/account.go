// +build !desktop

package cmd

import (
	"github.com/spf13/cobra"
)

// accountCmd represents the account command
var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Group of comands concerning loophole account",
	Long:  "Parent for commands concerning your loophole account. Always use with one of subcommands",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(accountCmd)
}
