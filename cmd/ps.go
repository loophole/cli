// +build !desktop

package cmd

import (
	"github.com/loophole/cli/internal/app/loopholed"
	"github.com/spf13/cobra"
)

// psCmd represents the account command
var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "Group of comands concerning loophole account",
	Long:  "Parent for commands concerning your loophole account. Always use with one of subcommands",
	Run: func(cmd *cobra.Command, args []string) {
		loopholedClient := &loopholed.LoopholedClient{}

		loopholedClient.Ps()
	},
}

func init() {
	rootCmd.AddCommand(psCmd)
}
