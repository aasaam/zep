package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMainFunction(t *testing.T) {

	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	tempDir := t.TempDir()

	templatePath := filepath.Join(tempDir, "template.txt")
	templateContent := "Hello {{asStringOr \"NAME\" \"World\"}}!"
	err := os.WriteFile(templatePath, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create template file: %v", err)
	}

	os.Args = []string{"zep", templatePath}

	main()
}
