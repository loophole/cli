package cmd

import (
	"os"
	"path"

	"github.com/loophole/cli/internal/pkg/cache"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// completionCmd represents the completion command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from your account",
	Long:  "Logout from your account",
	Run: func(cmd *cobra.Command, args []string) {
		if !isTokenSaved() {
			logger.Fatal("Not logged in, nothing to do")
		}

		deleteToken()
	},
}

func deleteToken() {
	storageDir := cache.GetLocalStorageDir()
	tokensLocation := path.Join(storageDir, "tokens.json")

	err := os.Remove(tokensLocation)
	if err != nil {
		logger.Fatal("There was a problem removing tokens file", zap.Error(err))
	}
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
