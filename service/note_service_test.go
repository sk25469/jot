package service

import (
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/sk25469/jot/models"
)

// Test helper functions that don't require database

func TestGenerateNoteContent(t *testing.T) {
	// Create a test service (we only need the method, not the DB)
	service := &NoteService{}

	now := time.Date(2023, 10, 15, 10, 30, 0, 0, time.UTC)

	testCases := []struct {
		name string
		note *models.Note
	}{
		{
			name: "note with tags",
			note: &models.Note{
				Title:     "Test Note",
				Tags:      []string{"golang", "test"},
				Mode:      "dev",
				CreatedAt: now,
			},
		},
		{
			name: "note without tags",
			note: &models.Note{
				Title:     "Simple Note",
				Tags:      []string{},
				Mode:      "journal",
				CreatedAt: now,
			},
		},
		{
			name: "note with special characters in title",
			note: &models.Note{
				Title:     "Note: With \"Special\" Characters & Symbols",
				Tags:      []string{"special"},
				Mode:      "dev",
				CreatedAt: now,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			content := service.generateNoteContent(tc.note)

			// Should contain frontmatter markers
			if !strings.Contains(content, "---") {
				t.Errorf("Content should contain frontmatter markers, got: %s", content)
			}

			// Should contain title
			if !strings.Contains(content, tc.note.Title) {
				t.Errorf("Content should contain title %q, got: %s", tc.note.Title, content)
			}

			// Should contain mode
			if !strings.Contains(content, tc.note.Mode) {
				t.Errorf("Content should contain mode %q, got: %s", tc.note.Mode, content)
			}

			// Should contain date
			expectedDate := tc.note.CreatedAt.Format(time.RFC3339)
			if !strings.Contains(content, expectedDate) {
				t.Errorf("Content should contain date %q, got: %s", expectedDate, content)
			}

			// Check tags handling
			if len(tc.note.Tags) > 0 {
				for _, tag := range tc.note.Tags {
					if !strings.Contains(content, tag) {
						t.Errorf("Content should contain tag %q, got: %s", tag, content)
					}
				}
			} else {
				// Should contain empty tags array
				if !strings.Contains(content, "[]") {
					t.Errorf("Content should contain empty tags array, got: %s", content)
				}
			}
		})
	}
}

func TestGenerateContentHash(t *testing.T) {
	service := &NoteService{}

	testCases := []struct {
		content  string
		expected int // expected length
	}{
		{"Hello, World!", 40}, // SHA1 hex is 40 characters
		{"", 40},
		{"Multi\nline\ncontent", 40},
	}

	for _, tc := range testCases {
		t.Run(tc.content, func(t *testing.T) {
			hash := service.generateContentHash(tc.content)

			if len(hash) != tc.expected {
				t.Errorf("Hash length = %d, expected %d", len(hash), tc.expected)
			}

			// Should be valid hex
			matched, err := regexp.MatchString("^[a-f0-9]+$", hash)
			if err != nil {
				t.Errorf("Error checking hex format: %v", err)
			}
			if !matched {
				t.Errorf("Hash should be hex format, got: %s", hash)
			}
		})
	}

	// Test consistency
	content := "test content"
	hash1 := service.generateContentHash(content)
	hash2 := service.generateContentHash(content)

	if hash1 != hash2 {
		t.Errorf("Hash should be deterministic, got %q and %q", hash1, hash2)
	}

	// Test different content produces different hashes
	hash3 := service.generateContentHash("different content")
	if hash1 == hash3 {
		t.Errorf("Different content should produce different hashes")
	}
}

func TestGeneratePreview(t *testing.T) {
	service := &NoteService{}

	testCases := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "content without frontmatter",
			content:  "This is a simple note content.",
			expected: "This is a simple note content.",
		},
		{
			name: "content with frontmatter",
			content: `---
title: Test Note
tags: [golang]
mode: dev
---

This is the actual content after frontmatter.`,
			expected: "This is the actual content after frontmatter.",
		},
		{
			name: "long content",
			content: `---
title: Long Note
---

` + strings.Repeat("This is a very long content that should be truncated. ", 10),
			expected: strings.Repeat("This is a very long content that should be truncated. ", 10)[:200],
		},
		{
			name:     "empty content",
			content:  "",
			expected: "",
		},
		{
			name: "only frontmatter",
			content: `---
title: Test
---`,
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			preview := service.generatePreview(tc.content)

			if tc.name == "long content" {
				if len(preview) > 200 {
					t.Errorf("Preview should be truncated to 200 chars, got %d", len(preview))
				}
			} else {
				if preview != tc.expected {
					t.Errorf("Preview = %q, expected %q", preview, tc.expected)
				}
			}

			// Preview should never contain frontmatter markers
			if strings.Contains(preview, "---") {
				t.Errorf("Preview should not contain frontmatter markers, got: %s", preview)
			}
		})
	}
}

func TestCountWords(t *testing.T) {
	service := &NoteService{}

	testCases := []struct {
		content  string
		expected int
	}{
		{"Hello world", 2},
		{"", 0},
		{"One", 1},
		{"Multiple    spaces   between    words", 4},
		{"Line1\nLine2\nLine3", 3},
		{"Word1,Word2;Word3.Word4", 1}, // Word count uses simple whitespace splitting
		{"   Leading and trailing spaces   ", 4},
		{"\n\n\nNewlines\nonly\n\n", 2},
	}

	for _, tc := range testCases {
		t.Run(tc.content, func(t *testing.T) {
			count := service.countWords(tc.content)
			if count != tc.expected {
				t.Errorf("countWords(%q) = %d, expected %d", tc.content, count, tc.expected)
			}
		})
	}
}
