package compiler

import (
	"fmt"
	"io"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Constants is a structure that manages a collection of objects and provides indexing functionality for efficient retrieval.
type Constants struct {
	constants  []objects.IObject
	cache      map[string]int
	loader     bytecode.ILoader
	builtinLen int
}

// NewConstants initializes and returns a new instance of the Constants struct with empty data structures.
func NewConstants(loader bytecode.ILoader, builtinLen int) *Constants {
	c := &Constants{
		constants:  nil,
		cache:      make(map[string]int),
		loader:     loader,
		builtinLen: builtinLen,
	}
	return c
}

// Setup initializes the constants slice with the built-in functions.
func (c *Constants) Setup() error {
	c.cache = make(map[string]int)
	c.constants = make([]objects.IObject, c.builtinLen)
	if c.builtinLen == 0 {
		return nil
	}
	for idx := 0; idx < c.builtinLen; idx++ {
		bi := c.loader.Builtin(idx)
		if bi == nil {
			return fmt.Errorf("builtin %d not found", idx)
		}
		c.constants[idx] = bi
		c.cache[bi.Name()] = idx
	}
	return nil
}

// Print prints the constants to the provided writer.
func (c *Constants) Print(writer io.Writer) {
	for idx, v := range c.constants {
		_, _ = fmt.Fprintf(writer, "%d => %s", idx, v.String())
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

// AddOrGet adds a new constant to the pool or returns the index of an existing identical constant.
// This prevents duplicating constants in the bytecode.
func (c *Constants) AddOrGet(id string, obj objects.IObject) int {
	for i, constant := range c.constants {
		if obj.Equals(constant) {
			return i
		}
	}
	return c.Add(id, obj)
}

// Get retrieves the index of the constant associated with the given id from the cache.
// It returns the index and true if found, otherwise 0 and false.
func (c *Constants) Get(id string) (int, bool) {
	index, ok := c.cache[id]
	return index, ok
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
