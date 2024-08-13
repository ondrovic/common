package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/ondrovic/common/types"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// The MockCmd type is a struct used for mocking commands in Go code.
// @property {error} err - The `err` property in the `MockCmd` struct is a field that holds an error
// value. It can be used to store an error that may occur during the execution of a command or
// operation.
type MockCmd struct {
	err error
}

// The above code snippet is defining a method named `Run` for a struct type `MockCmd`. This method
// takes a pointer receiver `m` of type `MockCmd` and returns an error. Inside the method, it simply
// returns the error `m.err`. This method is likely intended to simulate running a command and
// returning an error if any.
func (m *MockCmd) Run() error {
	return m.err
}

// The MockDirOps type is used for mocking directory operations in Go code.
// @property {error} readDirErr - The `readDirErr` property in the `MockDirOps` struct is used to store
// an error that may occur when atestempting to read a directory. This error could be related to issues
// such as permission problems, directory not found, or any other error that may occur during the
// directory reading operation.
// @property {error} removeErr - The `removeErr` property in the `MockDirOps` struct is used to store
// an error that may occur when atestempting to remove a directory. This error could be related to
// permissions, file system issues, or any other problem that prevents the directory from being removed
// successfully.
type MockDirOps struct {
	readDirErr error
	removeErr  error
}

func (m *MockDirOps) ReadDir(_ string) ([]os.DirEntry, error) {
	if m.readDirErr != nil {
		return nil, m.readDirErr
	}
	return []os.DirEntry{}, nil
}

func (m *MockDirOps) Remove(_ string) error {
	return m.removeErr
}

// CreateTempFile creates a temporary file for testing.
func CreateTempFile(t *testing.T) *os.File {
	t.Helper()
	file, err := os.CreateTemp("", "testfile-*.txt")
	if err != nil {
		t.Fatalf("failed to create file: %v", err)
	}
	return file
}

// CreateEmptyDir creates an empty directory for testing.
func CreateEmptyDir(t *testing.T) string {
	t.Helper()
	dir := FormatPath(filepath.Join(t.TempDir(), "empty-dir"), runtime.GOOS)
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}
	return dir
}

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

// TestValidateStringField tests the ValidateStringField func.
func TestValidateStringField(t *testing.T) {
	type InputStruct struct {
		fieldName string
		value     string
	}
	tests := []*types.TestLayout[InputStruct, error]{
		{Name: "Empty field", Input: InputStruct{fieldName: "TestField", value: ""}, Expected: fmt.Errorf("TestField cannot be empty")},
		{Name: "Non-empty field", Input: InputStruct{fieldName: "TestField", value: "value"}},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			err := validateStringField(test.Input.fieldName, test.Input.value)
			if test.Expected == nil {
				if err != nil {
					t.Errorf("validateStringField() - %v error = %v, expected nil", test.Name, err)
				}
			} else {
				if err == nil || err.Error() != test.Expected.Error() {
					t.Errorf("validateStringField() - %v error = %v, expected %v", test.Name, err, test.Expected)
				}
			}
		})
	}
}

// TestValidateStructField tests ValidateStructField func.
func TestValidateStructField(t *testing.T) {
	type InputStruct struct {
		fieldName string
		value     interface{}
	}
	tests := []*types.TestLayout[InputStruct, error]{
		{Name: "Empty struct", Input: InputStruct{fieldName: "TestStruct", value: struct{}{}}, Expected: fmt.Errorf("TestStruct cannot be an empty struct")},
		{Name: "Non-empty struct", Input: InputStruct{fieldName: "TestStruct", value: struct{ Field string }{Field: "value"}}},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			err := validateStructField(test.Input.fieldName, reflect.ValueOf(test.Input.value))
			if test.Expected == nil {
				if err != nil {
					t.Errorf("validateStructField() - %v error = %v, expected nil", test.Name, err)
				}
			} else {
				if err == nil || err.Error() != test.Expected.Error() {
					t.Errorf("validateStructField() - %v error = %v, expected %v", test.Name, err, test.Expected)
				}
			}
		})
	}
}

