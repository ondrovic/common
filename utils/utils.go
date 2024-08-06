package utils

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// FileType
type FileType string

// OperatorType
type OperatorType string

// SizeUnit struct for units
type SizeUnit struct {
	Label string
	Size  int64
}

// TestLayout
type TestLayout[InputT any, ExpectedT any] struct {
	Name     string
	Input    InputT
	Expected ExpectedT
	Err      error
}

// CommandExecutor is an interface for executing commands
type CommandExecutor interface {
	Run() error
}

// RealCmd is a wrapper for exec.Cmd that implements CommandExecutor
type RealCmd struct {
	cmd *exec.Cmd
}

// Run executes the command and returns any errors
func (r *RealCmd) Run() error {
	r.cmd.Stdout = os.Stdout
	return r.cmd.Run()
}

// FileTypes
const (
	Any       FileType = "Any"
	Video     FileType = "Video"
	Image     FileType = "Image"
	Archive   FileType = "Archive"
	Documents FileType = "Documents"
)

// OperatorTypes
const (
	EqualTo            OperatorType = "Equal To"
	GreaterThan        OperatorType = "Greater Than"
	GreaterThanEqualTo OperatorType = "Greater Than Or Equal To"
	LessThan           OperatorType = "Less Than"
	LessThanEqualTo    OperatorType = "Less Than Or Equal To"
)

var (
	// SizeUnits
	SizeUnits = []SizeUnit{
		{Label: "PB", Size: 1 << 50}, // Petabyte
		{Label: "TB", Size: 1 << 40}, // Terabyte
		{Label: "GB", Size: 1 << 30}, // Gigabyte
		{Label: "MB", Size: 1 << 20}, // Megabyte
		{Label: "KB", Size: 1 << 10}, // Kilobyte
		{Label: "B", Size: 1},        // Byte
	}
	// FileExtensions
	FileExtensions = map[FileType]map[string]bool{
		Any: {
			"*.*": true,
		},
		Video: {
			".mp4": true, ".avi": true, ".mkv": true, ".mov": true, ".wmv": true,
			".flv": true, ".webm": true, ".m4v": true, ".mpg": true, ".mpeg": true,
			".ts": true,
		},
		Image: {
			".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".bmp": true,
			".tiff": true, ".webp": true, ".svg": true, ".raw": true, ".heic": true,
			".ico": true,
		},
		Archive: {
			".zip": true, ".rar": true, ".7z": true, ".tar": true, ".gz": true,
			".bz2": true, ".xz": true, ".iso": true, ".tgz": true, ".tbz2": true,
		},
		Documents: {
			".docx": true, ".doc": true, ".pdf": true, ".txt": true, ".rtf": true,
			".odt": true, ".xlsx": true, ".xls": true, ".pptx": true, ".ppt": true,
			".csv": true, ".md": true, ".pages": true,
		},
	}

	execCommand = func(name string, arg ...string) CommandExecutor {
		return &RealCmd{cmd: exec.Command(name, arg...)}
	}
)

// ClearTerminalScreen clears the terminal based on the provided OS name
func ClearTerminalScreen(goos string) error {
	var cmd CommandExecutor
	var err error

	switch strings.ToLower(goos) {
	case "linux", "darwin":
		cmd = execCommand("clear")
	case "windows":
		cmd = execCommand("cmd", "/c", "cls")
	default:
		return fmt.Errorf("unsupported platform: %s", goos)
	}

	if cmd != nil {
		err = cmd.Run()
		if err != nil {
			fmt.Printf("failed to clear terminal: %s\n", err)
			return err
		}
	}

	return nil
}

// ToFileType
func ToFileType(fileType string) FileType {
	switch strings.ToLower(fileType) {
	case "any":
		return Any
	case "video":
		return Video
	case "image":
		return Image
	case "archive":
		return Archive
	case "documents":
		return Documents
	default:
		return ""
	}
}

// ToOperatorType
func ToOperatorType(operatorType string) OperatorType {
	switch strings.ToLower(operatorType) {
	case "equal to", "equalto", "equal", "==":
		return EqualTo
	case "greater than", "greaterthan", ">":
		return GreaterThan
	case "greater than or equal to", "greaterthanorequalto", ">=":
		return GreaterThanEqualTo
	case "less than", "lessthan", "<":
		return LessThan
	case "less than or equal to", "lessthanorequalto", "<=":
		return LessThanEqualTo
	default:
		return ""
	}
}

// FormatSize formats size to human readable
func FormatSize(bytes int64) string {
	for _, unit := range SizeUnits {
		if bytes >= unit.Size {
			value := float64(bytes) / float64(unit.Size)
			// Round the value to two decimal places
			roundedValue := math.Round(value*100) / 100
			return fmt.Sprintf("%.2f %s", roundedValue, unit.Label)
		}
	}

	return "0 B"
}

// IsExtensionValid checks if the file's extension is allowed for a given file type.
func IsExtensionValid(fileType FileType, path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	extensions, exists := FileExtensions[fileType]
	if !exists {
		return false
	}

	// Check for wildcard entry (Any)
	if _, found := extensions["*.*"]; found {
		return true
	}

	// Check if the file extension is explicitly allowed
	return extensions[ext]
}
