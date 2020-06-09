package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/loophole/cli/internal/app/loophole"
	lm "github.com/loophole/cli/internal/app/loophole/models"
	"github.com/mitchellh/go-homedir"

	"github.com/spf13/cobra"
)

var config lm.Config

var rootCmd = &cobra.Command{
	Use:   "loophole <port> [host]",
	Short: "Loophole exposes stuff over secure tunnels.",
	Long:  "Loophole exposes local servers to the public over secure tunnels.",
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
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	rootCmd.PersistentFlags().StringVarP(&config.IdentityFile, "identity-file", "i", fmt.Sprintf("%s/.ssh/id_rsa", home), "Private key path")
	rootCmd.PersistentFlags().StringVar(&config.GatewayEndpoint.Host, "gateway-url", "loophole.host", "Remote gateway URL")
	rootCmd.PersistentFlags().Int32Var(&config.GatewayEndpoint.Port, "gateway-port", 8022, "Remote gateway port")
	rootCmd.PersistentFlags().StringVar(&config.APIURL, "api-url", "https://api.loophole.cloud", "Remote gateway URL")
}

// Execute runs command parsing chain
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
