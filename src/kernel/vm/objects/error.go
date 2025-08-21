package objects

import (
	"fmt"
)

const (
	ErrorType = "error"
)

const (
	maxErrorLen = 1024
)

// Error represents an object that encapsulates an error and implements the IObject interface.
type Error struct {
	IObject
	err   string
	value IObject
}

// NewError creates and returns a new Error object with the specified values.
func newError(factory *GateKeeper, frame int, err string) *Error {
	if len(err) > maxErrorLen {
		err = err[:maxErrorLen]
	}
	return &Error{
		IObject: factory.newObject(frame),
		value:   factory.NewString(frame, err),
		err:     err,
	}
}

// Value returns the underlying IObject value of the Error object.
func (o *Error) Value() IObject {
	return o.value
}

// TypeName returns the name of the type as a string, which is "error".
func (o *Error) TypeName() string {
	return ErrorType
}

// String returns the string representation of the Error object. If the values is nil, it returns "error".
func (o *Error) String() string {
	if o.value != nil {
		return fmt.Sprintf("%s: %s", ErrorType, o.value.String())
	}
	return ErrorType
}

// Boolean returns true, indicating that the Error object is always considered falsy in a boolean context.
func (o *Error) Boolean() bool {
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
		err = ErrInvalidIndexOnError
		return
	}
	res = o.value
	return
}
