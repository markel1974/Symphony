package objects

import "encoding/gob"

const (
	ObjectPointerType  = "object_pointer"
	ObjectPointerLabel = "<" + ObjectPointerType + ">"
)

func init() {
	gob.Register(&ObjectPointer{})
}

// ObjectPointer is a wrapper around a pointer to an IObject, allowing additional behaviors and encapsulation of the values.
// It embeds Object, inheriting default behaviors for the IObject interface methods.
// The value field holds the actual IObject instance being wrapped.
type ObjectPointer struct {
	factory IGateKeeper
	frame   int
	value   *IObject
}

// NewObjectPointer creates a new ObjectPointer instance wrapping the provided IObject pointer.// NewObjectPointer creates a new ObjectPointer instance with the provided IObject values.
func newObjectPointer(factory IGateKeeper, frame int, value *IObject) IObject {
	return &ObjectPointer{
		factory: factory,
		frame:   frame,
		value:   value,
	}
}

// GateKeeper returns a reference to the GateKeeper associated with the Object.
func (o *ObjectPointer) GateKeeper() IGateKeeper {
	return o.factory
}

// AsBool returns the boolean representation of the ObjectPointer, defaulting to false.
func (o *ObjectPointer) AsBool() bool {
	return (*o.value).AsBool()
}

// AsInt64 returns the length of the array as an int64 value.
func (o *ObjectPointer) AsInt64() int64 {
	return (*o.value).AsInt64()
}

// AsFloat64 returns the length of the array as an int64 value.
func (o *ObjectPointer) AsFloat64() float64 {
	return (*o.value).AsFloat64()
}

// AsString returns the string representation of the ObjectPointer instance.
func (o *ObjectPointer) AsString() string {
	return (*o.value).AsString()
}

// Nil checks if the object is nil and always returns false.
func (o *ObjectPointer) Nil() bool {
	return false
}

// AssignValue sets the current object to the provided IObject, returning ErrNotAssignable if the operation is not supported.
func (o *ObjectPointer) AssignValue(v IObject) error {
	return (*o.value).AssignValue(v)
}

// SetStatic sets the frame to FrameStatic, marking it with a static execution context.
func (o *ObjectPointer) SetStatic() {
	o.frame = FrameStatic
}

// Frame returns the current frame value of the Object.
func (o *ObjectPointer) Frame() int {
	return o.frame
}

// LogicalOp performs a logical operation with the given operator and RHS object, returning the result or an error.
func (o *ObjectPointer) LogicalOp(_ int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		switch op {
		case OperatorLogicalEq:
			if o.value == nil {
				return o.GateKeeper().TrueValue(), nil
			} else {
				return o.GateKeeper().FalseValue(), nil
			}
		case OperatorLogicalNotEq:
			if o.value == nil {
				return o.GateKeeper().FalseValue(), nil
			} else {
				return o.GateKeeper().TrueValue(), nil
			}
		default:
			return o.GateKeeper().UndefinedValue(), ErrInvalidOperator
		}
	}
	return o.GateKeeper().UndefinedValue(), ErrInvalidOperator
}

// ArithmeticOp performs an arithmetic operation with the given operator and right-hand-side operand and returns the result.
// Returns an error if the operation is invalid.
func (o *ObjectPointer) ArithmeticOp(frame int, obj ArithmeticOperator, rhs IObject) (IObject, error) {
	return (*o.value).ArithmeticOp(frame, obj, rhs)
}

// IndexGet attempts to retrieve a value at the given index and returns an error if the object is not indexable.
func (o *ObjectPointer) IndexGet(frame int, obj IObject) (res IObject, err error) {
	return (*o.value).IndexGet(frame, obj)
}

// IndexSet attempts to assign a value to an index in the object but always returns ErrUnsupportedIndex,
// as this operation is unsupported.
func (o *ObjectPointer) IndexSet(frame, obj IObject) (err error) {
	return (*o.value).IndexSet(frame, obj)
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (o *ObjectPointer) Iterate(frame int) IIterator {
	return (*o.value).Iterate(frame)
}

// CanIterate determines if the object can be iterated over and returns false for this implementation.
func (o *ObjectPointer) CanIterate() bool {
	return (*o.value).CanIterate()
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *ObjectPointer) Call(frame int, v ...IObject) (ret IObject, err error) {
	return (*o.value).Call(frame, v...)
}

// CanCall determines if the object can be invoked as a callable. Returns false for non-callable objects.
func (o *ObjectPointer) CanCall() bool {
	return (*o.value).CanCall()
}

// Length returns the length of the Int object.
func (o *ObjectPointer) Length() int {
	return (*o.value).Length()
}

// Value returns the internal IObject pointer stored in the ObjectPointer instance.
func (o *ObjectPointer) Value() *IObject {
	return o.value
}

// SetValue sets the internal values field of the ObjectPointer to the provided IObject pointer.
func (o *ObjectPointer) SetValue(value IObject) {
	*o.value = value
}

// TypeName returns the type name of the ObjectPointer as a string.
func (o *ObjectPointer) TypeName() string {
	return (*o.value).TypeName()
}

// Copy creates and returns a duplicate of the object implementing the IObject interface.
func (o *ObjectPointer) Copy(_ int, _ int) IObject {
	return o
}

// Falsy returns true if the value of the ObjectPointer is nil.
func (o *ObjectPointer) Falsy() bool {
	return o.value == nil
}

// Equals checks if the current ObjectPointer is equal to the provided IObject by comparing their memory addresses.
func (o *ObjectPointer) Equals(x IObject) bool {
	return o == x
}
