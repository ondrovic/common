package results

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"testing"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/ondrovic/common/types"
	"github.com/pterm/pterm"
)

// TODO: move new tests over to []*types.TestLayout[] format
// TODO: consolidate tests if possible

// Define items for test cases.

type Event struct {
	Name      string
	Timestamp time.Time
}

type Person struct {
	Name   string
	Age    int
	Income float64
	Height uint64
}

type MockFileInfo struct {
	Name string
	Size int64
}

type EmbeddedStruct struct {
	Name string
	Age  int
}

type PersonWithEmbedded struct {
	EmbeddedStruct
	Income float64
}

type TestStruct struct {
	Name string
	Age  int
	Size int64
}

// MockWriter is a custom io.Writer for capturing output.
type MockWriter struct {
	buffer bytes.Buffer
}

func (m *MockWriter) Write(p []byte) (n int, err error) {
	return m.buffer.Write(p)
}

func (m *MockWriter) String() string {
	return m.buffer.String()
}

// TestGenericRenderResultsTableInterfaceNonSlice tests the case where a non-slice input is provided
func TestGenericRenderResultsTableInterfaceNonSlice(t *testing.T) {
	// Redirect pterm output to a buffer
	oldOutput := pterm.Output
	var buf bytes.Buffer

	// Assign the buffer to pterm.Output to capture the output
	pterm.Output = true
	defer func() { pterm.Output = oldOutput }()

	// Call the function with a non-slice input
	GenericRenderResultsTableInterface("Not a slice", nil)

	// Read the output from the buffer
	output := buf.String()

	// Define the expected message since pterm.Println() is just doing a st.out it will be blank
	expectedMessage := ""

	// Check if the expected error message is in the output
	if !strings.Contains(output, expectedMessage) {
		t.Errorf("Expected output to contain '%s', but got: %s", expectedMessage, output)
	}

	// Check that the function returned (indirectly, by ensuring no other output was produced)
	if strings.TrimSpace(output) != expectedMessage {
		t.Errorf("Expected only '%s' in output, but got additional content: %s", expectedMessage, output)
	}
}

// TestGenericRenderResultsTableInterface tests the GenericRenderResultsTableInterface function.
func TestGenericRenderResultsTableInterface(t *testing.T) {
	// Redirect stdout to capture output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Test cases
	testCases := []struct {
		name         string
		slice        interface{}
		totalValues  map[string]interface{}
		expectedKeys []string
	}{
		{name: "Test with Person struct", slice: []Person{{Name: "Alice", Age: 30, Income: 50000, Height: 170}, {Name: "Bob", Age: 35, Income: 60000, Height: 180}}, totalValues: map[string]interface{}{"Age": 65, "Income": 110000}, expectedKeys: []string{"NAME", "AGE", "INCOME", "HEIGHT"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			GenericRenderResultsTableInterface(tc.slice, tc.totalValues)

			// Read captured output
			w.Close()
			out, _ := io.ReadAll(r)
			output := string(out)

			// Check if all expected keys are in the output
			for _, key := range tc.expectedKeys {
				if !contains(output, key) {
					t.Errorf("Expected key %s not found in output", key)
				}
			}

			// If totalValues are provided, check if they are in the output
			if tc.totalValues != nil {
				for key, value := range tc.totalValues {
					if !containsFormattedValue(output, value) {
						t.Errorf("Expected total value %v for key %s not found in output", value, key)
					}
				}
			}
		})
	}

	// Restore stdout
	os.Stdout = old

}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// Helper function to check if a string contains a formatted value
func containsFormattedValue(s string, value interface{}) bool {
	strValue := fmt.Sprintf("%v", value)

	// Remove any formatting (commas, spaces) from the output
	cleanOutput := regexp.MustCompile(`[,\s]`).ReplaceAllString(s, "")

	// For floating-point numbers, we need to be more flexible
	if floatValue, err := strconv.ParseFloat(strValue, 64); err == nil {
		// Look for the float value with some tolerance
		tolerance := 0.001
		pattern := fmt.Sprintf(`%.*f`, 3, floatValue) // Use 3 decimal places
		re := regexp.MustCompile(pattern)
		matches := re.FindAllString(cleanOutput, -1)
		for _, match := range matches {
			matchFloat, _ := strconv.ParseFloat(match, 64)
			if math.Abs(matchFloat-floatValue) < tolerance {
				return true
			}
		}
	}

	// For integers and other types, we can do a simple contains check
	return strings.Contains(cleanOutput, strValue)
}

