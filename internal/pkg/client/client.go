package client

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

// SiteSpecification is struct containing information about site registration results
type SiteSpecification struct {
	SiteID     string
	ResultCode int
}

// SuccessResponse defines the json format in which the success response is returned
type SuccessResponse struct {
	SiteID string `json:"siteId"`
}

// ErrorResponse defines the json format in which the error response is returned
type ErrorResponse struct {
	StatusCode string `json:"statusCode"`
	Message    string `json:"message"`
	Error      string `json:"error"`
}

var isTokenSaved = token.IsTokenSaved
var getAccessToken = token.GetAccessToken

// RegisterSite is a funtion used to obtain site id and register keys in the gateway
func RegisterSite(apiURL string, publicKey ssh.PublicKey, siteID string) (SiteSpecification, error) {
	publicKeyString := publicKey.Type() + " " + base64.StdEncoding.EncodeToString(publicKey.Marshal())

	if !isTokenSaved() {
		return SiteSpecification{"", 601}, fmt.Errorf("Please log in before using Loophole")
	}

	accessToken, err := getAccessToken()
	if err != nil {
		return SiteSpecification{"", 600}, fmt.Errorf("There was a problem reading token")
	}

	data := map[string]string{
		"key": publicKeyString,
	}
	if siteID != "" {
		data["id"] = siteID
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return SiteSpecification{"", 0}, fmt.Errorf("There was a problem encoding request body")
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/register-site", apiURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return SiteSpecification{"", 0}, fmt.Errorf("There was a problem creating request")
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return SiteSpecification{"", 0}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		errorResponse := ErrorResponse{}
		json.NewDecoder(resp.Body).Decode(&errorResponse)

		return SiteSpecification{"", resp.StatusCode}, fmt.Errorf("%s", errorResponse.Message)
	}

	result := SuccessResponse{}
	json.NewDecoder(resp.Body).Decode(&result)

	if el := log.Debug(); el.Enabled() {
		fmt.Println()
		el.Interface("result", result).Msg("Response")
	}

	return SiteSpecification{result.SiteID, resp.StatusCode}, nil
}
