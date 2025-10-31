package cmd

import (
	"fmt"
	"sort"

	"github.com/charmbracelet/lipgloss"
	"github.com/sk25469/jot/app"
	"github.com/sk25469/jot/models"
	"github.com/sk25469/jot/styles"
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
	// Beautiful header
	header := styles.RenderHeader("Statistics")
	fmt.Println(header)
	
	// Create main stats section
	mainStats := createMainStatsSection(stats)
	fmt.Println(mainStats)
	fmt.Println()

	// Tags section
	if len(stats.TagCounts) > 0 {
		tagsSection := createTagsSection(stats.TagCounts)
		fmt.Println(tagsSection)
		fmt.Println()
	}

	// Modes section
	if len(stats.ModeStats) > 0 {
		modesSection := createModesSection(stats.ModeStats)
		fmt.Println(modesSection)
	}
}

func createMainStatsSection(stats *models.StatsResult) string {
	// Create formatted stat entries
	totalText := fmt.Sprintf("%s %s",
		styles.StatsLabelStyle.Render("Total notes:"),
		styles.StatsValueStyle.Render(fmt.Sprintf("%d", stats.TotalNotes)))
	
	weeklyText := fmt.Sprintf("%s %s",
		styles.StatsLabelStyle.Render("This week:"),
		styles.StatsValueStyle.Render(fmt.Sprintf("%d", stats.NotesThisWeek)))
	
	todayText := fmt.Sprintf("%s %s",
		styles.StatsLabelStyle.Render("Created today:"),
		styles.StatsValueStyle.Render(fmt.Sprintf("%d", stats.CreatedToday)))
	
	wordsText := fmt.Sprintf("%s %s",
		styles.StatsLabelStyle.Render("Total words:"),
		styles.StatsValueStyle.Render(fmt.Sprintf("%d", stats.WordCount)))
	
	// Progress bar for weekly activity
	progressBar := ""
	if stats.TotalNotes > 0 {
		progressBar = styles.RenderProgress(stats.NotesThisWeek, 10) // Assuming target of 10 notes per week
	}
	
	// Combine all stats
	statsContent := lipgloss.JoinVertical(
		lipgloss.Left,
		totalText,
		weeklyText,
		todayText,
		wordsText,
	)
	
	if progressBar != "" {
		statsContent = lipgloss.JoinVertical(
			lipgloss.Left,
			statsContent,
			"",
			lipgloss.NewStyle().Foreground(styles.Subtle).Render("Weekly Activity:"),
			progressBar,
		)
	}
	
	return styles.RenderBox("📊 Overview", statsContent)
}

func createTagsSection(tagCounts map[string]int) string {
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

	// Create styled tag entries
	var tagEntries []string
	limit := 8 // Show top 8 tags
	if len(tagList) < limit {
		limit = len(tagList)
	}
	
	for i := 0; i < limit; i++ {
		tag := tagList[i]
		tagStyle := styles.GetTagStyle(tag.tag)
		countStyle := styles.StatsValueStyle
		
		entry := lipgloss.JoinHorizontal(
			lipgloss.Left,
			tagStyle.Render(tag.tag),
			"  ",
			countStyle.Render(fmt.Sprintf("(%d)", tag.count)),
		)
		tagEntries = append(tagEntries, entry)
	}
	
	// Arrange tags in a grid (2 columns)
	var rows []string
	for i := 0; i < len(tagEntries); i += 2 {
		if i+1 < len(tagEntries) {
			row := lipgloss.JoinHorizontal(
				lipgloss.Left,
				tagEntries[i],
				"    ", // Spacing between columns
				tagEntries[i+1],
			)
			rows = append(rows, row)
		} else {
			rows = append(rows, tagEntries[i])
		}
	}
	
	tagsContent := lipgloss.JoinVertical(lipgloss.Left, rows...)
	return styles.RenderBox("🏷️  Popular Tags", tagsContent)
}

func createModesSection(modeStats map[string]int) string {
	var modeEntries []string
	
	for mode, count := range modeStats {
		modeStyle := styles.GetModeStyle(mode)
		countStyle := styles.StatsValueStyle
		
		entry := lipgloss.JoinHorizontal(
			lipgloss.Left,
			modeStyle.Render(mode),
			"  ",
			countStyle.Render(fmt.Sprintf("(%d)", count)),
		)
		modeEntries = append(modeEntries, entry)
	}
	
	modesContent := lipgloss.JoinVertical(lipgloss.Left, modeEntries...)
	return styles.RenderBox("📝 Note Modes", modesContent)
}