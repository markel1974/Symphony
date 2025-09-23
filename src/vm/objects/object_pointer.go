package objects

import (
	"bytes"
	"encoding/gob"
)

func init() {
	gob.Register(&ObjectPointer{})
}

// ObjectPointer is a wrapper around a pointer to an IObject, allowing additional behaviors and encapsulation of the Code.
// It embeds Object, inheriting default behaviors for the IObject interface methods.
// The Code field holds the actual IObject instance being wrapped.
type ObjectPointer struct {
	IAllocator
	data *IObject
}

// NewObjectPointer creates a new ObjectPointer instance wrapping the provided IObject pointer.// NewObjectPointer creates a new ObjectPointer instance with the provided IObject Code.
func newObjectPointer(allocator IAllocator, value *IObject) IObject {
	ptr := &ObjectPointer{
		IAllocator: allocator,
	}
	if value != nil {
		ptr.acquire(value)
	} else {
		undefined := allocator.GateKeeper().UndefinedValue()
		ptr.data = &undefined
	}
	return ptr
}

// setAllocator sets the allocator for the instance, defining its memory management and lifecycle behavior.
func (o *ObjectPointer) setAllocator(allocator IAllocator) {
	o.IAllocator = allocator
}

// AsInterface converts the object into a generic interface{} type and returns the underlying data.
func (o *ObjectPointer) AsInterface() interface{} {
	return (*o.data).AsInterface()
}

// AsBool returns the boolean representation of the ObjectPointer, defaulting to false.
func (o *ObjectPointer) AsBool() bool {
	return (*o.data).AsBool()
}

// AsInt64 returns the len of the array as an int64 Code.
func (o *ObjectPointer) AsInt64() int64 {
	return (*o.data).AsInt64()
}

// AsFloat64 returns the len of the array as an int64 Code.
func (o *ObjectPointer) AsFloat64() float64 {
	return (*o.data).AsFloat64()
}

// AsBytes converts the object elements into a single concatenated slice of bytes by calling AsBytes on each element.
func (o *ObjectPointer) AsBytes() []byte {
	return (*o.data).AsBytes()
}

// AsString returns the string representation of the ObjectPointer instance.
func (o *ObjectPointer) AsString() string {
	return (*o.data).AsString()
}

// Nil checks if the object is nil and always returns false.
func (o *ObjectPointer) Nil() bool {
	return false
}

// AssignValue sets the current object to the provided IObject, returning ErrNotAssignable if the operation is not supported.
func (o *ObjectPointer) AssignValue(v IObject) error {
	return (*o.data).AssignValue(v)
}

// LogicalOp performs a logical operation with the given operator and RHS object, returning the result or an error.
func (o *ObjectPointer) LogicalOp(frame int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		switch op {
		case OperatorLogicalEq:
			if o.data == nil {
				return o.GateKeeper().TrueValue(), nil
			} else {
				return o.GateKeeper().FalseValue(), nil
			}
		case OperatorLogicalNotEq:
			if o.data == nil {
				return o.GateKeeper().FalseValue(), nil
			} else {
				return o.GateKeeper().TrueValue(), nil
			}
		default:
			return o.GateKeeper().UndefinedValue(), ErrInvalidOperator
		}
	}
	return (*o.data).LogicalOp(frame, op, rhsIn)
}

// ArithmeticOp performs an arithmetic operation with the given operator and right-hand-side operand and returns the result.
// Returns an error if the operation is invalid.
func (o *ObjectPointer) ArithmeticOp(frame int, obj ArithmeticOperator, rhs IObject) (IObject, error) {
	return (*o.data).ArithmeticOp(frame, obj, rhs)
}

// UnaryOp applies the given unary operator to the ObjectPointer's wrapped data and returns the resulting IObject or an error.
func (o *ObjectPointer) UnaryOp(frame int, op UnaryOperator) (IObject, error) {
	return (*o.data).UnaryOp(frame, op)
}

// IndexGet attempts to retrieve a Code at the given index and returns an error if the object is not indexable.
func (o *ObjectPointer) IndexGet(frame int, obj IObject) (IObject, error) {
	return (*o.data).IndexGet(frame, obj)
}

// IndexSet attempts to assign a Code to an index in the object but always returns ErrIndexUnsupported,
// as this operation is unsupported.
func (o *ObjectPointer) IndexSet(frame, obj IObject) error {
	return (*o.data).IndexSet(frame, obj)
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (o *ObjectPointer) Iterate(frame int) IIterator {
	return (*o.data).Iterate(frame)
}

// Iterable determines if the object can be iterated over and returns false for this implementation.
func (o *ObjectPointer) Iterable() bool {
	return (*o.data).Iterable()
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *ObjectPointer) Call(frame int, v ...IObject) (retCount uint, ret IObject, err error) {
	return (*o.data).Call(frame, v...)
}

// Length returns the len of the Int object.
func (o *ObjectPointer) Length() int {
	return (*o.data).Length()
}

// Value returns the internal IObject pointer stored in the ObjectPointer instance.
func (o *ObjectPointer) Value() *IObject {
	return o.data
}

// TypeName returns the type name of the ObjectPointer as a string.
func (o *ObjectPointer) TypeName() string {
	return (*o.data).TypeName()
}

// Copy creates and returns a duplicate of the object implementing the IObject interface.
func (o *ObjectPointer) Copy(_ int, _ int) IObject {
	return o
}

// Falsy returns true if the Code of the ObjectPointer is nil.
func (o *ObjectPointer) Falsy() bool {
	return o.data == nil
}

// Equals checks if the current ObjectPointer is equal to the provided IObject by comparing their memory addresses.
func (o *ObjectPointer) Equals(x IObject) bool {
	return o == x
}

// acquire updates the ObjectPointer with a new IObject reference, sets its frame, and marks the object as static.
func (o *ObjectPointer) acquire(value *IObject) {
	o.data = value
	if (*o.data).Frame() != FrameStatic {
		(*o.data).AddRef()
	}
}

// Count returns the total number of elements in the instance and its sub-elements.
func (o *ObjectPointer) Count() int {
	return 1
}

func (o *ObjectPointer) SetObject(i IObject) {
	o.data = &i
}

// GobEncode serializes the ObjectPointer's data into a byte slice using gob encoding and returns the result or an error.
func (o *ObjectPointer) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(o.data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GobDecode decodes the provided byte slice into the ObjectPointer's data field using the gob package.
func (o *ObjectPointer) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	if err := decoder.Decode(&o.data); err != nil {
		return err
	}
	return nil
}
