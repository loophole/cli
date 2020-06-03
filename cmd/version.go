package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Loophole",
	Long:  `All software has versions. This is Loophole's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Loophole v0.0.1")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
