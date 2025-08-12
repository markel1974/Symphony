package objects

// ArrayIterator is an iterator type for traversing elements of an array.
// It implements the IIterator interface to provide sequential access to array elements.
type ArrayIterator struct {
	ObjectImpl
	v []IObject
	i int
	l int
}

// NewArrayIterator creates and returns a new ArrayIterator instance with the given slice of IObject.
func NewArrayIterator(v []IObject) *ArrayIterator {
	return &ArrayIterator{v: v, l: len(v), i: 0}
}

// TypeName returns the type name of the ArrayIterator as a string.
func (i *ArrayIterator) TypeName() string {
	return "array-iterator"
}

// String returns a string representation of the ArrayIterator instance.
func (i *ArrayIterator) String() string {
	return "<array-iterator>"
}

// Falsy determines whether the ArrayIterator should be considered a falsy values. Always returns true.
func (i *ArrayIterator) Boolean() bool {
	return true
}

// Equals checks whether the given IObject is equal to the current ArrayIterator instance by values comparison.
func (i *ArrayIterator) Equals(IObject) bool {
	return false
}

// Copy creates and returns a duplicate of the ArrayIterator, preserving its current state.
func (i *ArrayIterator) Copy() IObject {
	return &ArrayIterator{v: i.v, i: i.i, l: i.l}
}

// Next advances the iterator to the next element and returns true if the current position is within bounds.
func (i *ArrayIterator) Next() bool {
	i.i++
	return i.i <= i.l
}

// Key returns the index of the current element in the iteration as an IObject.
func (i *ArrayIterator) Key() IObject {
	return NewInt(int64(i.i - 1))
}

// Value returns the current element in the iteration based on the iterator's internal position.
func (i *ArrayIterator) Value() IObject {
	return i.v[i.i-1]
}
