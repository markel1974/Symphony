package objects

import (
	"bytes"
	"encoding/gob"
	"reflect"
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
	IAllocator
	data []IObject
}

// NewArray creates and returns a new Array object initialized with the provided slice of IObject elements.
func newArray(allocator IAllocator, value []IObject) IObject {
	if len(value) > MaxArrayLen {
		value = value[0:MaxArrayLen]
	}
	return &Array{
		IAllocator: allocator,
		data:       value,
	}
}

// setAllocator sets the allocator for the instance, defining its memory management and lifecycle behavior.
func (o *Array) setAllocator(allocator IAllocator) {
	o.IAllocator = allocator
}

// Setup initializes the Array with a specific frame and truncates the input slice if it exceeds MaxArrayLen.
func (o *Array) Setup(frame int, v []IObject) {
	o.setFrame(frame)
	if len(v) > MaxArrayLen {
		v = v[0:MaxArrayLen]
	}
	o.data = v
}

// AsInterface converts the object into a generic interface{} type and returns the underlying data.
func (o *Array) AsInterface() interface{} {
	res := make([]interface{}, len(o.data))
	for i, val := range o.data {
		res[i] = val.AsInterface()
	}
	return res
}

// AsValue converts the Array's elements to a slice of reflect.Value and returns it as a reflect.Value.
func (o *Array) AsValue(target reflect.Type) (reflect.Value, bool) {
	return _reflect(o, target)
}

// AsBool returns true if the array is not empty, otherwise false.
func (o *Array) AsBool() bool {
	return len(o.data) > 0
}

// AsInt64 returns the len of the array as an int64 data.
func (o *Array) AsInt64() int64 {
	return int64(len(o.data))
}

// AsFloat64 returns the len of the array as an int64 data.
func (o *Array) AsFloat64() float64 {
	return float64(len(o.data))
}

// AsBytes converts the object elements into a single concatenated slice of bytes by calling AsBytes on each element.
func (o *Array) AsBytes() []byte {
	var elements []byte
	for _, e := range o.data {
		elements = append(elements, e.AsBytes()...)
	}
	return elements
}

// AsString returns a string representation of the Array, formatting its elements in a comma-separated list enclosed in brackets.
func (o *Array) AsString() string {
	var elements []string
	for _, e := range o.data {
		elements = append(elements, e.AsString())
	}
	return strings.Join(elements, ";")
}

// AssignValue assigns the elements of another Array to the current Array if the input is of type *Array, otherwise returns an error.
func (o *Array) AssignValue(v IObject) error {
	target, ok := v.(*Array)
	if !ok {
		return ErrNotAssignable
	}
	o.Assign(target.data)
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

// ForEach applies the provided function `fn` to each element of the array, passing the index and the object as arguments.
func (o *Array) ForEach(fn func(x int, obj IObject)) {
	for x, v := range o.data {
		fn(x, v)
	}
}

// CopyRange creates and returns a subset of the Array's data as a new slice from the specified start to end indices.
func (o *Array) CopyRange(start uint, end uint) []IObject {
	if start > end || end > uint(len(o.data)) {
		return []IObject{}
	}
	ret := make([]IObject, end-start)
	copy(ret, o.data[start:end])
	return ret
}

// Index returns the element at the specified Index in the Array, or an error if the Index is out of bounds.
func (o *Array) Index(index int) (IObject, error) {
	if index < 0 || index >= len(o.data) {
		return nil, ErrIndexOutOfBounds
	}
	return o.data[index], nil
}

// Length returns the number of elements in the Array.
func (o *Array) Length() int {
	return len(o.data)
}

// SetValue assigns a given IObject data to the specified Index in the Array if the Index is within bounds
func (o *Array) SetValue(idx int, elem IObject) {
	if idx < 0 || idx >= len(o.data) {
		return
	}
	elem.AddRef()
	o.data[idx] = elem
}

// Append adds an element to the end of the Array.
func (o *Array) Append(elem IObject) {
	if len(o.data) >= MaxArrayLen {
		return
	}
	elem.AddRef()
	o.data = append(o.data, elem)
}

// Assign replaces the current slice of elements with the provided slice.
func (o *Array) Assign(v []IObject) {
	if len(v) > MaxArrayLen {
		v = v[0:MaxArrayLen]
	}
	o.data = v
}

// LogicalOp performs a logical operation on the array using the given operator and operand, returning an error if invalid.
func (o *Array) LogicalOp(_ int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		ret, err := logicalOpNil(op)
		return o.GateKeeper().FromBoolError(ret, err)
	}
	return nil, ErrInvalidOperator
}

