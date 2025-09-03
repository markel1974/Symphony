package bytecode

import "fmt"

// FilePos represents a source file position with filename, offset, line, and column information.
type FilePos struct {
	filename string
	offset   int
	line     int
	column   int
}

// NewFilePos creates a new FilePos instance with the specified filename, offset, line, and column values.
func NewFilePos(filename string, offset int, line int, column int) *FilePos {
	return &FilePos{filename, offset, line, column}
}

// Filename returns the filename associated with the FilePos instance.
func (p FilePos) Filename() string {
	return p.filename
}

// Offset returns the zero-based byte offset associated with the FilePos.
func (p FilePos) Offset() int {
	return p.offset
}

// Line returns the line number stored in the FilePos structure.
func (p FilePos) Line() int {
	return p.line
}

// Column returns the column number of the position, starting at 1.
func (p FilePos) Column() int {
	return p.column
}

// Valid reports whether the FilePos has a valid line number greater than 0.
func (p FilePos) Valid() bool {
	return p.line > 0
}

// String returns a string representation of the FilePos, including filename, line, and column if valid; otherwise, "-".
func (p FilePos) String() string {
	if p.Valid() {
		return fmt.Sprintf("<%s> [line: %d] [col: %d]", p.filename, p.line, p.column)
	}
	return "<>"
}
