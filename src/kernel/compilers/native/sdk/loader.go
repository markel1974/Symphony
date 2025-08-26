package sdk

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/compilers/native/common"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// registerBuiltinFn defines a function type for registering built-in functionalities using an IGateKeeper instance.
type registerBuiltinFn func(f objects.IGateKeeper) IBuiltin

// _registerBuiltin holds the function to register built-in functions or objects with an IGateKeeper instance.
var _registerBuiltin registerBuiltinFn

// RegisterBuiltin sets a function used to register built-in functionalities into the system.
func RegisterBuiltin(f registerBuiltinFn) {
	_registerBuiltin = f
}

// registerPackageFn is a function type defining a method that registers a package using IGateKeeper and returns an IPackage.
type registerPackageFn func(f objects.IGateKeeper) IPackage

// _registerPackage holds a list of functions for registering packages, allowing dynamic addition of IPackage instances.
var _registerPackage []registerPackageFn

// RegisterPackage registers a package by appending the provided package registration function to the internal list.
func RegisterPackage(f registerPackageFn) {
	_registerPackage = append(_registerPackage, f)
}

// Loader represents a mechanism to manage and load packages and built-in objects in the system.
type Loader struct {
	gk       objects.IGateKeeper
	packages map[string]IPackage
	builtins []*BuiltinEntry
}

// NewLoader initializes and returns a new Loader instance with predefined standard packages and built-in functions.
func NewLoader(gk objects.IGateKeeper) *Loader {
	builtinContainer := _registerBuiltin(gk).Container()
	packages := make([]IPackage, len(_registerPackage))
	for i, fn := range _registerPackage {
		packages[i] = fn(gk)
	}
	loader := &Loader{
		gk:       gk,
		packages: make(map[string]IPackage),
		builtins: make([]*BuiltinEntry, len(builtinContainer)),
	}
	for i, obj := range builtinContainer {
		fn, ok := obj.(*objects.FuncPackage)
		if !ok {
			continue
		}
		b := gk.NewBuiltin(objects.FrameStatic, fn.Name(), i)
		builtin, ok := b.(*objects.Builtin)
		if !ok {
			continue
		}
		loader.builtins[i] = NewBuiltinAdapter(builtin, fn)
	}
	for _, p := range packages {
		loader.packages[p.Name()] = p
	}
	return loader
}

// Id returns the unique identifier of the loader as defined in the common package.
func (l *Loader) Id() string {
	return common.Identifier
}

// AddPackage adds a package with the given Id and attributes to the Loader's packages map.
func (l *Loader) AddPackage(id string, attr map[string]objects.IObject) {
	l.packages[id] = NewExternalPackage(id, attr)
}

// BuiltinLen returns the number of built-in objects stored in the Loader instance.
func (l *Loader) BuiltinLen() int {
	return len(l.builtins)
}

// Builtin retrieves a built-in object by its index or returns nil if the index is out of range.
func (l *Loader) Builtin(idx int) *objects.Builtin {
	if idx < 0 || idx >= len(l.builtins) {
		return nil
	}
	return l.builtins[idx].builtin
}

// Resolve resolves a list of symbol references into concrete objects within the loader's context.
// It returns a slice of resolved objects or an error if any reference is invalid.
func (l *Loader) Resolve(symbols []objects.IObject) ([]objects.IObject, error) {
	references := make([]objects.IObject, len(symbols))
	for i, ref := range symbols {
		if ref == nil {
			return nil, fmt.Errorf("can't load symbols, invalid reference %d", i)
		}
		switch c := ref.(type) {
		case *objects.Builtin:
			symbol := l.resolveBuiltin(i)
			if symbol == nil {
				return nil, fmt.Errorf("builtin symbol not found: %s", c.Name())
			}
			references[i] = symbol
		default:
			symbol, ok := l.resolveReference(ref)
			if !ok {
				return nil, fmt.Errorf("can't load symbols, invalid reference %d", i)
			}
			references[i] = symbol
		}
	}
	return references, nil
}

// resolveBuiltin returns the object associated with the given index from the built-in list or nil if the index is invalid.
func (l *Loader) resolveBuiltin(idx int) objects.IObject {
	if idx < 0 || idx >= len(l.builtins) {
		return nil
	}
	return l.builtins[idx].object
}

// resolveReference retrieves a symbol from a package by decoding its reference array and returns the associated object if found.
func (l *Loader) resolveReference(in objects.IObject) (objects.IObject, bool) {
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
	module, ok := l.packages[packageName.Value()]
	if !ok {
		return nil, false
	}
	v, ok := module.Get(symbolName.Value())
	return v, ok
}
