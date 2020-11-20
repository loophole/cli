package apiclient

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/loophole/cli/internal/pkg/token"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh"
)

// SuccessResponse defines the json format in which the success response is returned
type SuccessResponse struct {
	SiteID string `json:"siteId"`
}

// ErrorResponse defines the json format in which the error response is returned
type ErrorResponse struct {
	StatusCode int32  `json:"statusCode"`
	Message    string `json:"message"`
	Error      string `json:"error"`
}

type RequestError struct {
	Message    string
	Details    string
	StatusCode int
}

func (err RequestError) Error() string {
	return err.Message
}

var isTokenSaved = token.IsTokenSaved
var getAccessToken = token.GetAccessToken
var tokenWasRefreshed = false

// RegisterSite is a funtion used to obtain site id and register keys in the gateway
func RegisterSite(apiURL string, publicKey ssh.PublicKey, siteID string) (string, error) {
	publicKeyString := publicKey.Type() + " " + base64.StdEncoding.EncodeToString(publicKey.Marshal())

	if !isTokenSaved() {
		return "", RequestError{
			Message:    "You're not logged in",
			Details:    "Cannot read locally stored token",
			StatusCode: http.StatusUnauthorized,
		}
	}

	accessToken, err := getAccessToken()
	if err != nil {
		return "", RequestError{
			Message:    "There was a problem reading token",
			Details:    err.Error(),
			StatusCode: http.StatusUnauthorized,
		}
	}

	data := map[string]string{
		"key": publicKeyString,
	}
	if siteID != "" {
		data["id"] = siteID
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/register-site", apiURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		errorResponse := ErrorResponse{}
		err := json.NewDecoder(resp.Body).Decode(&errorResponse)
		if err != nil {
			return "", err
		}

		switch resp.StatusCode {
		case http.StatusBadRequest:
			return "", RequestError{
				Message: errorResponse.Message,
				Details: `The given hostname didn't match the requirements:
- Starts with a letter
- Contains only small letters and numbers`,
				StatusCode: resp.StatusCode,
			}
		case http.StatusUnauthorized:
			if !tokenWasRefreshed {
				err := token.RefreshToken()
				if err != nil {
					return "", RequestError{
						Message:    "Authentication failed, then refreshing token failed",
						Details:    errorResponse.Message,
						StatusCode: resp.StatusCode,
					}
				}
				tokenWasRefreshed = true
				return RegisterSite(apiURL, publicKey, siteID)
			}
			return "", RequestError{
				Message:    "Authentication failed, try logging out and logging in again",
				Details:    errorResponse.Message,
				StatusCode: resp.StatusCode,
			}

		case http.StatusForbidden:
			return "", RequestError{
				Message:    "You don't have required permissions to establish tunnel with given parameters",
				Details:    errorResponse.Message,
				StatusCode: resp.StatusCode,
			}
		case http.StatusConflict:
			return "", RequestError{
				Message:    "The given hostname is already taken by different user",
				Details:    errorResponse.Message,
				StatusCode: resp.StatusCode,
			}
		case http.StatusUnprocessableEntity:
			return "", RequestError{
				Message: errorResponse.Message,
				Details: `The given hostname didn't match the requirements:
- Starts with a letter
- Contains only small letters and numbers
- Minimum 6 characters (not applicable for premium users`,
				StatusCode: resp.StatusCode,
			}
		default:
			return "", RequestError{
				Message:    errorResponse.Message,
				Details:    "Something unexpected happened, please let developers know",
				StatusCode: resp.StatusCode,
			}
		}
	}

	result := SuccessResponse{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	if el := log.Debug(); el.Enabled() {
		fmt.Println()
		el.Interface("result", result).Msg("Response")
	}

	return result.SiteID, nil
}
