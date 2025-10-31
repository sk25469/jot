package app

import (
	"fmt"
	"path/filepath"

	"github.com/sk25469/jot/config"
	"github.com/sk25469/jot/database"
	"github.com/sk25469/jot/service"
)

// App holds the application dependencies
type App struct {
	DB          *database.DB
	NoteService *service.NoteService
}

// Global app instance
var Instance *App

// Initialize sets up the application with database and services
func Initialize() error {
	// Initialize config first
	if err := config.InitConfig(); err != nil {
		return fmt.Errorf("failed to initialize config: %w", err)
	}

	// Initialize database
	dbPath := filepath.Join(config.GetJotDir(), "jot.db")
	db, err := database.New(database.Config{
		Path: dbPath,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize services
	noteService := service.NewNoteService(db)

	// Sync existing notes from filesystem to database
	if err := noteService.SyncFromFileSystem(); err != nil {
		return fmt.Errorf("failed to sync notes from filesystem: %w", err)
	}

	// Set global instance
	Instance = &App{
		DB:          db,
		NoteService: noteService,
	}

	return nil
}

// Cleanup closes database connections and performs cleanup
func Cleanup() error {
	if Instance != nil && Instance.DB != nil {
		return Instance.DB.Close()
	}
	return nil
}

// GetJotDir returns the jot directory path (same as config package function)
func GetJotDir() string {
	return config.GetJotDir()
}
