package objects

import (
	"bytes"
	"encoding/gob"
)

const (
	MapIteratorType  = "map_iterator"
	MapIteratorLabel = "<" + MapIteratorType + ">"
)

func init() {
	gob.Register(&MapIterator{})
}

// MapIterator is a type used for iterating over key-Code pairs in a map-like structure.
// It implements the IIterator interface and provides methods for traversal and element access.
// The embedded Object provides default implementations for methods from the IObject interface.
// The internal state includes the map (Code), keys (keys), current position index (index), and total keys length (length).
type MapIterator struct {
	IAllocator
	data   map[string]IObject
	keys   []string
	index  int
	length int
}

// NewMapIterator creates and returns a new instance of MapIterator.
func newMapIterator(allocator IAllocator, v map[string]IObject, index int) IIterator {
	var keys []string
	for k := range v {
		keys = append(keys, k)
	}
	return &MapIterator{
		IAllocator: allocator,
		data:       v,
		keys:       keys,
		length:     len(keys),
		index:      index,
	}
}

// setAllocator sets the allocator for the instance, defining its memory management and lifecycle behavior.
func (o *MapIterator) setAllocator(allocator IAllocator) {
	o.IAllocator = allocator
}

// AsInterface converts the object into a generic interface{} type and returns the underlying data.
func (o *MapIterator) AsInterface() interface{} {
	res := make(map[string]interface{})
	for key, v := range o.data {
		res[key] = v.AsInterface()
	}
	return res
}

// AsBool returns true if the map is not empty, otherwise false.
func (o *MapIterator) AsBool() bool {
	return o.length > 0
}

// AsInt64 returns the length of the array as an int64 Code.
func (o *MapIterator) AsInt64() int64 {
	return int64(o.length)
}

// AsFloat64 returns the length of the array as an int64 Code.
func (o *MapIterator) AsFloat64() float64 {
	return float64(o.length)
}

// AsBytes converts the object elements into a single concatenated slice of bytes by calling AsBytes on each element.
func (o *MapIterator) AsBytes() []byte {
	var res []byte
	for _, v := range o.data {
		res = append(res, v.AsBytes()...)
	}
	return res
}

// AsString returns the string representation of the MapIterator.
func (o *MapIterator) AsString() string {
	return MapIteratorLabel
}

// AssignValue sets the current object to the provided IObject, returning ErrNotAssignable if the operation is not supported.
func (o *MapIterator) AssignValue(_ IObject) error {
	return ErrNotAssignable
}

// Nil checks if the object is nil and always returns false.
func (o *MapIterator) Nil() bool {
	return false
}

// LogicalOp performs a logical operation on the MapIterator object using the specified operator and operand.
// Returns nil and ErrInvalidOperator as logical operations are unsupported for MapIterator.
func (o *MapIterator) LogicalOp(_ int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		return logicalOpNil(o.GateKeeper(), op)
	}
	return nil, ErrInvalidOperator
}

// ArithmeticOp performs an arithmetic operation specified by the operator on the current object and the given operand.
// Returns an error if the operation is unsupported.
func (o *MapIterator) ArithmeticOp(_ int, _ ArithmeticOperator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// IndexGet attempts to retrieve a Code at the given index and returns an error if the object is not indexable.
func (o *MapIterator) IndexGet(_ int, _ IObject) (IObject, error) {
	return o.GateKeeper().UndefinedValue(), ErrIndexNotIndexable
}

// IndexSet attempts to assign a Code to an index in the object but always returns ErrIndexUnsupported,
// as this operation is unsupported.
func (o *MapIterator) IndexSet(_, _ IObject) error {
	return ErrIndexUnsupported
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (o *MapIterator) Iterate(_ int) IIterator {
	return nil
}

// Iterable determines if the object can be iterated over and returns false for this implementation.
func (o *MapIterator) Iterable() bool {
	return false
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *MapIterator) Call(_ int, _ ...IObject) (retCount uint, ret IObject, err error) {
	return 0, nil, nil
}

// Length returns the length of the Int object.
func (o *MapIterator) Length() int {
	return 0
}

// Copy creates and returns a new instance of MapIterator, duplicating its current state.
func (o *MapIterator) Copy(frame int, _ int) IObject {
	ret := o.GateKeeper().NewMapIterator(frame, o.data, o.index)
	return ret
}

// TypeName returns the type name of the MapIterator as a string.
func (o *MapIterator) TypeName() string {
	return MapIteratorType
}

// Falsy returns true, indicating the MapIterator is considered falsy in a boolean context.
func (o *MapIterator) Falsy() bool {
	return true
}

// Equals determine if the current MapIterator is equal to another IObject. Returns false by default.
func (o *MapIterator) Equals(IObject) bool {
	return false
}

// Next advances the iterator to the next element and returns true if there are more elements to iterate over.
func (o *MapIterator) Next() bool {
	o.index++
	return o.index <= o.length
}

// Key retrieves the key of the current element in the iteration as an IObject.
func (o *MapIterator) Key(frame int) IObject {
	idx := o.index - 1
	if idx < 0 || idx >= o.length {
		return o.GateKeeper().UndefinedValue()
	}
	k := o.keys[idx]
	return o.GateKeeper().NewString(frame, k)
}

// Value retrieves the Code of the current element in the iteration based on the iterator's current position.
func (o *MapIterator) Value(_ int) IObject {
	idx := o.index - 1
	if idx < 0 || idx >= o.length {
		return o.GateKeeper().UndefinedValue()
	}
	k := o.keys[idx]
	return o.data[k]
}

// Count returns the total number of elements in the instance and its sub-elements.
func (o *MapIterator) Count() int {
	counter := 0
	for _, v := range o.data {
		counter += v.Count()
	}
	return counter
}

// GobEncode serializes the ArrayIterator's data into a byte slice using gob encoding and returns the result or an error.
func (o *MapIterator) GobEncode() ([]byte, error) {
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

// GobDecode decodes the provided byte slice into the MapIterator's data field using the gob package.
func (o *MapIterator) GobDecode(data []byte) error {
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
