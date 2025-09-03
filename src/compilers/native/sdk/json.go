package sdk

import (
	"bytes"
	"encoding/json"

	objects2 "github.com/markel1974/c64emu/src/vm/objects"
)

func init() {
	RegisterPackage(NewJson)
}

// Json represents a module containing JSON-related operations and utilities.
type Json struct {
	container map[string]objects2.IObject
}

// NewJson creates and returns a new instance of Json containing predefined JSON operation modules.
func NewJson(factory objects2.IGateKeeper) IPackage {
	j := &Json{}
	container := []objects2.IObject{
		factory.NewFuncImport(objects2.FrameStatic, "Unmarshal", j.unmarshal),
		factory.NewFuncImport(objects2.FrameStatic, "Marshal", j.marshal),
		factory.NewFuncImport(objects2.FrameStatic, "Indent", j.indent),
		factory.NewFuncImport(objects2.FrameStatic, "HTMLEscape", j.htmlEscape),
	}
	j.container = BuildContainer(container, nil)
	return j
}

// Name returns the string identifier "json" for the Json module.
func (j *Json) Name() string {
	return "json"
}

// Get retrieves an object associated with the given name from the container. It returns the object and a boolean indicating success.
func (j *Json) Get(name string) (objects2.IObject, bool) {
	v, ok := j.container[name]
	return v, ok
}

// Unmarshal parses a JSON-encoded string or byte slice into a Map object and returns it as IObject.
func (j *Json) unmarshal(gk objects2.IGateKeeper, frame int, args ...objects2.IObject) (objects2.IObject, error) {
	if len(args) != 1 {
		return nil, objects2.ErrInvalidArgumentsNumber
	}
	var data []byte
	switch o := args[0].(type) {
	case *objects2.Bytes:
		data = o.Value()
	case *objects2.String:
		data = []byte(o.Value())
	}
	if data == nil {
		return nil, objects2.NewInvalidArgumentError(0, "bytes/string", args[0].TypeName())
	}
	d := make(map[string]interface{})
	if err := json.Unmarshal(data, &d); err != nil {
		return gk.NewError(frame, err.Error()), nil
	}
	result := gk.NewMap(frame, gk.FromMap(frame, d))
	return result, nil
}

// Marshal serializes a single IObject into a JSON-encoded byte slice and returns it as a Bytes object.
func (j *Json) marshal(gk objects2.IGateKeeper, frame int, args ...objects2.IObject) (objects2.IObject, error) {
	if len(args) != 1 {
		return nil, objects2.ErrInvalidArgumentsNumber
	}
	result, err := json.Marshal(gk.ToInterface(args[0]))
	if err != nil {
		return gk.NewError(frame, err.Error()), nil
	}
	return gk.NewBytes(frame, result), nil
}

// Indent takes a JSON object (bytes or string), a prefix, and an indent string, and returns the indented JSON.
func (j *Json) indent(gk objects2.IGateKeeper, frame int, args ...objects2.IObject) (objects2.IObject, error) {
	if len(args) != 3 {
		return nil, objects2.ErrInvalidArgumentsNumber
	}
	prefix, err := gk.ToStringArg(1, args[1])
	if err != nil {
		return nil, err
	}
	indent, err := gk.ToStringArg(2, args[2])
	if err != nil {
		return nil, err
	}
	switch o := args[0].(type) {
	case *objects2.Bytes:
		var dst bytes.Buffer
		err = json.Indent(&dst, o.Value(), prefix, indent)
		if err != nil {
			return gk.NewError(frame, err.Error()), nil
		}
		return gk.NewBytes(frame, dst.Bytes()), nil
	case *objects2.String:
		var dst bytes.Buffer
		err = json.Indent(&dst, []byte(o.Value()), prefix, indent)
		if err != nil {
			return gk.NewError(frame, err.Error()), nil
		}
		return gk.NewBytes(frame, dst.Bytes()), nil
	default:
		return nil, objects2.NewInvalidArgumentError(0, "bytes/string", args[0].TypeName())
	}
}

// HTMLEscape escapes certain characters in a JSON string or byte slice to their HTML-safe equivalents.
// Accepts one argument of type `*objects.Bytes` or `*objects.String`.
// Returns a new `*objects.Bytes` containing the escaped output or an error if an invalid argument type is provided.
// Errors if the number of arguments is not exactly one.
func (j *Json) htmlEscape(gk objects2.IGateKeeper, frame int, args ...objects2.IObject) (objects2.IObject, error) {
	if len(args) != 1 {
		return nil, objects2.ErrInvalidArgumentsNumber
	}
	switch o := args[0].(type) {
	case *objects2.Bytes:
		var dst bytes.Buffer
		json.HTMLEscape(&dst, o.Value())
		return gk.NewBytes(frame, dst.Bytes()), nil
	case *objects2.String:
		var dst bytes.Buffer
		json.HTMLEscape(&dst, []byte(o.Value()))
		return gk.NewBytes(frame, dst.Bytes()), nil
	default:
		return nil, objects2.NewInvalidArgumentError(0, "bytes/string", args[0].TypeName())
	}
}
