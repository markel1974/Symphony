package bytecode

import (
	"encoding/gob"
	"fmt"
	"go/token"
	"sort"
)

func init() {
	gob.Register(&Files{})
}

// Files represents a collection of source files with positional metadata for efficient file and position lookups.
type Files struct {
	lastFile int
	files    []*FileSet
}

// NewFiles creates and returns a new instance of Files initialized with a base offset of 1.
func NewFiles() *Files {
	return &Files{lastFile: -1, files: []*FileSet{}}
}

// AddFile adds a new source file with the specified filename, base offset, and size to the Files.
// Returns a pointer to the created SourceFile. Panics if base or size is invalid or if offset overflows.
func (s *Files) AddFile(f *FileSet) {
	if f == nil {
		return
	}
	if f.Base() < 0 {
		return // > 2G
	}
	s.lastFile = len(s.files)
	s.files = append(s.files, f)
}

// Files returns a slice of all SourceFiles added to the Files.
func (s *Files) Files() []*FileSet {
	return s.files
}

// File retrieves the SourceFile corresponding to the specified Pos from the Files or nil if not found.
func (s *Files) File(p int) *FileSet {
	if p == 0 {
		return nil
	}
	if s.lastFile >= 0 {
		lf := s.files[s.lastFile]
		if lf.Base() <= p && p <= lf.Base()+lf.Size() {
			return lf
		}
	}
	if i := s.search(p); i >= 0 {
		f := s.files[i]
		if p <= f.Base()+f.Size() {
			s.lastFile = i
			return f
		}
	}
	return nil
}

// Position returns the detailed FilePos for a given Pos, translating it to filename, line, and column information.
// If the position is invalid or cannot be resolved to a specific file, a zero-value FilePos is returned.
func (s *Files) Position(p int) (token.Position, error) {
	f := s.File(p)
	if f == nil {
		return token.Position{}, fmt.Errorf("illegal pos")
	}
	return f.Position(p), nil
}

// searchFiles performs a binary search on a slice of SourceFile pointers to locate the index of the last file with a base offset <= x.
func (s *Files) search(x int) int {
	return sort.Search(len(s.files), func(i int) bool {
		return s.files[i].Base() > x
	}) - 1
}

// Encode serializes the Files structure, including its base, files, and lastFile properties using the provided gob.Encoder.
func (s *Files) Encode(enc *gob.Encoder) error {
	if err := enc.Encode(s.lastFile); err != nil {
		return err
	}
	if err := enc.Encode(s.files); err != nil {
		return err
	}
	return nil
}

// Decode deserializes the Files object, restoring its base, files slice, and lastFile from the provided gob.Decoder.
func (s *Files) Decode(dec *gob.Decoder) error {
	if err := dec.Decode(&s.lastFile); err != nil {
		return err
	}
	if err := dec.Decode(&s.files); err != nil {
		return err
	}
	return nil
}
