package stdlib

import "github.com/markel1974/c64emu/src/kernel/vm/objects"

// BuiltinModules are builtin type standard library modules.
var BuiltinModules = map[string]map[string]objects.IObject{
	//"os":   osModule,
	"fmt":   fmtModule,
	"math":  mathModule,
	"text":  textModule,
	"times": timesModule,
	"rand":  randModule,
	//"fmt":    fmtSafeModule,
	"json":   jsonModule,
	"base64": base64Module,
	"hex":    hexModule,
}
