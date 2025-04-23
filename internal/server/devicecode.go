package server

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/nromsdahl/squarephish2/internal/models"
	log "github.com/sirupsen/logrus"
)

// initDeviceCode will initialize the device code flow.
// Parameters:
//   - email: The email address of the authenticating user.
//   - entraConfig: The Entra configuration.
//   - requestConfig: The request configuration.
//
// It returns the device code resopnse object or an error if device code fails.
func initDeviceCode(email string, entraConfig models.EntraConfig, requestConfig models.RequestConfig) (models.DeviceCodeResponse, error) {
	var deviceCodeResult models.DeviceCodeResponse

	contentType := "application/x-www-form-urlencoded"
	deviceCodeURL := "https://login.microsoftonline.com/" + entraConfig.Tenant + "/oauth2/v2.0/deviceCode"

	deviceCodeParams := url.Values{}
	deviceCodeParams.Set("client_id", entraConfig.ClientID)
	deviceCodeParams.Set("scope", entraConfig.Scope)

	// Encode the parameters into a URL-encoded string
	// (e.g., "client_id=...&scope=...")
	deviceCodePostData := deviceCodeParams.Encode()

	// Create an io.Reader from the encoded data string.
	// http.Post requires the request body as an io.Reader.
	deviceCodeBodyReader := strings.NewReader(deviceCodePostData)

	log.Printf("[%s] Initializing device code flow...\n", email)
	log.Printf("[%s]     Client ID: %s\n", email, entraConfig.ClientID)
	log.Printf("[%s]     Scope:     %s\n", email, entraConfig.Scope)

	// Create a custom HTTP client
	client := &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// Create a custom HTTP request with the specified User Agent
	req, err := http.NewRequest("POST", deviceCodeURL, deviceCodeBodyReader)
	if err != nil {
		return deviceCodeResult, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("User-Agent", requestConfig.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return deviceCodeResult, fmt.Errorf("error sending POST request: %w", err)
	}

	// Ensure the response body is closed even if errors occur later.
	defer resp.Body.Close()

	// Check if the status code indicates failure
	if resp.StatusCode != 200 {
		return deviceCodeResult, fmt.Errorf("warning: Received non-success status code %d", resp.StatusCode)
	}

	responseBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return deviceCodeResult, fmt.Errorf("error reading response body: %w", err)
	}

	if err := json.Unmarshal(responseBodyBytes, &deviceCodeResult); err != nil {
		return deviceCodeResult, fmt.Errorf("error unmarshalling response body: %w", err)
	}

	return deviceCodeResult, nil
}
