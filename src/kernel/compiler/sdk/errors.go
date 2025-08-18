package sdk

import (
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Errors is a type that encapsulates a map of module functions accessible as objects.
type Errors struct {
	*Package
}

// NewErrors initializes and returns a new Errors instance with pre-defined function modules.
func NewErrors() *Errors {
	e := &Errors{}
	container := []*objects.FuncPackage{
		objects.NewFuncPackage(objects.FuncPackageDef, "New", e.New),
	}
	e.Package = NewPackage("errors", container, nil)
	return e
}

// New creates a new error object from the provided argument, ensuring it is a valid string and returning an error if not.
func (e *Errors) New(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	s, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	v, err := objects.NewString(s)
	if err != nil {
		return nil, err
	}
	return objects.NewError(v), nil
}
