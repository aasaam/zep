package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"io"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/template"
)

// Environment represents a mapping of environment variable keys to their values
type Environment map[string]string

// NewEnvironment creates a new Environment from a map of environment variables.
func NewEnvironment(envMap map[string]string) Environment {
	return Environment(envMap)
}

// AsString retrieves a string value for the given environment key
// Panics if the key is not found
func (env Environment) AsString(key string) string {
	value, ok := env[key]
	if !ok {
		panic(fmt.Errorf("environment variable '%s' not found", key))
	}
	return value
}

// AsStringOr retrieves a string value for the given environment key
// Returns the defaultValue if the key is not found
func (env Environment) AsStringOr(key, defaultValue string) string {
	value, ok := env[key]
	if !ok {
		return defaultValue
	}
	return value
}

// AsStringSlice retrieves a string value for the given environment key and splits it by delimiter
// Panics if the key is not found
func (env Environment) AsStringSlice(key, delimiter string) []string {
	value, ok := env[key]
	if !ok {
		panic(fmt.Errorf("environment variable '%s' not found", key))
	}
	return strings.Split(value, delimiter)
}

// AsStringSliceTrim retrieves a string value for the given environment key, splits it by delimiter,
// and optionally trims each element using the specified trim characters
// Panics if the key is not found
func (env Environment) AsStringSliceTrim(key, delimiter string, trimChars string) []string {
	value, ok := env[key]
	if !ok {
		panic(fmt.Errorf("environment variable '%s' not found", key))
	}

	elements := strings.Split(value, delimiter)
	for i, element := range elements {
		elements[i] = strings.Trim(element, trimChars)
	}
	return elements
}

// AsBool retrieves a boolean value for the given environment key
// Accepts "true", "1", "yes" as true and "false", "0", "no" as false (case insensitive)
// Panics if the key is not found or the value cannot be parsed as a boolean
func (env Environment) AsBool(key string) bool {
	value, ok := env[key]
	if !ok {
		panic(fmt.Errorf("environment variable '%s' not found", key))
	}

	lowerValue := strings.ToLower(value)
	switch lowerValue {
	case "true", "1", "yes", "on", "enable", "enabled":
		return true
	case "false", "0", "no", "off", "disable", "disabled":
		return false
	default:
		panic(fmt.Errorf("could not parse '%s' (value: '%s') as boolean", key, value))
	}
}

// AsBoolOr retrieves a boolean value for the given environment key
// Accepts "true", "1", "yes" as true and "false", "0", "no" as false (case insensitive)
// Returns the defaultValue if the key is not found or the value cannot be parsed
func (env Environment) AsBoolOr(key string, defaultValue bool) bool {
	value, ok := env[key]
	if !ok {
		return defaultValue
	}

	lowerValue := strings.ToLower(value)
	switch lowerValue {
	case "true", "1", "yes", "on", "enable", "enabled":
		return true
	case "false", "0", "no", "off", "disable", "disabled":
		return false
	default:
		return defaultValue
	}
}

// AsURL retrieves a URL value for the given environment key
// Panics if the key is not found or the value cannot be parsed as a URL
func (env Environment) AsURL(key string) string {
	value, ok := env[key]
	if !ok {
		panic(fmt.Errorf("environment variable '%s' not found", key))
	}
	u, err := url.ParseRequestURI(value)
	if err != nil || u.Scheme == "" {
		panic(fmt.Errorf("could not parse '%s' (value: '%s') as URL: %v", key, value, err))
	}
	return u.String()
}

// AsHostPort retrieves a host:port value for the given environment key
// Prefixes with "http://" before parsing to extract the host
// Panics if the key is not found or the value cannot be parsed
func (env Environment) AsHostPort(key string) string {
	value, ok := env[key]
	if !ok {
		panic(fmt.Errorf("environment variable '%s' not found", key))
	}
	u, err := url.ParseRequestURI("http://" + value)
	if err != nil {
		panic(fmt.Errorf("could not parse '%s' (value: '%s') as URL: %v", key, value, err))
	}
	port := 0
	uPort, err := strconv.Atoi(u.Port())
	if err == nil {
		port = uPort
	}
	if port < 1 || port > 65535 {
		panic(fmt.Errorf("port '%s' (value: '%s') is out of range (1-65535)", key, value))
	}
	return u.Host
}

// AsInt retrieves an integer value for the given environment key
// Panics if the key is not found or the value cannot be parsed as an integer
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

// AsIntOr retrieves an integer value for the given environment key
// Returns the defaultValue if the key is not found or the value cannot be parsed
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

