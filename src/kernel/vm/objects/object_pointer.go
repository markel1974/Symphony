package objects

const (
	ObjectPointerType  = "object_pointer"
	ObjectPointerLabel = "<" + ObjectPointerType + ">"
)

// ObjectPointer is a wrapper around a pointer to an IObject, allowing additional behaviors and encapsulation of the values.
// It embeds Object, inheriting default behaviors for the IObject interface methods.
// The value field holds the actual IObject instance being wrapped.
type ObjectPointer struct {
	IObject
	value *IObject
}

// NewObjectPointer creates a new ObjectPointer instance wrapping the provided IObject pointer.// NewObjectPointer creates a new ObjectPointer instance with the provided IObject values.
func newObjectPointer(factory *GateKeeper, frame int, value *IObject) *ObjectPointer {
	return &ObjectPointer{
		IObject: factory.newObject(frame),
		value:   value,
	}
}

// Value returns the internal IObject pointer stored in the ObjectPointer instance.
func (o *ObjectPointer) Value() *IObject {
	return o.value
}

// SetValue sets the internal values field of the ObjectPointer to the provided IObject pointer.
func (o *ObjectPointer) SetValue(value IObject) {
	*o.value = value
}

// String returns the string representation of the ObjectPointer instance.
func (o *ObjectPointer) String() string {
	return ObjectPointerType
}

// TypeName returns the type name of the ObjectPointer as a string.
func (o *ObjectPointer) TypeName() string {
	return ObjectPointerLabel
}

// Copy creates and returns a duplicate of the object implementing the IObject interface.
func (o *ObjectPointer) Copy(_ int, _ int) IObject {
	return o
}

// Boolean returns true if the value of the ObjectPointer is nil.
func (o *ObjectPointer) Boolean() bool {
	return o.value == nil
}

// Equals checks if the current ObjectPointer is equal to the provided IObject by comparing their memory addresses.
func (o *ObjectPointer) Equals(x IObject) bool {
	return o == x
}
