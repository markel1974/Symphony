package compiler

import (
	"fmt"
	"go/ast"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// MethodDescription represents metadata about a method, including its name, input parameters, and return types.
type MethodDescription struct {
	Name        string
	InputParams []string
	ReturnTypes []string
}

// InterfaceDescription represents the structure of an interface with its name and associated method descriptions.
type InterfaceDescription struct {
	Name    string
	Methods []*MethodDescription
}

// InterfaceTable provides a data structure for managing a collection of InterfaceDescription objects indexed by name.
// It relies on GateKeeper for handling object lifecycle and takes Scopes for scope context management.
// The container map holds the registered interface descriptions for lookup and manipulation.
type InterfaceTable struct {
	gk        objects.IGateKeeper
	scopes    *Scopes
	container map[string]*InterfaceDescription
}

// NewInterfaceTable creates and initializes a new InterfaceTable with a gatekeeper, scopes, and an empty container map.
func NewInterfaceTable(gk objects.IGateKeeper, scopes *Scopes) *InterfaceTable {
	return &InterfaceTable{
		gk:        gk,
		scopes:    scopes,
		container: make(map[string]*InterfaceDescription),
	}
}

// Add registers a new interface in the InterfaceTable with its methods and attributes defined in the given AST node.
// It returns an error if the interface name is already defined or if method parsing fails.
func (it *InterfaceTable) Add(name string, node *ast.InterfaceType) error {
	if _, exists := it.container[name]; exists {
		return fmt.Errorf("interface '%s' already defined", name)
	}

	var methods []*MethodDescription
	if node.Methods != nil {
		for _, field := range node.Methods.List {
			if len(field.Names) > 0 {
				if funcType, ok := field.Type.(*ast.FuncType); ok {
					// Estrai i tipi dei parametri e dei valori di ritorno
					inputParams, err := GetReceivers(funcType.Params)
					if err != nil {
						return fmt.Errorf("error parsing params for method %s in interface %s: %w", field.Names[0].Name, name, err)
					}
					returnTypes, err := GetReceivers(funcType.Results)
					if err != nil {
						return fmt.Errorf("error parsing return types for method %s in interface %s: %w", field.Names[0].Name, name, err)
					}

					method := &MethodDescription{
						Name:        field.Names[0].Name,
						InputParams: inputParams,
						ReturnTypes: returnTypes,
					}
					methods = append(methods, method)
				}
			}
		}
	}

	it.container[name] = &InterfaceDescription{
		Name:    name,
		Methods: methods,
	}

	return nil
}

// Get retrieves an InterfaceDescription and a boolean indicating its existence from the InterfaceTable by name.
func (it *InterfaceTable) Get(name string) (*InterfaceDescription, bool) {
	desc, ok := it.container[name]
	return desc, ok
}

// Has checks if the specified interface name exists in the InterfaceTable. Returns true if found, otherwise false.
func (it *InterfaceTable) Has(name string) bool {
	_, ok := it.container[name]
	return ok
}
