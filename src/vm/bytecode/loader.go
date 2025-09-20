package bytecode

import (
	"fmt"
	"strings"

	"github.com/markel1974/c64emu/src/vm/objects"
)

// Loader manages the registration and resolution of packages and their symbols within a runtime environment.
type Loader struct {
	gk       objects.IGateKeeper
	packages map[string]IPackage
}

// NewLoader creates and returns a new Loader instance initialized with the provided IGateKeeper.
func NewLoader(gk objects.IGateKeeper) *Loader {
	loader := &Loader{
		gk:       gk,
		packages: make(map[string]IPackage),
	}
	return loader
}

// RegisterPackage adds a list of packages to the loader using provided registration functions and returns an error if any fail.
func (l *Loader) RegisterPackage(registerPackage []RegisterPackageFn) error {
	for _, fn := range registerPackage {
		pkg := fn(l.gk)
		if err := l.AddPackage(pkg); err != nil {
			return err
		}
	}
	return nil
}

// AddPackage adds a new package to the loader's package map. Returns an error if the package name already exists.
func (l *Loader) AddPackage(p IPackage) error {
	if z, ok := l.packages[p.Name()]; ok {
		return fmt.Errorf("package %s => %s already exists", z.Name(), p.Name())
	}
	l.packages[p.Name()] = p
	return nil
}

// Resolve resolves a list of IObject symbols into their corresponding references from the registered packages.
// It returns a slice of resolved IObjects or an error if any symbol cannot be resolved.
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
