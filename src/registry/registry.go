package registry

import (
	"github.com/markel1974/c64emu/src/references"
	"log"
)

// _componentFactories holds a list of registered factories for creating and managing components in the system.
var _componentFactories []references.IFactory

// _componentFactoryHelper is a map used to store and manage component factories indexed by their unique identifiers.
var _componentFactoryHelper = make(map[string]references.IFactory)

// RegisterComponentFactory adds a new component factory to the global registry for structured component instantiation.
// The function logs an error and halts the program if the factory is nil, has an empty identifier, or is already registered.
// GateKeeper identifiers must be unique, and the registry updates the helper map and list upon successful registration.
func RegisterComponentFactory(factory references.IFactory) {
	if factory == nil {
		log.Fatal("cannot register nil component factory")
	}
	if len(factory.Identifier()) == 0 {
		log.Fatal("cannot register component factory with empty identifier")
	}
	if _, ok := _componentFactoryHelper[factory.Identifier()]; ok {
		log.Fatal("component factory with identifier " + factory.Identifier() + " already registered")
	}
	log.Println("registered component factory " + factory.Identifier())
	_componentFactoryHelper[factory.Identifier()] = factory
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
