package stdlib

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/compiler/modules"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Loader provides access to and management of modules, including built-in and source modules, for symbol resolution and imports.
type Loader struct {
	mod *modules.Modules
}

// NewLoader initializes and returns a new Loader instance with built-in modules preloaded.
func NewLoader() *Loader {
	mods := modules.NewModules()
	for name, mod := range BuiltinModules {
		mods.AddBuiltinModule(name, mod)
	}
	return &Loader{
		mod: mods,
	}
}

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

func (l *Loader) GetBuiltinSymbol(idx int) *objects.FunctionBuiltin {
	return GetBuiltin(idx)
}

func (l *Loader) GetSymbol(definition objects.IObject) (objects.IObject, bool) {
	return l.mod.GetSymbolFromDefinition(definition)
}