// AsIntSlice retrieves a string value, splits it by delimiter, and converts each element to an integer
// Panics if the key is not found or any element cannot be parsed as an integer
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

// AsFloat retrieves an integer value for the given environment key
// Panics if the key is not found or the value cannot be parsed as an integer
func (env Environment) AsFloat(key string) float64 {
	value, ok := env[key]
	if !ok {
		panic(fmt.Errorf("environment variable '%s' not found", key))
	}

	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		panic(fmt.Errorf("could not parse '%s' (value: '%s') as float: %v", key, value, err))
	}
	return floatValue
}

// AsFloatOr retrieves an integer value for the given environment key
// Returns the defaultValue if the key is not found or the value cannot be parsed
func (env Environment) AsFloatOr(key string, defaultValue float64) float64 {
	value, ok := env[key]
	if !ok {
		return defaultValue
	}

	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultValue
	}
	return floatValue
}

// AsFloatSlice retrieves a string value, splits it by delimiter, and converts each element to an integer
// Panics if the key is not found or any element cannot be parsed as an integer
func (env Environment) AsFloatSlice(key, delimiter string) []float64 {
	value, ok := env[key]
	if !ok {
		panic(fmt.Errorf("environment variable '%s' not found", key))
	}

	stringElements := strings.Split(value, delimiter)
	intSlice := make([]float64, 0, len(stringElements))

	for _, element := range stringElements {
		trimmedElement := strings.TrimSpace(element)
		floatValue, err := strconv.ParseFloat(trimmedElement, 64)
		if err != nil {
			panic(fmt.Errorf("on key '%s', could not parse '%s' as integer: %v", key, trimmedElement, err))
		}
		intSlice = append(intSlice, floatValue)
	}

	return intSlice
}

// AsPort retrieves a port number for the given environment key
// Validates that the port is in the valid range (1-65535)
// Panics if the key is not found, the value cannot be parsed, or is outside the valid range
func (env Environment) AsPort(key string) int {
	value, ok := env[key]
	if !ok {
		panic(fmt.Errorf("environment variable '%s' not found", key))
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		panic(fmt.Errorf("could not parse '%s' (value: '%s') as integer: %v", key, value, err))
	}
	if intValue < 1 || intValue > 65535 {
		panic(fmt.Errorf("port '%s' (value: '%s') is out of range (1-65535)", key, value))
	}
	return intValue
}

// AsPortOr retrieves a port number for the given environment key
// Returns the defaultPort if the key is not found, the value cannot be parsed, or is outside the valid range
// Panics if the defaultPort is outside the valid range (1-65535)
func (env Environment) AsPortOr(key string, defaultPort int) int {
	if defaultPort < 1 || defaultPort > 65535 {
		panic(fmt.Errorf("default port '%d' is out of range (1-65535)", defaultPort))
	}
	value, ok := env[key]
	if !ok {
		return defaultPort
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultPort
	}
	if intValue < 1 || intValue > 65535 {
		return defaultPort
	}
	return intValue
}

// All returns the entire environment map
func (env Environment) All() map[string]string {
	return env
}

// SortAll returns a new map with keys sorted alphabetically
func (env Environment) SortAll() map[string]string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	result := make(map[string]string, len(env))
	for _, k := range keys {
		result[k] = env[k]
	}
	return result
}

