package objects

import (
	"bytes"
	"encoding/gob"
	"reflect"
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
	data string
}

// NewString creates and returns a new String object initialized with the provided string Code.
func newString(allocator IAllocator, value string) IObject {
	if len(value) > MaxStringLen {
		value = value[0:MaxStringLen]
	}
	return &String{
		IAllocator: allocator,
		data:       value,
	}
}

// setAllocator sets the allocator for the instance, defining its memory management and lifecycle behavior.
func (o *String) setAllocator(allocator IAllocator) {
	o.IAllocator = allocator
}

// AsInterface converts the object into a generic interface{} type and returns the underlying data.
func (o *String) AsInterface() interface{} {
	return o.data
}

// AsValue attempts to convert the current object to a reflect.Value of the specified target type. Returns success status.
func (o *String) AsValue(target reflect.Type) (reflect.Value, bool) {
	return _reflect(o, target)
}

// AsBool converts the String object to a boolean. Returns true if the string has non-zero len, otherwise false.
func (o *String) AsBool() bool {
	return len(o.data) > 0
}

// AsInt64 returns the len of the array as an int64 data.
func (o *String) AsInt64() int64 {
	return int64(len(o.data))
}

// AsFloat64 returns the len of the array as an int64 data.
func (o *String) AsFloat64() float64 {
	return float64(len(o.data))
}

// AsBytes converts the object elements into a single concatenated slice of bytes by calling AsBytes on each element.
func (o *String) AsBytes() []byte {
	return []byte(o.data)
}

// AsString returns the quoted string representation of the String object.
func (o *String) AsString() string {
	return o.data
}

// Nil checks if the object is nil and always returns false.
func (o *String) Nil() bool {
	return false
}

// AssignValue assigns the data of another IObject to the current String object if the type is compatible, otherwise returns an error.
func (o *String) AssignValue(v IObject) error {
	target, ok := v.(*String)
	if !ok {
		return ErrNotAssignable
	}
	if len(target.data) > MaxStringLen {
		return ErrLimitExceed
	}
	o.data = target.data
	return nil
}

// IndexSet attempts to assign a data to an index in the object but always returns ErrIndexUnsupported,
// as this operation is unsupported.
func (o *String) IndexSet(_, _ IObject) (err error) {
	return ErrIndexUnsupported
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *String) Call(_ int, _ ...IObject) (retCount uint, ret IObject, err error) {
	return 0, nil, nil
}

// Value returns the string data of the String object.
func (o *String) Value() string {
	return o.data
}

// Length returns the len of the string data.
func (o *String) Length() int {
	return len(o.data)
}

// TypeName returns the name of the type "string".
func (o *String) TypeName() string {
	return StringType
}

// LogicalOp performs logical comparison between the String object and a right-hand side IObject based on the given operator.
// Returns a boolean IObject or an error for unsupported operations or invalid types.
func (o *String) LogicalOp(_ int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		ret, err := logicalOpNil(op)
		return o.GateKeeper().FromBoolError(ret, err)
	}

	switch op {
	case OperatorLogicalNotEq:
		if o.data != rhsIn.AsString() {
			return o.GateKeeper().TrueValue(), nil
		}
		return o.GateKeeper().FalseValue(), nil
	case OperatorLogicalEq:
		if o.data == rhsIn.AsString() {
			return o.GateKeeper().TrueValue(), nil
		}
		return o.GateKeeper().FalseValue(), nil
	case OperatorLogicalLess:
		if o.data < rhsIn.AsString() {
			return o.GateKeeper().TrueValue(), nil
		}
		return o.GateKeeper().FalseValue(), nil
	case OperatorLogicalLessEq:
		if o.data <= rhsIn.AsString() {
			return o.GateKeeper().TrueValue(), nil
		}
		return o.GateKeeper().FalseValue(), nil
	case OperatorLogicalGreater:
		if o.data > rhsIn.AsString() {
			return o.GateKeeper().TrueValue(), nil
		}
		return o.GateKeeper().FalseValue(), nil
	case OperatorLogicalGreaterEq:
		if o.data >= rhsIn.AsString() {
			return o.GateKeeper().TrueValue(), nil
		}
		return o.GateKeeper().FalseValue(), nil
	default:
		return nil, ErrInvalidOperator
	}
}

// ArithmeticOp performs an arithmetic operation based on the given operator and right-hand side operand.
// Supports concatenation with strings if the operator is OperatorAdd.
// Returns a new resulting object or an error for unsupported operations or excessive string len.
func (o *String) ArithmeticOp(frame int, op ArithmeticOperator, rhsIn IObject) (IObject, error) {
	switch op {
	case OperatorAdd:
		str := rhsIn.AsString()
		if len(o.data)+len(str) > MaxStringLen {
			return nil, ErrLimitExceed
		}
		return o.GateKeeper().NewString(frame, o.data+str), nil
	default:
		return nil, ErrInvalidOperator
	}
}

// UnaryOp performs a unary operation on the String object and always returns ErrInvalidOperator.
func (o *String) UnaryOp(_ int, _ UnaryOperator) (IObject, error) {
	return nil, ErrInvalidOperator
}

// Falsy returns true if the String's data is an empty string, indicating it is considered falsy in a boolean context.
func (o *String) Falsy() bool {
	return len(o.data) == 0
}

// Copy creates and returns a new String instance with the same data as the original.
func (o *String) Copy(frame int, _ int) IObject {
	return o.GateKeeper().NewString(frame, o.data)
}

// Equals checks whether the current String object is equal to the provided IObject by comparing their data.
func (o *String) Equals(x IObject) bool {
	t, ok := x.(*String)
	if !ok {
		return false
	}
	return o.data == t.data
}

// IndexGet retrieves the character at the specified index from the String object.
// Returns an error if the index is not of type Int or is out of bounds.
func (o *String) IndexGet(frame int, index IObject) (IObject, error) {
	intIdx, ok := index.(*Int)
	if !ok {
		return nil, ErrIndexInvalidType
	}
	idxVal := int(intIdx.data)
	r := []rune(o.data)
	if idxVal < 0 || idxVal >= len(r) {
		return o.GateKeeper().UndefinedValue(), nil
	}
	return o.GateKeeper().NewChar(frame, r[idxVal]), nil
}

// Count returns the total number of elements in the instance and its sub-elements.
func (o *String) Count() int {
	return 1
}

// SetValue assigns the given string data to the object, returning an error if the data exceeds the maximum allowed len.
func (o *String) SetValue(v string) error {
	if len(v) > MaxStringLen {
		return ErrLimitExceed
	}
	o.data = v
	return nil
}

// Iterable checks if the String object supports iteration and always returns true.
func (o *String) Iterable() bool {
	return true
}

// Iterate returns an IIterator for iterating over the runes of the String. It initializes runeStr if not already initialized.
func (o *String) Iterate(frame int) IIterator {
	return o.GateKeeper().NewStringIterator(frame, []rune(o.data), 0)
}

// GobEncode serializes the String's data into a byte slice using gob encoding and returns the result or an error.
func (o *String) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(o.data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GobDecode decodes the provided byte slice into the String's data field using the gob package.
func (o *String) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	if err := decoder.Decode(&o.data); err != nil {
		return err
	}
	return nil
}
