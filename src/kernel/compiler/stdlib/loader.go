package stdlib

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/compiler/modules"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// _builtinModules are builtin type standard library modules.
var _builtinModules = map[string]map[string]objects.IObject{
	//"os":   osModule,
	//"fmt":    _fmtSafeModule,
	"fmt":    _fmtModule,
	"math":   _mathModule,
	"text":   _textModule,
	"times":  _timesModule,
	"rand":   _randModule,
	"json":   _jsonModule,
	"base64": _base64Module,
	"hex":    _hexModule,
}

// GetAllBuiltin returns a slice containing all registered builtin functions.
func GetAllBuiltin() []*objects.FunctionBuiltin {
	return append([]*objects.FunctionBuiltin{}, _builtinFunctions...)
}

// GetBuiltin retrieves a FunctionBuiltin by its index from the predefined list of builtin functions.
func GetBuiltin(idx int) *objects.FunctionBuiltin {
	return _builtinFunctions[idx]
}

// Loader is responsible for managing modules and resolving symbols within those modules.
// It provides functionality for accessing and resolving both regular and built-in symbols.
// The mod field holds a collection of modules, allowing for import and retrieval of symbols.
type Loader struct {
	mod *modules.Modules
}

// NewLoader initializes and returns a new Loader instance with built-in modules preloaded.
func NewLoader() *Loader {
	mods := modules.NewModules()
	for name, mod := range _builtinModules {
		mods.AddBuiltinModule(name, mod)
	}
	return &Loader{
		mod: mods,
	}
}

// ResolveSymbols resolves a list of symbol references to their corresponding objects using the loader's symbol mapping.
func (l *Loader) ResolveSymbols(symbols []objects.IObject) ([]objects.IObject, error) {
	references := make([]objects.IObject, len(symbols))
	for i, ref := range symbols {
		symbol, ok := l.GetSymbol(ref)
		if !ok {
			return nil, fmt.Errorf("can't load symbols, invalid reference %d", i)
		}
		references[i] = symbol
	}
	return references, nil
}

// ResolveBuiltinSymbols resolves a list of IObject symbols into their corresponding built-in functions, returning an error if any fail.
func (l *Loader) ResolveBuiltinSymbols(symbols []objects.IObject) ([]*objects.FunctionBuiltin, error) {
	builtin := make([]*objects.FunctionBuiltin, len(symbols))
	for i := range symbols {
		v := l.GetBuiltinSymbol(i)
		if v == nil {
			return nil, fmt.Errorf("can't load builtin symbols, invalid reference %d", i)
		}
		builtin[i] = v
	}
	return builtin, nil
}

// GetBuiltinSymbol retrieves a built-in function by its index and returns it as a FunctionBuiltin instance.
func (l *Loader) GetBuiltinSymbol(idx int) *objects.FunctionBuiltin {
	return GetBuiltin(idx)
}

// GetSymbol retrieves a symbol from the module based on the provided object definition. Returns the symbol and success status.
func (l *Loader) GetSymbol(definition objects.IObject) (objects.IObject, bool) {
	return l.mod.GetSymbolFromDefinition(definition)
}
