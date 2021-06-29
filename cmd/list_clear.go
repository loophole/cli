// +build !desktop

package cmd

import (
	"github.com/loophole/cli/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listclearCommand represents the command that clears the list of previously used hostnames
var listclearCommand = &cobra.Command{
	Use:   "clear",
	Short: "Delete the current list of saved hostnames",
	Long:  `Delete the current list of saved hostnames`,
	Run: func(cmd *cobra.Command, args []string) {
		viper.Set("usedhostnames", []string{})
		config.SaveViperConfig()
	},
}

func init() {
	listCommand.AddCommand(listclearCommand)
}
