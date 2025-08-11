package objects

import "sort"

// Pos represents a position in the file set.
type Pos int

// NoPos represents an invalid position.
const NoPos Pos = 0

// IsValid returns true if the position is valid.
func (p Pos) IsValid() bool {
	return p != NoPos
}

// SourceFileSet represents a set of source files.
type SourceFileSet struct {
	Base     int           // base offset for the next file
	Files    []*SourceFile // list of files in the order added to the set
	LastFile *SourceFile   // cache of last file looked up
}

// NewFileSet creates a new file set.
func NewFileSet() *SourceFileSet {
	return &SourceFileSet{
		Base: 1, // 0 == NoPos
	}
}

// AddFile adds a new file in the file set.
func (s *SourceFileSet) AddFile(filename string, base, size int) *SourceFile {
	if base < 0 {
		base = s.Base
	}
	if base < s.Base || size < 0 {
		panic("illegal base or size")
	}
	f := &SourceFile{
		set:   s,
		Name:  filename,
		Base:  base,
		Size:  size,
		Lines: []int{0},
	}
	base += size + 1 // +1 because EOF also has a position
	if base < 0 {
		panic("offset overflow (> 2G of source code in file set)")
	}

	// add the file to the file set
	s.Base = base
	s.Files = append(s.Files, f)
	s.LastFile = f
	return f
}

// File returns the file that contains the position p. If no such file is
// found (for instance for p == NoPos), the result is nil.
func (s *SourceFileSet) File(p Pos) (f *SourceFile) {
	if p != NoPos {
		f = s.file(p)
	}
	return
}

// Position converts a SourcePos p in the fileset into a SourceFilePos value.
func (s *SourceFileSet) Position(p Pos) (pos SourceFilePos) {
	if p != NoPos {
		if f := s.file(p); f != nil {
			return f.position(p)
		}
	}
	return
}

func (s *SourceFileSet) file(p Pos) *SourceFile {
	// common case: p is in last file
	f := s.LastFile
	if f != nil && f.Base <= int(p) && int(p) <= f.Base+f.Size {
		return f
	}

	// p is not in last file - search all files
	if i := searchFiles(s.Files, int(p)); i >= 0 {
		f := s.Files[i]

		// f.base <= int(p) by definition of searchFiles
		if int(p) <= f.Base+f.Size {
			s.LastFile = f // race is ok - s.last is only a cache
			return f
		}
	}
	return nil
}

func searchFiles(a []*SourceFile, x int) int {
	return sort.Search(len(a), func(i int) bool { return a[i].Base > x }) - 1
}
