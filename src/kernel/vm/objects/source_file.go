package objects

// SourceFile represents a source file.
type SourceFile struct {
	// SourceFile set for the file
	set *SourceFileSet
	// SourceFile name as provided to AddFile
	Name string
	// SourcePos values range for this file is [base...base+size]
	Base int
	// SourceFile size as provided to AddFile
	Size int
	// Lines contains the offset of the first character for each line
	// (the first entry is always 0)
	Lines []int
}

// Set returns SourceFileSet.
func (f *SourceFile) Set() *SourceFileSet {
	return f.set
}

// LineCount returns the current number of lines.
func (f *SourceFile) LineCount() int {
	return len(f.Lines)
}

// AddLine adds a new line.
func (f *SourceFile) AddLine(offset int) {
	i := len(f.Lines)
	if (i == 0 || f.Lines[i-1] < offset) && offset < f.Size {
		f.Lines = append(f.Lines, offset)
	}
}

// LineStart returns the position of the first character in the line.
func (f *SourceFile) LineStart(line int) Pos {
	if line < 1 {
		panic("illegal line number (line numbering starts at 1)")
	}
	if line > len(f.Lines) {
		panic("illegal line number")
	}
	return Pos(f.Base + f.Lines[line-1])
}

// FileSetPos returns the position in the file set.
func (f *SourceFile) FileSetPos(offset int) Pos {
	if offset > f.Size {
		panic("illegal file offset")
	}
	return Pos(f.Base + offset)
}

// Offset translates the file set position into the file offset.
func (f *SourceFile) Offset(p Pos) int {
	if int(p) < f.Base || int(p) > f.Base+f.Size {
		panic("illegal SourcePos values")
	}
	return int(p) - f.Base
}

// Position translates the file set position into the file position.
func (f *SourceFile) Position(p Pos) (pos SourceFilePos) {
	if p != NoPos {
		if int(p) < f.Base || int(p) > f.Base+f.Size {
			panic("illegal SourcePos values")
		}
		pos = f.position(p)
	}
	return
}

func (f *SourceFile) position(p Pos) (pos SourceFilePos) {
	offset := int(p) - f.Base
	pos.Offset = offset
	pos.Filename, pos.Line, pos.Column = f.unpack(offset)
	return
}

func (f *SourceFile) unpack(offset int) (filename string, line, column int) {
	filename = f.Name
	if i := searchInts(f.Lines, offset); i >= 0 {
		line, column = i+1, offset-f.Lines[i]+1
	}
	return
}

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
