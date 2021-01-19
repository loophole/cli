// +build !desktop

package cmd

import (
	"fmt"
	"os"

	"github.com/loophole/cli/internal/pkg/communication"
	"github.com/loophole/cli/internal/pkg/token"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to use your account",
	Long: `Loophole service requires authentication, this command allows you to log in or set up one
in case you don't yet have it.

Running this command as not logged in user will prompt you to open URL and use the browser to verify your identity.

Running this command as logged in user will fail, in cae you want to relogin then you need to log out first`,
	Run: func(cmd *cobra.Command, args []string) {
		if token.IsTokenSaved() {
			communication.LoginFailure(fmt.Errorf("Already logged in, please use `%s account logout` first to re-login", os.Args[0]))
		}

		deviceCodeSpec, err := token.RegisterDevice()
		if err != nil {
			communication.LoginFailure(fmt.Errorf("Error obtaining device code: %s", err.Error()))
		}
		communication.LoginStart(*deviceCodeSpec)
		quitChannel := make(chan bool)
		tokens, err := token.PollForToken(deviceCodeSpec.DeviceCode, deviceCodeSpec.Interval, quitChannel)
		if err != nil {
			communication.LoginFailure(fmt.Errorf("Error obtaining token: %s", err.Error()))
		}
		err = token.SaveToken(tokens)
		if err != nil {
			communication.LoginFailure(fmt.Errorf("Error saving token: %s", err.Error()))

		}
		communication.LoginSuccess(tokens.IDToken)
	},
}

func init() {
	accountCmd.AddCommand(loginCmd)
}
