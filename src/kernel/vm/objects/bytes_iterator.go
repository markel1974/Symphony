package objects

import "encoding/gob"

const (
	BytesIteratorType  = "bytes_iterator"
	BytesIteratorLabel = "<" + BytesIteratorType + ">"
)

func init() {
	gob.Register(&BytesIterator{})
}

// BytesIterator is an iterator for traversing elements of a byte slice, implementing the IIterator interface.
type BytesIterator struct {
	gk     IGateKeeper
	frame  int
	values []byte
	index  int
	length int
}

func newBytesIterator(factory IGateKeeper, frame int, v []byte, index int) IIterator {
	return &BytesIterator{
		gk:     factory,
		frame:  frame,
		values: v,
		length: len(v),
		index:  index,
	}
}

// GateKeeper returns a reference to the GateKeeper associated with the Object.
func (o *BytesIterator) GateKeeper() IGateKeeper {
	return o.gk
}

// Frame returns the current frame value of the Object.
func (o *BytesIterator) Frame() int {
	return o.frame
}

// BinaryOp performs a binary operation on the current object and another object using the specified operator.
// Returns the result of the operation or an error if the operation is not supported.
func (o *BytesIterator) BinaryOp(_ int, _ Operator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// Boolean returns false for all objects.
func (o *BytesIterator) Boolean() bool {
	return false
}

// IndexGet attempts to retrieve a value at the given index and returns an error if the object is not indexable.
func (o *BytesIterator) IndexGet(_ int, _ IObject) (res IObject, err error) {
	return nil, ErrNotIndexable
}

// IndexSet attempts to assign a value to an index in the object but always returns ErrNotIndexAssignable,
// as this operation is unsupported.
func (o *BytesIterator) IndexSet(_, _ IObject) (err error) {
	return ErrNotIndexAssignable
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (o *BytesIterator) Iterate(_ int) IIterator {
	return nil
}

// CanIterate determines if the object can be iterated over and returns false for this implementation.
func (o *BytesIterator) CanIterate() bool {
	return false
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *BytesIterator) Call(_ int, _ ...IObject) (ret IObject, err error) {
	return nil, nil
}

// CanCall determines if the object can be invoked as a callable. Returns false for non-callable objects.
func (o *BytesIterator) CanCall() bool {
	return false
}

// Length returns the length of the Int object.
func (o *BytesIterator) Length() int {
	return 0
}

// TypeName returns the string representation of the type name, which is "bytes-iterator".
func (i *BytesIterator) TypeName() string {
	return BytesIteratorType
}

// String returns the string representation of the BytesIterator.
func (i *BytesIterator) String() string {
	return BytesIteratorLabel
}

// Equals checks whether the BytesIterator is equal to another object implementing the IObject interface.
func (i *BytesIterator) Equals(IObject) bool {
	return false
}

// Copy creates and returns a new instance of BytesIterator with the same state as the current instance.
func (i *BytesIterator) Copy(frame int, _ int) IObject {
	ret := i.GateKeeper().NewBytesIterator(frame, i.values, i.index)
	return ret
}

// Next advances the iterator to the next position and returns true if the current position is within bounds.
func (i *BytesIterator) Next() bool {
	i.index++
	return i.index <= i.length
}

// Key returns the current index of the iterator as an IObject, decremented by one from the internal index tracker.
func (i *BytesIterator) Key(frame int) IObject {
	idx := i.index - 1
	if idx < 0 || idx >= i.length {
		return i.GateKeeper().UndefinedValue()
	}
	return i.GateKeeper().NewInt(frame, int64(idx))
}

// Value returns the values of the current byte in the iteration as an IObject, wrapped in an Int struct.
func (i *BytesIterator) Value(frame int) IObject {
	idx := i.index - 1
	if idx < 0 || idx >= i.length {
		return i.GateKeeper().UndefinedValue()
	}
	return i.GateKeeper().NewInt(frame, int64(i.values[idx]))
}
