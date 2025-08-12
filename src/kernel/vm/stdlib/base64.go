package stdlib

import (
	"encoding/base64"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

var base64Module = map[string]objects.IObject{
	"encode":         objects.NewUserFunction("encode", FuncAYRS(base64.StdEncoding.EncodeToString)),
	"decode":         objects.NewUserFunction("decode", FuncASRYE(base64.StdEncoding.DecodeString)),
	"raw_encode":     objects.NewUserFunction("raw_encode", FuncAYRS(base64.RawStdEncoding.EncodeToString)),
	"raw_decode":     objects.NewUserFunction("raw_decode", FuncASRYE(base64.RawStdEncoding.DecodeString)),
	"url_encode":     objects.NewUserFunction("url_encode", FuncAYRS(base64.URLEncoding.EncodeToString)),
	"url_decode":     objects.NewUserFunction("url_decode", FuncASRYE(base64.URLEncoding.DecodeString)),
	"raw_url_encode": objects.NewUserFunction("raw_url_encode", FuncAYRS(base64.RawURLEncoding.EncodeToString)),
	"raw_url_decode": objects.NewUserFunction("raw_url_decode", FuncASRYE(base64.RawURLEncoding.DecodeString)),
}
