package registry

import "github.com/markel1974/c64emu/src/references"

// _componentFactories is a slice holding references to all registered component factories implementing IFactory.
var _componentFactories []references.IFactory

// RegisterComponentFactory appends a new component factory to the global list of registered factories.
func RegisterComponentFactory(factory references.IFactory) {
	_componentFactories = append(_componentFactories, factory)
}

// ComponentFactories returns a copy of the registered component factories or nil if no factories are registered.
func ComponentFactories() []references.IFactory {
	if len(_componentFactories) == 0 {
		return nil
	}
	factories := make([]references.IFactory, len(_componentFactories))
	copy(factories, _componentFactories)
	return factories
}
