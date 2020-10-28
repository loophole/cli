package cmd

import (
	"errors"
	"fmt"
	stdlog "log"
	"os"
	"strconv"
	"time"

	"github.com/loophole/cli/internal/app/loophole"
	lm "github.com/loophole/cli/internal/app/loophole/models"
	"github.com/loophole/cli/internal/app/metrics"
	"github.com/loophole/cli/internal/pkg/cache"
	"github.com/mattn/go-colorable"
	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

var config lm.Config
var verbose bool

var rootCmd = &cobra.Command{
	Use:   "loophole <port> [host]",
	Short: "Loophole - End to end TLS encrypted TCP communication between you and your clients",
	Long:  "Loophole - End to end TLS encrypted TCP communication between you and your clients",
	Run: func(cmd *cobra.Command, args []string) {
		config.Host = "127.0.0.1"
		if len(args) > 1 {
			config.Host = args[1]
		}
		port, _ := strconv.ParseInt(args[0], 10, 32)
		config.Port = int32(port)
		go metrics.Serve()
		loophole.Start(config)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Missing argument: port")
		}
		_, err := strconv.ParseInt(args[0], 10, 32)
		if err != nil {
			return fmt.Errorf("Invalid argument: port: %v", err)
		}
		return nil
	},
}

func init() {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	cobra.OnInitialize(initLogger)

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	rootCmd.Flags().StringVarP(&config.IdentityFile, "identity-file", "i", fmt.Sprintf("%s/.ssh/id_rsa", home), "private key path")
	config.GatewayEndpoint.Host = "gateway.loophole.host"
	config.GatewayEndpoint.Port = 8022
	rootCmd.Flags().StringVar(&config.SiteID, "hostname", "", "custom hostname you want to run service on")
	config.HTTPS = false
	rootCmd.Flags().BoolVar(&config.QR, "qr", false, "use if you want a QR version of your url to be shown")

}

func initLogger() {
	logLocation := cache.GetLocalStorageFile(fmt.Sprintf("%s, %s", time.Now().Format("2006-01-02--15-04-05"), ".log"), "logs")

	f, err := os.Create(logLocation)
	if err != nil {
		stdlog.Fatalln("Error creating log file:", err)
	}

	consoleWriter := zerolog.ConsoleWriter{Out: colorable.NewColorableStderr()}
	multi := zerolog.MultiLevelWriter(consoleWriter, f)
	log.Logger = zerolog.New(multi).With().Timestamp().Str("component", "loophole").Logger()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	stdlog.SetFlags(0)
	stdlog.SetOutput(log.Logger)
}

// Execute runs command parsing chain
func Execute(version string, commit string) {
	rootCmd.Version = fmt.Sprintf("%s (%s)", version, commit)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
