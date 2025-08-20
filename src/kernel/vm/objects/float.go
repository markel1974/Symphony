package objects

import (
	"math"
	"strconv"
)

const (
	FloatType = "float"
)

// Float represents a floating-point number and provides operations and behaviors specific to numeric types.
// It embeds Object to implement common interface methods and extends behavior where necessary.
// The value field holds the actual float64 values encapsulated by the Float type.
type Float struct {
	*Object
	value float64
}

// NewFloat creates and returns a pointer to a new Float object initialized with the specified float64 values.
func newFloat(factory *Factory, frame int, value float64) *Float {
	return &Float{
		Object: factory.NewObject(frame),
		value:  value,
	}
}

func (o *Float) Value() float64 {
	return o.value
}

// String returns the string representation of the Float object using its internal float64 values.
func (o *Float) String() string {
	return strconv.FormatFloat(o.value, 'f', -1, 64)
}

// TypeName returns the name of the type.
func (o *Float) TypeName() string {
	return FloatType
}

// BinaryOp performs a binary operation between the current Float and another IObject based on the specified operator.
// Returns the result of the operation as an IObject or an error if the operation is invalid.
func (o *Float) BinaryOp(frame int, op Operator, rhs IObject) (IObject, error) {
	switch rhs := rhs.(type) {
	case *Float:
		switch op {
		case OperatorAdd:
			r := o.value + rhs.value
			if r == o.value {
				return o, nil
			}
			return o.Factory().NewFloat(frame, r), nil
		case OperatorSub:
			r := o.value - rhs.value
			if r == o.value {
				return o, nil
			}
			return o.Factory().NewFloat(frame, r), nil
		case OperatorMul:
			r := o.value * rhs.value
			if r == o.value {
				return o, nil
			}
			return o.Factory().NewFloat(frame, r), nil
		case OperatorQuo:
			r := o.value / rhs.value
			if r == o.value {
				return o, nil
			}
			return o.Factory().NewFloat(frame, r), nil
		case OperatorLess:
			if o.value < rhs.value {
				return o.Factory().TrueValue(), nil
			}
			return o.Factory().FalseValue(), nil
		case OperatorGreater:
			if o.value > rhs.value {
				return o.Factory().TrueValue(), nil
			}
			return o.Factory().FalseValue(), nil
		case OperatorLessEq:
			if o.value <= rhs.value {
				return o.Factory().TrueValue(), nil
			}
			return o.Factory().FalseValue(), nil
		case OperatorGreaterEq:
			if o.value >= rhs.value {
				return o.Factory().TrueValue(), nil
			}
			return o.Factory().FalseValue(), nil
		default:
			return nil, ErrInvalidOperator
		}
	case *Int:
		switch op {
		case OperatorAdd:
			r := o.value + float64(rhs.value)
			if r == o.value {
				return o, nil
			}
			return o.Factory().NewFloat(frame, r), nil
		case OperatorSub:
			r := o.value - float64(rhs.value)
			if r == o.value {
				return o, nil
			}
			return o.Factory().NewFloat(frame, r), nil
		case OperatorMul:
			r := o.value * float64(rhs.value)
			if r == o.value {
				return o, nil
			}
			return o.Factory().NewFloat(frame, r), nil
		case OperatorQuo:
			r := o.value / float64(rhs.value)
			if r == o.value {
				return o, nil
			}
			return o.Factory().NewFloat(frame, r), nil
		case OperatorLess:
			if o.value < float64(rhs.value) {
				return o.Factory().TrueValue(), nil
			}
			return o.Factory().FalseValue(), nil
		case OperatorGreater:
			if o.value > float64(rhs.value) {
				return o.Factory().TrueValue(), nil
			}
			return o.Factory().FalseValue(), nil
		case OperatorLessEq:
			if o.value <= float64(rhs.value) {
				return o.Factory().TrueValue(), nil
			}
			return o.Factory().FalseValue(), nil
		case OperatorGreaterEq:
			if o.value >= float64(rhs.value) {
				return o.Factory().TrueValue(), nil
			}
			return o.Factory().FalseValue(), nil
		default:
			return nil, ErrInvalidOperator
		}
	}
	return nil, ErrInvalidOperator
}

// Copy creates and returns a new instance of the Float object, duplicating its current state.
func (o *Float) Copy(frame int) IObject {
	return o.Factory().NewFloat(frame, o.value)
}

// Boolean determines if the float object is considered falsy, returning true if the values is NaN; otherwise, false.
func (o *Float) Boolean() bool {
	return math.IsNaN(o.value)
}

// Equals checks if the current Float object is equal to another IObject by comparing their internal float64 values.
func (o *Float) Equals(x IObject) bool {
	t, ok := x.(*Float)
	if !ok {
		return false
	}
	return o.value == t.value
}
