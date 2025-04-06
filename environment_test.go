package main

import (
	"maps"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"testing"
	"time"
)

func TestNewEnvironment(t *testing.T) {
	envMap := map[string]string{"KEY1": "value1", "KEY2": "value2"}
	env := NewEnvironment(envMap)
	same := env.All()
	keys := slices.Collect(maps.Keys(env))
	sameKeys := slices.Collect(maps.Keys(env.SortAll()))
	if len(keys) != len(sameKeys) {
		t.Errorf("NewEnvironment did not create the expected keys")
	}

	if !reflect.DeepEqual(env, Environment(same)) {
		t.Errorf("NewEnvironment did not create the expected environment")
	}

	funcMap := GetTemplateFunctions(env)
	if funcMap == nil {
		t.Errorf("GetTemplateFunctions returned nil")
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
		trimChars string
		want      []string
		wantPanic bool
	}{
		{name: "trim spaces", key: "SPACES", delimiter: ",", trimChars: " ", want: []string{"a", "b", "c"}, wantPanic: false},
		{name: "trim quotes", key: "QUOTES", delimiter: ",", trimChars: `"`, want: []string{"x", "y", "z"}, wantPanic: false},
		{name: "non-existent key", key: "NONEXISTENT", delimiter: ",", trimChars: " ", want: nil, wantPanic: true},
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

			got := env.AsStringSliceTrim(tc.key, tc.delimiter, tc.trimChars)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("AsStringSliceTrim(%q, %q, %q) = %v, want %v",
					tc.key, tc.delimiter, tc.trimChars, got, tc.want)
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

func TestAsURL(t *testing.T) {
	env := Environment{
		"VALID_URL":   "http://example.com",
		"INVALID_URL": "invalid-url",
	}

	tests := []struct {
		name      string
		key       string
		want      string
		wantPanic bool
	}{
		{name: "not existing key", key: "NONEXISTENT", want: "", wantPanic: true},
		{name: "valid URL", key: "VALID_URL", want: "http://example.com", wantPanic: false},
		{name: "invalid URL", key: "INVALID_URL", want: "", wantPanic: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("AsURL did not panic for key %s", tc.key)
					}
				}()
			}

			got := env.AsURL(tc.key)
			if got != tc.want {
				t.Errorf("AsURL(%q) = %q, want %q", tc.key, got, tc.want)
			}
		})
	}
}

func TestAsHostPort(t *testing.T) {
	env := Environment{
		"VALID_HOST_PORT":    "localhost:8080",
		"INVALID_HOST_PORT1": "invalid-host-port:999999",
		"INVALID_HOST_PORT2": "/@#$%^&*():88",
	}

	tests := []struct {
		name      string
		key       string
		want      string
		wantPanic bool
	}{
		{name: "not existing key", key: "NONEXISTENT", want: "", wantPanic: true},
		{name: "valid host:port", key: "VALID_HOST_PORT", want: "localhost:8080", wantPanic: false},
		{name: "invalid host:port", key: "INVALID_HOST_PORT1", want: "", wantPanic: true},
		{name: "invalid host:port with special chars", key: "INVALID_HOST_PORT2", want: "", wantPanic: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("AsHostPort did not panic for key %s", tc.key)
					}
				}()
			}

			got := env.AsHostPort(tc.key)
			if got != tc.want {
				t.Errorf("AsHostPort(%q) = %q, want %q", tc.key, got, tc.want)
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

func TestAsFloat(t *testing.T) {
	env := Environment{
		"POSITIVE": "123.123",
		"NEGATIVE": "-42.24",
		"ZERO":     "0",
		"INVALID":  "not-an-int",
	}

	tests := []struct {
		name      string
		key       string
		want      float64
		wantPanic bool
	}{
		{name: "positive number", key: "POSITIVE", want: 123.123, wantPanic: false},
		{name: "negative number", key: "NEGATIVE", want: -42.24, wantPanic: false},
		{name: "zero", key: "ZERO", want: 0, wantPanic: false},
		{name: "invalid value", key: "INVALID", want: 0, wantPanic: true},
		{name: "non-existent key", key: "NONEXISTENT", want: 0, wantPanic: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("AsFloat did not panic for key %s", tc.key)
					}
				}()
			}

			got := env.AsFloat(tc.key)
			if got != tc.want {
				t.Errorf("AsFloat(%q) = %v, want %v", tc.key, got, tc.want)
			}
		})
	}
}

