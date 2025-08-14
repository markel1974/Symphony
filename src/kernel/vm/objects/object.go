package objects

// Object is a default implementation of the IObject interface with unimplemented or default behavior for methods.
type Object struct {
}

// TypeName returns the name of the type as a string. This method must be implemented by objects inheriting Object.
func (o *Object) TypeName() string {
	panic(ErrNotImplemented)
}

// String returns the string representation of the Object. Currently, it is not implemented and will panic.
func (o *Object) String() string {
	panic(ErrNotImplemented)
}

// BinaryOp performs a binary operation on the current object and another object using the specified operator.
// Returns the result of the operation or an error if the operation is not supported.
func (o *Object) BinaryOp(_ Operator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// Copy creates and returns a new instance of the object, duplicating its state.
func (o *Object) Copy() IObject {
	return nil
}

// Boolean returns false for all objects.
func (o *Object) Boolean() bool {
	return false
}

// Equals checks whether the current object is equal to another object of type IObject.
func (o *Object) Equals(x IObject) bool {
	return o == x
}

// IndexGet attempts to retrieve a value at the given index and returns an error if the object is not indexable.
func (o *Object) IndexGet(_ IObject) (res IObject, err error) {
	return nil, ErrNotIndexable
}

// IndexSet attempts to assign a value to an index in the object but always returns ErrNotIndexAssignable,
// as this operation is unsupported.
func (o *Object) IndexSet(_, _ IObject) (err error) {
	return ErrNotIndexAssignable
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (o *Object) Iterate() IIterator {
	return nil
}

// CanIterate determines if the object can be iterated over and returns false for this implementation.
func (o *Object) CanIterate() bool {
	return false
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *Object) Call(_ ...IObject) (ret IObject, err error) {
	return nil, nil
}

// CanCall determines if the object can be invoked as a callable. Returns false for non-callable objects.
func (o *Object) CanCall() bool {
	return false
}
