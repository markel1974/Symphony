package bytecode

import (
	"bytes"
	"encoding/gob"
	"go/token"
)

type fileData struct {
	Name  string
	Base  int
	Size  int
	Lines []int
}

func init() {
	gob.Register(&FileSet{})
}

// FileSet wraps a token.FileSet and provides methods to interact with file position and size metadata.
type FileSet struct {
	fSet *token.FileSet
}

// NewFileSet creates and returns a new instance of FileSet with the specified size and token.FileSet.
func NewFileSet(fSet *token.FileSet) *FileSet {
	return &FileSet{
		fSet: fSet,
	}
}

// Base returns the base offset of the underlying token.FileSet.
func (f *FileSet) Base() int {
	return f.fSet.Base()
}

// Position converts an integer position to a *bytecode.FilePos, including filename, offset, line, and column information.
func (f *FileSet) Position(p int) (*FilePos, error) {
	pos := f.fSet.Position(token.Pos(p))
	z := NewFilePos(pos.Filename, pos.Offset, pos.Line, pos.Column)
	return z, nil
}

// Size returns the number of files within the FileSet.
func (f *FileSet) Size() int {
	return 1
}

// GobEncode serializes the FileSet into a byte slice for use with the gob package, including metadata about contained files.
func (f *FileSet) GobEncode() ([]byte, error) {
	var files []fileData
	f.fSet.Iterate(func(f *token.File) bool {
		files = append(files, fileData{
			Name:  f.Name(),
			Base:  f.Base(),
			Size:  f.Size(),
			Lines: f.Lines(),
		})
		return true
	})
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(files); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GobDecode deserializes a FileSet object from a byte slice, restoring its internal file metadata.
func (f *FileSet) GobDecode(data []byte) error {
	var files []fileData
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)

	if err := decoder.Decode(&files); err != nil {
		return err
	}
	f.fSet = token.NewFileSet()
	for _, fd := range files {
		z := f.fSet.AddFile(fd.Name, fd.Base, fd.Size)
		if fd.Lines != nil {
			z.SetLines(fd.Lines)
		}
	}
	return nil
}
