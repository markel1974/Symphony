package objects

import (
	"math"
	"strconv"

	"github.com/markel1974/c64emu/src/kernel/vm/errors"
)

// Float represents a floating-point number type.
// It embeds ObjectImpl to provide default behavior for object operations.
// Value holds the actual float64 value of the Float type.
type Float struct {
	ObjectImpl
	Value float64
}

// NewFloat creates a new Float object with the specified float64 value.
func NewFloat(value float64) *Float {
	return &Float{Value: value}
}

// String returns the string representation of the Float value.
func (o *Float) String() string {
	return strconv.FormatFloat(o.Value, 'f', -1, 64)
}

// TypeName returns the string "float", indicating the type of the Float object.
func (o *Float) TypeName() string {
	return "float"
}

// BinaryOp applies a binary operator to the current Float and a right-hand side Object, returning the resulting Object or an error.
func (o *Float) BinaryOp(op Operator, rhs Object) (Object, error) {
	switch rhs := rhs.(type) {
	case *Float:
		switch op {
		case OperatorAdd:
			r := o.Value + rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Float{Value: r}, nil
		case OperatorSub:
			r := o.Value - rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Float{Value: r}, nil
		case OperatorMul:
			r := o.Value * rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Float{Value: r}, nil
		case OperatorQuo:
			r := o.Value / rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Float{Value: r}, nil
		case OperatorLess:
			if o.Value < rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case OperatorGreater:
			if o.Value > rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case OperatorLessEq:
			if o.Value <= rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case OperatorGreaterEq:
			if o.Value >= rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		default:
			return nil, errors.ErrInvalidOperator
		}
	case *Int:
		switch op {
		case OperatorAdd:
			r := o.Value + float64(rhs.Value)
			if r == o.Value {
				return o, nil
			}
			return &Float{Value: r}, nil
		case OperatorSub:
			r := o.Value - float64(rhs.Value)
			if r == o.Value {
				return o, nil
			}
			return &Float{Value: r}, nil
		case OperatorMul:
			r := o.Value * float64(rhs.Value)
			if r == o.Value {
				return o, nil
			}
			return &Float{Value: r}, nil
		case OperatorQuo:
			r := o.Value / float64(rhs.Value)
			if r == o.Value {
				return o, nil
			}
			return &Float{Value: r}, nil
		case OperatorLess:
			if o.Value < float64(rhs.Value) {
				return TrueValue, nil
			}
			return FalseValue, nil
		case OperatorGreater:
			if o.Value > float64(rhs.Value) {
				return TrueValue, nil
			}
			return FalseValue, nil
		case OperatorLessEq:
			if o.Value <= float64(rhs.Value) {
				return TrueValue, nil
			}
			return FalseValue, nil
		case OperatorGreaterEq:
			if o.Value >= float64(rhs.Value) {
				return TrueValue, nil
			}
			return FalseValue, nil
		default:
			return nil, errors.ErrInvalidOperator
		}
	}
	return nil, errors.ErrInvalidOperator
}

// Copy creates and returns a new instance of Float with the same value as the original.
func (o *Float) Copy() Object {
	return &Float{Value: o.Value}
}

// IsFalsy checks if the Float object's value is NaN, returning true if it is, otherwise false.
func (o *Float) IsFalsy() bool {
	return math.IsNaN(o.Value)
}

// Equals checks if two Float objects have the same value and returns true if they are equal.
func (o *Float) Equals(x Object) bool {
	t, ok := x.(*Float)
	if !ok {
		return false
	}
	return o.Value == t.Value
}

// ToFloat64 converts an Object to a float64 value if possible and returns a boolean indicating success.
func ToFloat64(o Object) (v float64, ok bool) {
	switch o := o.(type) {
	case *Int:
		v = float64(o.Value)
		ok = true
	case *Float:
		v = o.Value
		ok = true
	case *String:
		c, err := strconv.ParseFloat(o.Value, 64)
		if err == nil {
			v = c
			ok = true
		}
	}
	return
}
