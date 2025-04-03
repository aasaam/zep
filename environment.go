package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/template"
)

// Environment represents a collection of environment variables.
type Environment map[string]string

// NewEnvironment creates a new Environment from a map of environment variables.
func NewEnvironment(envMap map[string]string) Environment {
	return Environment(envMap)
}

// AsString returns the environment variable as a string.
// It panics if the variable doesn't exist.
func (env Environment) AsString(key string) string {
	value, ok := env[key]
	if !ok {
		panic(fmt.Errorf("environment variable '%s' not found", key))
	}
	return value
}

// AsStringOr returns the environment variable as a string or a default value if not found.
func (env Environment) AsStringOr(key, defaultValue string) string {
	value, ok := env[key]
	if !ok {
		return defaultValue
	}
	return value
}

// AsStringSlice returns the environment variable as a string slice.
// It panics if the variable doesn't exist.
func (env Environment) AsStringSlice(key, delimiter string) []string {
	value, ok := env[key]
	if !ok {
		panic(fmt.Errorf("environment variable '%s' not found", key))
	}
	return strings.Split(value, delimiter)
}

// AsStringSliceTrim returns the environment variable as a string slice with optional trimming.
// It panics if the variable doesn't exist.
func (env Environment) AsStringSliceTrim(key, delimiter string, trim bool, trimChars string) []string {
	value, ok := env[key]
	if !ok {
		panic(fmt.Errorf("environment variable '%s' not found", key))
	}

	elements := strings.Split(value, delimiter)
	if trim {
		for i, element := range elements {
			elements[i] = strings.Trim(element, trimChars)
		}
	}
	return elements
}

// AsBool returns the environment variable as a boolean.
// It recognizes "true", "1", "yes" as true and "false", "0", "no" as false (case-insensitive).
// It panics if the variable doesn't exist or cannot be parsed.
func (env Environment) AsBool(key string) bool {
	value, ok := env[key]
	if !ok {
		panic(fmt.Errorf("environment variable '%s' not found", key))
	}

	lowerValue := strings.ToLower(value)
	switch lowerValue {
	case "true", "1", "yes":
		return true
	case "false", "0", "no":
		return false
	default:
		panic(fmt.Errorf("could not parse '%s' (value: '%s') as boolean", key, value))
	}
}

// AsBoolOr returns the environment variable as a boolean or a default value if not found or invalid.
func (env Environment) AsBoolOr(key string, defaultValue bool) bool {
	value, ok := env[key]
	if !ok {
		return defaultValue
	}

	lowerValue := strings.ToLower(value)
	switch lowerValue {
	case "true", "1", "yes":
		return true
	case "false", "0", "no":
		return false
	default:
		return defaultValue
	}
}

// AsInt returns the environment variable as an integer.
// It panics if the variable doesn't exist or cannot be parsed.
func (env Environment) AsInt(key string) int {
	value, ok := env[key]
	if !ok {
		panic(fmt.Errorf("environment variable '%s' not found", key))
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		panic(fmt.Errorf("could not parse '%s' (value: '%s') as integer: %v", key, value, err))
	}
	return intValue
}

// AsIntOr returns the environment variable as an integer or a default value if not found or invalid.
func (env Environment) AsIntOr(key string, defaultValue int) int {
	value, ok := env[key]
	if !ok {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

// AsIntSlice returns the environment variable as an integer slice.
// It panics if the variable doesn't exist or any element cannot be parsed.
func (env Environment) AsIntSlice(key, delimiter string) []int {
	value, ok := env[key]
	if !ok {
		panic(fmt.Errorf("environment variable '%s' not found", key))
	}

	stringElements := strings.Split(value, delimiter)
	intSlice := make([]int, 0, len(stringElements))

	for _, element := range stringElements {
		trimmedElement := strings.TrimSpace(element)
		intValue, err := strconv.Atoi(trimmedElement)
		if err != nil {
			panic(fmt.Errorf("on key '%s', could not parse '%s' as integer: %v", key, trimmedElement, err))
		}
		intSlice = append(intSlice, intValue)
	}

	return intSlice
}

// IsEmpty checks if an environment variable is empty or doesn't exist.
func (env Environment) IsEmpty(key string) bool {
	value, ok := env[key]
	if !ok {
		return true // Consider non-existent as empty
	}
	return strings.TrimSpace(value) == ""
}

// RenderTemplate processes the template string with the given environment.
// It returns the rendered output or an error if template parsing or execution fails.
func RenderTemplate(templateContent string, env Environment) (string, error) {
	tmpl := template.New("envTemplate").Funcs(template.FuncMap{
		"seq": func(start, end int) []int {
			nums := make([]int, end-start+1)
			for i := range nums {
				nums[i] = start + i
			}
			return nums
		},
		"asString":          env.AsString,
		"asStringOr":        env.AsStringOr,
		"asStringSlice":     env.AsStringSlice,
		"asStringSliceTrim": env.AsStringSliceTrim,
		"asBool":            env.AsBool,
		"asBoolOr":          env.AsBoolOr,
		"asInt":             env.AsInt,
		"asIntOr":           env.AsIntOr,
		"asIntSlice":        env.AsIntSlice,
		"isEmpty":           env.IsEmpty,
	})

	parsedTmpl, err := tmpl.Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %w", err)
	}

	var buf bytes.Buffer
	if err := parsedTmpl.Execute(&buf, env); err != nil {
		return "", fmt.Errorf("error executing template: %w", err)
	}

	return buf.String(), nil
}
