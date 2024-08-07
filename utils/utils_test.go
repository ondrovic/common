package utils

import (
	"errors"
	"fmt"
	"os/exec"
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

// MockExecCommand replaces execCommand for testing
var MockExecCommand = func(name string, arg ...string) CommandExecutor {
	return &MockCmd{}
}

// TestExeCommand test the ExeCommand var function
func TestExecCommand(t *testing.T) {
	// Override execCommand with a mock implementation
	originalExecCommand := ExecCommand
	ExecCommand = MockExecCommand
	defer func() {
		ExecCommand = originalExecCommand
	}()

	// Define test cases
	tests := []*types.TestLayout[struct {
		CmdName string
		CmdArgs []string
	}, string]{
		{
			Name: "Successful command execution",
			Input: struct {
				CmdName string
				CmdArgs []string
			}{
				CmdName: "echo",
				CmdArgs: []string{"hello"},
			},
			Expected: "hello",
			Err:      nil,
		},
		{
			Name: "Command execution with error",
			Input: struct {
				CmdName string
				CmdArgs []string
			}{
				CmdName: "nonexistent_command",
				CmdArgs: []string{},
			},
			Expected: "command not found",
			Err:      nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// Update mock behavior for each test case
			mockCmd := &MockCmd{err: test.Err}
			ExecCommand = func(name string, arg ...string) CommandExecutor {
				return mockCmd
			}

			cmd := ExecCommand(test.Input.CmdName, test.Input.CmdArgs...)
			err := cmd.Run()

			if !errors.Is(err, test.Err) {
				t.Errorf("ExecCommand(%v, %v) = %v; want %v", test.Input.CmdName, test.Input.CmdArgs, err, test.Expected)
			}
		})
	}
}

