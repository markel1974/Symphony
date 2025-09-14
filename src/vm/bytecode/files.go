package bytecode

import (
	"encoding/gob"
	"fmt"
	"sort"
)

// Files represents a collection of source files with positional metadata for efficient file and position lookups.
type Files struct {
	base     int
	files    []IFile
	lastFile IFile
}

// NewFiles creates and returns a new instance of Files initialized with a base offset of 1.
func NewFiles() *Files {
	return &Files{base: 1}
}

// AddFile adds a new source file with the specified filename, base offset, and size to the Files.
// Returns a pointer to the created SourceFile. Panics if base or size is invalid or if offset overflows.
func (s *Files) AddFile(f IFile) {
	if f == nil {
		return
	}
	if f.Base() < 0 {
		return //overflow > 2G
	}
	s.base = f.Base()
	s.files = append(s.files, f)
	s.lastFile = f
}

// Files returns a slice of all SourceFiles added to the Files.
func (s *Files) Files() []IFile {
	return s.files
}

// File retrieves the SourceFile corresponding to the specified Pos from the Files or nil if not found.
func (s *Files) File(p int) IFile {
	if p == 0 {
		return nil
	}
	lf := s.lastFile
	if lf != nil && lf.Base() <= p && p <= lf.Base()+lf.Size() {
		return lf
	}
	if i := s.search(p); i >= 0 {
		f := s.files[i]
		if p <= f.Base()+f.Size() {
			s.lastFile = f
			return f
		}
	}
	return nil
}

// Position returns the detailed FilePos for a given Pos, translating it to filename, line, and column information.
// If the position is invalid or cannot be resolved to a specific file, a zero-value FilePos is returned.
func (s *Files) Position(p int) (*FilePos, error) {
	f := s.File(p)
	if f == nil {
		return nil, fmt.Errorf("illegal Pos value")
	}
	return f.Position(p)
}

// searchFiles performs a binary search on a slice of SourceFile pointers to locate the index of the last file with a base offset <= x.
func (s *Files) search(x int) int {
	return sort.Search(len(s.files), func(i int) bool {
		return s.files[i].Base() > x
	}) - 1
}

// Encode serializes the Files structure including its base, files, and lastFile properties using the provided gob.Encoder.
func (s *Files) Encode(enc *gob.Encoder) error {
	if err := enc.Encode(s.base); err != nil {
		return err
	}
	if err := enc.Encode(s.files); err != nil {
		return err
	}
	if err := enc.Encode(s.lastFile); err != nil {
		return err
	}
	return nil
}

// Decode deserializes the Files object, restoring its base, files slice, and lastFile from the provided gob.Decoder.
func (s *Files) Decode(dec *gob.Decoder) error {
	if err := dec.Decode(&s.base); err != nil {
		return err
	}
	if err := dec.Decode(&s.files); err != nil {
		return err
	}
	if err := dec.Decode(&s.lastFile); err != nil {
		return err
	}
	return nil
}
