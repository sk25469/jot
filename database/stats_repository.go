package database

import (
	"fmt"
	"strings"
	"time"

	"github.com/sahilsarwar/jot/models"
)

// StatsRepository handles database operations for statistics
type StatsRepository struct {
	db *DB
}

// NewStatsRepository creates a new stats repository
func NewStatsRepository(db *DB) *StatsRepository {
	return &StatsRepository{db: db}
}

// GetStats retrieves comprehensive statistics about notes
func (r *StatsRepository) GetStats() (*models.StatsResult, error) {
	stats := &models.StatsResult{
		TagCounts: make(map[string]int),
		ModeStats: make(map[string]int),
	}

	// Get total notes count
	err := r.db.conn.QueryRow("SELECT COUNT(*) FROM notes").Scan(&stats.TotalNotes)
	if err != nil {
		return nil, fmt.Errorf("failed to get total notes count: %w", err)
	}

	// Get notes created this week
	weekAgo := time.Now().AddDate(0, 0, -7)
	err = r.db.conn.QueryRow(
		"SELECT COUNT(*) FROM notes WHERE created_at >= ?", weekAgo).Scan(&stats.NotesThisWeek)
	if err != nil {
		return nil, fmt.Errorf("failed to get notes this week: %w", err)
	}

	// Get notes created today
	today := time.Now().Truncate(24 * time.Hour)
	err = r.db.conn.QueryRow(
		"SELECT COUNT(*) FROM notes WHERE created_at >= ?", today).Scan(&stats.CreatedToday)
	if err != nil {
		return nil, fmt.Errorf("failed to get notes created today: %w", err)
	}

	// Get total word count
	err = r.db.conn.QueryRow("SELECT COALESCE(SUM(word_count), 0) FROM notes").Scan(&stats.WordCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get total word count: %w", err)
	}

	// Get tag statistics
	tagRows, err := r.db.conn.Query(`
		SELECT t.name, COUNT(nt.note_id) as usage_count
		FROM tags t
		LEFT JOIN note_tags nt ON t.id = nt.tag_id
		GROUP BY t.id, t.name
		ORDER BY usage_count DESC`)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag statistics: %w", err)
	}
	defer tagRows.Close()

	for tagRows.Next() {
		var tagName string
		var count int
		if err := tagRows.Scan(&tagName, &count); err != nil {
			return nil, fmt.Errorf("failed to scan tag stats: %w", err)
		}
		stats.TagCounts[tagName] = count
	}

	// Get mode statistics
	modeRows, err := r.db.conn.Query(`
		SELECT mode, COUNT(*) as count
		FROM notes
		GROUP BY mode
		ORDER BY count DESC`)
	if err != nil {
		return nil, fmt.Errorf("failed to get mode statistics: %w", err)
	}
	defer modeRows.Close()

	for modeRows.Next() {
		var mode string
		var count int
		if err := modeRows.Scan(&mode, &count); err != nil {
			return nil, fmt.Errorf("failed to scan mode stats: %w", err)
		}
		stats.ModeStats[mode] = count
	}

	return stats, nil
}

// GetTagUsage returns the most used tags with their counts
func (r *StatsRepository) GetTagUsage(limit int) (map[string]int, error) {
	query := `
		SELECT t.name, COUNT(nt.note_id) as usage_count
		FROM tags t
		LEFT JOIN note_tags nt ON t.id = nt.tag_id
		GROUP BY t.id, t.name
		ORDER BY usage_count DESC`
	
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := r.db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag usage: %w", err)
	}
	defer rows.Close()

	tagCounts := make(map[string]int)
	for rows.Next() {
		var tagName string
		var count int
		if err := rows.Scan(&tagName, &count); err != nil {
			return nil, fmt.Errorf("failed to scan tag usage: %w", err)
		}
		tagCounts[tagName] = count
	}

	return tagCounts, nil
}

// GetRecentActivity returns notes created in the last N days
func (r *StatsRepository) GetRecentActivity(days int) ([]*models.Note, error) {
	since := time.Now().AddDate(0, 0, -days)
	
	query := `
		SELECT n.id, n.title, n.mode, n.file_path, n.file_name, n.content_hash,
			n.created_at, n.updated_at, n.content_preview, n.word_count,
			COALESCE(GROUP_CONCAT(t.name, ','), '') as tags
		FROM notes n
		LEFT JOIN note_tags nt ON n.id = nt.note_id
		LEFT JOIN tags t ON nt.tag_id = t.id
		WHERE n.created_at >= ?
		GROUP BY n.id
		ORDER BY n.created_at DESC`

	rows, err := r.db.conn.Query(query, since)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent activity: %w", err)
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
			return nil, fmt.Errorf("failed to scan recent note: %w", err)
		}

		if tagsStr != "" {
			note.Tags = strings.Split(tagsStr, ",")
		}

		notes = append(notes, note)
	}

	return notes, nil
}