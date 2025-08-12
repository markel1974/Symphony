package objects

import (
	"time"

	"github.com/markel1974/c64emu/src/kernel/vm/errors"
)

// Time represents a custom object encapsulating a Go time.Time values with extended behaviors and operations.
type Time struct {
	ObjectImpl
	value time.Time
}

// NewTime creates a new instance of Time wrapping the provided time.Time values.
func NewTime(value time.Time) *Time {
	return &Time{value: value}
}

// Value returns the underlying time.Time values of the Time object.
func (o *Time) Value() time.Time {
	return o.value
}

// String returns the string representation of the Time object by delegating to the underlying time.Time values.
func (o *Time) String() string {
	return o.value.String()
}

// TypeName returns the name of the type as a string, which is "time".
func (o *Time) TypeName() string {
	return "time"
}

// BinaryOp performs a binary operation between the Time object and another IObject using a specified Operator.
// Returns the resulting IObject or an error if the operation is invalid or unsupported.
func (o *Time) BinaryOp(op Operator, in IObject) (IObject, error) {
	switch rhs := in.(type) {
	case *Int:
		switch op {
		case OperatorAdd: // time + int => time
			if rhs.value == 0 {
				return o, nil
			}
			return &Time{value: o.value.Add(time.Duration(rhs.value))}, nil
		case OperatorSub: // time - int => time
			if rhs.value == 0 {
				return o, nil
			}
			return &Time{value: o.value.Add(time.Duration(-rhs.value))}, nil
		default:
			return nil, errors.ErrInvalidOperator
		}
	case *Time:
		switch op {
		case OperatorSub: // time - time => int (duration)
			return &Int{value: int64(o.value.Sub(rhs.value))}, nil
		case OperatorLess: // time < time => bool
			if o.value.Before(rhs.value) {
				return TrueValue, nil
			}
			return FalseValue, nil
		case OperatorGreater:
			if o.value.After(rhs.value) {
				return TrueValue, nil
			}
			return FalseValue, nil
		case OperatorLessEq:
			if o.value.Equal(rhs.value) || o.value.Before(rhs.value) {
				return TrueValue, nil
			}
			return FalseValue, nil
		case OperatorGreaterEq:
			if o.value.Equal(rhs.value) || o.value.After(rhs.value) {
				return TrueValue, nil
			}
			return FalseValue, nil
		default:
			return nil, errors.ErrInvalidOperator
		}
	}
	return nil, errors.ErrInvalidOperator
}

// Copy returns a new instance of the Time object with the same internal time values, duplicating its state.
func (o *Time) Copy() IObject {
	return &Time{value: o.value}
}

// Falsy returns true if the Time object's values is zero (indicating it is uninitialized or empty), otherwise false.
func (o *Time) Falsy() bool {
	return o.value.IsZero()
}

// Equals checks whether the Time object is equal to another object of type IObject, returning true if they match.
func (o *Time) Equals(x IObject) bool {
	t, ok := x.(*Time)
	if !ok {
		return false
	}
	return o.value.Equal(t.value)
}

// ToTime converts an IObject into a time.Time if it is time-compatible (e.g., *Time or *Int). Returns the time and a boolean.
func ToTime(o IObject) (v time.Time, ok bool) {
	switch o := o.(type) {
	case *Time:
		v = o.value
		ok = true
	case *Int:
		v = time.Unix(o.value, 0)
		ok = true
	}
	return
}
