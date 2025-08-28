package objects

import "encoding/gob"

const (
	StringIteratorType  = "string_iterator"
	StringIteratorLabel = "<" + StringIteratorType + ">"
)

func init() {
	gob.Register(&StringIterator{})
}

// StringIterator represents an iterator for traversing over the characters of a string, implemented as runes.
type StringIterator struct {
	factory IGateKeeper
	frame   int
	values  []rune
	index   int
	length  int
}

// NewStringIterator creates and returns a new instance of StringIterator with the given rune slice.
func newStringIterator(factory IGateKeeper, frame int, v []rune, index int) IIterator {
	return &StringIterator{
		factory: factory,
		frame:   frame,
		values:  v,
		length:  len(v),
		index:   index,
	}
}

// GateKeeper returns a reference to the GateKeeper associated with the Object.
func (o *StringIterator) GateKeeper() IGateKeeper {
	return o.factory
}

// Frame returns the current frame value of the Object.
func (o *StringIterator) Frame() int {
	return o.frame
}

// BinaryOp performs a binary operation on the current object and another object using the specified operator.
// Returns the result of the operation or an error if the operation is not supported.
func (o *StringIterator) BinaryOp(_ int, _ Operator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// IndexGet attempts to retrieve a value at the given index and returns an error if the object is not indexable.
func (o *StringIterator) IndexGet(_ int, _ IObject) (res IObject, err error) {
	return nil, ErrNotIndexable
}

// IndexSet attempts to assign a value to an index in the object but always returns ErrNotIndexAssignable,
// as this operation is unsupported.
func (o *StringIterator) IndexSet(_, _ IObject) (err error) {
	return ErrNotIndexAssignable
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (o *StringIterator) Iterate(_ int) IIterator {
	return nil
}

// CanIterate determines if the object can be iterated over and returns false for this implementation.
func (o *StringIterator) CanIterate() bool {
	return false
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *StringIterator) Call(_ int, _ ...IObject) (ret IObject, err error) {
	return nil, nil
}

// CanCall determines if the object can be invoked as a callable. Returns false for non-callable objects.
func (o *StringIterator) CanCall() bool {
	return false
}

// Length returns the length of the Int object.
func (o *StringIterator) Length() int {
	return 0
}

// Copy creates and returns a new instance of StringIterator with the same state as the current one.
func (i *StringIterator) Copy(frame int, _ int) IObject {
	ret := i.GateKeeper().NewStringIterator(frame, i.values, i.index)
	return ret
}

// TypeName returns the type name of the StringIterator as a string.
func (i *StringIterator) TypeName() string {
	return StringIteratorType
}

// String returns the string representation of the StringIterator, useful for debugging or logging.
func (i *StringIterator) String() string {
	return StringIteratorLabel
}

// Falsy returns true, indicating the iterator is considered falsy in a boolean context.
func (i *StringIterator) Falsy() bool {
	return true
}

// Equals compare the current StringIterator with another IObject and determine if they are equal.
func (i *StringIterator) Equals(IObject) bool {
	return false
}

// Next advances the iterator to the next position and returns true if the current position is within bounds.
func (i *StringIterator) Next() bool {
	i.index++
	return i.index <= i.length
}

// Key returns the zero-based index of the current element in the iteration as an Int object.
func (i *StringIterator) Key(frame int) IObject {
	idx := i.index - 1
	if idx < 0 || idx >= i.length {
		return i.GateKeeper().UndefinedValue()
	}
	return i.GateKeeper().NewInt(frame, int64(idx))
}

// Value returns the current character as an IObject wrapped in a Char instance from the internal rune slice.
func (i *StringIterator) Value(frame int) IObject {
	idx := i.index - 1
	if idx < 0 || idx >= i.length {
		return i.GateKeeper().UndefinedValue()
	}
	return i.GateKeeper().NewChar(frame, i.values[idx])
}
