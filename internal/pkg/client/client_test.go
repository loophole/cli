package client

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/crypto/ssh"
)

func TestRegisterSiteSuccessOKShouldPropagateWithoutIdProvided(t *testing.T) {
	oldIsTokenSaved := isTokenSaved
	defer func() { isTokenSaved = oldIsTokenSaved }()
	isTokenSaved = func() bool { return true }

	oldGetAccessToken := getAccessToken
	defer func() { getAccessToken = oldGetAccessToken }()
	getAccessToken = func() (string, error) { return "some-token", nil }

	expectedSiteID := "randomidwhichissuperlong"
	expectedStatus := http.StatusOK
	srv := serverMock(expectedStatus, fmt.Sprintf(`{
		"siteId": "%s"
	}`, expectedSiteID))
	defer srv.Close()
	publicKey, err := getPublicSSHKey()
	if err != nil {
		t.Fatal(err)
	}
	result, err := RegisterSite(srv.URL, publicKey, "")

	if result.ResultCode != expectedStatus {
		t.Fatalf("Status code is different than expected: %d", result.ResultCode)
	}
	if result.SiteID != expectedSiteID {
		t.Fatalf("Site ID '%s' is different than expected: %s", expectedSiteID, result.SiteID)
	}
}

func TestRegisterSiteSuccessCreatedShouldPropagateWithoutIdProvided(t *testing.T) {
	oldIsTokenSaved := isTokenSaved
	defer func() { isTokenSaved = oldIsTokenSaved }()
	isTokenSaved = func() bool { return true }

	oldGetAccessToken := getAccessToken
	defer func() { getAccessToken = oldGetAccessToken }()
	getAccessToken = func() (string, error) { return "some-token", nil }

	expectedSiteID := "randomidwhichissuperlong"
	expectedStatus := http.StatusCreated
	srv := serverMock(expectedStatus, fmt.Sprintf(`{
		"siteId": "%s"
	}`, expectedSiteID))
	defer srv.Close()
	publicKey, err := getPublicSSHKey()
	if err != nil {
		t.Fatal(err)
	}
	result, err := RegisterSite(srv.URL, publicKey, "")

	if result.ResultCode != expectedStatus {
		t.Fatalf("Status code is different than expected: %d", result.ResultCode)
	}
	if result.SiteID != expectedSiteID {
		t.Fatalf("Site ID '%s' is different than expected: %s", expectedSiteID, result.SiteID)
	}
}

func TestRegisterSiteSuccessOKShouldPropagateWithIdProvided(t *testing.T) {
	oldIsTokenSaved := isTokenSaved
	defer func() { isTokenSaved = oldIsTokenSaved }()
	isTokenSaved = func() bool { return true }

	oldGetAccessToken := getAccessToken
	defer func() { getAccessToken = oldGetAccessToken }()
	getAccessToken = func() (string, error) { return "some-token", nil }

	expectedSiteID := "providedhostname"
	expectedStatus := http.StatusOK
	srv := serverMock(expectedStatus, fmt.Sprintf(`{
		"siteId": "%s"
	}`, expectedSiteID))
	defer srv.Close()
	publicKey, err := getPublicSSHKey()
	if err != nil {
		t.Fatal(err)
	}
	result, err := RegisterSite(srv.URL, publicKey, expectedSiteID)

	if result.ResultCode != expectedStatus {
		t.Fatalf("Status code '%d' is different than expected: %d", result.ResultCode, expectedStatus)
	}
	if result.SiteID != expectedSiteID {
		t.Fatalf("Site ID '%s' is different than expected: %s", result.SiteID, expectedSiteID)
	}
}

func TestRegisterSiteSuccessCreatedShouldPropagateWithIdProvided(t *testing.T) {
	oldIsTokenSaved := isTokenSaved
	defer func() { isTokenSaved = oldIsTokenSaved }()
	isTokenSaved = func() bool { return true }

	oldGetAccessToken := getAccessToken
	defer func() { getAccessToken = oldGetAccessToken }()
	getAccessToken = func() (string, error) { return "some-token", nil }

	expectedSiteID := "providedhostname"
	expectedStatus := http.StatusCreated
	srv := serverMock(expectedStatus, fmt.Sprintf(`{
		"siteId": "%s"
	}`, expectedSiteID))
	defer srv.Close()
	publicKey, err := getPublicSSHKey()
	if err != nil {
		t.Fatal(err)
	}
	result, err := RegisterSite(srv.URL, publicKey, expectedSiteID)

	if result.ResultCode != expectedStatus {
		t.Fatalf("Status code '%d' is different than expected: %d", result.ResultCode, expectedStatus)
	}
	if result.SiteID != expectedSiteID {
		t.Fatalf("Site ID '%s' is different than expected: %s", result.SiteID, expectedSiteID)
	}
}

func TestRegisterSiteError400ShouldPropagateError(t *testing.T) {
	oldIsTokenSaved := isTokenSaved
	defer func() { isTokenSaved = oldIsTokenSaved }()
	isTokenSaved = func() bool { return true }

	oldGetAccessToken := getAccessToken
	defer func() { getAccessToken = oldGetAccessToken }()
	getAccessToken = func() (string, error) { return "some-token", nil }

	expectedErrorMessage := "You did something bad"
	expectedStatus := http.StatusBadRequest
	srv := serverMock(expectedStatus, fmt.Sprintf(`{
		"statusCode": "400",
		"error": "Bad request",
		"message": "%s"
	}`, expectedErrorMessage))
	defer srv.Close()
	publicKey, err := getPublicSSHKey()
	if err != nil {
		t.Fatal(err)
	}
	result, err := RegisterSite(srv.URL, publicKey, "")

	if err == nil {
		t.Fatalf("Expected an error to be returned")
	}
	if result.SiteID != "" {
		t.Fatalf("Expected site ID to be empty, got %s", result.SiteID)
	}
	if err.Error() != expectedErrorMessage {
		t.Fatalf("Expected '%s' error, got '%s", expectedErrorMessage, err.Error())
	}
}

