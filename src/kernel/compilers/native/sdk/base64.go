package sdk

import (
	"encoding/base64"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func init() {
	RegisterPackage(NewBase64)
}

// Base64 represents a type that provides a module map for Base64-related encoding and decoding operations.
type Base64 struct {
	container map[string]objects.IObject
}

// NewBase64 initializes a new Base64 instance with predefined encoding and decoding functions in the module map.
func NewBase64(f objects.IGateKeeper) IPackage {
	container := []objects.IObject{
		f.NewFuncImport(objects.FrameStatic, "EncodeToString", f.FuncIbSOs(base64.StdEncoding.EncodeToString)),
		f.NewFuncImport(objects.FrameStatic, "EncodeToString", f.FuncIsObSe(base64.StdEncoding.DecodeString)),
		f.NewFuncImport(objects.FrameStatic, "RawEncode", f.FuncIbSOs(base64.RawStdEncoding.EncodeToString)),
		f.NewFuncImport(objects.FrameStatic, "RawDecode", f.FuncIsObSe(base64.RawStdEncoding.DecodeString)),
		f.NewFuncImport(objects.FrameStatic, "UrlEncode", f.FuncIbSOs(base64.URLEncoding.EncodeToString)),
		f.NewFuncImport(objects.FrameStatic, "UrlDecode", f.FuncIsObSe(base64.URLEncoding.DecodeString)),
		f.NewFuncImport(objects.FrameStatic, "RawUrlEncode", f.FuncIbSOs(base64.RawURLEncoding.EncodeToString)),
		f.NewFuncImport(objects.FrameStatic, "rawUrlDecode", f.FuncIsObSe(base64.RawURLEncoding.DecodeString)),
	}
	b := &Base64{
		container: BuildContainer(container, nil),
	}
	return b
}

// Name returns the name identifier of the Base64 type, which is "base64".
func (b *Base64) Name() string {
	return "base64"
}

// Get retrieves an object associated with the given name from the container. It returns the object and a boolean indicating success.
func (b *Base64) Get(name string) (objects.IObject, bool) {
	v, ok := b.container[name]
	return v, ok
}
