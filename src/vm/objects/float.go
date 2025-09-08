package objects

import (
	"encoding/gob"
	"math"
	"strconv"
)

const (
	FloatType = "float"
)

func init() {
	gob.Register(&Float{})
}

// Float represents a floating-point number and provides operations and behaviors specific to numeric types.
// It embeds Object to implement common interface methods and extends behavior where necessary.
// The value field holds the actual float64 values encapsulated by the Float type.
type Float struct {
	Allocator
	value float64
}

// NewFloat creates and returns a pointer to a new Float object initialized with the specified float64 values.
func newFloat(gk IGateKeeper, frame int, value float64) IObject {
	return &Float{
		Allocator: Allocator{gk: gk, frame: frame},
		value:     value,
	}
}

// AsBool returns true if the object is not empty, otherwise false.
func (o *Float) AsBool() bool {
	return o.value != 0
}

// AsInt64 returns the length of the array as an int64 value.
func (o *Float) AsInt64() int64 {
	return int64(o.value)
}

// AsFloat64 returns the length of the array as an int64 value.
func (o *Float) AsFloat64() float64 {
	return o.value
}

// AssignValue assigns the value of another IObject to the current Float object if the type is compatible, otherwise returns an error.
func (o *Float) AssignValue(v IObject) error {
	target, ok := v.(*Float)
	if !ok {
		return ErrNotAssignable
	}
	o.value = target.value
	return nil
}

// AsString returns the string representation of the Float object using its internal float64 values.
func (o *Float) AsString() string {
	return strconv.FormatFloat(o.value, 'f', -1, 64)
}

// Nil checks if the object is nil and always returns false.
func (o *Float) Nil() bool {
	return false
}

// IndexGet attempts to retrieve a value at the given index and returns an error if the object is not indexable.
func (o *Float) IndexGet(_ int, _ IObject) (IObject, error) {
	return o.gk.UndefinedValue(), ErrNotIndexable
}

// IndexSet attempts to assign a value to an index in the object but always returns ErrUnsupportedIndex,
// as this operation is unsupported.
func (o *Float) IndexSet(_, _ IObject) error {
	return ErrUnsupportedIndex
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (o *Float) Iterate(_ int) IIterator {
	return nil
}

// Iterable determines if the object can be iterated over and returns false for this implementation.
func (o *Float) Iterable() bool {
	return false
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *Float) Call(_ int, _ ...IObject) (retCount uint, ret IObject, err error) {
	return 0, nil, nil
}

// Length returns the length of the Int object.
func (o *Float) Length() int {
	return 0
}

func (o *Float) Value() float64 {
	return o.value
}

// TypeName returns the name of the type.
func (o *Float) TypeName() string {
	return FloatType
}

// LogicalOp performs a logical operation between the Float object and another IObject using the specified operator.
// Returns a boolean IObject representing the result or an error if the operation is invalid.
func (o *Float) LogicalOp(_ int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		return logicalOpNil(o.gk, op)
	}
	ret, err := logicalOpFloat64(o.value, op, rhsIn.AsFloat64())
	if err != nil {
		return nil, err
	}
	if ret {
		return o.gk.TrueValue(), nil
	}
	return o.gk.FalseValue(), nil
}

// ArithmeticOp performs an arithmetic operation with the specified operator and returns the result or an error.
func (o *Float) ArithmeticOp(frame int, op ArithmeticOperator, rhsIn IObject) (IObject, error) {
	ret, err := arithmeticOpFloat64(o.value, op, rhsIn.AsFloat64())
	if err != nil {
		return nil, err
	}
	if ret == o.value {
		return o, nil
	}
	return o.GateKeeper().NewFloat(frame, ret), nil
}

// Copy creates and returns a new instance of the Float object, duplicating its current state.
func (o *Float) Copy(frame int, _ int) IObject {
	return o.GateKeeper().NewFloat(frame, o.value)
}

// Falsy determines if the float object is considered falsy, returning true if the values is NaN; otherwise, false.
func (o *Float) Falsy() bool {
	return math.IsNaN(o.value)
}

// Equals checks if the current Float object is equal to another IObject by comparing their internal float64 values.
func (o *Float) Equals(x IObject) bool {
	t, ok := x.(*Float)
	if !ok {
		return false
	}
	return o.value == t.value
}
