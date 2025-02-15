//go:build !desktop
// +build !desktop

package cmd

import (
	"context"
	"fmt"
	stdlog "log"
	"os"
	"time"

	"github.com/loophole/cli/config"
	"github.com/loophole/cli/internal/pkg/cache"
	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "loophole",
	Short: "Loophole - End to end TLS encrypted TCP communication between you and your clients",
	Long:  "Loophole - End to end TLS encrypted TCP communication between you and your clients",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cobra.OnInitialize(initLogger)

	rootCmd.PersistentFlags().BoolVarP(&config.Config.Display.Verbose, "verbose", "v", false, "verbose output")
}

func initLogger() {
	logLocation := cache.GetLocalStorageFile(fmt.Sprintf("%s.log", time.Now().Format("2006-01-02--15-04-05")), "logs")

	f, err := os.Create(logLocation)
	if err != nil {
		stdlog.Fatalln("Error creating log file:", err)
	}

	consoleWriter := zerolog.ConsoleWriter{Out: colorable.NewColorableStderr()}
	multi := zerolog.MultiLevelWriter(consoleWriter, f)
	log.Logger = zerolog.New(multi).With().Timestamp().Logger()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if config.Config.Display.Verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	stdlog.SetFlags(0)
	stdlog.SetOutput(log.Logger)
}

// Execute runs command parsing chain
func Execute() {
	rootCmd.Version = fmt.Sprintf("%s (%s)", config.Config.Version, config.Config.CommitHash)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// GoExecute runs command parsing chain from go.
// Look example ../example/main.go.
func GoExecute(ctx context.Context, version, commit, mode string, args ...string) error {
	config.Config.Version = version
	config.Config.CommitHash = commit
	config.Config.ClientMode = mode

	rootCmd.Version = fmt.Sprintf("%s (%s)", config.Config.Version, config.Config.CommitHash)

	rootCmd.SetArgs(args)
	return rootCmd.ExecuteContext(ctx)
}
