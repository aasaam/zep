package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunInvalidTemplate(t *testing.T) {
	tempDir := t.TempDir()

	templatePath := filepath.Join(tempDir, "invalid_template.txt")
	templateContent := "Hello {{.NAME!}"
	err := os.WriteFile(templatePath, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create template file: %v", err)
	}

	args := []string{"zep", templatePath}
	env := []string{"NAME=World"}

	output, err := Run(args, env)
	if err == nil {
		t.Errorf("Expected error for invalid template but got none")
	}

	expectedErrorMsg := "error rendering template"
	if err != nil && !contains(err.Error(), expectedErrorMsg) {
		t.Errorf("Expected error message to contain %q but got %q", expectedErrorMsg, err.Error())
	}

	if output != "" {
		t.Errorf("Expected empty output for invalid template but got %q", output)
	}
}

func TestRun(t *testing.T) {

	tempDir := t.TempDir()

	tests := []struct {
		name            string
		args            []string
		env             []string
		templateFile    string
		templateContent string
		expectedOutput  string
		expectError     bool
	}{
		{
			name:            "Basic template rendering",
			args:            []string{"zep", "template.txt"},
			env:             []string{"NAME=World", "COUNT=3"},
			templateFile:    "template.txt",
			templateContent: "Hello {{.NAME}}! Count: {{.COUNT}}",
			expectedOutput:  "Hello World! Count: 3",
			expectError:     false,
		},
		{
			name:            "Template with function",
			args:            []string{"zep", "template.txt"},
			env:             []string{"NAME=World", "ITEMS=a,b,c"},
			templateFile:    "template.txt",
			templateContent: "{{range asStringSlice \"ITEMS\" \",\"}}{{.}} {{end}}",
			expectedOutput:  "a b c ",
			expectError:     false,
		},
		{
			name:        "Missing template file",
			args:        []string{"zep", "nonexistent.txt"},
			env:         []string{"NAME=World"},
			expectError: true,
		},
		{
			name:        "Invalid arguments",
			args:        []string{"zep"},
			env:         []string{},
			expectError: true,
		},
		{
			name:            "Template with error",
			args:            []string{"zep", "template.txt"},
			env:             []string{"NAME=World"},
			templateFile:    "template.txt",
			templateContent: "{{asInt \"COUNT\"}}",
			expectError:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			if tc.templateFile != "" && tc.templateContent != "" {
				templatePath := filepath.Join(tempDir, tc.templateFile)
				err := os.WriteFile(templatePath, []byte(tc.templateContent), 0644)
				if err != nil {
					t.Fatalf("Failed to create template file: %v", err)
				}

				if len(tc.args) > 1 {
					tc.args[1] = templatePath
				}
			}

			output, err := Run(tc.args, tc.env)

			if tc.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Got unexpected error: %v", err)
			}

			if !tc.expectError && output != tc.expectedOutput {
				t.Errorf("Expected output %q but got %q", tc.expectedOutput, output)
			}
		})
	}
}

func TestRunWithEmptyEnvironment(t *testing.T) {

	tempDir := t.TempDir()

	templatePath := filepath.Join(tempDir, "template.txt")
	templateContent := "Hello {{asStringOr \"NAME\" \"default\"}}!"
	err := os.WriteFile(templatePath, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create template file: %v", err)
	}

	args := []string{"zep", templatePath}
	env := []string{}

	output, err := Run(args, env)
	if err != nil {
		t.Errorf("Unexpected error with empty environment: %v", err)
	}

	expectedOutput := "Hello default!"
	if output != expectedOutput {
		t.Errorf("Expected output %q with empty environment but got %q", expectedOutput, output)
	}
}

func TestRunWithMalformedEnvironment(t *testing.T) {

	tempDir := t.TempDir()

	templatePath := filepath.Join(tempDir, "template.txt")
	templateContent := "Hello {{asStringOr \"NAME\" \"default\"}}!"
	err := os.WriteFile(templatePath, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create template file: %v", err)
	}

	args := []string{"zep", templatePath}
	env := []string{"NAME=World", "INVALID_ENTRY", "=VALUE_WITHOUT_KEY"}

	output, err := Run(args, env)
	if err != nil {
		t.Errorf("Unexpected error with malformed environment: %v", err)
	}

	expectedOutput := "Hello World!"
	if output != expectedOutput {
		t.Errorf("Expected output %q with malformed environment but got %q", expectedOutput, output)
	}
}
