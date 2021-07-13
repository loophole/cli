// +build !desktop

package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/beevik/guid"
	"github.com/loophole/cli/config"
	lm "github.com/loophole/cli/internal/app/loophole/models"
	"github.com/loophole/cli/internal/pkg/cache"
	"github.com/loophole/cli/internal/pkg/communication"
	"github.com/loophole/cli/internal/pkg/inpututil"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/term"
)

var remoteEndpointSpecs lm.RemoteEndpointSpecs

var basicAuthUsernameFlagName = "basic-auth-username"
var basicAuthPasswordFlagName = "basic-auth-password"

func initServeCommand(serveCmd *cobra.Command) {
	sshDir := cache.GetLocalStorageDir(".ssh") // getting our sshDir and creating it, if it doesn't exist

	serveCmd.PersistentFlags().StringVarP(&remoteEndpointSpecs.IdentityFile, "identity-file", "i", fmt.Sprintf("%s/id_rsa", sshDir), "private key path")
	serveCmd.MarkFlagFilename("identity-file")

	serveCmd.PersistentFlags().StringVar(&remoteEndpointSpecs.SiteID, "hostname", "", "custom hostname you want to run service on")
	serveCmd.PersistentFlags().BoolVar(&config.Config.Display.QR, "qr", false, "use if you want a QR version of your url to be shown")

	serveCmd.PersistentFlags().StringVarP(&remoteEndpointSpecs.BasicAuthUsername, basicAuthUsernameFlagName, "u", "", "Basic authentication username to protect site with")
	serveCmd.PersistentFlags().StringVarP(&remoteEndpointSpecs.BasicAuthPassword, basicAuthPasswordFlagName, "p", "", "Basic authentication password to protect site with")

	remoteEndpointSpecs.TunnelID = guid.NewString()
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
		var password string
		if !inpututil.IsUsingPipe() { //only ask for password in terminal if not using pipe
			fmt.Print("Enter basic auth password: ")
			var err error
			passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
			password = string(passwordBytes)
			if err != nil {
				return err
			}
			fmt.Println()
		} else {
			reader := bufio.NewReader(os.Stdin)
			passwordBytes, err := reader.ReadBytes('\n')
			//if the reader encounters EOF before \n,
			//we assume that everything up until EOF is the intended password and continue
			if err != nil && err != io.EOF {
				communication.Warn("An error occured while reading the basic auth password from pipe.")
				communication.Fatal(err.Error())
			}
			password = strings.TrimSuffix(string(passwordBytes), "\n")

		}
		passwordFlag.Value.Set(password)
	}
	if passwordProvided && !usernameProvided {
		return fmt.Errorf("When using basic auth, both %s and %s have to be provided", basicAuthUsernameFlagName, basicAuthPasswordFlagName)
	}

	return nil
}
