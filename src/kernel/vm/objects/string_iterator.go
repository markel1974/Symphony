package objects

const (
	StringIteratorType  = "string_iterator"
	StringIteratorLabel = "<" + StringIteratorType + ">"
)

// StringIterator represents an iterator for traversing over the characters of a string, implemented as runes.
type StringIterator struct {
	*Object
	values []rune
	index  int
	length int
}

// NewStringIterator creates and returns a new instance of StringIterator with the given rune slice.
func newStringIterator(factory *GateKeeper, frame int, v []rune) *StringIterator {
	return &StringIterator{
		Object: factory.NewObject(frame),
		values: v,
		length: len(v),
		index:  0,
	}
}

// Copy creates and returns a new instance of StringIterator with the same state as the current one.
func (i *StringIterator) Copy(frame int) IObject {
	ret := i.Factory().NewStringIterator(frame, i.values)
	ret.index = i.index
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

// Boolean returns true, indicating the iterator is considered falsy in a boolean context.
func (i *StringIterator) Boolean() bool {
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
	return i.Factory().NewInt(frame, int64(i.index-1))
}

// Value returns the current character as an IObject wrapped in a Char instance from the internal rune slice.
func (i *StringIterator) Value(frame int) IObject {
	return i.Factory().NewChar(frame, i.values[i.index-1])
}
