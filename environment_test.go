package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestNewEnvironment(t *testing.T) {
	envMap := map[string]string{"KEY1": "value1", "KEY2": "value2"}
	env := NewEnvironment(envMap)
	if !reflect.DeepEqual(env, Environment(envMap)) {
		t.Errorf("NewEnvironment did not create the expected environment")
	}
}

func TestAsString(t *testing.T) {
	env := Environment{"KEY": "value"}

	tests := []struct {
		name      string
		key       string
		want      string
		wantPanic bool
	}{
		{name: "existing key", key: "KEY", want: "value", wantPanic: false},
		{name: "non-existent key", key: "NONEXISTENT", want: "", wantPanic: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("AsString did not panic for key %s", tc.key)
					}
				}()
			}

			got := env.AsString(tc.key)
			if got != tc.want {
				t.Errorf("AsString(%q) = %q, want %q", tc.key, got, tc.want)
			}
		})
	}
}

func TestAsStringOr(t *testing.T) {
	env := Environment{"KEY": "value"}

	tests := []struct {
		name         string
		key          string
		defaultValue string
		want         string
	}{
		{name: "existing key", key: "KEY", defaultValue: "default", want: "value"},
		{name: "non-existent key", key: "NONEXISTENT", defaultValue: "default", want: "default"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := env.AsStringOr(tc.key, tc.defaultValue)
			if got != tc.want {
				t.Errorf("AsStringOr(%q, %q) = %q, want %q", tc.key, tc.defaultValue, got, tc.want)
			}
		})
	}
}

func TestAsStringSlice(t *testing.T) {
	env := Environment{
		"COMMA_LIST": "a,b,c",
		"COLON_LIST": "x:y:z",
	}

	tests := []struct {
		name      string
		key       string
		delimiter string
		want      []string
		wantPanic bool
	}{
		{name: "comma delimiter", key: "COMMA_LIST", delimiter: ",", want: []string{"a", "b", "c"}, wantPanic: false},
		{name: "colon delimiter", key: "COLON_LIST", delimiter: ":", want: []string{"x", "y", "z"}, wantPanic: false},
		{name: "non-existent key", key: "NONEXISTENT", delimiter: ",", want: nil, wantPanic: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("AsStringSlice did not panic for key %s", tc.key)
					}
				}()
			}

			got := env.AsStringSlice(tc.key, tc.delimiter)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("AsStringSlice(%q, %q) = %v, want %v", tc.key, tc.delimiter, got, tc.want)
			}
		})
	}
}

func TestAsStringSliceTrim(t *testing.T) {
	env := Environment{
		"SPACES": " a , b , c ",
		"QUOTES": `"x","y","z"`,
	}

	tests := []struct {
		name      string
		key       string
		delimiter string
		trim      bool
		trimChars string
		want      []string
		wantPanic bool
	}{
		{name: "trim spaces", key: "SPACES", delimiter: ",", trim: true, trimChars: " ", want: []string{"a", "b", "c"}, wantPanic: false},
		{name: "no trim", key: "SPACES", delimiter: ",", trim: false, trimChars: " ", want: []string{" a ", " b ", " c "}, wantPanic: false},
		{name: "trim quotes", key: "QUOTES", delimiter: ",", trim: true, trimChars: `"`, want: []string{"x", "y", "z"}, wantPanic: false},
		{name: "non-existent key", key: "NONEXISTENT", delimiter: ",", trim: true, trimChars: " ", want: nil, wantPanic: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("AsStringSliceTrim did not panic for key %s", tc.key)
					}
				}()
			}

			got := env.AsStringSliceTrim(tc.key, tc.delimiter, tc.trim, tc.trimChars)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("AsStringSliceTrim(%q, %q, %v, %q) = %v, want %v",
					tc.key, tc.delimiter, tc.trim, tc.trimChars, got, tc.want)
			}
		})
	}
}

