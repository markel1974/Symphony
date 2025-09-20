package sdk

import (
	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/objects"
)

func init() {
	RegisterPackage(NewErrors)
}

// Errors is a type that encapsulates a map of module functions accessible as objects.
type Errors struct {
	*bytecode.Package
}

// NewErrors initializes and returns a new Errors instance with pre-defined function modules.
func NewErrors(factory objects.IGateKeeper) bytecode.IPackage {
	const (
		defNew = "New"
	)
	e := &Errors{Package: bytecode.NewPackage("errors")}
	e.Add(defNew, factory.NewFuncImport(objects.FrameStatic, defNew, 1, e.New))
	return e
}

// New creates a new error object from the provided argument, ensuring it is a valid string and returning an error if not.
func (e *Errors) New(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	s, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewError(frame, s), nil
}
