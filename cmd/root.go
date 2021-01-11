// +build !desktop

package cmd

import (
	"fmt"
	stdlog "log"
	"os"
	"time"

	"github.com/loophole/cli/config"
	lm "github.com/loophole/cli/internal/app/loophole/models"
	"github.com/loophole/cli/internal/pkg/cache"
	"github.com/loophole/cli/internal/pkg/closehandler"
	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

var displayOptions lm.DisplayOptions

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

	rootCmd.PersistentFlags().BoolVarP(&displayOptions.Verbose, "verbose", "v", false, "verbose output")
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
	if displayOptions.Verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	stdlog.SetFlags(0)
	stdlog.SetOutput(log.Logger)
}

// Execute runs command parsing chain
func Execute(config config.ApplicationConfig) {
	rootCmd.Version = fmt.Sprintf("%s (%s)", config.DisplayConfig.Version, config.DisplayConfig.CommitHash)
	displayOptions.Version = fmt.Sprintf("%s-%s", config.DisplayConfig.Version, config.DisplayConfig.CommitHash)

	closehandler.SetupCloseHandler(config.DisplayConfig.FeedbackFormURL)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
