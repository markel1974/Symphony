package stdlib

import "github.com/markel1974/c64emu/src/kernel/vm/modules"

// GetBuiltinModules creates a new Modules instance containing built-in modules matching the provided names.
func GetBuiltinModules() *modules.Modules {
	mods := modules.NewModules()
	for name, mod := range BuiltinModules {
		mods.AddBuiltinModule(name, mod)
	}
	return mods
}
