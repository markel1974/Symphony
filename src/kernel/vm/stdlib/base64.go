package stdlib

import (
	"encoding/base64"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

var base64Module = map[string]objects.IObject{
	"encode": &objects.UserFunction{
		Value: FuncAYRS(base64.StdEncoding.EncodeToString),
	},
	"decode": &objects.UserFunction{
		Value: FuncASRYE(base64.StdEncoding.DecodeString),
	},
	"raw_encode": &objects.UserFunction{
		Value: FuncAYRS(base64.RawStdEncoding.EncodeToString),
	},
	"raw_decode": &objects.UserFunction{
		Value: FuncASRYE(base64.RawStdEncoding.DecodeString),
	},
	"url_encode": &objects.UserFunction{
		Value: FuncAYRS(base64.URLEncoding.EncodeToString),
	},
	"url_decode": &objects.UserFunction{
		Value: FuncASRYE(base64.URLEncoding.DecodeString),
	},
	"raw_url_encode": &objects.UserFunction{
		Value: FuncAYRS(base64.RawURLEncoding.EncodeToString),
	},
	"raw_url_decode": &objects.UserFunction{
		Value: FuncASRYE(base64.RawURLEncoding.DecodeString),
	},
}
