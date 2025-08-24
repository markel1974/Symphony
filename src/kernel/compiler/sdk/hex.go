package sdk

import (
	"encoding/hex"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Hex represents a structure encapsulating a module of named IObject functions, typically for hex operations.
type Hex struct {
	factory objects.IGateKeeper
	*Package
}

// NewHex initializes and returns a new Hex instance with a predefined module containing encoding and decoding functions.
func NewHex(factory objects.IGateKeeper) *Hex {
	h := &Hex{
		factory: factory,
	}
	container := []objects.IObject{
		factory.NewFuncPackage(objects.FuncPackageDef, "EncodeToString", factory.FuncIbSOs(hex.EncodeToString)),
		factory.NewFuncPackage(objects.FuncPackageDef, "DecodeString", factory.FuncIsObSe(hex.DecodeString)),
	}
	h.Package = NewPackage("hex", container, nil)
	return h
}
