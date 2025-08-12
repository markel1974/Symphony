package objects

// ObjectPtr is a wrapper around a pointer to an IObject, allowing additional behaviors and encapsulation of the values.
// It embeds ObjectImpl, inheriting default behaviors for the IObject interface methods.
// The values field holds the actual IObject instance being wrapped.
type ObjectPtr struct {
	ObjectImpl
	value *IObject
}

// NewObjectPtr creates a new ObjectPtr instance wrapping the provided IObject pointer.// NewObjectPtr creates a new ObjectPtr instance with the provided IObject values.
func NewObjectPtr(value *IObject) *ObjectPtr {
	return &ObjectPtr{value: value}
}

// Value returns the internal IObject pointer stored in the ObjectPtr instance.
func (o *ObjectPtr) Value() *IObject {
	return o.value
}

// SetValue sets the internal values field of the ObjectPtr to the provided IObject pointer.
func (o *ObjectPtr) SetValue(value IObject) {
	*o.value = value
}

// String returns the string representation of the ObjectPtr as "free-var".
func (o *ObjectPtr) String() string {
	return "free-var"
}

// TypeName returns the type name of the object as a string, in this case, "<free-var>".
func (o *ObjectPtr) TypeName() string {
	return "<free-var>"
}

// Copy creates and returns a duplicate of the object implementing the IObject interface.
func (o *ObjectPtr) Copy() IObject {
	return o
}

// Falsy returns true if the values of the ObjectPtr is nil, indicating it is considered falsy in a boolean context.
func (o *ObjectPtr) Falsy() bool {
	return o.value == nil
}

// Equals checks if the current ObjectPtr is equal to the provided IObject by comparing their memory addresses.
func (o *ObjectPtr) Equals(x IObject) bool {
	return o == x
}
