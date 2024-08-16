package results

import (
	"os"
	"reflect"
	"sort"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/ondrovic/common/utils/formatters"
	"github.com/pterm/pterm"
)

// GenericRenderResultsTableInterface renders a table from a slice of structs or maps.
// It takes a slice of interface{} and a map of string to interface{} representing the total values.
// The function uses reflection to determine the headers and fields based on the struct or map keys.
// It then creates a table, appends the header and rows, and optionally appends the footer with total values.
// Finally, it renders the table using the pterm library.
//
// Example usage:
//
//	type Person struct {
//	    Name string
//	    Age  int
//	}
//
//	people := []Person{
//	    {Name: "John", Age: 30},
//	    {Name: "Jane", Age: 25},
//	}
//
//	totalValues := map[string]interface{}{
//	    "Age": 55,
//	}
func GenericRenderResultsTableInterface(slice interface{}, totalValues map[string]interface{}) {
	value := reflect.ValueOf(slice)
	if value.Kind() != reflect.Slice {
		pterm.Println("Expected a slice")
		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	// Get dynamic headers and fields based on slice type
	headers, fields := getHeadersAndFields(value.Index(0).Interface())

	// Append header
	t.AppendHeader(createHeaderRow(headers))

	// Append rows
	for i := 0; i < value.Len(); i++ {
		t.AppendRow(createDataRow(value.Index(i).Interface(), fields))
	}

	// Append footer if totalValues are provided
	if len(totalValues) > 0 {
		t.AppendFooter(createFooterRow(headers, totalValues))
	}

	t.SetStyle(table.StyleColoredDark)
	t.Render()
}

// GenericSortInterface sorts a slice of structs based on a specified field name.
// It takes three arguments:
//
//	slice: the slice of structs to be sorted
//	sortColumn: the name of the field to sort by
//	sortDescending: a boolean flag indicating whether to sort in descending order
//
// The function uses reflection to access the specified field of each struct in the slice.
// It then sorts the slice using the sort.Slice function and a custom comparison function.
// The comparison function compares the values of the specified field for each pair of structs
// and returns true if the first value is less than the second value (for ascending order),
// or false otherwise. If sortDescending is true, the comparison is reversed.
//
// The function supports sorting fields of type string, int, uint, and float.
// If the specified field is not found or has an unsupported type, the function returns without sorting.
func GenericSortInterface(slice interface{}, sortColumn string, sortDescending bool) {
	value := reflect.ValueOf(slice)
	if value.Kind() != reflect.Slice {
		return
	}

	sort.Slice(slice, func(i, j int) bool {
		vi := value.Index(i)
		vj := value.Index(j)

		fi := reflect.Indirect(vi).FieldByNameFunc(func(name string) bool {
			return strings.EqualFold(name, sortColumn)
		})
		fj := reflect.Indirect(vj).FieldByNameFunc(func(name string) bool {
			return strings.EqualFold(name, sortColumn)
		})

		if !fi.IsValid() || !fj.IsValid() {
			return false
		}

		less := false
		switch fi.Kind() {
		case reflect.String:
			less = fi.String() < fj.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			less = fi.Int() < fj.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			less = fi.Uint() < fj.Uint()
		case reflect.Float32, reflect.Float64:
			less = fi.Float() < fj.Float()
		default:
			return false
		}

		if sortDescending {
			return !less
		}
		return less
	})
}

// createHeaderRow creates a header row for the table.
func createHeaderRow(headers []string) table.Row {
	row := table.Row{}
	for _, header := range headers {
		row = append(row, header)
	}
	return row
}

// createDataRow creates a data row for the table based on struct fields.
func createDataRow(data interface{}, fields []string) table.Row {
	row := table.Row{}
	v := reflect.Indirect(reflect.ValueOf(data))
	for _, field := range fields {
		f := v.FieldByNameFunc(func(name string) bool {
			return strings.EqualFold(name, field)
		})
		if f.IsValid() {
			// Format Size if it's the Size field
			if strings.EqualFold(field, "Size") {
				row = append(row, formatters.FormatSize(f.Int())) // Adjust to your size formatting function
			} else {
				row = append(row, f.Interface())
			}
		} else {
			row = append(row, "")
		}
	}
	return row
}

// createFooterRow creates a footer row based on totalValues.
func createFooterRow(headers []string, totalValues map[string]interface{}) table.Row {
	row := table.Row{}
	for _, header := range headers {
		if totalValue, exists := totalValues[header]; exists {
			row = append(row, totalValue)
		} else {
			row = append(row, "")
		}
	}
	return row
}

// getHeadersAndFields returns the headers and fields of a struct dynamically.
func getHeadersAndFields(data interface{}) (headers []string, fields []string) {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Struct {
		return
	}
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Handle embedded structs (like FileInfo)
		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			embeddedHeaders, embeddedFields := getHeadersAndFields(reflect.New(field.Type).Elem().Interface())
			headers = append(headers, embeddedHeaders...)
			fields = append(fields, embeddedFields...)
		} else {
			headers = append(headers, field.Name)
			fields = append(fields, field.Name)
		}
	}
	return
}
