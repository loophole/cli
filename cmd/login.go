package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/loophole/cli/internal/pkg/cache"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	clientID       = "R569dcCOUErjw1xVZOzqc7OUCiGTYNqN"
	scope          = "openid"
	audience       = "https://owlsome.eu.auth0.com/api/v2/"
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

type deviceCodeSpec struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
	VeritificationURI       string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
}

type authError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type tokenSpec struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// completionCmd represents the completion command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to use your account",
	Long:  "Log in to use your account",
	Run: func(cmd *cobra.Command, args []string) {
		if isTokenSaved() {
			logger.Fatal("Already logged in, please logout first")
		}

		deviceCodeSpec, err := registerDevice()
		if err != nil {
			logger.Fatal("Error obtaining device code", zap.Error(err))
		}
		token, err := pollForToken(deviceCodeSpec.DeviceCode, deviceCodeSpec.Interval)
		if err != nil {
			logger.Fatal("Error obtaining token", zap.Error(err))
		}
		err = saveToken(token)
		if err != nil {
			logger.Fatal("Error saving token", zap.Error(err))
		}
		fmt.Println("Logged in succesfully")
	},
}

func registerDevice() (*deviceCodeSpec, error) {
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

	var jsonResponseBody deviceCodeSpec
	err = json.Unmarshal(body, &jsonResponseBody)
	if err != nil {
		return nil, fmt.Errorf("There was a problem decoding device token response body")
	}

	fmt.Printf("Please open %s and use %s code to log in\n", aurora.Yellow(jsonResponseBody.VeritificationURI), aurora.Yellow(jsonResponseBody.UserCode))

	return &jsonResponseBody, nil
}

func pollForToken(deviceCode string, interval int) (*tokenSpec, error) {
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
			var jsonResponseBody authError
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
			var jsonResponseBody tokenSpec
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

func saveToken(token *tokenSpec) error {
	storageDir := cache.GetLocalStorageDir()
	tokensLocation := path.Join(storageDir, "tokens.json")

	tokenBytes, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("There was a problem encoding token: %v", err)
	}
	err = ioutil.WriteFile(tokensLocation, tokenBytes, 0644)
	if err != nil {
		return fmt.Errorf("There was a problem writing tokens file: %v", err)
	}
	return nil
}

func isTokenSaved() bool {
	storageDir := cache.GetLocalStorageDir()
	tokensLocation := path.Join(storageDir, "tokens.json")
	if _, err := os.Stat(tokensLocation); os.IsNotExist(err) {
		return false
	} else if err != nil {
		logger.Fatal("There was a problem reading token file", zap.Error(err))
	}
	return true
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
