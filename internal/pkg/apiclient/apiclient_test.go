package apiclient

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
	srv := serverMock(http.StatusCreated, fmt.Sprintf(`{
		"siteId": "%s"
	}`, expectedSiteID))
	defer srv.Close()

	apiURL = srv.URL
	publicKey, err := getPublicSSHKey()
	if err != nil {
		t.Fatal(err)
	}
	result, err := RegisterSite(publicKey, "")

	if err != nil {
		t.Fatalf("Unexpected error returned: %v", err)
	}
	if result.SiteID != expectedSiteID {
		t.Fatalf("Site ID '%s' is different than expected: %s", expectedSiteID, result)
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
	srv := serverMock(http.StatusCreated, fmt.Sprintf(`{
		"siteId": "%s"
	}`, expectedSiteID))
	defer srv.Close()

	apiURL = srv.URL
	publicKey, err := getPublicSSHKey()
	if err != nil {
		t.Fatal(err)
	}
	result, err := RegisterSite(publicKey, "")

	if err != nil {
		t.Fatalf("Unexpected error returned: %v", err)
	}
	if result.SiteID != expectedSiteID {
		t.Fatalf("Site ID '%s' is different than expected: %s", expectedSiteID, result)
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

	apiURL = srv.URL
	publicKey, err := getPublicSSHKey()
	if err != nil {
		t.Fatal(err)
	}
	result, err := RegisterSite(publicKey, expectedSiteID)

	if err != nil {
		t.Fatalf("Unexpected error returned: %v", err)
	}
	if result.SiteID != expectedSiteID {
		t.Fatalf("Site ID '%s' is different than expected: %s", result, expectedSiteID)
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
	srv := serverMock(http.StatusCreated, fmt.Sprintf(`{
		"siteId": "%s"
	}`, expectedSiteID))
	defer srv.Close()

	apiURL = srv.URL
	publicKey, err := getPublicSSHKey()
	if err != nil {
		t.Fatal(err)
	}
	result, err := RegisterSite(publicKey, expectedSiteID)

	if err != nil {
		t.Fatalf("Unexpected error returned: %v", err)
	}
	if result.SiteID != expectedSiteID {
		t.Fatalf("Site ID '%s' is different than expected: %s", result, expectedSiteID)
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
		"statusCode": 400,
		"error": "Bad request",
		"message": "%s"
	}`, expectedErrorMessage))
	defer srv.Close()

	apiURL = srv.URL
	publicKey, err := getPublicSSHKey()
	if err != nil {
		t.Fatal(err)
	}
	result, err := RegisterSite(publicKey, "")

	if err == nil {
		t.Fatalf("Expected an error to be returned")
	}
	requestErr, ok := err.(RequestError)
	if !ok {
		t.Fatalf("Expected RequestError to be returned")
	}
	if result != nil {
		t.Fatalf("Expected result to be nil, got %v", result)
	}
	if requestErr.StatusCode != expectedStatus {
		t.Fatalf("Expected '%d' status, got '%d", expectedStatus, requestErr.StatusCode)
	}
	if requestErr.Message != expectedErrorMessage {
		t.Fatalf("Expected '%s' message, got '%s", expectedErrorMessage, requestErr.Message)
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
		"statusCode": 401,
		"error": "Unauthorized",
		"message": "%s"
	}`, expectedErrorMessage))
	defer srv.Close()

	apiURL = srv.URL
	publicKey, err := getPublicSSHKey()
	if err != nil {
		t.Fatal(err)
	}
	result, err := RegisterSite(publicKey, "")

	if err == nil {
		t.Fatalf("Expected an error to be returned")
	}
	requestErr, ok := err.(RequestError)
	if !ok {
		t.Fatalf("Expected RequestError to be returned")
	}
	if result != nil {
		t.Fatalf("Expected result to be nil, got %v", result)
	}
	if requestErr.StatusCode != expectedStatus {
		t.Fatalf("Expected '%d' status, got '%d", expectedStatus, requestErr.StatusCode)
	}
	if requestErr.Details != expectedErrorMessage {
		t.Fatalf("Expected '%s' message, got '%s", expectedErrorMessage, requestErr.Details)
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
		"statusCode": 403,
		"error": "Forbidden",
		"message": "%s"
	}`, expectedErrorMessage))
	defer srv.Close()

	apiURL = srv.URL
	publicKey, err := getPublicSSHKey()
	if err != nil {
		t.Fatal(err)
	}
	result, err := RegisterSite(publicKey, "")

	if err == nil {
		t.Fatalf("Expected an error to be returned")
	}
	requestErr, ok := err.(RequestError)
	if !ok {
		t.Fatalf("Expected RequestError to be returned")
	}
	if result != nil {
		t.Fatalf("Expected result to be nil, got %v", result)
	}
	if requestErr.StatusCode != expectedStatus {
		t.Fatalf("Expected '%d' status, got '%d", expectedStatus, requestErr.StatusCode)
	}
	if requestErr.Details != expectedErrorMessage {
		t.Fatalf("Expected '%s' message, got '%s", expectedErrorMessage, requestErr.Details)
	}
}

func TestRegisterTokenNotSavedReturns401(t *testing.T) {
	oldIsTokenSaved := isTokenSaved
	defer func() { isTokenSaved = oldIsTokenSaved }()
	isTokenSaved = func() bool { return false }

	oldGetAccessToken := getAccessToken
	defer func() { getAccessToken = oldGetAccessToken }()
	getAccessToken = func() (string, error) { return "some-token", nil }

	srv := serverMock(http.StatusOK, `{
		"siteId": "whateverrrr"
	}`)
	defer srv.Close()

	apiURL = srv.URL
	publicKey, err := getPublicSSHKey()
	if err != nil {
		t.Fatal(err)
	}
	result, err := RegisterSite(publicKey, "")

	if err == nil {
		t.Fatalf("Expected an error to be returned")
	}
	requestErr, ok := err.(RequestError)
	if !ok {
		t.Fatalf("Expected RequestError to be returned")
	}
	if result != nil {
		t.Fatalf("Expected result to be nil, got %v", result)
	}
	if requestErr.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Expected '%d' status, got '%d", http.StatusUnauthorized, requestErr.StatusCode)
	}
}

func TestRegisterTokenReadingProblemReturns401(t *testing.T) {
	oldIsTokenSaved := isTokenSaved
	defer func() { isTokenSaved = oldIsTokenSaved }()
	isTokenSaved = func() bool { return true }

	oldGetAccessToken := getAccessToken
	defer func() { getAccessToken = oldGetAccessToken }()
	getAccessToken = func() (string, error) { return "", fmt.Errorf("Something bad happened when reading token") }

	srv := serverMock(http.StatusOK, `{
		"siteId": "whateverrrr"
	}`)
	defer srv.Close()

	apiURL = srv.URL

	publicKey, err := getPublicSSHKey()
	if err != nil {
		t.Fatal(err)
	}
	result, err := RegisterSite(publicKey, "")

	if err == nil {
		t.Fatalf("Expected an error to be returned")
	}
	requestErr, ok := err.(RequestError)
	if !ok {
		t.Fatalf("Expected RequestError to be returned")
	}
	if result != nil {
		t.Fatalf("Expected result to be nil, got %v", result)
	}
	if requestErr.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Expected '%d' status, got '%d", http.StatusUnauthorized, requestErr.StatusCode)
	}
}

func serverMock(httpStatus int, expectedResponse string) *httptest.Server {
	handler := http.NewServeMux()
	registerSiteMock := getRegisterSiteHandler(httpStatus, expectedResponse)
	handler.HandleFunc("/api/site", registerSiteMock)

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
