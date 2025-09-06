package sdk

import (
	"fmt"
	"strings"

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
	gk       objects.IGateKeeper
	packages map[string]IPackage
}

// NewLoader initializes and returns a new Loader instance with predefined standard packages and built-in functions.
func NewLoader(gk objects.IGateKeeper) (*Loader, error) {
	container := make([]IPackage, len(_registerPackage))
	for i, fn := range _registerPackage {
		container[i] = fn(gk)
	}
	loader := &Loader{
		gk:       gk,
		packages: make(map[string]IPackage),
	}
	for _, p := range container {
		loader.packages[p.Name()] = p
	}
	return loader, nil
}

// Id returns the unique identifier of the loader as defined in the common package.
func (l *Loader) Id() string {
	return tables.Identifier
}

// AddPackage adds a package with the given Id and attributes to the Loader's packages map.
func (l *Loader) AddPackage(id string, functions []objects.IObject, constants map[string]objects.IObject) {
	l.packages[id] = NewExternalPackage(id, functions, constants)
}

// Resolve resolves a list of symbol references into concrete objects within the loader's context.
// It returns a slice of resolved objects or an error if any reference is invalid.
func (l *Loader) Resolve(symbols []objects.IObject) ([]objects.IObject, error) {
	references := make([]objects.IObject, len(symbols))
	for i, ref := range symbols {
		references[i] = l.gk.UndefinedValue()
		values := strings.Split(ref.AsString(), ".")
		if len(values) >= 2 {
			module, ok := l.packages[values[0]]
			if !ok {
				return nil, fmt.Errorf("can't load symbols, invalid package %s", values[0])
			}
			symbol, ok := module.Get(values[1])
			if !ok {
				return nil, fmt.Errorf("can't load symbols, invalid symbol %s", values[1])
			}
			references[i] = symbol
		}
	}
	return references, nil
}
