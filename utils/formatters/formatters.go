package formatters

import (
	"errors"
	"math"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/ondrovic/common/types"
	"github.com/pterm/pterm"
)

// ToLower converts an interface{} to a lowercase string.
// It returns an error if the input is not a string.
func ToLower(i interface{}) (string, error) {
	str, ok := i.(string)
	if !ok {
		return "", errors.New("input is not a string")
	}
	return strings.ToLower(str), nil
}

// ToUpper converts an interface{} to a uppercase string.
// It returns an error if the input is not a string.
func ToUpper(i interface{}) (string, error) {
	str, ok := i.(string)
	if !ok {
		return "", errors.New("input is not a string")
	}
	return strings.ToUpper(str), nil
}

// Contains checks if a string contains another substring or any substring from a slice of strings.
// It returns an error if the main string is empty or if subStr is neither a string nor a slice of strings.
func Contains(s string, subStr interface{}) (bool, error) {
	if s == "" {
		return false, errors.New("string cannot be empty")
	}

	switch sub := subStr.(type) {
	case string:
		if sub == "" {
			return false, errors.New("substring cannot be empty")
		}
		return strings.Contains(s, sub), nil
	case []string:
		for _, str := range sub {
			if str == "" {
				continue // Skip empty substrings in the slice
			}
			if strings.Contains(s, str) {
				return true, nil
			}
		}
		return false, nil
	default:
		return false, errors.New("substring must be a string or a slice of strings")
	}
}

// The Pluralize function takes a count and returns the singular or plural form of a word based on the.
// count.
func Pluralize(count interface{}, singular, plural string) (string, error) {
	// Validate that count is an integer type
	switch v := reflect.ValueOf(count); v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// Continue with the logic only if count is a valid integer
		if v.Int() < 0 {
			return "", errors.New("count cannot be negative")
		}
		if singular == "" || plural == "" {
			return "", errors.New("singular and plural forms cannot be empty")
		}
		if v.Int() <= 1 {
			return singular, nil
		}
		return plural, nil
	default:
		return "", errors.New("count must be an integer")
	}
}

// The `FormatPath` function converts file paths to either Windows or Unix style based on the operating.
// system specified.
func FormatPath(path, goos string) string {
	switch goos {
	case "windows":
		// Convert to Windows style paths (with backslashes)
		return filepath.FromSlash(path)
	case "linux", "darwin":
		// Convert to Unix style paths (with forward slashes)
		return filepath.ToSlash(path)
	default:
		// Default to Unix style paths
		return path
	}
}

// The `FormatSize` function converts a given size in bytes to a human-readable format with appropriate.
// units.
func FormatSize(bytes int64) string {
	for _, unit := range types.SizeUnits {
		if bytes >= unit.Size {
			value := float64(bytes) / float64(unit.Size)
			// Round the value to two decimal places
			roundedValue := math.Round(value*100) / 100
			return pterm.Sprintf("%.2f %s", roundedValue, unit.Label)
		}
	}

	return "0 B"
}

// The function `GetVersion` is used for setting the version.
func GetVersion(version, fallback string) string {
	if version == "" {
		return fallback
	}

	return version
}
