package cmd

import (
	"fmt"
	"strings"

	"github.com/sahilsarwar/jot/notes"
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

	results, err := notes.SearchNotes(query)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		fmt.Printf("No notes found matching: %s\n", query)
		return nil
	}

	printSearchResults(results, query)
	return nil
}

func printSearchResults(results []*notes.Note, query string) {
	fmt.Printf("Found %d note(s) matching: %s\n\n", len(results), query)

	for _, note := range results {
		fmt.Printf("→ %s %s\n", note.Date.Format("2006-01-02"), note.Title)
		if len(note.Tags) > 0 {
			fmt.Printf("  Tags: %s\n", strings.Join(note.Tags, ", "))
		}
		fmt.Println()
	}
}