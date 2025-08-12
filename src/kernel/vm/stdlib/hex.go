package stdlib

import (
	"encoding/hex"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

var hexModule = map[string]objects.IObject{
	"encode": objects.NewUserFunction("encode", objects.FuncAYRS(hex.EncodeToString)),
	"decode": objects.NewUserFunction("decode", objects.FuncASRYE(hex.DecodeString)),
}
