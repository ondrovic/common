package utils

import (
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"unicode"

	"github.com/ondrovic/common/types"
	"github.com/pterm/pterm"
)

// The CommandExecutor interface defines a method Run that executes a command and returns an error.
// @property {error} Run - The `Run` method is a function defined in the `CommandExecutor` interface.
// This method is expected to execute a command and return an error if any issues occur during the
// execution.
type CommandExecutor interface {
	Run() error
}

// The RealCmd type represents a command to be executed in Go.
// @property cmd - The `cmd` property in the `RealCmd` struct is a pointer to an `exec.Cmd` struct.
// This property is used to store information about a command that can be executed in the operating
// system.
type RealCmd struct {
	cmd *exec.Cmd
}

// The `func (r *RealCmd) Run() error` method is a method defined on the `RealCmd` struct. This method
// implements the `Run` function of the `CommandExecutor` interface.
func (r *RealCmd) Run() error {
	r.cmd.Stdout = os.Stdout
	return r.cmd.Run()
}

var (
	// The line `ExecCommand = func(name string, arg ...string) CommandExecutor {
	// 		return &RealCmd{cmd: exec.Command(name, arg...)}
	// 	}` is defining a variable `ExecCommand` as a function literal. This function takes a `name` string
	// parameter and a variadic parameter `arg` which is a slice of strings. It returns an instance of the
	// `RealCmd` struct that implements the `CommandExecutor` interface.
	ExecCommand = func(name string, arg ...string) CommandExecutor {
		return &RealCmd{cmd: exec.Command(name, arg...)}
	}

	// The line `osStatFunc = os.Stat` is creating a variable `osStatFunc` and assigning it the value of
	// the `os.Stat` function. This is essentially creating an alias for the `os.Stat` function, allowing
	// it to be referenced by the new variable name `osStatFunc` within the package.
	osStatFunc = os.Stat
)

// Function to check if a string field is empty.
func validateStringField(fieldName, value string) error {
	if value == "" {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}
	return nil
}

// Function to check if a struct field is empty.
func validateStructField(fieldName string, value reflect.Value) error {
	if value.IsZero() {
		return fmt.Errorf("%s cannot be an empty struct", fieldName)
	}
	return nil
}

// Validate function using reflection.
func validateApp(app interface{}) error {
	v := reflect.ValueOf(app)

	// If it's a pointer, dereference it
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("validateApp expects a struct")
	}

	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		var err error
		switch value.Kind() {
		case reflect.String:
			err = validateStringField(field.Name, value.String())
		case reflect.Struct:
			err = validateStructField(field.Name, value)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// The function `ApplicationBanner` validates and displays an application banner with specified styles.
func ApplicationBanner(app *types.Application, clearScreen func(string) error) error {
	if err := clearScreen(runtime.GOOS); err != nil {
		return err
	}

	if err := validateApp(app); err != nil {
		return err
	}

	pterm.DefaultHeader.
		WithFullWidth().
		WithBackgroundStyle(
			pterm.NewStyle(app.Style.Color.Background),
		).
		WithTextStyle(
			pterm.NewStyle(app.Style.Color.Foreground),
		).
		Println(app.Name)
	return nil
}

// The function `ClearTerminalScreen` clears the terminal screen based on the operating system provided
// as an argument.
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
			return fmt.Errorf("failed to clear terminal %w", err)
		}
	}

	return nil
}

// The function `GetVersion` is used for setting the version
func GetVersion(version, fallback string) string {
	if version == "" {
		return fallback
	}

	return version
}

// The function `ToFileType` converts a string representation of a file type to a corresponding enum
// value from the `types.FileType` enum.
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

// The function `ToOperatorType` converts a string representation of an operator type to its
// corresponding enum value.
func ToOperatorType(operatorType string) types.OperatorType {
	switch strings.ToLower(operatorType) {
	case "et", "equal to", "equalto", "equal", "==":
		return types.OperatorTypes.EqualTo
	case "gt", "greater", "greater than", "greaterthan", ">":
		return types.OperatorTypes.GreaterThan
	case "gte", "greater than or equal to", "greaterthanorequalto", ">=":
		return types.OperatorTypes.GreaterThanEqualTo
	case "lt", "less", "less than", "lessthan", "<":
		return types.OperatorTypes.LessThan
	case "lte", "less than or equal to", "lessthanorequalto", "<=":
		return types.OperatorTypes.LessThanEqualTo
	default:
		return ""
	}
}

// The `FormatSize` function converts a given size in bytes to a human-readable format with appropriate
// units.
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

// The `FormatPath` function converts file paths to either Windows or Unix style based on the operating
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

// The IsExtensionValid function checks if a given file extension is valid for a specified file type
// based on a predefined list of allowed extensions.
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

// The function `IsDirectoryEmpty` checks if a directory is empty by listing its entries.
func IsDirectoryEmpty(path string, ops types.DirOps) (bool, error) {
	fileInfo, err := osStatFunc(path)
	if err != nil {
		return false, err
	}

	if !fileInfo.IsDir() {
		return false, fmt.Errorf("%s is not a directory", path)
	}

	entries, err := ops.ReadDir(path)
	if err != nil {
		return false, err
	}

	return len(entries) == 0, nil
}

// The function `GetOperatorSizeMatches` determines if a file size matches a specified operator, wanted
// file size, and tolerance size.
func GetOperatorSizeMatches(operator types.OperatorType, wantedFileSize int64, toleranceSize float64, fileSize int64) (bool, error) {
	results, err := CalculateTolerances(wantedFileSize, toleranceSize)
	if err != nil {
		return false, fmt.Errorf("error calculating tolerances %w", err)
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

// The CalculateTolerances function calculates upper and lower bounds based on a wanted file size and
// tolerance size in bytes.
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

// The function `ConvertStringSizeToBytes` converts a string representation of size with units to
// bytes.
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

// The Pluralize function takes a count and returns the singular or plural form of a word based on the
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

// The function `RemoveEmptyDir` checks if a directory is empty and removes it if it is.
func RemoveEmptyDir(path string, ops types.DirOps) (bool, error) {
	// Check if the directory exists
	fileInfo, err := osStatFunc(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, fmt.Errorf("directory does not exist: %w", err)
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
