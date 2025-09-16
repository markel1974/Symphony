package objects

import (
	"bytes"
	"encoding/gob"
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
	Allocator
	values []byte
}

// NewBytes creates and returns a new Bytes object initialized with the provided byte slice.
func newBytes(factory IGateKeeper, frame int, value []byte) IObject {
	if len(value) > maxBytesLen {
		value = value[0:maxBytesLen]
	}
	return &Bytes{
		Allocator: Allocator{gk: factory, frame: frame},
		values:    value,
	}
}

// AsBool returns true if the object is not empty, otherwise false.
func (o *Bytes) AsBool() bool {
	return len(o.values) > 0
}

// AsInt64 returns the length of the array as an int64 value.
func (o *Bytes) AsInt64() int64 {
	return int64(len(o.values))
}

// AsFloat64 returns the length of the array as an int64 value.
func (o *Bytes) AsFloat64() float64 {
	return float64(len(o.values))
}

// AsString returns the string representation of the Bytes object by converting its underlying byte slice to a string.
func (o *Bytes) AsString() string {
	return string(o.values)
}

// AssignValue assigns the values from another `Bytes` object to the current instance, returning an error if the types are incompatible.
func (o *Bytes) AssignValue(v IObject) error {
	target, ok := v.(*Bytes)
	if !ok {
		return ErrNotAssignable
	}
	o.values = target.values
	return nil
}

// Nil checks if the object is nil and always returns false.
func (o *Bytes) Nil() bool {
	return false
}

// IndexSet attempts to assign a value to an index in the object but always returns ErrIndexUnsupported,
// as this operation is unsupported.
func (o *Bytes) IndexSet(_, _ IObject) (err error) {
	return ErrIndexUnsupported
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *Bytes) Call(_ int, _ ...IObject) (retCount uint, ret IObject, err error) {
	return 0, nil, nil
}

// Length returns the length of the Bytes object, which is the number of bytes in the underlying byte slice.
func (o *Bytes) Length() int {
	return len(o.values)
}

// Value returns the underlying byte slice of the Bytes object.
func (o *Bytes) Value() []byte {
	return o.values
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
			if len(o.values)+len(rhs.values) > maxBytesLen {
				return nil, ErrLimitExceed
			}
			return o.GateKeeper().NewBytes(frame, append(o.values, rhs.values...)), nil
		default:
			if len(o.values)+1 > maxBytesLen {
				return nil, ErrLimitExceed
			}
			v := byte(in.AsInt64())
			return o.GateKeeper().NewBytes(frame, append(o.values, v)), nil
		}
	default:
		return nil, ErrInvalidOperator
	}
}

// LogicalOp performs a logical operation on the Bytes object using the specified operator and operand, returning an error.
func (o *Bytes) LogicalOp(_ int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		return logicalOpNil(o.gk, op)
	}
	return nil, ErrInvalidOperator
}

// Copy creates and returns a new `Bytes` object with a duplicated values slice, ensuring no reference sharing.
func (o *Bytes) Copy(frame int, _ int) IObject {
	return o.GateKeeper().NewBytes(frame, append([]byte{}, o.values...))
}

// Falsy determines if the Bytes object is considered falsy by checking if it contains no values. Returns true if empty.
func (o *Bytes) Falsy() bool {
	return len(o.values) == 0
}

// Equals checks if the current Bytes object is equal to another IObject of type *Bytes, comparing their byte values.
func (o *Bytes) Equals(x IObject) bool {
	t, ok := x.(*Bytes)
	if !ok {
		return false
	}
	return bytes.Equal(o.values, t.values)
}

// IndexGet retrieves the values at the specified index from the Bytes object and returns an error if the index is invalid.
func (o *Bytes) IndexGet(frame int, index IObject) (IObject, error) {
	intIdx, ok := index.(*Int)
	if !ok {
		return nil, ErrIndexInvalidType
	}
	idxVal := int(intIdx.value)
	if idxVal < 0 || idxVal >= len(o.values) {
		return o.gk.UndefinedValue(), nil
	}
	return o.gk.NewInt(frame, int64(o.values[idxVal])), nil
}

// Iterable returns true if the object can be iterated over, otherwise false.
func (o *Bytes) Iterable() bool {
	return true
}

// Iterate returns an iterator for the Bytes object, enabling sequential access to its byte values.
func (o *Bytes) Iterate(frame int) IIterator {
	return o.GateKeeper().NewBytesIterator(frame, o.values, 0)
}

// Count returns the total number of elements in the instance and its sub-elements.
func (o *Bytes) Count() int {
	return 1
}
