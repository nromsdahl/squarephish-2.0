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

	"github.com/nromsdahl/squarephish2/internal/database"
	"github.com/nromsdahl/squarephish2/internal/models"
	log "github.com/sirupsen/logrus"
)

// authPoll continuously polls the token endpoint waiting for valid user
// authentication.
// Parameters:
//   - email: The email address of the authenticating user.
//   - deviceCode: The device code JSON object.
//   - entraConfig: The Entra configuration.
//   - requestConfig: The request configuration.
//
// It returns nothing.
func authPoll(email string, deviceCode models.DeviceCodeResponse, entraConfig models.EntraConfig, requestConfig models.RequestConfig) {
	tokenURL := "https://login.microsoftonline.com/" + entraConfig.Tenant + "/oauth2/v2.0/token"
	grantType := "urn:ietf:params:oauth:grant-type:device_code"
	contentType := "application/x-www-form-urlencoded"

	startTime := time.Now()
	duration := time.Duration(deviceCode.ExpiresIn) * time.Second
	expirationTime := startTime.Add(duration)

	tokenParams := url.Values{}
	tokenParams.Set("grant_type", grantType)
	tokenParams.Set("code", deviceCode.DeviceCode)
	tokenParams.Set("client_id", entraConfig.ClientID)

	// Encode the parameters into a URL-encoded string
	// (e.g., "client_id=...&scope=...")
	tokenPostData := tokenParams.Encode()

	// Initialize reusable variables for error handling
	// and response handling. 'resp' and 'responseBodyBytes' are
	// used outside of the for loop on valid authentication.
	var err error
	var resp *http.Response
	var responseBodyBytes []byte

	for {
		log.Printf("[%s] Polling for user authentication...\n", email)

		// Create an io.Reader from the encoded data string.
		// http.Post requires the request body as an io.Reader.
		tokenBodyReader := strings.NewReader(tokenPostData)

		// Create a custom HTTP client
		client := &http.Client{
			Timeout: 15 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}

		// Create a custom HTTP request with the specified User Agent
		req, err := http.NewRequest("POST", tokenURL, tokenBodyReader)
		if err != nil {
			return
		}

		req.Header.Set("Content-Type", contentType)
		req.Header.Set("User-Agent", requestConfig.UserAgent)

		resp, err = client.Do(req)
		if err != nil {
			return
		}

		// Ensure the response body is closed even if errors occur later.
		defer resp.Body.Close()

		// Authentication successful
		if resp.StatusCode == 200 {
			log.Printf("[%s] Authentication successful", email)
			break
		}

		// Check pending status
		responseBodyBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("[%s] Error reading token response body:\n%v", email, err)
			return
		}

		var pendingResult models.PendingTokenResponse
		if err = json.Unmarshal(responseBodyBytes, &pendingResult); err != nil {
			log.Printf("[%s] Error unmarshalling pending response:\n%v", email, err)
			return
		}

		// Check if the error is not authorization_pending
		if pendingResult.Error != "authorization_pending" {
			log.Printf("[%s] Invalid error response:\n%v", email, pendingResult)
			return
		}

		// Check if the polling time has expired
		if time.Now().After(expirationTime) {
			log.Printf("[%s] Device code authentication expired", email)
			return
		}

		time.Sleep(time.Duration(deviceCode.Interval) * time.Second)
	}

	responseBodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[%s] Error reading token response body:\n%v", email, err)
		return
	}

	var tokenResult models.TokenResponse
	if err = json.Unmarshal(responseBodyBytes, &tokenResult); err != nil {
		log.Printf("[%s] Error unmarshalling token response:\n%v", email, err)
		return
	}

	// log.Printf("AccessToken:  %s\n", tokenResult.AccessToken)
	// log.Printf("RefreshToken: %s\n", tokenResult.RefreshToken)
	log.Printf("[%s] Token retrieved and saved to database\n", email)

	// Save the token to the database
	var jsonCredential []byte
	jsonCredential, err = json.Marshal(tokenResult)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	err = database.SaveToken(email, string(jsonCredential))
	if err != nil {
		log.Printf("Error saving credential: %v", err)
	}
}
