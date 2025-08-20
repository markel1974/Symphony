package objects

// Object is a default implementation of the IObject interface with unimplemented or default behavior for methods.
type Object struct {
	factory *Factory
	frame   int
}

func newObject(factory *Factory, frame int) *Object {
	return &Object{
		factory: factory,
		frame:   frame,
	}
}

// Factory returns a reference to the Factory associated with the Object.
func (o *Object) Factory() *Factory {
	return o.factory
}

// Frame returns the current frame value of the Object.
func (o *Object) Frame() int {
	return o.frame
}

// SetFrame updates the frame field of the Object with the specified frame value.
func (o *Object) SetFrame(frame int) {
	o.frame = frame
}

// TypeName returns the name of the object type. This method must be implemented by objects inheriting Object.
func (o *Object) TypeName() string {
	return "not_implemented"
}

// String returns a string representation of the object. This method must be implemented by objects inheriting Object.
func (o *Object) String() string {
	return "not_implemented"
}

// BinaryOp performs a binary operation on the current object and another object using the specified operator.
// Returns the result of the operation or an error if the operation is not supported.
func (o *Object) BinaryOp(_ int, _ Operator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// Copy creates and returns a new instance of the object, duplicating its state.
func (o *Object) Copy(_ int) IObject {
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
func (o *Object) IndexGet(_ int, _ IObject) (res IObject, err error) {
	return nil, ErrNotIndexable
}

// IndexSet attempts to assign a value to an index in the object but always returns ErrNotIndexAssignable,
// as this operation is unsupported.
func (o *Object) IndexSet(_, _ IObject) (err error) {
	return ErrNotIndexAssignable
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (o *Object) Iterate(_ int) IIterator {
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

// Length returns the length of the Int object.
func (o *Object) Length() int {
	return 0
}
