package sdk

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// init initializes the package by registering the Os package using the RegisterPackage function.
func init() {
	RegisterPackage(NewOs)
}

// Os represents a package that provides access to a container of IObject instances, allowing retrieval by name.
type Os struct {
	container map[string]objects.IObject
}

// NewOs initializes a new instance of the Os package with provided gatekeeper functionalities and returns it as IPackage.
func NewOs(factory objects.IGateKeeper) IPackage {
	f := &Os{}
	container := []objects.IObject{
		factory.NewFuncImport(objects.FrameStatic, "Exit", f.exit),
	}
	f.container = BuildContainer(container, nil)
	return f
}

// Name returns the name of the operating system as a string.
func (f *Os) Name() string {
	return "os"
}

// Get retrieves the IObject associated with the specified name from the container and returns it along with a boolean indicating its presence.
func (f *Os) Get(name string) (objects.IObject, bool) {
	v, ok := f.container[name]
	return v, ok
}

// exit terminates the current process after logging the provided arguments, converting them to native types with IGateKeeper.
func (f *Os) exit(gk objects.IGateKeeper, _ int, args ...objects.IObject) (objects.IObject, error) {
	ret := int64(0)
	if len(args) > 0 {
		ret = args[0].AsInt64()
	}
	return nil, fmt.Errorf("exit with code %d", ret)
}