// isEmpty checks if a string is empty or contains only whitespace
func isEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// isNotEmpty checks if a string is not empty or contains only whitespace
func isNotEmpty(s string) bool {
	return strings.TrimSpace(s) != ""
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// containsCaseInsensitive checks if a string contains a substring, ignoring case
func containsCaseInsensitive(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// containsAny checks if a string contains any of the specified characters
func containsAny(s, chars string) bool {
	return strings.ContainsAny(s, chars)
}

// hasPrefix checks if a string starts with the specified prefix
func hasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

// hasSuffix checks if a string ends with the specified suffix
func hasSuffix(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

// toLower converts a string to lowercase
func toLower(s string) string {
	return strings.ToLower(s)
}

// toUpper converts a string to uppercase
func toUpper(s string) string {
	return strings.ToUpper(s)
}

// trim removes the specified characters from the beginning and end of a string
func trim(s, cutset string) string {
	return strings.Trim(s, cutset)
}

// trimLeft removes the specified characters from the beginning of a string
func trimLeft(s, cutset string) string {
	return strings.TrimLeft(s, cutset)
}

// trimRight removes the specified characters from the end of a string
func trimRight(s, cutset string) string {
	return strings.TrimRight(s, cutset)
}

// trimSpace removes whitespace from the beginning and end of a string
func trimSpace(s string) string {
	return strings.TrimSpace(s)
}

// base64Encode encodes a string to base64
func base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

// base64Decode decodes a base64 string
// Panics if the string cannot be decoded
func base64Decode(s string) string {
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(fmt.Errorf("could not decode base64 string: %v", err))
	}
	return string(decoded)
}

// hash computes a hash of the input string using the specified algorithm
// Supported algorithms: md5, sha1, sha224, sha256, sha512
// Panics if an unsupported algorithm is specified
func hash(input string, algorithm string) string {
	switch strings.ToLower(algorithm) {
	case "md5":
		return fmt.Sprintf("%x", md5.Sum([]byte(input)))
	case "sha1":
		return fmt.Sprintf("%x", sha1.Sum([]byte(input)))
	case "sha224":
		return fmt.Sprintf("%x", sha256.Sum224([]byte(input)))
	case "sha256":
		return fmt.Sprintf("%x", sha256.Sum256([]byte(input)))
	case "sha512":
		return fmt.Sprintf("%x", sha512.Sum512([]byte(input)))
	default:
		panic(fmt.Errorf("unsupported hash algorithm: %s", algorithm))
	}
}

// sequence generates a slice of integers from start to end (inclusive)
// Returns nil if start > end
func sequence(start, end int) []int {
	if start > end {
		return nil
	}
	seq := make([]int, end-start+1)
	for i := start; i <= end; i++ {
		seq[i-start] = i
	}
	return seq
}

// fileExistOrDefault copies a default file to the destination path if the destination does not exist
// Preserves the file mode of the default file
// Panics if any file operation fails
func fileExistOrDefault(destination string, defaultPath string) bool {
	if _, err := os.Stat(destination); os.IsNotExist(err) {
		fileInfo, err := os.Stat(defaultPath)
		if err != nil {
			panic(fmt.Errorf("could not read file permissions for '%s': %v", defaultPath, err))
		}

		r, err := os.Open(defaultPath)
		if err != nil {
			panic(fmt.Errorf("could not open file '%s': %v", defaultPath, err))
		}
		defer r.Close()

		w, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fileInfo.Mode())
		if err != nil {
			panic(fmt.Errorf("could not open file '%s': %v", destination, err))
		}
		defer w.Close()

		if _, err := io.Copy(w, r); err != nil {
			panic(fmt.Errorf("could not copy file '%s' to '%s': %v", defaultPath, destination, err))
		}
	}
	return true
}

// GetTemplateFunctions returns a map of functions that can be used in templates
// The functions provide access to environment variables and various string utilities
func GetTemplateFunctions(env Environment) template.FuncMap {
	return template.FuncMap{
		// Environment accessors
		"all":               env.All,
		"asString":          env.AsString,
		"asStringOr":        env.AsStringOr,
		"asStringSlice":     env.AsStringSlice,
		"asStringSliceTrim": env.AsStringSliceTrim,
		"asBool":            env.AsBool,
		"asBoolOr":          env.AsBoolOr,
		"asInt":             env.AsInt,
		"asIntOr":           env.AsIntOr,
		"asIntSlice":        env.AsIntSlice,
		"asFloat":           env.AsFloat,
		"asFloatOr":         env.AsFloatOr,
		"asFloatSlice":      env.AsFloatSlice,
		"asPort":            env.AsPort,
		"asPortOr":          env.AsPortOr,
		"asURL":             env.AsURL,
		"asHostPort":        env.AsHostPort,
		"sortAll":           env.SortAll,

		// String functions
		"contains":                contains,
		"containsAny":             containsAny,
		"containsCaseInsensitive": containsCaseInsensitive,
		"hasPrefix":               hasPrefix,
		"hasSuffix":               hasSuffix,
		"toLower":                 toLower,
		"toUpper":                 toUpper,
		"trim":                    trim,
		"trimLeft":                trimLeft,
		"trimRight":               trimRight,
		"trimSpace":               trimSpace,
		"isEmpty":                 isEmpty,
		"isNotEmpty":              isNotEmpty,

		// Encoding and utility functions
		"base64Decode": base64Decode,
		"base64Encode": base64Encode,
		"hash":         hash,
		"sequence":     sequence,

		// File
		"fileExistOrDefault": fileExistOrDefault,
	}
}

// RenderTemplate processes the template string with the given environment.
// It returns the rendered output or an error if template parsing or execution fails.
func RenderTemplate(templateContent string, env Environment) (string, error) {
	tmpl := template.New("envTemplate").Funcs(GetTemplateFunctions(env))
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
