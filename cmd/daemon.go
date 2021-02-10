// +build !desktop

package cmd

import (
	"github.com/spf13/cobra"
)

// daemonCmd represents the account command
var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Group of comands concerning loophole daemon",
	Long:  "Parent for commands concerning daemonized loophole. Always use with one of subcommands",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(daemonCmd)
}
