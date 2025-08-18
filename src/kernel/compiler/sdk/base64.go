package sdk

import (
	"encoding/base64"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// base64Module provides encoding and decoding utility functions for base64, raw, and URL-safe base64 formats.
var _base64Module = map[string]objects.IObject{
	"EncodeToString":       objects.NewFunctionModule(objects.FunctionModuleDef, "EncodeToString", objects.FuncIbSOs(base64.StdEncoding.EncodeToString)),
	"DecodeString":         objects.NewFunctionModule(objects.FunctionModuleDef, "EncodeToString", objects.FuncIsObSe(base64.StdEncoding.DecodeString)),
	"RawEncodeToString":    objects.NewFunctionModule(objects.FunctionModuleDef, "RawEncode", objects.FuncIbSOs(base64.RawStdEncoding.EncodeToString)),
	"RawDecodeString":      objects.NewFunctionModule(objects.FunctionModuleDef, "RawDecode", objects.FuncIsObSe(base64.RawStdEncoding.DecodeString)),
	"UrlEncodeToString":    objects.NewFunctionModule(objects.FunctionModuleDef, "UrlEncode", objects.FuncIbSOs(base64.URLEncoding.EncodeToString)),
	"UrlDecodeString":      objects.NewFunctionModule(objects.FunctionModuleDef, "UrlDecode", objects.FuncIsObSe(base64.URLEncoding.DecodeString)),
	"RawUrlEncodeToString": objects.NewFunctionModule(objects.FunctionModuleDef, "RawUrlEncode", objects.FuncIbSOs(base64.RawURLEncoding.EncodeToString)),
	"rawUrlDecodeStringe":  objects.NewFunctionModule(objects.FunctionModuleDef, "rawUrlDecode", objects.FuncIsObSe(base64.RawURLEncoding.DecodeString)),
}
