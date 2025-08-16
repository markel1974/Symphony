package objects

import (
	"errors"
	"fmt"
)

// ErrStackOverflow is a stack overflow error.
// ErrObjectAllocLimit is an object allocation limit error.
// ErrIndexOutOfBounds is an error where a given index is out of bounds.
// ErrInvalidIndexType represents an invalid index type.
// ErrInvalidIndexValueType represents an invalid index value type.
// ErrInvalidIndexOnError represents an invalid index on error.
// ErrInvalidOperator represents an error for invalid operator usage.
// ErrWrongNumArguments represents a wrong amount of argument error.
// ErrBytesLimit represents an error where the size of bytes value exceeds the limit.
// ErrStringLimit represents an error where the size of string value exceeds the limit.
// ErrNotIndexable is an error where an Object is not indexable.
// ErrNotIndexAssignable is an error where an Object is not index-assignable.
// ErrNotImplemented is an error where an Object has not implemented a required method.
// ErrInvalidRangeStep is an error where the range step is less than or equal to 0 in the range function.
var (
	ErrStackOverflow = errors.New("stack overflow")

	ErrObjectAllocLimit = errors.New("object allocation limit exceeded")

	ErrIndexOutOfBounds = errors.New("index out of bounds")

	ErrInvalidIndexType = errors.New("invalid index type")

	ErrInvalidIndexValueType = errors.New("invalid index value type")

	ErrInvalidIndexOnError = errors.New("invalid index on error")

	ErrInvalidOperator = errors.New("invalid operator")

	ErrWrongNumArguments = errors.New("wrong number of arguments")

	ErrBytesLimit = errors.New("exceeding bytes size limit")

	ErrStringLimit = errors.New("exceeding string size limit")

	ErrNotIndexable = errors.New("not indexable")

	ErrNotIndexAssignable = errors.New("not index-assignable")

	ErrNotImplemented = errors.New("not implemented")

	ErrInvalidRangeStep = errors.New("range step must be greater than 0")
)

// New creates and returns a new error with the specified message string.
func New(src string) error {
	return errors.New(src)
}

// Is reports whether any error in an error chain matches the target error. Uses errors.Is for comparison.
func Is(err error, target error) bool {
	return errors.Is(err, target)
}

// NewInvalidArgumentError creates an error indicating an argument has an unexpected type, providing its name, expected, and actual types.
func NewInvalidArgumentError(index int, expected string, found string) error {
	name := fmt.Sprintf("argument %d", index)
	return fmt.Errorf("invalid type for argument '%s': expected %s, found %s", name, expected, found)
}

// NewObjectError wraps a Go error into an IObject-compatible error object or returns TrueValue if the error is nil.
func NewObjectError(err error) IObject {
	if err == nil {
		return TrueValue
	}
	return NewError(NewStringNoSize(err.Error()))
}
