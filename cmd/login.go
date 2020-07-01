package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/loophole/cli/internal/pkg/token"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	clientID       = "R569dcCOUErjw1xVZOzqc7OUCiGTYNqN"
	scope          = "openid"
	audience       = "https://api.loophole.cloud"
	timeoutSeconds = 180
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

		deviceCodeSpec, err := registerDevice()
		if err != nil {
			logger.Fatal("Error obtaining device code", zap.Error(err))
		}
		tokens, err := pollForToken(deviceCodeSpec.DeviceCode, deviceCodeSpec.Interval)
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

func registerDevice() (*token.DeviceCodeSpec, error) {
	deviceCodeURL := "https://owlsome.eu.auth0.com/oauth/device/code"

	payload := strings.NewReader(fmt.Sprintf("client_id=%s&scope=%s&audience=%s", url.QueryEscape(clientID), url.QueryEscape(scope), url.QueryEscape(audience)))

	req, err := http.NewRequest("POST", deviceCodeURL, payload)
	if err != nil {
		return nil, fmt.Errorf("There was a problem creating HTTP POST request for device code")
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("There was a problem executing request for device code")
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("There was a problem reading device token response body")
	}

	var jsonResponseBody token.DeviceCodeSpec
	err = json.Unmarshal(body, &jsonResponseBody)
	if err != nil {
		return nil, fmt.Errorf("There was a problem decoding device token response body")
	}

	fmt.Printf("Please open %s and use %s code to log in\n", aurora.Yellow(jsonResponseBody.VeritificationURI), aurora.Yellow(jsonResponseBody.UserCode))

	return &jsonResponseBody, nil
}

func pollForToken(deviceCode string, interval int) (*token.TokenSpec, error) {
	tokenURL := "https://owlsome.eu.auth0.com/oauth/token"
	grantType := "urn:ietf:params:oauth:grant-type:device_code"

	pollingInterval := time.Duration(interval) * time.Second
	logger.Debug("Polling with interval", zap.Duration("interval", pollingInterval), zap.String("unit", "second"))

	for {
		payload := strings.NewReader(fmt.Sprintf("grant_type=%s&device_code=%s&client_id=%s", url.QueryEscape(grantType), url.QueryEscape(deviceCode), url.QueryEscape(clientID)))

		req, err := http.NewRequest("POST", tokenURL, payload)
		if err != nil {
			logger.Debug("There was a problem creating HTTP POST request for token", zap.Error(err))
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		time.Sleep(pollingInterval)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			logger.Debug("There was a problem executing request for token", zap.Error(err))
			continue
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			logger.Debug("There was a problem reading token response body", zap.Error(err), zap.ByteString("body", body))
			continue
		}

		if res.StatusCode > 400 && res.StatusCode < 500 {
			var jsonResponseBody token.AuthError
			err := json.Unmarshal(body, &jsonResponseBody)
			if err != nil {
				logger.Debug("There was a problem decoding token response body", zap.Error(err), zap.ByteString("body", body))
				continue
			}
			logger.Debug("Error response", zap.String("error", jsonResponseBody.Error), zap.String("errorDescription", jsonResponseBody.ErrorDescription))
			if jsonResponseBody.Error == "authorization_pending" || jsonResponseBody.Error == "slow_down" {
				continue
			} else if jsonResponseBody.Error == "expired_token" || jsonResponseBody.Error == "invalid_grand" {
				return nil, fmt.Errorf("The device token expired, please reinitialize the login")
			} else if jsonResponseBody.Error == "access_denied" {
				return nil, fmt.Errorf("The device token got denied, please reinitialize the login")
			}
		} else if res.StatusCode >= 200 && res.StatusCode <= 300 {
			var jsonResponseBody token.TokenSpec
			err := json.Unmarshal(body, &jsonResponseBody)
			if err != nil {
				logger.Debug("There was a problem decoding token response body", zap.Error(err))
				continue
			}
			return &jsonResponseBody, nil
		} else {
			return nil, fmt.Errorf("Unexpected response from authorization server: %s", body)
		}
	}
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
