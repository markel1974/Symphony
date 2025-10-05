package objects

import (
	"bytes"
	"encoding/gob"
	"reflect"
)

// BoolType defines the string representation of the boolean type. It is used as the type name for boolean objects.
const (
	BoolType = "bool"
)

// init registers the Bool type with the gob package for encoding and decoding.
func init() {
	gob.Register(&Bool{})
}

// Bool represents a boolean type with frame-specific and gatekeeper-managed data assignments.
type Bool struct {
	IAllocator
	data bool
}

// newBool creates and returns a new instance of Bool, initializing it with the given frame and boolean Code.
func newBool(allocator IAllocator, value bool) IObject {
	return &Bool{
		IAllocator: allocator,
		data:       value,
	}
}

// setAllocator sets the allocator for the instance, defining its memory management and lifecycle behavior.
func (o *Bool) setAllocator(allocator IAllocator) {
	o.IAllocator = allocator
}

// AsInterface converts the object into a generic interface{} type and returns the underlying data.
func (o *Bool) AsInterface() interface{} {
	return o == o.GateKeeper().TrueValue()
}

// AsValue attempts to convert the ArrayIterator's data into a reflect.Value of the specified type. Returns false if invalid.
func (o *Bool) AsValue(target reflect.Type) (reflect.Value, bool) {
	return _reflect(o, target)
}

// AsBool returns the boolean value stored in the Bool object.
func (o *Bool) AsBool() bool {
	return o.data
}

// AsInt64 returns the len of the array as an int64 data.
func (o *Bool) AsInt64() int64 {
	if o.data {
		return 1
	}
	return 0
}

// AsFloat64 returns the len of the array as an int64 data.
func (o *Bool) AsFloat64() float64 {
	if o.data {
		return 1
	}
	return 0
}

// AsBytes converts the object elements into a single concatenated slice of bytes by calling AsBytes on each element.
func (o *Bool) AsBytes() []byte {
	if o == o.GateKeeper().TrueValue() {
		return []byte{1}
	} else {
		return []byte{0}
	}
}

// AsString returns the string representation of the Bool object, either "true" or "false" based on its boolean data.
func (o *Bool) AsString() string {
	if o.data {
		return "true"
	}
	return "false"
}

// AssignValue assigns the data of another Bool instance to the current instance. Returns an error if the input is not a Bool.
func (o *Bool) AssignValue(v IObject) error {
	target, ok := v.(*Bool)
	if !ok {
		return ErrNotAssignable
	}
	o.data = target.data
	return nil
}

// Nil checks if the object is nil and always returns false.
func (o *Bool) Nil() bool {
	return false
}

// LogicalOp performs a logical operation on the Bool instance using the provided operator and right-hand operand.
// Returns the result of the operation as an IObject or an error if the operation is invalid.
func (o *Bool) LogicalOp(_ int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		ret, err := logicalOpNil(op)
		return o.GateKeeper().FromBoolError(ret, err)
	}
	lhsValue := int64(0)
	if o.data {
		lhsValue = 1
	}
	ret, err := logicalOpInt64(lhsValue, op, rhsIn.AsInt64())
	return o.GateKeeper().FromBoolError(ret, err)
}

// ArithmeticOp performs an arithmetic operation between the Bool object and a given IObject using the specified operator.
// Returns the result as an IObject and an error if the operation is not valid or executable.
func (o *Bool) ArithmeticOp(_ int, op ArithmeticOperator, rhsIn IObject) (IObject, error) {
	lhsValue := int64(0)
	if o.data {
		lhsValue = 1
	}
	ret, err := arithmeticOpInt64(lhsValue, op, rhsIn.AsInt64())
	if err != nil {
		return nil, err
	}
	if ret != 0 {
		return o.GateKeeper().TrueValue(), nil
	}
	return o.GateKeeper().FalseValue(), nil
}

// UnaryOp performs a unary operation using the specified UnaryOperator. Returns a new object or an error.
func (o *Bool) UnaryOp(_ int, op UnaryOperator) (IObject, error) {
	ret, err := unaryOpBool(op, o.data)
	return o.GateKeeper().FromBoolError(ret, err)
}

// IndexGet retrieves the data at a given index from the Bool object, but always returns an error as Bool is not indexable.
func (o *Bool) IndexGet(_ int, _ IObject) (IObject, error) {
	return o.GateKeeper().UndefinedValue(), ErrIndexNotIndexable
}

// IndexSet attempts to set an index on the Bool object but always returns ErrIndexUnsupported as Bool is not indexable.
func (o *Bool) IndexSet(_, _ IObject) (err error) {
	return ErrIndexUnsupported
}

// Iterate returns nil as Bool does not support iteration.
func (o *Bool) Iterate(_ int) IIterator {
	return nil
}

// Iterable indicates whether the Bool object can be iterated. Always returns false.
func (o *Bool) Iterable() bool {
	return false
}

// Call invokes the Bool object as a callable function with the provided arguments, returning nil and no error.
func (o *Bool) Call(_ int, _ ...IObject) (retCount uint, ret IObject, err error) {
	return 0, nil, nil
}

// Length returns the len of the Bool object, which is always 0.
func (o *Bool) Length() int {
	return 0
}

// TypeName returns the type name of the Bool object as a string.
func (o *Bool) TypeName() string {
	return BoolType
}

// Copy creates and returns a new Bool instance with the same data and the specified execution frame.
func (o *Bool) Copy(frame int, _ int) IObject {
	return o.GateKeeper().NewBool(frame, o.data)
}

// Falsy determines if the Bool's data is logically false by returning the negation of its `Code` field.
func (o *Bool) Falsy() bool {
	return !o.data
}

// Equals checks if the current Bool object is equal to the provided IObject.
func (o *Bool) Equals(x IObject) bool {
	return o == x
}

// Count returns the total number of elements in the instance and its sub-elements.
func (o *Bool) Count() int {
	return 1
}

// SetValue assigns the provided boolean data to the Bool object's internal data field.
func (o *Bool) SetValue(v bool) {
	o.data = v
}

// GobEncode serializes the Bool's data into a byte slice using gob encoding and returns the result or an error.
func (o *Bool) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(o.data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GobDecode decodes the provided byte slice into the Bool's data field using the gob package.
func (o *Bool) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	if err := decoder.Decode(&o.data); err != nil {
		return err
	}
	return nil
}
