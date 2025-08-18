package stdlib

import (
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

var _errorsModule = map[string]objects.IObject{
	"New": objects.NewFunctionModule(objects.FunctionModuleDef, "New", errorsNew),
}

func errorsNew(args ...objects.IObject) (ret objects.IObject, err error) {
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
