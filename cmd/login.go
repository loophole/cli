package cmd

import (
	"os"

	"github.com/loophole/cli/internal/pkg/token"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to use your account",
	Long:  "Log in to use your account",
	Run: func(cmd *cobra.Command, args []string) {
		if token.IsTokenSaved() {
			log.Fatal().Msg("Already logged in, please logout first to reinitialize login")
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
