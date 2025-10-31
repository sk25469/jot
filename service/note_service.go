package service

import (
	"crypto/sha1"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/sahilsarwar/jot/config"
	"github.com/sahilsarwar/jot/database"
	"github.com/sahilsarwar/jot/models"
)

// NoteService handles business logic for notes
type NoteService struct {
	noteRepo  *database.NoteRepository
	statsRepo *database.StatsRepository
}

// NewNoteService creates a new note service
func NewNoteService(db *database.DB) *NoteService {
	return &NoteService{
		noteRepo:  database.NewNoteRepository(db),
		statsRepo: database.NewStatsRepository(db),
	}
}

// CreateNote creates a new note with the given title and options
func (s *NoteService) CreateNote(title string, tags []string, mode string) (*models.Note, error) {
	if mode == "" {
		mode = config.AppConfig.DefaultMode
	}

	// Generate timestamp-based filename
	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05Z")
	slug := slugify(title)
	filename := fmt.Sprintf("%s-%s.md", timestamp, slug)
	
	notesDir := config.GetNotesDir()
	filePath := filepath.Join(notesDir, filename)

	// Generate short hash ID from filename
	id := generateShortID(filename)

	// Create note model
	note := &models.Note{
		ID:        id,
		Title:     title,
		Mode:      mode,
		FilePath:  filePath,
		FileName:  filename,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Tags:      tags,
	}

	// Create note content with metadata header
	content := s.generateNoteContent(note)
	note.ContentHash = s.generateContentHash(content)
	note.ContentPreview = s.generatePreview(content)
	note.WordCount = s.countWords(content)
	
	// Write file
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return nil, fmt.Errorf("failed to write note file: %w", err)
	}

	// Save to database
	if err := s.noteRepo.Create(note); err != nil {
		// Clean up file if database save fails
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save note to database: %w", err)
	}

	// Update FTS index
	if err := s.updateFTSIndex(note, content); err != nil {
		// Log warning but don't fail - FTS is optional
		fmt.Printf("Warning: failed to update FTS index: %v\n", err)
	}

	// Open in editor
	if err := s.openInEditor(filePath); err != nil {
		return nil, fmt.Errorf("failed to open editor: %w", err)
	}

	return note, nil
}

// ListNotes returns notes with optional filtering
func (s *NoteService) ListNotes(tagFilter, modeFilter string) ([]*models.Note, error) {
	filter := models.DefaultListFilter()
	
	if tagFilter != "" {
		filter.Tags = []string{tagFilter}
	}
	if modeFilter != "" {
		filter.Mode = modeFilter
	}

	return s.noteRepo.List(filter)
}

// SearchNotes searches for notes by query string
func (s *NoteService) SearchNotes(query string) ([]*models.SearchResult, error) {
	return s.noteRepo.Search(query)
}

// OpenNote opens a note by ID or title
func (s *NoteService) OpenNote(identifier string) error {
	// Try to find by exact ID first
	note, err := s.noteRepo.GetByID(identifier)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}

	// If not found by exact ID, try partial ID match
	if note == nil {
		notes, err := s.ListNotes("", "")
		if err != nil {
			return fmt.Errorf("failed to list notes for partial search: %w", err)
		}

		var matches []*models.Note
		for _, n := range notes {
			if strings.HasPrefix(n.ID, identifier) {
				matches = append(matches, n)
			}
		}

		if len(matches) == 1 {
			note = matches[0]
		} else if len(matches) > 1 {
			var ids []string
			for _, n := range matches {
				ids = append(ids, n.ID)
			}
			return fmt.Errorf("ambiguous ID '%s', could match: %s", 
				identifier, strings.Join(ids, ", "))
		}
	}

	// If still not found, try partial title match
	if note == nil {
		notes, err := s.ListNotes("", "")
		if err != nil {
			return fmt.Errorf("failed to list notes for title search: %w", err)
		}

		for _, n := range notes {
			if strings.Contains(strings.ToLower(n.Title), strings.ToLower(identifier)) {
				note = n
				break
			}
		}
	}

	if note == nil {
		return fmt.Errorf("note not found: %s", identifier)
	}

	return s.openInEditor(note.FilePath)
}

// GetStats returns statistics about notes
func (s *NoteService) GetStats() (*models.StatsResult, error) {
	return s.statsRepo.GetStats()
}

// SyncFromFileSystem scans the notes directory and syncs with database
func (s *NoteService) SyncFromFileSystem() error {
	notesDir := config.GetNotesDir()
	
	files, err := filepath.Glob(filepath.Join(notesDir, "*.md"))
	if err != nil {
		return fmt.Errorf("failed to scan notes directory: %w", err)
	}

	for _, file := range files {
		if err := s.syncNoteFromFile(file); err != nil {
			// Log error but continue with other files
			fmt.Printf("Warning: failed to sync %s: %v\n", file, err)
		}
	}

	return nil
}

// Helper functions

func (s *NoteService) syncNoteFromFile(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	filename := filepath.Base(filePath)
	id := generateShortID(filename)

	// Check if note exists in database
	existingNote, err := s.noteRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to check existing note: %w", err)
	}

	// Parse note from file content
	note, err := s.parseNoteFile(filePath, string(content))
	if err != nil {
		return fmt.Errorf("failed to parse note file: %w", err)
	}

	if existingNote == nil {
		// Create new note in database
		if err := s.noteRepo.Create(note); err != nil {
			return err
		}
		// Update FTS index
		return s.updateFTSIndex(note, string(content))
	} else {
		// Update if content changed
		newHash := s.generateContentHash(string(content))
		if existingNote.ContentHash != newHash {
			note.ID = existingNote.ID // Preserve ID
			note.CreatedAt = existingNote.CreatedAt // Preserve creation time
			note.ContentHash = newHash
			if err := s.noteRepo.Update(note); err != nil {
				return err
			}
			// Update FTS index
			return s.updateFTSIndex(note, string(content))
		}
	}

	return nil
}

