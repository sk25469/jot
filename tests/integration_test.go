package main

import (
	"os"
	"testing"
	"time"

	"github.com/sk25469/jot/models"
)

// Integration tests for the jot CLI
// These tests verify that components work together correctly

func TestModelsIntegration(t *testing.T) {
	// Test that we can create and manipulate all model types
	now := time.Now()

	// Create a note
	note := &models.Note{
		ID:             "abc1234",
		Title:          "Integration Test Note",
		Mode:           "dev",
		FilePath:       "/tmp/test.md",
		FileName:       "test.md",
		ContentHash:    "testhash",
		CreatedAt:      now,
		UpdatedAt:      now,
		ContentPreview: "This is a test note for integration testing",
		WordCount:      9,
		Tags:           []string{"integration", "test"},
	}

	// Verify note fields
	if note.ID != "abc1234" {
		t.Errorf("Note ID mismatch")
	}

	if len(note.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(note.Tags))
	}

	// Create a search result
	searchResult := &models.SearchResult{
		Note:      *note,
		Rank:      1.5,
		Snippet:   "integration testing",
		MatchType: "content",
	}

	// Verify search result
	if searchResult.Rank != 1.5 {
		t.Errorf("Search result rank mismatch")
	}

	if searchResult.Note.Title != note.Title {
		t.Errorf("Embedded note data mismatch")
	}

	// Test list filter
	filter := models.DefaultListFilter()
	filter.Tags = []string{"integration"}
	filter.Mode = "dev"

	if filter.Limit != 100 {
		t.Errorf("Default limit should be 100")
	}

	if len(filter.Tags) != 1 || filter.Tags[0] != "integration" {
		t.Errorf("Filter tags not set correctly")
	}

	// Test stats result
	stats := &models.StatsResult{
		TotalNotes:    1,
		NotesThisWeek: 1,
		TagCounts: map[string]int{
			"integration": 1,
			"test":        1,
		},
		ModeStats: map[string]int{
			"dev": 1,
		},
		WordCount:    9,
		CreatedToday: 1,
	}

	if stats.TotalNotes != 1 {
		t.Errorf("Stats total notes mismatch")
	}

	if stats.TagCounts["integration"] != 1 {
		t.Errorf("Tag count mismatch")
	}
}

func TestEnvironmentSetup(t *testing.T) {
	// Test that we can get environment information needed for the app

	// Test home directory access
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Errorf("Should be able to get user home directory: %v", err)
	}

	if homeDir == "" {
		t.Errorf("Home directory should not be empty")
	}

	// Test EDITOR environment variable
	editor := os.Getenv("EDITOR")
	// Editor may or may not be set, but we should be able to read it
	_ = editor // Just testing we can access it

	// Test that we can create temporary files (needed for notes)
	tempFile, err := os.CreateTemp("", "jot-test-*.md")
	if err != nil {
		t.Errorf("Should be able to create temporary files: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Test writing to file
	testContent := "# Test Note\n\nThis is a test."
	_, err = tempFile.WriteString(testContent)
	if err != nil {
		t.Errorf("Should be able to write to files: %v", err)
	}
}

func TestTimeHandling(t *testing.T) {
	// Test time operations that the app relies on

	now := time.Now()

	// Test RFC3339 formatting (used in note metadata)
	formatted := now.Format(time.RFC3339)
	if formatted == "" {
		t.Errorf("RFC3339 formatting should not be empty")
	}

	// Test parsing RFC3339
	parsed, err := time.Parse(time.RFC3339, formatted)
	if err != nil {
		t.Errorf("Should be able to parse RFC3339 time: %v", err)
	}

	// Should be approximately the same (within a second)
	if parsed.Sub(now).Abs() > time.Second {
		t.Errorf("Parsed time should be close to original")
	}

	// Test date formatting (used in display)
	dateFormatted := now.Format("2006-01-02")
	if len(dateFormatted) != 10 {
		t.Errorf("Date format should be 10 characters, got %d", len(dateFormatted))
	}

	// Test timestamp formatting (used in filenames)
	timestampFormatted := now.Format("2006-01-02T15-04-05Z")
	if timestampFormatted == "" {
		t.Errorf("Timestamp formatting should not be empty")
	}
}

func TestErrorHandling(t *testing.T) {
	// Test common error scenarios that the app should handle gracefully

	// Test reading non-existent file
	_, err := os.ReadFile("/non/existent/file.md")
	if err == nil {
		t.Errorf("Should get error when reading non-existent file")
	}

	// Test creating file in non-existent directory
	err = os.WriteFile("/non/existent/dir/file.md", []byte("test"), 0644)
	if err == nil {
		t.Errorf("Should get error when writing to non-existent directory")
	}

	// Test that we can check if errors are specific types
	if os.IsNotExist(err) {
		// This is the expected error type
	} else {
		t.Logf("Got different error type: %T", err)
	}
}
