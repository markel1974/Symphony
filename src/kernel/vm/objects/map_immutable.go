package objects

import (
	"fmt"
	"strings"
)

const (
	MapImmutableType = "map_immutable"
)

// MapImmutable represents a read-only map structure where keys are strings and values are of type IObject.
type MapImmutable struct {
	*Object
	values map[string]IObject
}

// NewMapImmutable creates a new instance of MapImmutable with the provided map of string keys and IObject values.
func newMapImmutable(factory *GateKeeper, frame int, value map[string]IObject) *MapImmutable {
	return &MapImmutable{
		Object: factory.NewObject(frame),
		values: value,
	}
}

// Values returns the underlying map of string keys and IObject values.
func (o *MapImmutable) Values() map[string]IObject {
	return o.values
}

// Length returns the number of elements in the MapImmutable.
func (o *MapImmutable) Length() int {
	return len(o.values)
}

func (o *MapImmutable) SetValue(k string, v IObject) {
	o.values[k] = v
}

func (o *MapImmutable) GetValue(k string) (IObject, bool) {
	v, ok := o.values[k]
	return v, ok
}

// TypeName returns the type name of the MapImmutable as a string.
func (o *MapImmutable) TypeName() string {
	return MapImmutableType
}

// String generates a string representation of the MapImmutable in key-values pair format.
func (o *MapImmutable) String() string {
	var pairs []string
	for k, v := range o.values {
		pairs = append(pairs, fmt.Sprintf("%s: %s", k, v.String()))
	}
	return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
}

// Copy creates and returns a deep copy of the MapImmutable, duplicating all key-values pairs.
func (o *MapImmutable) Copy(frame int) IObject {
	c := make(map[string]IObject)
	for k, v := range o.values {
		c[k] = v.Copy(frame)
	}
	return o.Factory().NewMap(frame, c)
}

// Boolean returns true if the map is empty, indicating it is considered "falsy", otherwise false.
func (o *MapImmutable) Boolean() bool {
	return len(o.values) == 0
}

// IndexGet retrieves the values associated with the given index in the MapImmutable. Returns an error for invalid index types.
func (o *MapImmutable) IndexGet(_ int, index IObject) (res IObject, err error) {
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

// Equals determine whether the current MapImmutable is equal to another IObject by comparing their key-values pairs.
func (o *MapImmutable) Equals(in IObject) bool {
	var xVal map[string]IObject
	switch x := in.(type) {
	case *Map:
		xVal = x.values
	case *MapImmutable:
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
func (o *MapImmutable) Iterate(frame int) IIterator {
	return o.Factory().NewMapIterator(frame, o.values)
}

// CanIterate returns true, indicating that the MapImmutable supports iteration.
func (o *MapImmutable) CanIterate() bool {
	return true
}
