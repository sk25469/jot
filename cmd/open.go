package cmd

import (
	"strings"

	"github.com/sahilsarwar/jot/app"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open <id or title>",
	Short: "Open a note in your editor",
	Long:  `Open a note by ID or partial title match in your configured editor.`,
	Args:  cobra.MinimumNArgs(1),
	RunE:  runOpenCommand,
}

func runOpenCommand(cmd *cobra.Command, args []string) error {
	identifier := strings.Join(args, " ")
	return app.Instance.NoteService.OpenNote(identifier)
}