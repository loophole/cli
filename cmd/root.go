package cmd

import (
	"errors"
	"fmt"
	stdlog "log"
	"os"
	"strconv"

	"github.com/loophole/cli/internal/app/loophole"
	lm "github.com/loophole/cli/internal/app/loophole/models"
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
	rootCmd.Version = "1.0.0"

	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	cobra.OnInitialize(initLogger)

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	rootCmd.Flags().StringVarP(&config.IdentityFile, "identity-file", "i", fmt.Sprintf("%s/.ssh/id_rsa", home), "private key path")
	rootCmd.Flags().StringVar(&config.GatewayEndpoint.Host, "gateway-url", "gateway.loophole.host", "remote gateway URL")
	rootCmd.Flags().Int32Var(&config.GatewayEndpoint.Port, "gateway-port", 8022, "remote gateway port")
	rootCmd.Flags().StringVar(&config.SiteID, "hostname", "", "custom hostname you want to run service on")
	rootCmd.Flags().BoolVar(&config.HTTPS, "https", false, "use if your local service is already using HTTPS")

}

func initLogger() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	stdlog.SetFlags(0)
	stdlog.SetOutput(log.Logger)
}

// Execute runs command parsing chain
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
