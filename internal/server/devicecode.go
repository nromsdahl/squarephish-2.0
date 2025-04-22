package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/nromsdahl/squarephish2/internal/models"
	log "github.com/sirupsen/logrus"
)

// initDeviceCode will initialize the device code flow.
// Parameters:
//   - email: The email address of the authenticating user.
//
// It returns the device code resopnse object or an error if device code fails.
func initDeviceCode(email string) (models.DeviceCodeResponse, error) {
	var deviceCodeResult models.DeviceCodeResponse

	scope := ".default offline_access profile openid"
	clientID := "29d9ed98-a469-4536-ade2-f981bc1d605e" // Microsoft Authentication Broker
	contentType := "application/x-www-form-urlencoded"
	deviceCodeURL := "https://login.microsoftonline.com/organizations/oauth2/v2.0/deviceCode"

	deviceCodeParams := url.Values{}
	deviceCodeParams.Set("client_id", clientID)
	deviceCodeParams.Set("scope", scope)

	// Encode the parameters into a URL-encoded string
	// (e.g., "client_id=...&scope=...")
	deviceCodePostData := deviceCodeParams.Encode()

	// Create an io.Reader from the encoded data string.
	// http.Post requires the request body as an io.Reader.
	deviceCodeBodyReader := strings.NewReader(deviceCodePostData)

	log.Printf("[%s] Initializing device code flow...\n", email)
	log.Printf("[%s]     Client ID: %s\n", email, clientID)
	log.Printf("[%s]     Scope:     %s\n", email, scope)

	resp, err := http.Post(deviceCodeURL, contentType, deviceCodeBodyReader)
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
