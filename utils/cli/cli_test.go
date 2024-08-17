package cli

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/ondrovic/common/types"
	"github.com/pterm/pterm"
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

// TestClearTerminalScreen tests the ClearTerminalScreen func.
func TestClearTerminalScreen(t *testing.T) {
	type InputStruct struct {
		goos interface{}
	}
	type ExpectedOutcome struct {
		shouldFail bool
		err        error
	}

	// helper func to determine outcome based on runtimeGOOS
	simulatedOutcome := func(inputOS interface{}) ExpectedOutcome {
		shouldFail := runtime.GOOS != inputOS

		var err error

		if inputOS == "unknown" {
			err = fmt.Errorf("unsupported platform: %s", inputOS)
		} else if inputOS == 123 {
			err = fmt.Errorf("error converting goos to lowercase: %v", inputOS)
		} else if shouldFail {
			err = fmt.Errorf("failed to clear terminal")
		}

		return ExpectedOutcome{
			shouldFail: shouldFail,
			err:        err,
		}
	}

	tests := []*types.TestLayout[InputStruct, ExpectedOutcome]{
		// ExpectedOutcome is calculated before the test runs
		{Name: "Test error converting", Input: InputStruct{goos: 123}},
		{Name: "Test Linux clear command", Input: InputStruct{goos: "linux"}},
		{Name: "Test Windows clear command", Input: InputStruct{goos: "windows"}},
		{Name: "Test unsupported OS", Input: InputStruct{goos: "unknown"}},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// simulated coverage so it works both locally and on gh actions
			results := simulatedOutcome(test.Input.goos)
			test.Expected.shouldFail = results.shouldFail
			test.Expected.err = results.err

			err := ClearTerminalScreen(test.Input.goos)
			switch {
			case err != nil && test.Expected.err == nil:
				t.Errorf("ClearTerminalScreen() - %v(%q) = %v; expected no error", test.Name, test.Input, err)
			case err == nil && test.Expected.err != nil:
				t.Errorf("ClearTerminalScreen() - %v(%q) = no error; expected %v", test.Name, test.Input, test.Expected.err)
			case err != nil && test.Expected.err != nil && !strings.Contains(err.Error(), test.Expected.err.Error()):
				t.Errorf("ClearTerminalScreen() - %v(%q) = %v; expected %v", test.Name, test.Input, err, test.Expected.err)
			}
		})
	}
}

// TestApplicationBanner tests ApplicationBanner func.
func TestApplicationBanner(t *testing.T) {
	type InputStruct struct {
		app         *types.Application
		clearScreen func(interface{}) error
	}

	tests := []*types.TestLayout[InputStruct, error]{
		{
			Name: "ClearTerminalScreen error",
			Input: InputStruct{
				app: &types.Application{
					Name:        "Test App",
					Description: "Test Description",
					Style:       types.Styles{Color: types.Colors{Background: pterm.BgRed, Foreground: pterm.FgWhite}},
					Usage:       "Test Usage",
					Version:     "1.0.0",
				},
				clearScreen: func(i interface{}) error {
					return fmt.Errorf("simulated clear screen error")
				},
			},
			Expected: fmt.Errorf("simulated clear screen error"),
		},
		{Name: "Application is not a struct", Input: InputStruct{app: nil, clearScreen: ClearTerminalScreen}, Expected: fmt.Errorf("validateApp expects a struct")},
		{Name: "Empty application name", Input: InputStruct{app: &types.Application{Description: "Test Description", Style: types.Styles{Color: types.Colors{}}, Usage: "Test Usage", Version: "1.0.0"}, clearScreen: ClearTerminalScreen}, Expected: fmt.Errorf("Name cannot be empty")},
		{Name: "Empty application description", Input: InputStruct{app: &types.Application{Name: "Test App", Style: types.Styles{Color: types.Colors{}}, Usage: "Test Usage", Version: "1.0.0"}, clearScreen: ClearTerminalScreen}, Expected: fmt.Errorf("Description cannot be empty")},
		{Name: "Empty application style struct", Input: InputStruct{app: &types.Application{Name: "Test App", Description: "Test Description", Style: types.Styles{Color: types.Colors{}}, Usage: "Test Usage", Version: "1.0.0"}, clearScreen: ClearTerminalScreen}, Expected: fmt.Errorf("Style cannot be an empty struct")},
		{Name: "Empty application usage", Input: InputStruct{app: &types.Application{Name: "Test App", Description: "Test Description", Style: types.Styles{Color: types.Colors{Background: pterm.BgRed, Foreground: pterm.FgWhite}}, Version: "1.0.0"}, clearScreen: ClearTerminalScreen}, Expected: fmt.Errorf("Usage cannot be empty")},
		{Name: "Empty application version", Input: InputStruct{app: &types.Application{Name: "Test App", Description: "Test Description", Style: types.Styles{Color: types.Colors{Background: pterm.BgRed, Foreground: pterm.FgWhite}}, Usage: "Test Usage"}, clearScreen: ClearTerminalScreen}, Expected: fmt.Errorf("Version cannot be empty")},
		{Name: "Valid application", Input: InputStruct{app: &types.Application{Name: "Test App", Description: "Test Description", Style: types.Styles{Color: types.Colors{Background: pterm.BgRed, Foreground: pterm.FgWhite}}, Usage: "Test Usage", Version: "1.0.0"}, clearScreen: ClearTerminalScreen}, Expected: nil},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			err := ApplicationBanner(test.Input.app, test.Input.clearScreen)
			if test.Expected == nil {
				if err != nil {
					t.Errorf("ApplicationBanner() - %v error = %v, expected nil", test.Name, err)
				}
			} else {
				if err == nil || err.Error() != test.Expected.Error() {
					t.Errorf("ApplicationBanner() - %v error = %v, expected %v", test.Name, err, test.Expected)
				}
			}
		})
	}
}
