package objects

import (
	"bytes"
)

// MaxBytesLen is the maximum allowed size for byte slices across all instances, ensuring consistency in size limits.
const (
	// MaxBytesLen is the maximum length for bytes values. Note this limit applies to all compiler/VM instances in the process.
	MaxBytesLen = 2147483647
)

// Bytes represents a data type for handling a sequence of bytes.
// It embeds ObjectImpl and provides behaviors like indexing, iteration, and binary operations.
type Bytes struct {
	ObjectImpl
	values []byte
}

// NewBytes creates and returns a new Bytes object initialized with the provided byte slice.
func NewBytes(value []byte) *Bytes {
	return &Bytes{values: value}
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
	return "bytes"
}

// BinaryOp performs a binary operation on the Bytes object based on the specified operator and operand.
// Supports addition for concatenating two Bytes objects, ensuring the combined length does not exceed MaxBytesLen.
// Returns the resulting Bytes object or an error if the operation or operand is invalid.
func (o *Bytes) BinaryOp(op Operator, in IObject) (IObject, error) {
	switch op {
	case OperatorAdd:
		switch rhs := in.(type) {
		case *Bytes:
			if len(o.values)+len(rhs.values) > MaxBytesLen {
				return nil, ErrBytesLimit
			}
			return &Bytes{values: append(o.values, rhs.values...)}, nil
		}
	default:
		return nil, ErrInvalidOperator
	}
	return nil, ErrInvalidOperator
}

// Copy creates and returns a new `Bytes` object with a duplicated values slice, ensuring no reference sharing.
func (o *Bytes) Copy() IObject {
	return &Bytes{values: append([]byte{}, o.values...)}
}

// Falsy determines if the Bytes object is considered falsy by checking if it contains no values. Returns true if empty.
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
func (o *Bytes) IndexGet(index IObject) (res IObject, err error) {
	intIdx, ok := index.(*Int)
	if !ok {
		err = ErrInvalidIndexType
		return
	}
	idxVal := int(intIdx.value)
	if idxVal < 0 || idxVal >= len(o.values) {
		res = UndefinedValue
		return
	}
	res = NewInt(int64(o.values[idxVal]))
	return
}

// Iterate returns an iterator for the Bytes object, enabling sequential access to its byte values.
func (o *Bytes) Iterate() IIterator {
	return &BytesIterator{
		v: o.values,
		l: len(o.values),
	}
}

// CanIterate returns true if the object can be iterated over, otherwise false.
func (o *Bytes) CanIterate() bool {
	return true
}

// BytesIterator is an iterator for traversing elements of a byte slice, implementing the IIterator interface.
type BytesIterator struct {
	ObjectImpl
	v []byte
	i int
	l int
}

// TypeName returns the string representation of the type name, which is "bytes-iterator".
func (i *BytesIterator) TypeName() string {
	return "bytes-iterator"
}

// String returns the string representation of the BytesIterator.
func (i *BytesIterator) String() string {
	return "<bytes-iterator>"
}

// Equals checks whether the BytesIterator is equal to another object implementing the IObject interface.
func (i *BytesIterator) Equals(IObject) bool {
	return false
}

// Copy creates and returns a new instance of BytesIterator with the same state as the current instance.
func (i *BytesIterator) Copy() IObject {
	return &BytesIterator{v: i.v, i: i.i, l: i.l}
}

// Next advances the iterator to the next position and returns true if the current position is within bounds.
func (i *BytesIterator) Next() bool {
	i.i++
	return i.i <= i.l
}

// Key returns the current index of the iterator as an IObject, decremented by one from the internal index tracker.
func (i *BytesIterator) Key() IObject {
	return NewInt(int64(i.i - 1))
}

// Value returns the values of the current byte in the iteration as an IObject, wrapped in an Int struct.
func (i *BytesIterator) Value() IObject {
	return NewInt(int64(i.v[i.i-1]))
}

// ToByteSlice converts an IObject to a byte slice if the object is of type *Bytes or *String.
// It returns the converted byte slice and a boolean indicating success.
func ToByteSlice(o IObject) (v []byte, ok bool) {
	switch o := o.(type) {
	case *Bytes:
		v = o.values
		ok = true
	case *String:
		v = []byte(o.value)
		ok = true
	}
	return
}
