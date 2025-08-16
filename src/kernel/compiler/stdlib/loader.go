package stdlib

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// _builtinModules are builtin type standard library modules.
var _builtinModules = map[string]map[string]objects.IObject{
	//"os":   osModule,
	//"fmt":    _fmtSafeModule,
	"fmt":     _fmtModule,
	"math":    _mathModule,
	"strings": _stringsModule,
	"time":    _timeModule,
	"rand":    _randModule,
	"json":    _jsonModule,
	"base64":  _base64Module,
	"hex":     _hexModule,
}

// Module represents a module with predefined attributes that can be imported or accessed at runtime.
// Attrs stores a map of string keys to IObject values, representing the module's predefined attributes.
type Module struct {
	attrs map[string]objects.IObject
}

// NewModule creates and returns a new instance of ModuleBuiltin with the given attributes.
func NewModule(attrs map[string]objects.IObject) *Module {
	return &Module{attrs: attrs}
}

// CompileModule transforms the ModuleBuiltin's attributes into an immutable map, embedding the given module name.
func (m *Module) CompileModule(moduleName string) *objects.MapImmutable {
	attrs := make(map[string]objects.IObject, len(m.attrs))
	for k, v := range m.attrs {
		attrs[k] = v.Copy()
	}
	attrs["__module_name__"] = objects.NewStringNoSize(moduleName)
	return objects.NewMapImmutable(attrs)
}

// Symbol returns the value of the attribute with the given name, if it exists.
func (m *Module) Symbol(name string) (objects.IObject, bool) {
	v, ok := m.attrs[name]
	return v, ok
}

// Loader is responsible for managing modules and resolving symbols within those modules.
// It provides functionality for accessing and resolving both regular and built-in symbols.
// The mod field holds a collection of modules, allowing for import and retrieval of symbols.
type Loader struct {
	modules          map[string]*Module
	builtinFunctions []*objects.FunctionModule
}

// NewLoader initializes and returns a new Loader instance with built-in modules preloaded.
func NewLoader() *Loader {
	modules := make(map[string]*Module)
	for name, mod := range _builtinModules {
		modules[name] = NewModule(mod)
	}
	return &Loader{
		modules:          modules,
		builtinFunctions: append([]*objects.FunctionModule{}, _builtinFunctions...),
	}
}

// GetBuiltinFunctions returns a slice containing all registered builtin functions.
func (l *Loader) GetBuiltinFunctions() []*objects.FunctionModule {
	return l.builtinFunctions
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

// ResolveBuiltinFunction resolves a list of IObject symbols into their corresponding built-in functions, returning an error if any fail.
func (l *Loader) ResolveBuiltinFunction(symbols []objects.IObject) ([]*objects.FunctionModule, error) {
	builtin := make([]*objects.FunctionModule, len(symbols))
	for i := range symbols {
		v := l.GetBuiltinFunction(i)
		if v == nil {
			return nil, fmt.Errorf("can't load builtin symbols, invalid reference %d", i)
		}
		builtin[i] = v
	}
	return builtin, nil
}

// GetBuiltinFunction retrieves a built-in function by its index and returns it as a FunctionBuiltin instance.
func (l *Loader) GetBuiltinFunction(idx int) *objects.FunctionModule {
	if idx < 0 || idx >= len(_builtinFunctions) {
		return nil
	}
	return _builtinFunctions[idx]
}

// GetSymbol retrieves a symbol from the module based on the provided object definition. Returns the symbol and success status.
func (l *Loader) GetSymbol(in objects.IObject) (objects.IObject, bool) {
	definition, ok := in.(*objects.Array)
	if !ok {
		return nil, false
	}
	pName, err := definition.Index(0)
	if err != nil {
		return nil, false
	}
	sName, err := definition.Index(1)
	if err != nil {
		return nil, false
	}
	packageName, ok := pName.(*objects.String)
	if !ok {
		return nil, false
	}
	symbolName, ok := sName.(*objects.String)
	if !ok {
		return nil, false
	}
	mod, ok := l.modules[packageName.Value()]
	if !ok {
		return nil, false
	}
	sym, ok := mod.Symbol(symbolName.Value())
	return sym, ok
}

// CompileModule compiles the specified module by its name, returning an immutable map of its attributes or an error if not found.
func (l *Loader) CompileModule(name string) (*objects.MapImmutable, error) {
	m, ok := l.modules[name]
	if !ok {
		return nil, fmt.Errorf("module %s not found", name)
	}
	return m.CompileModule(name), nil
}
