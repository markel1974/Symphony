package bytecode

import (
	"github.com/markel1974/c64emu/src/vm/objects"
)

// ICompiler defines an interface for compiling source code and managing constants, imports, and globals objects.
type ICompiler interface {
	Id() string

	Compile(filename string, source any) error

	FileSet() IFile

	Constants() []objects.IObject

	Imports() []objects.IObject

	Globals() []objects.IObject
}
