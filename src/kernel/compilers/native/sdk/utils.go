package sdk

import "github.com/markel1974/c64emu/src/kernel/vm/objects"

func BuildContainer(functions []objects.IObject, constants map[string]objects.IObject) map[string]objects.IObject {
	container := make(map[string]objects.IObject)
	for _, obj := range functions {
		fn, ok := obj.(*objects.FuncExternal)
		if ok {
			container[fn.Name()] = fn
		}
	}
	for id, c := range constants {
		container[id] = c
	}
	return container
}
