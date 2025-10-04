package objects

import (
	"bytes"
	"encoding/gob"
	"reflect"
)

const (
	BytesIteratorType = "bytes_iterator"
)

func init() {
	gob.Register(&BytesIterator{})
}

// BytesIterator is an iterator for traversing elements of a byte slice, implementing the IIterator interface.
type BytesIterator struct {
	IAllocator
	data   []byte
	index  int
	length int
}

// NewBytesIterator creates and returns a new BytesIterator instance with the given slice of bytes.
func newBytesIterator(allocator IAllocator, v []byte, index int) IIterator {
	return &BytesIterator{
		IAllocator: allocator,
		data:       v,
		length:     len(v),
		index:      index,
	}
}

// setAllocator sets the allocator for the instance, defining its memory management and lifecycle behavior.
func (o *BytesIterator) setAllocator(allocator IAllocator) {
	o.IAllocator = allocator
}

// AsInterface converts the object into a generic interface{} type and returns the underlying data.
func (o *BytesIterator) AsInterface() interface{} {
	return o.data
}

// AsValue attempts to convert the BytesIterator data to a reflect.Value of the specified target type and returns it with a success flag.
func (o *BytesIterator) AsValue(target reflect.Type) (reflect.Value, bool) {
	if target.Kind() == reflect.ValueOf(o.data).Kind() {
		return reflect.ValueOf(o.data), true
	}
	return reflect.Value{}, false
}

// AsBool returns true if the object is not empty, otherwise false.
func (o *BytesIterator) AsBool() bool {
	return o.length > 0
}

// AsInt64 returns the length of the array as an int64 data.
func (o *BytesIterator) AsInt64() int64 {
	return int64(o.length)
}

// AsFloat64 returns the length of the array as an int64 data.
func (o *BytesIterator) AsFloat64() float64 {
	return float64(o.length)
}

// AsBytes converts the object elements into a single concatenated slice of bytes by calling AsBytes on each element.
func (o *BytesIterator) AsBytes() []byte {
	return o.data
}

// AsString returns the string representation of the BytesIterator.
func (o *BytesIterator) AsString() string {
	return ""
}

// AssignValue sets the current object to the provided IObject, returning ErrNotAssignable if the operation is not supported.
func (o *BytesIterator) AssignValue(_ IObject) error {
	return ErrNotAssignable
}

// Nil checks if the object is nil and always returns false.
func (o *BytesIterator) Nil() bool {
	return false
}

// LogicalOp performs a logical operation on the BytesIterator object, always returning an error for invalid operations.
func (o *BytesIterator) LogicalOp(_ int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		ret, err := logicalOpNil(op)
		return o.GateKeeper().FromBoolError(ret, err)
	}
	return nil, ErrInvalidOperator
}

// ArithmeticOp performs an arithmetic operation with the given operator and operands; returns an error for invalid operators.
func (o *BytesIterator) ArithmeticOp(_ int, _ ArithmeticOperator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// UnaryOp applies a unary operation and returns an error indicating an invalid operation.
func (o *BytesIterator) UnaryOp(_ int, _ UnaryOperator) (IObject, error) {
	return nil, ErrInvalidOperator
}

// Falsy returns false for all objects.
func (o *BytesIterator) Falsy() bool {
	return false
}

// IndexGet attempts to retrieve a data at the given index and returns an error if the object is not indexable.
func (o *BytesIterator) IndexGet(_ int, _ IObject) (IObject, error) {
	return o.GateKeeper().UndefinedValue(), ErrIndexNotIndexable
}

// IndexSet attempts to assign a data to an index in the object but always returns ErrIndexUnsupported,
// as this operation is unsupported.
func (o *BytesIterator) IndexSet(_, _ IObject) error {
	return ErrIndexUnsupported
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (o *BytesIterator) Iterate(_ int) IIterator {
	return nil
}

// Iterable determines if the object can be iterated over and returns false for this implementation.
func (o *BytesIterator) Iterable() bool {
	return false
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *BytesIterator) Call(_ int, _ ...IObject) (retCount uint, ret IObject, err error) {
	return 0, nil, nil
}

// Length returns the length of the Int object.
func (o *BytesIterator) Length() int {
	return 0
}

// TypeName returns the string representation of the type name, which is "bytes-iterator".
func (o *BytesIterator) TypeName() string {
	return BytesIteratorType
}

// Equals checks whether the BytesIterator is equal to another object implementing the IObject interface.
func (o *BytesIterator) Equals(IObject) bool {
	return false
}

// Copy creates and returns a new instance of BytesIterator with the same state as the current instance.
func (o *BytesIterator) Copy(frame int, _ int) IObject {
	ret := o.GateKeeper().NewBytesIterator(frame, o.data, o.index)
	return ret
}

// Next advances the iterator to the next position and returns true if the current position is within bounds.
func (o *BytesIterator) Next() bool {
	o.index++
	return o.index <= o.length
}

// Key returns the current index of the iterator as an IObject, decremented by one from the internal index tracker.
func (o *BytesIterator) Key(frame int) IObject {
	idx := o.index - 1
	if idx < 0 || idx >= o.length {
		return o.GateKeeper().UndefinedValue()
	}
	return o.GateKeeper().NewInt(frame, int64(idx))
}

// Value returns the data of the current byte in the iteration as an IObject, wrapped in an Int struct.
func (o *BytesIterator) Value(frame int) IObject {
	idx := o.index - 1
	if idx < 0 || idx >= o.length {
		return o.GateKeeper().UndefinedValue()
	}
	return o.GateKeeper().NewInt(frame, int64(o.data[idx]))
}

// Count returns the total number of elements in the instance and its sub-elements.
func (o *BytesIterator) Count() int {
	return 1
}

// GobEncode serializes the ArrayIterator's data into a byte slice using gob encoding and returns the result or an error.
func (o *BytesIterator) GobEncode() ([]byte, error) {
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
func (o *BytesIterator) GobDecode(data []byte) error {
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
