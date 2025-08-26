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

// Frame returns the current frame value of the Object.
func (o *StructIterator) Frame() int {
	return o.frame
}

// BinaryOp performs a binary operation on the current object and another object using the specified operator.
// Returns the result of the operation or an error if the operation is not supported.
func (o *StructIterator) BinaryOp(_ int, _ Operator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// IndexGet attempts to retrieve a value at the given index and returns an error if the object is not indexable.
func (o *StructIterator) IndexGet(_ int, _ IObject) (res IObject, err error) {
	return nil, ErrNotIndexable
}

// IndexSet attempts to assign a value to an index in the object but always returns ErrNotIndexAssignable,
// as this operation is unsupported.
func (o *StructIterator) IndexSet(_, _ IObject) (err error) {
	return ErrNotIndexAssignable
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
func (i *StructIterator) Copy(frame int, _ int) IObject {
	ret := i.GateKeeper().NewStructIterator(frame, i.values, i.index)
	return ret
}

// TypeName returns the type name of the StructIterator object as a string.
func (i *StructIterator) TypeName() string {
	return StructIteratorType
}

// String returns the string representation of the StructIterator instance.
func (i *StructIterator) String() string {
	return StructIteratorLabel
}

// Boolean returns true, indicating the current StructIterator is truthy.
func (i *StructIterator) Boolean() bool {
	return true
}

// Equals compares the current StructIterator with another IObject for equality and returns true if they are equal.
func (i *StructIterator) Equals(IObject) bool {
	return false
}

// Next advances the iterator to the next position and returns true if there are more elements to iterate over.
func (i *StructIterator) Next() bool {
	i.index++
	return i.index <= i.length
}

// Key retrieves the key of the current element in the MapIterator as an IObject.
func (i *StructIterator) Key(frame int) IObject {
	idx := i.index - 1
	if idx < 0 || idx >= i.length {
		return i.GateKeeper().UndefinedValue()
	}
	k := i.keys[idx]
	return i.GateKeeper().NewString(frame, k)
}

// Value retrieves the value of the current element in the iteration based on the iterator's current position.
func (i *StructIterator) Value(_ int) IObject {
	idx := i.index - 1
	if idx < 0 || idx >= i.length {
		return i.GateKeeper().UndefinedValue()
	}
	k := i.keys[idx]
	return i.values[k]
}
