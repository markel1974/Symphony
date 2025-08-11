package stdlib

/*
//go:generate go run gensrcmods.go

import (
	"github.com/markel1974/injector/src/vm/modules"
)

// AllModuleNames returns a list of all default module names.
func AllModuleNames() []string {
	var names []string
	for name := range BuiltinModules {
		names = append(names, name)
	}
	for name := range SourceModules {
		names = append(names, name)
	}
	return names
}

// GetModuleMap returns the module map that includes all modules
// for the given module names.
func GetModuleMap(names ...string) *modules.ModuleMap {
	mods := modules.NewModuleMap()
	for _, name := range names {
		if mod := BuiltinModules[name]; mod != nil {
			mods.AddBuiltinModule(name, mod)
		}
		if mod := SourceModules[name]; mod != "" {
			mods.AddSourceModule(name, []byte(mod))
		}
	}
	return mods
}


*/
