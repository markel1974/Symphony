package compilers

import (
	_nativeLoader "github.com/markel1974/c64emu/src/kernel/compilers/native/sdk"
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"

	_nativeCompiler "github.com/markel1974/c64emu/src/kernel/compilers/native/compiler"
)

// NewCompiler creates a new compiler and loader based on the provided IGateKeeper, Opcodes, and identifier string.
// It returns an ICompiler, an ILoader, and an error if the creation fails.
func NewCompiler(gk objects.IGateKeeper, opcodes *bytecode.Opcodes, id string) (bytecode.ICompiler, bytecode.ILoader, error) {
	switch id {
	default:
		loader, err := NewLoader(gk, id)
		if err != nil {
			return nil, nil, err
		}
		return _nativeCompiler.New(gk, loader, opcodes), loader, nil
	}
}

// NewLoader initializes and returns a new ILoader instance based on the provided IGateKeeper and identifier string.
// Returns an ILoader interface or an error if initialization fails.
func NewLoader(gk objects.IGateKeeper, id string) (bytecode.ILoader, error) {
	switch id {
	default:
		loader := _nativeLoader.NewLoader(gk)
		return loader, nil
	}
}
