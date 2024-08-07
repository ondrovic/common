package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/ondrovic/common/types"
)

// MockCmd is a mock implementation of CommandExecutor for simulating errors
type MockCmd struct {
	err error
}

// Run simulates command execution by returning the predefined error
func (m *MockCmd) Run() error {
	return m.err
}

// TestClearTerminalScreen tests the ClearTerminalScreen function
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
			if err != nil && test.Expected.err == nil {
				t.Errorf("ClearTerminalScreen(%q) = %v; expected no error", test.Input, err)
			} else if err == nil && test.Expected.err != nil {
				t.Errorf("ClearTerminalScreen(%q) = no error; expected %v", test.Input, test.Expected.err)
			} else if err != nil && test.Expected.err != nil && !strings.Contains(err.Error(), test.Expected.err.Error()) {
				t.Errorf("ClearTerminalScreen(%q) = %v; expected %v", test.Input, err, test.Expected.err)
			}
		})
	}
}

// TestToFileType tests ToFileType func
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
				t.Errorf("ToFileType(%q) = %q; expected %q", test.Input, result, test.Expected)
			}
		})
	}
}

// TestToOperatorType tests ToOperatorType func
func TestToOperatorType(t *testing.T) {
	tests := []*types.TestLayout[string, types.OperatorType]{
		{Name: "Test equal to", Input: "equal to", Expected: types.OperatorTypes.EqualTo, Err: nil},
		{Name: `Test equalto`, Input: "equalto", Expected: types.OperatorTypes.EqualTo, Err: nil},
		{Name: "Test equal", Input: "equal", Expected: types.OperatorTypes.EqualTo, Err: nil},
		{Name: "Test ==", Input: "==", Expected: types.OperatorTypes.EqualTo, Err: nil},

		{Name: "Test greater than", Input: "greater than", Expected: types.OperatorTypes.GreaterThan, Err: nil},
		{Name: "Test greaterthan", Input: "greaterthan", Expected: types.OperatorTypes.GreaterThan, Err: nil},
		{Name: "Test >", Input: ">", Expected: types.OperatorTypes.GreaterThan, Err: nil},

		{Name: "Test greater than or equal to", Input: "greater than or equal to", Expected: types.OperatorTypes.GreaterThanEqualTo, Err: nil},
		{Name: "Test greaterthanorequalto", Input: "greaterthanorequalto", Expected: types.OperatorTypes.GreaterThanEqualTo, Err: nil},
		{Name: "Test >=", Input: ">=", Expected: types.OperatorTypes.GreaterThanEqualTo, Err: nil},

		{Name: "Test less than", Input: "less than", Expected: types.OperatorTypes.LessThan, Err: nil},
		{Name: "Test lessthan", Input: "lessthan", Expected: types.OperatorTypes.LessThan, Err: nil},
		{Name: "Test <", Input: "<", Expected: types.OperatorTypes.LessThan, Err: nil},

		{Name: "Test less than or equal to", Input: "less than or equal to", Expected: types.OperatorTypes.LessThanEqualTo, Err: nil},
		{Name: "Test lessthanorequalto", Input: "lessthanorequalto", Expected: types.OperatorTypes.LessThanEqualTo, Err: nil},
		{Name: "Test <=", Input: "<=", Expected: types.OperatorTypes.LessThanEqualTo, Err: nil},

		{Name: "Test default case", Input: "", Expected: "", Err: nil},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) { // Run each test case as a sub-test
			result := ToOperatorType(test.Input)
			if result != test.Expected {
				t.Errorf("ToFileType(%q) = %q; expected %q", test.Input, result, test.Expected)
			}
		})
	}
}

// TestFormatSize tests FormatSize func
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
				t.Errorf("ToFileType(%q) = %q; expected %q", test.Input, result, test.Expected)
			}
		})
	}
}

// TestFormatPath tests FormatPath func
func TestFormatPath(t *testing.T) {
	type InputStruct struct {
		Path string
		GOOS string
	}
	tests := []*types.TestLayout[InputStruct, string]{
		{Name: "Test Windows path formatting", Input: InputStruct{Path: `C:\path\to\file`, GOOS: "windows"}, Expected: `C:\path\to\file`, Err: nil},
		{Name: "Test Linux path formatting", Input: InputStruct{Path: `\path\to\file`, GOOS: "linux"}, Expected: `/path/to/file`, Err: nil},
		{Name: "Test macOS path formatting", Input: InputStruct{Path: `/path/to/file`, GOOS: "darwin"}, Expected: `/path/to/file`, Err: nil},
		{Name: "Test default case for Unix", Input: InputStruct{Path: `/path/to/file`, GOOS: "unknown"}, Expected: `/path/to/file`, Err: nil},
		{Name: "Test Windows path with forward slashes", Input: InputStruct{Path: `C:/path/to/file`, GOOS: "windows"}, Expected: `C:\path\to\file`, Err: nil},
		{Name: "Test Unix path with backslashes", Input: InputStruct{Path: `\path\to\file`, GOOS: "linux"}, Expected: `/path/to/file`, Err: nil},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) { // Run each test case as a sub-test
			result := FormatPath(test.Input.Path, test.Input.GOOS)
			if result != test.Expected {
				t.Errorf("FormatPath(%q, %q) = %q; expected %q", test.Input.Path, test.Input.GOOS, result, test.Expected)
			}
		})
	}
}

