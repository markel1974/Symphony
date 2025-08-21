package objects

import (
	"fmt"
	"strings"
)

// StructType represents the constant string identifier for a struct type in the application context.
const (
	StructType = "struct"
)

// Struct is a composite object that implements the IObject interface and stores a collection of key-value pairs.
type Struct struct {
	*Object
	values map[string]IObject
}

// NewStruct creates a new instance of MapImmutable with the provided map of string keys and IObject values.
func newStruct(factory *GateKeeper, frame int, value map[string]IObject) *Struct {
	return &Struct{
		Object: factory.NewObject(frame),
		values: value,
	}
}

// Values returns the underlying map of string keys to IObject values contained within the Struct.
func (o *Struct) Values() map[string]IObject {
	return o.values
}

// Length returns the number of key-value pairs stored in the Struct.
func (o *Struct) Length() int {
	return len(o.values)
}

// SetValue sets the specified key to the given value in the values map of the Struct.
func (o *Struct) SetValue(k string, v IObject) {
	o.values[k] = v
}

// GetValue retrieves the value associated with the given key in the values map and a boolean indicating its presence.
func (o *Struct) GetValue(k string) (IObject, bool) {
	v, ok := o.values[k]
	return v, ok
}

// TypeName returns the type name of the object as a string.
func (o *Struct) TypeName() string {
	return StructType
}

// String returns a string representation of the Struct object, formatting its key-value pairs into a map-like structure.
func (o *Struct) String() string {
	var pairs []string
	for k, v := range o.values {
		pairs = append(pairs, fmt.Sprintf("%s: %s", k, v.String()))
	}
	return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
}

// Copy creates and returns a new IObject by duplicating the internal state of the Struct instance.
func (o *Struct) Copy(frame int) IObject {
	c := make(map[string]IObject)
	for k, v := range o.values {
		c[k] = v.Copy(frame)
	}
	return o.GateKeeper().NewStruct(frame, c)
}

// Boolean returns true if the Struct contains no values, otherwise false.
func (o *Struct) Boolean() bool {
	return len(o.values) == 0
}

// IndexGet retrieves the value associated with the given index within the Struct. Returns an error for invalid index types.
func (o *Struct) IndexGet(_ int, index IObject) (res IObject, err error) {
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

// Equals checks if the current Struct is equal to another IObject by comparing their key-value pairs and lengths.
func (o *Struct) Equals(in IObject) bool {
	var xVal map[string]IObject
	switch x := in.(type) {
	case *Struct:
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

// Iterate returns an IIterator for traversing the key-value pairs in the Struct's internal map.
func (o *Struct) Iterate(frame int) IIterator {
	return o.GateKeeper().NewStructIterator(frame, o.values)
}

// CanIterate checks if the object can be iterated over. Always returns true for this implementation.
func (o *Struct) CanIterate() bool {
	return true
}
