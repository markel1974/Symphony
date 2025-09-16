package objects

import (
	"bytes"
	"encoding/gob"
	"math"
	"strconv"
)

const (
	FloatType = "float"
)

func init() {
	gob.Register(&Float{})
}

// Float represents a floating-point number and provides operations and behaviors specific to numeric types.
// It embeds Object to implement common interface methods and extends behavior where necessary.
// The data field holds the actual float64 data encapsulated by the Float type.
type Float struct {
	IAllocator
	data float64
}

// NewFloat creates and returns a pointer to a new Float object initialized with the specified float64 Code.
func newFloat(allocator IAllocator, value float64) IObject {
	return &Float{
		IAllocator: allocator,
		data:       value,
	}
}

// setAllocator sets the allocator for the instance, defining its memory management and lifecycle behavior.
func (o *Float) setAllocator(allocator IAllocator) {
	o.IAllocator = allocator
}

// AsBool returns true if the object is not empty, otherwise false.
func (o *Float) AsBool() bool {
	return o.data != 0
}

// AsInt64 returns the len of the array as an int64 data.
func (o *Float) AsInt64() int64 {
	return int64(o.data)
}

// AsFloat64 returns the len of the array as an int64 data.
func (o *Float) AsFloat64() float64 {
	return o.data
}

// AssignValue assigns the data of another IObject to the current Float object if the type is compatible, otherwise returns an error.
func (o *Float) AssignValue(v IObject) error {
	target, ok := v.(*Float)
	if !ok {
		return ErrNotAssignable
	}
	o.data = target.data
	return nil
}

// AsString returns the string representation of the Float object using its internal float64 data.
func (o *Float) AsString() string {
	return strconv.FormatFloat(o.data, 'f', -1, 64)
}

// Nil checks if the object is nil and always returns false.
func (o *Float) Nil() bool {
	return false
}

// IndexGet attempts to retrieve a data at the given index and returns an error if the object is not indexable.
func (o *Float) IndexGet(_ int, _ IObject) (IObject, error) {
	return o.GateKeeper().UndefinedValue(), ErrIndexNotIndexable
}

// IndexSet attempts to assign a data to an index in the object but always returns ErrIndexUnsupported,
// as this operation is unsupported.
func (o *Float) IndexSet(_, _ IObject) error {
	return ErrIndexUnsupported
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (o *Float) Iterate(_ int) IIterator {
	return nil
}

// Iterable determines if the object can be iterated over and returns false for this implementation.
func (o *Float) Iterable() bool {
	return false
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *Float) Call(_ int, _ ...IObject) (retCount uint, ret IObject, err error) {
	return 0, nil, nil
}

// Length returns the len of the Int object.
func (o *Float) Length() int {
	return 0
}

func (o *Float) Value() float64 {
	return o.data
}

// TypeName returns the name of the type.
func (o *Float) TypeName() string {
	return FloatType
}

// LogicalOp performs a logical operation between the Float object and another IObject using the specified operator.
// Returns a boolean IObject representing the result or an error if the operation is invalid.
func (o *Float) LogicalOp(_ int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		return logicalOpNil(o.GateKeeper(), op)
	}
	ret, err := logicalOpFloat64(o.data, op, rhsIn.AsFloat64())
	if err != nil {
		return nil, err
	}
	if ret {
		return o.GateKeeper().TrueValue(), nil
	}
	return o.GateKeeper().FalseValue(), nil
}

// ArithmeticOp performs an arithmetic operation with the specified operator and returns the result or an error.
func (o *Float) ArithmeticOp(frame int, op ArithmeticOperator, rhsIn IObject) (IObject, error) {
	ret, err := arithmeticOpFloat64(o.data, op, rhsIn.AsFloat64())
	if err != nil {
		return nil, err
	}
	if ret == o.data {
		return o, nil
	}
	return o.GateKeeper().NewFloat(frame, ret), nil
}

// Copy creates and returns a new instance of the Float object, duplicating its current state.
func (o *Float) Copy(frame int, _ int) IObject {
	return o.GateKeeper().NewFloat(frame, o.data)
}

// Falsy determines if the float object is considered falsy, returning true if the data is NaN; otherwise, false.
func (o *Float) Falsy() bool {
	return math.IsNaN(o.data)
}

// Equals checks if the current Float object is equal to another IObject by comparing their internal float64 data.
func (o *Float) Equals(x IObject) bool {
	t, ok := x.(*Float)
	if !ok {
		return false
	}
	return o.data == t.data
}

// Count returns the total number of elements in the instance and its sub-elements.
func (o *Float) Count() int {
	return 1
}

// SetValue assigns a new float64 data to the internal data field of the Float object.
func (o *Float) SetValue(value float64) {
	o.data = value
}

// GobEncode serializes the Float's data into a byte slice using gob encoding and returns the result or an error.
func (o *Float) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(o.data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GobDecode decodes the provided byte slice into the Float's data field using the gob package.
func (o *Float) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	if err := decoder.Decode(&o.data); err != nil {
		return err
	}
	return nil
}
