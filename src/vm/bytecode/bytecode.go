package bytecode

import (
	"encoding/gob"
	"go/token"
	"io"

	"github.com/markel1974/symphony/src/vm/objects"
)

// PreInitFunction represents the name of the pre-initialization function in the application.
// InitFunction represents the name of the standard initialization function in the application.
const (
	PreInitFunction = "__init__"
	InitFunction    = "init"
)

// Bytecode represents a compiled state containing source file metadata and categorized object containers.
type Bytecode struct {
	gk          objects.IGateKeeper
	fileHandler *Files
	containers  []*Container
}

// NewBytecodeEmpty creates and returns an empty Bytecode instance with initialized files and container values.
func NewBytecodeEmpty(gk objects.IGateKeeper) *Bytecode {
	bc := &Bytecode{
		gk:          gk,
		fileHandler: NewFiles(),
		containers:  make([]*Container, LastType),
	}
	for i := range bc.containers {
		bc.containers[i] = NewContainer(gk, ContainerType(i))
	}
	return bc
}

// NewBytecode initializes a new Bytecode instance with provided constants, imports, and globals data.
func NewBytecode(gk objects.IGateKeeper, constants []objects.IObject, imports []objects.IObject, globals []objects.IObject, file *FileSet) *Bytecode {
	bc := NewBytecodeEmpty(gk)
	bc.Assign(ConstantsType, constants)
	bc.Assign(ImportsType, imports)
	bc.Assign(GlobalsType, globals)
	bc.AddFile(file)
	return bc
}

// AddFile adds a new IFile instance to the collection of files in the Bytecode.
func (b *Bytecode) AddFile(f *FileSet) {
	b.fileHandler.AddFile(f)
}

// AddFiles appends multiple IFile instances to the Bytecode's internal file collection.
func (b *Bytecode) AddFiles(data []*FileSet) {
	for _, f := range data {
		b.fileHandler.AddFile(f)
	}
}

// Position returns the file position (FilePos) corresponding to a given bytecode offset (p).
// It translates the offset into filename, line, and column information or returns an error if invalid.
func (b *Bytecode) Position(p int) (token.Position, error) {
	return b.fileHandler.Position(p)
}

// Files retrieves a slice of IFile instances associated with the Bytecode.
func (b *Bytecode) Files() []*FileSet {
	return b.fileHandler.Files()
}

// Assign updates the specified container in `b.values` with new data.
func (b *Bytecode) Assign(idx ContainerType, data []objects.IObject) {
	b.containers[idx].data = data
}

// Append adds the provided data slice of objects to the container at the specified index in the Bytecode values.
func (b *Bytecode) Append(idx ContainerType, data []objects.IObject) {
	b.containers[idx].Append(data)
}

// Containers return the slice of Container instances associated with the Bytecode.
func (b *Bytecode) Containers() []*Container {
	return b.containers
}

// Constants retrieves the list of constants stored in the Bytecode instance.
func (b *Bytecode) Constants() []objects.IObject {
	return b.containers[ConstantsType].data
}

// Imports return the list of imported objects stored in the Bytecode instance.
func (b *Bytecode) Imports() []objects.IObject {
	return b.containers[ImportsType].data
}

// Globals returns the list of global objects stored in the bytecode.
func (b *Bytecode) Globals() []objects.IObject {
	return b.containers[GlobalsType].data
}

// Encode serializes the Bytecode instance, writing its data to the provided io.Writer using gob encoding.
func (b *Bytecode) Encode(w io.Writer) error {
	enc := gob.NewEncoder(w)
	for _, container := range b.containers {
		if err := container.Encode(enc); err != nil {
			return err
		}
	}
	if err := b.fileHandler.Encode(enc); err != nil {
		return err
	}
	return nil
}

// Decode deserializes bytecode data from the provided io.Reader into the Bytecode instance using gob decoding.
func (b *Bytecode) Decode(r io.Reader) error {
	dec := gob.NewDecoder(r)
	for _, container := range b.containers {
		if err := container.Decode(dec); err != nil {
			return err
		}
	}
	if err := b.fileHandler.Decode(dec); err != nil {
		return err
	}
	return nil
}
