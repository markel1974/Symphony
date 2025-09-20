package sdk

import (
	"github.com/markel1974/c64emu/src/vm/objects"
)

func BuildImporter(functions []objects.IObject, constants map[string]objects.IObject) map[string]objects.IObject {
	container := make(map[string]objects.IObject)
	for _, obj := range functions {
		switch fn := obj.(type) {
		case *objects.Func:
			container[fn.Name()] = fn
		case *objects.FuncImport:
			container[fn.Name()] = fn
		case *objects.FuncJit:
			container[fn.Name()] = fn
		}
	}
	for id, c := range constants {
		container[id] = c
	}
	return container
}
