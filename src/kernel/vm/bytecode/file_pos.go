package bytecode

import "fmt"

// FilePos represents a position in a source file with filename, offset, line, and column details.
type FilePos struct {
	filename string // filename, if any
	offset   int    // offset, starting at 0
	line     int    // line number, starting at 1
	column   int    // column number, starting at 1 (byte count)
}

func NewFilePos(filename string, offset int, line int, column int) *FilePos {
	return &FilePos{filename, offset, line, column}
}

func (p FilePos) Filename() string {
	return p.filename
}

func (p FilePos) Offset() int {
	return p.offset
}

func (p FilePos) Line() int {
	return p.line
}

func (p FilePos) Column() int {
	return p.column
}

// IsValid checks if the FilePos instance represents a valid position by ensuring the line number is greater than 0.
func (p FilePos) IsValid() bool {
	return p.line > 0
}

// String returns a string representation of the FilePos, including filename, line, and column if valid.
func (p FilePos) String() string {
	s := p.filename
	if p.IsValid() {
		if s != "" {
			s += ":"
		}
		s += fmt.Sprintf("%d", p.line)
		if p.column != 0 {
			s += fmt.Sprintf(":%d", p.column)
		}
	}
	if s == "" {
		s = "-"
	}
	return s
}
