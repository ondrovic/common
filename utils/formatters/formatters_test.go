package formatters

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ondrovic/common/types"
)

// CreateNonEmptyDir creates a non-empty directory for testing.
func CreateNonEmptyDir(t *testing.T) string {
	t.Helper()
	dir := filepath.Join(t.TempDir(), FormatPath("non-empty-dir", runtime.GOOS))
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}

	// Create a file inside the non-empty directory
	filePath := filepath.Join(dir, "testfile.txt")
	file, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("failed to create file: %v", err)
	}
	defer file.Close() // Ensure file is closed before directory cleanup

	// Write some content to the file
	if _, err := file.WriteString("test content"); err != nil {
		t.Fatalf("failed to write to file: %v", err)
	}

	return dir
}

// TestGetVersion tests GetVersion func.
func TestGetVersion(t *testing.T) {
	type InputStruct struct {
		version  string
		fallback string
	}

	tests := []*types.TestLayout[InputStruct, string]{
		{Name: "Empty version returns fallback value", Input: InputStruct{version: "", fallback: "test-ver"}, Expected: "test-ver"},
		{Name: "Non-empty version ignores fallback, returns version", Input: InputStruct{version: "1.0.0", fallback: "test-ver"}, Expected: "1.0.0"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) { // Run each test case as a sub-test
			result := GetVersion(test.Input.version, test.Input.fallback)
			if result != test.Expected {
				t.Errorf("GetVersion() - %v(%q) = %q; expected %q", test.Name, test.Input, result, test.Expected)
			}
		})
	}
}

// TestPluralize tests Pluralize func.
func TestPluralize(t *testing.T) {
	type InputStruct struct {
		count    interface{}
		singular string
		plural   string
	}

	tests := []*types.TestLayout[InputStruct, string]{
		{Name: "Test Negative Count", Input: InputStruct{count: -1, singular: "", plural: ""}, Expected: "", Err: fmt.Errorf("count cannot be negative")},
		{Name: "Singular cannot be empty", Input: InputStruct{count: 0, singular: "", plural: "wants"}, Expected: "", Err: fmt.Errorf("singular and plural forms cannot be empty")},
		{Name: "Plural cannot be empty", Input: InputStruct{count: 0, singular: "want", plural: ""}, Expected: "", Err: fmt.Errorf("singular and plural forms cannot be empty")},
		{Name: "Singular case", Input: InputStruct{count: 1, singular: "want", plural: "wants"}, Expected: "want", Err: nil},
		{Name: "Plural case", Input: InputStruct{count: 2, singular: "want", plural: "wants"}, Expected: "wants", Err: nil},
		{Name: "Default case", Input: InputStruct{count: 1.25, singular: "want", plural: "wants"}, Expected: "", Err: fmt.Errorf("count must be an integer")},
		{Name: "Zero count", Input: InputStruct{count: 0, singular: "want", plural: "wanted"}, Expected: "want", Err: nil},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			result, err := Pluralize(test.Input.count, test.Input.singular, test.Input.plural)

			if result != test.Expected {
				t.Errorf("Pluralize(%q) - %v = %v; want %v", test.Input, test.Name, result, test.Expected)
			}

			if (err != nil && test.Err == nil) || (err == nil && test.Err != nil) || (err != nil && test.Err != nil && err.Error() != test.Err.Error()) {
				t.Errorf("Pluralize(%q) - %v error = %v; want %v", test.Input, test.Name, err, test.Err)
			}
		})
	}
}

// TestToLower tests ToLower func.
func TestToLower(t *testing.T) {
	expectedResult := "hello world"
	tests := []*types.TestLayout[interface{}, string]{
		{Name: "Test all uppercase", Input: "HELLO WORLD", Expected: expectedResult},
		{Name: "Test title case", Input: "Hello world", Expected: expectedResult},
		{Name: "Test camel case", Input: "hello World", Expected: expectedResult},
		{Name: "Test lower case", Input: "hello world", Expected: expectedResult},
		{Name: "Test error", Input: 52, Expected: "", Err: fmt.Errorf("input is not a string")},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			result, err := ToLower(test.Input)

			if result != test.Expected {
				t.Errorf("ToLower(%q) - %v = %q; expected %q", test.Input, test.Name, result, test.Expected)
			}

			if (err != nil && test.Err == nil) || (err == nil && test.Err != nil) || (err != nil && test.Err != nil && err.Error() != test.Err.Error()) {
				t.Errorf("ToLower(%q) - %v = %q; expected %q", test.Input, test.Name, err, test.Err)
			}
		})
	}
}

// TestToUpper tests ToUpper func.
func TestToUpper(t *testing.T) {
	expectedResult := "HELLO WORLD"
	tests := []*types.TestLayout[interface{}, string]{
		{Name: "Test all uppercase", Input: "HELLO WORLD", Expected: expectedResult},
		{Name: "Test title case", Input: "Hello world", Expected: expectedResult},
		{Name: "Test camel case", Input: "hello World", Expected: expectedResult},
		{Name: "Test lower case", Input: "hello world", Expected: expectedResult},
		{Name: "Test error", Input: 52, Expected: "", Err: fmt.Errorf("input is not a string")},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			result, err := ToUpper(test.Input)

			if result != test.Expected {
				t.Errorf("ToUpper(%q) - %v = %q; expected %q", test.Input, test.Name, result, test.Expected)
			}

			if (err != nil && test.Err == nil) || (err == nil && test.Err != nil) || (err != nil && test.Err != nil && err.Error() != test.Err.Error()) {
				t.Errorf("ToUpper(%q) - %v = %q; expected %q", test.Input, test.Name, err, test.Err)
			}
		})
	}
}

