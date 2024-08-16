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
	"github.com/ondrovic/common/utils/formatters"
	"github.com/pterm/pterm"
)

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

// // Mock formatter to simulate the ToLower function behavior
// var formatter = struct {
// 	ToLower func(string) (string, error)
// }{
// 	ToLower: func(s string) (string, error) {
// 		if s == "error" {
// 			return "", fmt.Errorf("conversion error")
// 		}
// 		return strings.ToLower(s), nil
// 	},
// }

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
	dir := formatters.FormatPath(filepath.Join(t.TempDir(), "empty-dir"), runtime.GOOS)
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}
	return dir
}

// CreateNonEmptyDir creates a non-empty directory for testing.
func CreateNonEmptyDir(t *testing.T) string {
	t.Helper()
	dir := filepath.Join(t.TempDir(), formatters.FormatPath("non-empty-dir", runtime.GOOS))
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
	type InputStruct struct {
		goos string
	}
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
				fmt.Errorf("failed to clear terminal"),
			}
		}

		return ExpectedOutcome{
			false,
			nil,
		}
	}

	tests := []*types.TestLayout[InputStruct, ExpectedOutcome]{
		// ExpectedOutcome is calculated before the test runs
		{Name: "Test Linux clear command", Input: InputStruct{goos: "linux"}},
		{Name: "Test Windows clear command", Input: InputStruct{goos: "windows"}},
		{Name: "Test unsupported OS", Input: InputStruct{goos: "unknown"}},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// Determine if the test should fail based on the current OS
			results := expectedResultsBasedOnOS(test.Input.goos)

			test.Expected.shouldFail = results.shouldFail
			test.Expected.err = results.err

			err := ClearTerminalScreen(test.Input.goos)
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
	statFilePath := filepath.Join(emptyDir, "whoops-not-here.txt")
	tests := []*types.TestLayout[string, bool]{
		{Name: "Test empty directory", Input: emptyDir, Expected: true, Err: nil},
		{Name: "Test non-empty directory", Input: nonEmptyDir, Expected: false, Err: nil},
		{Name: "Test non-existing directory stat", Input: statFilePath, Expected: false, Err: fmt.Errorf("stat error - file not found: %s", statFilePath)},
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

// TestInRange tests the InRange function with different scenarios.
func TestInRange(t *testing.T) {
	type InputStruct struct {
		target  interface{}
		options interface{}
	}

	// Save the original ToLowerWrapper function
	originalToLower := ToLowerWrapper

	// Defer function to restore the original ToLowerWrapper after tests are done
	defer func() {
		ToLowerWrapper = originalToLower
	}()

	// Mock ToLowerWrapper function
	ToLowerWrapper = func(input interface{}) (string, error) {
		switch input {
		case "errorOption", "error":
			return "", fmt.Errorf("conversion error")
		default:
			return formatters.ToLower(input)
		}

	}

	tests := []*types.TestLayout[InputStruct, bool]{
		{Name: "Test target matches option - single match", Input: InputStruct{target: "Hello", options: []string{"hello", "world"}}, Expected: true},
		{Name: "Test target does not match any option", Input: InputStruct{target: "foo", options: []string{"bar", "baz"}}, Expected: false},
		{Name: "Test single option match", Input: InputStruct{target: "yes", options: "yes"}, Expected: true},
		{Name: "Test empty target", Input: InputStruct{target: "", options: []string{"a", "b"}}, Expected: false},
		{Name: "Test empty options", Input: InputStruct{target: "test", options: []string{}}, Expected: false},
		{Name: "Test invalid target type", Input: InputStruct{target: 123, options: []string{"a", "b"}}, Expected: false, Err: fmt.Errorf("error converting target to lowercase: 123")},
		{Name: "Test invalid options type", Input: InputStruct{target: "test", options: 123}, Expected: false, Err: fmt.Errorf("options must be a string or a slice of strings")},
		{Name: "Test conversion error in target", Input: InputStruct{target: "error", options: []string{"valid"}}, Expected: false, Err: fmt.Errorf("error converting target to lowercase: error")},
		{Name: "Test conversion error in options", Input: InputStruct{target: "valid", options: []string{"error"}}, Expected: false, Err: fmt.Errorf("error converting option to lowercase: error")},
		{Name: "Test conversion error in options string", Input: InputStruct{target: "valid", options: "errorOption"}, Expected: false, Err: fmt.Errorf("error converting option to lowercase: errorOption")},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			result, err := InRange(test.Input.target, test.Input.options)

			if result != test.Expected {
				t.Errorf("InRange(%v, %v) - %v = %v; expected %v", test.Input.target, test.Input.options, test.Name, result, test.Expected)
			}
			if (err != nil && test.Err == nil) || (err == nil && test.Err != nil) || (err != nil && test.Err != nil && err.Error() != test.Err.Error()) {
				t.Errorf("InRange(%v, %v) - %v = %v; expected %v", test.Input.target, test.Input.options, test.Name, err, test.Err)
			}
		})
	}
}

func TestDefaultToLowerWrapper(t *testing.T) {
	// Test the default behavior of ToLowerWrapper
	input := "TeSt"
	expected := "test"

	result, err := ToLowerWrapper(input)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result != expected {
		t.Errorf("ToLowerWrapper(%v) = %v; expected %v", input, result, expected)
	}
}
