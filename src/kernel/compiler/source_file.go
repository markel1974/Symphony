package compiler

// SourceFile represents a source file in a SourceFileSet.
// It contains metadata such as name, size, base offset, and line offsets.
type SourceFile struct {
	set   *SourceFileSet
	name  string
	base  int
	size  int
	lines []int
}

// Set returns the SourceFileSet associated with the SourceFile.
func (f *SourceFile) Set() *SourceFileSet {
	return f.set
}

// LineCount returns the total number of lines in the source file.
func (f *SourceFile) LineCount() int {
	return len(f.lines)
}

// AddLine adds a new line offset to the source file if it is valid and greater than the previous line offset.
func (f *SourceFile) AddLine(offset int) {
	i := len(f.lines)
	if (i == 0 || f.lines[i-1] < offset) && offset < f.size {
		f.lines = append(f.lines, offset)
	}
}

// LineStart returns the Pos value of the first character in the specified line number of the SourceFile.
// The provided line number must be greater than 0 and within the range of valid lines in the file.
func (f *SourceFile) LineStart(line int) Pos {
	if line < 1 {
		panic("illegal line number (line numbering starts at 1)")
	}
	if line > len(f.lines) {
		panic("illegal line number")
	}
	return Pos(f.base + f.lines[line-1])
}

// FileSetPos returns the positional value for a given offset in the source file.
// It panics if the provided offset exceeds the file size.
func (f *SourceFile) FileSetPos(offset int) Pos {
	if offset > f.size {
		panic("illegal file offset")
	}
	return Pos(f.base + offset)
}

// Offset returns the absolute offset for a given position within the SourceFile. Panics if the position is invalid.
func (f *SourceFile) Offset(p Pos) int {
	if int(p) < f.base || int(p) > f.base+f.size {
		panic("illegal SourcePos values")
	}
	return int(p) - f.base
}

// Position translates a Pos value into a SourceFilePos, containing detailed position information within the file.
// It panics if the given position is invalid or out of range for the source file.
func (f *SourceFile) Position(p Pos) (pos SourceFilePos) {
	if p != NoPos {
		if int(p) < f.base || int(p) > f.base+f.size {
			panic("illegal SourcePos values")
		}
		pos = f.position(p)
	}
	return
}

// position computes the SourceFilePos from a given Pos by unpacking its offset relative to the SourceFile base.
func (f *SourceFile) position(p Pos) (pos SourceFilePos) {
	offset := int(p) - f.base
	pos.offset = offset
	pos.filename, pos.line, pos.column = f.unpack(offset)
	return
}

// unpack determines the filename, line, and column for a given offset in the source file.
// It returns the file name, line number (starting at 1), and column number (starting at 1).
// The method uses searchInts to locate the corresponding line in the offset data.
func (f *SourceFile) unpack(offset int) (filename string, line, column int) {
	filename = f.name
	if i := searchInts(f.lines, offset); i >= 0 {
		line, column = i+1, offset-f.lines[i]+1
	}
	return
}

// searchInts searches for the largest index `i` in a sorted slice `a` such that `a[i] <= x` and returns that index.
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
