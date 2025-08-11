package objects

import (
	"fmt"
	"strings"

	"github.com/markel1974/c64emu/src/kernel/vm/errors"
)

// Array represents a collection of Object elements, providing methods for manipulation, indexing, and iteration.
type Array struct {
	ObjectImpl
	Value []Object
}

// NewArray creates and returns a new Array object initialized with the provided slice of Object elements.
func NewArray(value []Object) *Array {
	return &Array{Value: value}
}

// TypeName returns the string "array", representing the type name of the Array object.
func (o *Array) TypeName() string {
	return "array"
}

// String returns a string representation of the Array, formatting its elements in a comma-separated list enclosed in brackets.
func (o *Array) String() string {
	var elements []string
	for _, e := range o.Value {
		elements = append(elements, e.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}

// BinaryOp performs a binary operation on the Array with the given operator and right-hand side object.
// It supports the addition operator (OperatorAdd) for concatenating arrays. Returns the result or an error for invalid operations.
func (o *Array) BinaryOp(op Operator, rhs Object) (Object, error) {
	if rhs, ok := rhs.(*Array); ok {
		switch op {
		case OperatorAdd:
			if len(rhs.Value) == 0 {
				return o, nil
			}
			return &Array{Value: append(o.Value, rhs.Value...)}, nil
		default:
			return nil, errors.ErrInvalidOperator
		}
	}
	return nil, errors.ErrInvalidOperator
}

// Copy creates and returns a deep copy of the Array and its elements.
func (o *Array) Copy() Object {
	var c []Object
	for _, elem := range o.Value {
		c = append(c, elem.Copy())
	}
	return &Array{Value: c}
}

// Falsy returns true if the array is empty, otherwise false.
func (o *Array) Falsy() bool {
	return len(o.Value) == 0
}

// Equals compares the current Array with another Object and returns true if they have equivalent values and order.
func (o *Array) Equals(x Object) bool {
	var xVal []Object
	switch x := x.(type) {
	case *Array:
		xVal = x.Value
	case *ImmutableArray:
		xVal = x.Value
	default:
		return false
	}
	if len(o.Value) != len(xVal) {
		return false
	}
	for i, e := range o.Value {
		if !e.Equals(xVal[i]) {
			return false
		}
	}
	return true
}

// IndexGet retrieves the element at the given index from the Array. Returns an error if the index type is invalid or out of bounds.
func (o *Array) IndexGet(index Object) (res Object, err error) {
	intIdx, ok := index.(*Int)
	if !ok {
		err = errors.ErrInvalidIndexType
		return
	}
	idxVal := int(intIdx.Value)
	if idxVal < 0 || idxVal >= len(o.Value) {
		res = UndefinedValue
		return
	}
	res = o.Value[idxVal]
	return
}

// IndexSet assigns a given value to the specified index in the array, returning an error if the operation is invalid.
func (o *Array) IndexSet(index, value Object) (err error) {
	intIdx, ok := ToInt(index)
	if !ok {
		err = errors.ErrInvalidIndexType
		return
	}
	if intIdx < 0 || intIdx >= len(o.Value) {
		err = errors.ErrIndexOutOfBounds
		return
	}
	o.Value[intIdx] = value
	return nil
}

// Iterate returns an Iterator for the Array instance, allowing sequential access to its elements.
func (o *Array) Iterate() Iterator {
	return &ArrayIterator{
		v: o.Value,
		l: len(o.Value),
	}
}

// CanIterate checks if the Array is iterable and always returns true.
func (o *Array) CanIterate() bool {
	return true
}

// ArrayIterator is an iterator type for traversing elements of an array.
// It implements the Iterator interface to provide sequential access to array elements.
type ArrayIterator struct {
	ObjectImpl
	v []Object
	i int
	l int
}

// TypeName returns the type name of the ArrayIterator as a string.
func (i *ArrayIterator) TypeName() string {
	return "array-iterator"
}

// String returns a string representation of the ArrayIterator instance.
func (i *ArrayIterator) String() string {
	return "<array-iterator>"
}

// Falsy determines whether the ArrayIterator should be considered a falsy value. Always returns true.
func (i *ArrayIterator) Falsy() bool {
	return true
}

// Equals checks whether the given Object is equal to the current ArrayIterator instance by value comparison.
func (i *ArrayIterator) Equals(Object) bool {
	return false
}

// Copy creates and returns a duplicate of the ArrayIterator, preserving its current state.
func (i *ArrayIterator) Copy() Object {
	return &ArrayIterator{v: i.v, i: i.i, l: i.l}
}

// Next advances the iterator to the next element and returns true if the current position is within bounds.
func (i *ArrayIterator) Next() bool {
	i.i++
	return i.i <= i.l
}

// Key returns the index of the current element in the iteration as an Object.
func (i *ArrayIterator) Key() Object {
	return &Int{Value: int64(i.i - 1)}
}

// Value returns the current element in the iteration based on the iterator's internal position.
func (i *ArrayIterator) Value() Object {
	return i.v[i.i-1]
}
