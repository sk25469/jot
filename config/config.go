package config

import (
	"os"
	"path/filepath"

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
	if AppConfig.StoragePath[:2] == "~/" {
		homeDir, _ := os.UserHomeDir()
		AppConfig.StoragePath = filepath.Join(homeDir, AppConfig.StoragePath[2:])
	}

	return nil
}

func getJotDir() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".jot")
}

func getDefaultStoragePath() string {
	return filepath.Join(getJotDir(), "notes")
}

func getDefaultEditor() string {
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}
	return "vim"
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