// TestApplicationBanner tests ApplicationBanner func.
func TestApplicationBanner(t *testing.T) {
	type InputStruct struct {
		app         *types.Application
		clearScreen func(string) error
	}

	tests := []*types.TestLayout[InputStruct, error]{
		{
			Name: "ClearTerminalScreen error",
			Input: InputStruct{
				app: &types.Application{
					Name:        "Test App",
					Description: "Test Description",
					Style:       types.Styles{Color: types.Colors{Background: pterm.BgRed, Foreground: pterm.FgWhite}},
					Usage:       "Test Usage",
					Version:     "1.0.0",
				},
				clearScreen: func(string) error {
					return fmt.Errorf("simulated clear screen error")
				},
			},
			Expected: fmt.Errorf("simulated clear screen error"),
		},
		{Name: "Application is not a struct", Input: InputStruct{app: nil, clearScreen: ClearTerminalScreen}, Expected: fmt.Errorf("validateApp expects a struct")},
		{Name: "Empty application name", Input: InputStruct{app: &types.Application{Description: "Test Description", Style: types.Styles{Color: types.Colors{}}, Usage: "Test Usage", Version: "1.0.0"}, clearScreen: ClearTerminalScreen}, Expected: fmt.Errorf("Name cannot be empty")},
		{Name: "Empty application description", Input: InputStruct{app: &types.Application{Name: "Test App", Style: types.Styles{Color: types.Colors{}}, Usage: "Test Usage", Version: "1.0.0"}, clearScreen: ClearTerminalScreen}, Expected: fmt.Errorf("Description cannot be empty")},
		{Name: "Empty application style struct", Input: InputStruct{app: &types.Application{Name: "Test App", Description: "Test Description", Style: types.Styles{Color: types.Colors{}}, Usage: "Test Usage", Version: "1.0.0"}, clearScreen: ClearTerminalScreen}, Expected: fmt.Errorf("Style cannot be an empty struct")},
		{Name: "Empty application usage", Input: InputStruct{app: &types.Application{Name: "Test App", Description: "Test Description", Style: types.Styles{Color: types.Colors{Background: pterm.BgRed, Foreground: pterm.FgWhite}}, Version: "1.0.0"}, clearScreen: ClearTerminalScreen}, Expected: fmt.Errorf("Usage cannot be empty")},
		{Name: "Empty application version", Input: InputStruct{app: &types.Application{Name: "Test App", Description: "Test Description", Style: types.Styles{Color: types.Colors{Background: pterm.BgRed, Foreground: pterm.FgWhite}}, Usage: "Test Usage"}, clearScreen: ClearTerminalScreen}, Expected: fmt.Errorf("Version cannot be empty")},
		{Name: "Valid application", Input: InputStruct{app: &types.Application{Name: "Test App", Description: "Test Description", Style: types.Styles{Color: types.Colors{Background: pterm.BgRed, Foreground: pterm.FgWhite}}, Usage: "Test Usage", Version: "1.0.0"}, clearScreen: ClearTerminalScreen}, Expected: nil},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			err := ApplicationBanner(test.Input.app, test.Input.clearScreen)
			if test.Expected == nil {
				if err != nil {
					t.Errorf("ApplicationBanner() - %v error = %v, expected nil", test.Name, err)
				}
			} else {
				if err == nil || err.Error() != test.Expected.Error() {
					t.Errorf("ApplicationBanner() - %v error = %v, expected %v", test.Name, err, test.Expected)
				}
			}
		})
	}
}

// TestClearTerminalScreen tests the ClearTerminalScreen func.
func TestClearTerminalScreen(t *testing.T) {
	type ExpectedOutcome struct {
		shouldFail bool
		err        error
	}

	// helper to determine the fail and error based on os
	expectedResultsBasedOnOS := func(inputOS string) ExpectedOutcome {
		shouldFail := runtime.GOOS != inputOS

		if inputOS == "unknown" {
			return ExpectedOutcome{
				true,
				fmt.Errorf("unsupported platform: %s", "unknown"),
			}
		} else if shouldFail {
			return ExpectedOutcome{
				true,
				errors.New("failed to clear terminal exec"),
			}
		}

		return ExpectedOutcome{
			false,
			nil,
		}
	}

	tests := []*types.TestLayout[string, ExpectedOutcome]{
		// ExpectedOutcome is calculated before the test runs
		{Name: "Test Linux clear command", Input: "linux"},
		{Name: "Test macOS clear command", Input: "darwin"},
		{Name: "Test Windows clear command", Input: "windows"},
		{Name: "Test unsupported OS", Input: "unknown"},
		{Name: "Test Linux clear command failure", Input: "linux"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// Determine if the test should fail based on the current OS
			results := expectedResultsBasedOnOS(test.Input)

			test.Expected.shouldFail = results.shouldFail
			test.Expected.err = results.err

			err := ClearTerminalScreen(test.Input)
			switch {
			case err != nil && test.Expected.err == nil:
				t.Errorf("ClearTerminalScreen() - %v(%q) = %v; expected no error", test.Name, test.Input, err)
			case err == nil && test.Expected.err != nil:
				t.Errorf("ClearTerminalScreen() - %v(%q) = no error; expected %v", test.Name, test.Input, test.Expected.err)
			case err != nil && test.Expected.err != nil && !strings.Contains(err.Error(), test.Expected.err.Error()):
				t.Errorf("ClearTerminalScreen() - %v(%q) = %v; expected %v", test.Name, test.Input, err, test.Expected.err)
			}
		})
	}
}

