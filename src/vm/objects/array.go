package objects

import (
	"encoding/gob"
	"fmt"
	"strings"
)

const (
	ArrayType = "array"
)

func init() {
	gob.Register(&Array{})
}

// Array represents a collection of IObject elements, providing methods for manipulation, indexing, and iteration.
type Array struct {
	Allocator
	values []IObject
}

// NewArray creates and returns a new Array object initialized with the provided slice of IObject elements.
func newArray(gk IGateKeeper, frame int, value []IObject) IObject {
	if len(value) > maxArrayLen {
		value = value[0:maxArrayLen]
	}
	return &Array{
		Allocator: Allocator{gk: gk, frame: frame},
		values:    value,
	}
}

// AsBool returns true if the array is not empty, otherwise false.
func (o *Array) AsBool() bool {
	return len(o.values) > 0
}

// AsInt64 returns the length of the array as an int64 value.
func (o *Array) AsInt64() int64 {
	return int64(len(o.values))
}

// AsFloat64 returns the length of the array as an int64 value.
func (o *Array) AsFloat64() float64 {
	return float64(len(o.values))
}

// AsString returns a string representation of the Array, formatting its elements in a comma-separated list enclosed in brackets.
func (o *Array) AsString() string {
	var elements []string
	for _, e := range o.values {
		elements = append(elements, e.AsString())
	}
	return fmt.Sprintf("[%s]", strings.Join(elements, "; "))
}

// AssignValue assigns the elements of another Array to the current Array if the input is of type *Array, otherwise returns an error.
func (o *Array) AssignValue(v IObject) error {
	target, ok := v.(*Array)
	if !ok {
		return ErrNotAssignable
	}
	o.Assign(target.values)
	return nil
}

// Nil checks if the object is nil and always returns false.
func (o *Array) Nil() bool {
	return false
}

// Call executes the Array instance as a callable object, passing the given arguments, and returns the result or an error.
func (o *Array) Call(_ int, _ ...IObject) (retCount uint, ret IObject, err error) {
	return 0, nil, nil
}

// TypeName returns the string "array", representing the type name of the Array object.
func (o *Array) TypeName() string {
	return ArrayType
}

// Values return the slice of IObject elements stored in the Array.
func (o *Array) Values() []IObject {
	return o.values
}

// Index returns the element at the specified index in the Array, or an error if the index is out of bounds.
func (o *Array) Index(index int) (IObject, error) {
	if index < 0 || index >= len(o.values) {
		return nil, ErrIndexOutOfBounds
	}
	return o.values[index], nil
}

// Length returns the number of elements in the Array.
func (o *Array) Length() int {
	return len(o.values)
}

// SetValue assigns a given IObject value to the specified index in the Array if the index is within bounds
func (o *Array) SetValue(idx int, value IObject) {
	if idx < 0 || idx >= len(o.values) {
		return
	}
	o.values[idx] = value
}

// Append adds an element to the end of the Array.
func (o *Array) Append(elem IObject) {
	if len(o.values) >= maxArrayLen {
		return
	}
	o.values = append(o.values, elem)
}

// Assign replaces the current slice of elements with the provided slice.
func (o *Array) Assign(v []IObject) {
	if len(v) > maxArrayLen {
		v = v[0:maxArrayLen]
	}
	o.values = v
}

// LogicalOp performs a logical operation on the array using the given operator and operand, returning an error if invalid.
func (o *Array) LogicalOp(_ int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		return logicalOpNil(o.gk, op)
	}
	return nil, ErrInvalidOperator
}

// ArithmeticOp performs arithmetic operations between the current Array and another Array based on the specified operator.
func (o *Array) ArithmeticOp(frame int, op ArithmeticOperator, rhsIn IObject) (IObject, error) {
	switch op {
	case OperatorAdd:
		switch rhs := rhsIn.(type) {
		case *Bool, *Int, *Float, *Char, *String:
			if len(o.values)+rhsIn.Length() > maxArrayLen {
				return nil, ErrLimitExceed
			}
			return o.gk.NewArray(frame, append(o.values, rhsIn)), nil
		case *Array:
			if len(rhs.values) == 0 {
				return o, nil
			}
			if len(o.values)+len(rhs.values) > maxArrayLen {
				return nil, ErrLimitExceed
			}
			return o.gk.NewArray(frame, append(o.values, rhs.values...)), nil
		}
	default:
		return nil, ErrInvalidOperator
	}
	return nil, ErrInvalidOperator
}

// Copy creates and returns a deep copy of the Array and its elements.
func (o *Array) Copy(frame int, depth int) IObject {
	var c []IObject
	for _, elem := range o.values {
		if depth >= maxDepth {
			break
		}
		c = append(c, elem.Copy(frame, depth+1))
	}
	return o.gk.NewArray(frame, c)
}

// Falsy returns true if the array is empty, otherwise false.
func (o *Array) Falsy() bool {
	return len(o.values) == 0
}

// Equals compare the current Array with another IObject and return true if they have equivalent values and order.
func (o *Array) Equals(in IObject) bool {
	var xVal []IObject
	switch x := in.(type) {
	case *Array:
		xVal = x.values
	default:
		return false
	}
	if len(o.values) != len(xVal) {
		return false
	}
	for i, e := range o.values {
		if !e.Equals(xVal[i]) {
			return false
		}
	}
	return true
}

// IndexGet retrieves the element at the given index from the Array. Returns an error if the index type is invalid or out of bounds.
func (o *Array) IndexGet(_ int, index IObject) (IObject, error) {
	intIdx, ok := index.(*Int)
	if !ok {
		return o.gk.UndefinedValue(), ErrInvalidIndexType
	}
	idxVal := int(intIdx.value)
	if idxVal < 0 || idxVal >= len(o.values) {
		return o.gk.UndefinedValue(), nil
	}
	return o.values[idxVal], nil
}

// IndexSet assigns a given value to the specified index in the array, returning an error if the operation is invalid.
func (o *Array) IndexSet(index IObject, value IObject) (err error) {
	idx, ok := o.GateKeeper().ToInt64(index)
	if !ok {
		err = ErrInvalidIndexType
		return
	}
	intIdx := int(idx)
	if intIdx < 0 || intIdx >= len(o.values) {
		err = ErrIndexOutOfBounds
		return
	}
	o.values[intIdx] = value
	return nil
}

// Count returns the total number of elements in the instance and its sub-elements.
func (o *Array) Count() int {
	counter := 0
	for _, v := range o.values {
		counter += v.Count()
	}
	return counter
}

// Iterable checks if the Array is iterable and always returns true.
func (o *Array) Iterable() bool {
	return true
}

// Iterate returns an IIterator for the Array instance, allowing sequential access to its elements.
func (o *Array) Iterate(frame int) IIterator {
	return o.GateKeeper().NewArrayIterator(frame, o.values, 0)
}
