package notes

import (
	"crypto/sha1"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/sk25469/jot/config"
)

type Note struct {
	ID       string
	Title    string
	Tags     []string
	Mode     string
	Date     time.Time
	FilePath string
	Content  string
}

// CreateNote creates a new note with the given title and options
func CreateNote(title string, tags []string, mode string) (*Note, error) {
	if mode == "" {
		mode = config.AppConfig.DefaultMode
	}

	// Generate timestamp-based filename
	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05Z")
	slug := slugify(title)
	filename := fmt.Sprintf("%s-%s.md", timestamp, slug)

	notesDir := config.GetNotesDir()
	filePath := filepath.Join(notesDir, filename)

	// Create note metadata
	note := &Note{
		ID:       timestamp,
		Title:    title,
		Tags:     tags,
		Mode:     mode,
		Date:     time.Now().UTC(),
		FilePath: filePath,
	}

	// Create note content with metadata header
	content := generateNoteContent(note)

	// Write file
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return nil, err
	}

	// Open in editor
	if err := openInEditor(filePath); err != nil {
		return nil, err
	}

	return note, nil
}

// ListNotes returns all notes, optionally filtered by tag and mode
func ListNotes(tagFilter, modeFilter string) ([]*Note, error) {
	notesDir := config.GetNotesDir()

	files, err := filepath.Glob(filepath.Join(notesDir, "*.md"))
	if err != nil {
		return nil, err
	}

	var notes []*Note
	for _, file := range files {
		note, err := parseNoteFile(file)
		if err != nil {
			continue // Skip invalid files
		}

		// Apply filters
		if tagFilter != "" && !contains(note.Tags, tagFilter) {
			continue
		}
		if modeFilter != "" && note.Mode != modeFilter {
			continue
		}

		notes = append(notes, note)
	}

	// Sort by date (newest first)
	sort.Slice(notes, func(i, j int) bool {
		return notes[i].Date.After(notes[j].Date)
	})

	return notes, nil
}

// SearchNotes searches for notes by query string
func SearchNotes(query string) ([]*Note, error) {
	notes, err := ListNotes("", "")
	if err != nil {
		return nil, err
	}

	query = strings.ToLower(query)
	var matches []*Note

	for _, note := range notes {
		// Search in title
		if strings.Contains(strings.ToLower(note.Title), query) {
			matches = append(matches, note)
			continue
		}

		// Search in tags
		for _, tag := range note.Tags {
			if strings.Contains(strings.ToLower(tag), query) {
				matches = append(matches, note)
				break
			}
		}

		// Search in content
		if note.Content != "" && strings.Contains(strings.ToLower(note.Content), query) {
			matches = append(matches, note)
		}
	}

	return matches, nil
}

// OpenNote opens a note by ID or title
func OpenNote(identifier string) error {
	notes, err := ListNotes("", "")
	if err != nil {
		return err
	}

	var targetNote *Note

	// Try to find by exact ID first
	for _, note := range notes {
		if note.ID == identifier {
			targetNote = note
			break
		}
	}

	// If not found by exact ID, try partial ID match (like git)
	if targetNote == nil {
		matches := []*Note{}
		for _, note := range notes {
			if strings.HasPrefix(note.ID, identifier) {
				matches = append(matches, note)
			}
		}

		if len(matches) == 1 {
			targetNote = matches[0]
		} else if len(matches) > 1 {
			return fmt.Errorf("ambiguous ID '%s', could match: %s",
				identifier,
				func() string {
					var ids []string
					for _, n := range matches {
						ids = append(ids, n.ID)
					}
					return strings.Join(ids, ", ")
				}())
		}
	}

	// If still not found, try partial title match
	if targetNote == nil {
		for _, note := range notes {
			if strings.Contains(strings.ToLower(note.Title), strings.ToLower(identifier)) {
				targetNote = note
				break
			}
		}
	}

	if targetNote == nil {
		return fmt.Errorf("note not found: %s", identifier)
	}

	return openInEditor(targetNote.FilePath)
}

// GetStats returns basic statistics about notes
func GetStats() (map[string]interface{}, error) {
	notes, err := ListNotes("", "")
	if err != nil {
		return nil, err
	}

	// Count notes this week
	weekAgo := time.Now().AddDate(0, 0, -7)
	thisWeek := 0
	for _, note := range notes {
		if note.Date.After(weekAgo) {
			thisWeek++
		}
	}

	// Count tags
	tagCounts := make(map[string]int)
	for _, note := range notes {
		for _, tag := range note.Tags {
			tagCounts[tag]++
		}
	}

	stats := map[string]interface{}{
		"total_notes": len(notes),
		"this_week":   thisWeek,
		"tag_counts":  tagCounts,
	}

	return stats, nil
}

// Helper functions

func generateShortID(filename string) string {
	// Generate a short hash from the filename for a git-like ID
	h := sha1.New()
	h.Write([]byte(filename))
	hash := fmt.Sprintf("%x", h.Sum(nil))
	// Return first 7 characters like git
	return hash[:7]
}

func slugify(text string) string {
	// Convert to lowercase and replace spaces/special chars with hyphens
	reg := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	slug := reg.ReplaceAllString(strings.ToLower(text), "-")
	return strings.Trim(slug, "-")
}

func generateNoteContent(note *Note) string {
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
		note.Date.Format(time.RFC3339),
	)
}

func parseNoteFile(filePath string) (*Note, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	contentStr := string(content)

	// Generate short hash ID from filename
	filename := filepath.Base(filePath)
	id := generateShortID(filename)

	// Parse metadata from content
	lines := strings.Split(contentStr, "\n")
	note := &Note{
		ID:       id,
		FilePath: filePath,
		Content:  contentStr,
	}

	// Simple metadata parsing
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
					note.Date = date
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

	return note, nil
}

func openInEditor(filePath string) error {
	editor := config.AppConfig.Editor
	cmd := exec.Command(editor, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
