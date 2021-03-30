package cmd

import (
	"fmt"
	"os"

	"github.com/loophole/cli/internal/app/loopholed"
	"github.com/spf13/cobra"
)

// daemonStatusCommand represents the completion command
var daemonStatusCommand = &cobra.Command{
	Use:   "status",
	Short: "Displays loophole daemon status",
	Long:  "Displays loophole daemon status",
	Run: func(cmd *cobra.Command, args []string) {
		service := loopholed.New()
		status, err := service.Status()
		if err != nil {
			fmt.Println(status, "\nError: ", err)
			os.Exit(1)
		}
	},
}

func init() {
	daemonCmd.AddCommand(daemonStatusCommand)
}
