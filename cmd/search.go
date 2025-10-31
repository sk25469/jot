package cmd

import (
	"fmt"
	"strings"

	"github.com/sahilsarwar/jot/app"
	"github.com/sahilsarwar/jot/models"
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
		fmt.Printf("No notes found matching: %s\n", query)
		return nil
	}

	printSearchResults(results, query)
	return nil
}

func printSearchResults(results []*models.SearchResult, query string) {
	fmt.Printf("Found %d note(s) matching: %s\n\n", len(results), query)

	for _, result := range results {
		fmt.Printf("→ %s %s\n", result.CreatedAt.Format("2006-01-02"), result.Title)
		if len(result.Tags) > 0 {
			fmt.Printf("  Tags: %s\n", strings.Join(result.Tags, ", "))
		}
		if result.Snippet != "" && result.Snippet != result.Title {
			fmt.Printf("  %s\n", result.Snippet)
		}
		fmt.Println()
	}
}