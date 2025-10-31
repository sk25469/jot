package cmd

import (
	"fmt"
	"os"

	"github.com/sahilsarwar/jot/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "jot",
	Short: "A lightning-fast terminal-based note-taking and journaling CLI",
	Long: `jot is a terminal-based note-taking CLI that feels like git and fzf had a baby.
It provides lightning-fast capture and recall of thoughts, code, or reflections
without ever leaving the terminal.

All notes are stored as plain markdown files in ~/.jot/notes/`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return config.InitConfig()
	},
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(openCmd)
	rootCmd.AddCommand(statsCmd)
}