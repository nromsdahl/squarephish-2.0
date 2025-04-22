package database

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// migrate initializes the database schema.
// It creates the necessary tables if they do not exist.
//
// It returns an error if any of the queries fail.
func migrate() error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS configuration  (key TEXT PRIMARY KEY, value TEXT);`,
		`CREATE TABLE IF NOT EXISTS emails_sent    (id INTEGER PRIMARY KEY AUTOINCREMENT, timestamp DATETIME DEFAULT CURRENT_TIMESTAMP, email TEXT, subject TEXT);`,
		`CREATE TABLE IF NOT EXISTS emails_scanned (id INTEGER PRIMARY KEY AUTOINCREMENT, timestamp DATETIME DEFAULT CURRENT_TIMESTAMP, email TEXT);`,
		`CREATE TABLE IF NOT EXISTS credentials    (id INTEGER PRIMARY KEY AUTOINCREMENT, timestamp DATETIME DEFAULT CURRENT_TIMESTAMP, email TEXT, token TEXT);`,
	}

	for _, query := range migrations {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to initialize database: %w", err)
		}
	}

	return nil
}
