package bytecode

import (
	"encoding/gob"

	"github.com/markel1974/c64emu/src/vm/objects"
)

// ContainerType represents an index type used to access and manipulate containers within a Bytecode structure.
type ContainerType int

// ConstantsType represents the index for constant containers.
// ImportsType represents the index for import containers.
// GlobalsType represents the index for global containers.
// LastType indicates the last valid container type and must appear last.
const (
	ConstantsType ContainerType = iota
	ImportsType
	GlobalsType
	LastType //must be the last one
)

// String returns the string representation of a ContainerType value.
func (ct ContainerType) String() string {
	switch ct {
	case ConstantsType:
		return "Constants"
	case ImportsType:
		return "Imports"
	case GlobalsType:
		return "Globals"
	default:
		return "Unknown"
	}
}

// Container represents a struct that holds a named collection of IObject instances.
type Container struct {
	kind ContainerType
	data []objects.IObject
}

// NewContainer creates and returns a new Container instance of the specified ContainerType.
func NewContainer(kind ContainerType) *Container {
	return &Container{
		kind: kind,
		data: nil,
	}
}

// Type returns the type of the container as a ContainerType.
func (c *Container) Type() ContainerType {
	return c.kind
}

// Objects retrieve the slice of IObject instances stored within the Container.
func (c *Container) Objects() []objects.IObject {
	return c.data
}

// Append adds the provided slice of IObject to the Container's existing data slice.
func (c *Container) Append(data []objects.IObject) {
	c.data = append(c.data, data...)
}

// Encode serializes the Container's data slice into the provided gob.Encoder. It returns an error if the operation fails.
func (c *Container) Encode(enc *gob.Encoder) error {
	return enc.Encode(c.data)
}

// Decode deserializes the Container's data using the provided gob.Decoder instance.
func (c *Container) Decode(dec *gob.Decoder) error {
	return dec.Decode(&c.data)
}
