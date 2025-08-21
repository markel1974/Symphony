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
func NewBase64(f *objects.GateKeeper) *Base64 {
	b := &Base64{}
	container := []objects.IObject{
		f.NewFuncPackage(objects.FuncPackageDef, "EncodeToString", f.FuncIbSOs(base64.StdEncoding.EncodeToString)),
		f.NewFuncPackage(objects.FuncPackageDef, "EncodeToString", f.FuncIsObSe(base64.StdEncoding.DecodeString)),
		f.NewFuncPackage(objects.FuncPackageDef, "RawEncode", f.FuncIbSOs(base64.RawStdEncoding.EncodeToString)),
		f.NewFuncPackage(objects.FuncPackageDef, "RawDecode", f.FuncIsObSe(base64.RawStdEncoding.DecodeString)),
		f.NewFuncPackage(objects.FuncPackageDef, "UrlEncode", f.FuncIbSOs(base64.URLEncoding.EncodeToString)),
		f.NewFuncPackage(objects.FuncPackageDef, "UrlDecode", f.FuncIsObSe(base64.URLEncoding.DecodeString)),
		f.NewFuncPackage(objects.FuncPackageDef, "RawUrlEncode", f.FuncIbSOs(base64.RawURLEncoding.EncodeToString)),
		f.NewFuncPackage(objects.FuncPackageDef, "rawUrlDecode", f.FuncIsObSe(base64.RawURLEncoding.DecodeString)),
	}
	b.Package = NewPackage("base64", container, nil)
	return b
}
