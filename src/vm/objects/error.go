package objects

import (
	"encoding/gob"
	"fmt"
)

const (
	ErrorType = "error"
)

const (
	maxErrorLen = 1024
)

func init() {
	gob.Register(&Error{})
}

// Error represents an object that encapsulates an error and implements the IObject interface.
type Error struct {
	Allocator
	err   string
	value IObject
}

// NewError creates and returns a new Error object with the specified values.
func newError(factory IGateKeeper, frame int, err string) IObject {
	if len(err) > maxErrorLen {
		err = err[:maxErrorLen]
	}
	return &Error{
		Allocator: Allocator{gk: factory, frame: frame},
		value:     factory.NewString(frame, err),
		err:       err,
	}
}

// AsBool returns true if the object is not empty, otherwise false.
func (o *Error) AsBool() bool {
	return false
}

// AsInt64 returns the length of the array as an int64 value.
func (o *Error) AsInt64() int64 {
	return 0
}

// AsFloat64 returns the length of the array as an int64 value.
func (o *Error) AsFloat64() float64 {
	return 0
}

// AsString returns the string representation of the Error object. If the values is nil, it returns "error".
func (o *Error) AsString() string {
	if o.value != nil {
		return o.value.AsString()
	}
	return ErrorType
}

// AssignValue sets the current object to the provided IObject, returning ErrNotAssignable if the operation is not supported.
func (o *Error) AssignValue(_ IObject) error {
	return ErrNotAssignable
}

// Nil checks if the object is nil and always returns false.
func (o *Error) Nil() bool {
	return false
}

// LogicalOp performs a logical operation on the Error object with the provided operator and operand, returning an error.
func (o *Error) LogicalOp(_ int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		return logicalOpNil(o.gk, op)
	}
	rhs, ok := rhsIn.(*Error)
	if !ok {
		return nil, ErrInvalidOperator
	}
	switch op {
	case OperatorLogicalEq:
		if o.value == rhs.value {
			return o.gk.TrueValue(), nil
		}
		return o.gk.FalseValue(), nil
	case OperatorLogicalNotEq:
		if o.value != rhs.value {
			return o.gk.TrueValue(), nil
		}
		return o.gk.FalseValue(), nil
	default:
		return nil, ErrInvalidOperator
	}
}

// ArithmeticOp performs the specified arithmetic operation on the Error object and always returns ErrInvalidOperator.
func (o *Error) ArithmeticOp(_ int, _ ArithmeticOperator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// IndexSet attempts to assign a value to an index in the object but always returns ErrIndexUnsupported,
// as this operation is unsupported.
func (o *Error) IndexSet(_, _ IObject) (err error) {
	return ErrIndexUnsupported
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (o *Error) Iterate(_ int) IIterator {
	return nil
}

// Iterable determines if the object can be iterated over and returns false for this implementation.
func (o *Error) Iterable() bool {
	return false
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *Error) Call(_ int, params ...IObject) (retCount uint, ret IObject, err error) {
	if len(params) != 1 {
		return 0, o.gk.UndefinedValue(), fmt.Errorf("expected 1 param, got %d", len(params))
	}
	switch params[0].AsString() {
	case "Error":
		return 1, o.value, nil
	}
	return 0, o.gk.UndefinedValue(), fmt.Errorf("invalid param: %s", params[0].AsString())
}

// Length returns the length of the Int object.
func (o *Error) Length() int {
	return 0
}

// Value returns the underlying IObject value of the Error object.
func (o *Error) Value() IObject {
	return o.value
}

// TypeName returns the name of the type as a string, which is "error".
func (o *Error) TypeName() string {
	return ErrorType
}

// Falsy returns true, indicating that the Error object is always considered falsy in a boolean context.
func (o *Error) Falsy() bool {
	return true // error is always false.
}

// Copy creates and returns a new instance of the Error object with the same underlying values.
func (o *Error) Copy(frame int, _ int) IObject {
	return o.GateKeeper().NewError(frame, o.err)
}

// Equals checks if the current Error object is equal to another object using pointer equality.
func (o *Error) Equals(x IObject) bool {
	return o == x // pointer equality
}

// IndexGet retrieves the values associated with the "values" index in an Error object or returns an error for invalid indices.
func (o *Error) IndexGet(_ int, index IObject) (IObject, error) {
	if strIdx, _ := o.GateKeeper().ToString(index); strIdx != "values" {
		return nil, ErrIndexInvalidValueType
	}
	return o.value, nil
}

// Count returns the total number of elements in the instance and its sub-elements.
func (o *Error) Count() int {
	return o.value.Count()
}