// TestToFileType tests ToFileType func.
func TestToFileType(t *testing.T) {
	tests := []*types.TestLayout[string, types.FileType]{
		{Name: "Test any type", Input: "any", Expected: types.FileTypes.Any, Err: nil},
		{Name: "Test video type", Input: "video", Expected: types.FileTypes.Video, Err: nil},
		{Name: "Test image type", Input: "image", Expected: types.FileTypes.Image, Err: nil},
		{Name: "Test archive type", Input: "archive", Expected: types.FileTypes.Archive, Err: nil},
		{Name: "Test documents type", Input: "documents", Expected: types.FileTypes.Documents, Err: nil},
		{Name: "Test case-insensitive", Input: "ANY", Expected: types.FileTypes.Any, Err: nil}, // Case-insensitive check
		{Name: "Test mixed case", Input: "ViDeO", Expected: types.FileTypes.Video, Err: nil},   // Mixed case check
		{Name: "Test invalid input", Input: "invalid", Expected: "", Err: nil},                 // Invalid input check
		{Name: "Test empty input", Input: "", Expected: "", Err: nil},                          // Empty string check
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) { // Run each test case as a sub-test
			result := ToFileType(test.Input)
			if result != test.Expected {
				t.Errorf("TestToFileType() - %v(%q) = %q; expected %q", test.Name, test.Input, result, test.Expected)
			}
		})
	}
}

// TestToOperatorType tests ToOperatorType func.
func TestToOperatorType(t *testing.T) {
	tests := []*types.TestLayout[string, types.OperatorType]{
		{Name: "Test et", Input: "et", Expected: types.OperatorTypes.EqualTo, Err: nil},
		{Name: "Test equal to", Input: "equal to", Expected: types.OperatorTypes.EqualTo, Err: nil},
		{Name: `Test equalto`, Input: "equalto", Expected: types.OperatorTypes.EqualTo, Err: nil},
		{Name: "Test equal", Input: "equal", Expected: types.OperatorTypes.EqualTo, Err: nil},
		{Name: "Test ==", Input: "==", Expected: types.OperatorTypes.EqualTo, Err: nil},

		{Name: "Test gt", Input: "gt", Expected: types.OperatorTypes.GreaterThan, Err: nil},
		{Name: "Test greater than", Input: "greater than", Expected: types.OperatorTypes.GreaterThan, Err: nil},
		{Name: "Test greaterthan", Input: "greaterthan", Expected: types.OperatorTypes.GreaterThan, Err: nil},
		{Name: "Test >", Input: ">", Expected: types.OperatorTypes.GreaterThan, Err: nil},

		{Name: "Test gte", Input: "gte", Expected: types.OperatorTypes.GreaterThanEqualTo, Err: nil},
		{Name: "Test greater than or equal to", Input: "greater than or equal to", Expected: types.OperatorTypes.GreaterThanEqualTo, Err: nil},
		{Name: "Test greaterthanorequalto", Input: "greaterthanorequalto", Expected: types.OperatorTypes.GreaterThanEqualTo, Err: nil},
		{Name: "Test >=", Input: ">=", Expected: types.OperatorTypes.GreaterThanEqualTo, Err: nil},

		{Name: "Test lt", Input: "lt", Expected: types.OperatorTypes.LessThan, Err: nil},
		{Name: "Test less than", Input: "less than", Expected: types.OperatorTypes.LessThan, Err: nil},
		{Name: "Test lessthan", Input: "lessthan", Expected: types.OperatorTypes.LessThan, Err: nil},
		{Name: "Test <", Input: "<", Expected: types.OperatorTypes.LessThan, Err: nil},

		{Name: "Test lte", Input: "lte", Expected: types.OperatorTypes.LessThanEqualTo, Err: nil},
		{Name: "Test less than or equal to", Input: "less than or equal to", Expected: types.OperatorTypes.LessThanEqualTo, Err: nil},
		{Name: "Test lessthanorequalto", Input: "lessthanorequalto", Expected: types.OperatorTypes.LessThanEqualTo, Err: nil},
		{Name: "Test <=", Input: "<=", Expected: types.OperatorTypes.LessThanEqualTo, Err: nil},

		{Name: "Test default case", Input: "", Expected: "", Err: nil},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) { // Run each test case as a sub-test
			result := ToOperatorType(test.Input)
			if result != test.Expected {
				t.Errorf("ToOperatorType() - %v(%q) = %q; expected %q", test.Name, test.Input, result, test.Expected)
			}
		})
	}
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

