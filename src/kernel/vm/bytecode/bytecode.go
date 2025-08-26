package bytecode

import (
	"encoding/gob"
	"io"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// init registers various types with the gob package to enable serialization and deserialization.
func init() {
	gob.Register(&Files{})
}

// Bytecode represents a construct that encapsulates compiled code, associated constants, and object references.
// It aggregates information like source files, constant pool, and referenced objects required for execution.
type Bytecode struct {
	factory    objects.IGateKeeper
	opcodes    *Opcodes
	files      *Files
	constants  []objects.IObject
	references []objects.IObject
	global     []objects.IObject
}

// NewBytecode creates and returns a new instance of Bytecode with an initialized Files object.
func NewBytecode(factory objects.IGateKeeper, op *Opcodes, constants []objects.IObject, references []objects.IObject, global []objects.IObject) *Bytecode {
	return &Bytecode{
		factory:    factory,
		opcodes:    op,
		files:      NewFiles(),
		constants:  constants,
		references: references,
		global:     global,
	}
}

// AddFile adds an IFile to the internal files collection of the Bytecode.
// Returns an error if the file cannot be added.
func (b *Bytecode) AddFile(f IFile) error {
	return b.files.AddFile(f)
}

// Position retrieves the FilePos structure for a given position p in the bytecode's source files.
// Returns an error if the position is invalid or does not map to a specific file.
func (b *Bytecode) Position(p int) (*FilePos, error) {
	return b.files.Position(p)
}

// SourceFiles returns the collection of source files associated with the Bytecode.
func (b *Bytecode) SourceFiles() *Files {
	return b.files
}

// Constants returns the list of constant objects stored within the Bytecode instance.
func (b *Bytecode) Constants() []objects.IObject {
	return b.constants
}

// References retrieves the list of IObject references stored in the Bytecode.
func (b *Bytecode) References() []objects.IObject {
	return b.references
}

// Global retrieves the list of IObject references stored in the Bytecode.
func (b *Bytecode) Global() []objects.IObject {
	return b.global
}

// Encode serializes the Bytecode object and writes it to the provided io.Writer in gob format. Returns an error if encoding fails.
func (b *Bytecode) Encode(w io.Writer) error {
	enc := gob.NewEncoder(w)
	if err := enc.Encode(b.files); err != nil {
		return err
	}
	if err := enc.Encode(b.constants); err != nil {
		return err
	}
	if err := enc.Encode(b.references); err != nil {
		return err
	}
	if err := enc.Encode(b.global); err != nil {
		return err
	}
	return nil
}

// Decode deserializes the bytecode from the given io.Reader and resolves its components using the provided modules.
func (b *Bytecode) Decode(r io.Reader) error {
	dec := gob.NewDecoder(r)
	if err := dec.Decode(&b.files); err != nil {
		return err
	}
	if err := dec.Decode(&b.constants); err != nil {
		return err
	}
	if err := dec.Decode(&b.references); err != nil {
		return err
	}
	if err := dec.Decode(&b.global); err != nil {
		return err
	}
	return nil
}
