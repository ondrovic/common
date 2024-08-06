package utils

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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
	case "equal to", "equalto", "equal", "==":
		return types.OperatorTypes.EqualTo
	case "greater than", "greaterthan", ">":
		return types.OperatorTypes.GreaterThan
	case "greater than or equal to", "greaterthanorequalto", ">=":
		return types.OperatorTypes.GreaterThanEqualTo
	case "less than", "lessthan", "<":
		return types.OperatorTypes.LessThan
	case "less than or equal to", "lessthanorequalto", "<=":
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
