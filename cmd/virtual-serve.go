package cmd

import (
	"fmt"

	lm "github.com/loophole/cli/internal/app/loophole/models"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var remoteEndpointSpecs lm.RemoteEndpointSpecs

func initServeCommand(serveCmd *cobra.Command) {
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
}
