package stdlib

import (
	"encoding/base64"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// base64Module provides encoding and decoding utility functions for base64, raw, and URL-safe base64 formats.
var base64Module = map[string]objects.IObject{
	"encode":         objects.NewFunctionUser("encode", objects.FuncAYRS(base64.StdEncoding.EncodeToString)),
	"decode":         objects.NewFunctionUser("decode", objects.FuncASRYE(base64.StdEncoding.DecodeString)),
	"raw_encode":     objects.NewFunctionUser("raw_encode", objects.FuncAYRS(base64.RawStdEncoding.EncodeToString)),
	"raw_decode":     objects.NewFunctionUser("raw_decode", objects.FuncASRYE(base64.RawStdEncoding.DecodeString)),
	"url_encode":     objects.NewFunctionUser("url_encode", objects.FuncAYRS(base64.URLEncoding.EncodeToString)),
	"url_decode":     objects.NewFunctionUser("url_decode", objects.FuncASRYE(base64.URLEncoding.DecodeString)),
	"raw_url_encode": objects.NewFunctionUser("raw_url_encode", objects.FuncAYRS(base64.RawURLEncoding.EncodeToString)),
	"raw_url_decode": objects.NewFunctionUser("raw_url_decode", objects.FuncASRYE(base64.RawURLEncoding.DecodeString)),
}
