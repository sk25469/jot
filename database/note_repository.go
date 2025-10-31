package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/sk25469/jot/models"
)

// NoteRepository handles database operations for notes
type NoteRepository struct {
	db *DB
}

// NewNoteRepository creates a new note repository
func NewNoteRepository(db *DB) *NoteRepository {
	return &NoteRepository{db: db}
}

// GetDB returns the underlying database connection
func (r *NoteRepository) GetDB() *DB {
	return r.db
}

// Create creates a new note in the database
func (r *NoteRepository) Create(note *models.Note) error {
	tx, err := r.db.conn.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert note
	query := `
		INSERT INTO notes (id, title, mode, file_path, file_name, content_hash, 
			created_at, updated_at, content_preview, word_count)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = tx.Exec(query,
		note.ID, note.Title, note.Mode, note.FilePath, note.FileName,
		note.ContentHash, note.CreatedAt, note.UpdatedAt,
		note.ContentPreview, note.WordCount)
	if err != nil {
		return fmt.Errorf("failed to insert note: %w", err)
	}

	// Insert tags if any
	if len(note.Tags) > 0 {
		if err := r.insertTagsForNote(tx, note.ID, note.Tags); err != nil {
			return fmt.Errorf("failed to insert tags: %w", err)
		}
	}

	return tx.Commit()
}

// Update updates an existing note
func (r *NoteRepository) Update(note *models.Note) error {
	tx, err := r.db.conn.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update note
	query := `
		UPDATE notes 
		SET title = ?, mode = ?, content_hash = ?, updated_at = ?, 
			content_preview = ?, word_count = ?
		WHERE id = ?`

	_, err = tx.Exec(query,
		note.Title, note.Mode, note.ContentHash, note.UpdatedAt,
		note.ContentPreview, note.WordCount, note.ID)
	if err != nil {
		return fmt.Errorf("failed to update note: %w", err)
	}

	// Update tags - delete existing and insert new ones
	if _, err := tx.Exec("DELETE FROM note_tags WHERE note_id = ?", note.ID); err != nil {
		return fmt.Errorf("failed to delete existing tags: %w", err)
	}

	if len(note.Tags) > 0 {
		if err := r.insertTagsForNote(tx, note.ID, note.Tags); err != nil {
			return fmt.Errorf("failed to insert updated tags: %w", err)
		}
	}

	return tx.Commit()
}

// GetByID retrieves a note by its ID
func (r *NoteRepository) GetByID(id string) (*models.Note, error) {
	query := `
		SELECT n.id, n.title, n.mode, n.file_path, n.file_name, n.content_hash,
			n.created_at, n.updated_at, n.content_preview, n.word_count,
			COALESCE(GROUP_CONCAT(t.name, ','), '') as tags
		FROM notes n
		LEFT JOIN note_tags nt ON n.id = nt.note_id
		LEFT JOIN tags t ON nt.tag_id = t.id
		WHERE n.id = ?
		GROUP BY n.id`

	row := r.db.conn.QueryRow(query, id)

	note := &models.Note{}
	var tagsStr string

	err := row.Scan(
		&note.ID, &note.Title, &note.Mode, &note.FilePath, &note.FileName,
		&note.ContentHash, &note.CreatedAt, &note.UpdatedAt,
		&note.ContentPreview, &note.WordCount, &tagsStr)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get note by ID: %w", err)
	}

	// Parse tags
	if tagsStr != "" {
		note.Tags = strings.Split(tagsStr, ",")
	}

	return note, nil
}

// List retrieves notes with optional filtering
func (r *NoteRepository) List(filter models.ListFilter) ([]*models.Note, error) {
	query := `
		SELECT n.id, n.title, n.mode, n.file_path, n.file_name, n.content_hash,
			n.created_at, n.updated_at, n.content_preview, n.word_count,
			COALESCE(GROUP_CONCAT(t.name, ','), '') as tags
		FROM notes n
		LEFT JOIN note_tags nt ON n.id = nt.note_id
		LEFT JOIN tags t ON nt.tag_id = t.id`

	var conditions []string
	var args []interface{}

	// Build WHERE clause
	if filter.Mode != "" {
		conditions = append(conditions, "n.mode = ?")
		args = append(args, filter.Mode)
	}

	if filter.Since != nil {
		conditions = append(conditions, "n.created_at >= ?")
		args = append(args, filter.Since)
	}

	if filter.Until != nil {
		conditions = append(conditions, "n.created_at <= ?")
		args = append(args, filter.Until)
	}

	// Handle tag filtering
	if len(filter.Tags) > 0 {
		tagPlaceholders := strings.Repeat("?,", len(filter.Tags)-1) + "?"
		conditions = append(conditions, fmt.Sprintf(`n.id IN (
			SELECT nt.note_id FROM note_tags nt 
			JOIN tags t ON nt.tag_id = t.id 
			WHERE t.name IN (%s))`, tagPlaceholders))
		for _, tag := range filter.Tags {
			args = append(args, tag)
		}
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " GROUP BY n.id"

	// Add sorting
	sortColumn := "n.created_at"
	switch filter.SortBy {
	case "updated":
		sortColumn = "n.updated_at"
	case "title":
		sortColumn = "n.title"
	}

	sortOrder := "DESC"
	if filter.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	query += fmt.Sprintf(" ORDER BY %s %s", sortColumn, sortOrder)

	// Add pagination
	if filter.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, filter.Limit)

		if filter.Offset > 0 {
			query += " OFFSET ?"
			args = append(args, filter.Offset)
		}
	}

	rows, err := r.db.conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list notes: %w", err)
	}
	defer rows.Close()

	var notes []*models.Note
	for rows.Next() {
		note := &models.Note{}
		var tagsStr string

		err := rows.Scan(
			&note.ID, &note.Title, &note.Mode, &note.FilePath, &note.FileName,
			&note.ContentHash, &note.CreatedAt, &note.UpdatedAt,
			&note.ContentPreview, &note.WordCount, &tagsStr)

		if err != nil {
			return nil, fmt.Errorf("failed to scan note: %w", err)
		}

		// Parse tags
		if tagsStr != "" {
			note.Tags = strings.Split(tagsStr, ",")
		}

		notes = append(notes, note)
	}

	return notes, nil
}

// Search performs full-text search on notes
func (r *NoteRepository) Search(query string) ([]*models.SearchResult, error) {
	searchQuery := `
		SELECT n.id, n.title, n.mode, n.file_path, n.file_name, n.content_hash,
			n.created_at, n.updated_at, n.content_preview, n.word_count,
			COALESCE(GROUP_CONCAT(t.name, ','), '') as tags,
			fts.bm25(fts) as rank, 'fts' as match_type
		FROM notes_fts fts
		JOIN notes n ON fts.note_id = n.id
		LEFT JOIN note_tags nt ON n.id = nt.note_id
		LEFT JOIN tags t ON nt.tag_id = t.id
		WHERE notes_fts MATCH ?
		GROUP BY n.id
		ORDER BY rank ASC`

	rows, err := r.db.conn.Query(searchQuery, query)
	if err != nil {
		// Fallback to simple LIKE search if FTS fails
		return r.fallbackSearch(query)
	}
	defer rows.Close()

	var results []*models.SearchResult
	for rows.Next() {
		result := &models.SearchResult{}
		var tagsStr string

		err := rows.Scan(
			&result.ID, &result.Title, &result.Mode, &result.FilePath, &result.FileName,
			&result.ContentHash, &result.CreatedAt, &result.UpdatedAt,
			&result.ContentPreview, &result.WordCount, &tagsStr,
			&result.Rank, &result.MatchType)

		if err != nil {
			return nil, fmt.Errorf("failed to scan search result: %w", err)
		}

		// Parse tags
		if tagsStr != "" {
			result.Tags = strings.Split(tagsStr, ",")
		}

		// Generate snippet from content preview
		result.Snippet = r.generateSnippet(result.ContentPreview, query)

		results = append(results, result)
	}

	return results, nil
}

// generateSnippet creates a search snippet highlighting the query terms
func (r *NoteRepository) generateSnippet(content, query string) string {
	if content == "" {
		return ""
	}

	// Simple snippet generation - in a real implementation you might want
	// to find the query terms and show context around them
	if len(content) <= 200 {
		return content
	}

	// Try to find the query in the content for context
	queryLower := strings.ToLower(query)
	contentLower := strings.ToLower(content)

	if idx := strings.Index(contentLower, queryLower); idx != -1 {
		// Show context around the match
		start := idx - 50
		if start < 0 {
			start = 0
		}
		end := idx + len(query) + 150
		if end > len(content) {
			end = len(content)
		}

		snippet := content[start:end]
		if start > 0 {
			snippet = "..." + snippet
		}
		if end < len(content) {
			snippet = snippet + "..."
		}
		return snippet
	}

	// Fallback to beginning of content
	return content[:200] + "..."
}

// Delete removes a note from the database
func (r *NoteRepository) Delete(id string) error {
	tx, err := r.db.conn.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete note (cascades to note_tags due to foreign key)
	_, err = tx.Exec("DELETE FROM notes WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}

	// Update tag usage counts
	_, err = tx.Exec(`
		UPDATE tags SET usage_count = usage_count - 1 
		WHERE id IN (
			SELECT tag_id FROM note_tags WHERE note_id = ?
		)`, id)
	if err != nil {
		return fmt.Errorf("failed to update tag usage counts: %w", err)
	}

	return tx.Commit()
}

// insertTagsForNote handles tag insertion for a note
func (r *NoteRepository) insertTagsForNote(tx *sql.Tx, noteID string, tags []string) error {
	for _, tagName := range tags {
		// Get or create tag
		var tagID int
		err := tx.QueryRow("SELECT id FROM tags WHERE name = ?", tagName).Scan(&tagID)

		if err == sql.ErrNoRows {
			// Create new tag
			result, err := tx.Exec(
				"INSERT INTO tags (name, created_at, usage_count) VALUES (?, ?, 1)",
				tagName, time.Now())
			if err != nil {
				return fmt.Errorf("failed to create tag %s: %w", tagName, err)
			}

			id, err := result.LastInsertId()
			if err != nil {
				return fmt.Errorf("failed to get tag ID: %w", err)
			}
			tagID = int(id)
		} else if err != nil {
			return fmt.Errorf("failed to query tag %s: %w", tagName, err)
		} else {
			// Update usage count for existing tag
			_, err = tx.Exec("UPDATE tags SET usage_count = usage_count + 1 WHERE id = ?", tagID)
			if err != nil {
				return fmt.Errorf("failed to update tag usage count: %w", err)
			}
		}

		// Create note-tag relationship
		_, err = tx.Exec("INSERT INTO note_tags (note_id, tag_id) VALUES (?, ?)", noteID, tagID)
		if err != nil {
			return fmt.Errorf("failed to create note-tag relationship: %w", err)
		}
	}

	return nil
}

// fallbackSearch provides simple LIKE-based search when FTS is not available
func (r *NoteRepository) fallbackSearch(query string) ([]*models.SearchResult, error) {
	searchQuery := `
		SELECT n.id, n.title, n.mode, n.file_path, n.file_name, n.content_hash,
			n.created_at, n.updated_at, n.content_preview, n.word_count,
			COALESCE(GROUP_CONCAT(t.name, ','), '') as tags
		FROM notes n
		LEFT JOIN note_tags nt ON n.id = nt.note_id
		LEFT JOIN tags t ON nt.tag_id = t.id
		WHERE n.title LIKE ? OR n.content_preview LIKE ? OR t.name LIKE ?
		GROUP BY n.id
		ORDER BY n.updated_at DESC`

	likeQuery := "%" + query + "%"
	rows, err := r.db.conn.Query(searchQuery, likeQuery, likeQuery, likeQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to execute fallback search: %w", err)
	}
	defer rows.Close()

	var results []*models.SearchResult
	for rows.Next() {
		result := &models.SearchResult{}
		var tagsStr string

		err := rows.Scan(
			&result.ID, &result.Title, &result.Mode, &result.FilePath, &result.FileName,
			&result.ContentHash, &result.CreatedAt, &result.UpdatedAt,
			&result.ContentPreview, &result.WordCount, &tagsStr)

		if err != nil {
			return nil, fmt.Errorf("failed to scan fallback search result: %w", err)
		}

		// Parse tags
		if tagsStr != "" {
			result.Tags = strings.Split(tagsStr, ",")
		}

		result.Rank = 1.0
		result.MatchType = "fallback"
		result.Snippet = result.ContentPreview
		if len(result.Snippet) > 200 {
			result.Snippet = result.Snippet[:200] + "..."
		}

		results = append(results, result)
	}

	return results, nil
}
