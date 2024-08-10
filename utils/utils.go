package utils

import (
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	"github.com/ondrovic/common/types"
)

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

var (
	ExecCommand = func(name string, arg ...string) CommandExecutor {
		return &RealCmd{cmd: exec.Command(name, arg...)}
	}

	osStatFunc = os.Stat
)

// ClearTerminalScreen clears the terminal based on the provided OS name
func ClearTerminalScreen(goos string) error {
	var cmd CommandExecutor
	var err error

	switch strings.ToLower(goos) {
	case "linux", "darwin":
		cmd = ExecCommand("clear")
	case "windows":
		cmd = ExecCommand("cmd", "/c", "cls")
	default:
		return fmt.Errorf("unsupported platform: %s", goos)
	}

	if cmd != nil {
		err = cmd.Run()
		if err != nil {
			// fmt.Printf("failed to clear terminal: %s\n", err)
			return fmt.Errorf("failed to clear terminal %s", err)
		}
	}

	return nil
}

// ToFileType
func ToFileType(fileType string) types.FileType {
	switch strings.ToLower(fileType) {
	case "any":
		return types.FileTypes.Any
	case "video":
		return types.FileTypes.Video
	case "image":
		return types.FileTypes.Image
	case "archive":
		return types.FileTypes.Archive
	case "documents":
		return types.FileTypes.Documents
	default:
		return ""
	}
}

// ToOperatorType
func ToOperatorType(operatorType string) types.OperatorType {
	switch strings.ToLower(operatorType) {
	case "et", "equal to", "equalto", "equal", "==":
		return types.OperatorTypes.EqualTo
	case "gt", "greater", "greater than", "greaterthan", ">":
		return types.OperatorTypes.GreaterThan
	case "gte","greater than or equal to", "greaterthanorequalto", ">=":
		return types.OperatorTypes.GreaterThanEqualTo
	case "lt","less", "less than", "lessthan", "<":
		return types.OperatorTypes.LessThan
	case "lte","less than or equal to", "lessthanorequalto", "<=":
		return types.OperatorTypes.LessThanEqualTo
	default:
		return ""
	}
}

// FormatSize formats size to human readable
func FormatSize(bytes int64) string {
	for _, unit := range types.SizeUnits {
		if bytes >= unit.Size {
			value := float64(bytes) / float64(unit.Size)
			// Round the value to two decimal places
			roundedValue := math.Round(value*100) / 100
			return fmt.Sprintf("%.2f %s", roundedValue, unit.Label)
		}
	}

	return "0 B"
}

// FormatPath formats the path based on the operating system
func FormatPath(path string, goos string) string {
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

// IsExtensionValid checks if the file's extension is allowed for a given file type.
func IsExtensionValid(fileType types.FileType, path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	extensions, exists := types.FileExtensions[fileType]
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

// GetOperatorSizeMatches determines whether a file matches the size or falls within the tolerance range.
func GetOperatorSizeMatches(operator types.OperatorType, wantedFileSize int64, toleranceSize float64, fileSize int64) (bool, error) {
	results, err := CalculateTolerances(wantedFileSize, toleranceSize)
	if err != nil {
		return false, fmt.Errorf("error calculating tolerances %v", err)
	}

	switch operator {
	case types.OperatorTypes.EqualTo:
		return fileSize >= results.LowerBoundSize && fileSize <= results.UpperBoundSize, nil
	case types.OperatorTypes.LessThan:
		return fileSize < wantedFileSize, nil // Changed lowerBound to fileSize
	case types.OperatorTypes.LessThanEqualTo:
		return fileSize <= wantedFileSize, nil // Changed upperBound to fileSize
	case types.OperatorTypes.GreaterThan:
		return fileSize > wantedFileSize, nil // Changed upperBound to fileSize
	case types.OperatorTypes.GreaterThanEqualTo:
		return fileSize >= wantedFileSize, nil // Changed lowerBound to fileSize
	default:
		return fileSize >= results.LowerBoundSize && fileSize <= results.UpperBoundSize, nil
	}
}

// CalculateTolerances handles calculating the tolerance threshold based on the wantedFileSize and toleranceSize
func CalculateTolerances(wantedFileSize int64, toleranceSize float64) (types.ToleranceResults, error) {
    // Check for invalid input values
    if wantedFileSize < 0 {
        return types.ToleranceResults{}, fmt.Errorf("wantedFileSize cannot be negative")
    }
    if toleranceSize < 0 {
        return types.ToleranceResults{}, fmt.Errorf("toleranceSize cannot be negative")
    }

    // Calculate tolerance in bytes (using int64 directly)
    toleranceBytes := int64(toleranceSize * 1024)

    // Calculate upper and lower bounds
    upperBoundSize := wantedFileSize + toleranceBytes
    lowerBoundSize := wantedFileSize - toleranceBytes

    // Ensure lower bound does not go below zero
    if lowerBoundSize < 0 {
        lowerBoundSize = 0
    }

    // Return the calculated tolerance results
    return types.ToleranceResults{
        ToleranceSize:  toleranceBytes,
        UpperBoundSize: upperBoundSize,
        LowerBoundSize: lowerBoundSize,
    }, nil
}

// ConvertStringSizeToBytes converts a size string with a unit to bytes.
func ConvertStringSizeToBytes(sizeStr string) (int64, error) {
	sizeStr = strings.TrimSpace(sizeStr)
	if sizeStr == "" {
		return 0, errors.New("size cannot be empty")
	}

	// Separate the numeric part and the unit part
	var numStr, unitStr string
	for i, r := range sizeStr {
		if unicode.IsLetter(r) {
			numStr = strings.TrimSpace(sizeStr[:i])
			unitStr = strings.TrimSpace(sizeStr[i:])
			break
		}
	}

	// If no unit was found, return an error
	if unitStr == "" || numStr == "" {
		return 0, errors.New("invalid size format")
	}

	// Normalize the unit string to uppercase
	unitStr = strings.ToUpper(unitStr)

	// Parse the numeric part
	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, err
	}

	// Find the matching unit and convert to bytes
	// for _, unit := range types.Units {
	for _, unit := range types.SizeUnits {
		if unit.Label == unitStr {
			return int64(num * float64(unit.Size)), nil
		}
	}

	return 0, errors.New("invalid size unit")
}

// RemoveEmptyDir handles removing of empty directories
func RemoveEmptyDir(path string, ops types.DirOps) (bool, error) {
    // Check if the directory exists
    fileInfo, err := osStatFunc(path)
    if err != nil {
        if os.IsNotExist(err) {
            return false, fmt.Errorf("directory does not exist: %v", err)
        }
        return false, err
    }

    // Check if it's a directory
    if !fileInfo.IsDir() {
        return false, fmt.Errorf("%s is not a directory", path)
    }

    // Read the directory contents
    entries, err := ops.ReadDir(path)
    if err != nil {
        return false, err
    }

    // Check if the directory is empty
    if len(entries) > 0 {
        return false, fmt.Errorf("directory %s is not empty", path)
    }

    // Remove the directory
    err = ops.Remove(path)
    if err != nil {
        return false, err
    }

    return true, nil
}