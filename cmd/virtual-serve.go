package cmd

import (
	"fmt"

	lm "github.com/loophole/cli/internal/app/loophole/models"
	"github.com/loophole/cli/internal/pkg/cache"
	"github.com/spf13/cobra"
)

var remoteEndpointSpecs lm.RemoteEndpointSpecs

func initServeCommand(serveCmd *cobra.Command) {
	sshDir := cache.GetLocalStorageDir(".ssh") //getting our sshDir and creating it, if it doesn't exist

	remoteEndpointSpecs.GatewayEndpoint.Host = "gateway.loophole.host"
	remoteEndpointSpecs.GatewayEndpoint.Port = 8022
	remoteEndpointSpecs.APIEndpoint.Protocol = "https"
	remoteEndpointSpecs.APIEndpoint.Host = "api.loophole.cloud"
	remoteEndpointSpecs.APIEndpoint.Port = 443

	serveCmd.PersistentFlags().StringVarP(&remoteEndpointSpecs.IdentityFile, "identity-file", "i", fmt.Sprintf("%s/id_rsa", sshDir), "private key path")
	serveCmd.PersistentFlags().StringVar(&remoteEndpointSpecs.SiteID, "hostname", "", "custom hostname you want to run service on")
	serveCmd.PersistentFlags().BoolVar(&displayOptions.QR, "qr", false, "use if you want a QR version of your url to be shown")
}
