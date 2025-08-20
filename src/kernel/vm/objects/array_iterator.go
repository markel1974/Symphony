package objects

const (
	ArrayIteratorType  = "array_iterator"
	ArrayIteratorLabel = "<" + ArrayIteratorType + ">"
)

// ArrayIterator is an iterator type for traversing elements of an array.
// It implements the IIterator interface to provide sequential access to array elements.
type ArrayIterator struct {
	*Object
	values []IObject
	index  int
	length int
}

// NewArrayIterator creates and returns a new ArrayIterator instance with the given slice of IObject.
func _newArrayIterator(factory *Factory, v []IObject) *ArrayIterator {
	return &ArrayIterator{
		Object: factory.NewObject(),
		values: v,
		length: len(v),
		index:  0,
	}
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
func (i *ArrayIterator) Copy() IObject {
	ret := i.Factory().NewArrayIterator(i.values)
	ret.index = i.index
	return ret
}

// Next advances the iterator to the next element and returns true if the current position is within bounds.
func (i *ArrayIterator) Next() bool {
	i.index++
	return i.index <= i.length
}

// Key returns the index of the current element in the iteration as an IObject.
func (i *ArrayIterator) Key() IObject {
	return i.Factory().NewInt(int64(i.index - 1))
}

// Value returns the current element in the iteration based on the iterator's internal position.
func (i *ArrayIterator) Value() IObject {
	return i.values[i.index-1]
}
