package objects

import "encoding/gob"

// StructIteratorType represents the type for a struct iterator, used as a constant string identifier.
// StructIteratorLabel is a formatted label that includes the StructIteratorType constant within angle brackets.
const (
	StructIteratorType  = "struct_iterator"
	StructIteratorLabel = "<" + StructIteratorType + ">"
)

func init() {
	gob.Register(&StructIterator{})
}

// StructIterator is an iterator for traversing over the keys and values of a struct-like map of IObjects.
// It embeds Object and implements the IIterator interface.
// The keys and values are stored in separate slices to facilitate order-preserving iteration.
// The index tracks the current position within the iteration, and length is the total number of elements.
type StructIterator struct {
	factory IGateKeeper
	frame   int
	values  map[string]IObject
	keys    []string
	index   int
	length  int
}

// NewStructIterator initializes and returns a new StructIterator for the given map of string keys to IObject values.
func newStructIterator(factory IGateKeeper, frame int, v map[string]IObject, index int) IIterator {
	var keys []string
	for k := range v {
		keys = append(keys, k)
	}
	return &StructIterator{
		factory: factory,
		frame:   frame,
		values:  v,
		keys:    keys,
		length:  len(keys),
		index:   index,
	}
}

// GateKeeper returns a reference to the GateKeeper associated with the Object.
func (o *StructIterator) GateKeeper() IGateKeeper {
	return o.factory
}

// AsBool returns true if the array is not empty, otherwise false.
func (o *StructIterator) AsBool() bool {
	return o.length > 0
}

// AsInt64 returns the length of the array as an int64 value.
func (o *StructIterator) AsInt64() int64 {
	return int64(o.length)
}

// AsFloat64 returns the length of the array as an int64 value.
func (o *StructIterator) AsFloat64() float64 {
	return float64(o.length)
}

// AsString returns the string representation of the StructIterator instance.
func (o *StructIterator) AsString() string {
	return StructIteratorLabel
}

// AssignValue sets the current object to the provided IObject, returning ErrNotAssignable if the operation is not supported.
func (o *StructIterator) AssignValue(_ IObject) error {
	return ErrNotAssignable
}

// SetStatic sets the frame to FrameStatic, marking it with a static execution context.
func (o *StructIterator) SetStatic() {
	o.frame = FrameStatic
}

// Frame returns the current frame value of the Object.
func (o *StructIterator) Frame() int {
	return o.frame
}

// LogicalOp performs a logical operation between the StructIterator and another object using the specified operator.
// It returns ErrInvalidOperator as logical operations are unsupported for StructIterator.
func (o *StructIterator) LogicalOp(_ int, _ LogicalOperator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// ArithmeticOp performs an arithmetic operation specified by the Operator on the object and another IObject, returning the result or an error.
func (o *StructIterator) ArithmeticOp(_ int, _ ArithmeticOperator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// IndexGet attempts to retrieve a value at the given index and returns an error if the object is not indexable.
func (o *StructIterator) IndexGet(_ int, _ IObject) (res IObject, err error) {
	return nil, ErrNotIndexable
}

// IndexSet attempts to assign a value to an index in the object but always returns ErrUnsupportedIndex,
// as this operation is unsupported.
func (o *StructIterator) IndexSet(_, _ IObject) (err error) {
	return ErrUnsupportedIndex
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (o *StructIterator) Iterate(_ int) IIterator {
	return nil
}

// CanIterate determines if the object can be iterated over and returns false for this implementation.
func (o *StructIterator) CanIterate() bool {
	return false
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *StructIterator) Call(_ int, _ ...IObject) (ret IObject, err error) {
	return nil, nil
}

// CanCall determines if the object can be invoked as a callable. Returns false for non-callable objects.
func (o *StructIterator) CanCall() bool {
	return false
}

// Length returns the length of the Int object.
func (o *StructIterator) Length() int {
	return 0
}

// Copy creates and returns a duplicate instance of the current StructIterator with the same internal state.
func (o *StructIterator) Copy(frame int, _ int) IObject {
	ret := o.GateKeeper().NewStructIterator(frame, o.values, o.index)
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

// Value retrieves the value of the current element in the iteration based on the iterator's current position.
func (o *StructIterator) Value(_ int) IObject {
	idx := o.index - 1
	if idx < 0 || idx >= o.length {
		return o.GateKeeper().UndefinedValue()
	}
	k := o.keys[idx]
	return o.values[k]
}
