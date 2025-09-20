package bytecode

import (
	"fmt"
	"strings"

	"github.com/markel1974/c64emu/src/vm/objects"
)

// Loader represents a mechanism to manage and load packages and built-in objects in the system.
type Loader struct {
	gk       objects.IGateKeeper
	packages map[string]IPackage
}

// NewLoader initializes and returns a new Loader instance with predefined standard packages and built-in functions.
func NewLoader(gk objects.IGateKeeper) *Loader {
	loader := &Loader{
		gk:       gk,
		packages: make(map[string]IPackage),
	}
	return loader
}

func (l *Loader) RegisterPackage(registerPackage []RegisterPackageFn) error {
	container := make([]IPackage, len(registerPackage))
	for i, fn := range registerPackage {
		container[i] = fn(l.gk)
	}
	for _, p := range container {
		if z, ok := l.packages[p.Name()]; ok {
			return fmt.Errorf("package %s => %s already exists", z.Name(), p.Name())
		}
		l.packages[p.Name()] = p
	}
	return nil
}

// AddPackage adds a package with the given Id and attributes to the Loader's packages map.
func (l *Loader) AddPackage(id string, functions []objects.IObject, constants map[string]objects.IObject) {
	l.packages[id] = NewPackage(id, functions, constants)
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
