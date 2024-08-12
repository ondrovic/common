package types

import (
	"os"

	"github.com/pterm/pterm"
)

// The `DirOps` interface defines methods for reading directory entries and removing files.
// @property ReadDir - The `ReadDir` method reads the directory named by `name` and returns a list of
// directory entries sorted by filename.
// @property {error} Remove - The `Remove` method in the `DirOps` interface is used to delete a
// directory entry with the specified name. It takes the name of the directory entry as a parameter and
// returns an error if the operation fails.
type DirOps interface {
	ReadDir(name string) ([]os.DirEntry, error)
	Remove(name string) error
}

// The RealDirOps type is likely related to file system operations in the Go programming language.
type RealDirOps struct{}

type FileType string

type OperatorType string

// The type `Application` represents an application with various attributes such as name, description,
// style, usage, and version.
// @property Name - The `Name` property in the `Application` struct is a pointer to a string, which
// likely represents the name of the application.
// @property Description - The `Application` struct has the following properties:
// @property {Styles} Style - The `Style` property in the `Application` struct seems to be of type
// `Styles`. It is likely used to define the visual style or design theme of the application. The
// `Styles` type could be an enum or a custom type that specifies different styles that the application
// can have, such as
// @property Usage - The `Usage` property in the `Application` struct represents the usage or purpose
// of the application. It provides information on how the application is intended to be used or what it
// is designed for.
// @property Version - The `Version` property in the `Application` struct represents the version of the
// application. It is a pointer to a string, which means it can either be `nil` or point to a string
// value containing the version information.
type Application struct {
	Name        string
	Description string
	Style       Styles
	Usage       string
	Version     string
}

// The type Styles contains a field named Color of type Colors.
// @property {Colors} Color - The `Styles` struct has a property called `Color` of type `Colors`.
type Styles struct {
	Color Colors
}

// The type `Colors` defines a structure with two fields, `Background` and `Foreground`, both of type
// `pterm.Color`.
// @property Background - The `Background` property in the `Colors` struct represents the color used
// for the background of a user interface element or text. It is typically used to set the color behind
// the content to provide contrast and improve readability.
// @property Foreground - Foreground is a property of the Colors struct that represents the color used
// for the text or elements in the foreground of a user interface. It is typically the color that is
// most prominent or visible to the user.
type Colors struct {
	Background pterm.Color
	Foreground pterm.Color
}

// The SizeUnit type represents a unit of size with a label and size value.
// @property {string} Label - The `Label` property in the `SizeUnit` struct represents the label or
// name associated with a particular size unit. It is a string type field.
// @property {int64} Size - The `Size` property in the `SizeUnit` struct represents the size value of a
// unit, typically in bytes.
type SizeUnit struct {
	Label string
	Size  int64
}

// The ToleranceResults struct defines the size tolerance and bounds for a value.
// @property {int64} ToleranceSize - ToleranceSize represents the acceptable range or margin of error
// for a particular measurement or value.
// @property {int64} UpperBoundSize - UpperBoundSize represents the upper limit or maximum size allowed
// for a certain parameter or value.
// @property {int64} LowerBoundSize - The `LowerBoundSize` property in the `ToleranceResults` struct
// represents the lower limit or threshold size for a certain tolerance level. It is used to define the
// minimum acceptable size or value within the specified tolerance range.
type ToleranceResults struct {
	ToleranceSize  int64
	UpperBoundSize int64
	LowerBoundSize int64
}

// The TestLayout type is a generic struct used for storing test case information.
// @property {string} Name - The `Name` property in the `TestLayout` struct represents the name or
// description of the test case. It is used to identify and differentiate between different test cases.
// @property {InputT} Input - The `Input` property in the `TestLayout` struct represents the input
// value that will be used for testing a specific functionality or feature. It can be of any type
// specified when defining the `TestLayout` struct.
// @property {ExpectedT} Expected - The `TestLayout` struct has the following properties:
// @property {error} Err - The `Err` property in the `TestLayout` struct is used to store any error
// that may occur during the test execution. It allows you to capture and handle errors that may arise
// while running tests on the input data.
type TestLayout[InputT any, ExpectedT any] struct {
	Name     string
	Input    InputT
	Expected ExpectedT
	Err      error
}

// The `func (r RealDirOps) ReadDir(name string) ([]os.DirEntry, error)` function is a method defined
// on the `RealDirOps` struct. This method is implementing the `ReadDir` function of the `DirOps`
// interface.
func (r RealDirOps) ReadDir(name string) ([]os.DirEntry, error) {
	return os.ReadDir(name)
}

// The `func (r RealDirOps) Remove(name string) error {` function is a method defined on the
// `RealDirOps` struct in Go. This method is implementing the `Remove` function of the `DirOps`
// interface. It takes a `name` parameter which represents the name of the directory entry to be
// removed. Inside the function, it calls the `os.Remove(name)` function which attempts to remove the
// directory entry specified by the `name` parameter. If the operation is successful, it returns `nil`
// indicating no error. If there is an error during the removal operation, it returns an error object
// describing the issue encountered.
func (r RealDirOps) Remove(name string) error {
	return os.Remove(name)
}

