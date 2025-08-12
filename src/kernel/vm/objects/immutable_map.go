package objects

import (
	"fmt"
	"strings"

	"github.com/markel1974/c64emu/src/kernel/vm/errors"
)

// ImmutableMap represents a read-only map structure where keys are strings and values are of type IObject.
type ImmutableMap struct {
	ObjectImpl
	values map[string]IObject
}

// NewImmutableMap creates a new instance of ImmutableMap with the provided map of string keys and IObject values.
func NewImmutableMap(value map[string]IObject) *ImmutableMap {
	return &ImmutableMap{values: value}
}

// Values returns the underlying map of string keys and IObject values.
func (o *ImmutableMap) Values() map[string]IObject {
	return o.values
}

// Length returns the number of elements in the ImmutableMap.
func (o *ImmutableMap) Length() int {
	return len(o.values)
}

func (o *ImmutableMap) SetValue(k string, v IObject) {
	o.values[k] = v
}

func (o *ImmutableMap) GetValue(k string) (IObject, bool) {
	v, ok := o.values[k]
	return v, ok
}

// TypeName returns the type name of the ImmutableMap as a string.
func (o *ImmutableMap) TypeName() string {
	return "immutable-map"
}

// String generates a string representation of the ImmutableMap in key-values pair format.
func (o *ImmutableMap) String() string {
	var pairs []string
	for k, v := range o.values {
		pairs = append(pairs, fmt.Sprintf("%s: %s", k, v.String()))
	}
	return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
}

// Copy creates and returns a deep copy of the ImmutableMap, duplicating all key-values pairs.
func (o *ImmutableMap) Copy() IObject {
	c := make(map[string]IObject)
	for k, v := range o.values {
		c[k] = v.Copy()
	}
	return NewMap(c)
}

// Falsy returns true if the map is empty, indicating it is considered "falsy", otherwise false.
func (o *ImmutableMap) Falsy() bool {
	return len(o.values) == 0
}

// IndexGet retrieves the values associated with the given index in the ImmutableMap. Returns an error for invalid index types.
func (o *ImmutableMap) IndexGet(index IObject) (res IObject, err error) {
	strIdx, ok := ToString(index)
	if !ok {
		err = errors.ErrInvalidIndexType
		return
	}
	res, ok = o.values[strIdx]
	if !ok {
		res = UndefinedValue
	}
	return
}

// Equals determines whether the current ImmutableMap is equal to another IObject by comparing their key-values pairs.
func (o *ImmutableMap) Equals(in IObject) bool {
	var xVal map[string]IObject
	switch x := in.(type) {
	case *Map:
		xVal = x.values
	case *ImmutableMap:
		xVal = x.values
	default:
		return false
	}
	if len(o.values) != len(xVal) {
		return false
	}
	for k, v := range o.values {
		tv := xVal[k]
		if !v.Equals(tv) {
			return false
		}
	}
	return true
}

// Iterate returns an iterator for traversing the key-values pairs in the immutable map.
func (o *ImmutableMap) Iterate() IIterator {
	var keys []string
	for k := range o.values {
		keys = append(keys, k)
	}
	return &MapIterator{
		v: o.values,
		k: keys,
		l: len(keys),
	}
}

// CanIterate returns true, indicating that the ImmutableMap supports iteration.
func (o *ImmutableMap) CanIterate() bool {
	return true
}