func TestAsFloatOr(t *testing.T) {
	env := Environment{
		"VALID":   "123.345",
		"INVALID": "not-an-int",
	}

	tests := []struct {
		name         string
		key          string
		defaultValue float64
		want         float64
	}{
		{name: "existing valid", key: "VALID", defaultValue: 999.1, want: 123.345},
		{name: "existing invalid", key: "INVALID", defaultValue: 999.1, want: 999.1},
		{name: "non-existent key", key: "NONEXISTENT", defaultValue: 999.1, want: 999.1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := env.AsFloatOr(tc.key, tc.defaultValue)
			if got != tc.want {
				t.Errorf("AsFloatOr(%q, %v) = %v, want %v", tc.key, tc.defaultValue, got, tc.want)
			}
		})
	}
}

func TestAsFloatSlice(t *testing.T) {
	env := Environment{
		"VALID":  "1.1,2.2,3.3",
		"SPACES": " 4.4, 5.5, 6.6 ",
		"MIXED":  "7.7, invalid, 9.9",
	}

	tests := []struct {
		name      string
		key       string
		delimiter string
		want      []float64
		wantPanic bool
	}{
		{name: "valid values", key: "VALID", delimiter: ",", want: []float64{1.1, 2.2, 3.3}, wantPanic: false},
		{name: "values with spaces", key: "SPACES", delimiter: ",", want: []float64{4.4, 5.5, 6.6}, wantPanic: false},
		{name: "mixed valid and invalid", key: "MIXED", delimiter: ",", want: nil, wantPanic: true},
		{name: "non-existent key", key: "NONEXISTENT", delimiter: ",", want: nil, wantPanic: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("AsFloatSlice did not panic for key %s", tc.key)
					}
				}()
			}

			got := env.AsFloatSlice(tc.key, tc.delimiter)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("AsFloatSlice(%q, %q) = %v, want %v", tc.key, tc.delimiter, got, tc.want)
			}
		})
	}
}

func TestAsPort(t *testing.T) {
	env := Environment{
		"VALID_PORT":    "8080",
		"INVALID_PORT1": "invalid-port",
		"INVALID_PORT2": "999999",
	}

	tests := []struct {
		name      string
		key       string
		want      int
		wantPanic bool
	}{
		{name: "valid port", key: "VALID_PORT", want: 8080, wantPanic: false},
		{name: "invalid port", key: "INVALID_PORT1", want: 0, wantPanic: true},
		{name: "out of range port", key: "INVALID_PORT2", want: 0, wantPanic: true},
		{name: "non-existent key", key: "NONEXISTENT", want: 0, wantPanic: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("AsPort did not panic for key %s", tc.key)
					}
				}()
			}

			got := env.AsPort(tc.key)
			if got != tc.want {
				t.Errorf("AsPort(%q) = %d, want %d", tc.key, got, tc.want)
			}
		})
	}
}

func TestAsPortOr(t *testing.T) {
	env := Environment{
		"VALID_PORT":   "8080",
		"INVALID_PORT": "invalid-port",
		"OUT_OF_RANGE": "999999",
	}

	tests := []struct {
		name         string
		key          string
		defaultValue int
		want         int
	}{
		{name: "existing valid", key: "VALID_PORT", defaultValue: 999, want: 8080},
		{name: "existing invalid", key: "INVALID_PORT", defaultValue: 999, want: 999},
		{name: "out of range", key: "OUT_OF_RANGE", defaultValue: 999, want: 999},
		{name: "non-existent key", key: "NONEXISTENT", defaultValue: 999, want: 999},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := env.AsPortOr(tc.key, tc.defaultValue)
			if got != tc.want {
				t.Errorf("AsPortOr(%q, %v) = %v, want %v", tc.key, tc.defaultValue, got, tc.want)
			}
		})
	}
}

func Test_isEmpty(t *testing.T) {
	tests := []struct {
		name   string
		value  string
		wanted bool
	}{
		{name: "empty string", value: "", wanted: true},
		{name: "non-empty string", value: "not empty", wanted: false},
		{name: "whitespace string", value: "   ", wanted: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := isEmpty(tc.value)
			if got != tc.wanted {
				t.Errorf("isEmpty(%q) = %v, want %v", tc.value, got, tc.wanted)
			}
		})
	}
}

func Test_contains(t *testing.T) {
	tests := []struct {
		name   string
		value  string
		substr string
		wanted bool
	}{
		{name: "substring present", value: "hello world", substr: "world", wanted: true},
		{name: "substring not present", value: "hello world", substr: "foo", wanted: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := contains(tc.value, tc.substr)
			if got != tc.wanted {
				t.Errorf("contains(%q, %q) = %v, want %v", tc.value, tc.substr, got, tc.wanted)
			}
		})
	}
}

func Test_containsCaseInsensitive(t *testing.T) {
	tests := []struct {
		name   string
		value  string
		substr string
		wanted bool
	}{
		{name: "substring present", value: "hello world", substr: "WORLD", wanted: true},
		{name: "substring not present", value: "hello world", substr: "foo", wanted: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := containsCaseInsensitive(tc.value, tc.substr)
			if got != tc.wanted {
				t.Errorf("containsCaseInsensitive(%q, %q) = %v, want %v", tc.value, tc.substr, got, tc.wanted)
			}
		})
	}
}

