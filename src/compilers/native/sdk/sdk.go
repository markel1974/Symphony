package sdk

import (
	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/objects"
)

// registerPackageFn is a function type defining a method that registers a package using IGateKeeper and returns an IPackage.
type registerPackageFn func(f objects.IGateKeeper) bytecode.IPackage

// _registerPackage is a slice of functions used to register packages, where each function takes an IGateKeeper and returns an IPackage.
var _registerPackage []registerPackageFn

// register appends a given package registration function to the internal list of package registries.
func register(f registerPackageFn) {
	_registerPackage = append(_registerPackage, f)
}

// Packages initializes and returns a slice of IPackage instances by invoking registered package functions with a provided IGateKeeper.
func Packages(gk objects.IGateKeeper) []bytecode.IPackage {
	ret := make([]bytecode.IPackage, len(_registerPackage))
	for idx, fn := range _registerPackage {
		ret[idx] = fn(gk)
	}
	return ret
}
