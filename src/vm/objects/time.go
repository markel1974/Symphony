package objects

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"time"
)

const (
	TimeType = "time"
)

func init() {
	gob.Register(&Time{})
}

// Time represents a custom object encapsulating a Go time.Time Code with extended behaviors and operations.
type Time struct {
	IAllocator
	data time.Time
}

// NewTime creates a new instance of Time wrapping the provided time.Time Code.
func newTime(allocator IAllocator, value time.Time) IObject {
	return &Time{
		IAllocator: allocator,
		data:       value,
	}
}

// setAllocator sets the allocator for the instance, defining its memory management and lifecycle behavior.
func (o *Time) setAllocator(allocator IAllocator) {
	o.IAllocator = allocator
}

// AsInterface converts the object into a generic interface{} type and returns the underlying data.
func (o *Time) AsInterface() interface{} {
	return o.data
}

// AsBool returns the boolean representation of the Time object, which is true if the Code is not zero.
func (o *Time) AsBool() bool {
	return !o.data.IsZero()
}

// AsInt64 returns the len of the array as an int64 Code.
func (o *Time) AsInt64() int64 {
	return o.data.Unix()
}

// AsFloat64 returns the len of the array as an int64 Code.
func (o *Time) AsFloat64() float64 {
	return float64(o.data.Unix())
}

// AsString returns the string representation of the Time object by delegating to the underlying time.Time Code.
func (o *Time) AsString() string {
	return o.data.String()
}

// AsBytes converts the object elements into a single concatenated slice of bytes by calling AsBytes on each element.
func (o *Time) AsBytes() []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(o.data.Unix()))
	return b
}

// AssignValue assigns the Code of another IObject to the current Time object if the type is compatible, otherwise returns an error.
func (o *Time) AssignValue(v IObject) error {
	target, ok := v.(*Time)
	if !ok {
		return ErrNotAssignable
	}
	o.data = target.data
	return nil
}

// Nil checks if the object is nil and always returns false.
func (o *Time) Nil() bool {
	return false
}

// IndexGet attempts to retrieve a Code at the given index and returns an error if the object is not indexable.
func (o *Time) IndexGet(_ int, _ IObject) (IObject, error) {
	return o.GateKeeper().UndefinedValue(), ErrIndexNotIndexable
}

// IndexSet attempts to assign a Code to an index in the object but always returns ErrIndexUnsupported,
// as this operation is unsupported.
func (o *Time) IndexSet(_, _ IObject) error {
	return ErrIndexUnsupported
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (o *Time) Iterate(_ int) IIterator {
	return nil
}

// Iterable determines if the object can be iterated over and returns false for this implementation.
func (o *Time) Iterable() bool {
	return false
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *Time) Call(_ int, _ ...IObject) (retCount uint, ret IObject, err error) {
	return 0, nil, nil
}

// Length returns the len of the Int object.
func (o *Time) Length() int {
	return 0
}

// Value returns the underlying time.Time Code of the Time object.
func (o *Time) Value() time.Time {
	return o.data
}

// TypeName returns the name of the type as a string, which is "time".
func (o *Time) TypeName() string {
	return TimeType
}

// LogicalOp performs logical comparison operations (e.g., <, >, <=, >=) between the Time object and another Time object.
func (o *Time) LogicalOp(_ int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		return logicalOpNil(o.GateKeeper(), op)
	}
	switch rhs := rhsIn.(type) {
	case *Time:
		switch op {
		case OperatorLogicalLess:
			if o.data.Before(rhs.data) {
				return o.GateKeeper().TrueValue(), nil
			}
			return o.GateKeeper().FalseValue(), nil
		case OperatorLogicalGreater:
			if o.data.After(rhs.data) {
				return o.GateKeeper().TrueValue(), nil
			}
			return o.GateKeeper().FalseValue(), nil
		case OperatorLogicalLessEq:
			if o.data.Equal(rhs.data) || o.data.Before(rhs.data) {
				return o.GateKeeper().TrueValue(), nil
			}
			return o.GateKeeper().FalseValue(), nil
		case OperatorLogicalGreaterEq:
			if o.data.Equal(rhs.data) || o.data.After(rhs.data) {
				return o.GateKeeper().TrueValue(), nil
			}
			return o.GateKeeper().FalseValue(), nil
		default:
			return nil, ErrInvalidOperator
		}
	}
	return nil, ErrInvalidOperator
}

// ArithmeticOp performs arithmetic operations (addition, subtraction) on the Time object based on the given operator.
// Supports operations with Int and Time objects. Returns a new IObject or an error for invalid operators.
func (o *Time) ArithmeticOp(frame int, op ArithmeticOperator, in IObject) (IObject, error) {
	switch rhs := in.(type) {
	case *Int:
		switch op {
		case OperatorAdd:
			if rhs.data == 0 {
				return o, nil
			}
			return o.GateKeeper().NewTime(frame, o.data.Add(time.Duration(rhs.data))), nil
		case OperatorSub:
			if rhs.data == 0 {
				return o, nil
			}
			return o.GateKeeper().NewTime(frame, o.data.Add(time.Duration(-rhs.data))), nil
		default:
			return nil, ErrInvalidOperator
		}
	case *Time:
		switch op {
		case OperatorSub:
			return o.GateKeeper().NewInt(frame, int64(o.data.Sub(rhs.data))), nil
		default:
			return nil, ErrInvalidOperator
		}
	}
	return nil, ErrInvalidOperator
}

// Copy returns a new instance of the Time object with the same internal time Code, duplicating its state.
func (o *Time) Copy(frame int, _ int) IObject {
	return o.GateKeeper().NewTime(frame, o.data)
}

// Falsy returns true if the Time object's Code is zero (indicating it is uninitialized or empty), otherwise false.
func (o *Time) Falsy() bool {
	return o.data.IsZero()
}

// Equals checks whether the Time object is equal to another object of type IObject, returning true if they match.
func (o *Time) Equals(x IObject) bool {
	t, ok := x.(*Time)
	if !ok {
		return false
	}
	return o.data.Equal(t.data)
}

// Count returns the total number of elements in the instance and its sub-elements.
func (o *Time) Count() int {
	return 1
}

// SetValue sets the internal time.Time value of the Time object to the specified value.
func (o *Time) SetValue(value time.Time) {
	o.data = value
}

// GobEncode serializes the Time's data into a byte slice using gob encoding and returns the result or an error.
func (o *Time) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(o.data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GobDecode decodes the provided byte slice into the Time's data field using the gob package.
func (o *Time) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	if err := decoder.Decode(&o.data); err != nil {
		return err
	}
	return nil
}
