package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMainFunction(t *testing.T) {
	// Save original args and restore them after test
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// Create a temporary directory for test files
	tempDir := t.TempDir()
	defer os.RemoveAll(tempDir)

	// Create a template file
	templatePath := filepath.Join(tempDir, "template.txt")
	templateContent := "Hello {{asStringOr \"NAME\" \"World\"}}!"
	err := os.WriteFile(templatePath, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create template file: %v", err)
	}

	// Set command line args
	// Note: This test is limited since we can't easily capture stdout in the main function
	// We're just ensuring it doesn't panic
	os.Args = []string{"zep", templatePath}

	// If main() doesn't panic, test passes
	main()
}
