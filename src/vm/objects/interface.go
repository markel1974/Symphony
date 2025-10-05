package objects

import (
	"bytes"
	"encoding/gob"
	"reflect"
)

// init registers the Interface type with the gob package for encoding and decoding operations.
func init() {
	gob.Register(&Interface{})
}

// Interface represents an object with contextual execution information and dynamic properties managed within a frame.
type Interface struct {
	IAllocator
	data   IObject
	iTable map[string]IObject
}

// newInterface creates a new instance of Interface with the provided gk, frame ID, Code, and interface table.
func newInterface(allocator IAllocator, value IObject, iTable map[string]IObject) IObject {
	return &Interface{
		IAllocator: allocator,
		data:       value,
		iTable:     iTable,
	}
}

// setAllocator sets the allocator for the instance, defining its memory management and lifecycle behavior.
func (o *Interface) setAllocator(allocator IAllocator) {
	o.IAllocator = allocator
}

// AsInterface converts the object into a generic interface{} type and returns the underlying data.
func (o *Interface) AsInterface() interface{} {
	return nil
}

// AsValue attempts to convert the object to a reflect.Value of the specified target type and returns it along with a success flag.
func (o *Interface) AsValue(target reflect.Type) (reflect.Value, bool) {
	return _reflect(o.data, target)
}

// AsBool converts and returns the Interface's underlying Code as a boolean.
func (o *Interface) AsBool() bool {
	return o.AsBool()
}

// AsInt64 returns the len of the array as an int64 Code.
func (o *Interface) AsInt64() int64 {
	return o.AsInt64()
}

// AsFloat64 returns the len of the array as an int64 Code.
func (o *Interface) AsFloat64() float64 {
	return o.AsFloat64()
}

// AsBytes converts the object elements into a single concatenated slice of bytes by calling AsBytes on each element.
func (o *Interface) AsBytes() []byte {
	return nil
}

// AsString returns the string representation of the Interface instance by delegating to the underlying IObject Code.
func (o *Interface) AsString() string {
	return o.data.AsString()
}

// AssignValue sets the current object to the provided IObject, returning ErrNotAssignable if the operation is not supported.
func (o *Interface) AssignValue(v IObject) error {
	return o.data.AssignValue(v)
}

// Nil checks if the object is nil and always returns false.
func (o *Interface) Nil() bool {
	return false
}

// TypeName returns the type name of the underlying IObject.
func (o *Interface) TypeName() string {
	return o.data.TypeName()
}

// Falsy determines whether the underlying Code of the object evaluates to a falsy Code, returning true if it does.
func (o *Interface) Falsy() bool {
	return o.data.Falsy()
}

// Equals compares the current object with another IObject and returns true if they are equal.
func (o *Interface) Equals(other IObject) bool {
	return o.data.Equals(other)
}

// Copy creates and returns a deep copy of the current object with the specified frame and depth.
func (o *Interface) Copy(frame int, depth int) IObject {
	return o.data.Copy(frame, depth)
}

// LogicalOp applies a logical operation (e.g., AND, OR) between the current object and a right-hand-side object.
// It returns the result of the operation or an error if the operation cannot be performed.
func (o *Interface) LogicalOp(frame int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		ret, err := logicalOpNil(op)
		return o.GateKeeper().FromBoolError(ret, err)
	}
	return o.data.LogicalOp(frame, op, rhsIn)
}

// ArithmeticOp performs an arithmetic operation specified by the operator on the current object and the right-hand side operand.
// It uses the provided frame for managing execution context, returning the result or an error if the operation fails.
func (o *Interface) ArithmeticOp(frame int, op ArithmeticOperator, rightHandSide IObject) (IObject, error) {
	return o.data.ArithmeticOp(frame, op, rightHandSide)
}

// UnaryOp applies a unary operation specified by the operator to the current object using the given execution frame.
// It returns the result of the operation or an error if the operation is not supported.
func (o *Interface) UnaryOp(frame int, op UnaryOperator) (IObject, error) {
	return o.data.UnaryOp(frame, op)
}

// IndexGet retrieves the Code at the specified index from the object, using the provided execution frame and index.
func (o *Interface) IndexGet(frame int, index IObject) (IObject, error) {
	return o.data.IndexGet(frame, index)
}

// IndexSet sets a Code at the specified index in the IObject, returning an error if the operation fails.
func (o *Interface) IndexSet(index, value IObject) error {
	return o.data.IndexSet(index, value)
}

// Iterate returns an iterator for traversing over the elements of the Code associated with the interface.
func (o *Interface) Iterate(frame int) IIterator {
	return o.data.Iterate(frame)
}

// Iterable determines whether the wrapped IObject supports iteration, returning true if it does, or false otherwise.
func (o *Interface) Iterable() bool {
	return o.data.Iterable()
}

// Call invokes the object with the specified frame and arguments, returning the result or an error if unsupported.
func (o *Interface) Call(_ int, _ ...IObject) (retCount uint, ret IObject, err error) {
	return 0, nil, nil
}

// Length returns the len of the underlying Code represented by the IObject interface.
func (o *Interface) Length() int {
	return o.data.Length()
}

// Value returns the underlying IObject associated with the Interface instance.
func (o *Interface) Value() IObject {
	return o.data
}

func (o *Interface) Method(name string) (IObject, bool) {
	m, ok := o.iTable[name]
	if !ok || m == nil {
		return o.GateKeeper().UndefinedValue(), false
	}
	return m, ok
}

// Count returns the total number of elements in the instance and its sub-elements.
func (o *Interface) Count() int {
	return 1
}

// GobEncode serializes the Interface's data into a byte slice using gob encoding and returns the result or an error.
func (o *Interface) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(o.iTable); err != nil {
		return nil, err
	}
	if err := encoder.Encode(o.data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GobDecode decodes the provided byte slice into the Interface's data field using the gob package.
func (o *Interface) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	if err := decoder.Decode(&o.iTable); err != nil {
		return err
	}
	if err := decoder.Decode(&o.data); err != nil {
		return err
	}
	return nil
}