// TestGetHeadersAndFields tests the getHeadersAndFields function
func TestGetHeadersAndFields(t *testing.T) {
	testCases := []struct {
		name            string
		input           interface{}
		expectedHeaders []string
		expectedFields  []string
	}{
		{name: "Test with Event struct", input: Event{}, expectedHeaders: []string{"Name", "Timestamp"}, expectedFields: []string{"Name", "Timestamp"}},
		{name: "Test with Person struct", input: Person{}, expectedHeaders: []string{"Name", "Age", "Income", "Height"}, expectedFields: []string{"Name", "Age", "Income", "Height"}},
		{name: "Test with non-struct input", input: "Not a struct", expectedHeaders: nil, expectedFields: nil},
		{name: "Test with embedded struct", input: PersonWithEmbedded{}, expectedHeaders: []string{"Name", "Age", "Income"}, expectedFields: []string{"Name", "Age", "Income"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			headers, fields := getHeadersAndFields(tc.input)

			if !reflect.DeepEqual(headers, tc.expectedHeaders) {
				t.Errorf("Expected headers %v, got %v", tc.expectedHeaders, headers)
			}

			if !reflect.DeepEqual(fields, tc.expectedFields) {
				t.Errorf("Expected fields %v, got %v", tc.expectedFields, fields)
			}
		})
	}
}

// TestCreateDataRow tests the createDataRow function
func TestCreateDataRow(t *testing.T) {
	testCases := []struct {
		name           string
		input          interface{}
		fields         []string
		expectedOutput table.Row
	}{
		{name: "Test with Person struct", input: Person{Name: "Alice", Age: 30, Income: 50000, Height: 170}, fields: []string{"Name", "Age", "Income", "Height"}, expectedOutput: table.Row{"Alice", int(30), float64(50000), uint64(170)}},
		{name: "Test with MockFileInfo struct", input: MockFileInfo{Name: "test.txt", Size: 1024}, fields: []string{"Name", "Size"}, expectedOutput: table.Row{"test.txt", "1.00 KB"}},
		{name: "Test with non-existent field", input: Person{Name: "Bob", Age: 25}, fields: []string{"Name", "Age", "NonExistentField"}, expectedOutput: table.Row{"Bob", 25, ""}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output := createDataRow(tc.input, tc.fields)

			if !reflect.DeepEqual(output, tc.expectedOutput) {
				t.Errorf("Expected output %v, got %v", tc.expectedOutput, output)

				// Add detailed type information for debugging
				t.Logf("Expected types: %v", getTypes(tc.expectedOutput))
				t.Logf("Actual types: %v", getTypes(output))
			}
		})
	}
}

// Helper function to get types of slice elements
func getTypes(slice interface{}) []string {
	s := reflect.ValueOf(slice)
	types := make([]string, s.Len())
	for i := 0; i < s.Len(); i++ {
		types[i] = fmt.Sprintf("%T", s.Index(i).Interface())
	}
	return types
}

// Mock FormatSize function (replace with your actual implementation)
func FormatSize(size int64) string {
	return fmt.Sprintf("%.2f KB", float64(size)/1024)
}

