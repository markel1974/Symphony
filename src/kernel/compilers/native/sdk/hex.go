package sdk

import (
	"encoding/hex"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func init() {
	RegisterPackage(NewHex)
}

// Hex represents a structure encapsulating a module of named IObject functions, typically for hex operations.
type Hex struct {
	factory   objects.IGateKeeper
	container map[string]objects.IObject
}

// NewHex initializes and returns a new Hex instance with a predefined module containing encoding and decoding functions.
func NewHex(factory objects.IGateKeeper) IPackage {
	h := &Hex{
		factory: factory,
	}
	container := []objects.IObject{
		factory.NewFuncPackage(objects.FuncPackageDef, "EncodeToString", factory.FuncIbSOs(hex.EncodeToString)),
		factory.NewFuncPackage(objects.FuncPackageDef, "DecodeString", factory.FuncIsObSe(hex.DecodeString)),
	}
	h.container = BuildContainer(container, nil)
	return h
}

// Name returns the name of the Hex structure, which is "hex".
func (f *Hex) Name() string {
	return "hex"
}

// Get retrieves an object associated with the given name from the container. It returns the object and a boolean indicating success.
func (f *Hex) Get(name string) (objects.IObject, bool) {
	v, ok := f.container[name]
	return v, ok
}
