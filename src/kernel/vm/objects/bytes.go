package objects

import (
	"bytes"

	"github.com/markel1974/c64emu/src/kernel/vm/errors"
	"github.com/markel1974/c64emu/src/kernel/vm/tokens"
)

var (
	// MaxBytesLen is the maximum length for bytes value. Note this limit applies to all compiler/VM instances in the process.
	MaxBytesLen = 2147483647
)

// Bytes represents a byte array.
type Bytes struct {
	ObjectImpl
	Value []byte
}

func (o *Bytes) String() string {
	return string(o.Value)
}

// TypeName returns the name of the type.
func (o *Bytes) TypeName() string {
	return "bytes"
}

// BinaryOp returns another object that is the result of a given binary operator and a right-hand side object.
func (o *Bytes) BinaryOp(op tokens.Token, rhs Object) (Object, error) {
	switch op {
	case tokens.Add:
		switch rhs := rhs.(type) {
		case *Bytes:
			if len(o.Value)+len(rhs.Value) > MaxBytesLen {
				return nil, errors.ErrBytesLimit
			}
			return &Bytes{Value: append(o.Value, rhs.Value...)}, nil
		}
	default:
		return nil, errors.ErrInvalidOperator
	}
	return nil, errors.ErrInvalidOperator
}

// Copy returns a copy of the type.
func (o *Bytes) Copy() Object {
	return &Bytes{Value: append([]byte{}, o.Value...)}
}

// IsFalsy returns true if the value of the type is falsy.
func (o *Bytes) IsFalsy() bool {
	return len(o.Value) == 0
}

// Equals returns true if the value of the type is equal to the value of another object.
func (o *Bytes) Equals(x Object) bool {
	t, ok := x.(*Bytes)
	if !ok {
		return false
	}
	return bytes.Equal(o.Value, t.Value)
}

// IndexGet returns an element (as Int) at a given index.
func (o *Bytes) IndexGet(index Object) (res Object, err error) {
	intIdx, ok := index.(*Int)
	if !ok {
		err = errors.ErrInvalidIndexType
		return
	}
	idxVal := int(intIdx.Value)
	if idxVal < 0 || idxVal >= len(o.Value) {
		res = UndefinedValue
		return
	}
	res = &Int{Value: int64(o.Value[idxVal])}
	return
}

// Iterate creates a bytes iterator.
func (o *Bytes) Iterate() Iterator {
	return &BytesIterator{
		v: o.Value,
		l: len(o.Value),
	}
}

// CanIterate returns whether the Object can be Iterated.
func (o *Bytes) CanIterate() bool {
	return true
}

// BytesIterator represents an iterator for a string.
type BytesIterator struct {
	ObjectImpl
	v []byte
	i int
	l int
}

// TypeName returns the name of the type.
func (i *BytesIterator) TypeName() string {
	return "bytes-iterator"
}

func (i *BytesIterator) String() string {
	return "<bytes-iterator>"
}

// Equals returns true if the value of the type is equal to the value of another object.
func (i *BytesIterator) Equals(Object) bool {
	return false
}

// Copy returns a copy of the type.
func (i *BytesIterator) Copy() Object {
	return &BytesIterator{v: i.v, i: i.i, l: i.l}
}

// Next returns true if there are more elements to iterate.
func (i *BytesIterator) Next() bool {
	i.i++
	return i.i <= i.l
}

// Key returns the key or index value of the current element.
func (i *BytesIterator) Key() Object {
	return &Int{Value: int64(i.i - 1)}
}

// Value returns the value of the current element.
func (i *BytesIterator) Value() Object {
	return &Int{Value: int64(i.v[i.i-1])}
}

// ToByteSlice will try to convert object o to []byte value.
func ToByteSlice(o Object) (v []byte, ok bool) {
	switch o := o.(type) {
	case *Bytes:
		v = o.Value
		ok = true
	case *String:
		v = []byte(o.Value)
		ok = true
	}
	return
}
