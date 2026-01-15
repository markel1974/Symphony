package compiler

import (
	"go/ast"
	"go/token"
	"reflect"

	"github.com/markel1974/symphony/src/compilers/native/tables"
)

// TypeCompatibility is a structure used to manage type-checking, compatibility, and implementation relationships.
// It includes tables for structs, interfaces, and functions, and tracks implementations of interfaces by structs.
type TypeCompatibility struct {
	fileSet         *token.FileSet
	definitionTable *tables.DefinitionTable
	functionTable   *tables.FunctionTable
	implementations map[string][]string
}

// NewTypeCompatibility initializes a TypeCompatibility instance with provided struct, interface, and function tables.
func NewTypeCompatibility(definitionTable *tables.DefinitionTable, functionTable *tables.FunctionTable) *TypeCompatibility {
	return &TypeCompatibility{
		definitionTable: definitionTable,
		functionTable:   functionTable,
		implementations: make(map[string][]string),
	}
}

// Setup initializes the TypeCompatibility instance by setting the provided FileSet as its fileSet property.
func (tc *TypeCompatibility) Setup(fileSet *token.FileSet, _ func(node ast.Node) error) error {
	tc.fileSet = fileSet
	return nil
}

// Prepare analyzes and verifies struct implementation of interfaces, storing valid relationships for central access.
func (tc *TypeCompatibility) Prepare() error {
	//fmt.Println("Running interface implementation check...")
	// Iterate over each defined struct
	for _, structName := range tc.definitionTable.StructKeys() { //tc.structTable.Container() {
		// Iterate over each defined interface
		for interfaceName, interfaceDesc := range tc.definitionTable.InterfaceContainer() {
			implements, err := tc.checkStructImplementsInterface(structName, interfaceDesc)
			if err != nil {
				// TODO collect all errors
				return err
			}
			if implements {
				// Store the relationship
				tc.implementations[structName] = append(tc.implementations[structName], interfaceName)
				//fmt.Printf("=> SUCCESS: Struct '%s' implements interface '%s'\n", structName, interfaceName)
			}
		}
	}
	// add found implementations to StructTable for centralized access
	tc.definitionTable.StructSetImplementations(tc.implementations)
	return nil
}

// Compile executes the final step to confirm type compatibility checks and prepares the system for further operations.
func (tc *TypeCompatibility) Compile() error {
	return nil
}

// Finalize finalizes the TypeCompatibility structure by performing necessary cleanup or concluding operations. Returns an error if it fails.
func (tc *TypeCompatibility) Finalize() error {
	return nil
}

// checkStructImplementsInterface determines if a struct implements all methods defined by a given interface description.
func (tc *TypeCompatibility) checkStructImplementsInterface(structName string, interfaceDesc *tables.InterfaceDescription) (bool, error) {
	for _, requiredMethod := range interfaceDesc.Methods {
		// The "mangled" method name for the struct is "StructName.MethodName"
		mangledMethodName := tables.GetMangledName(structName, requiredMethod.Name)
		// Look up the function (method) description in the functionTable
		var structMethod *tables.FunctionDescription
		for i := 0; i < tc.functionTable.Len(); i++ {
			fd, _ := tc.functionTable.Get(i)
			if fd.Name == mangledMethodName {
				structMethod = fd
				break
			}
		}
		if structMethod == nil {
			// Required method not found on struct
			return false, nil
		}
		// Compare method signatures
		// NOTE: receiver counts as first parameter for struct method
		numStructParams := len(structMethod.InputNames)
		if len(requiredMethod.InputParams) != numStructParams {
			return false, nil
		}
		if !reflect.DeepEqual(requiredMethod.InputParams, structMethod.InputTypes) {
			return false, nil
		}
		if !reflect.DeepEqual(requiredMethod.ReturnTypes, structMethod.ReturnTypes) {
			return false, nil
		}
	}
	// If we get here, all required methods were found and signatures match
	return true, nil
}
