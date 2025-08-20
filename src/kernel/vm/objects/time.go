package objects

import (
	"time"
)

const (
	TimeType = "time"
)

// Time represents a custom object encapsulating a Go time.Time values with extended behaviors and operations.
type Time struct {
	*Object
	value time.Time
}

// NewTime creates a new instance of Time wrapping the provided time.Time values.
func _newTime(factory *Factory, value time.Time) *Time {
	return &Time{
		Object: factory.NewObject(),
		value:  value,
	}
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
	return TimeType
}

// BinaryOp performs a binary operation between the Time object and another IObject using a specified Operator.
// Returns the resulting IObject or an error if the operation is invalid or unsupported.
func (o *Time) BinaryOp(op Operator, in IObject) (IObject, error) {
	switch rhs := in.(type) {
	case *Int:
		switch op {
		case OperatorAdd:
			if rhs.value == 0 {
				return o, nil
			}
			return o.Factory().NewTime(o.value.Add(time.Duration(rhs.value))), nil
		case OperatorSub:
			if rhs.value == 0 {
				return o, nil
			}
			return o.Factory().NewTime(o.value.Add(time.Duration(-rhs.value))), nil
		default:
			return nil, ErrInvalidOperator
		}
	case *Time:
		switch op {
		case OperatorSub:
			return o.Factory().NewInt(int64(o.value.Sub(rhs.value))), nil
		case OperatorLess:
			if o.value.Before(rhs.value) {
				return o.Factory().TrueValue(), nil
			}
			return o.Factory().FalseValue(), nil
		case OperatorGreater:
			if o.value.After(rhs.value) {
				return o.Factory().TrueValue(), nil
			}
			return o.Factory().FalseValue(), nil
		case OperatorLessEq:
			if o.value.Equal(rhs.value) || o.value.Before(rhs.value) {
				return o.Factory().TrueValue(), nil
			}
			return o.Factory().FalseValue(), nil
		case OperatorGreaterEq:
			if o.value.Equal(rhs.value) || o.value.After(rhs.value) {
				return o.Factory().TrueValue(), nil
			}
			return o.Factory().FalseValue(), nil
		default:
			return nil, ErrInvalidOperator
		}
	}
	return nil, ErrInvalidOperator
}

// Copy returns a new instance of the Time object with the same internal time values, duplicating its state.
func (o *Time) Copy() IObject {
	return o.Factory().NewTime(o.value)
}

// Boolean returns true if the Time object's values is zero (indicating it is uninitialized or empty), otherwise false.
func (o *Time) Boolean() bool {
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