// TestIsExtensionValid tests IsExtensionValid func
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
				t.Errorf("IsExtensionValid(%q, %q) = %v; expected %v", test.Input.FileType, test.Input.Path, result, test.Expected)
			}
		})
	}
}

// TestGetOperatorSizeMatches tests GetOperatorSizeMatches func
func TestGetOperatorSizeMatches(t *testing.T) {
	type InputStruct struct {
		Operator      types.OperatorType
		FileSize      int64
		ToleranceSize int64
		InfoSize      int64
	}
	tests := []types.TestLayout[InputStruct, bool]{
		{
			Name: "EqualTo - Match FileSize",
			Input: InputStruct{
				Operator:      types.OperatorTypes.EqualTo,
				FileSize:      1024,
				ToleranceSize: 0,
				InfoSize:      1024,
			},
			Expected: true,
		},
		{
			Name: "EqualTo - Match ToleranceSize",
			Input: InputStruct{
				Operator:      types.OperatorTypes.EqualTo,
				FileSize:      1024,
				ToleranceSize: 1050,
				InfoSize:      1050,
			},
			Expected: true,
		},
		{
			Name: "LessThan - Less than FileSize",
			Input: InputStruct{
				Operator:      types.OperatorTypes.LessThan,
				FileSize:      1024,
				ToleranceSize: 1050,
				InfoSize:      512,
			},
			Expected: true,
		},
		{
			Name: "LessThan - Less than or equal to FileSize",
			Input: InputStruct{
				Operator:      types.OperatorTypes.LessThanEqualTo,
				FileSize:      1024,
				ToleranceSize: 1050,
				InfoSize:      512,
			},
			Expected: true,
		},
		{
			Name: "GreaterThan - Greater than FileSize",
			Input: InputStruct{
				Operator:      types.OperatorTypes.GreaterThan,
				FileSize:      1024,
				ToleranceSize: 0,
				InfoSize:      2048,
			},
			Expected: true,
		},
		{
			Name: "GreaterThanEqualTo - Greater than or equal to FileSize",
			Input: InputStruct{
				Operator:      types.OperatorTypes.GreaterThanEqualTo,
				FileSize:      1024,
				ToleranceSize: 0,
				InfoSize:      2048,
			},
			Expected: true,
		},
		{
			Name: "Default Case - Invalid Operator",
			Input: InputStruct{
				Operator:      types.OperatorType("invalid"), // an invalid operator
				FileSize:      1024,
				ToleranceSize: 1050,
				InfoSize:      1024,
			},
			Expected: true, // The default case returns true if infoSize equals either fileSize or toleranceSize
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			result := GetOperatorSizeMatches(test.Input.Operator, test.Input.FileSize, test.Input.ToleranceSize, test.Input.InfoSize)

			if result != test.Expected {
				t.Errorf("GetOperatorSizeMatches(%v, %v, %v, %v) = %v; want %v",
					test.Input.Operator,
					test.Input.FileSize,
					test.Input.ToleranceSize,
					test.Input.InfoSize,
					result,
					test.Expected,
				)
			}
		})
	}
}

