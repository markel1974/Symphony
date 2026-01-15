package compilers

import (
	"github.com/markel1974/symphony/src/vm/bytecode"
	"github.com/markel1974/symphony/src/vm/objects"
	"github.com/markel1974/symphony/src/vm/opcodes"

	_nativeCompiler "github.com/markel1974/symphony/src/compilers/native/compiler"
	_nativeLoader "github.com/markel1974/symphony/src/compilers/native/sdk"
)

// NewCompiler creates and initializes a new compiler instance, loading native packages and configuring compilation tools.
// It accepts an IGateKeeper, IOpcodes, and ILoader to manage object allocation, opcodes, and runtime/package loading.
// Returns a compiled ICompiler instance ready for use or an error if initialization fails.
func NewCompiler(gk objects.IGateKeeper, opcodes opcodes.IOpcodes, loader bytecode.ILoader) (bytecode.ICompiler, error) {
	switch opcodes.Id() {
	default:
		if err := loader.AddPackages(_nativeLoader.Packages(gk)); err != nil {
			return nil, err
		}
		compiler := _nativeCompiler.New(gk, loader, opcodes)
		return compiler, nil
	}
}