func TestRegisterSiteError401ShouldPropagateError(t *testing.T) {
	oldIsTokenSaved := isTokenSaved
	defer func() { isTokenSaved = oldIsTokenSaved }()
	isTokenSaved = func() bool { return true }

	oldGetAccessToken := getAccessToken
	defer func() { getAccessToken = oldGetAccessToken }()
	getAccessToken = func() (string, error) { return "some-token", nil }

	expectedErrorMessage := "You are not authenticated"
	expectedStatus := http.StatusUnauthorized
	srv := serverMock(expectedStatus, fmt.Sprintf(`{
		"statusCode": "401",
		"error": "Unauthorized",
		"message": "%s"
	}`, expectedErrorMessage))
	defer srv.Close()
	publicKey, err := getPublicSSHKey()
	if err != nil {
		t.Fatal(err)
	}
	result, err := RegisterSite(srv.URL, publicKey, "")

	if err == nil {
		t.Fatalf("Expected an error to be returned")
	}
	if result.SiteID != "" {
		t.Fatalf("Expected site ID to be empty, got %s", result.SiteID)
	}
	if err.Error() != expectedErrorMessage {
		t.Fatalf("Expected '%s' error, got '%s", expectedErrorMessage, err.Error())
	}
}

func TestRegisterSiteError403ShouldPropagateError(t *testing.T) {
	oldIsTokenSaved := isTokenSaved
	defer func() { isTokenSaved = oldIsTokenSaved }()
	isTokenSaved = func() bool { return true }

	oldGetAccessToken := getAccessToken
	defer func() { getAccessToken = oldGetAccessToken }()
	getAccessToken = func() (string, error) { return "some-token", nil }

	expectedErrorMessage := "You are not allowed to do this"
	expectedStatus := http.StatusForbidden
	srv := serverMock(expectedStatus, fmt.Sprintf(`{
		"statusCode": "403",
		"error": "Forbidden",
		"message": "%s"
	}`, expectedErrorMessage))
	defer srv.Close()
	publicKey, err := getPublicSSHKey()
	if err != nil {
		t.Fatal(err)
	}
	result, err := RegisterSite(srv.URL, publicKey, "")

	if err == nil {
		t.Fatalf("Expected an error to be returned")
	}
	if result.SiteID != "" {
		t.Fatalf("Expected site ID to be empty, got %s", result.SiteID)
	}
	if err.Error() != expectedErrorMessage {
		t.Fatalf("Expected '%s' error, got '%s", expectedErrorMessage, err.Error())
	}
}

func TestRegisterTokenNotSavedReturns600(t *testing.T) {
	oldIsTokenSaved := isTokenSaved
	defer func() { isTokenSaved = oldIsTokenSaved }()
	isTokenSaved = func() bool { return false }

	oldGetAccessToken := getAccessToken
	defer func() { getAccessToken = oldGetAccessToken }()
	getAccessToken = func() (string, error) { return "some-token", nil }

	expectedErrorMessage := "Please log in before using Loophole"
	expectedStatus := http.StatusOK
	srv := serverMock(expectedStatus, `{
		"siteId": "whateverrrr"
	}`)
	defer srv.Close()
	publicKey, err := getPublicSSHKey()
	if err != nil {
		t.Fatal(err)
	}
	result, err := RegisterSite(srv.URL, publicKey, "")

	if err == nil {
		t.Fatalf("Expected an error to be returned")
	}
	if result.SiteID != "" {
		t.Fatalf("Expected site ID to be empty, got %s", result.SiteID)
	}
	if err.Error() != expectedErrorMessage {
		t.Fatalf("Expected '%s' error, got '%s", expectedErrorMessage, err.Error())
	}
}

func TestRegisterTokenReadingProblemReturns601(t *testing.T) {
	oldIsTokenSaved := isTokenSaved
	defer func() { isTokenSaved = oldIsTokenSaved }()
	isTokenSaved = func() bool { return true }

	oldGetAccessToken := getAccessToken
	defer func() { getAccessToken = oldGetAccessToken }()
	getAccessToken = func() (string, error) { return "", fmt.Errorf("Something bad happened when reading token") }

	expectedErrorMessage := "There was a problem reading token"
	expectedStatus := http.StatusOK
	srv := serverMock(expectedStatus, `{
		"siteId": "whateverrrr"
	}`)
	defer srv.Close()
	publicKey, err := getPublicSSHKey()
	if err != nil {
		t.Fatal(err)
	}
	result, err := RegisterSite(srv.URL, publicKey, "")

	if err == nil {
		t.Fatalf("Expected an error to be returned")
	}
	if result.SiteID != "" {
		t.Fatalf("Expected site ID to be empty, got %s", result.SiteID)
	}
	if err.Error() != expectedErrorMessage {
		t.Fatalf("Expected '%s' error, got '%s", expectedErrorMessage, err.Error())
	}
}

func serverMock(httpStatus int, expectedResponse string) *httptest.Server {
	handler := http.NewServeMux()
	registerSiteMock := getRegisterSiteHandler(httpStatus, expectedResponse)
	handler.HandleFunc("/api/register-site", registerSiteMock)

	srv := httptest.NewServer(handler)

	return srv
}

func getRegisterSiteHandler(httpStatus int, expectedResponse string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(httpStatus)
		_, _ = w.Write([]byte(expectedResponse))
	}
}

func getPublicSSHKey() (ssh.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return nil, err
	}

	// generate and write public key
	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, err
	}
	return publicKey, nil
}
