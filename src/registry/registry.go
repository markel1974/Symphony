package registry

import "github.com/markel1974/c64emu/src/references"

var _componentFactories = make(map[string]references.IFactory)

func RegisterComponentFactory(factory references.IFactory) {
	_componentFactories[factory.Identifier()] = factory
}

func ComponentFactories() map[string]references.IFactory {
	container := make(map[string]references.IFactory)
	for k, v := range _componentFactories {
		container[k] = v
	}
	return container
}
