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
// ErrInvalidArgumentsNumber represents a wrong amount of argument error.
// ErrBytesLimit represents an error where the size of bytes value exceeds the limit.
// ErrStringLimit represents an error where the size of string value exceeds the limit.
// ErrNotIndexable is an error where an Object is not indexable.
// ErrUnsupportedIndex is an error where an Object is not index-assignable.
// ErrInvalidRangeStep is an error where the range step is less than or equal to 0 in the range function.
var (
	ErrDivisionByZero = errors.New("division by zero")

	ErrStackOverflow = errors.New("stack overflow")

	ErrInvalidIndexValueType = errors.New("invalid index value type")

	ErrInvalidOperator = errors.New("invalid operator")

	ErrInvalidArgumentsNumber = errors.New("invalid arguments number")

	ErrInvalidRangeStep = errors.New("invalid range step")

	ErrLimitExceed = errors.New("limit exceed")

	ErrIndexOutOfBounds = errors.New("index out of bounds")

	ErrInvalidIndexType = errors.New("invalid index type")

	ErrNotIndexable = errors.New("not indexable")

	ErrUnsupportedIndex = errors.New("unsupported index")

	ErrUnsupportedOperation = errors.New("unsupported operation")

	ErrSelectorNotProvided = errors.New("selector not provided")

	ErrAllocationLimit = errors.New("allocation limit")
)

// NewInvalidArgumentError creates an error indicating an argument has an unexpected type, providing its name, expected, and actual types.
func NewInvalidArgumentError(index int, expected string, found string) error {
	name := fmt.Sprintf("argument %d", index)
	return fmt.Errorf("invalid type for argument '%s': expected %s, found %s", name, expected, found)
}

func ComputeIndexGetError(err error, dst string, sel string) error {
	if errors.Is(err, ErrNotIndexable) {
		return fmt.Errorf("not indexable: %s", dst)
	}
	if errors.Is(err, ErrInvalidIndexType) {
		return fmt.Errorf("invalid index type: %s", sel)
	}
	return err
}

func ComputeIndexSetError(err error, dst string, src string) error {
	if errors.Is(err, ErrUnsupportedIndex) {
		return fmt.Errorf("not index-assignable: %s", dst)
	}
	if errors.Is(err, ErrInvalidIndexValueType) {
		return fmt.Errorf("invaid index values type: %s", src)
	}
	return err
}

func ComputeCallError(err error, fn string) error {
	if errors.Is(err, ErrInvalidArgumentsNumber) {
		return fmt.Errorf("wrong number of arguments in call to '%s'", fn)
	}
	return err
}
