package vm

import "github.com/markel1974/c64emu/src/kernel/vm/objects"

// ILoader represents an interface for loading built-in functions and resolving symbols from definitions.
type ILoader interface {
	GetBuiltinSymbol(idx int) *objects.FunctionBuiltin

	GetSymbol(definition objects.IObject) (objects.IObject, bool)
}
