package objects

import (
	"bytes"
	"encoding/gob"
)

// StructIteratorType represents the type for a struct iterator, used as a constant string identifier.
// StructIteratorLabel is a formatted label that includes the StructIteratorType constant within angle brackets.
const (
	StructIteratorType  = "struct_iterator"
	StructIteratorLabel = "<" + StructIteratorType + ">"
)

func init() {
	gob.Register(&StructIterator{})
}

// StructIterator is an iterator for traversing over the keys and Code of a struct-like map of IObjects.
// It embeds Object and implements the IIterator interface.
// The keys and Code are stored in separate slices to facilitate order-preserving iteration.
// The index tracks the current position within the iteration, and length is the total number of elements.
type StructIterator struct {
	IAllocator
	data   map[string]IObject
	keys   []string
	index  int
	length int
}

// NewStructIterator initializes and returns a new StructIterator for the given map of string keys to IObject Code.
func newStructIterator(allocator IAllocator, v map[string]IObject, index int) IIterator {
	var keys []string
	for k := range v {
		keys = append(keys, k)
	}
	return &StructIterator{
		IAllocator: allocator,
		data:       v,
		keys:       keys,
		length:     len(keys),
		index:      index,
	}
}

// setAllocator sets the allocator for the instance, defining its memory management and lifecycle behavior.
func (o *StructIterator) setAllocator(allocator IAllocator) {
	o.IAllocator = allocator
}

// GateKeeper returns a reference to the GateKeeper associated with the Object.
func (o *StructIterator) GateKeeper() IGateKeeper {
	return o.GateKeeper()
}

// AsInterface converts the object into a generic interface{} type and returns the underlying data.
func (o *StructIterator) AsInterface() interface{} {
	res := make(map[string]interface{})
	for key, v := range o.data {
		res[key] = o.GateKeeper().ToInterface(v)
	}
	return res
}

// AsBool returns true if the array is not empty, otherwise false.
func (o *StructIterator) AsBool() bool {
	return o.length > 0
}

// AsInt64 returns the length of the array as an int64 Code.
func (o *StructIterator) AsInt64() int64 {
	return int64(o.length)
}

// AsFloat64 returns the length of the array as an int64 Code.
func (o *StructIterator) AsFloat64() float64 {
	return float64(o.length)
}

// AsString returns the string representation of the StructIterator instance.
func (o *StructIterator) AsString() string {
	return StructIteratorLabel
}

// Nil checks if the object is nil and always returns false.
func (o *StructIterator) Nil() bool {
	return false
}

// AssignValue sets the current object to the provided IObject, returning ErrNotAssignable if the operation is not supported.
func (o *StructIterator) AssignValue(_ IObject) error {
	return ErrNotAssignable
}

// LogicalOp performs a logical operation between the StructIterator and another object using the specified operator.
// It returns ErrInvalidOperator as logical operations are unsupported for StructIterator.
func (o *StructIterator) LogicalOp(_ int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		return logicalOpNil(o.GateKeeper(), op)
	}
	return nil, ErrInvalidOperator
}

// ArithmeticOp performs an arithmetic operation specified by the Operator on the object and another IObject, returning the result or an error.
func (o *StructIterator) ArithmeticOp(_ int, _ ArithmeticOperator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// IndexGet attempts to retrieve a Code at the given index and returns an error if the object is not indexable.
func (o *StructIterator) IndexGet(_ int, _ IObject) (IObject, error) {
	return o.GateKeeper().UndefinedValue(), ErrIndexNotIndexable
}

// IndexSet attempts to assign a Code to an index in the object but always returns ErrIndexUnsupported,
// as this operation is unsupported.
func (o *StructIterator) IndexSet(_, _ IObject) (err error) {
	return ErrIndexUnsupported
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (o *StructIterator) Iterate(_ int) IIterator {
	return nil
}

// Iterable determines if the object can be iterated over and returns false for this implementation.
func (o *StructIterator) Iterable() bool {
	return false
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *StructIterator) Call(_ int, _ ...IObject) (retCount uint, ret IObject, err error) {
	return 0, nil, nil
}

// Length returns the length of the Int object.
func (o *StructIterator) Length() int {
	return 0
}

// Copy creates and returns a duplicate instance of the current StructIterator with the same internal state.
func (o *StructIterator) Copy(frame int, _ int) IObject {
	ret := o.GateKeeper().NewStructIterator(frame, o.data, o.index)
	return ret
}

// TypeName returns the type name of the StructIterator object as a string.
func (o *StructIterator) TypeName() string {
	return StructIteratorType
}

// Falsy returns true, indicating the current StructIterator is truthy.
func (o *StructIterator) Falsy() bool {
	return true
}

// Equals compares the current StructIterator with another IObject for equality and returns true if they are equal.
func (o *StructIterator) Equals(IObject) bool {
	return false
}

// Next advances the iterator to the next position and returns true if there are more elements to iterate over.
func (o *StructIterator) Next() bool {
	o.index++
	return o.index <= o.length
}

// Key retrieves the key of the current element in the MapIterator as an IObject.
func (o *StructIterator) Key(frame int) IObject {
	idx := o.index - 1
	if idx < 0 || idx >= o.length {
		return o.GateKeeper().UndefinedValue()
	}
	k := o.keys[idx]
	return o.GateKeeper().NewString(frame, k)
}

// Value retrieves the Code of the current element in the iteration based on the iterator's current position.
func (o *StructIterator) Value(_ int) IObject {
	idx := o.index - 1
	if idx < 0 || idx >= o.length {
		return o.GateKeeper().UndefinedValue()
	}
	k := o.keys[idx]
	return o.data[k]
}

// Count returns the total number of elements in the instance and its sub-elements.
func (o *StructIterator) Count() int {
	counter := 0
	for _, v := range o.data {
		counter += v.Count()
	}
	return counter
}

// GobEncode serializes the ArrayIterator's data into a byte slice using gob encoding and returns the result or an error.
func (o *StructIterator) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(o.data); err != nil {
		return nil, err
	}
	if err := encoder.Encode(o.index); err != nil {
		return nil, err
	}
	if err := encoder.Encode(o.keys); err != nil {
		return nil, err
	}
	if err := encoder.Encode(o.length); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GobDecode decodes the provided byte slice into the StructIterator's data field using the gob package.
func (o *StructIterator) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	if err := decoder.Decode(&o.data); err != nil {
		return err
	}
	if err := decoder.Decode(&o.index); err != nil {
		return err
	}
	if err := decoder.Decode(&o.keys); err != nil {
		return err
	}
	if err := decoder.Decode(&o.length); err != nil {
		return err
	}
	return nil
}