// TestClearTerminalScreen tests the ClearTerminalScreen func
func TestClearTerminalScreen(t *testing.T) {
	tests := []*types.TestLayout[string, error]{
		{Name: "Windows", Input: "windows", Expected: nil},
		// {Name: "Linux", Input: "linux", Expected: nil},
		// {Name: "macOS", Input: "darwin", Expected: nil},
		{Name: "Unsupported OS", Input: "unsupported", Expected: fmt.Errorf("unsupported platform: unsupported")},
		{Name: "Command Error", Input: "linux", Expected: fmt.Errorf("mock error")},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// Mock execCommand function based on the test case
			if test.Name == "Command Error" {
				ExecCommand = func(name string, arg ...string) CommandExecutor {
					return &MockCmd{err: fmt.Errorf("mock error")}
				}
			} else {
				ExecCommand = func(name string, arg ...string) CommandExecutor {
					return &RealCmd{cmd: exec.Command(name, arg...)}
				}
			}

			result := ClearTerminalScreen(test.Input)
			if result != nil && test.Expected != nil {
				if result.Error() != test.Expected.Error() {
					t.Errorf("ClearTerminalScreen(%q) = %v; expected %v", test.Input, result, test.Expected)
				}
			} else if result != test.Expected {
				t.Errorf("ClearTerminalScreen(%q) = %v; expected %v", test.Input, result, test.Expected)
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
	tests := []*types.TestLayout[struct{ Path, GOOS string }, string]{
		{Name: "Test Windows path formatting", Input: struct{ Path, GOOS string }{Path: `C:\path\to\file`, GOOS: "windows"}, Expected: `C:\path\to\file`, Err: nil},
		{Name: "Test Linux path formatting", Input: struct{ Path, GOOS string }{Path: `\path\to\file`, GOOS: "linux"}, Expected: `/path/to/file`, Err: nil},
		{Name: "Test macOS path formatting", Input: struct{ Path, GOOS string }{Path: `/path/to/file`, GOOS: "darwin"}, Expected: `/path/to/file`, Err: nil},
		{Name: "Test default case for Unix", Input: struct{ Path, GOOS string }{Path: `/path/to/file`, GOOS: "unknown"}, Expected: `/path/to/file`, Err: nil},
		{Name: "Test Windows path with forward slashes", Input: struct{ Path, GOOS string }{Path: `C:/path/to/file`, GOOS: "windows"}, Expected: `C:\path\to\file`, Err: nil},
		{Name: "Test Unix path with backslashes", Input: struct{ Path, GOOS string }{Path: `\path\to\file`, GOOS: "linux"}, Expected: `/path/to/file`, Err: nil},
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
	tests := []*types.TestLayout[struct {
		fileType types.FileType
		path     string
	}, bool]{
		// Tests for Any file type (wildcard)
		{Name: "Any - valid extension", Input: struct {
			fileType types.FileType
			path     string
		}{fileType: types.FileTypes.Any, path: "example.file"}, Expected: true},
		{Name: "Any - no extension", Input: struct {
			fileType types.FileType
			path     string
		}{fileType: types.FileTypes.Any, path: "example"}, Expected: true},

		// Tests for Video file type
		{Name: "Video - valid extension .mp4", Input: struct {
			fileType types.FileType
			path     string
		}{fileType: types.FileTypes.Video, path: "video.mp4"}, Expected: true},
		{Name: "Video - invalid extension .txt", Input: struct {
			fileType types.FileType
			path     string
		}{fileType: types.FileTypes.Video, path: "document.txt"}, Expected: false},
		{Name: "Video - empty extension", Input: struct {
			fileType types.FileType
			path     string
		}{fileType: types.FileTypes.Video, path: "video"}, Expected: false},

		// Tests for Image file type
		{Name: "Image - valid extension .jpg", Input: struct {
			fileType types.FileType
			path     string
		}{fileType: types.FileTypes.Image, path: "picture.jpg"}, Expected: true},
		{Name: "Image - valid extension .png", Input: struct {
			fileType types.FileType
			path     string
		}{fileType: types.FileTypes.Image, path: "picture.png"}, Expected: true},
		{Name: "Image - invalid extension .mp4", Input: struct {
			fileType types.FileType
			path     string
		}{fileType: types.FileTypes.Image, path: "video.mp4"}, Expected: false},
		{Name: "Image - empty extension", Input: struct {
			fileType types.FileType
			path     string
		}{fileType: types.FileTypes.Image, path: "picture"}, Expected: false},

		// Tests for Archive file type
		{Name: "Archive - valid extension .zip", Input: struct {
			fileType types.FileType
			path     string
		}{fileType: types.FileTypes.Archive, path: "archive.zip"}, Expected: true},
		{Name: "Archive - valid extension .tar", Input: struct {
			fileType types.FileType
			path     string
		}{fileType: types.FileTypes.Archive, path: "archive.tar"}, Expected: true},
		{Name: "Archive - invalid extension .jpg", Input: struct {
			fileType types.FileType
			path     string
		}{fileType: types.FileTypes.Archive, path: "image.jpg"}, Expected: false},
		{Name: "Archive - empty extension", Input: struct {
			fileType types.FileType
			path     string
		}{fileType: types.FileTypes.Archive, path: "archive"}, Expected: false},

		// Tests for Documents file type
		{Name: "Documents - valid extension .pdf", Input: struct {
			fileType types.FileType
			path     string
		}{fileType: types.FileTypes.Documents, path: "document.pdf"}, Expected: true},
		{Name: "Documents - valid extension .docx", Input: struct {
			fileType types.FileType
			path     string
		}{fileType: types.FileTypes.Documents, path: "document.docx"}, Expected: true},
		{Name: "Documents - invalid extension .mp4", Input: struct {
			fileType types.FileType
			path     string
		}{fileType: types.FileTypes.Documents, path: "video.mp4"}, Expected: false},
		{Name: "Documents - empty extension", Input: struct {
			fileType types.FileType
			path     string
		}{fileType: types.FileTypes.Documents, path: "document"}, Expected: false},

		// Test for fileType not present in FileExtensions map
		{Name: "Unknown type - any extension", Input: struct {
			fileType types.FileType
			path     string
		}{fileType: "unknown_type", path: "file.any"}, Expected: false},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) { // Run each test case as a sub-test
			result := IsExtensionValid(test.Input.fileType, test.Input.path)
			if result != test.Expected {
				t.Errorf("IsExtensionValid(%q, %q) = %v; expected %v", test.Input.fileType, test.Input.path, result, test.Expected)
			}
		})
	}
}

// TestGetOperatorSizeMatches tests GetOperatorSizeMatches func
func TestGetOperatorSizeMatches(t *testing.T) {
	tests := []types.TestLayout[struct {
		Operator      types.OperatorType
		FileSize      int64
		ToleranceSize int64
		InfoSize      int64
	}, bool]{
		{
			Name: "EqualTo - Match FileSize",
			Input: struct {
				Operator      types.OperatorType
				FileSize      int64
				ToleranceSize int64
				InfoSize      int64
			}{
				Operator:      types.OperatorTypes.EqualTo,
				FileSize:      1024,
				ToleranceSize: 0,
				InfoSize:      1024,
			},
			Expected: true,
		},
		{
			Name: "EqualTo - Match ToleranceSize",
			Input: struct {
				Operator      types.OperatorType
				FileSize      int64
				ToleranceSize int64
				InfoSize      int64
			}{
				Operator:      types.OperatorTypes.EqualTo,
				FileSize:      1024,
				ToleranceSize: 1050,
				InfoSize:      1050,
			},
			Expected: true,
		},
		{
			Name: "LessThan - Less than FileSize",
			Input: struct {
				Operator      types.OperatorType
				FileSize      int64
				ToleranceSize int64
				InfoSize      int64
			}{
				Operator:      types.OperatorTypes.LessThan,
				FileSize:      1024,
				ToleranceSize: 1050,
				InfoSize:      512,
			},
			Expected: true,
		},
		{
			Name: "LessThan - Less than or equal to FileSize",
			Input: struct {
				Operator      types.OperatorType
				FileSize      int64
				ToleranceSize int64
				InfoSize      int64
			}{
				Operator:      types.OperatorTypes.LessThanEqualTo,
				FileSize:      1024,
				ToleranceSize: 1050,
				InfoSize:      512,
			},
			Expected: true,
		},
		{
			Name: "GreaterThanEqualTo - Greater than FileSize",
			Input: struct {
				Operator      types.OperatorType
				FileSize      int64
				ToleranceSize int64
				InfoSize      int64
			}{
				Operator:      types.OperatorTypes.GreaterThan,
				FileSize:      1024,
				ToleranceSize: 0,
				InfoSize:      2048,
			},
			Expected: true,
		},
		{
			Name: "GreaterThanEqualTo - Greater than or equal to FileSize",
			Input: struct {
				Operator      types.OperatorType
				FileSize      int64
				ToleranceSize int64
				InfoSize      int64
			}{
				Operator:      types.OperatorTypes.GreaterThanEqualTo,
				FileSize:      1024,
				ToleranceSize: 0,
				InfoSize:      2048,
			},
			Expected: true,
		},
		{
			Name: "Default Case - Invalid Operator",
			Input: struct {
				Operator      types.OperatorType
				FileSize      int64
				ToleranceSize int64
				InfoSize      int64
			}{
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
	tests := []*types.TestLayout[struct {
		sizeStr   string
		tolerance float64
	}, int64]{
		{
			Name: "1 KB with 10% tolerance",
			Input: struct {
				sizeStr   string
				tolerance float64
			}{"1 KB", 10},
			Expected: 1126,
			Err:      nil,
		},
		{
			Name: "1 MB with 50% tolerance",
			Input: struct {
				sizeStr   string
				tolerance float64
			}{"1 MB", 50},
			Expected: 1572864,
			Err:      nil,
		},
		{
			Name: "100 B with 100% tolerance",
			Input: struct {
				sizeStr   string
				tolerance float64
			}{"100 B", 100},
			Expected: 200,
			Err:      nil,
		},
		{
			Name: "10 GB with 0% tolerance",
			Input: struct {
				sizeStr   string
				tolerance float64
			}{"10 GB", 0},
			Expected: 10737418240,
			Err:      nil,
		},
		{
			Name: "2.5 KB with 20% tolerance",
			Input: struct {
				sizeStr   string
				tolerance float64
			}{"2.5 KB", 20},
			Expected: 3072,
			Err:      nil,
		},
		{
			Name: "10 MB with 25% tolerance",
			Input: struct {
				sizeStr   string
				tolerance float64
			}{"10 MB", 25},
			Expected: 13107200,
			Err:      nil,
		},
		{
			Name: "1000 B with 0% tolerance",
			Input: struct {
				sizeStr   string
				tolerance float64
			}{"1000 B", 0},
			Expected: 1000,
			Err:      nil,
		},
		{
			Name: "5 GB with -10% tolerance",
			Input: struct {
				sizeStr   string
				tolerance float64
			}{"5 GB", -10},
			Expected: 4831838208,
			Err:      nil,
		},
		{
			Name: "500 KB with 200% tolerance",
			Input: struct {
				sizeStr   string
				tolerance float64
			}{"500 KB", 200},
			Expected: 1536000,
			Err:      nil,
		},
		{
			Name: "Empty size string",
			Input: struct {
				sizeStr   string
				tolerance float64
			}{"", 10},
			Expected: 0,
			Err:      errors.New("size string is empty"),
		},
		{
			Name: "1 KB with -10% tolerance",
			Input: struct {
				sizeStr   string
				tolerance float64
			}{"1 KB", -10},
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