// TestHandleCliFlags tests HandleCliFlags func.
func TestHandleCliFlags(t *testing.T) {
	type ExpectedOutcome struct {
		result bool
		err    error
	}
	tests := []*types.TestLayout[string, ExpectedOutcome]{
		{Name: "Help flag -h", Input: "-h", Expected: ExpectedOutcome{result: true, err: nil}},
		{Name: "Help flag --help", Input: "--help", Expected: ExpectedOutcome{result: true, err: nil}},
		{Name: "Version flag -v", Input: "-v", Expected: ExpectedOutcome{result: true, err: nil}},
		{Name: "Version flag --version", Input: "--version", Expected: ExpectedOutcome{result: true, err: nil}},
		{Name: "Any other flag", Input: "-o", Expected: ExpectedOutcome{result: false, err: nil}},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// Save and restore os.Args
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()

			// Set up test args
			os.Args = []string{"program", test.Input}

			// Create a mock cobra.Command
			cmd := &cobra.Command{
				Use: "test",
			}

			result, err := HandleCliFlags(cmd)

			if result != test.Expected.result {
				t.Errorf("HandleCliFlags() - %v(%q) result = %v; expected %v", test.Name, test.Input, result, test.Expected.result)
			}

			if (err == nil && test.Expected.err != nil) || (err != nil && test.Expected.err == nil) || (err != nil && err.Error() == test.Expected.err.Error()) {
				t.Errorf("HandleCliFlags() - %v(%q) error = %c; expected %v", test.Name, test.Input, err, test.Expected.err)
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

// TestFormatPath tests FormatPath func.
func TestFormatPath(t *testing.T) {
	type InputStruct struct {
		Path string
		GOOS string
	}
	tests := []*types.TestLayout[InputStruct, string]{
		{Name: "Test Windows path formatesting", Input: InputStruct{Path: `C:\path\to\file`, GOOS: "windows"}, Expected: `C:\path\to\file`, Err: nil},
		{Name: "Test Linux path formatesting", Input: InputStruct{Path: `\path\to\file`, GOOS: "linux"}, Expected: `/path/to/file`, Err: nil},
		{Name: "Test macOS path formatesting", Input: InputStruct{Path: `/path/to/file`, GOOS: "darwin"}, Expected: `/path/to/file`, Err: nil},
		{Name: "Test default case for Unix", Input: InputStruct{Path: `/path/to/file`, GOOS: "unknown"}, Expected: `/path/to/file`, Err: nil},
		{Name: "Test Windows path with forward slashes", Input: InputStruct{Path: `C:/path/to/file`, GOOS: "windows"}, Expected: `C:\path\to\file`, Err: nil},
		{Name: "Test Unix path with backslashes", Input: InputStruct{Path: `\path\to\file`, GOOS: "linux"}, Expected: `/path/to/file`, Err: nil},
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

// TestIsExtensionValid tests IsExtensionValid func.
func TestIsExtensionValid(t *testing.T) {
	type InputStruct struct {
		FileType types.FileType
		Path     string
	}
	tests := []*types.TestLayout[InputStruct, bool]{
		// Tests for Any file type (wildcard)
		{Name: "Any - valid extension", Input: InputStruct{FileType: types.FileTypes.Any, Path: "example.file"}, Expected: true},
		{Name: "Any - no extension", Input: InputStruct{FileType: types.FileTypes.Any, Path: "example"}, Expected: true},

		// Tests for Video file type
		{Name: "Video - valid extension .mp4", Input: InputStruct{FileType: types.FileTypes.Video, Path: "video.mp4"}, Expected: true},
		{Name: "Video - invalid extension .txt", Input: InputStruct{FileType: types.FileTypes.Video, Path: "document.txt"}, Expected: false},
		{Name: "Video - empty extension", Input: InputStruct{FileType: types.FileTypes.Video, Path: "video"}, Expected: false},

		// Tests for Image file type
		{Name: "Image - valid extension .jpg", Input: InputStruct{FileType: types.FileTypes.Image, Path: "picture.jpg"}, Expected: true},
		{Name: "Image - valid extension .png", Input: InputStruct{FileType: types.FileTypes.Image, Path: "picture.png"}, Expected: true},
		{Name: "Image - invalid extension .mp4", Input: InputStruct{FileType: types.FileTypes.Image, Path: "video.mp4"}, Expected: false},
		{Name: "Image - empty extension", Input: InputStruct{FileType: types.FileTypes.Image, Path: "picture"}, Expected: false},

		// Tests for Archive file type
		{Name: "Archive - valid extension .zip", Input: InputStruct{FileType: types.FileTypes.Archive, Path: "archive.zip"}, Expected: true},
		{Name: "Archive - valid extension .tar", Input: InputStruct{FileType: types.FileTypes.Archive, Path: "archive.tar"}, Expected: true},
		{Name: "Archive - invalid extension .jpg", Input: InputStruct{FileType: types.FileTypes.Archive, Path: "image.jpg"}, Expected: false},
		{Name: "Archive - empty extension", Input: InputStruct{FileType: types.FileTypes.Archive, Path: "archive"}, Expected: false},

		// Tests for Documents file type
		{Name: "Documents - valid extension .pdf", Input: InputStruct{FileType: types.FileTypes.Documents, Path: "document.pdf"}, Expected: true},
		{Name: "Documents - valid extension .docx", Input: InputStruct{FileType: types.FileTypes.Documents, Path: "document.docx"}, Expected: true},
		{Name: "Documents - invalid extension .mp4", Input: InputStruct{FileType: types.FileTypes.Documents, Path: "video.mp4"}, Expected: false},
		{Name: "Documents - empty extension", Input: InputStruct{FileType: types.FileTypes.Documents, Path: "document"}, Expected: false},

		// Test for fileType not present in FileExtensions map
		{Name: "Unknown type - any extension", Input: InputStruct{FileType: "unknown_type", Path: "file.any"}, Expected: false},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) { // Run each test case as a sub-test
			result := IsExtensionValid(test.Input.FileType, test.Input.Path)
			if result != test.Expected {
				t.Errorf("IsExtensionValid(%q, %q) - %v  = %v; expected %v", test.Input.FileType, test.Input.Path, test.Name, result, test.Expected)
			}
		})
	}
}

// TestIsDirectoryEmpty handles testing for IsDirectoryEmpty func.
func TestIsDirectoryEmpty(t *testing.T) {
	emptyDir := CreateEmptyDir(t)
	nonEmptyDir := CreateNonEmptyDir(t)
	file := CreateTempFile(t)
	readDirErrorPath := CreateEmptyDir(t)
	tests := []*types.TestLayout[string, bool]{
		{Name: "Test empty directory", Input: emptyDir, Expected: true, Err: nil},
		{Name: "Test non-empty directory", Input: nonEmptyDir, Expected: false, Err: nil},
		{Name: "Test non-existing directory", Input: "/path/to/non_existent_dir", Expected: false, Err: fmt.Errorf("The system cannot find the path specified")},
		{Name: "Test file instead of directory", Input: file.Name(), Expected: false, Err: errors.New("not a directory")},
		{Name: "Test Dir with read error", Input: readDirErrorPath, Expected: false, Err: fmt.Errorf("simulated ReadDir error")},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var ops types.DirOps
			switch test.Name {
			case "Test Dir with read error":
				ops = &MockDirOps{
					readDirErr: fmt.Errorf("simulated ReadDir error"),
				}
			default:
				ops = &types.RealDirOps{}
			}

			result, err := IsDirectoryEmpty(test.Input, ops)

			if result != test.Expected {
				t.Errorf("IsDirectoryEmpty(%q) - %v = %v; want %v", test.Input, test.Name, result, test.Expected)
			}

			if (err != nil && test.Err == nil) || (err == nil && test.Err != nil) || (err != nil && test.Err != nil && !strings.Contains(err.Error(), test.Err.Error())) {
				t.Errorf("IsDirectoryEmpty(%q) - %v error = %v; want error containing %v", test.Input, test.Name, err, test.Err)
			}
		})
	}
}

// TestGetOperatorSizeMatches tests GetOperatorSizeMatches func.
func TestGetOperatorSizeMatches(t *testing.T) {
	type InputStruct struct {
		Operator      types.OperatorType
		WantedSize    int64
		ToleranceSize float64
		FileSize      int64
	}

	tests := []*types.TestLayout[InputStruct, bool]{
		{Name: "EqualTo Match FileSize", Input: InputStruct{Operator: types.OperatorTypes.EqualTo, WantedSize: 1024, ToleranceSize: 0, FileSize: 1024}, Expected: true, Err: nil},
		{Name: "EqualTo Match Within Tolerance", Input: InputStruct{Operator: types.OperatorTypes.EqualTo, WantedSize: 1024, ToleranceSize: 1.0, FileSize: 1050}, Expected: true, Err: nil},
		{Name: "EqualTo Outside Tolerance (above)", Input: InputStruct{Operator: types.OperatorTypes.EqualTo, WantedSize: 1024, ToleranceSize: 1.0, FileSize: 2049}, Expected: false, Err: nil},
		{Name: "EqualTo Outside Tolerance (below)", Input: InputStruct{Operator: types.OperatorTypes.EqualTo, WantedSize: 1024, ToleranceSize: 0, FileSize: 0}, Expected: false, Err: nil},
		{Name: "LessThan Less than FileSize", Input: InputStruct{Operator: types.OperatorTypes.LessThan, WantedSize: 1024, ToleranceSize: 0, FileSize: 1023}, Expected: true, Err: nil},
		{Name: "LessThan Equal to FileSize", Input: InputStruct{Operator: types.OperatorTypes.LessThan, WantedSize: 1024, ToleranceSize: 0, FileSize: 1024}, Expected: false, Err: nil},
		{Name: "LessThan Equal to FileSize with Tolerance", Input: InputStruct{Operator: types.OperatorTypes.LessThan, WantedSize: 1024, ToleranceSize: 1.0, FileSize: 1025}, Expected: false, Err: nil},
		{Name: "LessThan with Tolerance Above", Input: InputStruct{Operator: types.OperatorTypes.LessThan, WantedSize: 1024, ToleranceSize: 1.0, FileSize: 1022}, Expected: true, Err: nil},
		{Name: "LessThan with Tolerance Below", Input: InputStruct{Operator: types.OperatorTypes.LessThan, WantedSize: 1024, ToleranceSize: 1.0, FileSize: 1026}, Expected: false, Err: nil},
		{Name: "LessThanEqualTo Equal to FileSize", Input: InputStruct{Operator: types.OperatorTypes.LessThanEqualTo, WantedSize: 1024, ToleranceSize: 0, FileSize: 1024}, Expected: true, Err: nil},
		{Name: "GreaterThan Greater than FileSize", Input: InputStruct{Operator: types.OperatorTypes.GreaterThan, WantedSize: 1024, ToleranceSize: 0, FileSize: 1025}, Expected: true, Err: nil},
		{Name: "GreaterThan Equal to FileSize", Input: InputStruct{Operator: types.OperatorTypes.GreaterThan, WantedSize: 1024, ToleranceSize: 0, FileSize: 1024}, Expected: false, Err: nil},
		{Name: "GreaterThanEqualTo Equal to FileSize", Input: InputStruct{Operator: types.OperatorTypes.GreaterThanEqualTo, WantedSize: 1024, ToleranceSize: 0, FileSize: 1024}, Expected: true, Err: nil},
		{Name: "GreaterThanEqualTo Less than FileSize", Input: InputStruct{Operator: types.OperatorTypes.GreaterThanEqualTo, WantedSize: 1024, ToleranceSize: 0, FileSize: 1023}, Expected: false, Err: nil},
		{Name: "Default Case Invalid Operator (behaves like EqualTo)", Input: InputStruct{Operator: types.OperatorType("invalid"), WantedSize: 1024, ToleranceSize: 1.0, FileSize: 1025}, Expected: true, Err: nil},
		// // Specific use case 315 KB
		{Name: "Equal to FileSize Within Tolerance", Input: InputStruct{Operator: types.OperatorTypes.EqualTo, WantedSize: 315000, ToleranceSize: 0.05, FileSize: 314950}, Expected: true, Err: nil},
		{Name: "Equal to FileSize Outside Tolerance (above)", Input: InputStruct{Operator: types.OperatorTypes.EqualTo, WantedSize: 315000, ToleranceSize: 0.05, FileSize: 330000}, Expected: false, Err: nil},
		{Name: "Equal to FileSize Outside Tolerance (below)", Input: InputStruct{Operator: types.OperatorTypes.EqualTo, WantedSize: 315000, ToleranceSize: 0.05, FileSize: 299000}, Expected: false, Err: nil},
		// Errors
		{Name: "CalculateTolerances Tolerance Cannot Be Negative", Input: InputStruct{Operator: types.OperatorTypes.EqualTo, WantedSize: 0, ToleranceSize: -1, FileSize: 1024}, Expected: false, Err: fmt.Errorf("error calculating tolerances toleranceSize cannot be negative")},
		{Name: "CalculateTolerances WantedSize Cannot Be Negative", Input: InputStruct{Operator: types.OperatorTypes.EqualTo, WantedSize: -1, ToleranceSize: 1, FileSize: 1024}, Expected: false, Err: fmt.Errorf("error calculating tolerances wantedFileSize cannot be negative")},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			result, err := GetOperatorSizeMatches(test.Input.Operator, test.Input.WantedSize, test.Input.ToleranceSize, test.Input.FileSize)

			if result != test.Expected {
				t.Errorf("GetOperatorSizeMatches(%v, %v, %v, %v) - %v = %v; want %v", test.Input.Operator, test.Input.WantedSize, test.Input.ToleranceSize, test.Input.FileSize, test.Name, result, test.Expected)
			}

			if (err != nil && test.Err == nil) || (err == nil && test.Err != nil) || (err != nil && test.Err != nil && err.Error() != test.Err.Error()) {
				t.Errorf("GetOperatorSizeMatches(%v, %v, %v, %v) - %v = %v; want %v", test.Input.Operator, test.Input.WantedSize, test.Input.ToleranceSize, test.Input.FileSize, test.Name, result, test.Expected)
			}
		})
	}
}

