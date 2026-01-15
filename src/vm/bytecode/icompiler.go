package bytecode

import (
	"github.com/markel1974/symphony/src/vm/objects"
)

// ICompiler defines an interface for compiling source code and managing constants, imports, and globals objects.
type ICompiler interface {
	Compile(filename string, source any) error

	FileSet() *FileSet

	Constants() []objects.IObject

	Imports() []objects.IObject

	Globals() []objects.IObject
}