var (
	// The `FileTypes` variable is a struct that defines different file types as constants. Each file type
	// is represented by a `FileType` value. The struct initializes these constants with specific string
	// values representing the file types.
	FileTypes = struct {
		Any       FileType
		Video     FileType
		Image     FileType
		Archive   FileType
		Documents FileType
	}{
		Any:       "Any",
		Video:     "Video",
		Image:     "Image",
		Archive:   "Archive",
		Documents: "Documents",
	}

	// The `OperatorTypes` variable is defining a struct that contains different comparison operator types
	// as constants. Each operator type is represented by an `OperatorType` value. The struct initializes
	// these constants with specific string values representing the comparison operators.
	OperatorTypes = struct {
		EqualTo            OperatorType
		GreaterThan        OperatorType
		GreaterThanEqualTo OperatorType
		LessThan           OperatorType
		LessThanEqualTo    OperatorType
	}{
		EqualTo:            "Equal To",
		GreaterThan:        "Greater Than",
		GreaterThanEqualTo: "Greater Than or Equal To",
		LessThan:           "Less Than",
		LessThanEqualTo:    "Less Than Or Equal To",
	}

	// The `SizeUnits` variable is a slice of `SizeUnit` structs that defines different size units along
	// with their corresponding values in bytes. Each `SizeUnit` struct in the slice represents a specific
	// size unit such as Petabyte (PB), Terabyte (TB), Gigabyte (GB), Megabyte (MB), Kilobyte (KB), and
	// Byte (B).
	SizeUnits = []SizeUnit{
		{Label: "PB", Size: 1 << 50}, // Petabyte
		{Label: "TB", Size: 1 << 40}, // Terabyte
		{Label: "GB", Size: 1 << 30}, // Gigabyte
		{Label: "MB", Size: 1 << 20}, // Megabyte
		{Label: "KB", Size: 1 << 10}, // Kilobyte
		{Label: "B", Size: 1},        // Byte
	}

	// The `FileExtensions` variable is a map in Go that associates each `FileType` with a map of file
	// extensions and a boolean value. Here's what it does:.
	FileExtensions = map[FileType]map[string]bool{
		// The `FileTypes.Any: {"*.*": true},` entry in the `FileExtensions` variable is associating the
		// `FileType` constant `Any` with a map of file extensions. In this case, the file extension `*.*` is
		// associated with a boolean value `true`.
		FileTypes.Any: {
			"*.*": true,
		},
		// The `FileTypes.Video` constant is associated with a map of file extensions and boolean values.
		// Each file extension key represents a specific video file format, and the boolean value `true`
		// indicates that files with those extensions are considered to be of the video type.
		FileTypes.Video: {
			".mp4": true, ".avi": true, ".mkv": true, ".mov": true, ".wmv": true,
			".flv": true, ".webm": true, ".m4v": true, ".mpg": true, ".mpeg": true,
			".ts": true,
		},
		// The `FileTypes.Image` constant is associated with a map of file extensions and boolean values.
		// Each file extension key represents a specific image file format, and the boolean value `true`
		// indicates that files with those extensions are considered to be of the image type. In this case,
		// the image file formats include `.jpg`, `.jpeg`, `.png`, `.gif`, `.bmp`, `.tiff`, `.webp`, `.svg`,
		// `.raw`, `.heic`, and `.ico`.
		FileTypes.Image: {
			".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".bmp": true,
			".tiff": true, ".webp": true, ".svg": true, ".raw": true, ".heic": true,
			".ico": true,
		},
		// The `FileTypes.Archive` constant is associated with a map of file extensions and boolean values.
		// Each file extension key represents a specific archive file format, and the boolean value `true`
		// indicates that files with those extensions are considered to be of the archive type. In this case,
		// the archive file formats include `.zip`, `.rar`, `.7z`, `.tar`, `.gz`, `.bz2`, `.xz`, `.iso`,
		// `.tgz`, and `.tbz2`. This mapping allows for easy identification of archive files based on their
		// file extensions within the application or system.
		FileTypes.Archive: {
			".zip": true, ".rar": true, ".7z": true, ".tar": true, ".gz": true,
			".bz2": true, ".xz": true, ".iso": true, ".tgz": true, ".tbz2": true,
		},
		// The `FileTypes.Documents` constant is associated with a map of file extensions and boolean values.
		// Each file extension key represents a specific document file format, and the boolean value `true`
		// indicates that files with those extensions are considered to be of the document type.
		FileTypes.Documents: {
			".docx": true, ".doc": true, ".pdf": true, ".txt": true, ".rtf": true,
			".odt": true, ".xlsx": true, ".xls": true, ".pptx": true, ".ppt": true,
			".csv": true, ".md": true, ".pages": true,
		},
	}
)
