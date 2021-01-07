// +build !desktop

package cmd

import (
	"errors"

	"github.com/loophole/cli/internal/app/loophole"
	lm "github.com/loophole/cli/internal/app/loophole/models"
	"github.com/loophole/cli/internal/pkg/communication"
	"github.com/loophole/cli/internal/pkg/token"
	"github.com/spf13/cobra"
)

var dirEndpointSpecs lm.LocalDirectorySpecs

var dirCmd = &cobra.Command{
	Use:     "path <path>",
	Aliases: []string{"dir", "directory"},
	Short:   "Expose given directory to the public",
	Long: `Exposes local directory to the public via loophole tunnel.

To expose local directory (e.g. /data/my-data) simply use 'loophole path /data/my-data'.`,
	Run: func(cmd *cobra.Command, args []string) {
		loggedIn := token.IsTokenSaved()
		communication.PrintWelcomeMessage(loggedIn)
		dirEndpointSpecs.Path = args[0]
		quitChannel := make(chan bool)
		loophole.ForwardDirectory(lm.ExposeDirectoryConfig{
			Local:   dirEndpointSpecs,
			Remote:  remoteEndpointSpecs,
			Display: displayOptions,
		}, quitChannel)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Missing argument: path")
		}
		return nil
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return parseBasicAuthFlags(cmd.Flags())
	},
}

func init() {
	initServeCommand(dirCmd)
	rootCmd.AddCommand(dirCmd)
}
