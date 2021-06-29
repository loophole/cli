// +build !desktop

package cmd

import (
	"fmt"

	"github.com/loophole/cli/internal/pkg/communication"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCommand represents the command that lists previously used hostnames
var listCommand = &cobra.Command{
	Use:   "list",
	Short: "Show used hostnames",
	Long:  `Show previously used hostnames that were successfully used.`,
	Run: func(cmd *cobra.Command, args []string) {
		hostnames := viper.GetStringSlice("usedhostnames")
		communication.Info(fmt.Sprintf("The following %d hostnames have been used:", len(hostnames)))
		for _, hostname := range hostnames {
			communication.Info(hostname)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCommand)
}
