package compiler

import "fmt"

// SourceFilePos represents a position in a source file with filename, offset, line, and column details.
type SourceFilePos struct {
	filename string // filename, if any
	offset   int    // offset, starting at 0
	line     int    // line number, starting at 1
	column   int    // column number, starting at 1 (byte count)
}

// IsValid checks if the SourceFilePos instance represents a valid position by ensuring the line number is greater than 0.
func (p SourceFilePos) IsValid() bool {
	return p.line > 0
}

// String returns a string representation of the SourceFilePos, including filename, line, and column if valid.
func (p SourceFilePos) String() string {
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
