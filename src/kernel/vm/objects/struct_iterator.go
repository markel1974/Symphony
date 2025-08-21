package objects

// StructIteratorType represents the type for a struct iterator, used as a constant string identifier.
// StructIteratorLabel is a formatted label that includes the StructIteratorType constant within angle brackets.
const (
	StructIteratorType  = "struct_iterator"
	StructIteratorLabel = "<" + StructIteratorType + ">"
)

// StructIterator is an iterator for traversing over the keys and values of a struct-like map of IObjects.
// It embeds Object and implements the IIterator interface.
// The keys and values are stored in separate slices to facilitate order-preserving iteration.
// The index tracks the current position within the iteration, and length is the total number of elements.
type StructIterator struct {
	IObject
	values map[string]IObject
	keys   []string
	index  int
	length int
}

// NewStructIterator initializes and returns a new StructIterator for the given map of string keys to IObject values.
func newStructIterator(factory *GateKeeper, frame int, v map[string]IObject, index int) *StructIterator {
	var keys []string
	for k := range v {
		keys = append(keys, k)
	}
	return &StructIterator{
		IObject: factory.newObject(frame),
		values:  v,
		keys:    keys,
		length:  len(keys),
		index:   index,
	}
}

// Copy creates and returns a duplicate instance of the current StructIterator with the same internal state.
func (i *StructIterator) Copy(frame int, _ int) IObject {
	ret := i.GateKeeper().newStructIterator(frame, i.values, i.index)
	return ret
}

// TypeName returns the type name of the StructIterator object as a string.
func (i *StructIterator) TypeName() string {
	return StructIteratorType
}

// String returns the string representation of the StructIterator instance.
func (i *StructIterator) String() string {
	return StructIteratorLabel
}

// Boolean returns true, indicating the current StructIterator is truthy.
func (i *StructIterator) Boolean() bool {
	return true
}

// Equals compares the current StructIterator with another IObject for equality and returns true if they are equal.
func (i *StructIterator) Equals(IObject) bool {
	return false
}

// Next advances the iterator to the next position and returns true if there are more elements to iterate over.
func (i *StructIterator) Next() bool {
	i.index++
	return i.index <= i.length
}

// Key retrieves the key of the current element in the MapIterator as an IObject.
func (i *StructIterator) Key(frame int) IObject {
	idx := i.index - 1
	if idx < 0 || idx >= i.length {
		return i.GateKeeper().UndefinedValue()
	}
	k := i.keys[idx]
	return i.GateKeeper().NewString(frame, k)
}

// Value retrieves the value of the current element in the iteration based on the iterator's current position.
func (i *StructIterator) Value(_ int) IObject {
	idx := i.index - 1
	if idx < 0 || idx >= i.length {
		return i.GateKeeper().UndefinedValue()
	}
	k := i.keys[idx]
	return i.values[k]
}
