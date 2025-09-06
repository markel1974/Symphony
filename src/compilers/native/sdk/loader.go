package sdk

import (
	"fmt"

	"github.com/markel1974/c64emu/src/compilers/native/tables"
	"github.com/markel1974/c64emu/src/vm/objects"
)

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
	gk             objects.IGateKeeper
	packagesByName map[string]IPackage
	packagesByID   map[uint32]map[uint32]IPackage
}

// NewLoader initializes and returns a new Loader instance with predefined standard packages and built-in functions.
func NewLoader(gk objects.IGateKeeper) (*Loader, error) {
	packages := make([]IPackage, len(_registerPackage))
	for i, fn := range _registerPackage {
		packages[i] = fn(gk)
	}
	loader := &Loader{
		gk:             gk,
		packagesByName: make(map[string]IPackage),
		packagesByID:   make(map[uint32]map[uint32]IPackage),
	}
	for _, p := range packages {
		pId := PackageIDFromString(p.Name())
		if _, exists := loader.packagesByID[pId]; exists {
			return nil, fmt.Errorf("hash collision detected for package ID %d", pId)
		}
		loader.packagesByName[p.Name()] = p
		z := make(map[uint32]IPackage)
		loader.packagesByID[pId] = z
		//for k, v := range p {
		//	z[k] = v
		//}
		//loader.packagesByID[pId] = p
	}
	return loader, nil
}

// Id returns the unique identifier of the loader as defined in the common package.
func (l *Loader) Id() string {
	return tables.Identifier
}

// AddPackage adds a package with the given Id and attributes to the Loader's packages map.
func (l *Loader) AddPackage(id string, functions []objects.IObject, constants map[string]objects.IObject) {
	l.packagesByName[id] = NewExternalPackage(id, functions, constants)
}

// Resolve resolves a list of symbol references into concrete objects within the loader's context.
// It returns a slice of resolved objects or an error if any reference is invalid.
func (l *Loader) Resolve(symbols []objects.IObject) ([]objects.IObject, error) {
	references := make([]objects.IObject, len(symbols))
	for i, ref := range symbols {
		if ref == nil {
			return nil, fmt.Errorf("can't load symbols, invalid reference %d", i)
		}
		switch ref.(type) {
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
	module, ok := l.packagesByName[packageName.Value()]
	if !ok {
		return nil, false
	}
	v, ok := module.Get(symbolName.Value())
	return v, ok
}
