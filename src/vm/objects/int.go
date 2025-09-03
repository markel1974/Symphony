package objects

import (
	"encoding/gob"
	"strconv"
)

const (
	IntType = "int"
)

func init() {
	gob.Register(&Int{})
}

// Int represents an integer type with a 64-bit value and methods for operations, equality, and object behavior.
type Int struct {
	gk    IGateKeeper
	frame int
	value int64
}

// NewInt creates and returns a new instance of the Int struct initialized with the specified int64 value.
func newInt(factory IGateKeeper, frame int, value int64) IObject {
	return &Int{
		gk:    factory,
		frame: frame,
		value: value,
	}
}

// GateKeeper returns a reference to the GateKeeper associated with the Object.
func (o *Int) GateKeeper() IGateKeeper {
	return o.gk
}

// AsBool converts the integer value to a boolean, returning true if the value is non-zero, otherwise false.
func (o *Int) AsBool() bool {
	return o.value != 0
}

// AsInt64 returns the length of the array as an int64 value.
func (o *Int) AsInt64() int64 {
	return o.value
}

// AsFloat64 returns the length of the array as an int64 value.
func (o *Int) AsFloat64() float64 {
	return float64(o.value)
}

// AsString returns the string representation of the Int value using base 10 format.
func (o *Int) AsString() string {
	return strconv.FormatInt(o.value, 10)
}

// AssignValue assigns the value of another IObject to the current Int object if the type is compatible, otherwise returns an error.
func (o *Int) AssignValue(v IObject) error {
	target, ok := v.(*Int)
	if !ok {
		return ErrNotAssignable
	}
	o.value = target.value
	return nil
}

// Nil checks if the object is nil and always returns false.
func (o *Int) Nil() bool {
	return false
}

// SetStatic sets the frame to FrameStatic, marking it with a static execution context.
func (o *Int) SetStatic() {
	o.frame = FrameStatic
}

// Frame returns the current frame value of the Object.
func (o *Int) Frame() int {
	return o.frame
}

// IndexGet attempts to retrieve a value at the given index and returns an error if the object is not indexable.
func (o *Int) IndexGet(_ int, _ IObject) (res IObject, err error) {
	return nil, ErrNotIndexable
}

// IndexSet attempts to assign a value to an index in the object but always returns ErrUnsupportedIndex,
// as this operation is unsupported.
func (o *Int) IndexSet(_, _ IObject) (err error) {
	return ErrUnsupportedIndex
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (o *Int) Iterate(_ int) IIterator {
	return nil
}

// CanIterate determines if the object can be iterated over and returns false for this implementation.
func (o *Int) CanIterate() bool {
	return false
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *Int) Call(_ int, _ ...IObject) (ret IObject, err error) {
	return nil, nil
}

// CanCall determines if the object can be invoked as a callable. Returns false for non-callable objects.
func (o *Int) CanCall() bool {
	return false
}

// Length returns the length of the Int object.
func (o *Int) Length() int {
	return 0
}

// Value returns the underlying int64 value of the Int object.
func (o *Int) Value() int64 {
	return o.value
}

// SetValue sets the underlying int64 value of the Int object.
func (o *Int) SetValue(value int64) {
	o.value = value
}

// TypeName returns the name of the type as a string, which is "int" for this object.
func (o *Int) TypeName() string {
	return IntType
}

// LogicalOp performs a logical operation between the object's value and a given operand, using the specified operator.
func (o *Int) LogicalOp(_ int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		return logicalOpNil(o.gk, op)
	}
	ret, err := logicalOpInt64(o.value, op, rhsIn.AsInt64())
	if err != nil {
		return nil, err
	}
	if ret {
		return o.gk.TrueValue(), nil
	}
	return o.gk.FalseValue(), nil
}

// ArithmeticOp performs a binary arithmetic operation on the Int object using the specified operator and right-hand operand.
func (o *Int) ArithmeticOp(frame int, op ArithmeticOperator, rhsIn IObject) (IObject, error) {
	ret, err := arithmeticOpInt64(o.value, op, rhsIn.AsInt64())
	if err != nil {
		return nil, err
	}
	if ret == o.value {
		return o, nil
	}
	return o.GateKeeper().NewInt(frame, ret), nil
}

// Copy creates and returns a new instance of the Int object with the same value as the current instance.
func (o *Int) Copy(frame int, _ int) IObject {
	return o.GateKeeper().NewInt(frame, o.value)
}

// Falsy checks whether the integer value is considered falsy. Returns true if the value is 0, otherwise false.
func (o *Int) Falsy() bool {
	return o.value == 0
}

// Equals checks if the current Int object is equal to another IObject of type Int by comparing their values.
func (o *Int) Equals(x IObject) bool {
	t, ok := x.(*Int)
	if !ok {
		return false
	}
	return o.value == t.value
}
