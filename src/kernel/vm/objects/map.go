package objects

import (
	"fmt"
	"strings"

	"github.com/markel1974/c64emu/src/kernel/vm/errors"
)

// Map represents a collection of key-value pairs where keys are strings and values implement the IObject interface.
type Map struct {
	ObjectImpl
	values map[string]IObject
}

// NewMap creates and returns a new instance of Map initialized with the provided map of string keys to IObject values.
func NewMap(value map[string]IObject) *Map {
	return &Map{values: value}
}

// Get retrieves the value associated with the specified key from the map. If the key is not found, it returns nil.
func (o *Map) Get(key string) IObject {
	return o.values[key]
}

// Set assigns the specified value to the given key in the Map. Overrides the value if the key already exists.
func (o *Map) Set(key string, value IObject) {
	o.values[key] = value
}

// Delete removes the entry associated with the specified key from the map.
func (o *Map) Delete(key string) {
	delete(o.values, key)
}

// Has checks if the specified key exists in the Map and returns true if found, otherwise false.
func (o *Map) Has(key string) bool {
	_, ok := o.values[key]
	return ok
}

// Values returns the internal map of key-value pairs stored in the Map object.
func (o *Map) Values() map[string]IObject {
	return o.values
}

// Length returns the number of key-value pairs in the Map object.
func (o *Map) Length() int {
	return len(o.values)
}

// TypeName returns the string "map", which represents the type name of the Map object.
func (o *Map) TypeName() string {
	return "map"
}

// String returns the string representation of the Map object in the format of key-value pairs enclosed in braces.
func (o *Map) String() string {
	var pairs []string
	for k, v := range o.values {
		pairs = append(pairs, fmt.Sprintf("%s: %s", k, v.String()))
	}
	return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
}

// Copy creates and returns a deep copy of the Map object, duplicating all key-value pairs recursively.
func (o *Map) Copy() IObject {
	c := make(map[string]IObject)
	for k, v := range o.values {
		c[k] = v.Copy()
	}
	return &Map{values: c}
}

// Falsy returns true if the map contains no key-value pairs, indicating it is empty.
func (o *Map) Falsy() bool {
	return len(o.values) == 0
}

// Equals checks if the Map is equal to another IObject by comparing their key-value pairs. Returns true if equal.
func (o *Map) Equals(in IObject) bool {
	var xVal map[string]IObject
	switch x := in.(type) {
	case *Map:
		xVal = x.values
	case *ImmutableMap:
		xVal = x.Value
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

// IndexGet retrieves the value associated with the given index in the map. Returns UndefinedValue if the index does not exist.
// An error is returned if the index type is invalid.
func (o *Map) IndexGet(index IObject) (res IObject, err error) {
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

// IndexSet sets the specified value at the given string-convertible index in the Map. Returns an error for invalid index types.
func (o *Map) IndexSet(index, value IObject) (err error) {
	strIdx, ok := ToString(index)
	if !ok {
		err = errors.ErrInvalidIndexType
		return
	}
	o.values[strIdx] = value
	return nil
}

// Iterate creates and returns an iterator for the Map, allowing iteration over its keys and associated values.
func (o *Map) Iterate() IIterator {
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

// CanIterate returns true, indicating that the Map object supports iteration over its elements.
func (o *Map) CanIterate() bool {
	return true
}

// MapIterator is a struct used to iterate over a map's key-value pairs.
// It implements the IIterator interface for sequential traversal.
// ObjectImpl is embedded to provide default behaviors defined by the IObject interface.
// v is the map being iterated over.
// k stores the keys of the map in a slice for ordered access.
// i tracks the current iteration index.
// l holds the length of the keys slice.
type MapIterator struct {
	ObjectImpl
	v map[string]IObject
	k []string
	i int
	l int
}

// TypeName returns the type name of the MapIterator as a string.
func (i *MapIterator) TypeName() string {
	return "map-iterator"
}

// String returns the string representation of the MapIterator.
func (i *MapIterator) String() string {
	return "<map-iterator>"
}

// Falsy returns true, indicating the MapIterator object is considered falsy in a boolean context.
func (i *MapIterator) Falsy() bool {
	return true
}

// Equals checks if the given IObject is equal to the current MapIterator instance.
func (i *MapIterator) Equals(IObject) bool {
	return false
}

// Copy creates and returns a duplicate of the MapIterator with the same internal state.
func (i *MapIterator) Copy() IObject {
	return &MapIterator{v: i.v, k: i.k, i: i.i, l: i.l}
}

// Next advances the iterator to the next element and returns true if the current position is within range.
func (i *MapIterator) Next() bool {
	i.i++
	return i.i <= i.l
}

// Key returns the key of the current element in the iteration as an IObject.
func (i *MapIterator) Key() IObject {
	k := i.k[i.i-1]
	return NewString(k)
}

// Value retrieves the value associated with the current key in the map iterator.
func (i *MapIterator) Value() IObject {
	k := i.k[i.i-1]
	return i.v[k]
}

// ToMap converts an IObject to a Go native map[string]interface{} representation if the object is a Map type.
func ToMap(o IObject) (res map[string]interface{}) {
	switch o := o.(type) {
	case *Map:
		res = make(map[string]interface{})
		for key, v := range o.values {
			res[key] = ToInterface(v)
		}
	}
	return
}
