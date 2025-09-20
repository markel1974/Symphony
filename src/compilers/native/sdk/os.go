package sdk

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/objects"
)

// init initializes the package by registering the Os package using the RegisterPackage function.
func init() {
	RegisterPackage(NewOs)
}

// Os represents a package that provides access to a container of IObject instances, allowing retrieval by name.
type Os struct {
	*bytecode.Package
}

// NewOs initializes a new instance of the Os package with provided gatekeeper functionalities and returns it as IPackage.
func NewOs(factory objects.IGateKeeper) bytecode.IPackage {
	const (
		defExit = "Exit"
	)
	f := &Os{Package: bytecode.NewPackage("os")}
	f.Add(defExit, factory.NewFuncImport(objects.FrameStatic, defExit, 1, f.exit))
	return f
}

// exit terminates the current process after logging the provided arguments, converting them to native types with IGateKeeper.
func (f *Os) exit(gk objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	ret, err := gk.ToInt64Arg(0, args)
	if err != nil {
		return 0, nil, err
	}
	return 0, nil, fmt.Errorf("exit with code %d", ret)
}
