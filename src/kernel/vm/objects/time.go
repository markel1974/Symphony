package objects

import (
	"encoding/gob"
	"time"
)

const (
	TimeType = "time"
)

func init() {
	gob.Register(&Time{})
}

// Time represents a custom object encapsulating a Go time.Time values with extended behaviors and operations.
type Time struct {
	factory IGateKeeper
	frame   int
	value   time.Time
}

// NewTime creates a new instance of Time wrapping the provided time.Time values.
func newTime(factory IGateKeeper, frame int, value time.Time) IObject {
	return &Time{
		factory: factory,
		frame:   frame,
		value:   value,
	}
}

// GateKeeper returns a reference to the GateKeeper associated with the Object.
func (o *Time) GateKeeper() IGateKeeper {
	return o.factory
}

// Frame returns the current frame value of the Object.
func (o *Time) Frame() int {
	return o.frame
}

// IndexGet attempts to retrieve a value at the given index and returns an error if the object is not indexable.
func (o *Time) IndexGet(_ int, _ IObject) (res IObject, err error) {
	return nil, ErrNotIndexable
}

// IndexSet attempts to assign a value to an index in the object but always returns ErrUnsupportedIndex,
// as this operation is unsupported.
func (o *Time) IndexSet(_, _ IObject) (err error) {
	return ErrUnsupportedIndex
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (o *Time) Iterate(_ int) IIterator {
	return nil
}

// CanIterate determines if the object can be iterated over and returns false for this implementation.
func (o *Time) CanIterate() bool {
	return false
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *Time) Call(_ int, _ ...IObject) (ret IObject, err error) {
	return nil, nil
}

// CanCall determines if the object can be invoked as a callable. Returns false for non-callable objects.
func (o *Time) CanCall() bool {
	return false
}

// Length returns the length of the Int object.
func (o *Time) Length() int {
	return 0
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

// LogicalOp performs logical comparison operations (e.g., <, >, <=, >=) between the Time object and another Time object.
func (o *Time) LogicalOp(_ int, op LogicalOperator, in IObject) (IObject, error) {
	switch rhs := in.(type) {
	case *Time:
		switch op {
		case OperatorLogicalLess:
			if o.value.Before(rhs.value) {
				return o.GateKeeper().TrueValue(), nil
			}
			return o.GateKeeper().FalseValue(), nil
		case OperatorLogicalGreater:
			if o.value.After(rhs.value) {
				return o.GateKeeper().TrueValue(), nil
			}
			return o.GateKeeper().FalseValue(), nil
		case OperatorLogicalLessEq:
			if o.value.Equal(rhs.value) || o.value.Before(rhs.value) {
				return o.GateKeeper().TrueValue(), nil
			}
			return o.GateKeeper().FalseValue(), nil
		case OperatorLogicalGreaterEq:
			if o.value.Equal(rhs.value) || o.value.After(rhs.value) {
				return o.GateKeeper().TrueValue(), nil
			}
			return o.GateKeeper().FalseValue(), nil
		default:
			return nil, ErrInvalidOperator
		}
	}
	return nil, ErrInvalidOperator
}

// ArithmeticOp performs arithmetic operations (addition, subtraction) on the Time object based on the given operator.
// Supports operations with Int and Time objects. Returns a new IObject or an error for invalid operators.
func (o *Time) ArithmeticOp(frame int, op ArithmeticOperator, in IObject) (IObject, error) {
	switch rhs := in.(type) {
	case *Int:
		switch op {
		case OperatorAdd:
			if rhs.value == 0 {
				return o, nil
			}
			return o.GateKeeper().NewTime(frame, o.value.Add(time.Duration(rhs.value))), nil
		case OperatorSub:
			if rhs.value == 0 {
				return o, nil
			}
			return o.GateKeeper().NewTime(frame, o.value.Add(time.Duration(-rhs.value))), nil
		default:
			return nil, ErrInvalidOperator
		}
	case *Time:
		switch op {
		case OperatorSub:
			return o.GateKeeper().NewInt(frame, int64(o.value.Sub(rhs.value))), nil
		default:
			return nil, ErrInvalidOperator
		}
	}
	return nil, ErrInvalidOperator
}

// Copy returns a new instance of the Time object with the same internal time values, duplicating its state.
func (o *Time) Copy(frame int, _ int) IObject {
	return o.GateKeeper().NewTime(frame, o.value)
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
