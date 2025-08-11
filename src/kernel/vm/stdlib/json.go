package stdlib

import (
	"bytes"
	gojson "encoding/json"

	"github.com/markel1974/c64emu/src/kernel/vm/errors"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
	"github.com/markel1974/c64emu/src/kernel/vm/stdlib/json"
)

var jsonModule = map[string]objects.IObject{
	"decode": &objects.UserFunction{
		Name:  "decode",
		Value: jsonDecode,
	},
	"encode": &objects.UserFunction{
		Name:  "encode",
		Value: jsonEncode,
	},
	"indent": &objects.UserFunction{
		Name:  "encode",
		Value: jsonIndent,
	},
	"html_escape": &objects.UserFunction{
		Name:  "html_escape",
		Value: jsonHTMLEscape,
	},
}

func jsonDecode(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 1 {
		return nil, errors.ErrWrongNumArguments
	}

	switch o := args[0].(type) {
	case *objects.Bytes:
		v, err := json.Decode(o.Value)
		if err != nil {
			return &objects.Error{
				Value: &objects.String{Value: err.Error()},
			}, nil
		}
		return v, nil
	case *objects.String:
		v, err := json.Decode([]byte(o.Value))
		if err != nil {
			return &objects.Error{
				Value: &objects.String{Value: err.Error()},
			}, nil
		}
		return v, nil
	default:
		return nil, errors.NewInvalidArgumentType("first", "bytes/string", args[0].TypeName())
	}
}

func jsonEncode(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 1 {
		return nil, errors.ErrWrongNumArguments
	}

	b, err := json.Encode(args[0])
	if err != nil {
		return &objects.Error{Value: &objects.String{Value: err.Error()}}, nil
	}

	return &objects.Bytes{Value: b}, nil
}

func jsonIndent(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 3 {
		return nil, errors.ErrWrongNumArguments
	}

	prefix, ok := objects.ToString(args[1])
	if !ok {
		return nil, errors.NewInvalidArgumentType("prefix", "string(compatible)", args[1].TypeName())
	}

	indent, ok := objects.ToString(args[2])
	if !ok {
		return nil, errors.NewInvalidArgumentType("indent", "string(compatible)", args[2].TypeName())
	}

	switch o := args[0].(type) {
	case *objects.Bytes:
		var dst bytes.Buffer
		err := gojson.Indent(&dst, o.Value, prefix, indent)
		if err != nil {
			return &objects.Error{
				Value: &objects.String{Value: err.Error()},
			}, nil
		}
		return &objects.Bytes{Value: dst.Bytes()}, nil
	case *objects.String:
		var dst bytes.Buffer
		err := gojson.Indent(&dst, []byte(o.Value), prefix, indent)
		if err != nil {
			return &objects.Error{
				Value: &objects.String{Value: err.Error()},
			}, nil
		}
		return &objects.Bytes{Value: dst.Bytes()}, nil
	default:
		return nil, errors.NewInvalidArgumentType("first", "bytes/string", args[0].TypeName())
	}
}

func jsonHTMLEscape(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 1 {
		return nil, errors.ErrWrongNumArguments
	}

	switch o := args[0].(type) {
	case *objects.Bytes:
		var dst bytes.Buffer
		gojson.HTMLEscape(&dst, o.Value)
		return &objects.Bytes{Value: dst.Bytes()}, nil
	case *objects.String:
		var dst bytes.Buffer
		gojson.HTMLEscape(&dst, []byte(o.Value))
		return &objects.Bytes{Value: dst.Bytes()}, nil
	default:
		return nil, errors.NewInvalidArgumentType("first", "bytes/string", args[0].TypeName())
	}
}
