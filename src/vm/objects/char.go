package objects

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
)

const (
	CharType = "char"
)

func init() {
	gob.Register(&Char{})
}

// Char represents a character type, encapsulating a single rune data and inheriting behavior from Object.
type Char struct {
	IAllocator
	data rune
}

// NewChar creates and returns a new Char object with the specified rune Code.
func newChar(allocator IAllocator, value rune) IObject {
	return &Char{
		IAllocator: allocator,
		data:       value,
	}
}

// setAllocator sets the allocator for the instance, defining its memory management and lifecycle behavior.
func (o *Char) setAllocator(allocator IAllocator) {
	o.IAllocator = allocator
}

// AsInterface converts the object into a generic interface{} type and returns the underlying data.
func (o *Char) AsInterface() interface{} {
	return o.data
}

// AsBool returns true if the object is not empty, otherwise false.
func (o *Char) AsBool() bool {
	return o.data != 0
}

// AsInt64 returns the len of the array as an int64 data.
func (o *Char) AsInt64() int64 {
	return int64(o.data)
}

// AsFloat64 returns the len of the array as an int64 data.
func (o *Char) AsFloat64() float64 {
	return float64(o.data)
}

// AsBytes converts the object elements into a single concatenated slice of bytes by calling AsBytes on each element.
func (o *Char) AsBytes() []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(o.data))
	return b
}

// AsString returns the string representation of the Char object's data.
func (o *Char) AsString() string {
	return string(o.data)
}

// AssignValue assigns the data of another IObject to the current Char object if the type is compatible, otherwise returns an error.
func (o *Char) AssignValue(v IObject) error {
	target, ok := v.(*Char)
	if !ok {
		return ErrNotAssignable
	}
	o.data = target.data
	return nil
}

// Nil checks if the object is nil and always returns false.
func (o *Char) Nil() bool {
	return false
}

// IndexGet attempts to retrieve a data at the given index and returns an error if the object is not indexable.
func (o *Char) IndexGet(_ int, _ IObject) (IObject, error) {
	return o.GateKeeper().UndefinedValue(), ErrIndexNotIndexable
}

// IndexSet attempts to assign a data to an index in the object but always returns ErrIndexUnsupported,
// as this operation is unsupported.
func (o *Char) IndexSet(_, _ IObject) (err error) {
	return ErrIndexUnsupported
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (o *Char) Iterate(_ int) IIterator {
	return nil
}

// Iterable determines if the object can be iterated over and returns false for this implementation.
func (o *Char) Iterable() bool {
	return false
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *Char) Call(_ int, _ ...IObject) (retCount uint, ret IObject, err error) {
	return 0, nil, nil
}

// Length returns the len of the Int object.
func (o *Char) Length() int {
	return 0
}

// Value returns the rune data stored in the Char object.
func (o *Char) Value() rune {
	return o.data
}

// TypeName returns the name of the type as a string.
func (o *Char) TypeName() string {
	return CharType
}

// LogicalOp performs a logical operation between the current Char object and another IObject using the specified operator.
func (o *Char) LogicalOp(_ int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		return logicalOpNil(o.GateKeeper(), op)
	}
	ret, err := logicalOpInt64(int64(o.data), op, rhsIn.AsInt64())
	if err != nil {
		return nil, err
	}
	if ret {
		return o.GateKeeper().TrueValue(), nil
	}
	return o.GateKeeper().FalseValue(), nil
}

// ArithmeticOp applies the specified arithmetic operation between a Char object and another IObject, returning the result.
// Returns an error if the operation is invalid or unsupported.
func (o *Char) ArithmeticOp(frame int, op ArithmeticOperator, rhsIn IObject) (IObject, error) {
	ret, err := arithmeticOpInt64(int64(o.data), op, rhsIn.AsInt64())
	if err != nil {
		return nil, err
	}
	if ret == int64(o.data) {
		return o, nil
	}
	return o.GateKeeper().NewChar(frame, rune(ret)), nil
}

// UnaryOp applies a unary operator to the data in the Char object and returns a new IObject or an error if unsupported.
func (o *Char) UnaryOp(frame int, op UnaryOperator) (IObject, error) {
	val := o.AsInt64()
	r, err := unaryOpInt64(op, val)
	if err != nil {
		return nil, err
	}
	return o.GateKeeper().NewChar(frame, rune(r)), nil
}

// Copy creates and returns a new instance of the Char object with the same data.
func (o *Char) Copy(frame int, _ int) IObject {
	return o.GateKeeper().NewChar(frame, o.data)
}

// Falsy checks whether the Char object represents a falsy state, returning true if the underlying data is 0.
func (o *Char) Falsy() bool {
	return o.data == 0
}

// Equals checks if the current Char object is equal to another IObject. Returns true if both objects are equal.
func (o *Char) Equals(x IObject) bool {
	t, ok := x.(*Char)
	if !ok {
		return false
	}
	return o.data == t.data
}

// Count returns the total number of elements in the instance and its sub-elements.
func (o *Char) Count() int {
	return 1
}

// SetValue updates the data stored in the Char object with the given rune.
func (o *Char) SetValue(value rune) {
	o.data = value
}

// GobEncode serializes the Char's data into a byte slice using gob encoding and returns the result or an error.
func (o *Char) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(o.data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GobDecode decodes the provided byte slice into the Char's data field using the gob package.
func (o *Char) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	if err := decoder.Decode(&o.data); err != nil {
		return err
	}
	return nil
}
