package stdlib

import (
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func wrapError(err error) objects.IObject {
	if err == nil {
		return objects.TrueValue
	}
	return &objects.Error{Value: &objects.String{Value: err.Error()}}
}
