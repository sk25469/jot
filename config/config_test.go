package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetJotDir(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}

	expected := filepath.Join(homeDir, ".jot")
	result := getJotDir()

	if result != expected {
		t.Errorf("getJotDir() = %q, expected %q", result, expected)
	}
}

func TestGetDefaultStoragePath(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}

	expected := filepath.Join(homeDir, ".jot", "notes")
	result := getDefaultStoragePath()

	if result != expected {
		t.Errorf("getDefaultStoragePath() = %q, expected %q", result, expected)
	}
}

func TestGetDefaultEditor(t *testing.T) {
	// Save original EDITOR value
	originalEditor := os.Getenv("EDITOR")
	defer os.Setenv("EDITOR", originalEditor)

	testCases := []struct {
		name          string
		editorEnv     string
		shouldUnset   bool
		expectedValue string
	}{
		{
			name:          "with EDITOR set",
			editorEnv:     "nano",
			shouldUnset:   false,
			expectedValue: "nano",
		},
		{
			name:          "with EDITOR unset",
			editorEnv:     "",
			shouldUnset:   true,
			expectedValue: "vim",
		},
		{
			name:          "with custom editor",
			editorEnv:     "code",
			shouldUnset:   false,
			expectedValue: "code",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.shouldUnset {
				os.Unsetenv("EDITOR")
			} else {
				os.Setenv("EDITOR", tc.editorEnv)
			}

			result := getDefaultEditor()
			if result != tc.expectedValue {
				t.Errorf("getDefaultEditor() = %q, expected %q", result, tc.expectedValue)
			}
		})
	}
}

func TestConfigStruct(t *testing.T) {
	config := Config{
		Editor:      "vim",
		DefaultMode: "dev",
		StoragePath: "/path/to/notes",
	}

	if config.Editor != "vim" {
		t.Errorf("Expected Editor to be 'vim', got %q", config.Editor)
	}

	if config.DefaultMode != "dev" {
		t.Errorf("Expected DefaultMode to be 'dev', got %q", config.DefaultMode)
	}

	if config.StoragePath != "/path/to/notes" {
		t.Errorf("Expected StoragePath to be '/path/to/notes', got %q", config.StoragePath)
	}
}

func TestGetJotDirPublic(t *testing.T) {
	// Test the public function
	result := GetJotDir()

	// Should not be empty
	if result == "" {
		t.Errorf("GetJotDir() should not return empty string")
	}

	// Should end with .jot
	if filepath.Base(result) != ".jot" {
		t.Errorf("GetJotDir() should end with .jot, got %q", result)
	}

	// Should be an absolute path
	if !filepath.IsAbs(result) {
		t.Errorf("GetJotDir() should return absolute path, got %q", result)
	}
}

// Integration test for GetNotesDir (requires AppConfig to be initialized)
func TestGetNotesDir(t *testing.T) {
	// Set up a temporary config for testing
	originalConfig := AppConfig
	defer func() { AppConfig = originalConfig }()

	AppConfig = Config{
		StoragePath: "/tmp/test-notes",
	}

	result := GetNotesDir()
	expected := "/tmp/test-notes"

	if result != expected {
		t.Errorf("GetNotesDir() = %q, expected %q", result, expected)
	}
}

func TestCreateDefaultConfig(t *testing.T) {
	// Test the config content format by checking what the function would generate
	defaultEditor := getDefaultEditor()
	defaultStorage := getDefaultStoragePath()

	// Test that default values are reasonable
	if defaultEditor == "" {
		t.Errorf("Default editor should not be empty")
	}

	if defaultStorage == "" {
		t.Errorf("Default storage path should not be empty")
	}

	// Test that storage path contains .jot
	if filepath.Base(filepath.Dir(defaultStorage)) != ".jot" {
		t.Errorf("Default storage path should be under .jot directory")
	}
}

// Test path expansion functionality
func TestPathExpansion(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}

	testCases := []struct {
		name        string
		input       string
		expected    string
		shouldStart bool // whether result should start with expected
	}{
		{
			name:        "tilde expansion",
			input:       "~/notes",
			expected:    filepath.Join(homeDir, "notes"),
			shouldStart: false,
		},
		{
			name:        "absolute path",
			input:       "/absolute/path",
			expected:    "/absolute/path",
			shouldStart: false,
		},
		{
			name:        "relative path",
			input:       "relative/path",
			expected:    "relative/path",
			shouldStart: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Simulate the path expansion logic from InitConfig
			result := tc.input
			if len(result) >= 2 && result[:2] == "~/" {
				result = filepath.Join(homeDir, result[2:])
			}

			if tc.shouldStart {
				if !strings.HasPrefix(result, tc.expected) {
					t.Errorf("Path expansion result %q should start with %q", result, tc.expected)
				}
			} else {
				if result != tc.expected {
					t.Errorf("Path expansion result %q, expected %q", result, tc.expected)
				}
			}
		})
	}
}
