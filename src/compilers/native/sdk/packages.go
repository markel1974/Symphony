package sdk

import (
	"github.com/markel1974/c64emu/src/vm/bytecode"
)

// _registerPackage is a slice of functions used to register packages, where each function takes an IGateKeeper and returns an IPackage.
var _registerPackage []bytecode.RegisterPackageFn

// RegisterPackage appends a given package registration function to the internal list of package registries.
func RegisterPackage(f bytecode.RegisterPackageFn) {
	_registerPackage = append(_registerPackage, f)
}

// Packages returns the name of the package and a slice of RegisterPackageFn to register native functionality.
func Packages() []bytecode.RegisterPackageFn {
	return _registerPackage
}
