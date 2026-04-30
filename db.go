package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// DB wraps a sql.DB connection to the counter database.
var db *sql.DB

// initDB opens (or creates) the SQLite database at the given path
// and ensures the counter table exists with an initial row.
func initDB(dbPath string) error {
	var err error
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}

	// Enable WAL mode for better concurrent read/write performance.
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return fmt.Errorf("set WAL mode: %w", err)
	}

	// Create the counter table if it doesn't exist.
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS counter (
			id INTEGER PRIMARY KEY CHECK (id = 1),
			value INTEGER NOT NULL DEFAULT 0
		)
	`); err != nil {
		return fmt.Errorf("create table: %w", err)
	}

	// Ensure the single counter row exists.
	if _, err := db.Exec(`
		INSERT OR IGNORE INTO counter (id, value) VALUES (1, 0)
	`); err != nil {
		return fmt.Errorf("insert initial row: %w", err)
	}

	return nil
}

// getCounter returns the current counter value.
func getCounter() (int, error) {
	var value int
	err := db.QueryRow("SELECT value FROM counter WHERE id = 1").Scan(&value)
	if err != nil {
		return 0, fmt.Errorf("get counter: %w", err)
	}
	return value, nil
}

// incrementCounter atomically increments the counter and returns the new value.
func incrementCounter() (int, error) {
	var value int
	err := db.QueryRow(`
		UPDATE counter SET value = value + 1 WHERE id = 1
		RETURNING value
	`).Scan(&value)
	if err != nil {
		return 0, fmt.Errorf("increment counter: %w", err)
	}
	return value, nil
}

// closeDB closes the database connection.
func closeDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
