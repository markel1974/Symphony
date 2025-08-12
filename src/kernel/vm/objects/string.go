package objects

import (
	"strconv"
)

// MaxStringLen defines the maximum allowed byte-length for string values across all compiler/VM instances in the process.
var (
	// MaxStringLen is the maximum byte-length for string values. Note this
	// limit applies to all compiler/VM instances in the process.
	MaxStringLen = 2147483647
)

// String represents a wrapper around a standard string with additional behavior and methods for runtime operations.
// This type embeds ObjectImpl and supports operations like indexing, iteration, comparison, and copying.
// It implements IObject and provides a richer functionality for string manipulation within the runtime system.
type String struct {
	ObjectImpl
	value   string
	runeStr []rune
}

// NewString creates and returns a new String object initialized with the provided string values.
func NewString(value string) (*String, error) {
	if len(value) > MaxStringLen {
		return nil, ErrStringLimit
	}
	return &String{
		value: value,
	}, nil
}

func NewStringNoSize(value string) *String {
	return &String{
		value: value,
	}
}

// Value returns the string values of the String object.
func (o *String) Value() string {
	return o.value
}

// Length returns the length of the string values.
func (o *String) Length() int {
	return len(o.value)
}

// TypeName returns the name of the type "string".
func (o *String) TypeName() string {
	return "string"
}

// String returns the quoted string representation of the String object.
func (o *String) String() string {
	return strconv.Quote(o.value)
}

// BinaryOp performs the specified binary operation on the calling string object and a right-hand operand.
// Supported operations include addition and comparison (e.g., less than, greater than, and their equal variants).
// Returns the result of the operation or an error if the operation is invalid or exceeds size limits.
func (o *String) BinaryOp(op Operator, rhs IObject) (IObject, error) {
	switch op {
	case OperatorAdd:
		switch rhs := rhs.(type) {
		case *String:
			if len(o.value)+len(rhs.value) > MaxStringLen {
				return nil, ErrStringLimit
			}
			return &String{value: o.value + rhs.value}, nil
		default:
			rhsStr := rhs.String()
			if len(o.value)+len(rhsStr) > MaxStringLen {
				return nil, ErrStringLimit
			}
			return &String{value: o.value + rhsStr}, nil
		}
	case OperatorLess:
		switch rhs := rhs.(type) {
		case *String:
			if o.value < rhs.value {
				return TrueValue, nil
			}
			return FalseValue, nil
		}
	case OperatorLessEq:
		switch rhs := rhs.(type) {
		case *String:
			if o.value <= rhs.value {
				return TrueValue, nil
			}
			return FalseValue, nil
		}
	case OperatorGreater:
		switch rhs := rhs.(type) {
		case *String:
			if o.value > rhs.value {
				return TrueValue, nil
			}
			return FalseValue, nil
		}
	case OperatorGreaterEq:
		switch rhs := rhs.(type) {
		case *String:
			if o.value >= rhs.value {
				return TrueValue, nil
			}
			return FalseValue, nil
		}
	default:
		return nil, ErrInvalidOperator
	}
	return nil, ErrInvalidOperator
}

// Falsy returns true if the String's values is an empty string, indicating it is considered falsy in a boolean context.
func (o *String) Falsy() bool {
	return len(o.value) == 0
}

// Copy creates and returns a new String instance with the same values as the original.
func (o *String) Copy() IObject {
	return &String{value: o.value}
}

// Equals checks whether the current String object is equal to the provided IObject by comparing their values.
func (o *String) Equals(x IObject) bool {
	t, ok := x.(*String)
	if !ok {
		return false
	}
	return o.value == t.value
}

// IndexGet retrieves the character at the specified index from the String object.
// Returns an error if the index is not of type Int or is out of bounds.
func (o *String) IndexGet(index IObject) (res IObject, err error) {
	intIdx, ok := index.(*Int)
	if !ok {
		err = ErrInvalidIndexType
		return
	}
	idxVal := int(intIdx.value)
	if o.runeStr == nil {
		o.runeStr = []rune(o.value)
	}
	if idxVal < 0 || idxVal >= len(o.runeStr) {
		res = UndefinedValue
		return
	}
	res = &Char{value: o.runeStr[idxVal]}
	return
}

// Iterate returns an IIterator for iterating over the runes of the String. It initializes runeStr if not already initialized.
func (o *String) Iterate() IIterator {
	if o.runeStr == nil {
		o.runeStr = []rune(o.value)
	}
	return &StringIterator{
		v: o.runeStr,
		l: len(o.runeStr),
	}
}

// CanIterate checks if the String object supports iteration and always returns true.
func (o *String) CanIterate() bool {
	return true
}

// StringIterator represents an iterator for traversing over the characters of a string, implemented as runes.
type StringIterator struct {
	ObjectImpl
	v []rune
	i int
	l int
}

// TypeName returns the type name of the StringIterator as a string.
func (i *StringIterator) TypeName() string {
	return "string-iterator"
}

// String returns the string representation of the StringIterator, useful for debugging or logging.
func (i *StringIterator) String() string {
	return "<string-iterator>"
}

// Falsy returns true, indicating the iterator is considered falsy in a boolean context.
func (i *StringIterator) Falsy() bool {
	return true
}

// Equals compares the current StringIterator with another IObject and determines if they are equal.
func (i *StringIterator) Equals(IObject) bool {
	return false
}

// Copy creates and returns a new instance of StringIterator with the same state as the current one.
func (i *StringIterator) Copy() IObject {
	return &StringIterator{v: i.v, i: i.i, l: i.l}
}

// Next advances the iterator to the next position and returns true if the current position is within bounds.
func (i *StringIterator) Next() bool {
	i.i++
	return i.i <= i.l
}

// Key returns the zero-based index of the current element in the iteration as an Int object.
func (i *StringIterator) Key() IObject {
	return &Int{value: int64(i.i - 1)}
}

// Value returns the current character as an IObject wrapped in a Char instance from the internal rune slice.
func (i *StringIterator) Value() IObject {
	return &Char{value: i.v[i.i-1]}
}

// ToRune attempts to convert an IObject to a rune if it is of type *Int or *Char, returning the rune and a boolean success flag.
func ToRune(o IObject) (v rune, ok bool) {
	switch o := o.(type) {
	case *Int:
		v = rune(o.value)
		ok = true
	case *Char:
		v = o.value
		ok = true
	}
	return
}

// ToString converts an IObject to its string representation and determines whether the conversion is valid.
func ToString(o IObject) (string, bool) {
	if o == nil {
		return "", false
	}
	if o == UndefinedValue {
		return "", false
	}
	if str, isStr := o.(*String); isStr {
		return str.value, true
	}
	return o.String(), true
}

func ToStringArg(name string, o IObject) (string, error) {
	v, ok := ToString(o)
	if !ok {
		return "", NewInvalidArgumentType(name, "string(compatible)", o.TypeName())
	}
	return v, nil
}
