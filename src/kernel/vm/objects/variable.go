package objects

import (
	"errors"
)

// Variable represents a named entity with a string identifier and an associated value of type objects.IObject.
type Variable struct {
	name  string
	value IObject
}

// NewVariable creates a new Variable instance with a given name and value, converting the value to an internal object type.
// It returns the newly created *Variable and an error if any issue occurs during value conversion.
func NewVariable(name string, value interface{}) (*Variable, error) {
	obj := FromInterface(value)
	return &Variable{
		name:  name,
		value: obj,
	}, nil
}

// Name returns the name of the Variable as a string.
func (v *Variable) Name() string {
	return v.name
}

// Value retrieves the value of the Variable as an interface{}.
func (v *Variable) Value() interface{} {
	return ToInterface(v.value)
}

// ValueType returns the type name of the value held by the Variable.
func (v *Variable) ValueType() string {
	return v.value.TypeName()
}

// Int converts the value of the Variable to an integer using objects.ToInt and returns the result.
func (v *Variable) Int() int {
	c, _ := ToInt64(v.value)
	return int(c)
}

// Int64 converts the value of the Variable to an int64 if possible and returns it.
func (v *Variable) Int64() int64 {
	c, _ := ToInt64(v.value)
	return c
}

// Float converts the Variable's underlying value to a float64 if possible and returns the result.
func (v *Variable) Float() float64 {
	c, _ := ToFloat64(v.value)
	return c
}

// Char converts the value of the Variable to a rune and returns it. If conversion fails, a zero-value rune is returned.
func (v *Variable) Char() rune {
	c, _ := ToRune(v.value)
	return c
}

// Bool converts the stored value of the Variable to a boolean and returns it.
func (v *Variable) Bool() bool {
	c, _ := ToBool(v.value)
	return c
}

// Array converts the Variable's value to a slice of interface{} if it is of type *objects.Array, otherwise returns nil.
func (v *Variable) Array() []interface{} {
	switch val := v.value.(type) {
	case *Array:
		var arr []interface{}
		for _, e := range val.Values() {
			arr = append(arr, ToInterface(e))
		}
		return arr
	}
	return nil
}

// ArrayInt extracts integer elements from the Variable's value if it is an array and returns them as a slice of integers.
func (v *Variable) ArrayInt() []int {
	switch val := v.value.(type) {
	case *Array:
		var arr []int
		for _, e := range val.Values() {
			i := ToInterface(e)
			if val, ok := i.(int); ok {
				arr = append(arr, val)
			}
		}
		return arr
	}
	return nil
}

// ArrayBool extracts and returns a slice of bool values from the underlying array if it contains boolean elements.
// Returns nil if the value is not an array or contains non-boolean elements.
func (v *Variable) ArrayBool() []bool {
	switch val := v.value.(type) {
	case *Array:
		var arr []bool
		for _, e := range val.Values() {
			i := ToInterface(e)
			if val, ok := i.(bool); ok {
				arr = append(arr, val)
			}
		}
		return arr
	}
	return nil
}

// ArrayString returns a slice of strings extracted from an array value in the Variable if it contains string elements.
func (v *Variable) ArrayString() []string {
	switch val := v.value.(type) {
	case *Array:
		var arr []string
		for _, e := range val.Values() {
			i := ToInterface(e)
			if val, ok := i.(string); ok {
				arr = append(arr, val)
			}
		}
		return arr
	}
	return nil
}

// Map converts the underlying value of the Variable to a map[string]interface{} if it is of type *objects.Map.
func (v *Variable) Map() map[string]interface{} {
	switch val := v.value.(type) {
	case *Map:
		kv := make(map[string]interface{})
		for mk, mv := range val.Values() {
			kv[mk] = ToInterface(mv)
		}
		return kv
	}
	return nil
}

// String returns the string representation of the value stored in the Variable.
func (v *Variable) String() string {
	c, _ := ToString(v.value)
	return c
}

// Bytes converts the internal value of the Variable to a byte slice, returning nil if conversion fails.
func (v *Variable) Bytes() []byte {
	c, _ := ToByteSlice(v.value)
	return c
}

// Error checks if the value of the Variable is of type *objects.Error and returns it as a Go error if true.
func (v *Variable) Error() error {
	err, ok := v.value.(*Error)
	if ok {
		return errors.New(err.String())
	}
	return nil
}

// Object returns the underlying objects.IObject instance held by the Variable.
func (v *Variable) Object() IObject {
	return v.value
}

// IsUndefined checks if the variable's value is equal to the predefined undefined value.
func (v *Variable) IsUndefined() bool {
	return v.value == UndefinedValue
}
