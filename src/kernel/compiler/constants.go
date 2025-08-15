package compiler

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Constants is a structure that manages a collection of objects and provides indexing functionality for efficient retrieval.
type Constants struct {
	constants []objects.IObject
	cache     map[string]int
}

// NewConstants initializes and returns a new instance of the Constants struct with empty data structures.
func NewConstants() *Constants {
	return &Constants{
		constants: []objects.IObject{},
		cache:     make(map[string]int),
	}
}

// Add appends the given object to the constants slice and returns its index; caches the index if an ID is provided.
func (c *Constants) Add(id string, obj objects.IObject) int {
	c.constants = append(c.constants, obj)
	nameIndex := len(c.constants) - 1
	if len(id) > 0 {
		c.cache[id] = nameIndex
	}
	return nameIndex
}

// Get retrieves the index of the constant associated with the given id from the cache.
// It returns the index and true if found, otherwise 0 and false.
func (c *Constants) Get(id string) (int, bool) {
	index, ok := c.cache[id]
	if !ok {
		return 0, false
	}
	return index, true
}

// SetIndex updates an object at a specific index.
func (c *Constants) SetIndex(index int, obj objects.IObject) error {
	if index < 0 || index >= len(c.constants) {
		return fmt.Errorf("constants: index %d out of bounds", index)
	}
	c.constants[index] = obj
	return nil
}

// Len returns the number of constants currently stored in the Constants structure.
func (c *Constants) Len() int {
	return len(c.constants)
}

// GetByIndex retrieves the object at the specified index and a boolean indicating success or failure of the operation.
func (c *Constants) GetByIndex(index int) (objects.IObject, bool) {
	if index < 0 || index >= len(c.constants) {
		return nil, false
	}
	return c.constants[index], true
}

// Retrieve returns the list of all constant objects stored in the Constants structure.
func (c *Constants) Retrieve() []objects.IObject {
	return c.constants
}
