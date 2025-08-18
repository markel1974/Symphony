package sdk

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// _builtinModules are builtin type standard library modules.
var _builtinModules = map[string]map[string]objects.IObject{
	//"os":   osModule,
	//"fmt":    _fmtSafeModule,
	"errors":  _errorsModule,
	"fmt":     _fmtModule,
	"math":    _mathModule,
	"strings": _stringsModule,
	"regexp":  _regexpModule,
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

type Builtin struct {
	wrapper *objects.Builtin
	object  objects.IObject
}

// Loader is responsible for managing modules and resolving symbols within those modules.
// It provides functionality for accessing and resolving both regular and built-in symbols.
// The mod field holds a collection of modules, allowing for import and retrieval of symbols.
type Loader struct {
	modules map[string]*Module
	builtin []*Builtin
}

// NewLoader initializes and returns a new Loader instance with built-in modules preloaded.
func NewLoader() *Loader {
	modules := make(map[string]*Module)
	for name, mod := range _builtinModules {
		modules[name] = &Module{attrs: mod}
	}
	builtin := make([]*Builtin, len(_builtinFunctions))
	for i, fn := range _builtinFunctions {
		wrapper := objects.NewBuiltin(fn.Name(), i)
		builtin[i] = &Builtin{wrapper: wrapper, object: fn}
	}
	return &Loader{
		modules: modules,
		builtin: builtin,
	}
}

func (l *Loader) BuiltinLen() int {
	return len(l.builtin)
}

func (l *Loader) Builtin(idx int) *objects.Builtin {
	if idx < 0 || idx >= len(l.builtin) {
		return nil
	}
	return l.builtin[idx].wrapper
}

// BuiltinResolve retrieves a built-in function by its index and returns it as a FunctionBuiltin instance.
func (l *Loader) BuiltinResolve(idx int) objects.IObject {
	if idx < 0 || idx >= len(l.builtin) {
		return nil
	}
	return l.builtin[idx].object
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
	module, ok := l.modules[packageName.Value()]
	if !ok {
		return nil, false
	}
	v, ok := module.attrs[symbolName.Value()]
	return v, ok
}

// CompileModule compiles the specified module by its name, returning an immutable map of its attributes or an error if not found.
func (l *Loader) CompileModule(name string) (*objects.MapImmutable, error) {
	module, ok := l.modules[name]
	if !ok {
		return nil, fmt.Errorf("module %s not found", name)
	}
	attrs := make(map[string]objects.IObject, len(module.attrs))
	for k, v := range module.attrs {
		attrs[k] = v.Copy()
	}
	attrs[bytecode.ModuleKey] = objects.NewStringNoSize(name)
	return objects.NewMapImmutable(attrs), nil
}
