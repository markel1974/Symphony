package vm

import "github.com/markel1974/c64emu/src/kernel/vm/objects"

// ILoader represents an interface for loading built-in functions and resolving symbols from definitions.
type ILoader interface {
	GetBuiltin(idx int) *objects.FunctionBuiltin

	GetSymbolFromDefinition(definition objects.IObject) (objects.IObject, bool)
}
