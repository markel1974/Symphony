package bytecode

import (
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
func (s *Files) AddFile(f IFile) error {
	if f.Base() < 0 {
		return fmt.Errorf("offset overflow (> 2G of source code in file set)")
	}
	s.base = f.Base()
	s.files = append(s.files, f)
	s.lastFile = f
	return nil
}

// File retrieves the SourceFile containing the given position p, or nil if p is NoPos or not valid for any file.
func (s *Files) File(p int) IFile {
	if p == 0 {
		return nil
	}
	return s.file(p)
}

// Position returns the detailed FilePos for a given Pos, translating it to filename, line, and column information.
// If the position is invalid or cannot be resolved to a specific file, a zero-value FilePos is returned.
func (s *Files) Position(p int) (*FilePos, error) {
	if p == 0 {
		return nil, fmt.Errorf("illegal NoPos value")
	}
	f := s.file(p)
	if f == nil {
		return nil, fmt.Errorf("illegal Pos value")
	}
	return f.Position(p)
}

// file retrieves the SourceFile corresponding to the specified Pos from the Files or nil if not found.
func (s *Files) file(p int) IFile {
	lf := s.lastFile
	if lf != nil && lf.Base() <= p && p <= lf.Base()+lf.Size() {
		return lf
	}
	if i := searchFiles(s.files, p); i >= 0 {
		f := s.files[i]
		if p <= f.Base()+f.Size() {
			s.lastFile = f
			return f
		}
	}
	return nil
}

// searchFiles performs a binary search on a slice of SourceFile pointers to locate the index of the last file with a base offset <= x.
func searchFiles(a []IFile, x int) int {
	return sort.Search(len(a), func(i int) bool { return a[i].Base() > x }) - 1
}
