package stdlib

import (
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// wrapError wraps a Go error into an IObject-compatible error object or returns TrueValue if the error is nil.
func wrapError(err error) objects.IObject {
	if err == nil {
		return objects.TrueValue
	}
	return objects.NewError(objects.NewStringNoSize(err.Error()))
}
