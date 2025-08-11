package objects

import (
	"fmt"
	"strings"

	"github.com/markel1974/c64emu/src/kernel/vm/errors"
)

// Array represents an array of objects.
type Array struct {
	ObjectImpl
	Value []Object
}

// TypeName returns the name of the type.
func (o *Array) TypeName() string {
	return "array"
}

func (o *Array) String() string {
	var elements []string
	for _, e := range o.Value {
		elements = append(elements, e.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}

// BinaryOp returns another object that is the result of a given binary
// operator and a right-hand side object.
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

// Copy returns a copy of the type.
func (o *Array) Copy() Object {
	var c []Object
	for _, elem := range o.Value {
		c = append(c, elem.Copy())
	}
	return &Array{Value: c}
}

// IsFalsy returns true if the value of the type is falsy.
func (o *Array) IsFalsy() bool {
	return len(o.Value) == 0
}

// Equals returns true if the value of the type is equal to the value of another object.
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

// IndexGet returns an element at a given index.
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

// IndexSet sets an element at a given index.
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

// Iterate creates an array iterator.
func (o *Array) Iterate() Iterator {
	return &ArrayIterator{
		v: o.Value,
		l: len(o.Value),
	}
}

// CanIterate returns whether the Object can be Iterated.
func (o *Array) CanIterate() bool {
	return true
}

// ArrayIterator is an iterator for an array.
type ArrayIterator struct {
	ObjectImpl
	v []Object
	i int
	l int
}

// TypeName returns the name of the type.
func (i *ArrayIterator) TypeName() string {
	return "array-iterator"
}

func (i *ArrayIterator) String() string {
	return "<array-iterator>"
}

// IsFalsy returns true if the value of the type is falsy.
func (i *ArrayIterator) IsFalsy() bool {
	return true
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (i *ArrayIterator) Equals(Object) bool {
	return false
}

// Copy returns a copy of the type.
func (i *ArrayIterator) Copy() Object {
	return &ArrayIterator{v: i.v, i: i.i, l: i.l}
}

// Next returns true if there are more elements to iterate.
func (i *ArrayIterator) Next() bool {
	i.i++
	return i.i <= i.l
}

// Key returns the key or index value of the current element.
func (i *ArrayIterator) Key() Object {
	return &Int{Value: int64(i.i - 1)}
}

// Value returns the value of the current element.
func (i *ArrayIterator) Value() Object {
	return i.v[i.i-1]
}
