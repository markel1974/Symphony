package objects

import (
	"time"

	"github.com/markel1974/c64emu/src/kernel/vm/errors"
)

// Time represents a time values.
type Time struct {
	ObjectImpl
	value time.Time
}

func NewTime(value time.Time) *Time {
	return &Time{value: value}
}

func (o *Time) String() string {
	return o.value.String()
}

// TypeName returns the name of the type.
func (o *Time) TypeName() string {
	return "time"
}

// BinaryOp returns another object that is the result of a given binary
// operator and a right-hand side object.
func (o *Time) BinaryOp(op Operator, rhs IObject) (IObject, error) {
	switch rhs := rhs.(type) {
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

// Copy returns a copy of the type.
func (o *Time) Copy() IObject {
	return &Time{value: o.value}
}

// IsFalsy returns true if the values of the type is falsy.
func (o *Time) Falsy() bool {
	return o.value.IsZero()
}

// Equals returns true if the values of the type is equal to the values of
// another object.
func (o *Time) Equals(x IObject) bool {
	t, ok := x.(*Time)
	if !ok {
		return false
	}
	return o.value.Equal(t.value)
}

// ToTime will try to convert object o to time.Time values.
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
