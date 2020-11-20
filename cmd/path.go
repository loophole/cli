package cmd

import (
	"errors"

	"github.com/loophole/cli/internal/app/loophole"
	lm "github.com/loophole/cli/internal/app/loophole/models"
	"github.com/spf13/cobra"
)

var dirEndpointSpecs lm.LocalDirectorySpecs

var dirCmd = &cobra.Command{
	Use:   "dir <path>",
	Short: "Expose given directory to the public",
	Long:  "Expose directory to the public",
	Run: func(cmd *cobra.Command, args []string) {
		dirEndpointSpecs.Path = args[0]
		loophole.ForwardDirectory(lm.ExposeDirectoryConfig{
			Local:   dirEndpointSpecs,
			Remote:  remoteEndpointSpecs,
			Display: displayOptions,
		})
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Missing argument: path")
		}
		return nil
	},
}

func init() {
	initServeCommand(dirCmd)
	rootCmd.AddCommand(dirCmd)
}
