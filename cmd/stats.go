package cmd

import (
	"fmt"
	"sort"

	"github.com/sahilsarwar/jot/notes"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show note statistics",
	Long:  `Display statistics about your notes including total count, weekly activity, and popular tags.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		stats, err := notes.GetStats()
		if err != nil {
			return err
		}

		totalNotes := stats["total_notes"].(int)
		thisWeek := stats["this_week"].(int)
		tagCounts := stats["tag_counts"].(map[string]int)

		fmt.Printf("📊 Note Statistics\n\n")
		fmt.Printf("Total notes: %d\n", totalNotes)
		fmt.Printf("This week: %d\n", thisWeek)

		if len(tagCounts) > 0 {
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

		return nil
	},
}