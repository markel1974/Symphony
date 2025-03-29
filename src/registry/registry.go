package registry

import "github.com/markel1974/c64emu/src/references"

var _componentFactories []references.IFactory

func RegisterComponentFactory(factory references.IFactory) {
	_componentFactories = append(_componentFactories, factory)
}

func ComponentFactories() []references.IFactory {
	if len(_componentFactories) == 0 {
		return nil
	}
	factories := make([]references.IFactory, len(_componentFactories))
	copy(factories, _componentFactories)
	return factories
}
