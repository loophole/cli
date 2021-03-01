// +build !desktop

package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/loophole/cli/internal/app/loophole"
	lm "github.com/loophole/cli/internal/app/loophole/models"
	"github.com/loophole/cli/internal/pkg/communication"
	"github.com/loophole/cli/internal/pkg/token"

	"github.com/spf13/cobra"
)

var localEndpointSpecs lm.LocalHTTPEndpointSpecs

var httpCmd = &cobra.Command{
	Use:   "http <port> [host]",
	Short: "Expose http server on given port to the public",
	Long: `Exposes http server running locally, or on locally available machine to the public via loophole tunnel.

To expose server running locally on port 3000 simply use 'loophole http 3000'.
To expose port running on some local host e.g. 192.168.1.20 use 'loophole http <port> 192.168.1.20'`,
	Run: func(cmd *cobra.Command, args []string) {
		loggedIn := token.IsTokenSaved()
		idToken := token.GetIdToken()
		communication.ApplicationStart(loggedIn, idToken)

		checkVersion()

		localEndpointSpecs.Host = "127.0.0.1"
		if len(args) > 1 {
			localEndpointSpecs.Host = args[1]
		}
		port, _ := strconv.ParseInt(args[0], 10, 32)
		localEndpointSpecs.Port = int32(port)
		quitChannel := make(chan bool)

		exposeConfig := lm.ExposeHTTPConfig{
			Local:  localEndpointSpecs,
			Remote: remoteEndpointSpecs,
		}

		authMethod, err := loophole.RegisterTunnel(&exposeConfig.Remote)
		if err != nil {
			communication.Fatal(err.Error())
		}

		loophole.ForwardPort(exposeConfig, authMethod, quitChannel)
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
	httpCmd.Flags().BoolVar(&localEndpointSpecs.HTTPS, "https", false, "use if your server is already using HTTPS")
	httpCmd.Flags().BoolVar(&remoteEndpointSpecs.DisableProxyErrorPage, "disable-proxy-error-page", false, "disable proxy error page and return 502 when your server is not available")
	httpCmd.Flags().StringVar(&localEndpointSpecs.Path, "path", "", "specify path you wish to expose")

	rootCmd.AddCommand(httpCmd)
}
