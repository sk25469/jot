package cmd

import (
	"fmt"
	"strings"

	"github.com/sahilsarwar/jot/notes"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all notes",
	Long:  `List all notes with optional filtering by tag and mode.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		tagFilter, _ := cmd.Flags().GetString("tag")
		modeFilter, _ := cmd.Flags().GetString("mode")

		notesList, err := notes.ListNotes(tagFilter, modeFilter)
		if err != nil {
			return err
		}

		if len(notesList) == 0 {
			fmt.Println("No notes found.")
			return nil
		}

		// Print header
		fmt.Printf("%-3s %-12s %-30s %s\n", "ID", "DATE", "TITLE", "TAGS")
		fmt.Printf("%-3s %-12s %-30s %s\n", "---", "----", "-----", "----")

		// Print notes
		for i, note := range notesList {
			idShort := fmt.Sprintf("%03d", i+1)
			dateShort := note.Date.Format("2006-01-02")
			title := note.Title
			if len(title) > 30 {
				title = title[:27] + "..."
			}
			tagsStr := strings.Join(note.Tags, ", ")

			fmt.Printf("%-3s %-12s %-30s %s\n", idShort, dateShort, title, tagsStr)
		}

		return nil
	},
}

func init() {
	listCmd.Flags().StringP("tag", "t", "", "Filter by tag")
	listCmd.Flags().StringP("mode", "m", "", "Filter by mode")
}