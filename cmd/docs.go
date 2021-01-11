// +build !desktop

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// docsCommand represents the completion command
var docsCommand = &cobra.Command{
	Use:   "docs",
	Short: "Generates docs",
	Long: `Generates docs for existing loophole commands and saves them in './docs'.

Mainly for developers.`,
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		docsPath := "./docs"
		if _, err := os.Stat(docsPath); os.IsNotExist(err) {
			os.Mkdir(docsPath, os.ModePerm)
		}
		err := doc.GenMarkdownTree(rootCmd, docsPath)
		if err != nil {
			fmt.Printf("Failed to generate docs: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Docs succesfully generated in %s\n", docsPath)
	},
}

func init() {
	rootCmd.AddCommand(docsCommand)
}