// TestCalculateToleranceToBytes test CalculateToleranceToBytes func
func TestCalculateToleranceToBytes(t *testing.T) {
	type InputStruct struct {
		sizeStr   string
		tolerance float64
	}

	tests := []*types.TestLayout[InputStruct, int64]{
		{
			Name: "1 KB with 10% tolerance",
			Input: InputStruct{
				sizeStr:   "1 KB",
				tolerance: 10,
			},
			Expected: 1126,
			Err:      nil,
		},
		{
			Name: "1 MB with 50% tolerance",
			Input: InputStruct{
				sizeStr:   "1 MB",
				tolerance: 50,
			},
			Expected: 1572864,
			Err:      nil,
		},
		{
			Name: "100 B with 100% tolerance",
			Input: InputStruct{
				sizeStr:   "100 B",
				tolerance: 100,
			},
			Expected: 200,
			Err:      nil,
		},
		{
			Name: "10 GB with 0% tolerance",
			Input: InputStruct{
				sizeStr:   "10 GB",
				tolerance: 0,
			},
			Expected: 10737418240,
			Err:      nil,
		},
		{
			Name: "2.5 KB with 20% tolerance",
			Input: InputStruct{
				sizeStr:   "2.5 KB",
				tolerance: 20,
			},
			Expected: 3072,
			Err:      nil,
		},
		{
			Name: "10 MB with 25% tolerance",
			Input: InputStruct{
				sizeStr:   "10 MB",
				tolerance: 25,
			},
			Expected: 13107200,
			Err:      nil,
		},
		{
			Name: "1000 B with 0% tolerance",
			Input: InputStruct{
				sizeStr:   "1000 B",
				tolerance: 0,
			},
			Expected: 1000,
			Err:      nil,
		},
		{
			Name: "5 GB with -10% tolerance",
			Input: InputStruct{
				sizeStr:   "5 GB",
				tolerance: -10,
			},
			Expected: 4831838208,
			Err:      nil,
		},
		{
			Name: "500 KB with 200% tolerance",
			Input: InputStruct{
				sizeStr:   "500 KB",
				tolerance: 200,
			},
			Expected: 1536000,
			Err:      nil,
		},
		{
			Name: "Empty size string",
			Input: InputStruct{
				sizeStr:   "",
				tolerance: 10,
			},
			Expected: 0,
			Err:      errors.New("size string is empty"),
		},
		{
			Name: "1 KB with -10% tolerance",
			Input: InputStruct{
				sizeStr:   "1 KB",
				tolerance: -10,
			},
			Expected: 921,
			Err:      nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			result, err := CalculateToleranceToBytes(test.Input.sizeStr, test.Input.tolerance)
			if (err != nil) != (test.Err != nil) {
				t.Errorf("CalculateToleranceToBytes(%v, %v) error = %v; wantErr %v", test.Input.sizeStr, test.Input.tolerance, err, test.Err)
				return
			}
			if result != test.Expected {
				t.Errorf("CalculateToleranceToBytes(%v, %v) = %v; want %v", test.Input.sizeStr, test.Input.tolerance, result, test.Expected)
			}
		})
	}
}


// TestConvertStringSizeToBytes tests ConvertStringSizeToBytes
func TestConvertStringSizeToBytes(t *testing.T) {
	tests := []*types.TestLayout[string, int64]{
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
		t.Run(test.Input, func(t *testing.T) {
			result, err := ConvertStringSizeToBytes(test.Input)

			if result != test.Expected {
				t.Errorf("ConvertStringSizeToBytes(%q) = %v; want %v", test.Input, result, test.Expected)
			}

			if (err != nil && test.Err == nil) || (err == nil && test.Err != nil) || (err != nil && test.Err != nil && err.Error() != test.Err.Error()) {
				t.Errorf("ConvertStringSizeToBytes(%q) error = %v; want %v", test.Input, err, test.Err)
			}
		})
	}
}

// TestRemoveEmptyDir tests RemoveEmptyDir function
func TestRemoveEmptyDir(t *testing.T) {
	// Helper function to create a temporary file for testing
	createTempFile := func(t *testing.T) *os.File {
		file, err := os.CreateTemp("", "testfile-*.txt")
		if err != nil {
			t.Fatalf("failed to create file: %v", err)
		}
		return file
	}

	// Helper function to create an empty directory for testing
	createEmptyDir := func(t *testing.T) string {
		dir := t.TempDir() + "/empty-dir"
		if err := os.Mkdir(dir, 0755); err != nil {
			t.Fatalf("failed to create directory: %v", err)
		}
		return dir
	}

	// Helper function to create a non-empty directory for testing
	createNonEmptyDir := func(t *testing.T) string {
		dir := filepath.Join(t.TempDir(), FormatPath("/non-empty-dir", runtime.GOOS))
		if err := os.Mkdir(dir, 0755); err != nil {
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

	// Create paths for test cases
	nonExistentDir := "nonexistent-dir"
	file := createTempFile(t)
	filePath := file.Name()
	file.Close() // Ensure the file is closed before using it in tests
	nonEmptyDirPath := createNonEmptyDir(t)
	emptyDirPath := createEmptyDir(t)

	// Test cases
	tests := []*types.TestLayout[string, bool]{
		{
			Name:     "Directory does not exist",
			Input:    nonExistentDir,
			Expected: false,
			Err:      fmt.Errorf("directory does not exist: CreateFile %v: The system cannot find the file specified.", nonExistentDir),
		},
		{
			Name:     "Path is not a directory",
			Input:    filePath,
			Expected: false,
			Err:      fmt.Errorf("%s is not a directory", filePath),
		},
		{
			Name:     "Directory is not empty",
			Input:    nonEmptyDirPath,
			Expected: false,
			Err:      fmt.Errorf("directory %s is not empty", nonEmptyDirPath),
		},
		{
			Name:     "Successfully remove empty directory",
			Input:    emptyDirPath,
			Expected: true,
			Err:      nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Input, func(t *testing.T) {
			result, err := RemoveEmptyDir(test.Input)

			if result != test.Expected {
				t.Errorf("RemoveEmptyDir(%q) = %v; want %v", test.Input, result, test.Expected)
			}

			if (err != nil && test.Err == nil) || (err == nil && test.Err != nil) || (err != nil && test.Err != nil && err.Error() != test.Err.Error()) {
				t.Errorf("RemoveEmptyDir(%q) error = %v; want %v", test.Input, err, test.Err)
			}
		})
	}
}
