package objects

// UndefinedIteratorType represents a type label for an undefined iterator.
// UndefinedIteratorLabel represents a label used when the iterator type is undefined.
const (
	UndefinedIteratorType  = "undefined_iterator"
	UndefinedIteratorLabel = "<" + UndefinedIteratorType + ">"
)

// UndefinedIterator represents an iterator that is implicitly undefined and cannot traverse any elements.
type UndefinedIterator struct {
	IObject
}

// newUndefinedIterator creates a new instance of UndefinedIterator initialized with an undefined object from the GateKeeper.
func newUndefinedIterator(factory *GateKeeper, frame int) *UndefinedIterator {
	return &UndefinedIterator{
		IObject: factory.newObject(frame),
	}
}

// Copy creates and returns a reference to the same UndefinedIterator, ignoring the input parameters.
func (i *UndefinedIterator) Copy(_ int, _ int) IObject {
	return i
}

// TypeName returns the type name of the UndefinedIterator object as a string.
func (i *UndefinedIterator) TypeName() string {
	return UndefinedIteratorType
}

// String returns a string representation of the UndefinedIterator object.
func (i *UndefinedIterator) String() string {
	return UndefinedIteratorLabel
}

// Boolean returns the boolean representation of the UndefinedIterator, which is always true.
func (i *UndefinedIterator) Boolean() bool {
	return true
}

// Equals checks if the current UndefinedIterator is equal to the provided IObject and always returns false.
func (i *UndefinedIterator) Equals(IObject) bool {
	return false
}

// Next advances the iterator and returns false, as the UndefinedIterator does not support iteration.
func (i *UndefinedIterator) Next() bool {
	return false
}

// Key returns the key of the current element in the iteration, which is always the undefined value for this iterator type.
func (i *UndefinedIterator) Key(frame int) IObject {
	return i.GateKeeper().UndefinedValue()
}

// Value retrieves the undefined value associated with the current UndefinedIterator.
func (i *UndefinedIterator) Value(_ int) IObject {
	return i.GateKeeper().UndefinedValue()
}