// ArithmeticOp performs arithmetic operations between the current Array and another Array based on the specified operator.
func (o *Array) ArithmeticOp(frame int, op ArithmeticOperator, rhsIn IObject) (IObject, error) {
	switch op {
	case OperatorAdd:
		switch rhs := rhsIn.(type) {
		case *Bool, *Int, *Float, *Char, *String:
			if len(o.data)+rhsIn.Length() > MaxArrayLen {
				return nil, ErrLimitExceed
			}
			return o.GateKeeper().NewArray(frame, append(o.data, rhsIn)), nil
		case *Array:
			if len(rhs.data) == 0 {
				return o, nil
			}
			if len(o.data)+len(rhs.data) > MaxArrayLen {
				return nil, ErrLimitExceed
			}
			return o.GateKeeper().NewArray(frame, append(o.data, rhs.data...)), nil
		}
	default:
		return nil, ErrInvalidOperator
	}
	return nil, ErrInvalidOperator
}

// UnaryOp applies a unary operation on the Array and returns an error indicating an invalid operation.
func (o *Array) UnaryOp(_ int, _ UnaryOperator) (IObject, error) {
	return nil, ErrInvalidOperator
}

// Copy creates and returns a deep copy of the Array and its elements.
func (o *Array) Copy(frame int, depth int) IObject {
	var c []IObject
	for _, elem := range o.data {
		if depth >= MaxDepth {
			break
		}
		c = append(c, elem.Copy(frame, depth+1))
	}
	return o.GateKeeper().NewArray(frame, c)
}

// Falsy returns true if the array is empty, otherwise false.
func (o *Array) Falsy() bool {
	return len(o.data) == 0
}

// Equals compare the current Array with another IObject and return true if they have equivalent data and order.
func (o *Array) Equals(in IObject) bool {
	var xVal []IObject
	switch x := in.(type) {
	case *Array:
		xVal = x.data
	default:
		return false
	}
	if len(o.data) != len(xVal) {
		return false
	}
	for i, e := range o.data {
		if !e.Equals(xVal[i]) {
			return false
		}
	}
	return true
}

// IndexGet retrieves the element at the given Index from the Array. Returns an error if the Index type is invalid or out of bounds.
func (o *Array) IndexGet(_ int, index IObject) (IObject, error) {
	idx := index.AsInt64()
	if idx < 0 || idx >= int64(len(o.data)) {
		return o.GateKeeper().UndefinedValue(), ErrIndexOutOfBounds
	}
	return o.data[idx], nil
}

// IndexSet assigns a given data to the specified Index in the array, returning an error if the operation is invalid.
func (o *Array) IndexSet(index IObject, value IObject) error {
	idx := index.AsInt64()
	if idx < 0 || idx >= int64(len(o.data)) {
		return ErrIndexOutOfBounds
	}
	o.data[idx] = value
	return nil
}

// Count returns the total number of elements in the instance and its sub-elements.
func (o *Array) Count() int {
	counter := 0
	for _, v := range o.data {
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
	return o.GateKeeper().NewArrayIterator(frame, o.data, 0)
}

// GobEncode serializes the Array's data into a byte slice using gob encoding and returns the result or an error.
func (o *Array) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(o.data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GobDecode decodes the provided byte slice into the Array's data field using the gob package.
func (o *Array) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	if err := decoder.Decode(&o.data); err != nil {
		return err
	}
	return nil
}
