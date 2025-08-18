package sdk

import (
	"bytes"
	"encoding/json"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Json represents a module containing JSON-related operations and utilities.
type Json struct {
	*Module
}

// NewJson creates and returns a new instance of Json containing predefined JSON operation modules.
func NewJson() *Json {
	j := &Json{
		Module: NewModule(),
	}
	j.attrs = map[string]objects.IObject{
		"Unmarshal":  objects.NewFunctionModule(objects.FunctionModuleDef, "Unmarshal", j.Unmarshal),
		"Marshal":    objects.NewFunctionModule(objects.FunctionModuleDef, "Marshal", j.Marshal),
		"Indent":     objects.NewFunctionModule(objects.FunctionModuleDef, "Indent", j.Indent),
		"HTMLEscape": objects.NewFunctionModule(objects.FunctionModuleDef, "html_escape", j.HTMLEscape),
	}
	return j
}

// Name returns the name of Json module.
func (j *Json) Name() string {
	return "json"
}

// Unmarshal parses a JSON-encoded string or byte slice into a Map object and returns it as IObject.
func (j *Json) Unmarshal(args ...objects.IObject) (ret objects.IObject, err error) {
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
		return objects.NewError(objects.NewStringNoSize(err.Error())), nil
	}
	result := objects.NewMap(objects.FromMap(d))
	return result, nil
}

// Marshal serializes a single IObject into a JSON-encoded byte slice and returns it as a Bytes object.
func (j *Json) Marshal(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	result, err := json.Marshal(objects.ToInterface(args[0]))
	if err != nil {
		return objects.NewError(objects.NewStringNoSize(err.Error())), nil
	}
	return objects.NewBytes(result), nil
}

// Indent takes a JSON object (bytes or string), a prefix, and an indent string, and returns the indented JSON.
func (j *Json) Indent(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 3 {
		return nil, objects.ErrWrongNumArguments
	}
	prefix, err := objects.ToStringArg(1, args[1])
	if err != nil {
		return nil, err
	}
	indent, err := objects.ToStringArg(2, args[2])
	if err != nil {
		return nil, err
	}
	switch o := args[0].(type) {
	case *objects.Bytes:
		var dst bytes.Buffer
		err = json.Indent(&dst, o.Value(), prefix, indent)
		if err != nil {
			return objects.NewError(objects.NewStringNoSize(err.Error())), nil
		}
		return objects.NewBytes(dst.Bytes()), nil
	case *objects.String:
		var dst bytes.Buffer
		err = json.Indent(&dst, []byte(o.Value()), prefix, indent)
		if err != nil {
			return objects.NewError(objects.NewStringNoSize(err.Error())), nil
		}
		return objects.NewBytes(dst.Bytes()), nil
	default:
		return nil, objects.NewInvalidArgumentError(0, "bytes/string", args[0].TypeName())
	}
}

// HTMLEscape escapes certain characters in a JSON string or byte slice to their HTML-safe equivalents.
// Accepts one argument of type `*objects.Bytes` or `*objects.String`.
// Returns a new `*objects.Bytes` containing the escaped output or an error if an invalid argument type is provided.
// Errors if the number of arguments is not exactly one.
func (j *Json) HTMLEscape(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	switch o := args[0].(type) {
	case *objects.Bytes:
		var dst bytes.Buffer
		json.HTMLEscape(&dst, o.Value())
		return objects.NewBytes(dst.Bytes()), nil
	case *objects.String:
		var dst bytes.Buffer
		json.HTMLEscape(&dst, []byte(o.Value()))
		return objects.NewBytes(dst.Bytes()), nil
	default:
		return nil, objects.NewInvalidArgumentError(0, "bytes/string", args[0].TypeName())
	}
}
