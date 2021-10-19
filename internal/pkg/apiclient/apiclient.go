package apiclient

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"time"

	"github.com/loophole/cli/config"
	"github.com/loophole/cli/internal/pkg/communication"
	"github.com/loophole/cli/internal/pkg/token"
	"golang.org/x/crypto/ssh"
)

// RegistrationSuccessResponse defines the json format in which the registration success response is returned
type RegistrationSuccessResponse struct {
	SiteID string `json:"siteId"`
	Domain string `json:"domain"`
}

// InfoSuccessResponse defines the json format in which the info success response is returned
type InfoSuccessResponse struct {
	Version string `json:"version"`
}

// ErrorResponse defines the json format in which the error response is returned
type ErrorResponse struct {
	StatusCode int32  `json:"statusCode"`
	Message    string `json:"message"`
	Error      string `json:"error"`
}

// RequestError is an error returned when the request finished with failure
type RequestError struct {
	Message    string
	Details    string
	StatusCode int
}

func (err RequestError) Error() string {
	return fmt.Sprintf("Request Error (%d): %s - %s", err.StatusCode, err.Message, err.Details)
}

var isTokenSaved = token.IsTokenSaved
var getAccessToken = token.GetAccessToken
var tokenWasRefreshed = false
var apiURL = config.Config.APIEndpoint.URI()

// RegisterSite is a funtion used to obtain site id and register keys in the gateway
func RegisterSite(publicKey ssh.PublicKey, requestedSiteID string) (*RegistrationSuccessResponse, error) {
	publicKeyString := publicKey.Type() + " " + base64.StdEncoding.EncodeToString(publicKey.Marshal())

	if !isTokenSaved() {
		return nil, RequestError{
			Message:    "You're not logged in",
			Details:    "Cannot read locally stored token",
			StatusCode: http.StatusUnauthorized,
		}
	}

	accessToken, err := getAccessToken()
	if err != nil {
		return nil, RequestError{
			Message:    "There was a problem reading token",
			Details:    err.Error(),
			StatusCode: http.StatusUnauthorized,
		}
	}

	data := map[string]string{
		"key": publicKeyString,
	}
	if requestedSiteID != "" {
		data["id"] = requestedSiteID
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/site", apiURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	var netClient = &http.Client{
		Timeout:   time.Second * 30,
		Transport: netTransport,
	}

	resp, err := netClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		errorResponse := ErrorResponse{}
		err := json.NewDecoder(resp.Body).Decode(&errorResponse)
		if err != nil {
			return nil, err
		}

		switch resp.StatusCode {
		case http.StatusBadRequest:
			return nil, RequestError{
				Message: errorResponse.Message,
				Details: `The given hostname didn't match the requirements:
- Starts with a letter
- Contains only small letters, numbers and single dashes (-) between them
- Ends with a small letter or number`,
				StatusCode: resp.StatusCode,
			}
		case http.StatusUnauthorized:
			if !tokenWasRefreshed {
				err := token.RefreshToken()
				if err != nil {
					return nil, RequestError{
						Message:    "Authentication failed, then refreshing token failed",
						Details:    errorResponse.Message,
						StatusCode: resp.StatusCode,
					}
				}
				tokenWasRefreshed = true
				return RegisterSite(publicKey, requestedSiteID)
			}
			return nil, RequestError{
				Message:    "Authentication failed, try logging out and logging in again",
				Details:    errorResponse.Message,
				StatusCode: resp.StatusCode,
			}

		case http.StatusForbidden:
			return nil, RequestError{
				Message:    "You don't have required permissions to establish tunnel with given parameters",
				Details:    errorResponse.Message,
				StatusCode: resp.StatusCode,
			}
		case http.StatusConflict:
			return nil, RequestError{
				Message:    "The given hostname is already taken by different user",
				Details:    errorResponse.Message,
				StatusCode: resp.StatusCode,
			}
		case http.StatusUnprocessableEntity:
			return nil, RequestError{
				Message: errorResponse.Message,
				Details: `The given hostname didn't match the requirements:
- Starts with a letter
- Contains only small letters, numbers and single dashes (-) between them
- Ends with a small letter or number
- Minimum 6 characters (not applicable for premium users`,
				StatusCode: resp.StatusCode,
			}
		default:
			return nil, RequestError{
				Message:    errorResponse.Message,
				Details:    "Something unexpected happened, please let developers know",
				StatusCode: resp.StatusCode,
			}
		}
	}

	result := RegistrationSuccessResponse{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	communication.Debug(fmt.Sprintf("Site registration response: %v", result))

	return &result, nil
}

func GetLatestAvailableVersion() (*InfoSuccessResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/info", apiURL), bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent())

	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	var netClient = &http.Client{
		Timeout:   time.Second * 30,
		Transport: netTransport,
	}

	resp, err := netClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Request for information failed, skipping")
	}

	result := InfoSuccessResponse{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	communication.Debug(fmt.Sprintf("Info response: %v", result))

	return &result, nil
}

func userAgent() string {
	return fmt.Sprintf("loophole-%s/%s-%s (%s/%s)",
		config.Config.ClientMode,
		config.Config.Version,
		config.Config.CommitHash,
		runtime.GOOS,
		runtime.GOARCH,
	)
}
