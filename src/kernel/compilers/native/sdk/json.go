package sdk

import (
	"bytes"
	"encoding/json"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func init() {
	RegisterPackage(NewJson)
}

// Json represents a module containing JSON-related operations and utilities.
type Json struct {
	factory   objects.IGateKeeper
	container map[string]objects.IObject
}

// NewJson creates and returns a new instance of Json containing predefined JSON operation modules.
func NewJson(factory objects.IGateKeeper) IPackage {
	j := &Json{
		factory: factory,
	}
	container := []objects.IObject{
		factory.NewFuncPackage(objects.FuncPackageDef, "Unmarshal", j.unmarshal),
		factory.NewFuncPackage(objects.FuncPackageDef, "Marshal", j.marshal),
		factory.NewFuncPackage(objects.FuncPackageDef, "Indent", j.indent),
		factory.NewFuncPackage(objects.FuncPackageDef, "HTMLEscape", j.htmlEscape),
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
func (j *Json) unmarshal(frame int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	var data []byte
	switch o := args[0].(type) {
	case *objects.Bytes:
		data = o.Value()
	case *objects.String:
		data = []byte(o.Value())
	}
	if data == nil {
		return nil, objects.NewInvalidArgumentError(0, "bytes/string", args[0].TypeName())
	}
	d := make(map[string]interface{})
	if err = json.Unmarshal(data, &d); err != nil {
		return j.factory.NewError(frame, err.Error()), nil
	}
	result := j.factory.NewMap(frame, j.factory.FromMap(frame, d))
	return result, nil
}

// Marshal serializes a single IObject into a JSON-encoded byte slice and returns it as a Bytes object.
func (j *Json) marshal(frame int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	result, err := json.Marshal(j.factory.ToInterface(args[0]))
	if err != nil {
		return j.factory.NewError(frame, err.Error()), nil
	}
	return j.factory.NewBytes(frame, result), nil
}

// Indent takes a JSON object (bytes or string), a prefix, and an indent string, and returns the indented JSON.
func (j *Json) indent(frame int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 3 {
		return nil, objects.ErrWrongNumArguments
	}
	prefix, err := j.factory.ToStringArg(1, args[1])
	if err != nil {
		return nil, err
	}
	indent, err := j.factory.ToStringArg(2, args[2])
	if err != nil {
		return nil, err
	}
	switch o := args[0].(type) {
	case *objects.Bytes:
		var dst bytes.Buffer
		err = json.Indent(&dst, o.Value(), prefix, indent)
		if err != nil {
			return j.factory.NewError(frame, err.Error()), nil
		}
		return j.factory.NewBytes(frame, dst.Bytes()), nil
	case *objects.String:
		var dst bytes.Buffer
		err = json.Indent(&dst, []byte(o.Value()), prefix, indent)
		if err != nil {
			return j.factory.NewError(frame, err.Error()), nil
		}
		return j.factory.NewBytes(frame, dst.Bytes()), nil
	default:
		return nil, objects.NewInvalidArgumentError(0, "bytes/string", args[0].TypeName())
	}
}

// HTMLEscape escapes certain characters in a JSON string or byte slice to their HTML-safe equivalents.
// Accepts one argument of type `*objects.Bytes` or `*objects.String`.
// Returns a new `*objects.Bytes` containing the escaped output or an error if an invalid argument type is provided.
// Errors if the number of arguments is not exactly one.
func (j *Json) htmlEscape(frame int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	switch o := args[0].(type) {
	case *objects.Bytes:
		var dst bytes.Buffer
		json.HTMLEscape(&dst, o.Value())
		return j.factory.NewBytes(frame, dst.Bytes()), nil
	case *objects.String:
		var dst bytes.Buffer
		json.HTMLEscape(&dst, []byte(o.Value()))
		return j.factory.NewBytes(frame, dst.Bytes()), nil
	default:
		return nil, objects.NewInvalidArgumentError(0, "bytes/string", args[0].TypeName())
	}
}
