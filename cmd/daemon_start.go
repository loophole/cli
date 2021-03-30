package cmd

import (
	"fmt"
	"os"

	"github.com/loophole/cli/internal/app/loopholed"
	"github.com/spf13/cobra"
)

// daemonStartCommand represents the completion command
var daemonStartCommand = &cobra.Command{
	Use:   "start",
	Short: "Starts stopped loophole daemon",
	Long:  "Starts stopped loophole daemon",
	Run: func(cmd *cobra.Command, args []string) {
		service := loopholed.New()
		status, err := service.Start()
		if err != nil {
			fmt.Println(status, "\nError: ", err)
			os.Exit(1)
		}
	},
}

func init() {
	daemonCmd.AddCommand(daemonStartCommand)
}
