package cmd

import (
	"fmt"
	"os"

	"github.com/loophole/cli/internal/pkg/token"
	"github.com/rs/zerolog/log"
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
			log.Fatal().Msg(fmt.Sprintf("Already logged in, please use `%s account logout` first to re-login", os.Args[0]))
			os.Exit(1)
		}

		deviceCodeSpec, err := token.RegisterDevice()
		if err != nil {
			log.Fatal().Err(err).Msg("Error obtaining device code")
		}
		tokens, err := token.PollForToken(deviceCodeSpec.DeviceCode, deviceCodeSpec.Interval)
		if err != nil {
			log.Fatal().Err(err).Msg("Error obtaining token")
		}
		err = token.SaveToken(tokens)
		if err != nil {
			log.Fatal().Err(err).Msg("Error saving token")
		}
		log.Info().Msg("Logged in succesfully")
	},
}

func init() {
	accountCmd.AddCommand(loginCmd)
}
