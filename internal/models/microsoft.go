package models

// These are the models for the Microsoft device code flow.

// https://login.microsoftonline.com/organizations/oauth2/v2.0/deviceCode
// DeviceCodeResponse contains the response for the initialization of the
// Microsoft device code flow.
type DeviceCodeResponse struct {
	UserCode        string `json:"user_code"`
	DeviceCode      string `json:"device_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
	Message         string `json:"message"`
}

// https://login.microsoftonline.com/organizations/oauth2/v2.0/token
// PendingTokenResponse contains the response for the pending token request.
type PendingTokenResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorCodes       []int  `json:"error_codes"`
	Timestamp        string `json:"timestamp"`
	TraceID          string `json:"trace_id"`
	CorrelationID    string `json:"correlation_id"`
}

// https://login.microsoftonline.com/organizations/oauth2/v2.0/token
// TokenResponse contains the response for a valid authentication token
// request.
type TokenResponse struct {
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
}