// TestGenericSortInterface tests the GenericSortInterface function.
func TestGenericSortInterface(t *testing.T) {
	type InputStruct struct {
		input          interface{}
		sortColumn     string
		sortDescending bool
	}

	type ExpectedResults struct {
		expected interface{}
	}

	tests := []*types.TestLayout[InputStruct, ExpectedResults]{
		{Name: "Sort by string field ascending", Input: InputStruct{input: []Person{{Name: "John", Age: 30}, {Name: "Alice", Age: 25}}, sortColumn: "Name"}, Expected: ExpectedResults{expected: []Person{{Name: "Alice", Age: 25}, {Name: "John", Age: 30}}}},
		{Name: "Sort by string field descending", Input: InputStruct{input: []Person{{Name: "John", Age: 30}, {Name: "Alice", Age: 25}}, sortColumn: "Name", sortDescending: true}, Expected: ExpectedResults{expected: []Person{{Name: "John", Age: 30}, {Name: "Alice", Age: 25}}}},
		{Name: "Sort by int field ascending", Input: InputStruct{input: []Person{{Name: "John", Age: 30}, {Name: "Alice", Age: 25}}, sortColumn: "Age"}, Expected: ExpectedResults{expected: []Person{{Name: "Alice", Age: 25}, {Name: "John", Age: 30}}}},
		{Name: "Sort by int field descending", Input: InputStruct{input: []Person{{Name: "John", Age: 30}, {Name: "Alice", Age: 25}}, sortColumn: "Age", sortDescending: true}, Expected: ExpectedResults{expected: []Person{{Name: "John", Age: 30}, {Name: "Alice", Age: 25}}}},
		{Name: "Sort by float64 field ascending", Input: InputStruct{input: []Person{{Name: "John", Income: 50000.0}, {Name: "Alice", Income: 45000.0}}, sortColumn: "Income"}, Expected: ExpectedResults{expected: []Person{{Name: "Alice", Income: 45000.0}, {Name: "John", Income: 50000.0}}}},
		{Name: "Sort by uint64 field ascending", Input: InputStruct{input: []Person{{Name: "John", Height: 180}, {Name: "Alice", Height: 165}}, sortColumn: "Height"}, Expected: ExpectedResults{expected: []Person{{Name: "Alice", Height: 165}, {Name: "John", Height: 180}}}},
		{Name: "Sort with case-insensitive column name", Input: InputStruct{input: []Person{{Name: "John", Age: 30}, {Name: "Alice", Age: 25}}, sortColumn: "aGe"}, Expected: ExpectedResults{expected: []Person{{Name: "Alice", Age: 25}, {Name: "John", Age: 30}}}},
		{Name: "Sort with non-existent column", Input: InputStruct{input: []Person{{Name: "John", Age: 30}, {Name: "Alice", Age: 25}}, sortColumn: "NonExistent"}, Expected: ExpectedResults{expected: []Person{{Name: "John", Age: 30}, {Name: "Alice", Age: 25}}}},
		{Name: "Sort empty slice", Input: InputStruct{input: []Person{}, sortColumn: "Name"}, Expected: ExpectedResults{expected: []Person{}}},
		{Name: "Sort single-element slice", Input: InputStruct{input: []Person{{Name: "John", Age: 30}}, sortColumn: "Name"}, Expected: ExpectedResults{expected: []Person{{Name: "John", Age: 30}}}},
		{Name: "Sort slice of structs with embedded fields", Input: InputStruct{input: []PersonWithEmbedded{{EmbeddedStruct: EmbeddedStruct{Name: "John", Age: 30}}, {EmbeddedStruct: EmbeddedStruct{Name: "Alice", Age: 25}}}, sortColumn: "Name"}, Expected: ExpectedResults{expected: []PersonWithEmbedded{{EmbeddedStruct: EmbeddedStruct{Name: "Alice", Age: 25}}, {EmbeddedStruct: EmbeddedStruct{Name: "John", Age: 30}}}}},
		{Name: "Non-slice input", Input: InputStruct{input: 123, sortColumn: "Age"}, Expected: ExpectedResults{expected: 123}},
		{Name: "Nil input", Input: InputStruct{input: nil, sortColumn: "Age"}, Expected: ExpectedResults{expected: nil}},
		{Name: "Sort by unhandled type (time.Time)", Input: InputStruct{input: []Event{{Name: "Event 2", Timestamp: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)}, {Name: "Event 1", Timestamp: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)}, {Name: "Event 3", Timestamp: time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)}}, sortColumn: "Timestamp"}, Expected: ExpectedResults{expected: []Event{{Name: "Event 2", Timestamp: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)}, {Name: "Event 1", Timestamp: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)}, {Name: "Event 3", Timestamp: time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)}}}},
	}

	for _, test := range tests {
		t.Run(test.Input.sortColumn, func(t *testing.T) {
			GenericSortInterface(test.Input.input, test.Input.sortColumn, test.Input.sortDescending)

			if !reflect.DeepEqual(test.Input.input, test.Expected.expected) {
				t.Errorf("GenericSortInterface() = %v, want %v", test.Input.input, test.Expected.expected)
			}
		})
	}
}
