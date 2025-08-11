package compiler

import "sort"

// Pos represents a source code position as an offset value within the file set.
type Pos int

// NoPos denotes an invalid position and is used as a default or sentinel value for unset positions.
const NoPos Pos = 0

// IsValid reports whether the position is valid, meaning it is not equal to NoPos.
func (p Pos) IsValid() bool {
	return p != NoPos
}

// SourceFileSet represents a collection of source files with positional metadata for efficient file and position lookups.
type SourceFileSet struct {
	base     int           // base offset for the next file
	files    []*SourceFile // list of files in the order added to the set
	lastFile *SourceFile   // cache of last file looked up
}

// NewFileSet creates and returns a new instance of SourceFileSet initialized with a base offset of 1.
func NewFileSet() *SourceFileSet {
	return &SourceFileSet{
		base: 1, // 0 == NoPos
	}
}

// AddFile adds a new source file with the specified filename, base offset, and size to the SourceFileSet.
// Returns a pointer to the created SourceFile. Panics if base or size is invalid or if offset overflows.
func (s *SourceFileSet) AddFile(filename string, base, size int) *SourceFile {
	if base < 0 {
		base = s.base
	}
	if base < s.base || size < 0 {
		panic("illegal base or size")
	}
	f := &SourceFile{
		set:   s,
		name:  filename,
		base:  base,
		size:  size,
		lines: []int{0},
	}
	base += size + 1
	if base < 0 {
		panic("offset overflow (> 2G of source code in file set)")
	}
	s.base = base
	s.files = append(s.files, f)
	s.lastFile = f
	return f
}

// File retrieves the SourceFile containing the given position p, or nil if p is NoPos or not valid for any file.
func (s *SourceFileSet) File(p Pos) (f *SourceFile) {
	if p != NoPos {
		f = s.file(p)
	}
	return
}

// Position returns the detailed SourceFilePos for a given Pos, translating it to filename, line, and column information.
// If the position is invalid or cannot be resolved to a specific file, a zero-value SourceFilePos is returned.
func (s *SourceFileSet) Position(p Pos) (pos SourceFilePos) {
	if p != NoPos {
		if f := s.file(p); f != nil {
			return f.position(p)
		}
	}
	return
}

// file retrieves the SourceFile corresponding to the specified Pos from the SourceFileSet or nil if not found.
func (s *SourceFileSet) file(p Pos) *SourceFile {
	// common case: p is in last file
	lf := s.lastFile
	if lf != nil && lf.base <= int(p) && int(p) <= lf.base+lf.size {
		return lf
	}
	if i := searchFiles(s.files, int(p)); i >= 0 {
		f := s.files[i]
		if int(p) <= f.base+f.size {
			s.lastFile = f
			return f
		}
	}
	return nil
}

// searchFiles performs a binary search on a slice of SourceFile pointers to locate the index of the last file with a base offset <= x.
func searchFiles(a []*SourceFile, x int) int {
	return sort.Search(len(a), func(i int) bool { return a[i].base > x }) - 1
}
