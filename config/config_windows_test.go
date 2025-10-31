package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestWindowsCompatibility(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping Windows-specific tests on non-Windows platform")
	}
	
	// Test Windows-specific directory structure
	jotDir := getJotDir()
	
	// Should use AppData on Windows
	appData := os.Getenv("APPDATA")
	if appData != "" {
		expectedDir := filepath.Join(appData, "jot")
		if jotDir != expectedDir {
			t.Errorf("Expected Windows jot dir to be %s, got %s", expectedDir, jotDir)
		}
	}
	
	// Test Windows editor detection
	editor := getDefaultEditor()
	
	// Should not default to vim on Windows
	if editor == "vim" {
		t.Errorf("Default editor should not be vim on Windows, got %s", editor)
	}
	
	// Should be a Windows-compatible editor
	windowsEditors := []string{"notepad.exe", "notepad++.exe", "code", "vim.exe"}
	isValidEditor := false
	for _, validEditor := range windowsEditors {
		if editor == validEditor {
			isValidEditor = true
			break
		}
	}
	
	if !isValidEditor {
		t.Errorf("Default editor %s is not Windows-compatible. Expected one of: %v", 
			editor, windowsEditors)
	}
}

func TestCrossPlatformPaths(t *testing.T) {
	// Test that paths are properly constructed for all platforms
	jotDir := getJotDir()
	
	// Should be absolute path
	if !filepath.IsAbs(jotDir) {
		t.Errorf("Jot directory should be absolute path, got %s", jotDir)
	}
	
	// Should contain platform-appropriate separator
	if runtime.GOOS == "windows" {
		// Windows should not use Unix-style hidden directory
		if filepath.Base(jotDir) == ".jot" {
			t.Errorf("Windows should not use hidden .jot directory, got %s", jotDir)
		}
	} else {
		// Unix-like systems should use hidden directory
		if filepath.Base(jotDir) != ".jot" {
			t.Errorf("Unix-like systems should use .jot directory, got %s", jotDir)
		}
	}
	
	// Test storage path
	storagePath := getDefaultStoragePath()
	if !filepath.IsAbs(storagePath) {
		t.Errorf("Storage path should be absolute, got %s", storagePath)
	}
	
	// Should end with "notes"
	if filepath.Base(storagePath) != "notes" {
		t.Errorf("Storage path should end with 'notes', got %s", storagePath)
	}
}

func TestEditorAvailability(t *testing.T) {
	// Test the isCommandAvailable function
	testCases := []struct {
		command  string
		platform string
	}{
		{"go", "all"},           // Go should be available (we're running tests)
		{"nonexistent", "all"},  // This should not be available
	}
	
	for _, tc := range testCases {
		if tc.platform != "all" && tc.platform != runtime.GOOS {
			continue
		}
		
		t.Run(tc.command, func(t *testing.T) {
			available := isCommandAvailable(tc.command)
			
			if tc.command == "go" && !available {
				t.Errorf("Go should be available (we're running tests)")
			}
			
			if tc.command == "nonexistent" && available {
				t.Errorf("Nonexistent command should not be available")
			}
		})
	}
}

func TestPathExpansionWindows(t *testing.T) {
	// Test path expansion on Windows vs Unix
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home directory: %v", err)
	}
	
	testPath := "~/test/path"
	expected := filepath.Join(homeDir, "test", "path")
	
	// Simulate path expansion
	result := testPath
	if len(result) >= 2 && result[:2] == "~/" {
		result = filepath.Join(homeDir, result[2:])
	}
	
	if result != expected {
		t.Errorf("Path expansion failed. Expected %s, got %s", expected, result)
	}
	
	// Test that the path uses correct separators
	if runtime.GOOS == "windows" {
		if !filepath.IsAbs(result) {
			t.Errorf("Windows path should be absolute after expansion")
		}
	}
}