// TestCalculateTolerances tests the CalculateTolerances func.
func TestCalculateTolerances(t *testing.T) {
	type InputStruct struct {
		wantedFileSize int64
		toleranceSize  float64
	}
	type ExpectedResults struct {
		results types.ToleranceResults
		err     error
	}

	tests := []*types.TestLayout[InputStruct, ExpectedResults]{
		{Name: "Basic tolerance calculation", Input: InputStruct{wantedFileSize: 1024, toleranceSize: 10}, Expected: ExpectedResults{results: types.ToleranceResults{ToleranceSize: 10240, UpperBoundSize: 11264, LowerBoundSize: 0}, err: nil}},
		{Name: "Zero tolerance size", Input: InputStruct{wantedFileSize: 2048, toleranceSize: 0}, Expected: ExpectedResults{results: types.ToleranceResults{ToleranceSize: 0, UpperBoundSize: 2048, LowerBoundSize: 2048}, err: nil}},
		{Name: "Negative tolerance size", Input: InputStruct{wantedFileSize: 2048, toleranceSize: -5}, Expected: ExpectedResults{results: types.ToleranceResults{ToleranceSize: 0, UpperBoundSize: 0, LowerBoundSize: 0}, err: fmt.Errorf("toleranceSize cannot be negative")}},
		{Name: "Negative wanted file size", Input: InputStruct{wantedFileSize: -1024, toleranceSize: 10}, Expected: ExpectedResults{results: types.ToleranceResults{ToleranceSize: 0, UpperBoundSize: 0, LowerBoundSize: 0}, err: fmt.Errorf("wantedFileSize cannot be negative")}},
		{Name: "Negative tolerance size with large wanted file size", Input: InputStruct{wantedFileSize: 10737418240, toleranceSize: -50}, Expected: ExpectedResults{results: types.ToleranceResults{ToleranceSize: 0, UpperBoundSize: 0, LowerBoundSize: 0}, err: fmt.Errorf("toleranceSize cannot be negative")}},
		{Name: "Lower bound clamped to zero", Input: InputStruct{wantedFileSize: 500, toleranceSize: 600}, Expected: ExpectedResults{results: types.ToleranceResults{ToleranceSize: 614400, UpperBoundSize: 614900, LowerBoundSize: 0}, err: nil}},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			result, err := CalculateTolerances(test.Input.wantedFileSize, test.Input.toleranceSize)
			if (err != nil) != (test.Expected.err != nil) {
				t.Errorf("CalculateTolerances(%v, %v) - %v error = %v; wantErr %v", test.Input.wantedFileSize, test.Input.toleranceSize, test.Name, err, test.Expected.err)
				return
			}

			if result != test.Expected.results {
				t.Errorf("CalculateTolerances(%v, %v) - %v = %v; want %v", test.Input.wantedFileSize, test.Input.toleranceSize, test.Name, result, test.Expected.results)
			}
		})
	}
}

