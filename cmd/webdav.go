package cmd

import (
	"errors"

	"github.com/loophole/cli/internal/app/loophole"
	lm "github.com/loophole/cli/internal/app/loophole/models"

	"github.com/spf13/cobra"
)

var webdavEndpointSpecs lm.LocalDirectorySpecs

var webdavCmd = &cobra.Command{
	Use:   "webdav <path>",
	Short: "Expose given directory to the public via WebDav",
	Long: `Exposes local directory to the public via WebDav via loophole tunnel.

This can then be even mounted on other machines in the Windows Explorer, macOS Finder, Linux Gnome Files or Linux KDE Konqueror etc.

To expose local directory via webdav (e.g. /data/my-data) simply use 'loophole webdav /data/my-data'.`,
	Run: func(cmd *cobra.Command, args []string) {
		webdavEndpointSpecs.Path = args[0]
		loophole.ForwardDirectoryViaWebdav(lm.ExposeWebdavConfig{
			Local:   webdavEndpointSpecs,
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
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return parseBasicAuthFlags(cmd.Flags())
	},
}

func init() {
	initServeCommand(webdavCmd)

	rootCmd.AddCommand(webdavCmd)
}
