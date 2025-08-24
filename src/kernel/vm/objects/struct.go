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
	factory IGateKeeper
	frame   int
	values  map[string]IObject
}

// NewStruct creates a new instance of MapImmutable with the provided map of string keys and IObject values.
func newStruct(factory IGateKeeper, frame int, value map[string]IObject) IObject {
	if len(value) > maxStructLen {
		nv := make(map[string]IObject)
		for k, v := range value {
			if len(nv) > maxStructLen {
				break
			}
			nv[k] = v
		}
		value = nv
	}
	return &Struct{
		factory: factory,
		frame:   frame,
		values:  value,
	}
}

// GateKeeper returns a reference to the GateKeeper associated with the Object.
func (o *Struct) GateKeeper() IGateKeeper {
	return o.factory
}

// Frame returns the current frame value of the Object.
func (o *Struct) Frame() int {
	return o.frame
}

// BinaryOp performs a binary operation on the current object and another object using the specified operator.
// Returns the result of the operation or an error if the operation is not supported.
func (o *Struct) BinaryOp(_ int, _ Operator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *Struct) Call(_ int, _ ...IObject) (ret IObject, err error) {
	return nil, nil
}

// CanCall determines if the object can be invoked as a callable. Returns false for non-callable objects.
func (o *Struct) CanCall() bool {
	return false
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
	if len(o.values) > maxStructLen {
		return
	}
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
func (o *Struct) Copy(frame int, depth int) IObject {
	c := make(map[string]IObject)
	for k, v := range o.values {
		if depth >= maxDepth {
			break
		}
		c[k] = v.Copy(frame, depth+1)
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

// IndexSet updates or assigns a value to the specified index within the Struct. Returns an error for invalid index types.
func (o *Struct) IndexSet(index, value IObject) (err error) {
	strIdx, ok := o.GateKeeper().ToString(index)
	if !ok {
		err = ErrInvalidIndexType
		return
	}
	o.values[strIdx] = value
	return nil
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
	return o.GateKeeper().NewStructIterator(frame, o.values, 0)
}

// CanIterate checks if the object can be iterated over. Always returns true for this implementation.
func (o *Struct) CanIterate() bool {
	return true
}