// TestContains tests Contains func.
func TestContains(t *testing.T) {
	type InputStruct struct {
		s      string
		subStr interface{}
	}

	testArrayOpts := []string{"a", "b", "c", "yes", ""}
	testStrOpts := "yes"
	tests := []*types.TestLayout[InputStruct, bool]{
		{Name: "Test subStr not a string", Input: InputStruct{s: "yes", subStr: 1}, Expected: false, Err: fmt.Errorf("substring must be a string or a slice of strings")},
		{Name: "Test str cannot be empty", Input: InputStruct{s: "", subStr: "yes"}, Expected: false, Err: fmt.Errorf("string cannot be empty")},
		{Name: "Test subStr cannot be empty", Input: InputStruct{s: "yes", subStr: ""}, Expected: false, Err: fmt.Errorf("substring cannot be empty")},

		{Name: "Test str found in subStr - single string", Input: InputStruct{s: "yes", subStr: testStrOpts}, Expected: true},
		{Name: "Test str found in subStr - string array", Input: InputStruct{s: "yes", subStr: testArrayOpts}, Expected: true},
		{Name: "Test str not found in subStr", Input: InputStruct{s: "no", subStr: testArrayOpts}, Expected: false},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			result, err := Contains(test.Input.s, test.Input.subStr)

			if result != test.Expected {
				t.Errorf("Contains(%v, %v) - %v = %v; expected %v", test.Input.s, test.Input.subStr, test.Name, result, test.Expected)
			}

			if (err != nil && test.Err == nil) || (err == nil && test.Err != nil) || (err != nil && test.Err != nil && err.Error() != test.Err.Error()) {
				// t.Errorf("ToLower(%q) - %v = %q; expected %q", test.Input, test.Name, err, test.Err)
				t.Errorf("Contains(%v, %v) - %v = %v; expected %v", test.Input.s, test.Input.subStr, test.Name, err, test.Err)
			}
		})
	}
}

// TestFormatPath tests FormatPath func.
func TestFormatPath(t *testing.T) {
	type InputStruct struct {
		Path string
		GOOS string
	}

	var dir = CreateNonEmptyDir(t)
	var filePath = filepath.Join(dir, "testfile.txt")
	tests := []*types.TestLayout[InputStruct, string]{
		{Name: "Test Linux path", Input: InputStruct{Path: filePath, GOOS: "linux"}, Expected: FormatPath(filePath, "linux"), Err: nil},
		{Name: "Test Windows path", Input: InputStruct{Path: filePath, GOOS: "windows"}, Expected: FormatPath(filePath, "windows"), Err: nil},
		{Name: "Test Default path", Input: InputStruct{Path: filePath, GOOS: "unknown"}, Expected: FormatPath(filePath, "unknown"), Err: nil},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) { // Run each test case as a sub-test
			result := FormatPath(test.Input.Path, test.Input.GOOS)
			if result != test.Expected {
				t.Errorf("FormatPath(%q, %q) - %v = %q; expected %q", test.Input.Path, test.Input.GOOS, test.Name, result, test.Expected)
			}
		})
	}
}

// TestFormatSize tests FormatSize func.
func TestFormatSize(t *testing.T) {
	tests := []*types.TestLayout[int64, string]{
		{Name: "Test 0 B", Input: 0, Expected: "0 B", Err: nil},
		{Name: "Test 1 B", Input: 1, Expected: "1.00 B", Err: nil},
		{Name: "Test 2 B", Input: 2, Expected: "2.00 B", Err: nil},
		{Name: "Test 1 KB", Input: 1024, Expected: "1.00 KB", Err: nil},
		{Name: "Test 2 KB", Input: 2048, Expected: "2.00 KB", Err: nil},
		{Name: "Test 1 MB", Input: 1048576, Expected: "1.00 MB", Err: nil},
		{Name: "Test 2 MB", Input: 2097152, Expected: "2.00 MB", Err: nil},
		{Name: "Test 1 GB", Input: 1073741824, Expected: "1.00 GB", Err: nil},
		{Name: "Test 2 GB", Input: 2147483648, Expected: "2.00 GB", Err: nil},
		{Name: "Test 1 TB", Input: 1099511627776, Expected: "1.00 TB", Err: nil},
		{Name: "Test 2 TB", Input: 2199023255552, Expected: "2.00 TB", Err: nil},
		{Name: "Test 1 PB", Input: 1125899906842624, Expected: "1.00 PB", Err: nil},
		{Name: "Test 2 PB", Input: 2251799813685248, Expected: "2.00 PB", Err: nil},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) { // Run each test case as a sub-test
			result := FormatSize(test.Input)
			if result != test.Expected {
				t.Errorf("FormatSize(%q) - %v = %q; expected %q", test.Input, test.Name, result, test.Expected)
			}
		})
	}
}
