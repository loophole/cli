package config

import (
	"fmt"
	"time"

	"github.com/loophole/cli/internal/app/loophole/models"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// OAuthConfig defined OAuth settings shape
type OAuthConfig struct {
	DeviceCodeURL string `json:"deviceCodeUrl"`
	TokenURL      string `json:"tokenUrl"`
	ClientID      string `json:"clientId"`
	Scope         string `json:"scope"`
	Audience      string `json:"audience"`
}

// DisplayConfig defines the display switches shape
type DisplayConfig struct {
	Verbose bool `json:"verbose"`
	QR      bool `json:"qr"`
}

// ApplicationConfig defines the application config shape
type ApplicationConfig struct {
	Version    string `json:"version"`
	CommitHash string `json:"commitHash"`
	ClientMode string `json:"clientMode"`

	FeedbackFormURL string `json:"feedbackFormUrl"`

	OAuth   OAuthConfig   `json:"oauthConfig"`
	Display DisplayConfig `json:"displayConfig"`

	APIEndpoint     models.Endpoint `json:"apiConfig"`
	GatewayEndpoint models.Endpoint `json:"gatewayConfig"`
}

func SetupViperConfig() error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	viper.SetDefault("lastreminder", time.Time{})         //date of last reminder, default is zero value for time
	viper.SetDefault("availableversion", "1.0.0-beta.14") //last seen latest version
	viper.SetDefault("remindercount", 3)                  //counts to zero, then switches from prompt to notification reminder
	viper.SetDefault("savehostnames", true)
	viper.SetDefault("usedhostnames", []string{})
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(fmt.Sprintf("%s/.loophole/", home))
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok { //create a config if none exist yet
			err = SaveViperConfig()
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func SaveViperConfig() error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	err = viper.WriteConfigAs(fmt.Sprintf("%s/.loophole/config.json", home))
	if err != nil {
		return err
	}
	return nil
}
