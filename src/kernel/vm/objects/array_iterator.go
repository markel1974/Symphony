package objects

import "encoding/gob"

const (
	ArrayIteratorType  = "array_iterator"
	ArrayIteratorLabel = "<" + ArrayIteratorType + ">"
)

func init() {
	gob.Register(&ArrayIterator{})
}

// ArrayIterator is an iterator type for traversing elements of an array.
// It implements the IIterator interface to provide sequential access to array elements.
type ArrayIterator struct {
	gk     IGateKeeper
	frame  int
	values []IObject
	index  int
	length int
}

// NewArrayIterator creates and returns a new ArrayIterator instance with the given slice of IObject.
func newArrayIterator(factory IGateKeeper, frame int, v []IObject, index int) IIterator {
	return &ArrayIterator{
		gk:     factory,
		frame:  frame,
		values: v,
		length: len(v),
		index:  index,
	}
}

// GateKeeper returns the GateKeeper instance associated with the ArrayIterator.
func (o *ArrayIterator) GateKeeper() IGateKeeper {
	return o.gk
}

// Frame returns the current execution frame associated with the ArrayIterator.
func (o *ArrayIterator) Frame() int {
	return o.frame
}

// Call invokes the ArrayIterator as a callable function with the provided frame and arguments, returning a result or error.
func (o *ArrayIterator) Call(_ int, _ ...IObject) (ret IObject, err error) {
	return nil, nil
}

// CanCall checks whether the ArrayIterator can be called as a function. Always returns false.
func (o *ArrayIterator) CanCall() bool {
	return false
}

// BinaryOp attempts a binary operation on the ArrayIterator but always returns nil and ErrInvalidOperator.
func (o *ArrayIterator) BinaryOp(_ int, _ Operator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// IndexGet retrieves an element at a specific index from the array but always returns ErrNotIndexable for this implementation.
func (o *ArrayIterator) IndexGet(_ int, _ IObject) (res IObject, err error) {
	return nil, ErrNotIndexable
}

// IndexSet attempts to assign a value to an index in the object but always returns ErrNotIndexAssignable,
// as this operation is unsupported.
func (o *ArrayIterator) IndexSet(_, _ IObject) (err error) {
	return ErrNotIndexAssignable
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (o *ArrayIterator) Iterate(_ int) IIterator {
	return nil
}

// CanIterate determines if the object can be iterated over and returns false for this implementation.
func (o *ArrayIterator) CanIterate() bool {
	return false
}

// Length returns the length of the Int object.
func (o *ArrayIterator) Length() int {
	return 0
}

// TypeName returns the type name of the ArrayIterator as a string.
func (i *ArrayIterator) TypeName() string {
	return ArrayIteratorType
}

// String returns a string representation of the ArrayIterator instance.
func (i *ArrayIterator) String() string {
	return ArrayIteratorLabel
}

// Boolean determines whether the ArrayIterator should be considered a falsy values. Always returns true.
func (i *ArrayIterator) Boolean() bool {
	return true
}

// Equals checks whether the given IObject is equal to the current ArrayIterator instance by values comparison.
func (i *ArrayIterator) Equals(IObject) bool {
	return false
}

// Copy creates and returns a duplicate of the ArrayIterator, preserving its current state.
func (i *ArrayIterator) Copy(frame int, _ int) IObject {
	ret := i.gk.NewArrayIterator(frame, i.values, i.index)
	return ret
}

// Next advances the iterator to the next element and returns true if the current position is within bounds.
func (i *ArrayIterator) Next() bool {
	i.index++
	return i.index <= i.length
}

// Key returns the index of the current element in the iteration as an IObject.
func (i *ArrayIterator) Key(frame int) IObject {
	idx := int64(i.index - 1)
	if idx < 0 || idx >= int64(i.length) {
		return i.gk.UndefinedValue()
	}
	return i.gk.NewInt(frame, idx)
}

// Value returns the current element in the iteration based on the iterator's internal position.
func (i *ArrayIterator) Value(frame int) IObject {
	idx := int64(i.index - 1)
	if idx < 0 || idx >= int64(i.length) {
		return i.gk.UndefinedValue()
	}
	return i.values[idx]
}
