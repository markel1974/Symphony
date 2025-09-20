package compilers

import (
	_nativeLoader "github.com/markel1974/c64emu/src/compilers/native/sdk"
	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"

	_nativeCompiler "github.com/markel1974/c64emu/src/compilers/native/compiler"
)

// NewCompiler creates a new compiler and loader based on the provided IGateKeeper, Opcodes, and identifier string.
// It returns an ICompiler, an ILoader, and an error if the creation fails.
func NewCompiler(gk objects.IGateKeeper, opcodes opcodes.IOpcodes) (bytecode.ICompiler, bytecode.ILoader, error) {
	switch opcodes.Id() {
	default:
		loader := bytecode.NewLoader(gk)
		if err := loader.RegisterPackage(_nativeLoader.Packages()); err != nil {
			return nil, nil, err
		}
		compiler := _nativeCompiler.New(gk, loader, opcodes)
		return compiler, loader, nil
	}
}
