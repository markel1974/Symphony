package objects

import (
	"strconv"
)

const (
	IntType = "int"
)

// Int represents an integer type with a 64-bit value and methods for operations, equality, and object behavior.
type Int struct {
	gk    *GateKeeper
	frame int
	value int64
}

// NewInt creates and returns a new instance of the Int struct initialized with the specified int64 value.
func newInt(factory *GateKeeper, frame int, value int64) IObject {
	return &Int{
		gk:    factory,
		frame: frame,
		value: value,
	}
}

// GateKeeper returns a reference to the GateKeeper associated with the Object.
func (o *Int) GateKeeper() *GateKeeper {
	return o.gk
}

// Frame returns the current frame value of the Object.
func (o *Int) Frame() int {
	return o.frame
}

// IndexGet attempts to retrieve a value at the given index and returns an error if the object is not indexable.
func (o *Int) IndexGet(_ int, _ IObject) (res IObject, err error) {
	return nil, ErrNotIndexable
}

// IndexSet attempts to assign a value to an index in the object but always returns ErrNotIndexAssignable,
// as this operation is unsupported.
func (o *Int) IndexSet(_, _ IObject) (err error) {
	return ErrNotIndexAssignable
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (o *Int) Iterate(_ int) IIterator {
	return nil
}

// CanIterate determines if the object can be iterated over and returns false for this implementation.
func (o *Int) CanIterate() bool {
	return false
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *Int) Call(_ int, _ ...IObject) (ret IObject, err error) {
	return nil, nil
}

// CanCall determines if the object can be invoked as a callable. Returns false for non-callable objects.
func (o *Int) CanCall() bool {
	return false
}

// Length returns the length of the Int object.
func (o *Int) Length() int {
	return 0
}

// Value returns the underlying int64 value of the Int object.
func (o *Int) Value() int64 {
	return o.value
}

// SetValue sets the underlying int64 value of the Int object.
func (o *Int) SetValue(value int64) {
	o.value = value
}

// String returns the string representation of the Int value using base 10 format.
func (o *Int) String() string {
	return strconv.FormatInt(o.value, 10)
}

// TypeName returns the name of the type as a string, which is "int" for this object.
func (o *Int) TypeName() string {
	return IntType
}

// BinaryOp performs a binary operation using the specified operator and right-hand side operand, returning the result.
func (o *Int) BinaryOp(frame int, op Operator, rhs IObject) (IObject, error) {
	switch rhs := rhs.(type) {
	case *Int:
		switch op {
		case OperatorAdd:
			r := o.value + rhs.value
			if r == o.value {
				return o, nil
			}
			return o.GateKeeper().NewInt(frame, r), nil
		case OperatorSub:
			r := o.value - rhs.value
			if r == o.value {
				return o, nil
			}
			return o.GateKeeper().NewInt(frame, r), nil
		case OperatorMul:
			r := o.value * rhs.value
			if r == o.value {
				return o, nil
			}
			return o.GateKeeper().NewInt(frame, r), nil
		case OperatorQuo:
			if rhs.value == 0 {
				return nil, ErrDivisionByZero
			}
			r := o.value / rhs.value
			if r == o.value {
				return o, nil
			}
			return o.GateKeeper().NewInt(frame, r), nil
		case OperatorRem:
			if rhs.value == 0 {
				return nil, ErrDivisionByZero
			}
			r := o.value % rhs.value
			if r == o.value {
				return o, nil
			}
			return o.GateKeeper().NewInt(frame, r), nil
		case OperatorAnd:
			r := o.value & rhs.value
			if r == o.value {
				return o, nil
			}
			return o.GateKeeper().NewInt(frame, r), nil
		case OperatorOr:
			r := o.value | rhs.value
			if r == o.value {
				return o, nil
			}
			return o.GateKeeper().NewInt(frame, r), nil
		case OperatorXor:
			r := o.value ^ rhs.value
			if r == o.value {
				return o, nil
			}
			return o.GateKeeper().NewInt(frame, r), nil
		case OperatorAndNot:
			r := o.value &^ rhs.value
			if r == o.value {
				return o, nil
			}
			return o.GateKeeper().NewInt(frame, r), nil
		case OperatorShl:
			r := o.value << uint64(rhs.value)
			if r == o.value {
				return o, nil
			}
			return o.GateKeeper().NewInt(frame, r), nil
		case OperatorShr:
			r := o.value >> uint64(rhs.value)
			if r == o.value {
				return o, nil
			}
			return o.GateKeeper().NewInt(frame, r), nil
		case OperatorLess:
			if o.value < rhs.value {
				return o.GateKeeper().TrueValue(), nil
			}
			return o.GateKeeper().FalseValue(), nil
		case OperatorGreater:
			if o.value > rhs.value {
				return o.GateKeeper().TrueValue(), nil
			}
			return o.GateKeeper().FalseValue(), nil
		case OperatorLessEq:
			if o.value <= rhs.value {
				return o.GateKeeper().TrueValue(), nil
			}
			return o.GateKeeper().FalseValue(), nil
		case OperatorGreaterEq:
			if o.value >= rhs.value {
				return o.GateKeeper().TrueValue(), nil
			}
			return o.GateKeeper().FalseValue(), nil
		default:
			return nil, ErrInvalidOperator
		}
	case *Float:
		switch op {
		case OperatorAdd:
			return o.GateKeeper().NewFloat(frame, float64(o.value)+rhs.value), nil
		case OperatorSub:
			return o.GateKeeper().NewFloat(frame, float64(o.value)-rhs.value), nil
		case OperatorMul:
			return o.GateKeeper().NewFloat(frame, float64(o.value)*rhs.value), nil
		case OperatorQuo:
			return o.GateKeeper().NewFloat(frame, float64(o.value)/rhs.value), nil
		case OperatorLess:
			if float64(o.value) < rhs.value {
				return o.GateKeeper().TrueValue(), nil
			}
			return o.GateKeeper().FalseValue(), nil
		case OperatorGreater:
			if float64(o.value) > rhs.value {
				return o.GateKeeper().TrueValue(), nil
			}
			return o.GateKeeper().FalseValue(), nil
		case OperatorLessEq:
			if float64(o.value) <= rhs.value {
				return o.GateKeeper().TrueValue(), nil
			}
			return o.GateKeeper().FalseValue(), nil
		case OperatorGreaterEq:
			if float64(o.value) >= rhs.value {
				return o.GateKeeper().TrueValue(), nil
			}
			return o.GateKeeper().FalseValue(), nil
		default:
			return nil, ErrInvalidOperator
		}
	case *Char:
		switch op {
		case OperatorAdd:
			return o.GateKeeper().NewChar(frame, rune(o.value)+rhs.value), nil
		case OperatorSub:
			return o.GateKeeper().NewChar(frame, rune(o.value)-rhs.value), nil
		case OperatorLess:
			if o.value < int64(rhs.value) {
				return o.GateKeeper().TrueValue(), nil
			}
			return o.GateKeeper().FalseValue(), nil
		case OperatorGreater:
			if o.value > int64(rhs.value) {
				return o.GateKeeper().TrueValue(), nil
			}
			return o.GateKeeper().FalseValue(), nil
		case OperatorLessEq:
			if o.value <= int64(rhs.value) {
				return o.GateKeeper().TrueValue(), nil
			}
			return o.GateKeeper().FalseValue(), nil
		case OperatorGreaterEq:
			if o.value >= int64(rhs.value) {
				return o.GateKeeper().TrueValue(), nil
			}
			return o.GateKeeper().FalseValue(), nil
		default:
			return nil, ErrInvalidOperator
		}
	}
	return nil, ErrInvalidOperator
}

// Copy creates and returns a new instance of the Int object with the same value as the current instance.
func (o *Int) Copy(frame int, _ int) IObject {
	return o.GateKeeper().NewInt(frame, o.value)
}

// Boolean checks whether the integer value is considered falsy. Returns true if the value is 0, otherwise false.
func (o *Int) Boolean() bool {
	return o.value == 0
}

// Equals checks if the current Int object is equal to another IObject of type Int by comparing their values.
func (o *Int) Equals(x IObject) bool {
	t, ok := x.(*Int)
	if !ok {
		return false
	}
	return o.value == t.value
}
