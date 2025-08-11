package objects

import (
	"strconv"

	"github.com/markel1974/c64emu/src/kernel/vm/errors"
)

// Int represents a wrapper for a 64-bit signed integer, providing methods for operations and comparisons.
type Int struct {
	ObjectImpl
	value int64
}

// NewInt creates and returns a new instance of the Int type with the specified int64 value.
func NewInt(value int64) *Int {
	return &Int{value: value}
}

// Value returns the int64 value of the Int type.
func (o *Int) Value() int64 {
	return o.value
}

// String returns the string representation of the Int type, converting its int64 value to a base-10 formatted string.
func (o *Int) String() string {
	return strconv.FormatInt(o.value, 10)
}

// TypeName returns the name of the type as a string, which is "int".
func (o *Int) TypeName() string {
	return "int"
}

// BinaryOp performs a binary operation using the specified operator between the current Int instance and another IObject.
// Returns the resulting IObject or an error if the operation is invalid.
func (o *Int) BinaryOp(op Operator, rhs IObject) (IObject, error) {
	switch rhs := rhs.(type) {
	case *Int:
		switch op {
		case OperatorAdd:
			r := o.value + rhs.value
			if r == o.value {
				return o, nil
			}
			return &Int{value: r}, nil
		case OperatorSub:
			r := o.value - rhs.value
			if r == o.value {
				return o, nil
			}
			return &Int{value: r}, nil
		case OperatorMul:
			r := o.value * rhs.value
			if r == o.value {
				return o, nil
			}
			return &Int{value: r}, nil
		case OperatorQuo:
			r := o.value / rhs.value
			if r == o.value {
				return o, nil
			}
			return &Int{value: r}, nil
		case OperatorRem:
			r := o.value % rhs.value
			if r == o.value {
				return o, nil
			}
			return &Int{value: r}, nil
		case OperatorAnd:
			r := o.value & rhs.value
			if r == o.value {
				return o, nil
			}
			return &Int{value: r}, nil
		case OperatorOr:
			r := o.value | rhs.value
			if r == o.value {
				return o, nil
			}
			return &Int{value: r}, nil
		case OperatorXor:
			r := o.value ^ rhs.value
			if r == o.value {
				return o, nil
			}
			return &Int{value: r}, nil
		case OperatorAndNot:
			r := o.value &^ rhs.value
			if r == o.value {
				return o, nil
			}
			return &Int{value: r}, nil
		case OperatorShl:
			r := o.value << uint64(rhs.value)
			if r == o.value {
				return o, nil
			}
			return &Int{value: r}, nil
		case OperatorShr:
			r := o.value >> uint64(rhs.value)
			if r == o.value {
				return o, nil
			}
			return &Int{value: r}, nil
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
	case *Float:
		switch op {
		case OperatorAdd:
			return &Float{value: float64(o.value) + rhs.value}, nil
		case OperatorSub:
			return &Float{value: float64(o.value) - rhs.value}, nil
		case OperatorMul:
			return &Float{value: float64(o.value) * rhs.value}, nil
		case OperatorQuo:
			return &Float{value: float64(o.value) / rhs.value}, nil
		case OperatorLess:
			if float64(o.value) < rhs.value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case OperatorGreater:
			if float64(o.value) > rhs.value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case OperatorLessEq:
			if float64(o.value) <= rhs.value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case OperatorGreaterEq:
			if float64(o.value) >= rhs.value {
				return TrueValue, nil
			}
			return FalseValue, nil
		default:
			return nil, errors.ErrInvalidOperator
		}
	case *Char:
		switch op {
		case OperatorAdd:
			return &Char{value: rune(o.value) + rhs.value}, nil
		case OperatorSub:
			return &Char{value: rune(o.value) - rhs.value}, nil
		case OperatorLess:
			if o.value < int64(rhs.value) {
				return TrueValue, nil
			}
			return FalseValue, nil
		case OperatorGreater:
			if o.value > int64(rhs.value) {
				return TrueValue, nil
			}
			return FalseValue, nil
		case OperatorLessEq:
			if o.value <= int64(rhs.value) {
				return TrueValue, nil
			}
			return FalseValue, nil
		case OperatorGreaterEq:
			if o.value >= int64(rhs.value) {
				return TrueValue, nil
			}
			return FalseValue, nil
		default:
			return nil, errors.ErrInvalidOperator
		}
	}
	return nil, errors.ErrInvalidOperator
}

// Copy creates and returns a new instance of the Int object, duplicating its value.
func (o *Int) Copy() IObject {
	return &Int{value: o.value}
}

// Falsy checks if the integer value is considered falsy (equal to 0) and returns true if it is, otherwise false.
func (o *Int) Falsy() bool {
	return o.value == 0
}

// Equals compares the Int object with another IObject for equality and returns true if both are of type Int and values are equal.
func (o *Int) Equals(x IObject) bool {
	t, ok := x.(*Int)
	if !ok {
		return false
	}
	return o.value == t.value
}

// ToInt64 converts an IObject implementation to an int64 if possible and returns whether the conversion was successful.
func ToInt64(o IObject) (v int64, ok bool) {
	switch o := o.(type) {
	case *Int:
		v = o.value
		ok = true
	case *Float:
		v = int64(o.value)
		ok = true
	case *Char:
		v = int64(o.value)
		ok = true
	case *Bool:
		if o == TrueValue {
			v = 1
		}
		ok = true
	case *String:
		c, err := strconv.ParseInt(o.value, 10, 64)
		if err == nil {
			v = c
			ok = true
		}
	}
	return
}

// ToInt attempts to convert the given IObject into an integer type.
// Returns the integer values and a boolean indicating success.
func ToInt(o IObject) (v int, ok bool) {
	switch o := o.(type) {
	case *Int:
		v = int(o.value)
		ok = true
	case *Float:
		v = int(o.value)
		ok = true
	case *Char:
		v = int(o.value)
		ok = true
	case *Bool:
		if o == TrueValue {
			v = 1
		}
		ok = true
	case *String:
		c, err := strconv.ParseInt(o.value, 10, 64)
		if err == nil {
			v = int(c)
			ok = true
		}
	}
	return
}
