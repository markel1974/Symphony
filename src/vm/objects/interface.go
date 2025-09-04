package objects

import "encoding/gob"

// init registers the Interface type with the gob package for encoding and decoding operations.
func init() {
	gob.Register(&Interface{})
}

//*Allocator
//Allocator: NewAllocator(gk, frame),

// Interface represents an object with contextual execution information and dynamic properties managed within a frame.
type Interface struct {
	Allocator
	value  IObject
	iTable map[string]IObject
}

// newInterface creates a new instance of Interface with the provided gk, frame ID, value, and interface table.
func newInterface(gk IGateKeeper, frame int, value IObject, itable map[string]IObject) IObject {
	return &Interface{
		Allocator: Allocator{gk: gk, frame: frame},
		value:     value,
		iTable:    itable,
	}
}

// AsBool converts and returns the Interface's underlying value as a boolean.
func (o *Interface) AsBool() bool {
	return o.AsBool()
}

// AsInt64 returns the length of the array as an int64 value.
func (o *Interface) AsInt64() int64 {
	return o.AsInt64()
}

// AsFloat64 returns the length of the array as an int64 value.
func (o *Interface) AsFloat64() float64 {
	return o.AsFloat64()
}

// AsString returns the string representation of the Interface instance by delegating to the underlying IObject value.
func (o *Interface) AsString() string {
	return o.value.AsString()
}

// AssignValue sets the current object to the provided IObject, returning ErrNotAssignable if the operation is not supported.
func (o *Interface) AssignValue(v IObject) error {
	return o.value.AssignValue(v)
}

// Nil checks if the object is nil and always returns false.
func (o *Interface) Nil() bool {
	return false
}

// TypeName returns the type name of the underlying IObject.
func (o *Interface) TypeName() string {
	return o.value.TypeName()
}

// Falsy determines whether the underlying value of the object evaluates to a falsy value, returning true if it does.
func (o *Interface) Falsy() bool {
	return o.value.Falsy()
}

// Equals compares the current object with another IObject and returns true if they are equal.
func (o *Interface) Equals(other IObject) bool {
	return o.value.Equals(other)
}

// Copy creates and returns a deep copy of the current object with the specified frame and depth.
func (o *Interface) Copy(frame int, depth int) IObject {
	return o.value.Copy(frame, depth)
}

// LogicalOp applies a logical operation (e.g., AND, OR) between the current object and a right-hand-side object.
// It returns the result of the operation or an error if the operation cannot be performed.
func (o *Interface) LogicalOp(frame int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		return logicalOpNil(o.gk, op)
	}
	return o.value.LogicalOp(frame, op, rhsIn)
}

// ArithmeticOp performs an arithmetic operation specified by the operator on the current object and the right-hand side operand.
// It uses the provided frame for managing execution context, returning the result or an error if the operation fails.
func (o *Interface) ArithmeticOp(frame int, op ArithmeticOperator, rightHandSide IObject) (IObject, error) {
	return o.value.ArithmeticOp(frame, op, rightHandSide)
}

// IndexGet retrieves the value at the specified index from the object, using the provided execution frame and index.
func (o *Interface) IndexGet(frame int, index IObject) (value IObject, err error) {
	return o.value.IndexGet(frame, index)
}

// IndexSet sets a value at the specified index in the IObject, returning an error if the operation fails.
func (o *Interface) IndexSet(index, value IObject) error {
	return o.value.IndexSet(index, value)
}

// Iterate returns an iterator for traversing over the elements of the value associated with the interface.
func (o *Interface) Iterate(frame int) IIterator {
	return o.value.Iterate(frame)
}

// Iterable determines whether the wrapped IObject supports iteration, returning true if it does, or false otherwise.
func (o *Interface) Iterable() bool {
	return o.value.Iterable()
}

// Call invokes the object with the specified frame and arguments, returning the result or an error if unsupported.
func (o *Interface) Call(_ int, _ ...IObject) (ret IObject, err error) {
	return nil, ErrUnsupportedOperation
}

// Length returns the length of the underlying value represented by the IObject interface.
func (o *Interface) Length() int {
	return o.value.Length()
}

// Value returns the underlying IObject associated with the Interface instance.
func (o *Interface) Value() IObject {
	return o.value
}

// ITable returns the internal map representation `iTable` of the Interface instance.
func (o *Interface) ITable() map[string]IObject {
	return o.iTable
}
