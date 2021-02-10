package cmd

import (
	"fmt"
	"os"

	"github.com/loophole/cli/internal/app/loopholed"
	"github.com/spf13/cobra"
)

// daemonStopCommand represents the completion command
var daemonStopCommand = &cobra.Command{
	Use:   "stop",
	Short: "Stops running loophole daemon",
	Long:  "Stops running loophole daemon",
	Run: func(cmd *cobra.Command, args []string) {
		service := loopholed.New()
		status, err := service.Stop()
		if err != nil {
			fmt.Println(status, "\nError: ", err)
			os.Exit(1)
		}
	},
}

func init() {
	daemonCmd.AddCommand(daemonStopCommand)
}
