package stdlib

import (
	"encoding/base64"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// base64Module provides encoding and decoding utility functions for base64, raw, and URL-safe base64 formats.
var _base64Module = map[string]objects.IObject{
	"encode":         objects.NewFunctionModule(objects.FunctionModuleDef, "encode", objects.FuncIbSOs(base64.StdEncoding.EncodeToString)),
	"decode":         objects.NewFunctionModule(objects.FunctionModuleDef, "decode", objects.FuncIsObSe(base64.StdEncoding.DecodeString)),
	"raw_encode":     objects.NewFunctionModule(objects.FunctionModuleDef, "raw_encode", objects.FuncIbSOs(base64.RawStdEncoding.EncodeToString)),
	"raw_decode":     objects.NewFunctionModule(objects.FunctionModuleDef, "raw_decode", objects.FuncIsObSe(base64.RawStdEncoding.DecodeString)),
	"url_encode":     objects.NewFunctionModule(objects.FunctionModuleDef, "url_encode", objects.FuncIbSOs(base64.URLEncoding.EncodeToString)),
	"url_decode":     objects.NewFunctionModule(objects.FunctionModuleDef, "url_decode", objects.FuncIsObSe(base64.URLEncoding.DecodeString)),
	"raw_url_encode": objects.NewFunctionModule(objects.FunctionModuleDef, "raw_url_encode", objects.FuncIbSOs(base64.RawURLEncoding.EncodeToString)),
	"raw_url_decode": objects.NewFunctionModule(objects.FunctionModuleDef, "raw_url_decode", objects.FuncIsObSe(base64.RawURLEncoding.DecodeString)),
}
