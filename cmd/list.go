package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/sk25469/jot/app"
	"github.com/sk25469/jot/models"
	"github.com/sk25469/jot/styles"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all notes",
	Long:  `List all notes with optional filtering by tag and mode.`,
	RunE:  runListCommand,
}

func runListCommand(cmd *cobra.Command, args []string) error {
	tagFilter, _ := cmd.Flags().GetString("tag")
	modeFilter, _ := cmd.Flags().GetString("mode")

	notesList, err := app.Instance.NoteService.ListNotes(tagFilter, modeFilter)
	if err != nil {
		return err
	}

	if len(notesList) == 0 {
		fmt.Println("No notes found.")
		return nil
	}

	printNotesList(notesList)
	return nil
}

func printNotesList(notesList []*models.Note) {
	if len(notesList) == 0 {
		fmt.Println(styles.WarningStyle.Render("No notes found."))
		return
	}

	// Print beautiful header
	header := styles.RenderHeader(fmt.Sprintf("Notes (%d)", len(notesList)))
	fmt.Println(header)
	fmt.Println()

	// Create a table-like layout using lipgloss
	for i, note := range notesList {
		// Create the main note entry
		noteEntry := createNoteEntry(note)
		fmt.Println(noteEntry)

		// Add spacing between notes, but not after the last one
		if i < len(notesList)-1 {
			fmt.Println()
		}
	}

	// Footer with total count
	fmt.Println()
	fmt.Println(styles.RenderSeparator())
	totalText := fmt.Sprintf("Total: %d notes", len(notesList))
	fmt.Println(styles.StatsLabelStyle.Render(totalText))
}

func createNoteEntry(note *models.Note) string {
	// Format the ID with style
	idText := styles.IDStyle.Render(note.ID)

	// Format the date
	dateText := styles.DateStyle.Render(note.CreatedAt.Format("2006-01-02"))

	// Format the title
	title := note.Title
	if len(title) > 50 {
		title = title[:47] + "..."
	}
	titleText := styles.ContentStyle.Render(title)

	// Format the mode badge
	modeText := styles.GetModeStyle(note.Mode).Render(note.Mode)

	// Format tags with colors
	tagsText := ""
	if len(note.Tags) > 0 {
		tagsText = styles.RenderTags(note.Tags)
	} else {
		tagsText = lipgloss.NewStyle().Foreground(styles.Subtle).Render("no tags")
	}

	// Create the first line with ID, date, and title
	firstLine := lipgloss.JoinHorizontal(
		lipgloss.Left,
		idText,
		"  ",
		dateText,
		"  ",
		modeText,
		"  ",
		titleText,
	)

	// Create the second line with tags (indented)
	secondLine := lipgloss.NewStyle().
		MarginLeft(2).
		Render("tags: " + tagsText)

	// Combine both lines
	return lipgloss.JoinVertical(
		lipgloss.Left,
		firstLine,
		secondLine,
	)
}

func formatNoteTitle(title string) string {
	if len(title) > 30 {
		return title[:27] + "..."
	}
	return title
}

func init() {
	listCmd.Flags().StringP("tag", "t", "", "Filter by tag")
	listCmd.Flags().StringP("mode", "m", "", "Filter by mode")
}
