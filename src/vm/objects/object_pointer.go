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
	Allocator
	valuePtr *IObject
}

// NewObjectPointer creates a new ObjectPointer instance wrapping the provided IObject pointer.// NewObjectPointer creates a new ObjectPointer instance with the provided IObject values.
func newObjectPointer(factory IGateKeeper, frame int, value *IObject) IObject {
	ptr := &ObjectPointer{
		Allocator: Allocator{gk: factory, frame: frame},
	}
	if value != nil {
		ptr.acquire(value)
	} else {
		undefined := factory.UndefinedValue()
		ptr.valuePtr = &undefined
	}
	return ptr
}

// AsBool returns the boolean representation of the ObjectPointer, defaulting to false.
func (o *ObjectPointer) AsBool() bool {
	return (*o.valuePtr).AsBool()
}

// AsInt64 returns the length of the array as an int64 value.
func (o *ObjectPointer) AsInt64() int64 {
	return (*o.valuePtr).AsInt64()
}

// AsFloat64 returns the length of the array as an int64 value.
func (o *ObjectPointer) AsFloat64() float64 {
	return (*o.valuePtr).AsFloat64()
}

// AsString returns the string representation of the ObjectPointer instance.
func (o *ObjectPointer) AsString() string {
	return (*o.valuePtr).AsString()
}

// Nil checks if the object is nil and always returns false.
func (o *ObjectPointer) Nil() bool {
	return false
}

// AssignValue sets the current object to the provided IObject, returning ErrNotAssignable if the operation is not supported.
func (o *ObjectPointer) AssignValue(v IObject) error {
	return (*o.valuePtr).AssignValue(v)
}

// LogicalOp performs a logical operation with the given operator and RHS object, returning the result or an error.
func (o *ObjectPointer) LogicalOp(_ int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		switch op {
		case OperatorLogicalEq:
			if o.valuePtr == nil {
				return o.GateKeeper().TrueValue(), nil
			} else {
				return o.GateKeeper().FalseValue(), nil
			}
		case OperatorLogicalNotEq:
			if o.valuePtr == nil {
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
	return (*o.valuePtr).ArithmeticOp(frame, obj, rhs)
}

// IndexGet attempts to retrieve a value at the given index and returns an error if the object is not indexable.
func (o *ObjectPointer) IndexGet(frame int, obj IObject) (res IObject, err error) {
	return (*o.valuePtr).IndexGet(frame, obj)
}

// IndexSet attempts to assign a value to an index in the object but always returns ErrUnsupportedIndex,
// as this operation is unsupported.
func (o *ObjectPointer) IndexSet(frame, obj IObject) (err error) {
	return (*o.valuePtr).IndexSet(frame, obj)
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (o *ObjectPointer) Iterate(frame int) IIterator {
	return (*o.valuePtr).Iterate(frame)
}

// CanIterate determines if the object can be iterated over and returns false for this implementation.
func (o *ObjectPointer) CanIterate() bool {
	return (*o.valuePtr).CanIterate()
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *ObjectPointer) Call(frame int, v ...IObject) (ret IObject, err error) {
	return (*o.valuePtr).Call(frame, v...)
}

// CanCall determines if the object can be invoked as a callable. Returns false for non-callable objects.
func (o *ObjectPointer) CanCall() bool {
	return (*o.valuePtr).CanCall()
}

// Length returns the length of the Int object.
func (o *ObjectPointer) Length() int {
	return (*o.valuePtr).Length()
}

// Value returns the internal IObject pointer stored in the ObjectPointer instance.
func (o *ObjectPointer) Value() *IObject {
	return o.valuePtr
}

// TypeName returns the type name of the ObjectPointer as a string.
func (o *ObjectPointer) TypeName() string {
	return (*o.valuePtr).TypeName()
}

// Copy creates and returns a duplicate of the object implementing the IObject interface.
func (o *ObjectPointer) Copy(_ int, _ int) IObject {
	return o
}

// Falsy returns true if the value of the ObjectPointer is nil.
func (o *ObjectPointer) Falsy() bool {
	return o.valuePtr == nil
}

// Equals checks if the current ObjectPointer is equal to the provided IObject by comparing their memory addresses.
func (o *ObjectPointer) Equals(x IObject) bool {
	return o == x
}

// release resets the current value pointer to an undefined value and updates its frame to static, returning the previous value and frame.
//func (o *ObjectPointer) release2() (IObject, bool) {
//	retPtr := *o.valuePtr
//	release := false
//	if retPtr.Frame() != FrameStatic {
//		release = retPtr.ReleaseRef() <= 0
//	}
//	undefined := o.GateKeeper().UndefinedValue()
//	o.valuePtr = &undefined
//	return retPtr, release
//}

// acquire updates the ObjectPointer with a new IObject reference, sets its frame, and marks the object as static.
func (o *ObjectPointer) acquire(value *IObject) {
	o.valuePtr = value
	if (*o.valuePtr).Frame() != FrameStatic {
		(*o.valuePtr).AddRef()
	}
}
