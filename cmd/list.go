package cmd

import (
	"fmt"
	"strings"

	"github.com/sahilsarwar/jot/app"
	"github.com/sahilsarwar/jot/models"
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
	// Print header
	fmt.Printf("%-8s %-12s %-30s %s\n", "ID", "DATE", "TITLE", "TAGS")
	fmt.Printf("%-8s %-12s %-30s %s\n", "--------", "----", "-----", "----")

	// Print notes
	for _, note := range notesList {
		dateShort := note.CreatedAt.Format("2006-01-02")
		title := formatNoteTitle(note.Title)
		tagsStr := strings.Join(note.Tags, ", ")

		fmt.Printf("%-8s %-12s %-30s %s\n", note.ID, dateShort, title, tagsStr)
	}
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