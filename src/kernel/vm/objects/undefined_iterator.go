package objects

// UndefinedIteratorType represents a type label for an undefined iterator.
// UndefinedIteratorLabel represents a label used when the iterator type is undefined.
const (
	UndefinedIteratorType  = "undefined_iterator"
	UndefinedIteratorLabel = "<" + UndefinedIteratorType + ">"
)

// UndefinedIterator represents an iterator that is implicitly undefined and cannot traverse any elements.
type UndefinedIterator struct {
	factory IGateKeeper
	frame   int
}

// newUndefinedIterator creates a new instance of UndefinedIterator initialized with an undefined object from the GateKeeper.
func newUndefinedIterator(factory IGateKeeper, frame int) IIterator {
	return &UndefinedIterator{
		factory: factory,
		frame:   frame,
	}
}

// GateKeeper returns a reference to the GateKeeper associated with the Object.
func (o *UndefinedIterator) GateKeeper() IGateKeeper {
	return o.factory
}

// Frame returns the current frame value of the Object.
func (o *UndefinedIterator) Frame() int {
	return o.frame
}

// BinaryOp performs a binary operation on the current object and another object using the specified operator.
// Returns the result of the operation or an error if the operation is not supported.
func (o *UndefinedIterator) BinaryOp(_ int, _ Operator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// IndexGet attempts to retrieve a value at the given index and returns an error if the object is not indexable.
func (o *UndefinedIterator) IndexGet(_ int, _ IObject) (res IObject, err error) {
	return nil, ErrNotIndexable
}

// IndexSet attempts to assign a value to an index in the object but always returns ErrNotIndexAssignable,
// as this operation is unsupported.
func (o *UndefinedIterator) IndexSet(_, _ IObject) (err error) {
	return ErrNotIndexAssignable
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (o *UndefinedIterator) Iterate(_ int) IIterator {
	return nil
}

// CanIterate determines if the object can be iterated over and returns false for this implementation.
func (o *UndefinedIterator) CanIterate() bool {
	return false
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *UndefinedIterator) Call(_ int, _ ...IObject) (ret IObject, err error) {
	return nil, nil
}

// CanCall determines if the object can be invoked as a callable. Returns false for non-callable objects.
func (o *UndefinedIterator) CanCall() bool {
	return false
}

// Length returns the length of the Int object.
func (o *UndefinedIterator) Length() int {
	return 0
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
