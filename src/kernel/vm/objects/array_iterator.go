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

// AsBool returns true if the array is not empty, otherwise false.
func (o *ArrayIterator) AsBool() bool {
	return o.length > 0
}

// AsInt64 returns the length of the array as an int64 value.
func (o *ArrayIterator) AsInt64() int64 {
	return int64(o.length)
}

// AsFloat64 returns the length of the array as an int64 value.
func (o *ArrayIterator) AsFloat64() float64 {
	return float64(o.length)
}

// AsString returns a string representation of the ArrayIterator instance.
func (o *ArrayIterator) AsString() string {
	return ArrayIteratorLabel
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

// LogicalOp attempts to perform a logical operation, returning an error as this operation is not supported.
func (o *ArrayIterator) LogicalOp(_ int, _ LogicalOperator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// ArithmeticOp performs an arithmetic operation on the ArrayIterator using a specified operator and operand, returning an error.
func (o *ArrayIterator) ArithmeticOp(_ int, _ ArithmeticOperator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// IndexGet retrieves an element at a specific index from the array but always returns ErrNotIndexable for this implementation.
func (o *ArrayIterator) IndexGet(_ int, _ IObject) (res IObject, err error) {
	return nil, ErrNotIndexable
}

// IndexSet attempts to assign a value to an index in the object but always returns ErrUnsupportedIndex,
// as this operation is unsupported.
func (o *ArrayIterator) IndexSet(_, _ IObject) (err error) {
	return ErrUnsupportedIndex
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
func (o *ArrayIterator) TypeName() string {
	return ArrayIteratorType
}

// Falsy determines whether the ArrayIterator should be considered a falsy value. Always returns true.
func (o *ArrayIterator) Falsy() bool {
	return true
}

// Equals checks whether the given IObject is equal to the current ArrayIterator instance by values comparison.
func (o *ArrayIterator) Equals(IObject) bool {
	return false
}

// Copy creates and returns a duplicate of the ArrayIterator, preserving its current state.
func (o *ArrayIterator) Copy(frame int, _ int) IObject {
	ret := o.gk.NewArrayIterator(frame, o.values, o.index)
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
		return o.gk.UndefinedValue()
	}
	return o.gk.NewInt(frame, idx)
}

// Value returns the current element in the iteration based on the iterator's internal position.
func (o *ArrayIterator) Value(_ int) IObject {
	idx := int64(o.index - 1)
	if idx < 0 || idx >= int64(o.length) {
		return o.gk.UndefinedValue()
	}
	return o.values[idx]
}
