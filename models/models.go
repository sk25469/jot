package models

import (
	"time"
)

// Note represents a note in the database
type Note struct {
	ID             string    `db:"id" json:"id"`
	Title          string    `db:"title" json:"title"`
	Mode           string    `db:"mode" json:"mode"`
	FilePath       string    `db:"file_path" json:"file_path"`
	FileName       string    `db:"file_name" json:"file_name"`
	ContentHash    string    `db:"content_hash" json:"content_hash"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
	ContentPreview string    `db:"content_preview" json:"content_preview"`
	WordCount      int       `db:"word_count" json:"word_count"`
	Tags           []string  `json:"tags"` // Populated by joins, not stored directly
}

// Tag represents a tag in the database
type Tag struct {
	ID         int       `db:"id" json:"id"`
	Name       string    `db:"name" json:"name"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UsageCount int       `db:"usage_count" json:"usage_count"`
}

// NoteTag represents the many-to-many relationship between notes and tags
type NoteTag struct {
	NoteID string `db:"note_id" json:"note_id"`
	TagID  int    `db:"tag_id" json:"tag_id"`
}

// Config represents a configuration setting
type Config struct {
	Key       string    `db:"key" json:"key"`
	Value     string    `db:"value" json:"value"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// SearchResult represents a full-text search result
type SearchResult struct {
	Note
	Rank      float64 `json:"rank"`       // Search relevance rank
	Snippet   string  `json:"snippet"`    // Highlighted content snippet
	MatchType string  `json:"match_type"` // "title", "content", "tags"
}

// ListFilter represents filtering options for listing notes
type ListFilter struct {
	Tags      []string
	Mode      string
	Since     *time.Time
	Until     *time.Time
	Limit     int
	Offset    int
	SortBy    string // "created", "updated", "title"
	SortOrder string // "asc", "desc"
}

// DefaultListFilter returns a filter with sensible defaults
func DefaultListFilter() ListFilter {
	return ListFilter{
		Limit:     100,
		Offset:    0,
		SortBy:    "created",
		SortOrder: "desc",
	}
}

// StatsResult represents statistics about notes
type StatsResult struct {
	TotalNotes    int            `json:"total_notes"`
	NotesThisWeek int            `json:"notes_this_week"`
	TagCounts     map[string]int `json:"tag_counts"`
	ModeStats     map[string]int `json:"mode_stats"`
	WordCount     int            `json:"word_count"`
	CreatedToday  int            `json:"created_today"`
}