func Test_containsAny(t *testing.T) {
	tests := []struct {
		name   string
		value  string
		chars  string
		wanted bool
	}{
		{name: "substring present", value: "hello world", chars: "w", wanted: true},
		{name: "substring not present", value: "hello world", chars: "z.u", wanted: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := containsAny(tc.value, tc.chars)
			if got != tc.wanted {
				t.Errorf("containsAny(%q, %q) = %v, want %v", tc.value, tc.chars, got, tc.wanted)
			}
		})
	}
}

func Test_hasPrefix(t *testing.T) {
	tests := []struct {
		name   string
		value  string
		prefix string
		wanted bool
	}{
		{name: "prefix present", value: "hello world", prefix: "hello", wanted: true},
		{name: "prefix not present", value: "hello world", prefix: "world", wanted: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := hasPrefix(tc.value, tc.prefix)
			if got != tc.wanted {
				t.Errorf("hasPrefix(%q, %q) = %v, want %v", tc.value, tc.prefix, got, tc.wanted)
			}
		})
	}
}

func Test_hasSuffix(t *testing.T) {
	tests := []struct {
		name   string
		value  string
		suffix string
		wanted bool
	}{
		{name: "suffix present", value: "hello world", suffix: "world", wanted: true},
		{name: "suffix not present", value: "hello world", suffix: "hello", wanted: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := hasSuffix(tc.value, tc.suffix)
			if got != tc.wanted {
				t.Errorf("hasSuffix(%q, %q) = %v, want %v", tc.value, tc.suffix, got, tc.wanted)
			}
		})
	}
}

func Test_toLower(t *testing.T) {
	tests := []struct {
		name   string
		value  string
		wanted string
	}{
		{name: "lowercase", value: "hello", wanted: "hello"},
		{name: "uppercase", value: "HELLO", wanted: "hello"},
		{name: "mixed case", value: "HeLLo", wanted: "hello"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := toLower(tc.value)
			if got != tc.wanted {
				t.Errorf("toLower(%q) = %q, want %q", tc.value, got, tc.wanted)
			}
		})
	}
}

func Test_toUpper(t *testing.T) {
	tests := []struct {
		name   string
		value  string
		wanted string
	}{
		{name: "lowercase", value: "hello", wanted: "HELLO"},
		{name: "uppercase", value: "HELLO", wanted: "HELLO"},
		{name: "mixed case", value: "HeLLo", wanted: "HELLO"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := toUpper(tc.value)
			if got != tc.wanted {
				t.Errorf("toUpper(%q) = %q, want %q", tc.value, got, tc.wanted)
			}
		})
	}
}

func Test_trim(t *testing.T) {
	tests := []struct {
		name   string
		value  string
		cutset string
		wanted string
	}{
		{name: "trim spaces", value: "  hello  ", cutset: " ", wanted: "hello"},
		{name: "trim asterisks", value: "***hello***", cutset: "*", wanted: "hello"},
		{name: "trim mixed", value: "  ***hello***  ", cutset: " *", wanted: "hello"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := trim(tc.value, tc.cutset)
			if got != tc.wanted {
				t.Errorf("trim(%q, %q) = %q, want %q", tc.value, tc.cutset, got, tc.wanted)
			}
		})
	}
}

func Test_trimLeft(t *testing.T) {
	tests := []struct {
		name   string
		value  string
		cutset string
		wanted string
	}{
		{name: "trim spaces", value: "  hello  ", cutset: " ", wanted: "hello  "},
		{name: "trim asterisks", value: "***hello***", cutset: "*", wanted: "hello***"},
		{name: "trim mixed", value: "  ***hello***  ", cutset: " *", wanted: "hello***  "},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := trimLeft(tc.value, tc.cutset)
			if got != tc.wanted {
				t.Errorf("trimLeft(%q, %q) = %q, want %q", tc.value, tc.cutset, got, tc.wanted)
			}
		})
	}
}

func Test_trimRight(t *testing.T) {
	tests := []struct {
		name   string
		value  string
		cutset string
		wanted string
	}{
		{name: "trim spaces", value: "  hello  ", cutset: " ", wanted: "  hello"},
		{name: "trim asterisks", value: "***hello***", cutset: "*", wanted: "***hello"},
		{name: "trim mixed", value: "  ***hello***  ", cutset: " *", wanted: "  ***hello"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := trimRight(tc.value, tc.cutset)
			if got != tc.wanted {
				t.Errorf("trimRight(%q, %q) = %q, want %q", tc.value, tc.cutset, got, tc.wanted)
			}
		})
	}
}

