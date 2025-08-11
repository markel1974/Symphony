package objects

import (
	"github.com/markel1974/c64emu/src/kernel/vm/errors"
)

// Char represents a character type, encapsulating a single rune value and inheriting behavior from ObjectImpl.
type Char struct {
	ObjectImpl
	value rune
}

// NewChar creates and returns a new Char object with the specified rune value.
func NewChar(value rune) *Char {
	return &Char{value: value}
}

// Value returns the rune value stored in the Char object.
func (o *Char) Value() rune {
	return o.value
}

// String returns the string representation of the Char object's value.
func (o *Char) String() string {
	return string(o.value)
}

// TypeName returns the name of the type as a string, which is "char" for the Char type.
func (o *Char) TypeName() string {
	return "char"
}

// BinaryOp performs a binary operation between the Char object and another IObject using the specified operator.
func (o *Char) BinaryOp(op Operator, in IObject) (IObject, error) {
	switch rhs := in.(type) {
	case *Char:
		switch op {
		case OperatorAdd:
			r := o.value + rhs.value
			if r == o.value {
				return o, nil
			}
			return NewChar(r), nil
		case OperatorSub:
			r := o.value - rhs.value
			if r == o.value {
				return o, nil
			}
			return NewChar(r), nil
		case OperatorLess:
			if o.value < rhs.value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case OperatorGreater:
			if o.value > rhs.value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case OperatorLessEq:
			if o.value <= rhs.value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case OperatorGreaterEq:
			if o.value >= rhs.value {
				return TrueValue, nil
			}
			return FalseValue, nil
		default:
			return nil, errors.ErrInvalidOperator
		}
	case *Int:
		switch op {
		case OperatorAdd:
			r := o.value + rune(rhs.value)
			if r == o.value {
				return o, nil
			}
			return NewChar(r), nil
		case OperatorSub:
			r := o.value - rune(rhs.value)
			if r == o.value {
				return o, nil
			}
			return NewChar(r), nil
		case OperatorLess:
			if int64(o.Value()) < rhs.value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case OperatorGreater:
			if int64(o.Value()) > rhs.value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case OperatorLessEq:
			if int64(o.Value()) <= rhs.value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case OperatorGreaterEq:
			if int64(o.Value()) >= rhs.value {
				return TrueValue, nil
			}
			return FalseValue, nil
		default:
			return nil, errors.ErrInvalidOperator
		}
	}
	return nil, errors.ErrInvalidOperator
}

// Copy creates and returns a new instance of the Char object with the same value.
func (o *Char) Copy() IObject {
	return &Char{value: o.value}
}

// Falsy checks whether the Char object represents a falsy state, returning true if the underlying value is 0.
func (o *Char) Falsy() bool {
	return o.value == 0
}

// Equals checks if the current Char object is equal to another IObject. Returns true if both objects are equal.
func (o *Char) Equals(x IObject) bool {
	t, ok := x.(*Char)
	if !ok {
		return false
	}
	return o.value == t.value
}
