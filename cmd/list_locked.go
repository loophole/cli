// +build !desktop

package cmd

import (
	"fmt"

	"github.com/loophole/cli/internal/pkg/apiclient"
	"github.com/loophole/cli/internal/pkg/communication"
	"github.com/spf13/cobra"
)

// listlockedCommand represents the command that lists the hostnames that are currently locked for the user
var listlockedCommand = &cobra.Command{
	Use:   "locked",
	Short: "List the hostnames that are currently locked for you.",
	Long:  `List the hostnames that are currently locked for you.`,
	Run: func(cmd *cobra.Command, args []string) {
		hostnames, err := apiclient.GetLockedHostnames()
		if err != nil {
			communication.Error(fmt.Sprintf("Error while trying to retrieve locked hostnames: %s", err.Error()))
			return
		}
		communication.Info(fmt.Sprintf("The following %d hostnames are currently locked for you:", len(hostnames)))
		for _, hostname := range hostnames {
			communication.Info(hostname)
		}
	},
}

func init() {
	listCommand.AddCommand(listlockedCommand)
}
