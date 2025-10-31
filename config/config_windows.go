package config

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

type Config struct {
	Editor      string `mapstructure:"editor"`
	DefaultMode string `mapstructure:"default_mode"`
	StoragePath string `mapstructure:"storage_path"`
}

var AppConfig Config

func InitConfig() error {
	// Set default values
	viper.SetDefault("editor", getDefaultEditor())
	viper.SetDefault("default_mode", "dev")
	viper.SetDefault("storage_path", getDefaultStoragePath())

	// Config file settings
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(getJotDir())

	// Create jot directory if it doesn't exist
	jotDir := getJotDir()
	if err := os.MkdirAll(jotDir, 0755); err != nil {
		return err
	}

	// Create notes directory if it doesn't exist
	notesDir := filepath.Join(jotDir, "notes")
	if err := os.MkdirAll(notesDir, 0755); err != nil {
		return err
	}

	// Read config file if it exists
	if err := viper.ReadInConfig(); err != nil {
		// Config file not found, create a default one
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := createDefaultConfig(); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	// Unmarshal config into struct
	if err := viper.Unmarshal(&AppConfig); err != nil {
		return err
	}

	// Expand storage path if it contains ~
	if len(AppConfig.StoragePath) >= 2 && AppConfig.StoragePath[:2] == "~/" {
		homeDir, _ := os.UserHomeDir()
		AppConfig.StoragePath = filepath.Join(homeDir, AppConfig.StoragePath[2:])
	}

	return nil
}

func getJotDir() string {
	homeDir, _ := os.UserHomeDir()
	
	// Windows-specific directory
	if runtime.GOOS == "windows" {
		// Use AppData/Roaming on Windows
		appData := os.Getenv("APPDATA")
		if appData != "" {
			return filepath.Join(appData, "jot")
		}
		// Fallback to user home
		return filepath.Join(homeDir, "jot")
	}
	
	// Unix-style hidden directory for other platforms
	return filepath.Join(homeDir, ".jot")
}

func getDefaultStoragePath() string {
	return filepath.Join(getJotDir(), "notes")
}

func getDefaultEditor() string {
	// Check EDITOR environment variable first
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}
	
	// Platform-specific defaults
	switch runtime.GOOS {
	case "windows":
		// Try common Windows editors in order of preference
		editors := []string{"code", "notepad.exe", "notepad++.exe", "vim.exe"}
		for _, editor := range editors {
			if isCommandAvailable(editor) {
				return editor
			}
		}
		// Ultimate fallback to notepad (always available on Windows)
		return "notepad.exe"
	case "darwin":
		// macOS defaults
		editors := []string{"code", "vim", "nano", "open -t"}
		for _, editor := range editors {
			if isCommandAvailable(editor) {
				return editor
			}
		}
		return "vim"
	default:
		// Linux/Unix defaults
		editors := []string{"code", "vim", "nano", "gedit"}
		for _, editor := range editors {
			if isCommandAvailable(editor) {
				return editor
			}
		}
		return "vim"
	}
}

// isCommandAvailable checks if a command is available in PATH
func isCommandAvailable(cmd string) bool {
	_, err := os.Stat(cmd)
	if err == nil {
		return true
	}
	
	// Check in PATH
	paths := filepath.SplitList(os.Getenv("PATH"))
	for _, path := range paths {
		if runtime.GOOS == "windows" {
			// On Windows, try both with and without .exe extension
			exePath := filepath.Join(path, cmd)
			if _, err := os.Stat(exePath); err == nil {
				return true
			}
			if !filepath.Ext(cmd) == ".exe" {
				exePath = filepath.Join(path, cmd+".exe")
				if _, err := os.Stat(exePath); err == nil {
					return true
				}
			}
		} else {
			exePath := filepath.Join(path, cmd)
			if _, err := os.Stat(exePath); err == nil {
				return true
			}
		}
	}
	return false
}

func createDefaultConfig() error {
	configPath := filepath.Join(getJotDir(), "config.yaml")
	
	defaultConfig := `editor: "` + getDefaultEditor() + `"
default_mode: "dev"
storage_path: "` + getDefaultStoragePath() + `"
`
	
	return os.WriteFile(configPath, []byte(defaultConfig), 0644)
}

func GetNotesDir() string {
	return AppConfig.StoragePath
}

func GetJotDir() string {
	return getJotDir()
}