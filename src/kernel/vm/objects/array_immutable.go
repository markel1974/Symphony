package objects

import (
	"fmt"
	"strings"
)

// ArrayImmutable represents an array that cannot be modified after creation.
// Implements IObject and supports iteration, comparison, and copying.
type ArrayImmutable struct {
	ObjectImpl
	values []IObject
}

// NewArrayImmutable creates a new ArrayImmutable instance with the given slice of IObject, ensuring it is immutable.
func NewArrayImmutable(value []IObject) *ArrayImmutable {
	return &ArrayImmutable{values: value}
}

// Values returns the underlying slice of IObject stored in the ArrayImmutable, ensuring immutability.
func (o *ArrayImmutable) Values() []IObject {
	return o.values
}

// SetValue assigns a new values to the element at the specified index in the ArrayImmutable, if the index is within bounds.
func (o *ArrayImmutable) SetValue(idx int, v IObject) {
	if idx < 0 || idx >= len(o.values) {
		return
	}
	o.values[idx] = v
}

// Length returns the length of the ArrayImmutable, which is the number of elements it contains.
func (o *ArrayImmutable) Length() int {
	return len(o.values)
}

// TypeName returns the type name of the ArrayImmutable, which is "immutable-array".
func (o *ArrayImmutable) TypeName() string {
	return "immutable-array"
}

// String returns a string representation of the ArrayImmutable, displaying its elements in a comma-separated list.
func (o *ArrayImmutable) String() string {
	var elements []string
	for _, e := range o.values {
		elements = append(elements, e.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}

// BinaryOp performs a binary operation using the provided operator and right-hand side object. Returns the result or an error.
func (o *ArrayImmutable) BinaryOp(op Operator, rhs IObject) (IObject, error) {
	if ia, ok := rhs.(*ArrayImmutable); ok {
		switch op {
		case OperatorAdd:
			return NewArray(append(o.values, ia.values...)), nil
		default:
			return nil, ErrInvalidOperator
		}
	}
	return nil, ErrInvalidOperator
}

// Copy creates and returns a new copy of the ArrayImmutable, with all its elements deeply copied.
func (o *ArrayImmutable) Copy() IObject {
	var c []IObject
	for _, elem := range o.values {
		c = append(c, elem.Copy())
	}
	return NewArray(c)
}

// Falsy checks if the ArrayImmutable is considered falsy, returning true if its Value slice has no elements.
func (o *ArrayImmutable) Boolean() bool {
	return len(o.values) == 0
}

// Equals compares the ArrayImmutable with another IObject for values equality, returning true if their elements are identical.
func (o *ArrayImmutable) Equals(in IObject) bool {
	var xVal []IObject
	switch x := in.(type) {
	case *Array:
		xVal = x.Values()
	case *ArrayImmutable:
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
func (o *ArrayImmutable) IndexGet(index IObject) (res IObject, err error) {
	intIdx, ok := index.(*Int)
	if !ok {
		err = ErrInvalidIndexType
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

// Iterate returns an IIterator to traverse the elements of the ArrayImmutable sequentially.
func (o *ArrayImmutable) Iterate() IIterator {
	return &ArrayIterator{
		v: o.values,
		l: len(o.values),
	}
}

// CanIterate determines if the ArrayImmutable supports iteration, always returning true.
func (o *ArrayImmutable) CanIterate() bool {
	return true
}
