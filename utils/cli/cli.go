package cli

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/ondrovic/common/types"
	"github.com/ondrovic/common/utils"
	"github.com/ondrovic/common/utils/formatters"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	cmd            *exec.Cmd
	ToLowerWrapper = func(input interface{}) (string, error) {
		return formatters.ToLower(input)
	}
)

// The function `HandleCliFlags` is used to handle cobra cli flags.
func HandleCliFlags(cmd *cobra.Command) (bool, error) {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-h", "--help":
			err := cmd.Help()
			return true, err
		case "-v", "--version":
			pterm.Printf("Version: %s", cmd.Version)
			return true, nil
		}
	}
	return false, nil
}

// The function `ClearTerminalScreen` clears the terminal screen based on the operating system specified.
// by the `goos` parameter.
func ClearTerminalScreen(i interface{}) error {
	goosToLower, err := ToLowerWrapper(i)
	if err != nil {
		return fmt.Errorf("error converting goos to lowercase: %v", i)
	}
	switch goosToLower {
	case "linux", "darwin":
		cmd = exec.Command("clear")
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		return fmt.Errorf("unsupported platform: %s", goosToLower)
	}

	// attach to current process
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// Run the command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clear terminal")
	}

	return nil
}

// The function `ApplicationBanner` validates and displays an application banner with specified styles.
func ApplicationBanner(app *types.Application, clearScreen func(interface{}) error) error {
	if err := clearScreen(runtime.GOOS); err != nil {
		return err
	}

	if err := utils.ValidateStruct(app); err != nil {
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
