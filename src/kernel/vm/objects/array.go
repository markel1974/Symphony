package objects

import (
	"fmt"
	"strings"

	"github.com/markel1974/c64emu/src/kernel/vm/errors"
)

// Array represents a collection of IObject elements, providing methods for manipulation, indexing, and iteration.
type Array struct {
	ObjectImpl
	values []IObject
}

// NewArray creates and returns a new Array object initialized with the provided slice of IObject elements.
func NewArray(value []IObject) *Array {
	return &Array{values: value}
}

// TypeName returns the string "array", representing the type name of the Array object.
func (o *Array) TypeName() string {
	return "array"
}

// Values returns the slice of IObject elements stored in the Array.
func (o *Array) Values() []IObject {
	return o.values
}

// Length returns the number of elements in the Array.
func (o *Array) Length() int {
	return len(o.values)
}

// SetValue assigns a given IObject values to the specified index in the Array if the index is within bounds. No action otherwise.
func (o *Array) SetValue(idx int, value IObject) {
	if idx < 0 || idx >= len(o.values) {
		return
	}
	o.values[idx] = value
}

// Append adds an element to the end of the Array.
func (o *Array) Append(elem IObject) {
	o.values = append(o.values, elem)
}

// Assign replaces the current slice of elements with the provided slice.
func (o *Array) Assign(v []IObject) {
	o.values = v
}

// String returns a string representation of the Array, formatting its elements in a comma-separated list enclosed in brackets.
func (o *Array) String() string {
	var elements []string
	for _, e := range o.values {
		elements = append(elements, e.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}

// BinaryOp performs a binary operation on the Array with the given operator and right-hand side object.
// It supports the addition operator (OperatorAdd) for concatenating arrays. Returns the result or an error for invalid operations.
func (o *Array) BinaryOp(op Operator, in IObject) (IObject, error) {
	if rhs, ok := in.(*Array); ok {
		switch op {
		case OperatorAdd:
			if len(rhs.values) == 0 {
				return o, nil
			}
			return &Array{values: append(o.values, rhs.values...)}, nil
		default:
			return nil, errors.ErrInvalidOperator
		}
	}
	return nil, errors.ErrInvalidOperator
}

// Copy creates and returns a deep copy of the Array and its elements.
func (o *Array) Copy() IObject {
	var c []IObject
	for _, elem := range o.values {
		c = append(c, elem.Copy())
	}
	return &Array{values: c}
}

// Falsy returns true if the array is empty, otherwise false.
func (o *Array) Falsy() bool {
	return len(o.values) == 0
}

// Equals compares the current Array with another IObject and returns true if they have equivalent values and order.
func (o *Array) Equals(in IObject) bool {
	var xVal []IObject
	switch x := in.(type) {
	case *Array:
		xVal = x.values
	case *ImmutableArray:
		xVal = x.Values()
	default:
		return false
	}
	if len(o.values) != len(xVal) {
		return false
	}
	for i, e := range o.values {
		if !e.Equals(xVal[i]) {
			return false
		}
	}
	return true
}

// IndexGet retrieves the element at the given index from the Array. Returns an error if the index type is invalid or out of bounds.
func (o *Array) IndexGet(index IObject) (res IObject, err error) {
	intIdx, ok := index.(*Int)
	if !ok {
		err = errors.ErrInvalidIndexType
		return
	}
	idxVal := int(intIdx.value)
	if idxVal < 0 || idxVal >= len(o.values) {
		res = UndefinedValue
		return
	}
	res = o.values[idxVal]
	return
}

// IndexSet assigns a given values to the specified index in the array, returning an error if the operation is invalid.
func (o *Array) IndexSet(index, value IObject) (err error) {
	intIdx, ok := ToInt(index)
	if !ok {
		err = errors.ErrInvalidIndexType
		return
	}
	if intIdx < 0 || intIdx >= len(o.values) {
		err = errors.ErrIndexOutOfBounds
		return
	}
	o.values[intIdx] = value
	return nil
}

// Iterate returns an IIterator for the Array instance, allowing sequential access to its elements.
func (o *Array) Iterate() IIterator {
	return NewArrayIterator(o.values)
}

// CanIterate checks if the Array is iterable and always returns true.
func (o *Array) CanIterate() bool {
	return true
}
