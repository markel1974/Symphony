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
	Allocator
	values []rune
	index  int
	length int
}

// NewStringIterator creates and returns a new instance of StringIterator with the given rune slice.
func newStringIterator(factory IGateKeeper, frame int, v []rune, index int) IIterator {
	return &StringIterator{
		Allocator: Allocator{gk: factory, frame: frame},
		values:    v,
		length:    len(v),
		index:     index,
	}
}

// AsBool returns true if the array is not empty, otherwise false.
func (o *StringIterator) AsBool() bool {
	return o.length > 0
}

// AsInt64 returns the length of the array as an int64 value.
func (o *StringIterator) AsInt64() int64 {
	return int64(o.length)
}

// AsFloat64 returns the length of the array as an int64 value.
func (o *StringIterator) AsFloat64() float64 {
	return float64(o.length)
}

// AsString returns the string representation of the StringIterator, useful for debugging or logging.
func (o *StringIterator) AsString() string {
	return StringIteratorLabel
}

// AssignValue sets the current object to the provided IObject, returning ErrNotAssignable if the operation is not supported.
func (o *StringIterator) AssignValue(_ IObject) error {
	return ErrNotAssignable
}

// Nil checks if the object is nil and always returns false.
func (o *StringIterator) Nil() bool {
	return false
}

// LogicalOp performs a logical operation using a specified operator and right-hand operand, but always returns an error.
func (o *StringIterator) LogicalOp(_ int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		return logicalOpNil(o.gk, op)
	}
	return nil, ErrInvalidOperator
}

// ArithmeticOp performs an arithmetic operation on the StringIterator but always returns ErrInvalidOperator and no result.
func (o *StringIterator) ArithmeticOp(_ int, _ ArithmeticOperator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// IndexGet attempts to retrieve a value at the given index and returns an error if the object is not indexable.
func (o *StringIterator) IndexGet(_ int, _ IObject) (IObject, error) {
	return o.gk.UndefinedValue(), ErrNotIndexable
}

// IndexSet attempts to assign a value to an index in the object but always returns ErrUnsupportedIndex,
// as this operation is unsupported.
func (o *StringIterator) IndexSet(_, _ IObject) error {
	return ErrUnsupportedIndex
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (o *StringIterator) Iterate(_ int) IIterator {
	return nil
}

// Iterable determines if the object can be iterated over and returns false for this implementation.
func (o *StringIterator) Iterable() bool {
	return false
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *StringIterator) Call(_ int, _ ...IObject) (retCount uint, ret IObject, err error) {
	return 0, nil, nil
}

// Length returns the length of the Int object.
func (o *StringIterator) Length() int {
	return 0
}

// Copy creates and returns a new instance of StringIterator with the same state as the current one.
func (o *StringIterator) Copy(frame int, _ int) IObject {
	ret := o.GateKeeper().NewStringIterator(frame, o.values, o.index)
	return ret
}

// TypeName returns the type name of the StringIterator as a string.
func (o *StringIterator) TypeName() string {
	return StringIteratorType
}

// Falsy returns true, indicating the iterator is considered falsy in a boolean context.
func (o *StringIterator) Falsy() bool {
	return true
}

// Equals compare the current StringIterator with another IObject and determine if they are equal.
func (o *StringIterator) Equals(IObject) bool {
	return false
}

// Next advances the iterator to the next position and returns true if the current position is within bounds.
func (o *StringIterator) Next() bool {
	o.index++
	return o.index <= o.length
}

// Key returns the zero-based index of the current element in the iteration as an Int object.
func (o *StringIterator) Key(frame int) IObject {
	idx := o.index - 1
	if idx < 0 || idx >= o.length {
		return o.GateKeeper().UndefinedValue()
	}
	return o.GateKeeper().NewInt(frame, int64(idx))
}

// Value returns the current character as an IObject wrapped in a Char instance from the internal rune slice.
func (o *StringIterator) Value(frame int) IObject {
	idx := o.index - 1
	if idx < 0 || idx >= o.length {
		return o.GateKeeper().UndefinedValue()
	}
	return o.GateKeeper().NewChar(frame, o.values[idx])
}

// Count returns the total number of elements in the instance and its sub-elements.
func (o *StringIterator) Count() int {
	return 1
}
