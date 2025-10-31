package styles

import (
	"crypto/sha1"
	"fmt"
	"math"

	"github.com/charmbracelet/lipgloss"
)

// Color palette
var (
	// Primary colors
	Primary   = lipgloss.Color("#00D4AA")
	Secondary = lipgloss.Color("#6B73FF")
	Accent    = lipgloss.Color("#FFB86C")
	
	// Status colors
	Success = lipgloss.Color("#50FA7B")
	Warning = lipgloss.Color("#F1FA8C")
	Error   = lipgloss.Color("#FF5555")
	
	// Neutral colors
	Subtle = lipgloss.Color("#6272A4")
	Muted  = lipgloss.Color("#44475A")
	Text   = lipgloss.Color("#F8F8F2")
)

// Base styles
var (
	// Title styles
	TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(Primary).
		MarginLeft(1).
		MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
		Foreground(Secondary).
		Bold(true)

	// ID styles
	IDStyle = lipgloss.NewStyle().
		Foreground(Accent).
		Bold(true)

	// Date styles
	DateStyle = lipgloss.NewStyle().
		Foreground(Subtle).
		Italic(true)

	// Content styles
	ContentStyle = lipgloss.NewStyle().
		Foreground(Text)

	PreviewStyle = lipgloss.NewStyle().
		Foreground(Muted).
		Italic(true).
		MarginLeft(2)

	// List styles
	ListItemStyle = lipgloss.NewStyle().
		PaddingLeft(2).
		MarginBottom(1)

	ListHeaderStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(Primary).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(Primary).
		MarginBottom(1).
		PaddingBottom(1)

	// Stats styles
	StatsLabelStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(Secondary)

	StatsValueStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(Accent)

	// Search styles
	SearchQueryStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(Warning)

	SearchResultStyle = lipgloss.NewStyle().
		MarginLeft(1).
		PaddingLeft(1).
		BorderStyle(lipgloss.NormalBorder()).
		BorderLeft(true).
		BorderForeground(Subtle)

	SnippetStyle = lipgloss.NewStyle().
		Foreground(Muted).
		Italic(true).
		MarginLeft(4)

	// Error styles
	ErrorStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(Error).
		Background(lipgloss.Color("#44475A")).
		Padding(0, 1)

	WarningStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(Warning).
		Background(lipgloss.Color("#44475A")).
		Padding(0, 1)

	SuccessStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(Success).
		Background(lipgloss.Color("#44475A")).
		Padding(0, 1)
)

// Tag colors - deterministic based on tag name
var tagColors = []lipgloss.Color{
	lipgloss.Color("#FF79C6"), // Pink
	lipgloss.Color("#8BE9FD"), // Cyan
	lipgloss.Color("#50FA7B"), // Green
	lipgloss.Color("#FFB86C"), // Orange
	lipgloss.Color("#BD93F9"), // Purple
	lipgloss.Color("#F1FA8C"), // Yellow
	lipgloss.Color("#FF5555"), // Red
	lipgloss.Color("#6272A4"), // Blue
}

// GetTagStyle returns a consistent style for a tag based on its name
func GetTagStyle(tagName string) lipgloss.Style {
	// Hash the tag name to get consistent color
	h := sha1.New()
	h.Write([]byte(tagName))
	hash := h.Sum(nil)
	colorIndex := int(hash[0]) % len(tagColors)
	
	return lipgloss.NewStyle().
		Background(tagColors[colorIndex]).
		Foreground(lipgloss.Color("#282A36")). // Dark background for readability
		Bold(true).
		Padding(0, 1).
		MarginRight(1)
}

// GetModeStyle returns a style for note modes
func GetModeStyle(mode string) lipgloss.Style {
	switch mode {
	case "dev":
		return lipgloss.NewStyle().
			Background(Primary).
			Foreground(lipgloss.Color("#282A36")).
			Bold(true).
			Padding(0, 1)
	case "journal":
		return lipgloss.NewStyle().
			Background(Secondary).
			Foreground(lipgloss.Color("#F8F8F2")).
			Bold(true).
			Padding(0, 1)
	case "meeting":
		return lipgloss.NewStyle().
			Background(Accent).
			Foreground(lipgloss.Color("#282A36")).
			Bold(true).
			Padding(0, 1)
	default:
		return lipgloss.NewStyle().
			Background(Muted).
			Foreground(Text).
			Bold(true).
			Padding(0, 1)
	}
}

// Utility functions

// RenderTags renders a slice of tags with consistent colors
func RenderTags(tags []string) string {
	if len(tags) == 0 {
		return ""
	}
	
	var rendered []string
	for _, tag := range tags {
		rendered = append(rendered, GetTagStyle(tag).Render(tag))
	}
	
	return lipgloss.JoinHorizontal(lipgloss.Left, rendered...)
}

// RenderHeader creates a styled header
func RenderHeader(title string) string {
	return TitleStyle.Render("📝 " + title)
}

// RenderSeparator creates a visual separator
func RenderSeparator() string {
	return lipgloss.NewStyle().
		Foreground(Subtle).
		Render("────────────────────────────────────────")
}

// RenderProgress creates a simple progress indicator
func RenderProgress(current, total int) string {
	if total == 0 {
		return ""
	}
	
	percentage := float64(current) / float64(total) * 100
	filled := int(math.Round(percentage / 10)) // 10 segments
	
	progress := ""
	for i := 0; i < 10; i++ {
		if i < filled {
			progress += "█"
		} else {
			progress += "░"
		}
	}
	
	return lipgloss.NewStyle().
		Foreground(Primary).
		Render(fmt.Sprintf("[%s] %d/%d (%.1f%%)", progress, current, total, percentage))
}

// RenderBox creates a bordered box around content
func RenderBox(title, content string) string {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(Primary).
		Padding(1, 2).
		Width(60).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				lipgloss.NewStyle().Bold(true).Foreground(Primary).Render(title),
				"",
				content,
			),
		)
}