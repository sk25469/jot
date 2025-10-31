package cmd

import (
	"strings"

	"github.com/sk25469/jot/app"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [title]",
	Short: "Create a new note",
	Long:  `Create a new markdown note with optional title, tags, and mode.`,
	Args:  cobra.ArbitraryArgs,
	RunE:  runNewCommand,
}

func runNewCommand(cmd *cobra.Command, args []string) error {
	title := getNoteTitleFromArgs(args)
	tags, _ := cmd.Flags().GetStringSlice("tag")
	mode, _ := cmd.Flags().GetString("mode")

	_, err := app.Instance.NoteService.CreateNote(title, tags, mode)
	return err
}

func getNoteTitleFromArgs(args []string) string {
	title := strings.Join(args, " ")
	if title == "" {
		title = "Untitled"
	}
	return title
}

func init() {
	newCmd.Flags().StringSliceP("tag", "t", []string{}, "Tags for the note")
	newCmd.Flags().StringP("mode", "m", "", "Mode for the note (defaults to config default)")
}
