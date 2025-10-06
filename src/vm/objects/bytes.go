package objects

import (
	"bytes"
	"encoding/gob"
	"reflect"
)

func init() {
	gob.Register(&Bytes{})
}

const (
	BytesType = "bytes"
)

// Bytes represents a data type for handling a sequence of bytes.
// It embeds Object and provides behaviors like indexing, iteration, and binary operations.
type Bytes struct {
	IAllocator
	data []byte
}

// NewBytes creates and returns a new Bytes object initialized with the provided byte slice.
func newBytes(allocator IAllocator, value []byte) IObject {
	if len(value) > MaxBytesLen {
		value = value[0:MaxBytesLen]
	}
	return &Bytes{
		IAllocator: allocator,
		data:       value,
	}
}

// setAllocator sets the allocator for the instance, defining its memory management and lifecycle behavior.
func (o *Bytes) setAllocator(allocator IAllocator) {
	o.IAllocator = allocator
}

// AsBool returns true if the object is not empty, otherwise false.
func (o *Bytes) AsBool() bool {
	return len(o.data) > 0
}

// AsInterface converts the object into a generic interface{} type and returns the underlying data.
func (o *Bytes) AsInterface() interface{} {
	return o.data
}

// AsValue attempts to convert the Bytes object into a reflect.Value matching the given target type. Returns success status.
func (o *Bytes) AsValue(target reflect.Type) (reflect.Value, bool) {
	if target.Kind() == reflect.ValueOf(o.data).Kind() {
		return reflect.ValueOf(o.data), true
	}
	return reflect.Value{}, false
}

// AsInt64 returns the len of the array as an int64 data.
func (o *Bytes) AsInt64() int64 {
	return int64(len(o.data))
}

// AsFloat64 returns the len of the array as an int64 data.
func (o *Bytes) AsFloat64() float64 {
	return float64(len(o.data))
}

// AsBytes converts the object elements into a single concatenated slice of bytes by calling AsBytes on each element.
func (o *Bytes) AsBytes() []byte {
	return o.data
}

// AsString returns the string representation of the Bytes object by converting its underlying byte slice to a string.
func (o *Bytes) AsString() string {
	return string(o.data)
}

// AssignValue assigns the data from another `Bytes` object to the current instance, returning an error if the types are incompatible.
func (o *Bytes) AssignValue(v IObject) error {
	target, ok := v.(*Bytes)
	if !ok {
		return ErrNotAssignable
	}
	o.data = target.data
	return nil
}

// Nil checks if the object is nil and always returns false.
func (o *Bytes) Nil() bool {
	return false
}

// IndexSet attempts to assign a data to an index in the object but always returns ErrIndexUnsupported,
// as this operation is unsupported.
func (o *Bytes) IndexSet(_, _ IObject) (err error) {
	return ErrIndexUnsupported
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *Bytes) Call(_ int, _ ...IObject) (retCount uint, ret IObject, err error) {
	return 0, nil, nil
}

// Length returns the len of the Bytes object, which is the number of bytes in the underlying byte slice.
func (o *Bytes) Length() int {
	return len(o.data)
}

// GetValue returns the underlying byte slice of the Bytes object.
func (o *Bytes) GetValue() []byte {
	return o.data
}

// TypeName returns the name of the type as a string, which is "bytes".
func (o *Bytes) TypeName() string {
	return BytesType
}

// ArithmeticOp performs an arithmetic operation on the Bytes object using the specified operator and operand.
// Returns the result of the operation or an error if the operation is invalid or exceeds limitations.
func (o *Bytes) ArithmeticOp(frame int, op ArithmeticOperator, in IObject) (IObject, error) {
	switch op {
	case OperatorAdd:
		switch rhs := in.(type) {
		case *Bytes:
			if len(o.data)+len(rhs.data) > MaxBytesLen {
				return nil, ErrLimitExceed
			}
			return o.GateKeeper().NewBytes(frame, append(o.data, rhs.data...)), nil
		default:
			if len(o.data)+1 > MaxBytesLen {
				return nil, ErrLimitExceed
			}
			v := byte(in.AsInt64())
			return o.GateKeeper().NewBytes(frame, append(o.data, v)), nil
		}
	default:
		return nil, ErrInvalidOperator
	}
}

// LogicalOp performs a logical operation on the Bytes object using the specified operator and operand, returning an error.
func (o *Bytes) LogicalOp(_ int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		ret, err := logicalOpNil(op)
		return o.GateKeeper().FromBoolError(ret, err)
	}
	return nil, ErrInvalidOperator
}

// UnaryOp applies a unary operation and returns an error indicating an invalid operation.
func (o *Bytes) UnaryOp(_ int, _ UnaryOperator) (IObject, error) {
	return nil, ErrInvalidOperator
}

// Copy creates and returns a new `Bytes` object with a duplicated data slice, ensuring no reference sharing.
func (o *Bytes) Copy(frame int, _ int) IObject {
	return o.GateKeeper().NewBytes(frame, append([]byte{}, o.data...))
}

// Falsy determines if the Bytes object is considered falsy by checking if it contains no data. Returns true if empty.
func (o *Bytes) Falsy() bool {
	return len(o.data) == 0
}

// Equals checks if the current Bytes object is equal to another IObject of type *Bytes, comparing their byte data.
func (o *Bytes) Equals(x IObject) bool {
	t, ok := x.(*Bytes)
	if !ok {
		return false
	}
	return bytes.Equal(o.data, t.data)
}

// IndexGet retrieves the data at the specified index from the Bytes object and returns an error if the index is invalid.
func (o *Bytes) IndexGet(frame int, index IObject) (IObject, error) {
	intIdx, ok := index.(*Int)
	if !ok {
		return nil, ErrIndexInvalidType
	}
	idxVal := int(intIdx.data)
	if idxVal < 0 || idxVal >= len(o.data) {
		return o.GateKeeper().UndefinedValue(), nil
	}
	return o.GateKeeper().NewInt(frame, int64(o.data[idxVal])), nil
}

// Iterable returns true if the object can be iterated over, otherwise false.
func (o *Bytes) Iterable() bool {
	return true
}

// Iterate returns an iterator for the Bytes object, enabling sequential access to its byte data.
func (o *Bytes) Iterate(frame int) IIterator {
	return o.GateKeeper().NewBytesIterator(frame, o.data, 0)
}

// Count returns the total number of elements in the instance and its sub-elements.
func (o *Bytes) Count() int {
	return 1
}

// GobEncode serializes the Bool's data into a byte slice using gob encoding and returns the result or an error.
func (o *Bytes) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(o.data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GobDecode decodes the provided byte slice into the Bool's data field using the gob package.
func (o *Bytes) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	if err := decoder.Decode(&o.data); err != nil {
		return err
	}
	return nil
}
