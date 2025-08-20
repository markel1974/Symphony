package objects

const (
	BytesIteratorType  = "bytes_iterator"
	BytesIteratorLabel = "<" + BytesIteratorType + ">"
)

// BytesIterator is an iterator for traversing elements of a byte slice, implementing the IIterator interface.
type BytesIterator struct {
	*Object
	values []byte
	index  int
	length int
}

func _newBytesIterator(factory *Factory, v []byte) *BytesIterator {
	return &BytesIterator{
		Object: factory.NewObject(),
		values: v,
		length: len(v),
		index:  0,
	}
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
func (i *BytesIterator) Copy() IObject {
	ret := i.Factory().NewBytesIterator(i.values)
	ret.index = i.index
	return ret
}

// Next advances the iterator to the next position and returns true if the current position is within bounds.
func (i *BytesIterator) Next() bool {
	i.index++
	return i.index <= i.length
}

// Key returns the current index of the iterator as an IObject, decremented by one from the internal index tracker.
func (i *BytesIterator) Key() IObject {
	return i.Factory().NewInt(int64(i.index - 1))
}

// Value returns the values of the current byte in the iteration as an IObject, wrapped in an Int struct.
func (i *BytesIterator) Value() IObject {
	return i.Factory().NewInt(int64(i.values[i.index-1]))
}
