package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/loophole/cli/internal/app/loophole"
	"github.com/mitchellh/go-homedir"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var secure bool
var identityFile string
var rootCmd = &cobra.Command{
	Use:   "loophole <port> [host]",
	Short: "Loophole exposes stuff over secure tunnels.",
	Long:  "Loophole exposes local servers to the public over secure tunnels.",
	Run: func(cmd *cobra.Command, args []string) {
		host := "127.0.0.1"
		if len(args) > 1 {
			host = args[1]
		}
		port, _ := strconv.ParseInt(args[0], 10, 32)
		loophole.Start(int(port), host, secure, identityFile)
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
	rootCmd.PersistentFlags().BoolVarP(&secure, "secure", "s", false, "Exposed service is TLS protected")
	rootCmd.PersistentFlags().StringVarP(&identityFile, "identity-file", "i", "", "Private key path")

	initConfig()
}

func initConfig() {
	if identityFile == "" {
		home, err := homedir.Dir()
		if err != nil {
			panic(err)
		}

		identityFile = fmt.Sprintf("%s/.ssh/id_rsa", home)
		log.Printf("No identity key provided, using default: %s", identityFile)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// Execute runs command parsing chain
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
