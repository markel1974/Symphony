package objects

import (
	"encoding/gob"
)

const (
	StringType = "string"
)

func init() {
	gob.Register(&String{})
}

// String represents a wrapper around a standard string with additional behavior and methods for runtime operations.
// This type embeds Object and supports operations like indexing, iteration, comparison, and copying.
// It implements IObject and provides a richer functionality for string manipulation within the runtime system.
type String struct {
	IAllocator
	value   string
	runeStr []rune
}

// NewString creates and returns a new String object initialized with the provided string values.
func newString(allocator IAllocator, value string) IObject {
	if len(value) > MaxStringLen {
		value = value[0:MaxStringLen]
	}
	return &String{
		IAllocator: allocator,
		value:      value,
	}
}

// setAllocator sets the allocator for the instance, defining its memory management and lifecycle behavior.
func (o *String) setAllocator(allocator IAllocator) {
	o.IAllocator = allocator
}

// AsBool converts the String object to a boolean. Returns true if the string has non-zero length, otherwise false.
func (o *String) AsBool() bool {
	return len(o.value) > 0
}

// AsInt64 returns the length of the array as an int64 value.
func (o *String) AsInt64() int64 {
	return int64(len(o.value))
}

// AsFloat64 returns the length of the array as an int64 value.
func (o *String) AsFloat64() float64 {
	return float64(len(o.value))
}

// AsString returns the quoted string representation of the String object.
func (o *String) AsString() string {
	return o.value
}

// Nil checks if the object is nil and always returns false.
func (o *String) Nil() bool {
	return false
}

// AssignValue assigns the value of another IObject to the current String object if the type is compatible, otherwise returns an error.
func (o *String) AssignValue(v IObject) error {
	target, ok := v.(*String)
	if !ok {
		return ErrNotAssignable
	}
	if len(target.value) > MaxStringLen {
		return ErrLimitExceed
	}
	o.value = target.value
	return nil
}

// IndexSet attempts to assign a value to an index in the object but always returns ErrIndexUnsupported,
// as this operation is unsupported.
func (o *String) IndexSet(_, _ IObject) (err error) {
	return ErrIndexUnsupported
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *String) Call(_ int, _ ...IObject) (retCount uint, ret IObject, err error) {
	return 0, nil, nil
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

// LogicalOp performs logical comparison between the String object and a right-hand side IObject based on the given operator.
// Returns a boolean IObject or an error for unsupported operations or invalid types.
func (o *String) LogicalOp(_ int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		return logicalOpNil(o.GateKeeper(), op)
	}

	switch op {
	case OperatorLogicalNotEq:
		if o.value != rhsIn.AsString() {
			return o.GateKeeper().TrueValue(), nil
		}
		return o.GateKeeper().FalseValue(), nil
	case OperatorLogicalEq:
		if o.value == rhsIn.AsString() {
			return o.GateKeeper().TrueValue(), nil
		}
		return o.GateKeeper().FalseValue(), nil
	case OperatorLogicalLess:
		if o.value < rhsIn.AsString() {
			return o.GateKeeper().TrueValue(), nil
		}
		return o.GateKeeper().FalseValue(), nil
	case OperatorLogicalLessEq:
		if o.value <= rhsIn.AsString() {
			return o.GateKeeper().TrueValue(), nil
		}
		return o.GateKeeper().FalseValue(), nil
	case OperatorLogicalGreater:
		if o.value > rhsIn.AsString() {
			return o.GateKeeper().TrueValue(), nil
		}
		return o.GateKeeper().FalseValue(), nil
	case OperatorLogicalGreaterEq:
		if o.value >= rhsIn.AsString() {
			return o.GateKeeper().TrueValue(), nil
		}
		return o.GateKeeper().FalseValue(), nil
	default:
		return nil, ErrInvalidOperator
	}
}

// ArithmeticOp performs an arithmetic operation based on the given operator and right-hand side operand.
// Supports concatenation with strings if the operator is OperatorAdd.
// Returns a new resulting object or an error for unsupported operations or excessive string length.
func (o *String) ArithmeticOp(frame int, op ArithmeticOperator, rhsIn IObject) (IObject, error) {
	switch op {
	case OperatorAdd:
		str := rhsIn.AsString()
		if len(o.value)+len(str) > MaxStringLen {
			return nil, ErrLimitExceed
		}
		return o.GateKeeper().NewString(frame, o.value+str), nil
	default:
		return nil, ErrInvalidOperator
	}
}

// Falsy returns true if the String's values is an empty string, indicating it is considered falsy in a boolean context.
func (o *String) Falsy() bool {
	return len(o.value) == 0
}

// Copy creates and returns a new String instance with the same values as the original.
func (o *String) Copy(frame int, _ int) IObject {
	return o.GateKeeper().NewString(frame, o.value)
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
		return nil, ErrIndexInvalidType
	}
	idxVal := int(intIdx.value)
	if o.runeStr == nil {
		o.runeStr = []rune(o.value)
	}
	if idxVal < 0 || idxVal >= len(o.runeStr) {
		return o.GateKeeper().UndefinedValue(), nil
	}
	return o.GateKeeper().NewChar(frame, o.runeStr[idxVal]), nil
}

// Count returns the total number of elements in the instance and its sub-elements.
func (o *String) Count() int {
	return 1
}

// SetValue assigns the given string value to the object, returning an error if the value exceeds the maximum allowed length.
func (o *String) SetValue(v string) error {
	if len(v) > MaxStringLen {
		return ErrLimitExceed
	}
	o.value = v
	return nil
}

// Iterable checks if the String object supports iteration and always returns true.
func (o *String) Iterable() bool {
	return true
}

// Iterate returns an IIterator for iterating over the runes of the String. It initializes runeStr if not already initialized.
func (o *String) Iterate(frame int) IIterator {
	if o.runeStr == nil {
		o.runeStr = []rune(o.value)
	}
	return o.GateKeeper().NewStringIterator(frame, o.runeStr, 0)
}
