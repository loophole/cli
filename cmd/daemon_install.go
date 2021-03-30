package cmd

import (
	"fmt"
	"os"

	"github.com/loophole/cli/internal/app/loopholed"
	"github.com/spf13/cobra"
)

// daemonInstallCommand represents the completion command
var daemonInstallCommand = &cobra.Command{
	Use:   "install",
	Short: "Installs daemon in your OS",
	Long:  "Installs daemon in your OS",
	Run: func(cmd *cobra.Command, args []string) {
		service := loopholed.New()
		fmt.Println("Installing loophole daemon...")
		status, err := service.Install("daemon", "run")
		if err != nil {
			fmt.Println(status, "\nError: ", err)
			os.Exit(1)
		}
		fmt.Println("Daemon succesfully installed")
		fmt.Println("Attempting to start...")
		status, err = service.Start()
		if err != nil {
			fmt.Println(status, "\nError: ", err)
			os.Exit(1)
		}
		fmt.Println("Daemon succesfully started")
	},
}

func init() {
	daemonCmd.AddCommand(daemonInstallCommand)
}
