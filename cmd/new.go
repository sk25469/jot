package cmd

import (
	"strings"

	"github.com/sahilsarwar/jot/notes"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [title]",
	Short: "Create a new note",
	Long:  `Create a new markdown note with optional title, tags, and mode.`,
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		title := strings.Join(args, " ")
		if title == "" {
			title = "Untitled"
		}

		tags, _ := cmd.Flags().GetStringSlice("tag")
		mode, _ := cmd.Flags().GetString("mode")

		_, err := notes.CreateNote(title, tags, mode)
		return err
	},
}

func init() {
	newCmd.Flags().StringSliceP("tag", "t", []string{}, "Tags for the note")
	newCmd.Flags().StringP("mode", "m", "", "Mode for the note (defaults to config default)")
}