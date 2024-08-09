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

type MockDirOps struct {
	readDirErr error
	removeErr  error
}

func (m *MockDirOps) ReadDir(name string) ([]os.DirEntry, error) {
	if m.readDirErr != nil {
		return nil, m.readDirErr
	}
	return []os.DirEntry{}, nil
}

func (m *MockDirOps) Remove(name string) error {
	return m.removeErr
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
				t.Errorf("GetOperatorSizeMatches(%v, %v, %v, %v) = %v; want %v", test.Input.Operator, test.Input.WantedSize, test.Input.ToleranceSize, test.Input.FileSize, result, test.Expected)
			}

			if (err != nil && test.Err == nil) || (err == nil && test.Err != nil) || (err != nil && test.Err != nil && err.Error() != test.Err.Error()) {
				t.Errorf("GetOperatorSizeMatches(%v, %v, %v, %v) = %v; want %v", test.Input.Operator, test.Input.WantedSize, test.Input.ToleranceSize, test.Input.FileSize, result, test.Expected)
			}
		})
	}
}

// TestCalculateTolerances tests the CalculateTolerances function
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
				t.Errorf("CalculateTolerances(%v, %v) error = %v; wantErr %v", test.Input.wantedFileSize, test.Input.toleranceSize, err, test.Expected.err)
				return
			}

			if result != test.Expected.results {
				t.Errorf("CalculateTolerances(%v, %v) = %v; want %v", test.Input.wantedFileSize, test.Input.toleranceSize, result, test.Expected.results)
			}
		})
	}
}

// TestConvertStringSizeToBytes tests ConvertStringSizeToBytes
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
		dir := FormatPath(t.TempDir()+"/empty-dir", runtime.GOOS)
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
	file.Close()
	nonEmptyDirPath := createNonEmptyDir(t)
	emptyDirPath := createEmptyDir(t)
	readDirErrorPath := createEmptyDir(t)
	removeErrorPath := createEmptyDir(t)
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
				osStatFunc = func(name string) (os.FileInfo, error) {
					return nil, fmt.Errorf("simulated Stat error")
				}
			default:
				ops = &types.RealDirOps{}
			}

			result, err := RemoveEmptyDir(test.Input, ops)

			if result != test.Expected {
				t.Errorf("RemoveEmptyDir(%q) = %v; want %v", test.Input, result, test.Expected)
			}

			if (err != nil && test.Err == nil) || (err == nil && test.Err != nil) || (err != nil && test.Err != nil && !strings.Contains(err.Error(), test.Err.Error())) {
				t.Errorf("RemoveEmptyDir(%q) error = %v; want error containing %v", test.Input, err, test.Err)
			}
		})
	}
}
