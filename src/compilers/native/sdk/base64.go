package sdk

import (
	"encoding/base64"

	"github.com/markel1974/symphony/src/vm/bytecode"
	"github.com/markel1974/symphony/src/vm/objects"
)

// init initializes the package by registering the Base64 package using the register function.
func init() {
	register(NewBase64)
}

// Base64 provides methods for working with Base64 encoding and decoding operations within a container of IObject mappings.
type Base64 struct {
	*bytecode.Package
}

// NewBase64 initializes and returns a new Base64 package implementing IPackage with predefined encoding/decoding functions.
func NewBase64(f objects.IGateKeeper) bytecode.IPackage {
	const (
		defEncodeToString = "EncodeToString"
		defDecodeToString = "DecodeToString"
		defRawEncode      = "RawEncode"
		defRawDecode      = "RawDecode"
		defUrlEncode      = "UrlEncode"
		defUrlDecode      = "UrlDecode"
		defRawUrlEncode   = "RawUrlEncode"
		defRawUrlDecode   = "RawUrlDecode"
	)
	b := &Base64{Package: bytecode.NewPackage("base64")}
	b.Add(defEncodeToString, f.NewFuncImport(objects.FrameStatic, defEncodeToString, 1, b.stdEncodeToString))
	b.Add(defDecodeToString, f.NewFuncImport(objects.FrameStatic, defDecodeToString, 1, b.stdDecodeString))
	b.Add(defRawEncode, f.NewFuncImport(objects.FrameStatic, defRawEncode, 1, b.rawEncodeToString))
	b.Add(defRawDecode, f.NewFuncImport(objects.FrameStatic, defRawDecode, 1, b.rawDecodeString))
	b.Add(defUrlEncode, f.NewFuncImport(objects.FrameStatic, defUrlEncode, 1, b.urlEncodeToString))
	b.Add(defUrlDecode, f.NewFuncImport(objects.FrameStatic, defUrlDecode, 1, b.urlDecodeString))
	b.Add(defRawUrlEncode, f.NewFuncImport(objects.FrameStatic, defRawUrlEncode, 1, b.rawUrlEncodeToString))
	b.Add(defRawUrlDecode, f.NewFuncImport(objects.FrameStatic, defRawUrlDecode, 1, b.rawUrlDecodeString))
	return b
}

// stdEncodeToString encodes the first argument to a Base64 string using standard encoding and returns the result as a new string.
func (b *Base64) stdEncodeToString(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	bs1, err := gk.ToBytesArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	v := gk.NewString(frame, base64.StdEncoding.EncodeToString(bs1))
	return 1, v, nil
}

// stdDecodeString decodes a Base64-encoded string into its original byte representation using standard Base64 encoding.
func (b *Base64) stdDecodeString(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	s1, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	res, err := base64.StdEncoding.DecodeString(s1)
	if err != nil {
		return 0, gk.NewError(frame, err.Error()), nil
	}
	return 1, gk.NewBytes(frame, res), nil
}

// rawEncodeToString encodes the first argument as a Base64 raw encoded string and returns the result as an IObject.
// It utilizes base64.RawStdEncoding without padding and supports multiple arguments for error handling.
func (b *Base64) rawEncodeToString(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	bs1, err := gk.ToBytesArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	v := gk.NewString(frame, base64.RawStdEncoding.EncodeToString(bs1))
	return 1, v, nil
}

// rawDecodeString decodes a raw Base64-encoded string without padding and returns the resulting bytes as an IObject.
func (b *Base64) rawDecodeString(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	s1, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	res, err := base64.RawStdEncoding.DecodeString(s1)
	if err != nil {
		return 0, gk.NewError(frame, err.Error()), nil
	}
	return 1, gk.NewBytes(frame, res), nil
}

// urlEncodeToString encodes the first argument to a URL-safe Base64 string and returns it as an IObject.
func (b *Base64) urlEncodeToString(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	bs1, err := gk.ToBytesArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	v := gk.NewString(frame, base64.URLEncoding.EncodeToString(bs1))
	return 1, v, nil
}

// urlDecodeString decodes a URL-safe Base64 encoded string and returns the decoded byte array as an IObject.
func (b *Base64) urlDecodeString(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	s1, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	res, err := base64.URLEncoding.DecodeString(s1)
	if err != nil {
		return 0, gk.NewError(frame, err.Error()), nil
	}
	return 1, gk.NewBytes(frame, res), nil
}

// rawUrlDecodeString decodes a raw URL-safe base64 encoded string and returns the decoded bytes as an IObject.
func (b *Base64) rawUrlDecodeString(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	s1, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	res, err := base64.RawURLEncoding.DecodeString(s1)
	if err != nil {
		return 0, gk.NewError(frame, err.Error()), nil
	}
	return 1, gk.NewBytes(frame, res), nil
}

// rawUrlEncodeToString encodes the input bytes into a Base64 raw URL encoded string and returns it as an IObject.
func (b *Base64) rawUrlEncodeToString(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	bs1, err := gk.ToBytesArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	v := gk.NewString(frame, base64.RawURLEncoding.EncodeToString(bs1))
	return 1, v, nil
}
