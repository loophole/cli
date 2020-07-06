package cmd

import (
	"fmt"
	"os"

	"github.com/loophole/cli/internal/pkg/token"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func init() {
	atomicLevel := zap.NewAtomicLevel()
	encoderCfg := zap.NewProductionEncoderConfig()
	logger = zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atomicLevel,
	))

	atomicLevel.SetLevel(zap.DebugLevel)
}

// completionCmd represents the completion command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to use your account",
	Long:  "Log in to use your account",
	Run: func(cmd *cobra.Command, args []string) {
		if token.IsTokenSaved() {
			logger.Fatal("Already logged in, please logout first")
		}

		deviceCodeSpec, err := token.RegisterDevice()
		if err != nil {
			logger.Fatal("Error obtaining device code", zap.Error(err))
		}
		tokens, err := token.PollForToken(deviceCodeSpec.DeviceCode, deviceCodeSpec.Interval)
		if err != nil {
			logger.Fatal("Error obtaining token", zap.Error(err))
		}
		err = token.SaveToken(tokens)
		if err != nil {
			logger.Fatal("Error saving token", zap.Error(err))
		}
		fmt.Println("Logged in succesfully")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
