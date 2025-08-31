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
	gk    IGateKeeper
	frame int
	err   string
	value IObject
}

// NewError creates and returns a new Error object with the specified values.
func newError(factory IGateKeeper, frame int, err string) IObject {
	if len(err) > maxErrorLen {
		err = err[:maxErrorLen]
	}
	return &Error{
		gk:    factory,
		frame: frame,
		value: factory.NewString(frame, err),
		err:   err,
	}
}

// GateKeeper returns a reference to the GateKeeper associated with the Object.
func (o *Error) GateKeeper() IGateKeeper {
	return o.gk
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
		return fmt.Sprintf("%s: %s", ErrorType, o.value.AsString())
	}
	return ErrorType
}

// AssignValue sets the current object to the provided IObject, returning ErrNotAssignable if the operation is not supported.
func (o *Error) AssignValue(_ IObject) error {
	return ErrNotAssignable
}

// Frame returns the current frame value of the Object.
func (o *Error) Frame() int {
	return o.frame
}

// LogicalOp performs a logical operation on the Error object with the provided operator and operand, returning an error.
func (o *Error) LogicalOp(_ int, _ LogicalOperator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// ArithmeticOp performs the specified arithmetic operation on the Error object and always returns ErrInvalidOperator.
func (o *Error) ArithmeticOp(_ int, _ ArithmeticOperator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// IndexSet attempts to assign a value to an index in the object but always returns ErrUnsupportedIndex,
// as this operation is unsupported.
func (o *Error) IndexSet(_, _ IObject) (err error) {
	return ErrUnsupportedIndex
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (o *Error) Iterate(_ int) IIterator {
	return nil
}

// CanIterate determines if the object can be iterated over and returns false for this implementation.
func (o *Error) CanIterate() bool {
	return false
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *Error) Call(_ int, _ ...IObject) (ret IObject, err error) {
	return nil, nil
}

// CanCall determines if the object can be invoked as a callable. Returns false for non-callable objects.
func (o *Error) CanCall() bool {
	return false
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
func (o *Error) IndexGet(_ int, index IObject) (res IObject, err error) {
	if strIdx, _ := o.GateKeeper().ToString(index); strIdx != "values" {
		err = ErrInvalidIndexValueType
		return
	}
	res = o.value
	return
}