func Test_trimSpace(t *testing.T) {
	tests := []struct {
		name   string
		value  string
		wanted string
	}{
		{name: "trim spaces", value: "  hello  ", wanted: "hello"},
		{name: "trim asterisks", value: "***hello***", wanted: "***hello***"},
		{name: "trim mixed", value: "  ***hello***  ", wanted: "***hello***"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := trimSpace(tc.value)
			if got != tc.wanted {
				t.Errorf("trimSpace(%q) = %q, want %q", tc.value, got, tc.wanted)
			}
		})
	}
}

func Test_base64Encode(t *testing.T) {
	tests := []struct {
		name   string
		value  string
		wanted string
	}{
		{name: "simple string", value: "hello", wanted: "aGVsbG8="},
		{name: "empty string", value: "", wanted: ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := base64Encode(tc.value)
			if got != tc.wanted {
				t.Errorf("base64Encode(%q) = %q, want %q", tc.value, got, tc.wanted)
			}
		})
	}
}

func Test_base64Decode(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		wanted    string
		wantPanic bool
	}{
		{name: "valid base64", value: "aGVsbG8=", wanted: "hello", wantPanic: false},
		{name: "invalid base64", value: "invalid-base64", wanted: "", wantPanic: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("base64Decode did not panic for value %s", tc.value)
					}
				}()
			}

			got := base64Decode(tc.value)
			if got != tc.wanted {
				t.Errorf("base64Decode(%q) = %q, want %q", tc.value, got, tc.wanted)
			}
		})
	}
}

func Test_hash(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		algorithm string // md5, sha1, sha224, sha256, sha512 otherwise panic
		wantPanic bool
	}{
		{name: "valid md5", input: "hello", algorithm: "md5", wantPanic: false},
		{name: "valid sha1", input: "hello", algorithm: "sha1", wantPanic: false},
		{name: "valid sha224", input: "hello", algorithm: "sha224", wantPanic: false},
		{name: "valid sha256", input: "hello", algorithm: "sha256", wantPanic: false},
		{name: "valid sha512", input: "hello", algorithm: "sha512", wantPanic: false},
		{name: "invalid algorithm", input: "hello", algorithm: "invalid", wantPanic: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("hash did not panic for algorithm %s", tc.algorithm)
					}
				}()
			}

			got := hash(tc.input, tc.algorithm)
			if got == "" {
				t.Errorf("hash(%q, %q) = empty string", tc.input, tc.algorithm)
			}
		})
	}
}

func Test_sequence(t *testing.T) {
	se := sequence(1, 10)
	if len(se) != 10 {
		t.Errorf("expected length of sequence to be 10, got %d", len(se))
	}
	if sequence(10, 1) != nil {
		t.Errorf("expected sequence(10, 1) to be nil")
	}
}

func Test_fileExistOrDefault(t *testing.T) {

	t.Run("destination file exists", func(t *testing.T) {
		destination := t.TempDir() + "/testfile.txt"
		os.WriteFile(destination, []byte("test"), 0644)
		fileExistOrDefault(destination, "/no/matter/what.txt")
	})

	t.Run("destination file not exist exists", func(t *testing.T) {
		dir := t.TempDir()
		destination := filepath.Join(dir, "testfile.txt")
		os.WriteFile(destination, []byte("will be deleted"), 0644)
		os.Remove(destination)

		testString := time.Now().String()

		defaultPath := filepath.Join(dir, "defaultfile.txt")
		f, e := os.OpenFile(defaultPath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0777)
		if e != nil {
			t.Fatalf("failed to create default file: %v", e)
		}
		f.Write([]byte(testString))
		f.Close()

		fileInfo, err := os.Stat(defaultPath)
		t.Logf("permissions of defaultPath is: %o", fileInfo.Mode())
		if err != nil {
			t.Fatalf("failed to stat default file: %v", err)
		}

		fileExistOrDefault(destination, defaultPath)
		// read file data
		data, err := os.ReadFile(destination)
		if err != nil {
			t.Fatalf("failed to read file: %v", err)
		}
		if string(data) != testString {
			t.Fatalf("expected file content to be '%s', got '%s'", string(data), testString)
		}

		// check file permissions
		info, err := os.Stat(destination)
		if err != nil {
			t.Fatalf("failed to stat file: %v", err)
		}
		if info.Mode() != fileInfo.Mode() {
			t.Fatalf("expected file permissions to be %o, got %o", fileInfo.Mode(), info.Mode())
		}
	})

	t.Run("destination file not exist exists", func(t *testing.T) {
		dir := t.TempDir()
		destination := filepath.Join(dir, "testfile.txt")
		defaultPath := filepath.Join(dir, "defaultPath.txt")
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("fileExistOrDefault did not panic for non-existent destination file")
			}
		}()

		fileExistOrDefault(destination, defaultPath)
	})
}
