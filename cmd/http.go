package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/loophole/cli/internal/app/loophole"
	lm "github.com/loophole/cli/internal/app/loophole/models"

	"github.com/spf13/cobra"
)

var localEndpointSpecs lm.LocalHttpEndpointSpecs

var httpCmd = &cobra.Command{
	Use:   "http <port> [host]",
	Short: "Expose http server on given port to the public",
	Long: `Exposes http server running locally, or on locally available machine to the public via loophole tunnel.

To expose server running locally on port 3000 simply use 'loophole http 3000'.
To expose port running on some local host e.g. 192.168.1.20 use 'loophole http <port> 192.168.1.20'`,
	Run: func(cmd *cobra.Command, args []string) {
		localEndpointSpecs.Host = "127.0.0.1"
		if len(args) > 1 {
			localEndpointSpecs.Host = args[1]
		}
		port, _ := strconv.ParseInt(args[0], 10, 32)
		localEndpointSpecs.Port = int32(port)
		loophole.ForwardPort(lm.ExposeHttpConfig{
			Local:   localEndpointSpecs,
			Remote:  remoteEndpointSpecs,
			Display: displayOptions,
		})
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Missing argument: port")
		}
		_, err := strconv.ParseInt(args[0], 10, 32)
		if err != nil {
			return fmt.Errorf("Invalid argument: port: %v", err)
		}
		return nil
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return parseBasicAuthFlags(cmd.Flags())
	},
}

func init() {
	initServeCommand(httpCmd)
	localEndpointSpecs.HTTPS = false

	rootCmd.AddCommand(httpCmd)
}
