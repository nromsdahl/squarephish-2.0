package database

import (
	"context"
	"time"

	"github.com/nromsdahl/squarephish2/internal/models"
	log "github.com/sirupsen/logrus"
)

// This file contains the public API for the database package.
// It provides functions to initialize the database, load and save data,
// and close the database connection.

// --- Database Initialization ---

// InitializeDatabase initializes the database
// It returns an error if the database fails to initialize.
func Initialize() error {
	err := connect()
	if err != nil {
		log.Errorf("failed to connect to database: %v", err)
	}

	err = migrate()
	if err != nil {
		log.Errorf("failed to migrate database: %v", err)
	}

	return nil
}

func Close() error {
	err := close()
	if err != nil {
		log.Errorf("failed to close database: %v", err)
	}

	return nil
}

// --- Data Loading ---

// LoadConfig fetches the current configuration from the database
// It returns a ConfigData struct containing the configuration values.
func LoadConfig() models.ConfigData {
	ctx, cancel := context.WithTimeout(context.Background(), (3 * time.Second))
	defer cancel()

	// Fetch string values using the helper function
	smtpHost := queryConfigValue(ctx, "smtpHost", "")
	smtpPort := queryConfigValue(ctx, "smtpPort", "")
	smtpUsername := queryConfigValue(ctx, "smtpUsername", "")
	smtpPassword := queryConfigValue(ctx, "smtpPassword", "")

	phishHost := queryConfigValue(ctx, "phishHost", "")
	phishPort := queryConfigValue(ctx, "phishPort", "")

	emailSender := queryConfigValue(ctx, "emailSender", "")
	emailSubject := queryConfigValue(ctx, "emailSubject", "")
	emailBody := queryConfigValue(ctx, "emailBody", "")

	// Default to Microsoft Authentication Broker
	entraClientID := queryConfigValue(ctx, "entraClientID", "29d9ed98-a469-4536-ade2-f981bc1d605e")
	entraScope := queryConfigValue(ctx, "entraScope", ".default offline_access profile openid")
	entraTenant := queryConfigValue(ctx, "entraTenant", "organizations")

	// Default to Microsoft Edge on Windows User Agent
	userAgent := queryConfigValue(ctx, "userAgent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36 Edg/135.0.3179.85")

	// --- Populate the struct ---
	return models.ConfigData{
		SMTPConfig: models.SMTPConfig{
			Host:     smtpHost,
			Port:     smtpPort,
			Username: smtpUsername,
			Password: smtpPassword,
		},
		SquarePhishConfig: models.SquarePhishConfig{
			Host: phishHost,
			Port: phishPort,
		},
		EmailConfig: models.EmailConfig{
			Sender:  emailSender,
			Subject: emailSubject,
			Body:    emailBody,
		},
		EntraConfig: models.EntraConfig{
			ClientID: entraClientID,
			Scope:    entraScope,
			Tenant:   entraTenant,
		},
		RequestConfig: models.RequestConfig{
			UserAgent: userAgent,
		},
		ActivePage: "config",
		Title:      "Configuration Settings",
	}
}

// LoadDashboardData fetches metrics from the database
// It returns a DashboardData struct containing the metrics.
func LoadDashboardData() models.DashboardData {
	ctx, cancel := context.WithTimeout(context.Background(), (3 * time.Second))
	defer cancel()

	// Fetch each metric using the helper function
	emailsSent := queryMetric(ctx, "emails_sent")
	emailsScanned := queryMetric(ctx, "emails_scanned")

	credentials, err := queryCredentials(ctx)
	if err != nil {
		log.Printf("error querying credentials: %v", err)
		credentials = []models.Credential{}
	}

	return models.DashboardData{
		EmailsSent:    emailsSent,
		EmailsScanned: emailsScanned,
		Credentials:   credentials,
		ActivePage:    "dashboard",
		Title:         "Dashboard",
	}
}

// LoadToken fetches the token from the database
// Parameters:
//   - email: The email of the credential
//
// It returns the token of the credential.
func LoadToken(email string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), (3 * time.Second))
	defer cancel()

	token, err := queryToken(ctx, email)
	if err != nil {
		log.Printf("error querying token: %v", err)
		return "", err
	}

	return token, nil
}

// --- Data Saving ---

// SaveConfig saves the current configuration to the database
// Parameters:
//   - config: The configuration to save
//
// It returns an error if the statement fails to prepare or execute.
func SaveConfig(config models.ConfigData) error {
	return insertConfig(config)
}

// SaveToken saves a new credential to the database
// Parameters:
//   - email: The email of the credential
//   - token: The token of the credential
//
// It returns an error if the statement fails to prepare or execute.
func SaveToken(email, token string) error {
	return insertCredential(email, token)
}

// SaveEmailSent saves an email to the database
// Parameters:
//   - recipient: The recipient of the email
//   - subject: The subject of the email
//
// It returns an error if the statement fails to prepare or execute.
func SaveEmailSent(recipient, subject string) error {
	return insertEmailSent(recipient, subject)
}

// SaveEmailScanned saves an email to the database
// Parameters:
//   - recipient: The recipient of the email
//
// It returns an error if the statement fails to prepare or execute.
func SaveEmailScanned(recipient string) error {
	return insertEmailScanned(recipient)
}
