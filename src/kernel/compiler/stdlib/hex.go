package stdlib

import (
	"encoding/hex"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

var _hexModule = map[string]objects.IObject{
	"encode": objects.NewFunctionUser("encode", objects.FuncAYRS(hex.EncodeToString)),
	"decode": objects.NewFunctionUser("decode", objects.FuncASRYE(hex.DecodeString)),
}
