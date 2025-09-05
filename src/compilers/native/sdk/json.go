package sdk

import (
	"bytes"
	"encoding/json"

	"github.com/markel1974/c64emu/src/vm/objects"
)

func init() {
	RegisterPackage(NewJson)
}

// Json represents a module containing JSON-related operations and utilities.
type Json struct {
	container map[string]objects.IObject
}

// NewJson creates and returns a new instance of Json containing predefined JSON operation modules.
func NewJson(factory objects.IGateKeeper) IPackage {
	j := &Json{}
	container := []objects.IObject{
		factory.NewFuncImport(objects.FrameStatic, "Unmarshal", j.unmarshal),
		factory.NewFuncImport(objects.FrameStatic, "Marshal", j.marshal),
		factory.NewFuncImport(objects.FrameStatic, "Indent", j.indent),
		factory.NewFuncImport(objects.FrameStatic, "HTMLEscape", j.htmlEscape),
	}
	j.container = BuildContainer(container, nil)
	return j
}

// Name returns the string identifier "json" for the Json module.
func (j *Json) Name() string {
	return "json"
}

// Get retrieves an object associated with the given name from the container. It returns the object and a boolean indicating success.
func (j *Json) Get(name string) (objects.IObject, bool) {
	v, ok := j.container[name]
	return v, ok
}

// Unmarshal parses a JSON-encoded string or byte slice into a Map object and returns it as IObject.
func (j *Json) unmarshal(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) != 1 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	var data []byte
	switch o := args[0].(type) {
	case *objects.Bytes:
		data = o.Value()
	case *objects.String:
		data = []byte(o.Value())
	}
	if data == nil {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	d := make(map[string]interface{})
	if err := json.Unmarshal(data, &d); err != nil {
		return 0, gk.NewError(frame, err.Error()), nil
	}
	result := gk.NewMap(frame, gk.FromMap(frame, d))
	return 1, result, nil
}

// Marshal serializes a single IObject into a JSON-encoded byte slice and returns it as a Bytes object.
func (j *Json) marshal(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) != 1 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	result, err := json.Marshal(gk.ToInterface(args[0]))
	if err != nil {
		return 0, gk.NewError(frame, err.Error()), nil
	}
	return 1, gk.NewBytes(frame, result), nil
}

// Indent takes a JSON object (bytes or string), a prefix, and an indent string, and returns the indented JSON.
func (j *Json) indent(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	prefix, err := gk.ToStringArg(1, args)
	if err != nil {
		return 0, nil, err
	}
	indent, err := gk.ToStringArg(2, args)
	if err != nil {
		return 0, nil, err
	}
	switch o := args[0].(type) {
	case *objects.Bytes:
		var dst bytes.Buffer
		if err = json.Indent(&dst, o.Value(), prefix, indent); err != nil {
			return 0, gk.NewError(frame, err.Error()), nil
		}
		return 1, gk.NewBytes(frame, dst.Bytes()), nil
	case *objects.String:
		var dst bytes.Buffer
		if err = json.Indent(&dst, []byte(o.Value()), prefix, indent); err != nil {
			return 0, gk.NewError(frame, err.Error()), nil
		}
		return 1, gk.NewBytes(frame, dst.Bytes()), nil
	default:
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
}

// HTMLEscape escapes certain characters in a JSON string or byte slice to their HTML-safe equivalents.
// Accepts one argument of type `*objects.Bytes` or `*objects.String`.
// Returns a new `*objects.Bytes` containing the escaped output or an error if an invalid argument type is provided.
// Errors if the number of arguments is not exactly one.
func (j *Json) htmlEscape(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) != 1 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	switch o := args[0].(type) {
	case *objects.Bytes:
		var dst bytes.Buffer
		json.HTMLEscape(&dst, o.Value())
		return 1, gk.NewBytes(frame, dst.Bytes()), nil
	case *objects.String:
		var dst bytes.Buffer
		json.HTMLEscape(&dst, []byte(o.Value()))
		return 1, gk.NewBytes(frame, dst.Bytes()), nil
	default:
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
}
