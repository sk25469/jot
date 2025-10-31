package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/sahilsarwar/jot/app"
	"github.com/sahilsarwar/jot/models"
	"github.com/sahilsarwar/jot/styles"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search notes",
	Long:  `Search notes by title, content, or tags using fuzzy matching.`,
	Args:  cobra.MinimumNArgs(1),
	RunE:  runSearchCommand,
}

func runSearchCommand(cmd *cobra.Command, args []string) error {
	query := strings.Join(args, " ")

	results, err := app.Instance.NoteService.SearchNotes(query)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		noResultsMsg := fmt.Sprintf("No notes found matching: %s", query)
		fmt.Println(styles.WarningStyle.Render(noResultsMsg))
		return nil
	}

	printSearchResults(results, query)
	return nil
}

func printSearchResults(results []*models.SearchResult, query string) {
	// Beautiful header with search query
	header := styles.RenderHeader(fmt.Sprintf("Search Results (%d)", len(results)))
	queryText := lipgloss.NewStyle().
		Foreground(styles.Warning).
		Bold(true).
		Render(fmt.Sprintf("Query: \"%s\"", query))
	
	fmt.Println(header)
	fmt.Println(queryText)
	fmt.Println()

	for i, result := range results {
		// Create styled search result entry
		searchEntry := createSearchResultEntry(result, i+1)
		fmt.Println(searchEntry)
		
		// Add spacing between results
		if i < len(results)-1 {
			fmt.Println()
		}
	}
	
	// Footer
	fmt.Println()
	fmt.Println(styles.RenderSeparator())
	totalText := fmt.Sprintf("Found %d matches", len(results))
	fmt.Println(styles.StatsLabelStyle.Render(totalText))
}

func createSearchResultEntry(result *models.SearchResult, index int) string {
	// Rank indicator
	rankText := lipgloss.NewStyle().
		Foreground(styles.Primary).
		Bold(true).
		Render(fmt.Sprintf("#%d", index))
	
	// ID and date
	idText := styles.IDStyle.Render(result.ID)
	dateText := styles.DateStyle.Render(result.CreatedAt.Format("2006-01-02"))
	
	// Title with highlighting potential
	titleText := styles.ContentStyle.Bold(true).Render(result.Title)
	
	// Mode badge
	modeText := styles.GetModeStyle(result.Mode).Render(result.Mode)
	
	// Relevance score (if available)
	scoreText := ""
	if result.Rank > 0 {
		score := fmt.Sprintf("%.2f", result.Rank)
		scoreText = lipgloss.NewStyle().
			Foreground(styles.Accent).
			Render(fmt.Sprintf("score: %s", score))
	}
	
	// First line: rank, ID, date, mode, title
	firstLine := lipgloss.JoinHorizontal(
		lipgloss.Left,
		rankText,
		"  ",
		idText,
		"  ",
		dateText,
		"  ",
		modeText,
		"  ",
		titleText,
	)
	
	if scoreText != "" {
		firstLine = lipgloss.JoinHorizontal(
			lipgloss.Left,
			firstLine,
			"  ",
			scoreText,
		)
	}
	
	// Second line: tags
	secondLine := ""
	if len(result.Tags) > 0 {
		tagsText := styles.RenderTags(result.Tags)
		secondLine = lipgloss.NewStyle().
			MarginLeft(4).
			Render("tags: " + tagsText)
	}
	
	// Third line: snippet (if available)
	thirdLine := ""
	if result.Snippet != "" && result.Snippet != result.Title {
		snippetText := result.Snippet
		if len(snippetText) > 100 {
			snippetText = snippetText[:97] + "..."
		}
		thirdLine = lipgloss.NewStyle().
			MarginLeft(4).
			Foreground(styles.Muted).
			Italic(true).
			Render("\"" + snippetText + "\"")
	}
	
	// Combine all lines
	lines := []string{firstLine}
	if secondLine != "" {
		lines = append(lines, secondLine)
	}
	if thirdLine != "" {
		lines = append(lines, thirdLine)
	}
	
	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}