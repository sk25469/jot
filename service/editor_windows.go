package service

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/sk25469/jot/config"
)

// openInEditor opens a file in the configured editor with Windows compatibility
func (s *NoteService) openInEditor(filePath string) error {
	editor := config.AppConfig.Editor
	
	return openFileInEditor(editor, filePath)
}

// openFileInEditor handles cross-platform editor launching
func openFileInEditor(editor, filePath string) error {
	var cmd *exec.Cmd
	
	switch runtime.GOOS {
	case "windows":
		cmd = createWindowsCommand(editor, filePath)
	case "darwin":
		cmd = createMacCommand(editor, filePath)
	default:
		cmd = createUnixCommand(editor, filePath)
	}
	
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	return cmd.Run()
}

func createWindowsCommand(editor, filePath string) *exec.Cmd {
	// Handle special cases for Windows
	switch {
	case strings.Contains(editor, "notepad++"):
		return exec.Command("cmd", "/c", editor, filePath)
	case strings.Contains(editor, "notepad"):
		return exec.Command(editor, filePath)
	case editor == "code":
		return exec.Command("code", filePath)
	case strings.Contains(editor, "vim"):
		// Try to find vim in common locations
		locations := []string{
			"vim.exe",
			"C:\\Program Files\\Vim\\vim90\\vim.exe",
			"C:\\Program Files (x86)\\Vim\\vim90\\vim.exe",
			"C:\\tools\\vim\\vim90\\vim.exe",
		}
		for _, loc := range locations {
			if _, err := os.Stat(loc); err == nil {
				return exec.Command(loc, filePath)
			}
		}
		// Fallback to PATH
		return exec.Command("vim.exe", filePath)
	default:
		// Try as-is first
		if strings.Contains(editor, " ") {
			// Handle commands with arguments
			parts := strings.Fields(editor)
			if len(parts) > 1 {
				args := append(parts[1:], filePath)
				return exec.Command(parts[0], args...)
			}
		}
		return exec.Command(editor, filePath)
	}
}

func createMacCommand(editor, filePath string) *exec.Cmd {
	switch {
	case strings.HasPrefix(editor, "open"):
		// Handle "open -t" or similar
		parts := strings.Fields(editor)
		args := append(parts[1:], filePath)
		return exec.Command("open", args...)
	case editor == "code":
		return exec.Command("code", filePath)
	default:
		return exec.Command(editor, filePath)
	}
}

func createUnixCommand(editor, filePath string) *exec.Cmd {
	// Handle commands with arguments (like "code --wait")
	if strings.Contains(editor, " ") {
		parts := strings.Fields(editor)
		if len(parts) > 1 {
			args := append(parts[1:], filePath)
			return exec.Command(parts[0], args...)
		}
	}
	return exec.Command(editor, filePath)
}

// GetAvailableEditors returns a list of available editors on the current platform
func GetAvailableEditors() []string {
	var editors []string
	
	switch runtime.GOOS {
	case "windows":
		candidates := []string{
			"notepad.exe",
			"notepad++.exe",
			"code",
			"vim.exe",
			"sublime_text.exe",
			"atom.exe",
		}
		for _, editor := range candidates {
			if isEditorAvailable(editor) {
				editors = append(editors, editor)
			}
		}
	case "darwin":
		candidates := []string{
			"code",
			"vim",
			"nano",
			"subl",
			"atom",
			"open -t",
		}
		for _, editor := range candidates {
			if isEditorAvailable(editor) {
				editors = append(editors, editor)
			}
		}
	default:
		candidates := []string{
			"code",
			"vim",
			"nano",
			"gedit",
			"subl",
			"atom",
			"emacs",
		}
		for _, editor := range candidates {
			if isEditorAvailable(editor) {
				editors = append(editors, editor)
			}
		}
	}
	
	return editors
}

func isEditorAvailable(editor string) bool {
	// Handle commands with spaces (like "open -t")
	if strings.Contains(editor, " ") {
		parts := strings.Fields(editor)
		editor = parts[0]
	}
	
	// Try to execute with --version or similar to test availability
	cmd := exec.Command(editor, "--version")
	if err := cmd.Run(); err == nil {
		return true
	}
	
	// Try --help as fallback
	cmd = exec.Command(editor, "--help")
	if err := cmd.Run(); err == nil {
		return true
	}
	
	// For Windows, also try without .exe extension
	if runtime.GOOS == "windows" && !strings.HasSuffix(editor, ".exe") {
		return isEditorAvailable(editor + ".exe")
	}
	
	return false
}

// ValidateEditor checks if the configured editor is available
func ValidateEditor(editor string) error {
	if !isEditorAvailable(editor) {
		available := GetAvailableEditors()
		if len(available) > 0 {
			return fmt.Errorf("editor '%s' not found. Available editors: %s", 
				editor, strings.Join(available, ", "))
		}
		return fmt.Errorf("editor '%s' not found and no alternative editors detected", editor)
	}
	return nil
}