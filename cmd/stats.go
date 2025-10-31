package cmd

import (
	"fmt"
	"sort"

	"github.com/sahilsarwar/jot/app"
	"github.com/sahilsarwar/jot/models"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show note statistics",
	Long:  `Display statistics about your notes including total count, weekly activity, and popular tags.`,
	RunE:  runStatsCommand,
}

func runStatsCommand(cmd *cobra.Command, args []string) error {
	stats, err := app.Instance.NoteService.GetStats()
	if err != nil {
		return err
	}

	printNoteStatistics(stats)
	return nil
}

func printNoteStatistics(stats *models.StatsResult) {
	fmt.Printf("📊 Note Statistics\n\n")
	fmt.Printf("Total notes: %d\n", stats.TotalNotes)
	fmt.Printf("This week: %d\n", stats.NotesThisWeek)
	fmt.Printf("Created today: %d\n", stats.CreatedToday)
	fmt.Printf("Total words: %d\n", stats.WordCount)

	if len(stats.TagCounts) > 0 {
		printTopTags(stats.TagCounts)
	}

	if len(stats.ModeStats) > 0 {
		fmt.Printf("\nModes:\n")
		for mode, count := range stats.ModeStats {
			fmt.Printf("  %s (%d)\n", mode, count)
		}
	}
}

func printTopTags(tagCounts map[string]int) {
	fmt.Printf("\nMost used tags:\n")
	
	// Sort tags by count
	type tagCount struct {
		tag   string
		count int
	}
	var tagList []tagCount
	for tag, count := range tagCounts {
		tagList = append(tagList, tagCount{tag, count})
	}
	sort.Slice(tagList, func(i, j int) bool {
		return tagList[i].count > tagList[j].count
	})

	// Show top 5 tags
	limit := 5
	if len(tagList) < limit {
		limit = len(tagList)
	}
	for i := 0; i < limit; i++ {
		fmt.Printf("  %s (%d)\n", tagList[i].tag, tagList[i].count)
	}
}