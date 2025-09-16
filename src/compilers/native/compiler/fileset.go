package compiler

import (
	"encoding/gob"
	"go/token"

	"github.com/markel1974/c64emu/src/vm/bytecode"
)

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
func (f *FileSet) Position(p int) (*bytecode.FilePos, error) {
	pos := f.fSet.Position(token.Pos(p))
	z := bytecode.NewFilePos(pos.Filename, pos.Offset, pos.Line, pos.Column)
	return z, nil
}

// Size returns the number of files within the FileSet.
func (f *FileSet) Size() int {
	return 1
}
