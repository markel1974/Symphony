package objects

import (
	"strconv"
)

const (
	StringType = "string"
)

// MaxStringLen defines the maximum allowed byte-length for string values across all compiler/VM instances in the process.
const (
	// MaxStringLen is the maximum byte-length for string values. Note this
	// limit applies to all compiler/VM instances in the process.
	MaxStringLen = 2147483647
)

// String represents a wrapper around a standard string with additional behavior and methods for runtime operations.
// This type embeds Object and supports operations like indexing, iteration, comparison, and copying.
// It implements IObject and provides a richer functionality for string manipulation within the runtime system.
type String struct {
	*Object
	value   string
	runeStr []rune
}

// NewString creates and returns a new String object initialized with the provided string values.
func _newString(factory *Factory, frame int, value string) (*String, error) {
	if len(value) > MaxStringLen {
		return nil, ErrStringLimit
	}
	return &String{
		Object: factory.NewObject(frame),
		value:  value,
	}, nil
}

func _newStringNoSize(factory *Factory, frame int, value string) *String {
	return &String{
		Object: factory.NewObject(frame),
		value:  value,
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
	return StringType
}

// String returns the quoted string representation of the String object.
func (o *String) String() string {
	return strconv.Quote(o.value)
}

// BinaryOp performs the specified binary operation on the calling string object and a right-hand operand.
// Supported operations include addition and comparison (e.g., less than, greater than, and their equal variants).
// Returns the result of the operation or an error if the operation is invalid or exceeds size limits.
func (o *String) BinaryOp(frame int, op Operator, rhs IObject) (IObject, error) {
	switch op {
	case OperatorAdd:
		switch rhs := rhs.(type) {
		case *String:
			if len(o.value)+len(rhs.value) > MaxStringLen {
				return nil, ErrStringLimit
			}
			return o.Factory().NewString(frame, o.value+rhs.value)
		default:
			rhsStr := rhs.String()
			if len(o.value)+len(rhsStr) > MaxStringLen {
				return nil, ErrStringLimit
			}
			return o.Factory().NewString(frame, o.value+rhsStr)
		}
	case OperatorLess:
		switch rhs := rhs.(type) {
		case *String:
			if o.value < rhs.value {
				return o.Factory().TrueValue(), nil
			}
			return o.Factory().FalseValue(), nil
		}
	case OperatorLessEq:
		switch rhs := rhs.(type) {
		case *String:
			if o.value <= rhs.value {
				return o.Factory().TrueValue(), nil
			}
			return o.Factory().FalseValue(), nil
		}
	case OperatorGreater:
		switch rhs := rhs.(type) {
		case *String:
			if o.value > rhs.value {
				return o.Factory().TrueValue(), nil
			}
			return o.Factory().FalseValue(), nil
		}
	case OperatorGreaterEq:
		switch rhs := rhs.(type) {
		case *String:
			if o.value >= rhs.value {
				return o.Factory().TrueValue(), nil
			}
			return o.Factory().FalseValue(), nil
		}
	default:
		return nil, ErrInvalidOperator
	}
	return nil, ErrInvalidOperator
}

// Boolean returns true if the String's values is an empty string, indicating it is considered falsy in a boolean context.
func (o *String) Boolean() bool {
	return len(o.value) == 0
}

// Copy creates and returns a new String instance with the same values as the original.
func (o *String) Copy(frame int) IObject {
	return o.Factory().NewStringNoSize(frame, o.value)
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
func (o *String) IndexGet(frame int, index IObject) (IObject, error) {
	intIdx, ok := index.(*Int)
	if !ok {
		return nil, ErrInvalidIndexType
	}
	idxVal := int(intIdx.value)
	if o.runeStr == nil {
		o.runeStr = []rune(o.value)
	}
	if idxVal < 0 || idxVal >= len(o.runeStr) {
		return o.Factory().UndefinedValue(), nil
	}
	return o.Factory().NewChar(frame, o.runeStr[idxVal]), nil
}

// CanIterate checks if the String object supports iteration and always returns true.
func (o *String) CanIterate() bool {
	return true
}

// Iterate returns an IIterator for iterating over the runes of the String. It initializes runeStr if not already initialized.
func (o *String) Iterate(frame int) IIterator {
	if o.runeStr == nil {
		o.runeStr = []rune(o.value)
	}
	return o.Factory().NewStringIterator(frame, o.runeStr)
}
