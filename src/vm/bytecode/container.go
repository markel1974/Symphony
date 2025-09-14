package bytecode

import "github.com/markel1974/c64emu/src/vm/objects"

// containerType represents an enumeration used to index and manage container data structures.
type containerType int

// constantsType represents a container for constants.
// importsType represents a container for imports.
// globalsType represents a container for global variables.
// lastType indicates the last valid container type and must always be the final value in the enumeration.
const (
	constantsType containerType = iota
	importsType
	globalsType
	lastType //must be the last one
)

// Container represents a structure for holding a collection of IObjects along with a descriptive name.
type Container struct {
	name string
	data []objects.IObject
}

// Name returns the name of the container as a string.
func (c *Container) Name() string {
	return c.name
}

// Data retrieves the slice of IObject instances stored within the Container.
func (c *Container) Data() []objects.IObject {
	return c.data
}

// Append adds the provided slice of IObject data to the Container's existing data slice.
func (c *Container) Append(data []objects.IObject) {
	c.data = append(c.data, data...)
}

// Values returns all objects stored in the container as a slice of IObject.
func (c *Container) Values() []objects.IObject {
	return c.data
}

// ContainerData represents a collection of categorized data containers and associated source files.
type ContainerData struct {
	values      []Container
	sourceFiles []IFile
}

// NewContainerData creates and initializes a ContainerData instance with constants, imports, and globals from the given Bytecode.
func NewContainerData(b *Bytecode) *ContainerData {
	var constants []objects.IObject
	var imports []objects.IObject
	var globals []objects.IObject
	if b != nil {
		constants = b.Constants()
		imports = b.Imports()
		globals = b.Globals()
	}
	out := &ContainerData{
		values: make([]Container, lastType),
	}
	out.values[constantsType] = Container{name: "Constants", data: constants}
	out.values[importsType] = Container{name: "Imports", data: imports}
	out.values[globalsType] = Container{name: "Globals", data: globals}
	return out
}

// Bytecode constructs and returns a new Bytecode instance, aggregating data from constants, imports, globals, and source files.
func (cd *ContainerData) Bytecode() *Bytecode {
	out := NewBytecode(cd.values[constantsType].Data(), cd.values[importsType].Data(), cd.values[globalsType].Data())
	for _, sf := range cd.sourceFiles {
		out.AddFile(sf)
	}
	return out
}

// Assign updates the data of a specified container type within the ContainerData with the provided slice of IObject instances.
func (cd *ContainerData) Assign(idx containerType, data []objects.IObject) {
	cd.values[idx].data = data
}

// Append adds the provided slice of IObject to the Container at the specified containerType index in the ContainerData.
func (cd *ContainerData) Append(idx containerType, data []objects.IObject) {
	cd.values[idx].Append(data)
}

// AppendSourceFiles appends a slice of IFile instances to the sourceFiles field of the ContainerData instance.
func (cd *ContainerData) AppendSourceFiles(data []IFile) {
	cd.sourceFiles = append(cd.sourceFiles, data...)
}

// Data retrieves a slice of IObject instances from the Container at the specified containerType index.
func (cd *ContainerData) Data(idx containerType) []objects.IObject {
	return cd.values[idx].Data()
}

// Values return a slice of Container objects stored in the ContainerData instance.
func (cd *ContainerData) Values() []Container {
	return cd.values
}
