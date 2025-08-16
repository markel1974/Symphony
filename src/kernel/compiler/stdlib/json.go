package stdlib

import (
	"bytes"
	"encoding/json"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// jsonModule is a map of JSON-related functions (decode, encode, indent, html_escape) implementing the IObject interface.
var _jsonModule = map[string]objects.IObject{
	"Unmarshal":  objects.NewFunctionModule(objects.FunctionModuleDef, "Unmarshal", jsonUnmarshal),
	"Marshal":    objects.NewFunctionModule(objects.FunctionModuleDef, "Marshal", jsonMarshal),
	"Indent":     objects.NewFunctionModule(objects.FunctionModuleDef, "Indent", jsonIndent),
	"HTMLEscape": objects.NewFunctionModule(objects.FunctionModuleDef, "html_escape", jsonHTMLEscape),
}

// jsonUnmarshal parses a JSON-encoded bytes or string argument into a map-like object or returns an error if decoding fails.
func jsonUnmarshal(args ...objects.IObject) (ret objects.IObject, err error) {
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

// jsonMarshal serializes an IObject into its JSON-encoded byte representation and returns it as a Bytes object.
// jsonMarshal expects exactly one argument; otherwise, it returns an error for incorrect argument count.
// jsonMarshal returns a serialized result or an error object if JSON marshalling fails.
func jsonMarshal(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	result, err := json.Marshal(objects.ToInterface(args[0]))
	if err != nil {
		return objects.NewError(objects.NewStringNoSize(err.Error())), nil
	}
	return objects.NewBytes(result), nil
}

// jsonIndent formats JSON input with the specified prefix and indentation, returning the formatted JSON as Bytes.
// It accepts three arguments: a string/bytes JSON input, a string prefix, and a string indentation.
// Returns an error if the input is invalid or the formatting fails.
func jsonIndent(args ...objects.IObject) (ret objects.IObject, err error) {
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

// jsonHTMLEscape escapes special HTML characters in a bytes or string object using JSON encoding.
// Returns a new bytes object with the escaped data or an error if the input type is invalid.
func jsonHTMLEscape(args ...objects.IObject) (ret objects.IObject, err error) {
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
