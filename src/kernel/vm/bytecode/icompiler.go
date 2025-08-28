package bytecode

import "github.com/markel1974/c64emu/src/kernel/vm/objects"

// ICompiler defines an interface for compiling source code and managing constants, references, and global objects.
type ICompiler interface {
	Id() string

	Compile(filename string, source any) error

	FileSet() IFile

	Constants() []objects.IObject

	References() []objects.IObject

	Globals() []objects.IObject
}