func TestAsBool(t *testing.T) {
	env := Environment{
		"TRUE1":   "true",
		"TRUE2":   "TRUE",
		"TRUE3":   "1",
		"TRUE4":   "yes",
		"FALSE1":  "false",
		"FALSE2":  "FALSE",
		"FALSE3":  "0",
		"FALSE4":  "no",
		"INVALID": "invalid",
	}

	tests := []struct {
		name      string
		key       string
		want      bool
		wantPanic bool
	}{
		{name: "lowercase true", key: "TRUE1", want: true, wantPanic: false},
		{name: "uppercase TRUE", key: "TRUE2", want: true, wantPanic: false},
		{name: "numeric 1", key: "TRUE3", want: true, wantPanic: false},
		{name: "yes", key: "TRUE4", want: true, wantPanic: false},
		{name: "lowercase false", key: "FALSE1", want: false, wantPanic: false},
		{name: "uppercase FALSE", key: "FALSE2", want: false, wantPanic: false},
		{name: "numeric 0", key: "FALSE3", want: false, wantPanic: false},
		{name: "no", key: "FALSE4", want: false, wantPanic: false},
		{name: "invalid value", key: "INVALID", want: false, wantPanic: true},
		{name: "non-existent key", key: "NONEXISTENT", want: false, wantPanic: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("AsBool did not panic for key %s", tc.key)
					}
				}()
			}

			got := env.AsBool(tc.key)
			if got != tc.want {
				t.Errorf("AsBool(%q) = %v, want %v", tc.key, got, tc.want)
			}
		})
	}
}

func TestAsBoolOr(t *testing.T) {
	env := Environment{
		"TRUE":    "true",
		"FALSE":   "false",
		"INVALID": "invalid",
	}

	tests := []struct {
		name         string
		key          string
		defaultValue bool
		want         bool
	}{
		{name: "existing true", key: "TRUE", defaultValue: false, want: true},
		{name: "existing false", key: "FALSE", defaultValue: true, want: false},
		{name: "invalid value", key: "INVALID", defaultValue: true, want: true},
		{name: "non-existent key", key: "NONEXISTENT", defaultValue: true, want: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := env.AsBoolOr(tc.key, tc.defaultValue)
			if got != tc.want {
				t.Errorf("AsBoolOr(%q, %v) = %v, want %v", tc.key, tc.defaultValue, got, tc.want)
			}
		})
	}
}

func TestAsInt(t *testing.T) {
	env := Environment{
		"POSITIVE": "123",
		"NEGATIVE": "-42",
		"ZERO":     "0",
		"INVALID":  "not-an-int",
	}

	tests := []struct {
		name      string
		key       string
		want      int
		wantPanic bool
	}{
		{name: "positive number", key: "POSITIVE", want: 123, wantPanic: false},
		{name: "negative number", key: "NEGATIVE", want: -42, wantPanic: false},
		{name: "zero", key: "ZERO", want: 0, wantPanic: false},
		{name: "invalid value", key: "INVALID", want: 0, wantPanic: true},
		{name: "non-existent key", key: "NONEXISTENT", want: 0, wantPanic: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("AsInt did not panic for key %s", tc.key)
					}
				}()
			}

			got := env.AsInt(tc.key)
			if got != tc.want {
				t.Errorf("AsInt(%q) = %v, want %v", tc.key, got, tc.want)
			}
		})
	}
}

func TestAsIntOr(t *testing.T) {
	env := Environment{
		"VALID":   "123",
		"INVALID": "not-an-int",
	}

	tests := []struct {
		name         string
		key          string
		defaultValue int
		want         int
	}{
		{name: "existing valid", key: "VALID", defaultValue: 999, want: 123},
		{name: "existing invalid", key: "INVALID", defaultValue: 999, want: 999},
		{name: "non-existent key", key: "NONEXISTENT", defaultValue: 999, want: 999},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := env.AsIntOr(tc.key, tc.defaultValue)
			if got != tc.want {
				t.Errorf("AsIntOr(%q, %v) = %v, want %v", tc.key, tc.defaultValue, got, tc.want)
			}
		})
	}
}

func TestAsIntSlice(t *testing.T) {
	env := Environment{
		"VALID":  "1,2,3",
		"SPACES": " 4, 5, 6 ",
		"MIXED":  "7, invalid, 9",
	}

	tests := []struct {
		name      string
		key       string
		delimiter string
		want      []int
		wantPanic bool
	}{
		{name: "valid values", key: "VALID", delimiter: ",", want: []int{1, 2, 3}, wantPanic: false},
		{name: "values with spaces", key: "SPACES", delimiter: ",", want: []int{4, 5, 6}, wantPanic: false},
		{name: "mixed valid and invalid", key: "MIXED", delimiter: ",", want: nil, wantPanic: true},
		{name: "non-existent key", key: "NONEXISTENT", delimiter: ",", want: nil, wantPanic: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("AsIntSlice did not panic for key %s", tc.key)
					}
				}()
			}

			got := env.AsIntSlice(tc.key, tc.delimiter)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("AsIntSlice(%q, %q) = %v, want %v", tc.key, tc.delimiter, got, tc.want)
			}
		})
	}
}

