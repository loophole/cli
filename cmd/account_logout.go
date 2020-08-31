package cmd

import (
	"github.com/loophole/cli/internal/pkg/token"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from your account",
	Long:  "Logout from your account",
	Run: func(cmd *cobra.Command, args []string) {
		if !token.IsTokenSaved() {
			log.Fatal().Msg("Not logged in, nothing to do")
		}

		err := token.DeleteTokens()
		if err != nil {
			log.Fatal().Err(err).Msg("There as a problem logging out")
		}
		log.Info().Msg("Logged out succesfully")
	},
}

func init() {
	accountCmd.AddCommand(logoutCmd)
}
