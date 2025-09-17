package objects

import (
	"bytes"
	"encoding/gob"
	"strconv"
)

const (
	IntType = "int"
)

func init() {
	gob.Register(&Int{})
}

// Int represents an integer type with a 64-bit Code and methods for operations, equality, and object behavior.
type Int struct {
	IAllocator
	data int64
}

// NewInt creates and returns a new instance of the Int struct initialized with the specified int64 Code.
func newInt(allocator IAllocator, value int64) IObject {
	return &Int{
		IAllocator: allocator,
		data:       value,
	}
}

// setAllocator sets the allocator for the instance, defining its memory management and lifecycle behavior.
func (o *Int) setAllocator(allocator IAllocator) {
	o.IAllocator = allocator
}

// AsInterface converts the object into a generic interface{} type and returns the underlying data.
func (o *Int) AsInterface() interface{} {
	return o.data
}

// AsBool converts the integer Code to a boolean, returning true if the Code is non-zero, otherwise false.
func (o *Int) AsBool() bool {
	return o.data != 0
}

// AsInt64 returns the len of the array as an int64 Code.
func (o *Int) AsInt64() int64 {
	return o.data
}

// AsFloat64 returns the len of the array as an int64 Code.
func (o *Int) AsFloat64() float64 {
	return float64(o.data)
}

// AsString returns the string representation of the Int Code using base 10 format.
func (o *Int) AsString() string {
	return strconv.FormatInt(o.data, 10)
}

// AssignValue assigns the Code of another IObject to the current Int object if the type is compatible, otherwise returns an error.
func (o *Int) AssignValue(v IObject) error {
	target, ok := v.(*Int)
	if !ok {
		return ErrNotAssignable
	}
	o.data = target.data
	return nil
}

// Nil checks if the object is nil and always returns false.
func (o *Int) Nil() bool {
	return false
}

// IndexGet attempts to retrieve a Code at the given index and returns an error if the object is not indexable.
func (o *Int) IndexGet(_ int, _ IObject) (IObject, error) {
	return o.GateKeeper().UndefinedValue(), ErrIndexNotIndexable
}

// IndexSet attempts to assign a Code to an index in the object but always returns ErrIndexUnsupported,
// as this operation is unsupported.
func (o *Int) IndexSet(_, _ IObject) error {
	return ErrIndexUnsupported
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (o *Int) Iterate(_ int) IIterator {
	return nil
}

// Iterable determines if the object can be iterated over and returns false for this implementation.
func (o *Int) Iterable() bool {
	return false
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *Int) Call(_ int, _ ...IObject) (retCount uint, ret IObject, err error) {
	return 0, nil, nil
}

// Length returns the len of the Int object.
func (o *Int) Length() int {
	return 0
}

// Value returns the underlying int64 Code of the Int object.
func (o *Int) Value() int64 {
	return o.data
}

// SetValue sets the underlying int64 Code of the Int object.
func (o *Int) SetValue(value int64) {
	o.data = value
}

// TypeName returns the name of the type as a string, which is "int" for this object.
func (o *Int) TypeName() string {
	return IntType
}

// LogicalOp performs a logical operation between the object's Code and a given operand, using the specified operator.
func (o *Int) LogicalOp(_ int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		return logicalOpNil(o.GateKeeper(), op)
	}
	ret, err := logicalOpInt64(o.data, op, rhsIn.AsInt64())
	if err != nil {
		return nil, err
	}
	if ret {
		return o.GateKeeper().TrueValue(), nil
	}
	return o.GateKeeper().FalseValue(), nil
}

// ArithmeticOp performs a binary arithmetic operation on the Int object using the specified operator and right-hand operand.
func (o *Int) ArithmeticOp(frame int, op ArithmeticOperator, rhsIn IObject) (IObject, error) {
	ret, err := arithmeticOpInt64(o.data, op, rhsIn.AsInt64())
	if err != nil {
		return nil, err
	}
	if ret == o.data {
		return o, nil
	}
	return o.GateKeeper().NewInt(frame, ret), nil
}

// Copy creates and returns a new instance of the Int object with the same Code as the current instance.
func (o *Int) Copy(frame int, _ int) IObject {
	return o.GateKeeper().NewInt(frame, o.data)
}

// Falsy checks whether the integer Code is considered falsy. Returns true if the Code is 0, otherwise false.
func (o *Int) Falsy() bool {
	return o.data == 0
}

// Equals checks if the current Int object is equal to another IObject of type Int by comparing their Code.
func (o *Int) Equals(x IObject) bool {
	t, ok := x.(*Int)
	if !ok {
		return false
	}
	return o.data == t.data
}

// Count returns the total number of elements in the instance and its sub-elements.
func (o *Int) Count() int {
	return 1
}

// GobEncode serializes the Int's data into a byte slice using gob encoding and returns the result or an error.
func (o *Int) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(o.data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GobDecode decodes the provided byte slice into the Int's data field using the gob package.
func (o *Int) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	if err := decoder.Decode(&o.data); err != nil {
		return err
	}
	return nil
}
