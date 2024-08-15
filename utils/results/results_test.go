package results

import (
	// "github.com/jedib0t/go-pretty/v6/table"
	// "github.com/stretchr/testify/assert"
	"reflect"
	// "sort"
	"testing"
	"time"

	"github.com/ondrovic/common/types"
)

// Define test cases.
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
// TODO: fix these tests
// func TestGenericRenderResultsTableInterface(t *testing.T) {
// 	// Test with valid slice
// 	slice := []TestStruct{
// 		{Name: "John", Age: 30, Size: 1024},
// 		{Name: "Jane", Age: 25, Size: 2048},
// 	}
// 	totalValues := map[string]interface{}{
// 		"Name": "Total",
// 		"Age":  55,
// 		"Size": int64(3072),
// 	}

// 	// Capture stdout
// 	// You might need to implement a way to capture stdout for testing

// 	GenericRenderResultsTableInterface(slice, totalValues)

// 	// Assert captured stdout contains expected table content
// 	// This part depends on how you capture stdout

// 	// Test with non-slice input
// 	GenericRenderResultsTableInterface("not a slice", nil)
// 	// Assert that error message is printed
// }

// func TestCreateHeaderRow(t *testing.T) {
// 	headers := []string{"Name", "Age", "Size"}
// 	row := createHeaderRow(headers)
// 	assert.Equal(t, table.Row{"Name", "Age", "Size"}, row)
// }

// // func TestCreateDataRow(t *testing.T) {
// // 	data := TestStruct{Name: "John", Age: 30, Size: 1024}
// // 	fields := []string{"Name", "Age", "Size"}
// // 	row := createDataRow(data, fields)
// // 	assert.Equal(t, table.Row{"John", 30, "1.0 KB"}, row)

// // 	// Test with missing field
// // 	fields = append(fields, "NonExistent")
// // 	row = createDataRow(data, fields)
// // 	assert.Equal(t, table.Row{"John", 30, "1.0 KB", ""}, row)
// // }

// func TestCreateFooterRow(t *testing.T) {
// 	headers := []string{"Name", "Age", "Size"}
// 	totalValues := map[string]interface{}{
// 		"Name": "Total",
// 		"Age":  55,
// 		"Size": int64(3072),
// 	}
// 	row := createFooterRow(headers, totalValues)
// 	assert.Equal(t, table.Row{"Total", 55, int64(3072)}, row)

// 	// Test with missing total value
// 	headers = append(headers, "Extra")
// 	row = createFooterRow(headers, totalValues)
// 	assert.Equal(t, table.Row{"Total", 55, int64(3072), ""}, row)
// }

// func TestGetHeadersAndFields(t *testing.T) {
// 	data := TestStruct{Name: "John", Age: 30, Size: 1024}
// 	headers, fields := getHeadersAndFields(data)
// 	assert.Equal(t, []string{"Name", "Age", "Size"}, headers)
// 	assert.Equal(t, []string{"Name", "Age", "Size"}, fields)

// 	// Test with embedded struct
// 	embeddedData := struct {
// 		MockFileInfo
// 		Extra string
// 	}{
// 		MockFileInfo: MockFileInfo{Name: "test.txt", Size: 1024},
// 		Extra:        "extra",
// 	}
// 	headers, fields = getHeadersAndFields(embeddedData)
// 	assert.Equal(t, []string{"Name", "Size", "Extra"}, headers)
// 	assert.Equal(t, []string{"Name", "Size", "Extra"}, fields)

// 	// Test with non-struct input
// 	headers, fields = getHeadersAndFields("not a struct")
// 	assert.Empty(t, headers)
// 	assert.Empty(t, fields)
// }

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
