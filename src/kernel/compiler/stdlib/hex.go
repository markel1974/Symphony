package stdlib

import (
	"encoding/hex"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

var _hexModule = map[string]objects.IObject{
	"encode": objects.NewFunctionModule(objects.FunctionModuleDef, "encode", objects.FuncIbSOs(hex.EncodeToString)),
	"decode": objects.NewFunctionModule(objects.FunctionModuleDef, "decode", objects.FuncIsObSe(hex.DecodeString)),
}