func (s *NoteService) parseNoteFile(filePath, content string) (*models.Note, error) {
	filename := filepath.Base(filePath)
	id := generateShortID(filename)
	
	note := &models.Note{
		ID:             id,
		FilePath:       filePath,
		FileName:       filename,
		ContentHash:    s.generateContentHash(content),
		ContentPreview: s.generatePreview(content),
		WordCount:      s.countWords(content),
		UpdatedAt:      time.Now().UTC(),
	}

	// Parse metadata from content
	lines := strings.Split(content, "\n")
	inMetadata := false
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "---" {
			if !inMetadata {
				inMetadata = true
				continue
			} else {
				break
			}
		}
		
		if inMetadata {
			if strings.HasPrefix(line, "title:") {
				note.Title = strings.TrimSpace(strings.TrimPrefix(line, "title:"))
				note.Title = strings.Trim(note.Title, "\"")
			} else if strings.HasPrefix(line, "mode:") {
				note.Mode = strings.TrimSpace(strings.TrimPrefix(line, "mode:"))
				note.Mode = strings.Trim(note.Mode, "\"")
			} else if strings.HasPrefix(line, "date:") {
				dateStr := strings.TrimSpace(strings.TrimPrefix(line, "date:"))
				dateStr = strings.Trim(dateStr, "\"")
				if date, err := time.Parse(time.RFC3339, dateStr); err == nil {
					note.CreatedAt = date
				}
			} else if strings.HasPrefix(line, "tags:") {
				tagsStr := strings.TrimSpace(strings.TrimPrefix(line, "tags:"))
				tagsStr = strings.Trim(tagsStr, "[]")
				if tagsStr != "" {
					tags := strings.Split(tagsStr, ",")
					for i, tag := range tags {
						tags[i] = strings.TrimSpace(tag)
					}
					note.Tags = tags
				}
			}
		}
	}

	// Set defaults if not found
	if note.Title == "" {
		note.Title = "Untitled"
	}
	if note.Mode == "" {
		note.Mode = config.AppConfig.DefaultMode
	}
	if note.CreatedAt.IsZero() {
		note.CreatedAt = time.Now().UTC()
	}

	return note, nil
}

func generateShortID(filename string) string {
	h := sha1.New()
	h.Write([]byte(filename))
	hash := fmt.Sprintf("%x", h.Sum(nil))
	return hash[:7]
}

func slugify(text string) string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	slug := reg.ReplaceAllString(strings.ToLower(text), "-")
	return strings.Trim(slug, "-")
}

func (s *NoteService) generateNoteContent(note *models.Note) string {
	tagsStr := ""
	if len(note.Tags) > 0 {
		tagsStr = fmt.Sprintf("[%s]", strings.Join(note.Tags, ", "))
	} else {
		tagsStr = "[]"
	}

	return fmt.Sprintf(`---
title: %s
tags: %s
mode: %s
date: %s
---

`,
		note.Title,
		tagsStr,
		note.Mode,
		note.CreatedAt.Format(time.RFC3339),
	)
}

func (s *NoteService) generateContentHash(content string) string {
	h := sha1.New()
	h.Write([]byte(content))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (s *NoteService) generatePreview(content string) string {
	// Remove YAML frontmatter
	lines := strings.Split(content, "\n")
	inMetadata := false
	contentStart := 0
	
	for i, line := range lines {
		if strings.TrimSpace(line) == "---" {
			if !inMetadata {
				inMetadata = true
			} else {
				contentStart = i + 1
				break
			}
		}
	}

	if contentStart < len(lines) {
		preview := strings.Join(lines[contentStart:], "\n")
		preview = strings.TrimSpace(preview)
		if len(preview) > 200 {
			return preview[:200]
		}
		return preview
	}

	return ""
}

func (s *NoteService) countWords(content string) int {
	// Simple word count - split by whitespace
	words := strings.Fields(content)
	return len(words)
}

func (s *NoteService) openInEditor(filePath string) error {
	// This is the same as the original openInEditor function
	editor := config.AppConfig.Editor
	cmd := exec.Command(editor, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// updateFTSIndex updates the full-text search index for a note
func (s *NoteService) updateFTSIndex(note *models.Note, content string) error {
	// Get database connection
	db := s.noteRepo.GetDB()
	
	// Prepare tags string for FTS
	tagsStr := strings.Join(note.Tags, " ")
	
	// Insert or replace in FTS table
	query := `INSERT OR REPLACE INTO notes_fts (note_id, title, content, tags) VALUES (?, ?, ?, ?)`
	_, err := db.Connection().Exec(query, note.ID, note.Title, content, tagsStr)
	
	return err
}

// deleteFTSIndex removes a note from the FTS index
func (s *NoteService) deleteFTSIndex(noteID string) error {
	// Get database connection  
	db := s.noteRepo.GetDB()
	
	query := `DELETE FROM notes_fts WHERE note_id = ?`
	_, err := db.Connection().Exec(query, noteID)
	
	return err
}