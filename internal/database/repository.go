package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nromsdahl/squarephish2/internal/models"
)

// --- Data Querying ---

// queryToken is a helper function to query a single token from the DB
// Parameters:
//   - ctx: The context for the query
//   - email: The email of the credential
//
// It returns the token of the credential.
func queryToken(ctx context.Context, email string) (string, error) {
	var token string

	row := db.QueryRowContext(
		ctx,
		"SELECT token FROM credentials WHERE email = ? ORDER BY timestamp DESC LIMIT 1",
		email,
	)

	err := row.Scan(&token)
	if err != nil {
		return "", err
	}

	return token, nil
}

// queryMetric is a helper function to query a single metric value from the DB
// Parameters:
//   - ctx: The context for the query
//   - metricName: The name of the metric table to query
//
// It returns the value of the metric.
func queryMetric(ctx context.Context, metricName string) int {
	var value int

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", metricName)

	row := db.QueryRowContext(ctx, query)
	err := row.Scan(&value)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0
		}

		return 0
	}

	return value
}

// queryConfigValue is a helper function to query a single configuration value
// (TEXT) from the DB
// Parameters:
//   - ctx: The context for the query
//   - key: The key of the configuration value to query
//   - defaultValue: The default value to return if the key is not found
//
// It returns the value of the configuration key.
func queryConfigValue(ctx context.Context, key, defaultValue string) string {
	var value string

	row := db.QueryRowContext(
		ctx,
		"SELECT value FROM configuration WHERE key = ?",
		key,
	)

	err := row.Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return defaultValue
		}

		return defaultValue
	}

	return value
}

// queryCredentials is a helper function to query all credentials from the DB
// Parameters:
//   - ctx: The context for the query
//
// It returns a slice of Credential structs.
func queryCredentials(ctx context.Context) ([]models.Credential, error) {
	var credentials []models.Credential

	query := "SELECT email, token FROM credentials"

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var emailStr string
		var tokenStr string
		var token models.TokenResponse
		var credential models.Credential

		err := rows.Scan(&emailStr, &tokenStr)
		if err != nil {
			return nil, err
		}

		byteData := []byte(tokenStr)
		err = json.Unmarshal(byteData, &token)
		if err != nil {
			return nil, err
		}

		credential.Email = emailStr
		credential.Token = token
		credentials = append(credentials, credential)
	}

	return credentials, nil
}

// --- Data Saving ---

// insertEmailSent inserts a new email record into the emails_sent table
// Parameters:
//   - recipient: The recipient of the email
//   - subject: The subject of the email
//
// It returns an error if the statement fails to prepare or execute.
func insertEmailSent(recipient, subject string) error {
	stmt, err := db.Prepare("INSERT INTO emails_sent(email, subject) VALUES(?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(recipient, subject)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	return nil
}

// insertEmailScanned inserts a new email record into the emails_scanned table
// Parameters:
//   - recipient: The recipient of the email
//
// It returns an error if the statement fails to prepare or execute.
func insertEmailScanned(recipient string) error {
	stmt, err := db.Prepare("INSERT INTO emails_scanned(email) VALUES(?)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(recipient)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	return nil
}

// insertCredential saves a new credential to the database
// Parameters:
//   - email: The email of the credential
//   - token: The token of the credential
//
// It returns an error if the statement fails to prepare or execute.
func insertCredential(email, token string) error {
	stmt, err := db.Prepare("INSERT INTO credentials(email, token) VALUES(?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(email, token)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	return nil
}

// insertConfig saves the current configuration to the database
// Parameters:
//   - config: The configuration to save
//
// It returns an error if the statement fails to prepare or execute.
func insertConfig(config models.ConfigData) error {
	// Prepare a SQL insert or update statement
	stmt, err := db.Prepare(`
		INSERT INTO configuration (key, value) VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET value=excluded.value;
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	configItems := map[string]string{
		"smtpHost":     config.SMTPConfig.Host,
		"smtpPort":     config.SMTPConfig.Port,
		"smtpUsername": config.SMTPConfig.Username,
		"smtpPassword": config.SMTPConfig.Password,
		"phishHost":    config.SquarePhishConfig.Host,
		"phishPort":    config.SquarePhishConfig.Port,
		"emailSender":  config.EmailConfig.Sender,
		"emailSubject": config.EmailConfig.Subject,
		"emailBody":    config.EmailConfig.Body,
	}

	for key, value := range configItems {
		_, err = stmt.Exec(key, value)
		if err != nil {
			return fmt.Errorf("failed to execute statement for %s: %w", key, err)
		}
	}

	return nil
}
