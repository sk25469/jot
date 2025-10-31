package models

import (
	"testing"
	"time"
)

func TestDefaultListFilter(t *testing.T) {
	filter := DefaultListFilter()

	// Test default values
	if filter.Limit != 100 {
		t.Errorf("Expected Limit to be 100, got %d", filter.Limit)
	}

	if filter.Offset != 0 {
		t.Errorf("Expected Offset to be 0, got %d", filter.Offset)
	}

	if filter.SortBy != "created" {
		t.Errorf("Expected SortBy to be 'created', got %s", filter.SortBy)
	}

	if filter.SortOrder != "desc" {
		t.Errorf("Expected SortOrder to be 'desc', got %s", filter.SortOrder)
	}

	// Test nil values
	if filter.Since != nil {
		t.Errorf("Expected Since to be nil, got %v", filter.Since)
	}

	if filter.Until != nil {
		t.Errorf("Expected Until to be nil, got %v", filter.Until)
	}

	if filter.Mode != "" {
		t.Errorf("Expected Mode to be empty, got %s", filter.Mode)
	}

	if len(filter.Tags) != 0 {
		t.Errorf("Expected Tags to be empty, got %v", filter.Tags)
	}
}

func TestNoteStruct(t *testing.T) {
	now := time.Now()

	note := &Note{
		ID:             "abc1234",
		Title:          "Test Note",
		Mode:           "dev",
		FilePath:       "/path/to/note.md",
		FileName:       "note.md",
		ContentHash:    "hashvalue",
		CreatedAt:      now,
		UpdatedAt:      now,
		ContentPreview: "This is a test note",
		WordCount:      5,
		Tags:           []string{"test", "golang"},
	}

	// Test that all fields are set correctly
	if note.ID != "abc1234" {
		t.Errorf("Expected ID to be 'abc1234', got %s", note.ID)
	}

	if note.Title != "Test Note" {
		t.Errorf("Expected Title to be 'Test Note', got %s", note.Title)
	}

	if note.Mode != "dev" {
		t.Errorf("Expected Mode to be 'dev', got %s", note.Mode)
	}

	if note.WordCount != 5 {
		t.Errorf("Expected WordCount to be 5, got %d", note.WordCount)
	}

	if len(note.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(note.Tags))
	}

	if note.Tags[0] != "test" || note.Tags[1] != "golang" {
		t.Errorf("Expected tags [test, golang], got %v", note.Tags)
	}
}

func TestSearchResultStruct(t *testing.T) {
	now := time.Now()

	searchResult := &SearchResult{
		Note: Note{
			ID:        "abc1234",
			Title:     "Test Note",
			Mode:      "dev",
			CreatedAt: now,
			Tags:      []string{"test"},
		},
		Rank:      1.5,
		Snippet:   "This is a test snippet",
		MatchType: "title",
	}

	// Test SearchResult specific fields
	if searchResult.Rank != 1.5 {
		t.Errorf("Expected Rank to be 1.5, got %f", searchResult.Rank)
	}

	if searchResult.Snippet != "This is a test snippet" {
		t.Errorf("Expected Snippet to be 'This is a test snippet', got %s", searchResult.Snippet)
	}

	if searchResult.MatchType != "title" {
		t.Errorf("Expected MatchType to be 'title', got %s", searchResult.MatchType)
	}

	// Test inherited Note fields
	if searchResult.Note.ID != "abc1234" {
		t.Errorf("Expected inherited ID to be 'abc1234', got %s", searchResult.Note.ID)
	}
}

func TestStatsResultStruct(t *testing.T) {
	tagCounts := map[string]int{
		"golang": 5,
		"test":   3,
	}

	modeStats := map[string]int{
		"dev":     10,
		"journal": 2,
	}

	stats := &StatsResult{
		TotalNotes:    12,
		NotesThisWeek: 3,
		TagCounts:     tagCounts,
		ModeStats:     modeStats,
		WordCount:     1500,
		CreatedToday:  1,
	}

	// Test all fields
	if stats.TotalNotes != 12 {
		t.Errorf("Expected TotalNotes to be 12, got %d", stats.TotalNotes)
	}

	if stats.NotesThisWeek != 3 {
		t.Errorf("Expected NotesThisWeek to be 3, got %d", stats.NotesThisWeek)
	}

	if stats.WordCount != 1500 {
		t.Errorf("Expected WordCount to be 1500, got %d", stats.WordCount)
	}

	if stats.CreatedToday != 1 {
		t.Errorf("Expected CreatedToday to be 1, got %d", stats.CreatedToday)
	}

	// Test maps
	if stats.TagCounts["golang"] != 5 {
		t.Errorf("Expected golang tag count to be 5, got %d", stats.TagCounts["golang"])
	}

	if stats.ModeStats["dev"] != 10 {
		t.Errorf("Expected dev mode count to be 10, got %d", stats.ModeStats["dev"])
	}
}

func TestTagStruct(t *testing.T) {
	now := time.Now()

	tag := &Tag{
		ID:         1,
		Name:       "golang",
		CreatedAt:  now,
		UsageCount: 5,
	}

	if tag.ID != 1 {
		t.Errorf("Expected ID to be 1, got %d", tag.ID)
	}

	if tag.Name != "golang" {
		t.Errorf("Expected Name to be 'golang', got %s", tag.Name)
	}

	if tag.UsageCount != 5 {
		t.Errorf("Expected UsageCount to be 5, got %d", tag.UsageCount)
	}
}

func TestNoteTagStruct(t *testing.T) {
	noteTag := &NoteTag{
		NoteID: "abc1234",
		TagID:  1,
	}

	if noteTag.NoteID != "abc1234" {
		t.Errorf("Expected NoteID to be 'abc1234', got %s", noteTag.NoteID)
	}

	if noteTag.TagID != 1 {
		t.Errorf("Expected TagID to be 1, got %d", noteTag.TagID)
	}
}

func TestConfigStruct(t *testing.T) {
	now := time.Now()

	config := &Config{
		Key:       "editor",
		Value:     "vim",
		UpdatedAt: now,
	}

	if config.Key != "editor" {
		t.Errorf("Expected Key to be 'editor', got %s", config.Key)
	}

	if config.Value != "vim" {
		t.Errorf("Expected Value to be 'vim', got %s", config.Value)
	}

	if !config.UpdatedAt.Equal(now) {
		t.Errorf("Expected UpdatedAt to be %v, got %v", now, config.UpdatedAt)
	}
}
