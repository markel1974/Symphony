package objects

// ObjectPointer is a wrapper around a pointer to an IObject, allowing additional behaviors and encapsulation of the values.
// It embeds ObjectImpl, inheriting default behaviors for the IObject interface methods.
// The values field holds the actual IObject instance being wrapped.
type ObjectPointer struct {
	ObjectImpl
	value *IObject
}

// NewObjectPointer creates a new ObjectPointer instance wrapping the provided IObject pointer.// NewObjectPointer creates a new ObjectPointer instance with the provided IObject values.
func NewObjectPointer(value *IObject) *ObjectPointer {
	return &ObjectPointer{value: value}
}

// Value returns the internal IObject pointer stored in the ObjectPointer instance.
func (o *ObjectPointer) Value() *IObject {
	return o.value
}

// SetValue sets the internal values field of the ObjectPointer to the provided IObject pointer.
func (o *ObjectPointer) SetValue(value IObject) {
	*o.value = value
}

// String returns the string representation of the ObjectPointer as "free-var".
func (o *ObjectPointer) String() string {
	return "free-var"
}

// TypeName returns the type name of the object as a string, in this case, "<free-var>".
func (o *ObjectPointer) TypeName() string {
	return "<free-var>"
}

// Copy creates and returns a duplicate of the object implementing the IObject interface.
func (o *ObjectPointer) Copy() IObject {
	return o
}

// Falsy returns true if the values of the ObjectPointer is nil, indicating it is considered falsy in a boolean context.
func (o *ObjectPointer) Boolean() bool {
	return o.value == nil
}

// Equals checks if the current ObjectPointer is equal to the provided IObject by comparing their memory addresses.
func (o *ObjectPointer) Equals(x IObject) bool {
	return o == x
}
