package models

type DeviceCodeSpec struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
	VerificationURI         string `json:"verification_uri"`
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
