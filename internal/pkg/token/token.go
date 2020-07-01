package token

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/loophole/cli/internal/pkg/cache"
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

type DeviceCodeSpec struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
	VeritificationURI       string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
}

type AuthError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type TokenSpec struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

func IsTokenSaved() bool {
	tokensLocation := cache.GetLocalStorageFile("tokens.json")

	if _, err := os.Stat(tokensLocation); os.IsNotExist(err) {
		return false
	} else if err != nil {
		logger.Fatal("There was a problem reading tokens file", zap.Error(err))
	}
	return true
}

func SaveToken(token *TokenSpec) error {
	tokensLocation := cache.GetLocalStorageFile("tokens.json")

	tokenBytes, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("There was a problem encoding tokens: %v", err)
	}
	err = ioutil.WriteFile(tokensLocation, tokenBytes, 0644)
	if err != nil {
		return fmt.Errorf("There was a problem writing tokens file: %v", err)
	}
	return nil
}

func DeleteTokens() {
	tokensLocation := cache.GetLocalStorageFile("tokens.json")

	err := os.Remove(tokensLocation)
	if err != nil {
		logger.Fatal("There was a problem removing tokens file", zap.Error(err))
	}
}

func GetAccessToken() (string, error) {
	tokensLocation := cache.GetLocalStorageFile("tokens.json")

	tokens, err := ioutil.ReadFile(tokensLocation)
	if err != nil {
		return "", fmt.Errorf("There was a problem reading tokens: %v", err)
	}
	var token TokenSpec
	err = json.Unmarshal(tokens, &token)
	if err != nil {
		return "", fmt.Errorf("There was a problem decoding tokens: %v", err)
	}
	return token.AccessToken, nil
}
