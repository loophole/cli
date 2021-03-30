package cmd

import (
	"github.com/loophole/cli/internal/app/loopholed"
	"github.com/spf13/cobra"
)

// daemonRunCommand represents the completion command
var daemonRunCommand = &cobra.Command{
	Use:    "run",
	Short:  "Runs daemonized loophole actions",
	Long:   "Runs daemonized loophole actions, used only internally by daemon",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		service := loopholed.New()
		loopholeDaemon := &loopholed.LoopholeService{}

		service.Run(loopholeDaemon)
	},
}

func init() {
	daemonCmd.AddCommand(daemonRunCommand)
}
