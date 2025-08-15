package bytecode

import "github.com/markel1974/c64emu/src/kernel/vm/objects"

// ILoader defines an interface for loading and resolving symbols and built-in functions required by the system.
// GetBuiltinSymbol retrieves a built-in function by its index and returns a pointer to FunctionBuiltin.
// GetSymbol resolves a symbol based on a definition and returns the resolved object along with a success status.
// ResolveBuiltinSymbols resolves a slice of objects into their corresponding built-in function representations.
// ResolveSymbols resolves a slice of objects into their interpreted or defined values.
type ILoader interface {
	GetBuiltinSymbol(idx int) *objects.FunctionBuiltin

	GetSymbol(definition objects.IObject) (objects.IObject, bool)

	ResolveBuiltinSymbols([]objects.IObject) ([]*objects.FunctionBuiltin, error)

	ResolveSymbols([]objects.IObject) ([]objects.IObject, error)

	CompileModule(name string) (*objects.MapImmutable, error)
}
