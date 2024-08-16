package cli

import (
	"os"
	"testing"

	"github.com/ondrovic/common/types"
	"github.com/spf13/cobra"
)

// TestHandleCliFlags tests HandleCliFlags func.
func TestHandleCliFlags(t *testing.T) {
	type ExpectedOutcome struct {
		result bool
		err    error
	}
	tests := []*types.TestLayout[string, ExpectedOutcome]{
		{Name: "Help flag -h", Input: "-h", Expected: ExpectedOutcome{result: true, err: nil}},
		{Name: "Help flag --help", Input: "--help", Expected: ExpectedOutcome{result: true, err: nil}},
		{Name: "Version flag -v", Input: "-v", Expected: ExpectedOutcome{result: true, err: nil}},
		{Name: "Version flag --version", Input: "--version", Expected: ExpectedOutcome{result: true, err: nil}},
		{Name: "Any other flag", Input: "-o", Expected: ExpectedOutcome{result: false, err: nil}},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// Save and restore os.Args
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()

			// Set up test args
			os.Args = []string{"program", test.Input}

			// Create a mock cobra.Command
			cmd := &cobra.Command{
				Use: "test",
			}

			result, err := HandleCliFlags(cmd)

			if result != test.Expected.result {
				t.Errorf("HandleCliFlags() - %v(%q) result = %v; expected %v", test.Name, test.Input, result, test.Expected.result)
			}

			if (err == nil && test.Expected.err != nil) || (err != nil && test.Expected.err == nil) || (err != nil && err.Error() == test.Expected.err.Error()) {
				t.Errorf("HandleCliFlags() - %v(%q) error = %c; expected %v", test.Name, test.Input, err, test.Expected.err)
			}
		})
	}
}
