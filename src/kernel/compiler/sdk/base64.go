package sdk

import (
	"encoding/base64"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Base64 represents a type that provides a module map for Base64-related encoding and decoding operations.
type Base64 struct {
	*Package
}

// NewBase64 initializes a new Base64 instance with predefined encoding and decoding functions in the module map.
func NewBase64() *Base64 {
	b := &Base64{}
	container := []*objects.FuncPackage{
		objects.NewFuncPackage(objects.FuncPackageDef, "EncodeToString", objects.FuncIbSOs(base64.StdEncoding.EncodeToString)),
		objects.NewFuncPackage(objects.FuncPackageDef, "EncodeToString", objects.FuncIsObSe(base64.StdEncoding.DecodeString)),
		objects.NewFuncPackage(objects.FuncPackageDef, "RawEncode", objects.FuncIbSOs(base64.RawStdEncoding.EncodeToString)),
		objects.NewFuncPackage(objects.FuncPackageDef, "RawDecode", objects.FuncIsObSe(base64.RawStdEncoding.DecodeString)),
		objects.NewFuncPackage(objects.FuncPackageDef, "UrlEncode", objects.FuncIbSOs(base64.URLEncoding.EncodeToString)),
		objects.NewFuncPackage(objects.FuncPackageDef, "UrlDecode", objects.FuncIsObSe(base64.URLEncoding.DecodeString)),
		objects.NewFuncPackage(objects.FuncPackageDef, "RawUrlEncode", objects.FuncIbSOs(base64.RawURLEncoding.EncodeToString)),
		objects.NewFuncPackage(objects.FuncPackageDef, "rawUrlDecode", objects.FuncIsObSe(base64.RawURLEncoding.DecodeString)),
	}
	b.Package = NewPackage("base64", container, nil)
	return b
}