func TestIsEmpty(t *testing.T) {
	env := Environment{
		"NORMAL": "value",
		"EMPTY":  "",
		"SPACES": "   ",
	}

	tests := []struct {
		name string
		key  string
		want bool
	}{
		{name: "normal value", key: "NORMAL", want: false},
		{name: "empty value", key: "EMPTY", want: true},
		{name: "whitespace only", key: "SPACES", want: true},
		{name: "non-existent key", key: "NONEXISTENT", want: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := env.IsEmpty(tc.key)
			if got != tc.want {
				t.Errorf("IsEmpty(%q) = %v, want %v", tc.key, got, tc.want)
			}
		})
	}
}

func TestRenderTemplate(t *testing.T) {
	env := Environment{
		"NAME":  "World",
		"FLAG":  "true",
		"COUNT": "3",
		"LIST":  "1,2,3",
		"ITEMS": "a,b,c",
	}

	tests := []struct {
		name           string
		templateString string
		want           string
		wantErr        bool
	}{
		{
			name:           "simple variable",
			templateString: "Hello {{.NAME}}!",
			want:           "Hello World!",
			wantErr:        false,
		},
		{
			name:           "asString function",
			templateString: "Name: {{asString \"NAME\"}}",
			want:           "Name: World",
			wantErr:        false,
		},
		{
			name:           "asStringOr function",
			templateString: "Missing: {{asStringOr \"MISSING\" \"default\"}}",
			want:           "Missing: default",
			wantErr:        false,
		},
		{
			name:           "asStringSlice function",
			templateString: "Items: {{range asStringSlice \"ITEMS\" \",\"}}{{.}} {{end}}",
			want:           "Items: a b c ",
			wantErr:        false,
		},
		{
			name:           "asBool function",
			templateString: "Flag is {{asBool \"FLAG\"}}",
			want:           "Flag is true",
			wantErr:        false,
		},
		{
			name:           "asInt function",
			templateString: "Count: {{asInt \"COUNT\"}}",
			want:           "Count: 3",
			wantErr:        false,
		},
		{
			name:           "asIntSlice function",
			templateString: "Numbers: {{range asIntSlice \"LIST\" \",\"}}{{.}} {{end}}",
			want:           "Numbers: 1 2 3 ",
			wantErr:        false,
		},
		{
			name:           "isEmpty function",
			templateString: "Name empty: {{isEmpty \"NAME\"}}",
			want:           "Name empty: false",
			wantErr:        false,
		},
		{
			name:           "seq function",
			templateString: "{{range seq 1 3}}{{.}} {{end}}",
			want:           "1 2 3 ",
			wantErr:        false,
		},
		{
			name:           "invalid template",
			templateString: "{{if .NAME}}Incomplete",
			want:           "",
			wantErr:        true,
		},
		{
			name:           "error in template execution",
			templateString: "{{asInt \"INVALID\"}}",
			want:           "",
			wantErr:        true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := RenderTemplate(tc.templateString, env)

			if (err != nil) != tc.wantErr {
				t.Errorf("RenderTemplate() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr && got != tc.want {
				t.Errorf("RenderTemplate() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestRenderTemplateErrors(t *testing.T) {
	env := Environment{
		"NAME": "World",
	}

	// Test with a template that will definitely cause an error
	_, err := RenderTemplate("{{.NAME.NonExistentField}}", env)
	if err == nil {
		t.Errorf("Expected error for accessing non-existent field")
	} else if !strings.Contains(err.Error(), "executing template") {
		t.Errorf("Unexpected error message: %v", err)
	}

	// Test with a function that panics
	_, err = RenderTemplate("{{asInt \"INVALID\"}}", env)
	if err == nil {
		t.Errorf("Expected error when accessing non-existent variable with asInt")
	} else if !strings.Contains(err.Error(), "executing template") {
		t.Errorf("Unexpected error message: %v", err)
	}
}
