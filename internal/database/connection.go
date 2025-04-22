// filepath: /internal/database/connection.go
package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

const dbFile = "./squarephish.db"

// connect opens a connection to the SQLite database.
// If the database file does not exist, it will be created.
//
// It returns an error if any.
func connect() error {
	var err error

	db, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

// close closes the connection to the SQLite database.
//
// It returns an error if the database fails to close.
func close() error {
	if db == nil {
		return nil
	}

	if err := db.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}

	return nil
}
