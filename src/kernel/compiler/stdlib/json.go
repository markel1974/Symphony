package stdlib

import (
	"bytes"
	"encoding/json"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// jsonModule is a map of JSON-related functions (decode, encode, indent, html_escape) implementing the IObject interface.
var _jsonModule = map[string]objects.IObject{
	"decode":      objects.NewFunctionUser("decode", jsonDecode),
	"encode":      objects.NewFunctionUser("encode", jsonEncode),
	"indent":      objects.NewFunctionUser("encode", jsonIndent),
	"html_escape": objects.NewFunctionUser("html_escape", jsonHTMLEscape),
}

// jsonDecode parses a JSON-encoded bytes or string argument into a map-like object or returns an error if decoding fails.
func jsonDecode(args ...objects.IObject) (ret objects.IObject, err error) {
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
		return nil, objects.NewInvalidArgumentError("first", "bytes/string", args[0].TypeName())
	}
	d := make(map[string]interface{})
	if err = json.Unmarshal(data, &d); err != nil {
		return objects.NewError(objects.NewStringNoSize(err.Error())), nil
	}
	result := objects.NewMap(objects.FromMap(d))
	return result, nil
}

// jsonEncode serializes an IObject into its JSON-encoded byte representation and returns it as a Bytes object.
// jsonEncode expects exactly one argument; otherwise, it returns an error for incorrect argument count.
// jsonEncode returns a serialized result or an error object if JSON marshalling fails.
func jsonEncode(args ...objects.IObject) (ret objects.IObject, err error) {
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
	prefix, ok := objects.ToString(args[1])
	if !ok {
		return nil, objects.NewInvalidArgumentError("prefix", "string(compatible)", args[1].TypeName())
	}
	indent, ok := objects.ToString(args[2])
	if !ok {
		return nil, objects.NewInvalidArgumentError("indent", "string(compatible)", args[2].TypeName())
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
		return nil, objects.NewInvalidArgumentError("first", "bytes/string", args[0].TypeName())
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
		return nil, objects.NewInvalidArgumentError("first", "bytes/string", args[0].TypeName())
	}
}
