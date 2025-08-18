package sdk

import (
	"encoding/hex"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

var _hexModule = map[string]objects.IObject{
	"EncodeToString": objects.NewFunctionModule(objects.FunctionModuleDef, "EncodeToString", objects.FuncIbSOs(hex.EncodeToString)),
	"DecodeString":   objects.NewFunctionModule(objects.FunctionModuleDef, "DecodeString", objects.FuncIsObSe(hex.DecodeString)),
}
