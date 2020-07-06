package cmd

import (
	"github.com/loophole/cli/internal/pkg/token"
	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from your account",
	Long:  "Logout from your account",
	Run: func(cmd *cobra.Command, args []string) {
		if !token.IsTokenSaved() {
			logger.Fatal("Not logged in, nothing to do")
		}

		token.DeleteTokens()
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
