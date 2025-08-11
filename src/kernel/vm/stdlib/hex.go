package stdlib

import (
	"encoding/hex"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

var hexModule = map[string]objects.IObject{
	"encode": &objects.UserFunction{Value: FuncAYRS(hex.EncodeToString)},
	"decode": &objects.UserFunction{Value: FuncASRYE(hex.DecodeString)},
}
