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

// NewBytesIterator creates and returns a new BytesIterator instance with the given slice of bytes.
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

// AsBool returns true if the object is not empty, otherwise false.
func (o *BytesIterator) AsBool() bool {
	return o.length > 0
}

// AsInt64 returns the length of the array as an int64 value.
func (o *BytesIterator) AsInt64() int64 {
	return int64(o.length)
}

// AsFloat64 returns the length of the array as an int64 value.
func (o *BytesIterator) AsFloat64() float64 {
	return float64(o.length)
}

// AsString returns the string representation of the BytesIterator.
func (o *BytesIterator) AsString() string {
	return BytesIteratorLabel
}

// Frame returns the current frame value of the Object.
func (o *BytesIterator) Frame() int {
	return o.frame
}

// LogicalOp performs a logical operation on the BytesIterator object, always returning an error for invalid operations.
func (o *BytesIterator) LogicalOp(_ int, _ LogicalOperator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// ArithmeticOp performs an arithmetic operation with the given operator and operands; returns an error for invalid operators.
func (o *BytesIterator) ArithmeticOp(_ int, _ ArithmeticOperator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// Falsy returns false for all objects.
func (o *BytesIterator) Falsy() bool {
	return false
}

// IndexGet attempts to retrieve a value at the given index and returns an error if the object is not indexable.
func (o *BytesIterator) IndexGet(_ int, _ IObject) (res IObject, err error) {
	return nil, ErrNotIndexable
}

// IndexSet attempts to assign a value to an index in the object but always returns ErrUnsupportedIndex,
// as this operation is unsupported.
func (o *BytesIterator) IndexSet(_, _ IObject) (err error) {
	return ErrUnsupportedIndex
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
func (o *BytesIterator) TypeName() string {
	return BytesIteratorType
}

// Equals checks whether the BytesIterator is equal to another object implementing the IObject interface.
func (o *BytesIterator) Equals(IObject) bool {
	return false
}

// Copy creates and returns a new instance of BytesIterator with the same state as the current instance.
func (o *BytesIterator) Copy(frame int, _ int) IObject {
	ret := o.GateKeeper().NewBytesIterator(frame, o.values, o.index)
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

// Value returns the values of the current byte in the iteration as an IObject, wrapped in an Int struct.
func (o *BytesIterator) Value(frame int) IObject {
	idx := o.index - 1
	if idx < 0 || idx >= o.length {
		return o.GateKeeper().UndefinedValue()
	}
	return o.GateKeeper().NewInt(frame, int64(o.values[idx]))
}
