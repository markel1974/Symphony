package registry

import (
	"github.com/markel1974/c64emu/src/references"
	"log"
)

// _componentFactories is a slice holding references to all registered component factories implementing IFactory.
var _componentFactories []references.IFactory

var _componentFactoryHelp = make(map[string]references.IFactory)

// RegisterComponentFactory appends a new component factory to the global list of registered factories.
func RegisterComponentFactory(factory references.IFactory) {
	if factory == nil {
		log.Fatal("cannot register nil component factory")
	}
	if len(factory.Identifier()) == 0 {
		log.Fatal("cannot register component factory with empty identifier")
	}
	if _, ok := _componentFactoryHelp[factory.Identifier()]; ok {
		log.Fatal("component factory with identifier " + factory.Identifier() + " already registered")
	}
	log.Println("registered component factory " + factory.Identifier())
	_componentFactoryHelp[factory.Identifier()] = factory
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
