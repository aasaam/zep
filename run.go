package main

import (
	"fmt"
	"os"
	"strings"
)

// Run executes the template rendering process.
func Run(args []string, environ []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("usage: %s <template-file>", args[0])
	}

	templateFile := args[1]

	envMap := make(map[string]string)
	for _, e := range environ {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) == 2 {
			envMap[pair[0]] = pair[1]
		}
	}
	env := NewEnvironment(envMap)

	templateContent, err := os.ReadFile(templateFile)
	if err != nil {
		return "", fmt.Errorf("error reading template file '%s': %v", templateFile, err)
	}

	output, err := RenderTemplate(string(templateContent), env)
	if err != nil {
		return "", fmt.Errorf("error rendering template: %v", err)
	}

	return output, nil
}
