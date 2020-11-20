package cmd

import (
	"fmt"

	lm "github.com/loophole/cli/internal/app/loophole/models"
	"github.com/mitchellh/go-homedir"

	"github.com/spf13/cobra"
)

var remoteEndpointSpecs lm.RemoteEndpointSpecs

// accountCmd represents the account command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Parent for commands concerning tunnels",
	Long:  "Parent for commands concerning tunnels. Always use with one of the following: port, dir",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	remoteEndpointSpecs.GatewayEndpoint.Host = "gateway.loophole.host"
	remoteEndpointSpecs.GatewayEndpoint.Port = 8022
	remoteEndpointSpecs.APIEndpoint.Protocol = "https"
	remoteEndpointSpecs.APIEndpoint.Host = "api.loophole.cloud"
	remoteEndpointSpecs.APIEndpoint.Port = 443

	serveCmd.PersistentFlags().StringVarP(&remoteEndpointSpecs.IdentityFile, "identity-file", "i", fmt.Sprintf("%s/.ssh/id_rsa", home), "private key path")
	serveCmd.PersistentFlags().StringVar(&remoteEndpointSpecs.SiteID, "hostname", "", "custom hostname you want to run service on")
	serveCmd.PersistentFlags().BoolVar(&displayOptions.QR, "qr", false, "use if you want a QR version of your url to be shown")

	rootCmd.AddCommand(serveCmd)
}
