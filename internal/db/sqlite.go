package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

// DB represents the database connection
type DB struct {
	*sql.DB
}

// New creates a new database connection
func New(dbPath string) (*DB, error) {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if dir != "." && dir != "/" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create database directory: %w", err)
		}
	}

	// Open database connection
	sqlDB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Enable foreign keys
	if _, err := sqlDB.Exec("PRAGMA foreign_keys = ON"); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	db := &DB{sqlDB}

	// Initialize database schema
	if err := db.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return db, nil
}

// initSchema creates the database schema if it doesn't exist
func (db *DB) initSchema() error {
	log.Info().Msg("Initializing database schema")

	// Create users table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			role TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// Create api_keys table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS api_keys (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			key TEXT NOT NULL UNIQUE,
			created_at TIMESTAMP NOT NULL,
			expires_at TIMESTAMP NOT NULL,
			last_used TIMESTAMP NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create api_keys table: %w", err)
	}

	// Create events table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS events (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT,
			start_time TIMESTAMP NOT NULL,
			end_time TIMESTAMP NOT NULL,
			rewards TEXT,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create events table: %w", err)
	}

	// Create index on event start and end times
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_events_time 
		ON events(start_time, end_time)
	`)
	if err != nil {
		return fmt.Errorf("failed to create events index: %w", err)
	}

	// Check if admin user exists, create if not
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'admin'").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check for admin user: %w", err)
	}

	if count == 0 {
		log.Info().Msg("Creating default admin user")
		_, err = db.Exec(`
			INSERT INTO users (id, username, role, created_at)
			VALUES (?, ?, ?, datetime('now'))
		`, "00000000-0000-0000-0000-000000000000", "admin", "admin")
		if err != nil {
			return fmt.Errorf("failed to create admin user: %w", err)
		}

		// Create an API key for the admin user
		_, err = db.Exec(`
			INSERT INTO api_keys (id, user_id, key, created_at, expires_at, last_used)
			VALUES (?, ?, ?, datetime('now'), datetime('now', '+365 days'), datetime('now'))
		`, "00000000-0000-0000-0000-000000000001", "00000000-0000-0000-0000-000000000000", "admin-api-key-00000000-0000-0000-0000-000000000000")
		if err != nil {
			return fmt.Errorf("failed to create admin API key: %w", err)
		}
	}

	return nil
}
