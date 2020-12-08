package cmd

import (
	"fmt"
	"os"

	lm "github.com/loophole/cli/internal/app/loophole/models"
	"github.com/loophole/cli/internal/pkg/cache"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/crypto/ssh/terminal"
)

var remoteEndpointSpecs lm.RemoteEndpointSpecs

var basicAuthUsernameFlagName = "basic-auth-username"
var basicAuthPasswordFlagName = "basic-auth-password"

func initServeCommand(serveCmd *cobra.Command) {
	sshDir := cache.GetLocalStorageDir(".ssh") //getting our sshDir and creating it, if it doesn't exist

	remoteEndpointSpecs.GatewayEndpoint.Host = "gateway.loophole.host"
	remoteEndpointSpecs.GatewayEndpoint.Port = 8022
	remoteEndpointSpecs.APIEndpoint.Protocol = "https"
	remoteEndpointSpecs.APIEndpoint.Host = "api.loophole.cloud"
	remoteEndpointSpecs.APIEndpoint.Port = 443

	serveCmd.PersistentFlags().StringVarP(&remoteEndpointSpecs.IdentityFile, "identity-file", "i", fmt.Sprintf("%s/id_rsa", sshDir), "private key path")
	serveCmd.MarkFlagFilename("identity-file")
	serveCmd.PersistentFlags().StringVar(&remoteEndpointSpecs.SiteID, "hostname", "", "custom hostname you want to run service on")
	serveCmd.PersistentFlags().BoolVar(&displayOptions.QR, "qr", false, "use if you want a QR version of your url to be shown")

	serveCmd.PersistentFlags().StringVarP(&remoteEndpointSpecs.BasicAuthUsername, basicAuthUsernameFlagName, "u", "", "Basic authentication username to protect site with")
	serveCmd.PersistentFlags().StringVarP(&remoteEndpointSpecs.BasicAuthPassword, basicAuthPasswordFlagName, "p", "", "Basic authentication password to protect site with")
}

func parseBasicAuthFlags(flagset *pflag.FlagSet) error {
	usernameProvided := false
	passwordProvided := false
	var passwordFlag *pflag.Flag

	flagset.VisitAll(func(flag *pflag.Flag) {
		if flag.Name == basicAuthUsernameFlagName && flag.Value.String() != "" {
			usernameProvided = true
		}
		if flag.Name == basicAuthPasswordFlagName {
			passwordFlag = flag
			if flag.Value.String() != "" {
				passwordProvided = true
			}
		}
	})

	if usernameProvided && !passwordProvided {
		fmt.Print("Enter basic auth password:")

		password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return err
		}
		fmt.Println()
		passwordFlag.Value.Set(string(password))
	}
	if passwordProvided && !usernameProvided {
		return fmt.Errorf("When using basic auth, both %s and %s have to be provided", basicAuthUsernameFlagName, basicAuthPasswordFlagName)
	}

	return nil
}
