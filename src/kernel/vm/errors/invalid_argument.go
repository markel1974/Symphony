package errors

import (
	"fmt"
)

/*
// ErrInvalidArgumentType represents an invalid argument value type error.
type ErrInvalidArgumentType struct {
	Name     string
	Expected string
	Found    string
}

func (e ErrInvalidArgumentType) Error() string {
	return fmt.Sprintf("invalid type for argument '%s': expected %s, found %s", e.Name, e.Expected, e.Found)
}
*/

func NewInvalidArgumentType(name string, expected string, found string) error {
	return fmt.Errorf("invalid type for argument '%s': expected %s, found %s", name, expected, found)
}
