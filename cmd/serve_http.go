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

var portCmd = &cobra.Command{
	Use:   "port <port> [host]",
	Short: "Expose given port to the public",
	Long:  "Expose host:port to the public",
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
}

func init() {
	localEndpointSpecs.HTTPS = false

	serveCmd.AddCommand(portCmd)
}
