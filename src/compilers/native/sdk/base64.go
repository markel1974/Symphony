package sdk

import (
	"encoding/base64"

	"github.com/markel1974/c64emu/src/vm/objects"
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
	b := &Base64{}
	container := []objects.IObject{
		f.NewFuncImport(objects.FrameStatic, "EncodeToString", b.funcBytesToString(base64.StdEncoding.EncodeToString)),
		f.NewFuncImport(objects.FrameStatic, "EncodeToString", b.funcStringToBytesError(base64.StdEncoding.DecodeString)),
		f.NewFuncImport(objects.FrameStatic, "RawEncode", b.funcBytesToString(base64.RawStdEncoding.EncodeToString)),
		f.NewFuncImport(objects.FrameStatic, "RawDecode", b.funcStringToBytesError(base64.RawStdEncoding.DecodeString)),
		f.NewFuncImport(objects.FrameStatic, "UrlEncode", b.funcBytesToString(base64.URLEncoding.EncodeToString)),
		f.NewFuncImport(objects.FrameStatic, "UrlDecode", b.funcStringToBytesError(base64.URLEncoding.DecodeString)),
		f.NewFuncImport(objects.FrameStatic, "RawUrlEncode", b.funcBytesToString(base64.RawURLEncoding.EncodeToString)),
		f.NewFuncImport(objects.FrameStatic, "rawUrlDecode", b.funcStringToBytesError(base64.RawURLEncoding.DecodeString)),
	}
	b.container = BuildContainer(container, nil)
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

// funcBytesToString wraps a given function to process a byte slice argument and returns a callable function for the system.
func (b *Base64) funcBytesToString(fn func([]byte) string) objects.FuncCallable {
	return func(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
		if len(args) != 1 {
			return 0, nil, objects.ErrInvalidArgumentsNumber
		}
		bs1, err := gk.ToBytesArg(0, args[0])
		if err != nil {
			return 0, nil, err
		}
		v := gk.NewString(frame, fn(bs1))
		return 1, v, nil
	}
}

// funcStringToBytesError wraps a function accepting a string and returning []byte and error to conform to the FuncCallable type.
// It validates the number of arguments, ensures the first argument is a string, and transforms outputs to IObject types.
func (b *Base64) funcStringToBytesError(fn func(string) ([]byte, error)) objects.FuncCallable {
	return func(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
		if len(args) != 1 {
			return 0, nil, objects.ErrInvalidArgumentsNumber
		}
		s1, err := gk.ToStringArg(0, args[0])
		if err != nil {
			return 0, nil, err
		}
		res, err := fn(s1)
		if err != nil {
			return 0, gk.NewError(frame, err.Error()), nil
		}
		return 1, gk.NewBytes(frame, res), nil
	}
}
