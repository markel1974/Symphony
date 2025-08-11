package objects

import (
	"math"
	"strconv"

	"github.com/markel1974/c64emu/src/kernel/vm/errors"
)

// Float represents a floating point number value.
type Float struct {
	ObjectImpl
	Value float64
}

func (o *Float) String() string {
	return strconv.FormatFloat(o.Value, 'f', -1, 64)
}

// TypeName returns the name of the type.
func (o *Float) TypeName() string {
	return "float"
}

// BinaryOp returns another object that is the result of a given binary operator and a right-hand side object.
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
		}
	}
	return nil, errors.ErrInvalidOperator
}

// Copy returns a copy of the type.
func (o *Float) Copy() Object {
	return &Float{Value: o.Value}
}

// IsFalsy returns true if the value of the type is falsy.
func (o *Float) IsFalsy() bool {
	return math.IsNaN(o.Value)
}

// Equals returns true if the value of the type is equal to the value of  another object.
func (o *Float) Equals(x Object) bool {
	t, ok := x.(*Float)
	if !ok {
		return false
	}
	return o.Value == t.Value
}

// ToFloat64 will try to convert object o to float64 value.
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
