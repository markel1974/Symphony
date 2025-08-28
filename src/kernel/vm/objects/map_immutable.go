package objects

import (
	"encoding/gob"
	"fmt"
	"strings"
)

const (
	MapImmutableType = "map_immutable"
)

func init() {
	gob.Register(&MapImmutable{})
}

// MapImmutable represents a read-only map structure where keys are strings and values are of type IObject.
type MapImmutable struct {
	factory IGateKeeper
	frame   int
	values  map[string]IObject
}

// NewMapImmutable creates a new instance of MapImmutable with the provided map of string keys and IObject values.
func newMapImmutable(factory IGateKeeper, frame int, value map[string]IObject) IObject {
	return &MapImmutable{
		factory: factory,
		frame:   frame,
		values:  value,
	}
}

// GateKeeper returns a reference to the GateKeeper associated with the Object.
func (o *MapImmutable) GateKeeper() IGateKeeper {
	return o.factory
}

// Frame returns the current frame value of the Object.
func (o *MapImmutable) Frame() int {
	return o.frame
}

// BinaryOp performs a binary operation on the current object and another object using the specified operator.
// Returns the result of the operation or an error if the operation is not supported.
func (o *MapImmutable) BinaryOp(_ int, _ Operator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// IndexSet attempts to assign a value to an index in the object but always returns ErrNotIndexAssignable,
// as this operation is unsupported.
func (o *MapImmutable) IndexSet(_, _ IObject) (err error) {
	return ErrNotIndexAssignable
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *MapImmutable) Call(_ int, _ ...IObject) (ret IObject, err error) {
	return nil, nil
}

// CanCall determines if the object can be invoked as a callable. Returns false for non-callable objects.
func (o *MapImmutable) CanCall() bool {
	return false
}

// Values returns the underlying map of string keys and IObject values.
func (o *MapImmutable) Values() map[string]IObject {
	return o.values
}

// Length returns the number of elements in the MapImmutable.
func (o *MapImmutable) Length() int {
	return len(o.values)
}

// SetValue assigns the given value to the specified key in the MapImmutable.
func (o *MapImmutable) SetValue(k string, v IObject) {
	o.values[k] = v
}

// GetValue retrieves the value associated with the specified key from the MapImmutable.
// Returns the value and a boolean indicating if the key exists.
func (o *MapImmutable) GetValue(k string) (IObject, bool) {
	v, ok := o.values[k]
	return v, ok
}

// TypeName returns the string representation of the type of the MapImmutable object.
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
func (o *MapImmutable) Copy(frame int, depth int) IObject {
	c := make(map[string]IObject)
	for k, v := range o.values {
		if depth >= maxDepth {
			break
		}
		c[k] = v.Copy(frame, depth+1)
	}
	return o.GateKeeper().NewMap(frame, c)
}

// Falsy returns true if the map is empty, indicating it is considered "falsy", otherwise false.
func (o *MapImmutable) Falsy() bool {
	return len(o.values) == 0
}

// IndexGet retrieves the values associated with the given index in the MapImmutable. Returns an error for invalid index types.
func (o *MapImmutable) IndexGet(_ int, index IObject) (res IObject, err error) {
	strIdx, ok := o.GateKeeper().ToString(index)
	if !ok {
		err = ErrInvalidIndexType
		return
	}
	res, ok = o.values[strIdx]
	if !ok {
		res = o.GateKeeper().UndefinedValue()
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
	return o.GateKeeper().NewMapIterator(frame, o.values, 0)
}

// CanIterate returns true, indicating that the MapImmutable supports iteration.
func (o *MapImmutable) CanIterate() bool {
	return true
}
