package types

// FileType
type FileType string

// OperatorType
type OperatorType string

// SizeUnit struct for units
type SizeUnit struct {
	Label string
	Size  int64
}

// TestLayout
type TestLayout[InputT any, ExpectedT any] struct {
	Name     string
	Input    InputT
	Expected ExpectedT
	Err      error
}

var (
	// FileTypes
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
	// OperatorTypes
	OperatorTypes = struct {
		EqualTo OperatorType
		GreaterThan OperatorType
		GreaterThanEqualTo OperatorType
		LessThan OperatorType
		LessThanEqualTo OperatorType
	}{
		EqualTo: "Equal To",
		GreaterThan: "Greater Than",
		GreaterThanEqualTo: "Greater Than or Equal To",
		LessThan: "Less Than",
		LessThanEqualTo: "Less Than Or Equal To",
	}

	// SizeUnits
	SizeUnits = []SizeUnit{
		{Label: "PB", Size: 1 << 50}, // Petabyte
		{Label: "TB", Size: 1 << 40}, // Terabyte
		{Label: "GB", Size: 1 << 30}, // Gigabyte
		{Label: "MB", Size: 1 << 20}, // Megabyte
		{Label: "KB", Size: 1 << 10}, // Kilobyte
		{Label: "B", Size: 1},        // Byte
	}
	// FileExtensions
	FileExtensions = map[FileType]map[string]bool{
		FileTypes.Any: {
			"*.*": true,
		},
		FileTypes.Video: {
			".mp4": true, ".avi": true, ".mkv": true, ".mov": true, ".wmv": true,
			".flv": true, ".webm": true, ".m4v": true, ".mpg": true, ".mpeg": true,
			".ts": true,
		},
		FileTypes.Image: {
			".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".bmp": true,
			".tiff": true, ".webp": true, ".svg": true, ".raw": true, ".heic": true,
			".ico": true,
		},
		FileTypes.Archive: {
			".zip": true, ".rar": true, ".7z": true, ".tar": true, ".gz": true,
			".bz2": true, ".xz": true, ".iso": true, ".tgz": true, ".tbz2": true,
		},
		FileTypes.Documents: {
			".docx": true, ".doc": true, ".pdf": true, ".txt": true, ".rtf": true,
			".odt": true, ".xlsx": true, ".xls": true, ".pptx": true, ".ppt": true,
			".csv": true, ".md": true, ".pages": true,
		},
	}
)
