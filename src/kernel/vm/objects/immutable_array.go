package objects

import (
	"fmt"
	"strings"

	"github.com/markel1974/c64emu/src/kernel/vm/errors"
)

// ImmutableArray represents an array that cannot be modified after creation.
// Implements IObject and supports iteration, comparison, and copying.
type ImmutableArray struct {
	ObjectImpl
	values []IObject
}

// NewImmutableArray creates a new ImmutableArray instance with the given slice of IObject, ensuring it is immutable.
func NewImmutableArray(value []IObject) *ImmutableArray {
	return &ImmutableArray{values: value}
}

// Values returns the underlying slice of IObject stored in the ImmutableArray, ensuring immutability.
func (o *ImmutableArray) Values() []IObject {
	return o.values
}

// SetValue assigns a new values to the element at the specified index in the ImmutableArray, if the index is within bounds.
func (o *ImmutableArray) SetValue(idx int, v IObject) {
	if idx < 0 || idx >= len(o.values) {
		return
	}
	o.values[idx] = v
}

// Length returns the length of the ImmutableArray, which is the number of elements it contains.
func (o *ImmutableArray) Length() int {
	return len(o.values)
}

// TypeName returns the type name of the ImmutableArray, which is "immutable-array".
func (o *ImmutableArray) TypeName() string {
	return "immutable-array"
}

// String returns a string representation of the ImmutableArray, displaying its elements in a comma-separated list.
func (o *ImmutableArray) String() string {
	var elements []string
	for _, e := range o.values {
		elements = append(elements, e.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}

// BinaryOp performs a binary operation using the provided operator and right-hand side object. Returns the result or an error.
func (o *ImmutableArray) BinaryOp(op Operator, rhs IObject) (IObject, error) {
	if ia, ok := rhs.(*ImmutableArray); ok {
		switch op {
		case OperatorAdd:
			return NewArray(append(o.values, ia.values...)), nil
		default:
			return nil, errors.ErrInvalidOperator
		}
	}
	return nil, errors.ErrInvalidOperator
}

// Copy creates and returns a new copy of the ImmutableArray, with all its elements deeply copied.
func (o *ImmutableArray) Copy() IObject {
	var c []IObject
	for _, elem := range o.values {
		c = append(c, elem.Copy())
	}
	return NewArray(c)
}

// Falsy checks if the ImmutableArray is considered falsy, returning true if its Value slice has no elements.
func (o *ImmutableArray) Falsy() bool {
	return len(o.values) == 0
}

// Equals compares the ImmutableArray with another IObject for values equality, returning true if their elements are identical.
func (o *ImmutableArray) Equals(in IObject) bool {
	var xVal []IObject
	switch x := in.(type) {
	case *Array:
		xVal = x.Values()
	case *ImmutableArray:
		xVal = x.values
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

// IndexGet retrieves an element from the array at the specified index. Returns error for invalid index type or out of bounds.
func (o *ImmutableArray) IndexGet(index IObject) (res IObject, err error) {
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

// Iterate returns an IIterator to traverse the elements of the ImmutableArray sequentially.
func (o *ImmutableArray) Iterate() IIterator {
	return &ArrayIterator{
		v: o.values,
		l: len(o.values),
	}
}

// CanIterate determines if the ImmutableArray supports iteration, always returning true.
func (o *ImmutableArray) CanIterate() bool {
	return true
}
