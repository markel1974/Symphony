package compiler

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
)

// SourceFile represents a source file in a collection, storing its positional metadata and line offsets.
type SourceFile struct {
	set   *bytecode.Files
	name  string
	base  int
	size  int
	lines []int
}

// Set returns the Files associated with the SourceFile.
func (f *SourceFile) Set() *bytecode.Files {
	return f.set
}

// LineCount returns the total number of lines in the SourceFile.
func (f *SourceFile) LineCount() int {
	return len(f.lines)
}

// AddLine appends a new line offset to the SourceFile if valid, ensuring offsets are in ascending order and within size.
func (f *SourceFile) AddLine(offset int) {
	i := len(f.lines)
	if (i == 0 || f.lines[i-1] < offset) && offset < f.size {
		f.lines = append(f.lines, offset)
	}
}

// LineStart returns the start position (offset) of the specified line within the source file. Panics for invalid line numbers.
func (f *SourceFile) LineStart(line int) int {
	if line < 1 {
		panic("illegal line number (line numbering starts at 1)")
	}
	if line > len(f.lines) {
		panic("illegal line number")
	}
	return f.base + f.lines[line-1]
}

// FileSetPos calculates and returns the absolute position in the file set for the given file-relative offset.
// Panics if the offset exceeds the file size.
func (f *SourceFile) FileSetPos(offset int) int {
	if offset > f.size {
		panic("illegal file offset")
	}
	return f.base + offset
}

// Offset computes the offset of a position p relative to the start of the SourceFile, panic on illegal values.
func (f *SourceFile) Offset(p int) int {
	if int(p) < f.base || int(p) > f.base+f.size {
		panic("illegal SourcePos values")
	}
	return int(p) - f.base
}

// Position returns the position information (filename, offset, line, and column) for a given identifier `p` in the SourceFile.
// It returns an error if `p` is invalid, out of range, or set to `NoPos`.
// The resulting position is represented using a `FilePos` object populated with the corresponding details.
func (f *SourceFile) Position(p int) (*bytecode.FilePos, error) {
	if p == 0 {
		return nil, fmt.Errorf("illegal NoPos value")
	}
	if p < f.base || p > f.base+f.size {
		return nil, fmt.Errorf("illegal SourcePos values")
	}
	offset := p - f.base
	filename, line, column := f.unpack(offset)
	pos := bytecode.NewSourceFilePos(filename, offset, line, column)
	return pos, nil
}

// Base returns the base offset of the SourceFile in its Files.
func (f *SourceFile) Base() int {
	return f.base
}

// Size returns the size of the source file in bytes.
func (f *SourceFile) Size() int {
	return f.size
}

// unpack calculates the filename, line, and column for a given offset in the source file.
func (f *SourceFile) unpack(offset int) (filename string, line, column int) {
	filename = f.name
	if i := searchInts(f.lines, offset); i >= 0 {
		line, column = i+1, offset-f.lines[i]+1
	}
	return
}

// searchInts performs a binary search on a sorted slice of integers `a` to find the largest index `i` where a[i] <= x.
// Returns -1 if no such index exists. The input slice `a` must be sorted in ascending order.
func searchInts(a []int, x int) int {
	// This function body is a manually inlined version of:
	//   return sort.Search(len(a), func(i int) bool { return a[i] > x }) - 1
	i, j := 0, len(a)
	for i < j {
		h := i + (j-i)/2 // avoid overflow when computing h
		// i ≤ h < j
		if a[h] <= x {
			i = h + 1
		} else {
			j = h
		}
	}
	return i - 1
}
