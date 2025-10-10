package objects

import (
	"bytes"
	"encoding/gob"
	"reflect"
)

const (
	ArrayIteratorType = "array_iterator"
)

func init() {
	gob.Register(&ArrayIterator{})
}

// ArrayIterator is an iterator type for traversing elements of an array.
// It implements the IIterator interface to provide sequential access to array elements.
type ArrayIterator struct {
	IAllocator
	data   []IObject
	index  int
	length int
}

// NewArrayIterator creates and returns a new ArrayIterator instance with the given slice of IObject.
func newArrayIterator(allocator IAllocator, v []IObject, index int) IIterator {
	return &ArrayIterator{
		IAllocator: allocator,
		data:       v,
		length:     len(v),
		index:      index,
	}
}

// setAllocator sets the allocator for the instance, defining its memory management and lifecycle behavior.
func (o *ArrayIterator) setAllocator(allocator IAllocator) {
	o.IAllocator = allocator
}

// Setup initializes the ArrayIterator with the given frame, data slice, starting index, and calculates its length.
func (o *ArrayIterator) Setup(frame int, v []IObject, index int) {
	o.setFrame(frame)
	o.data = v
	o.index = index
	o.length = len(v)
}

// AsInterface converts the object into a generic interface{} type and returns the underlying data.
func (o *ArrayIterator) AsInterface() interface{} {
	res := make([]interface{}, len(o.data))
	for i, val := range o.data {
		res[i] = val.AsInterface()
	}
	return res
}

// AsValue attempts to convert the ArrayIterator's data into a reflect.Value of the specified type. Returns false if invalid.
func (o *ArrayIterator) AsValue(target reflect.Type) (reflect.Value, bool) {
	return _reflect(o, target)
}

// AsBool returns true if the array is not empty, otherwise false.
func (o *ArrayIterator) AsBool() bool {
	return o.length > 0
}

// AsInt64 returns the len of the array as an int64 data.
func (o *ArrayIterator) AsInt64() int64 {
	return int64(o.length)
}

// AsFloat64 returns the len of the array as an int64 data.
func (o *ArrayIterator) AsFloat64() float64 {
	return float64(o.length)
}

// AsBytes converts the object elements into a single concatenated slice of bytes by calling AsBytes on each element.
func (o *ArrayIterator) AsBytes() []byte {
	var elements []byte
	for _, e := range o.data {
		elements = append(elements, e.AsBytes()...)
	}
	return elements
}

// AsString returns a string representation of the ArrayIterator instance.
func (o *ArrayIterator) AsString() string {
	return ""
}

// AssignValue assigns the elements of another Array to the current Array if the input is of type *Array, otherwise returns an error.
func (o *ArrayIterator) AssignValue(_ IObject) error {
	return ErrNotAssignable
}

// Nil checks if the object is nil and always returns false.
func (o *ArrayIterator) Nil() bool {
	return false
}

// Call invokes the ArrayIterator as a callable function with the provided frame and arguments, returning a result or error.
func (o *ArrayIterator) Call(_ int, _ ...IObject) (retCount uint, ret IObject, err error) {
	return 0, nil, nil
}

// LogicalOp attempts to perform a logical operation, returning an error as this operation is not supported.
func (o *ArrayIterator) LogicalOp(_ int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		ret, err := logicalOpNil(op)
		return o.GateKeeper().FromBoolError(ret, err)
	}
	return nil, ErrInvalidOperator
}

// ArithmeticOp performs an arithmetic operation on the ArrayIterator using a specified operator and operand, returning an error.
func (o *ArrayIterator) ArithmeticOp(_ int, _ ArithmeticOperator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// UnaryOp applies a unary operation and returns an error indicating an invalid operation.
func (o *ArrayIterator) UnaryOp(_ int, _ UnaryOperator) (IObject, error) {
	return nil, ErrInvalidOperator
}

// IndexGet retrieves an element at a specific index from the array but always returns ErrIndexNotIndexable for this implementation.
func (o *ArrayIterator) IndexGet(_ int, _ IObject) (IObject, error) {
	return o.GateKeeper().UndefinedValue(), ErrIndexNotIndexable
}

// IndexSet attempts to assign a data to an index in the object but always returns ErrIndexUnsupported,
// as this operation is unsupported.
func (o *ArrayIterator) IndexSet(_, _ IObject) (err error) {
	return ErrIndexUnsupported
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (o *ArrayIterator) Iterate(_ int) IIterator {
	return nil
}

// Iterable determines if the object can be iterated over and returns false for this implementation.
func (o *ArrayIterator) Iterable() bool {
	return false
}

// Length returns the len of the Int object.
func (o *ArrayIterator) Length() int {
	return 0
}

// TypeName returns the type name of the ArrayIterator as a string.
func (o *ArrayIterator) TypeName() string {
	return ArrayIteratorType
}

// Falsy determines whether the ArrayIterator should be considered a falsy data. Always returns true.
func (o *ArrayIterator) Falsy() bool {
	return true
}

// Equals checks whether the given IObject is equal to the current ArrayIterator instance by data comparison.
func (o *ArrayIterator) Equals(IObject) bool {
	return false
}

// Copy creates and returns a duplicate of the ArrayIterator, preserving its current state.
func (o *ArrayIterator) Copy(frame int, _ int) IObject {
	ret := o.GateKeeper().NewArrayIterator(frame, o.data, o.index)
	return ret
}

// Next advances the iterator to the next element and returns true if the current position is within bounds.
func (o *ArrayIterator) Next() bool {
	o.index++
	return o.index <= o.length
}

// Key returns the index of the current element in the iteration as an IObject.
func (o *ArrayIterator) Key(frame int) IObject {
	idx := int64(o.index - 1)
	if idx < 0 || idx >= int64(o.length) {
		return o.GateKeeper().UndefinedValue()
	}
	return o.GateKeeper().NewInt(frame, idx)
}

// Value returns the current element in the iteration based on the iterator's internal position.
func (o *ArrayIterator) Value(_ int) IObject {
	idx := int64(o.index - 1)
	if idx < 0 || idx >= int64(o.length) {
		return o.GateKeeper().UndefinedValue()
	}
	return o.data[idx]
}

// Count returns the total number of elements in the instance and its sub-elements.
func (o *ArrayIterator) Count() int {
	counter := 0
	for _, v := range o.data {
		counter += v.Count()
	}
	return counter
}

// GobEncode serializes the ArrayIterator's data into a byte slice using gob encoding and returns the result or an error.
func (o *ArrayIterator) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(o.data); err != nil {
		return nil, err
	}
	if err := encoder.Encode(o.index); err != nil {
		return nil, err
	}
	if err := encoder.Encode(o.length); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GobDecode decodes the provided byte slice into the ArrayIterator's data field using the gob package.
func (o *ArrayIterator) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	if err := decoder.Decode(&o.data); err != nil {
		return err
	}
	if err := decoder.Decode(&o.index); err != nil {
		return err
	}
	if err := decoder.Decode(&o.length); err != nil {
		return err
	}
	return nil
}
