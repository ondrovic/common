package cli

import (
	"os"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
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
