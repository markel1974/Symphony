package sdk

import (
	"encoding/hex"

	"github.com/markel1974/c64emu/src/vm/objects"
)

// init registers the "Hex" package by appending its registration function to the internal package list.
func init() {
	RegisterPackage(NewHex)
}

// Hex represents a type that provides functionality for encoding and decoding hexadecimal strings.
// It contains a map-based container for managing associated objects and functions.
type Hex struct {
	container map[string]objects.IObject
}

// NewHex initializes a new Hex instance and returns it as an IPackage. It registers encodeToString and decodeString functions.
func NewHex(gk objects.IGateKeeper) IPackage {
	h := &Hex{}
	container := []objects.IObject{
		gk.NewFuncImport(objects.FrameStatic, "EncodeToString", 1, h.encodeToString),
		gk.NewFuncImport(objects.FrameStatic, "DecodeString", 1, h.decodeString),
	}
	h.container = BuildContainer(container, nil)
	return h
}

// Name returns the static name "hex" for the Hex type.
func (f *Hex) Name() string {
	return "hex"
}

// Get retrieves an object by its name from the container and indicates whether it exists.
func (f *Hex) Get(name string) (objects.IObject, bool) {
	v, ok := f.container[name]
	return v, ok
}

// encodeToString converts the given IObject into a hexadecimal string representation using the provided IGateKeeper.
// It requires exactly one argument. Returns an error if argument count is invalid or conversion fails.
func (f *Hex) encodeToString(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	bs1, err := gk.ToBytesArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	v := gk.NewString(frame, hex.EncodeToString(bs1))
	return 1, v, nil
}

// decodeString decodes a hexadecimal string to its byte representation, returning the result as an IObject or an error.
func (f *Hex) decodeString(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	s1, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	res, err := hex.DecodeString(s1)
	if err != nil {
		return 0, gk.NewError(frame, err.Error()), nil
	}
	return 1, gk.NewBytes(frame, res), nil
}
