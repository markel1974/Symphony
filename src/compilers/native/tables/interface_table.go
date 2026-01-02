package tables

import (
	"fmt"
	"go/ast"

	"github.com/markel1974/c64emu/src/vm/objects"
)

// MethodDescription represents a method with its name, input parameters, and return types.
type MethodDescription struct {
	Name        string
	InputParams []string
	ReturnTypes []string
}

// InterfaceDescription represents metadata about an interface, including its name, methods, and struct-related properties.
type InterfaceDescription struct {
	name      string
	Methods   []*MethodDescription
	fieldName string
	fieldNode ast.Node
	container string
	kind      string
	offset    int // Offset in bytes if used as a field in a struct
	totalSize int // Fixed size of the interface (e.g. 16 bytes: type + ptr)
}

// NewInterfaceDescription creates and returns a new InterfaceDescription with the given name and method descriptions.
func NewInterfaceDescription(id string, methods []*MethodDescription) *InterfaceDescription {
	v := &InterfaceDescription{
		name:      id,
		Methods:   methods,
		totalSize: 16,
	}
	return v
}

// Name returns the name of the interface description instance.
func (id *InterfaceDescription) Name() string {
	return id.name
}

// FieldName returns the name of the field associated with the interface description.
func (id *InterfaceDescription) FieldName() string {
	return id.fieldName
}

// SetFieldName updates the internal field name of the InterfaceDescription instance.
func (id *InterfaceDescription) SetFieldName(name string) {
	id.fieldName = name
}

// FieldBase returns the base name of the interface represented by the InterfaceDescription.
func (id *InterfaceDescription) FieldBase() string {
	return id.name
}

// FieldClone creates a deep copy of the InterfaceDescription as an IStructField, retaining all method and field attributes.
func (id *InterfaceDescription) FieldClone() IStructField {
	methods := make([]*MethodDescription, len(id.Methods))
	for idx, v := range id.Methods {
		// Fix: Preserve the method name from the source method 'v', not 'id.name'
		md := &MethodDescription{
			Name:        v.Name,
			InputParams: make([]string, len(v.InputParams)),
			ReturnTypes: make([]string, len(v.ReturnTypes)),
		}
		copy(md.InputParams, v.InputParams)
		copy(md.ReturnTypes, v.ReturnTypes)
		methods[idx] = md
	}
	out := NewInterfaceDescription(id.name, methods)
	out.container = id.container
	out.kind = id.kind
	out.fieldName = id.fieldName
	out.fieldNode = id.fieldNode
	out.offset = id.offset
	return out
}

// SetOptions updates the container and kind of the interface description, optionally indicating if it's a pointer.
func (id *InterfaceDescription) SetOptions(_ bool, container string, kind string) {
	id.container = container
	id.kind = kind
}

// Options retrieves the pointer flag, container name, and kind associated with the interface description instance.
func (id *InterfaceDescription) Options() (bool, string, string) {
	return false, id.container, id.kind
}

// FieldNode returns the underlying abstract syntax tree (AST) node associated with the interface description.
func (id *InterfaceDescription) FieldNode() ast.Node {
	return id.fieldNode
}

// SetFieldNode assigns the given AST node to the fieldNode property of the InterfaceDescription.
func (id *InterfaceDescription) SetFieldNode(node ast.Node) {
	id.fieldNode = node
}

// Offset returns the offset in bytes of the interface when used as a field in a struct.
func (id *InterfaceDescription) Offset() int {
	return id.offset
}

// SetOffset updates the offset value for the interface description in bytes.
func (id *InterfaceDescription) SetOffset(offset int) {
	id.offset = offset
}

// Definition returns the InterfaceDescription instance as an empty interface.
func (id *InterfaceDescription) Definition() IStructField {
	return id
}

// BindDefinition associates an external definition with the interface; this method is currently a no-operation (No-op).
func (id *InterfaceDescription) BindDefinition(def IStructField) {
	// No-op
}

// IsPlaceholder checks if the interface description is a placeholder and returns false in its current implementation.
func (id *InterfaceDescription) IsPlaceholder() bool {
	return false
}

// IsFinalized checks if the interface description is in a finalized state.
func (id *InterfaceDescription) IsFinalized() bool {
	return true
}

// SetFinalized marks the interface description as finalized or not; currently, this method performs no operation (no-op).
func (id *InterfaceDescription) SetFinalized(finalized bool) {
	// No-op
}

// TotalSize returns the total memory size in bytes for this interface, including its type and pointer.
func (id *InterfaceDescription) TotalSize() int {
	return id.totalSize
}

// SetTotalSize updates the totalSize field of the InterfaceDescription to the specified size in bytes.
func (id *InterfaceDescription) SetTotalSize(size int) {
	id.totalSize = size
}

// InterfaceTable manages a collection of interface descriptions associated with specific scopes and a gatekeeper.
type InterfaceTable struct {
	gk        objects.IGateKeeper
	scopes    *Scopes
	container map[string]*InterfaceDescription
}

// NewInterfaceTable creates and returns a new InterfaceTable instance initialized with the provided gatekeeper and scopes.
func NewInterfaceTable(gk objects.IGateKeeper, scopes *Scopes) *InterfaceTable {
	it := &InterfaceTable{
		gk:        gk,
		scopes:    scopes,
		container: make(map[string]*InterfaceDescription),
	}
	return it
}

// CreateInterface creates a new InterfaceDescription with the given id and methods and stores it in the container.
func (it *InterfaceTable) CreateInterface(id string, methods []*MethodDescription) *InterfaceDescription {
	v := NewInterfaceDescription(id, methods)
	it.container[id] = v
	return v
}

// Keys returns a slice of all keys present in the container map of the InterfaceTable.
func (it *InterfaceTable) Keys() []string {
	keys := make([]string, 0, len(it.container))
	for k := range it.container {
		keys = append(keys, k)
	}
	return keys
}

// Container returns the internal map of interface descriptions indexed by their names.
func (it *InterfaceTable) Container() map[string]*InterfaceDescription {
	return it.container
}

// Add registers a new interface in the table using its name and methods, returning an error if the interface already exists.
func (it *InterfaceTable) Add(name string, node *ast.InterfaceType) error {
	if _, exists := it.container[name]; exists {
		return fmt.Errorf("interface '%s' already defined", name)
	}

	var methods []*MethodDescription
	if node.Methods != nil {
		for _, field := range node.Methods.List {
			if len(field.Names) > 0 {
				if funcType, ok := field.Type.(*ast.FuncType); ok {
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
	it.container[name] = NewInterfaceDescription(name, methods)
	return nil
}

// Get retrieves an InterfaceDescription by its name from the InterfaceTable.
// Returns the description and a boolean indicating if it was found.
func (it *InterfaceTable) Get(name string) (*InterfaceDescription, bool) {
	desc, ok := it.container[name]
	return desc, ok
}

// Has checks if an interface with the given name exists in the container and returns true if it does, otherwise false.
func (it *InterfaceTable) Has(name string) bool {
	_, ok := it.container[name]
	return ok
}
