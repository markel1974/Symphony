package sdk

import (
	"encoding/hex"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Hex represents a structure encapsulating a module of named IObject functions, typically for hex operations.
type Hex struct {
	module map[string]objects.IObject
}

// NewHex initializes and returns a new Hex instance with a predefined module containing encoding and decoding functions.
func NewHex() *Hex {
	h := &Hex{}
	h.module = map[string]objects.IObject{
		"EncodeToString": objects.NewFunctionModule(objects.FunctionModuleDef, "EncodeToString", objects.FuncIbSOs(hex.EncodeToString)),
		"DecodeString":   objects.NewFunctionModule(objects.FunctionModuleDef, "DecodeString", objects.FuncIsObSe(hex.DecodeString)),
	}
	return h
}

// Module returns the map of string keys to IObject values stored in the Hex struct.
func (h *Hex) Module() map[string]objects.IObject {
	return h.module
}
