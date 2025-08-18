package sdk

import (
	"encoding/base64"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Base64 represents a type that provides a module map for Base64-related encoding and decoding operations.
type Base64 struct {
	module map[string]objects.IObject
}

// NewBase64 initializes a new Base64 instance with predefined encoding and decoding functions in the module map.
func NewBase64() *Base64 {
	b := &Base64{}
	b.module = map[string]objects.IObject{
		"EncodeToString":       objects.NewFunctionModule(objects.FunctionModuleDef, "EncodeToString", objects.FuncIbSOs(base64.StdEncoding.EncodeToString)),
		"DecodeString":         objects.NewFunctionModule(objects.FunctionModuleDef, "EncodeToString", objects.FuncIsObSe(base64.StdEncoding.DecodeString)),
		"RawEncodeToString":    objects.NewFunctionModule(objects.FunctionModuleDef, "RawEncode", objects.FuncIbSOs(base64.RawStdEncoding.EncodeToString)),
		"RawDecodeString":      objects.NewFunctionModule(objects.FunctionModuleDef, "RawDecode", objects.FuncIsObSe(base64.RawStdEncoding.DecodeString)),
		"UrlEncodeToString":    objects.NewFunctionModule(objects.FunctionModuleDef, "UrlEncode", objects.FuncIbSOs(base64.URLEncoding.EncodeToString)),
		"UrlDecodeString":      objects.NewFunctionModule(objects.FunctionModuleDef, "UrlDecode", objects.FuncIsObSe(base64.URLEncoding.DecodeString)),
		"RawUrlEncodeToString": objects.NewFunctionModule(objects.FunctionModuleDef, "RawUrlEncode", objects.FuncIbSOs(base64.RawURLEncoding.EncodeToString)),
		"rawUrlDecodeStringe":  objects.NewFunctionModule(objects.FunctionModuleDef, "rawUrlDecode", objects.FuncIsObSe(base64.RawURLEncoding.DecodeString)),
	}
	return b
}

// Name returns the name of the Base64 module.
func (b *Base64) Name() string {
	return "base64"
}

// Module returns the map of string keys to objects.IObject representing the Base64 module implementation.
func (b *Base64) Module() map[string]objects.IObject {
	return b.module
}
