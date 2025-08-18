package sdk

import (
	"encoding/hex"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Hex represents a structure encapsulating a module of named IObject functions, typically for hex operations.
type Hex struct {
	*Package
}

// NewHex initializes and returns a new Hex instance with a predefined module containing encoding and decoding functions.
func NewHex() *Hex {
	h := &Hex{}
	container := []*objects.FuncPackage{
		objects.NewFuncPackage(objects.FuncPackageDef, "EncodeToString", objects.FuncIbSOs(hex.EncodeToString)),
		objects.NewFuncPackage(objects.FuncPackageDef, "DecodeString", objects.FuncIsObSe(hex.DecodeString)),
	}
	h.Package = NewPackage("hex", container, nil)
	return h
}
