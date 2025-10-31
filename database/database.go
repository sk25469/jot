package database

import (
	"database/sql"
	"embed"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schemaFS embed.FS

// DB represents the application database
type DB struct {
	conn *sql.DB
	path string
}

// Config holds database configuration
type Config struct {
	Path string
}

// New creates a new database connection
func New(config Config) (*DB, error) {
	// Ensure database directory exists
	dbDir := filepath.Dir(config.Path)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database connection
	conn, err := sql.Open("sqlite", config.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure SQLite for optimal performance
	if err := configureSQLite(conn); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to configure SQLite: %w", err)
	}

	db := &DB{
		conn: conn,
		path: config.Path,
	}

	// Initialize schema if needed
	if err := db.initializeSchema(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return db, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}

// Connection returns the underlying sql.DB connection
func (db *DB) Connection() *sql.DB {
	return db.conn
}

// configureSQLite sets up SQLite for optimal performance
func configureSQLite(conn *sql.DB) error {
	pragmas := []string{
		"PRAGMA foreign_keys = ON",
		"PRAGMA journal_mode = WAL",
		"PRAGMA synchronous = NORMAL",
		"PRAGMA cache_size = -64000", // 64MB cache
		"PRAGMA temp_store = MEMORY",
		"PRAGMA mmap_size = 268435456", // 256MB mmap
	}

	for _, pragma := range pragmas {
		if _, err := conn.Exec(pragma); err != nil {
			return fmt.Errorf("failed to execute pragma %s: %w", pragma, err)
		}
	}

	return nil
}

// initializeSchema creates tables if they don't exist
func (db *DB) initializeSchema() error {
	// Check if database is already initialized
	var count int
	err := db.conn.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='notes'").Scan(&count)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check if schema exists: %w", err)
	}

	// If tables exist, check version and potentially migrate
	if count > 0 {
		return db.checkAndMigrate()
	}

	// Read schema from embedded file
	schemaBytes, err := schemaFS.ReadFile("schema.sql")
	if err != nil {
		return fmt.Errorf("failed to read embedded schema: %w", err)
	}

	// Execute schema
	if _, err := db.conn.Exec(string(schemaBytes)); err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	return nil
}

// checkAndMigrate checks the database version and runs migrations if needed
func (db *DB) checkAndMigrate() error {
	var version string
	err := db.conn.QueryRow("SELECT value FROM config WHERE key = 'db_version'").Scan(&version)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to get database version: %w", err)
	}

	// For now, we assume version 1.0 is current
	// Future migrations would go here
	if version == "" {
		// Insert version if missing
		_, err = db.conn.Exec("INSERT OR REPLACE INTO config (key, value) VALUES ('db_version', '1.0')")
		if err != nil {
			return fmt.Errorf("failed to set database version: %w", err)
		}
	}

	return nil
}
