package cmd

import (
	"fmt"
	"os"

	"github.com/loophole/cli/internal/app/loopholed"
	"github.com/spf13/cobra"
)

// daemonRemoveCommand represents the completion command
var daemonRemoveCommand = &cobra.Command{
	Use:   "remove",
	Short: "Removes loophole daemon from your OS",
	Long:  "Removes loophole daemon from your OS",
	Run: func(cmd *cobra.Command, args []string) {
		service := loopholed.New()
		status, err := service.Remove()
		if err != nil {
			fmt.Println(status, "\nError: ", err)
			os.Exit(1)
		}
	},
}

func init() {
	daemonCmd.AddCommand(daemonRemoveCommand)
}
