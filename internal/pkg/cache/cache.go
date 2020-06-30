package cache

import (
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
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

// GetLocalStorageDir returns local directory for loophole cache purposes
func GetLocalStorageDir() string {
	home, err := homedir.Dir()
	if err != nil {
		logger.Fatal("Error reading user home directory ", zap.Error(err))
	}

	return path.Join(home, ".local", "loophole")
}
