package sdk

import (
	"github.com/markel1974/c64emu/src/vm/objects"
)

func init() {
	RegisterPackage(NewErrors)
}

// Errors is a type that encapsulates a map of module functions accessible as objects.
type Errors struct {
	container map[string]objects.IObject
}

// NewErrors initializes and returns a new Errors instance with pre-defined function modules.
func NewErrors(factory objects.IGateKeeper) IPackage {
	e := &Errors{}
	container := []objects.IObject{
		factory.NewFuncImport(objects.FrameStatic, "New", e.New),
	}
	e.container = BuildContainer(container, nil)
	return e
}

// New creates a new error object from the provided argument, ensuring it is a valid string and returning an error if not.
func (e *Errors) New(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) != 1 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	s, err := gk.ToStringArg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewError(frame, s), nil
}

// Name returns the name identifier of the Errors type.
func (e *Errors) Name() string {
	return "errors"
}

// Get retrieves an object associated with the given name from the container. It returns the object and a boolean indicating success.
func (e *Errors) Get(name string) (objects.IObject, bool) {
	v, ok := e.container[name]
	return v, ok
}
