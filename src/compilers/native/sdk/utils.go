package sdk

import (
	"hash/fnv"

	"github.com/markel1974/c64emu/src/vm/objects"
)

func BuildContainer2(functions []objects.IObject, constants map[string]objects.IObject) map[string]objects.IObject {
	container := make(map[string]objects.IObject)
	for _, obj := range functions {
		fn, ok := obj.(*objects.FuncImport)
		if ok {
			container[fn.Name()] = fn
		}
	}
	for id, c := range constants {
		container[id] = c
	}
	return container
}

func PackageIDFromString(name string) uint32 {
	hash := fnv.New32a()
	_, _ = hash.Write([]byte(name))
	return hash.Sum32()
}
