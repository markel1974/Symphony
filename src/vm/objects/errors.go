package objects

import (
	"errors"
	"fmt"
)

var (
	// ErrDivisionByZero indicates an error caused by an attempt to divide a numeric value by zero.
	ErrDivisionByZero = errors.New("division by zero")

	// ErrStackOverflow is returned when the execution stack exceeds its allowed limit, indicating a stack overflow condition.
	ErrStackOverflow = errors.New("stack overflow")

	// ErrInvalidIndexValueType represents an error occurring when an index operation involves an invalid index value type.
	ErrInvalidIndexValueType = errors.New("invalid index value type")

	// ErrInvalidOperator represents an error returned when an unsupported or invalid operator is encountered.
	ErrInvalidOperator = errors.New("invalid operator")

	// ErrNotAssignable is used to indicate an error when attempting to assign an incompatible or unsupported value.
	ErrNotAssignable = errors.New("not assignable")

	// ErrInvalidArgument represents an error when an invalid argument is provided to a function or operation.
	ErrInvalidArgument = errors.New("invalid argument")

	// ErrInvalidArgumentsNumber indicates that the number of arguments provided is invalid for the operation or function.
	ErrInvalidArgumentsNumber = errors.New("invalid arguments number")

	// ErrLimitExceed is returned when an operation exceeds a predefined limit, such as maximum length or capacity.
	ErrLimitExceed = errors.New("limit exceed")

	// ErrIndexOutOfBounds indicates an error where an index is out of the allowable range for an operation.
	ErrIndexOutOfBounds = errors.New("index out of bounds")

	// ErrInvalidIndexType represents an error occurring when an index is of an invalid or unsupported type for the given operation.
	ErrInvalidIndexType = errors.New("invalid index type")

	// ErrNotIndexable represents an error indicating that an object is not indexable.
	ErrNotIndexable = errors.New("not indexable")

	// ErrUnsupportedIndex indicates that an index operation is unsupported for the given object.
	ErrUnsupportedIndex = errors.New("unsupported index")

	// ErrSelectorNotProvided indicates that no selector was provided for an operation requiring at least one selector.
	ErrSelectorNotProvided = errors.New("selector not provided")

	// ErrUnimplemented indicates that the requested functionality has not been implemented.
	ErrUnimplemented = errors.New("unimplemented")
)

// NewInvalidArgumentError creates and returns an error indicating that an argument has an invalid type.
func NewInvalidArgumentError(index int, expected string, found string) error {
	name := fmt.Sprintf("argument %d", index)
	return fmt.Errorf("invalid type for argument '%s': expected %s, found %s", name, expected, found)
}

// ComputeIndexGetError converts indexing-related errors into descriptive error messages for non-indexable or invalid types.
func ComputeIndexGetError(err error, dst string, sel string) error {
	if errors.Is(err, ErrNotIndexable) {
		return fmt.Errorf("not indexable: %s", dst)
	}
	if errors.Is(err, ErrInvalidIndexType) {
		return fmt.Errorf("invalid index type: %s", sel)
	}
	return err
}

// ComputeIndexSetError transforms specific index operation errors into user-friendly error messages.
// Returns an error indicating why the index operation failed based on the type of the given error.
func ComputeIndexSetError(err error, dst string, src string) error {
	if errors.Is(err, ErrUnsupportedIndex) {
		return fmt.Errorf("not index-assignable: %s", dst)
	}
	if errors.Is(err, ErrInvalidIndexValueType) {
		return fmt.Errorf("invaid index values type: %s", src)
	}
	return err
}

// ComputeCallError checks if the provided error matches ErrInvalidArgumentsNumber and returns a descriptive error, if so.
func ComputeCallError(err error, fn string) error {
	if errors.Is(err, ErrInvalidArgumentsNumber) {
		return fmt.Errorf("wrong number of arguments in call to '%s'", fn)
	}
	return err
}
