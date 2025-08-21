package sdk

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Package represents a modular collection of objects, offering access to its resources via a container map.
// It is identified by a unique name and supports storing objects implementing the IObject interface.
type Package struct {
	name      string
	container map[string]objects.IObject
}

// NewPackage creates a new Package instance with the specified name, functions, and constants.
// It initializes the container map and adds provided functions and constants to it.
// name is the name of the package.
// functions is a slice of FuncPackage objects to include in the package.
// constants is a map of constant objects to include in the package, keyed by their IDs.
// Returns the newly created Package instance.
func NewPackage(name string, functions []*objects.FuncPackage, constants map[string]objects.IObject) *Package {
	p := &Package{
		name:      name,
		container: make(map[string]objects.IObject),
	}
	for _, fn := range functions {
		p.container[fn.Name()] = fn
	}
	for id, c := range constants {
		p.container[id] = c
	}
	return p
}

// Name returns the unique name of the Package.
func (p *Package) Name() string {
	return p.name
}

// BuiltinWrapper is a struct designed to wrap a Builtin instance and its associated IObject for additional functionality.
type BuiltinWrapper struct {
	wrapper *objects.Builtin
	object  objects.IObject
}

// Loader represents a mechanism to manage and load packages and built-in objects in the system.
type Loader struct {
	factory  *objects.GateKeeper
	packages map[string]*Package
	builtin  []*BuiltinWrapper
}

// NewLoader initializes and returns a new Loader instance with predefined standard packages and built-in functions.
func NewLoader(factory *objects.GateKeeper) *Loader {
	builtin := NewBuiltinFunctions(factory).Package()
	packages := []*Package{
		NewErrors(factory).Package,
		NewFmt(factory).Package,
		NewMath(factory).Package,
		NewStrings(factory).Package,
		NewStrconv(factory).Package,
		NewRegexp(factory).Package,
		NewTime(factory).Package,
		NewRand(factory).Package,
		NewJson(factory).Package,
		NewBase64(factory).Package,
		NewHex(factory).Package,
	}
	loader := &Loader{
		factory:  factory,
		packages: make(map[string]*Package),
		builtin:  make([]*BuiltinWrapper, len(builtin)),
	}
	for i, fn := range builtin {
		wrapper := factory.NewBuiltin(objects.FrameStatic, fn.Name(), i)
		loader.builtin[i] = &BuiltinWrapper{wrapper: wrapper, object: fn}
	}
	for _, p := range packages {
		loader.packages[p.Name()] = p
	}
	return loader
}

// AddPackage adds a package with the given ID and attributes to the Loader's packages map.
func (l *Loader) AddPackage(id string, attr map[string]objects.IObject) {
	m := &Package{container: attr}
	l.packages[id] = m
}

// BuiltinLen returns the number of built-in objects stored in the Loader instance.
func (l *Loader) BuiltinLen() int {
	return len(l.builtin)
}

// Builtin retrieves a built-in object by its index or returns nil if the index is out of range.
func (l *Loader) Builtin(idx int) *objects.Builtin {
	if idx < 0 || idx >= len(l.builtin) {
		return nil
	}
	return l.builtin[idx].wrapper
}

// BuiltinResolve returns the object associated with the given index from the built-in list or nil if the index is invalid.
func (l *Loader) BuiltinResolve(idx int) objects.IObject {
	if idx < 0 || idx >= len(l.builtin) {
		return nil
	}
	return l.builtin[idx].object
}

// ResolveSymbols resolves a list of symbol references into concrete objects within the loader's context.
// It returns a slice of resolved objects or an error if any reference is invalid.
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

// GetSymbol retrieves a symbol from a package by decoding its reference array and returns the associated object if found.
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
	module, ok := l.packages[packageName.Value()]
	if !ok {
		return nil, false
	}
	v, ok := module.container[symbolName.Value()]
	return v, ok
}

// CompilePackage compiles a package by name, returning an immutable map of its attributes or an error if not found.
func (l *Loader) CompilePackage(name string) (*objects.MapImmutable, error) {
	module, ok := l.packages[name]
	if !ok {
		return nil, fmt.Errorf("module %s not found", name)
	}
	attrs := make(map[string]objects.IObject, len(module.container))
	for k, v := range module.container {
		attrs[k] = v.Copy(objects.FrameStatic, 0)
	}
	attrs[bytecode.ModuleKey] = l.factory.NewString(objects.FrameStatic, name)
	obj := l.factory.NewMapImmutable(objects.FrameStatic, attrs)
	mImm, ok := obj.(*objects.MapImmutable)
	if !ok {
		return nil, fmt.Errorf("can't compile module %s", name)
	}
	return mImm, nil
}
