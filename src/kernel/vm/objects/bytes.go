package objects

import (
	"bytes"
)

const (
	BytesType = "bytes"
)

// MaxBytesLen is the maximum allowed size for byte slices across all instances, ensuring consistency in size limits.
const (
	MaxBytesLen = 2147483647
)

// Bytes represents a data type for handling a sequence of bytes.
// It embeds Object and provides behaviors like indexing, iteration, and binary operations.
type Bytes struct {
	*Object
	values []byte
}

// NewBytes creates and returns a new Bytes object initialized with the provided byte slice.
func newBytes(factory *GateKeeper, frame int, value []byte) *Bytes {
	return &Bytes{
		Object: factory.NewObject(frame),
		values: value,
	}
}

// Length returns the length of the Bytes object, which is the number of bytes in the underlying byte slice.
func (o *Bytes) Length() int {
	return len(o.values)
}

// Value returns the underlying byte slice of the Bytes object.
func (o *Bytes) Value() []byte {
	return o.values
}

// String returns the string representation of the Bytes object by converting its underlying byte slice to a string.
func (o *Bytes) String() string {
	return string(o.values)
}

// TypeName returns the name of the type as a string, which is "bytes".
func (o *Bytes) TypeName() string {
	return BytesType
}

// BinaryOp performs a binary operation on the Bytes object based on the specified operator and operand.
// Supports addition for concatenating two Bytes objects, ensuring the combined length does not exceed MaxBytesLen.
// Returns the resulting Bytes object or an error if the operation or operand is invalid.
func (o *Bytes) BinaryOp(frame int, op Operator, in IObject) (IObject, error) {
	switch op {
	case OperatorAdd:
		switch rhs := in.(type) {
		case *Bytes:
			if len(o.values)+len(rhs.values) > MaxBytesLen {
				return nil, ErrBytesLimit
			}
			return o.Factory().NewBytes(frame, append(o.values, rhs.values...)), nil
		}
	default:
		return nil, ErrInvalidOperator
	}
	return nil, ErrInvalidOperator
}

// Copy creates and returns a new `Bytes` object with a duplicated values slice, ensuring no reference sharing.
func (o *Bytes) Copy(frame int) IObject {
	return o.Factory().NewBytes(frame, append([]byte{}, o.values...))
}

// Boolean determines if the Bytes object is considered falsy by checking if it contains no values. Returns true if empty.
func (o *Bytes) Boolean() bool {
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
func (o *Bytes) IndexGet(frame int, index IObject) (res IObject, err error) {
	intIdx, ok := index.(*Int)
	if !ok {
		err = ErrInvalidIndexType
		return
	}
	idxVal := int(intIdx.value)
	if idxVal < 0 || idxVal >= len(o.values) {
		res = o.Factory().UndefinedValue()
		return
	}
	res = o.Factory().NewInt(frame, int64(o.values[idxVal]))
	return
}

// CanIterate returns true if the object can be iterated over, otherwise false.
func (o *Bytes) CanIterate() bool {
	return true
}

// Iterate returns an iterator for the Bytes object, enabling sequential access to its byte values.
func (o *Bytes) Iterate(frame int) IIterator {
	return o.Factory().NewBytesIterator(frame, o.values)
}
