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

type SiteSpecification struct {
	SiteID     string
	ResultCode int
}

func RegisterSite(apiURL string, publicKey ssh.PublicKey, siteID string) (SiteSpecification, error) {
	publicKeyString := publicKey.Type() + " " + base64.StdEncoding.EncodeToString(publicKey.Marshal())

	if !token.IsTokenSaved() {
		return SiteSpecification{"", 601}, fmt.Errorf("Please log in before using Loophole")
	}

	accessToken, err := token.GetAccessToken()
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
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/site", apiURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return SiteSpecification{"", 0}, fmt.Errorf("There was a problem creating request")
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return SiteSpecification{"", 0}, err
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return SiteSpecification{"", resp.StatusCode}, fmt.Errorf("Site registration request ended with %d status", resp.StatusCode)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if el := log.Debug(); el.Enabled() {
		fmt.Println()
		el.Interface("result", result).Msg("Response")
	}
	defer resp.Body.Close()

	siteID, ok := result["siteId"].(string)
	if !ok {
		return SiteSpecification{"", 400}, fmt.Errorf("Error converting siteId to string")
	}
	return SiteSpecification{siteID, resp.StatusCode}, nil
}
