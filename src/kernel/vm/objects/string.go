package objects

import (
	"strconv"

	"github.com/markel1974/c64emu/src/kernel/vm/errors"
)

var (
	// MaxStringLen is the maximum byte-length for string value. Note this
	// limit applies to all compiler/VM instances in the process.
	MaxStringLen = 2147483647
)

// String represents a string value.
type String struct {
	ObjectImpl
	Value   string
	runeStr []rune
}

func NewString(value string) *String {
	return &String{
		Value: value,
	}
}

// TypeName returns the name of the type.
func (o *String) TypeName() string {
	return "string"
}

func (o *String) String() string {
	return strconv.Quote(o.Value)
}

// BinaryOp returns another object that is the result of a given binary
// operator and a right-hand side object.
func (o *String) BinaryOp(op Operator, rhs Object) (Object, error) {
	switch op {
	case OperatorAdd:
		switch rhs := rhs.(type) {
		case *String:
			if len(o.Value)+len(rhs.Value) > MaxStringLen {
				return nil, errors.ErrStringLimit
			}
			return &String{Value: o.Value + rhs.Value}, nil
		default:
			rhsStr := rhs.String()
			if len(o.Value)+len(rhsStr) > MaxStringLen {
				return nil, errors.ErrStringLimit
			}
			return &String{Value: o.Value + rhsStr}, nil
		}
	case OperatorLess:
		switch rhs := rhs.(type) {
		case *String:
			if o.Value < rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		}
	case OperatorLessEq:
		switch rhs := rhs.(type) {
		case *String:
			if o.Value <= rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		}
	case OperatorGreater:
		switch rhs := rhs.(type) {
		case *String:
			if o.Value > rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		}
	case OperatorGreaterEq:
		switch rhs := rhs.(type) {
		case *String:
			if o.Value >= rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		}
	default:
		return nil, errors.ErrInvalidOperator
	}
	return nil, errors.ErrInvalidOperator
}

// IsFalsy returns true if the value of the type is falsy.
func (o *String) Falsy() bool {
	return len(o.Value) == 0
}

// Copy returns a copy of the type.
func (o *String) Copy() Object {
	return &String{Value: o.Value}
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *String) Equals(x Object) bool {
	t, ok := x.(*String)
	if !ok {
		return false
	}
	return o.Value == t.Value
}

// IndexGet returns a character at a given index.
func (o *String) IndexGet(index Object) (res Object, err error) {
	intIdx, ok := index.(*Int)
	if !ok {
		err = errors.ErrInvalidIndexType
		return
	}
	idxVal := int(intIdx.Value)
	if o.runeStr == nil {
		o.runeStr = []rune(o.Value)
	}
	if idxVal < 0 || idxVal >= len(o.runeStr) {
		res = UndefinedValue
		return
	}
	res = &Char{Value: o.runeStr[idxVal]}
	return
}

// Iterate creates a string iterator.
func (o *String) Iterate() Iterator {
	if o.runeStr == nil {
		o.runeStr = []rune(o.Value)
	}
	return &StringIterator{
		v: o.runeStr,
		l: len(o.runeStr),
	}
}

// CanIterate returns whether the Object can be Iterated.
func (o *String) CanIterate() bool {
	return true
}

// StringIterator represents an iterator for a string.
type StringIterator struct {
	ObjectImpl
	v []rune
	i int
	l int
}

// TypeName returns the name of the type.
func (i *StringIterator) TypeName() string {
	return "string-iterator"
}

func (i *StringIterator) String() string {
	return "<string-iterator>"
}

// IsFalsy returns true if the value of the type is falsy.
func (i *StringIterator) Falsy() bool {
	return true
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (i *StringIterator) Equals(Object) bool {
	return false
}

// Copy returns a copy of the type.
func (i *StringIterator) Copy() Object {
	return &StringIterator{v: i.v, i: i.i, l: i.l}
}

// Next returns true if there are more elements to iterate.
func (i *StringIterator) Next() bool {
	i.i++
	return i.i <= i.l
}

// Key returns the key or index value of the current element.
func (i *StringIterator) Key() Object {
	return &Int{Value: int64(i.i - 1)}
}

// Value returns the value of the current element.
func (i *StringIterator) Value() Object {
	return &Char{Value: i.v[i.i-1]}
}

// ToRune will try to convert object o to rune value.
func ToRune(o Object) (v rune, ok bool) {
	switch o := o.(type) {
	case *Int:
		v = rune(o.Value)
		ok = true
	case *Char:
		v = o.Value
		ok = true
	}
	return
}

// ToString will try to convert object o to string value.
func ToString(o Object) (v string, ok bool) {
	if o == UndefinedValue {
		return
	}
	ok = true
	if str, isStr := o.(*String); isStr {
		v = str.Value
	} else {
		v = o.String()
	}
	return
}
