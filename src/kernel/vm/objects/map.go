package objects

import (
	"fmt"
	"strings"
)

const (
	MapType = "map"
)

// Map represents a collection of key-values pairs where keys are strings and values implement the IObject interface.
type Map struct {
	*Object
	values map[string]IObject
}

// NewMap creates and returns a new instance of Map initialized with the provided map of string keys to IObject values.
func newMap(factory *GateKeeper, frame int, value map[string]IObject) *Map {
	return &Map{
		Object: factory.NewObject(frame),
		values: value,
	}
}

// Get retrieves the values associated with the specified key from the map. If the key is not found, it returns nil.
func (o *Map) Get(key string) IObject {
	return o.values[key]
}

// Set assigns the specified values to the given key in the Map. Overrides the values if the key already exists.
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

// Values returns the internal map of key-values pairs stored in the Map object.
func (o *Map) Values() map[string]IObject {
	return o.values
}

// Length returns the number of key-values pairs in the Map object.
func (o *Map) Length() int {
	return len(o.values)
}

// TypeName returns the string "map", which represents the type name of the Map object.
func (o *Map) TypeName() string {
	return MapType
}

// String returns the string representation of the Map object in the format of key-values pairs enclosed in braces.
func (o *Map) String() string {
	var pairs []string
	for k, v := range o.values {
		pairs = append(pairs, fmt.Sprintf("%s: %s", k, v.String()))
	}
	return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
}

// Copy creates and returns a deep copy of the Map object, duplicating all key-values pairs recursively.
func (o *Map) Copy(frame int) IObject {
	c := make(map[string]IObject)
	for k, v := range o.values {
		c[k] = v.Copy(frame)
	}
	return o.Factory().NewMap(frame, c)
}

// Boolean returns true if the map contains no key-values pairs, indicating it is empty.
func (o *Map) Boolean() bool {
	return len(o.values) == 0
}

// Equals checks if the Map is equal to another IObject by comparing their key-values pairs. Returns true if equal.
func (o *Map) Equals(in IObject) bool {
	var xVal map[string]IObject
	switch x := in.(type) {
	case *Map:
		xVal = x.values
	case *MapImmutable:
		xVal = x.Values()
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

// IndexGet retrieves the values associated with the given index in the map. Returns UndefinedValue if the index does not exist.
// An error is returned if the index type is invalid.
func (o *Map) IndexGet(_ int, index IObject) (res IObject, err error) {
	strIdx, ok := o.Factory().ToString(index)
	if !ok {
		err = ErrInvalidIndexType
		return
	}
	res, ok = o.values[strIdx]
	if !ok {
		res = o.Factory().UndefinedValue()
	}
	return
}

// IndexSet sets the specified values at the given string-convertible index in the Map. Returns an error for invalid index types.
func (o *Map) IndexSet(index, value IObject) (err error) {
	strIdx, ok := o.Factory().ToString(index)
	if !ok {
		err = ErrInvalidIndexType
		return
	}
	o.values[strIdx] = value
	return nil
}

// Iterate creates and returns an iterator for the Map, allowing iteration over its keys and associated values.
func (o *Map) Iterate(frame int) IIterator {
	return o.Factory().NewMapIterator(frame, o.values)
}

// CanIterate returns true, indicating that the Map object supports iteration over its elements.
func (o *Map) CanIterate() bool {
	return true
}
