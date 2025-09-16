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
	fileSet *token.FileSet
	size    int
}

// NewFileSet creates and returns a new instance of FileSet with the specified size and token.FileSet.
func NewFileSet(fileSet *token.FileSet) *FileSet {
	size := 0
	fileSet.Iterate(func(f *token.File) bool {
		size++
		return true
	})
	return &FileSet{
		size:    size,
		fileSet: fileSet,
	}
}

// Base returns the base offset of the underlying token.FileSet.
func (f *FileSet) Base() int {
	return f.fileSet.Base()
}

// Position converts an integer position to a *bytecode.FilePos, including filename, offset, line, and column information.
func (f *FileSet) Position(p int) token.Position {
	return f.fileSet.Position(token.Pos(p))
}

// Size returns the number of files within the FileSet.
func (f *FileSet) Size() int {
	return f.size
}

// GobEncode serializes the FileSet into a byte slice for use with the gob package, including metadata about contained files.
func (f *FileSet) GobEncode() ([]byte, error) {
	var files []fileData
	f.fileSet.Iterate(func(f *token.File) bool {
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
	if err := encoder.Encode(f.size); err != nil {
		return nil, err
	}
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

	if err := decoder.Decode(&f.size); err != nil {
		return err
	}

	if err := decoder.Decode(&files); err != nil {
		return err
	}
	f.fileSet = token.NewFileSet()
	for _, fd := range files {
		z := f.fileSet.AddFile(fd.Name, fd.Base, fd.Size)
		if fd.Lines != nil {
			z.SetLines(fd.Lines)
		}
	}
	return nil
}
