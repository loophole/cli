// +build !desktop

package cmd

import (
	"github.com/loophole/cli/config"
	"github.com/loophole/cli/internal/pkg/communication"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listtoggleCommand represents the command that toggles whether we save previously used hostnames
var listtoggleCommand = &cobra.Command{
	Use:   "toggle",
	Short: "Stop or start saving used hostnames.",
	Long: `Stop or start saving used hostnames.
By default, a list of hostnames you successfully used is saved locally for your convenience.
With this command, you can stop it or resume it if you stopped it before.
This function is not related to the timed reservation of hostnames.`,
	Run: func(cmd *cobra.Command, args []string) {
		oldState := viper.GetBool("savehostnames")
		viper.Set("savehostnames", !oldState)
		config.SaveViperConfig()
		if oldState {
			communication.Info("Hostname saving is now turned off.")
		} else {
			communication.Info("Hostname saving is now turned on.")
		}
	},
}

func init() {
	listCommand.AddCommand(listtoggleCommand)
}
