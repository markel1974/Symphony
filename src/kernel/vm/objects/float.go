package objects

import (
	"encoding/gob"
	"math"
	"strconv"
)

const (
	FloatType = "float"
)

func init() {
	gob.Register(&Float{})
}

// Float represents a floating-point number and provides operations and behaviors specific to numeric types.
// It embeds Object to implement common interface methods and extends behavior where necessary.
// The value field holds the actual float64 values encapsulated by the Float type.
type Float struct {
	gk    IGateKeeper
	frame int
	value float64
}

// NewFloat creates and returns a pointer to a new Float object initialized with the specified float64 values.
func newFloat(gk IGateKeeper, frame int, value float64) IObject {
	return &Float{
		gk:    gk,
		frame: frame,
		value: value,
	}
}

// GateKeeper returns a reference to the GateKeeper associated with the Object.
func (o *Float) GateKeeper() IGateKeeper {
	return o.gk
}

// Frame returns the current frame value of the Object.
func (o *Float) Frame() int {
	return o.frame
}

// IndexGet attempts to retrieve a value at the given index and returns an error if the object is not indexable.
func (o *Float) IndexGet(_ int, _ IObject) (res IObject, err error) {
	return nil, ErrNotIndexable
}

// IndexSet attempts to assign a value to an index in the object but always returns ErrUnsupportedIndex,
// as this operation is unsupported.
func (o *Float) IndexSet(_, _ IObject) (err error) {
	return ErrUnsupportedIndex
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (o *Float) Iterate(_ int) IIterator {
	return nil
}

// CanIterate determines if the object can be iterated over and returns false for this implementation.
func (o *Float) CanIterate() bool {
	return false
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *Float) Call(_ int, _ ...IObject) (ret IObject, err error) {
	return nil, nil
}

// CanCall determines if the object can be invoked as a callable. Returns false for non-callable objects.
func (o *Float) CanCall() bool {
	return false
}

// Length returns the length of the Int object.
func (o *Float) Length() int {
	return 0
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
			return o.GateKeeper().NewFloat(frame, r), nil
		case OperatorSub:
			r := o.value - rhs.value
			if r == o.value {
				return o, nil
			}
			return o.GateKeeper().NewFloat(frame, r), nil
		case OperatorMul:
			r := o.value * rhs.value
			if r == o.value {
				return o, nil
			}
			return o.GateKeeper().NewFloat(frame, r), nil
		case OperatorQuo:
			r := o.value / rhs.value
			if r == o.value {
				return o, nil
			}
			return o.GateKeeper().NewFloat(frame, r), nil
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
	case *Int:
		switch op {
		case OperatorAdd:
			r := o.value + float64(rhs.value)
			if r == o.value {
				return o, nil
			}
			return o.GateKeeper().NewFloat(frame, r), nil
		case OperatorSub:
			r := o.value - float64(rhs.value)
			if r == o.value {
				return o, nil
			}
			return o.GateKeeper().NewFloat(frame, r), nil
		case OperatorMul:
			r := o.value * float64(rhs.value)
			if r == o.value {
				return o, nil
			}
			return o.GateKeeper().NewFloat(frame, r), nil
		case OperatorQuo:
			r := o.value / float64(rhs.value)
			if r == o.value {
				return o, nil
			}
			return o.GateKeeper().NewFloat(frame, r), nil
		case OperatorLess:
			if o.value < float64(rhs.value) {
				return o.GateKeeper().TrueValue(), nil
			}
			return o.GateKeeper().FalseValue(), nil
		case OperatorGreater:
			if o.value > float64(rhs.value) {
				return o.GateKeeper().TrueValue(), nil
			}
			return o.GateKeeper().FalseValue(), nil
		case OperatorLessEq:
			if o.value <= float64(rhs.value) {
				return o.GateKeeper().TrueValue(), nil
			}
			return o.GateKeeper().FalseValue(), nil
		case OperatorGreaterEq:
			if o.value >= float64(rhs.value) {
				return o.GateKeeper().TrueValue(), nil
			}
			return o.GateKeeper().FalseValue(), nil
		default:
			return nil, ErrInvalidOperator
		}
	}
	return nil, ErrInvalidOperator
}

// Copy creates and returns a new instance of the Float object, duplicating its current state.
func (o *Float) Copy(frame int, _ int) IObject {
	return o.GateKeeper().NewFloat(frame, o.value)
}

// Falsy determines if the float object is considered falsy, returning true if the values is NaN; otherwise, false.
func (o *Float) Falsy() bool {
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