// TestConvertStringSizeToBytes tests ConvertStringSizeToBytes.
func TestConvertStringSizeToBytes(t *testing.T) {
	tests := []*types.TestLayout[string, int64]{
		{Name: "Test Empty Size", Input: " ", Expected: 0, Err: errors.New("size cannot be empty")},
		{Name: "Test 1 B", Input: "1 B", Expected: 1, Err: nil},
		{Name: "Test 10 KB", Input: "10 KB", Expected: 10 * 1024, Err: nil},
		{Name: "Test 1 MB", Input: "1 MB", Expected: 1 * 1024 * 1024, Err: nil},
		{Name: "Test 5 GB", Input: "5 GB", Expected: 5 * 1024 * 1024 * 1024, Err: nil},
		{Name: "Test 100 GB", Input: "100 GB", Expected: 100 * 1024 * 1024 * 1024, Err: nil},
		{Name: "Test 2.5 TB", Input: "2.5 TB", Expected: 2.5 * 1024 * 1024 * 1024 * 1024, Err: nil},
		{Name: "Test 1 kB", Input: "1 kB", Expected: 1 * 1024, Err: nil},
		{Name: "Test 1000", Input: "1000", Expected: 0, Err: errors.New("invalid size format")},
		{Name: "Test 1000 M", Input: "1000 M", Expected: 0, Err: errors.New("invalid size unit")},
		{Name: "Test 1000 XYZ", Input: "1000 XYZ", Expected: 0, Err: errors.New("invalid size unit")},
		{Name: "Test not a size", Input: "not a size", Expected: 0, Err: errors.New("invalid size format")},
		{Name: "Test 12.34.56 MB", Input: "12.34.56 MB", Expected: 0, Err: errors.New(`strconv.ParseFloat: parsing "12.34.56": invalid syntax`)},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			result, err := ConvertStringSizeToBytes(test.Input)

			if result != test.Expected {
				t.Errorf("ConvertStringSizeToBytes(%q) - %v = %v; want %v", test.Input, test.Name, result, test.Expected)
			}

			if (err != nil && test.Err == nil) || (err == nil && test.Err != nil) || (err != nil && test.Err != nil && err.Error() != test.Err.Error()) {
				t.Errorf("ConvertStringSizeToBytes(%q) - %v error = %v; want %v", test.Input, test.Name, err, test.Err)
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

// TestRemoveEmptyDir tests RemoveEmptyDir func.
func TestRemoveEmptyDir(t *testing.T) {
	// Create paths for test cases
	nonExistentDir := "nonexistent-dir"
	file := CreateTempFile(t)
	filePath := file.Name()
	file.Close()
	nonEmptyDirPath := CreateNonEmptyDir(t)
	emptyDirPath := CreateEmptyDir(t)
	readDirErrorPath := CreateEmptyDir(t)
	removeErrorPath := CreateEmptyDir(t)
	statErrorPath := filepath.Join(t.TempDir(), "stat-error-dir")

	// Test cases
	tests := []*types.TestLayout[string, bool]{
		{Name: "Directory does not exist", Input: nonExistentDir, Expected: false, Err: fmt.Errorf("directory does not exist")},
		{Name: "Path is not a directory", Input: filePath, Expected: false, Err: fmt.Errorf("is not a directory")},
		{Name: "Directory is not empty", Input: nonEmptyDirPath, Expected: false, Err: fmt.Errorf("is not empty")},
		{Name: "Successfully remove empty directory", Input: emptyDirPath, Expected: true, Err: nil},
		{Name: "ReadDir error", Input: readDirErrorPath, Expected: false, Err: fmt.Errorf("simulated ReadDir error")},
		{Name: "Remove error", Input: removeErrorPath, Expected: false, Err: fmt.Errorf("simulated Remove error")},
		{Name: "Stat error", Input: statErrorPath, Expected: false, Err: fmt.Errorf("simulated Stat error")},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var ops types.DirOps
			switch test.Name {
			case "ReadDir error":
				ops = &MockDirOps{
					readDirErr: fmt.Errorf("simulated ReadDir error"),
				}
			case "Remove error":
				ops = &MockDirOps{
					removeErr: fmt.Errorf("simulated Remove error"),
				}
			case "Stat error":
				originalStatFunc := osStatFunc
				defer func() { osStatFunc = originalStatFunc }()
				osStatFunc = func(_ string) (os.FileInfo, error) {
					return nil, fmt.Errorf("simulated Stat error")
				}
			default:
				ops = &types.RealDirOps{}
			}

			result, err := RemoveEmptyDir(test.Input, ops)

			if result != test.Expected {
				t.Errorf("RemoveEmptyDir(%q) - %v = %v; want %v", test.Input, test.Name, result, test.Expected)
			}

			if (err != nil && test.Err == nil) || (err == nil && test.Err != nil) || (err != nil && test.Err != nil && !strings.Contains(err.Error(), test.Err.Error())) {
				t.Errorf("RemoveEmptyDir(%q) - %v error = %v; want error containing %v", test.Input, test.Name, err, test.Err)
			}
		})
	}
}
