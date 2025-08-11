package objects

import (
	"fmt"
	"strings"

	"github.com/markel1974/c64emu/src/kernel/vm/errors"
)

// ImmutableMap represents a read-only map structure where keys are strings and values are of type IObject.
type ImmutableMap struct {
	ObjectImpl
	Value map[string]IObject
}

// NewImmutableMap creates a new instance of ImmutableMap with the provided map of string keys and IObject values.
func NewImmutableMap(value map[string]IObject) *ImmutableMap {
	return &ImmutableMap{Value: value}
}

// Length returns the number of elements in the ImmutableMap.
func (o *ImmutableMap) Length() int {
	return len(o.Value)
}

// TypeName returns the type name of the ImmutableMap as a string.
func (o *ImmutableMap) TypeName() string {
	return "immutable-map"
}

// String generates a string representation of the ImmutableMap in key-value pair format.
func (o *ImmutableMap) String() string {
	var pairs []string
	for k, v := range o.Value {
		pairs = append(pairs, fmt.Sprintf("%s: %s", k, v.String()))
	}
	return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
}

// Copy creates and returns a deep copy of the ImmutableMap, duplicating all key-value pairs.
func (o *ImmutableMap) Copy() IObject {
	c := make(map[string]IObject)
	for k, v := range o.Value {
		c[k] = v.Copy()
	}
	return NewMap(c)
}

// Falsy returns true if the map is empty, indicating it is considered "falsy", otherwise false.
func (o *ImmutableMap) Falsy() bool {
	return len(o.Value) == 0
}

// IndexGet retrieves the value associated with the given index in the ImmutableMap. Returns an error for invalid index types.
func (o *ImmutableMap) IndexGet(index IObject) (res IObject, err error) {
	strIdx, ok := ToString(index)
	if !ok {
		err = errors.ErrInvalidIndexType
		return
	}
	res, ok = o.Value[strIdx]
	if !ok {
		res = UndefinedValue
	}
	return
}

// Equals determines whether the current ImmutableMap is equal to another IObject by comparing their key-value pairs.
func (o *ImmutableMap) Equals(in IObject) bool {
	var xVal map[string]IObject
	switch x := in.(type) {
	case *Map:
		xVal = x.values
	case *ImmutableMap:
		xVal = x.Value
	default:
		return false
	}
	if len(o.Value) != len(xVal) {
		return false
	}
	for k, v := range o.Value {
		tv := xVal[k]
		if !v.Equals(tv) {
			return false
		}
	}
	return true
}

// Iterate returns an iterator for traversing the key-value pairs in the immutable map.
func (o *ImmutableMap) Iterate() IIterator {
	var keys []string
	for k := range o.Value {
		keys = append(keys, k)
	}
	return &MapIterator{
		v: o.Value,
		k: keys,
		l: len(keys),
	}
}

// CanIterate returns true, indicating that the ImmutableMap supports iteration.
func (o *ImmutableMap) CanIterate() bool {
	return true
}
