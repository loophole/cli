// +build !desktop

package cmd

import (
	"github.com/loophole/cli/internal/pkg/communication"
	"github.com/loophole/cli/internal/pkg/token"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out from your account",
	Long: `This command deletes all the locally stored tokens which allows you to re-login or simply stay logged out.

In regular scenario you should not need to use it, as tokens are getting refreshed automatically.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !token.IsTokenSaved() {
			communication.LogFatalMsg("Not logged in, nothing to do")
		}

		err := token.DeleteTokens()
		if err != nil {
			communication.LogFatalErr("There as a problem logging out", err)
		}
		communication.LogInfo("Logged out succesfully")
	},
}

func init() {
	accountCmd.AddCommand(logoutCmd)
}
