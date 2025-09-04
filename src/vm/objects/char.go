package objects

import "encoding/gob"

const (
	CharType = "char"
)

func init() {
	gob.Register(&Char{})
}

// Char represents a character type, encapsulating a single rune values and inheriting behavior from Object.
type Char struct {
	Allocator
	value rune
}

// NewChar creates and returns a new Char object with the specified rune values.
func newChar(factory IGateKeeper, frame int, value rune) IObject {
	return &Char{
		Allocator: Allocator{gk: factory, frame: frame},
		value:     value,
	}
}

// AsBool returns true if the object is not empty, otherwise false.
func (o *Char) AsBool() bool {
	return o.value != 0
}

// AsInt64 returns the length of the array as an int64 value.
func (o *Char) AsInt64() int64 {
	return int64(o.value)
}

// AsFloat64 returns the length of the array as an int64 value.
func (o *Char) AsFloat64() float64 {
	return float64(o.value)
}

// AsString returns the string representation of the Char object's values.
func (o *Char) AsString() string {
	return string(o.value)
}

// AssignValue assigns the value of another IObject to the current Char object if the type is compatible, otherwise returns an error.
func (o *Char) AssignValue(v IObject) error {
	target, ok := v.(*Char)
	if !ok {
		return ErrNotAssignable
	}
	o.value = target.value
	return nil
}

// Nil checks if the object is nil and always returns false.
func (o *Char) Nil() bool {
	return false
}

// IndexGet attempts to retrieve a value at the given index and returns an error if the object is not indexable.
func (o *Char) IndexGet(_ int, _ IObject) (res IObject, err error) {
	return nil, ErrNotIndexable
}

// IndexSet attempts to assign a value to an index in the object but always returns ErrUnsupportedIndex,
// as this operation is unsupported.
func (o *Char) IndexSet(_, _ IObject) (err error) {
	return ErrUnsupportedIndex
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (o *Char) Iterate(_ int) IIterator {
	return nil
}

// Iterable determines if the object can be iterated over and returns false for this implementation.
func (o *Char) Iterable() bool {
	return false
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *Char) Call(_ int, _ ...IObject) (retCount uint, ret IObject, err error) {
	return 0, nil, nil
}

// Length returns the length of the Int object.
func (o *Char) Length() int {
	return 0
}

// Value returns the rune values stored in the Char object.
func (o *Char) Value() rune {
	return o.value
}

// TypeName returns the name of the type as a string.
func (o *Char) TypeName() string {
	return CharType
}

// LogicalOp performs a logical operation between the current Char object and another IObject using the specified operator.
func (o *Char) LogicalOp(_ int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		return logicalOpNil(o.gk, op)
	}
	ret, err := logicalOpInt64(int64(o.value), op, rhsIn.AsInt64())
	if err != nil {
		return nil, err
	}
	if ret {
		return o.gk.TrueValue(), nil
	}
	return o.gk.FalseValue(), nil
}

// ArithmeticOp applies the specified arithmetic operation between a Char object and another IObject, returning the result.
// Returns an error if the operation is invalid or unsupported.
func (o *Char) ArithmeticOp(frame int, op ArithmeticOperator, rhsIn IObject) (IObject, error) {
	ret, err := arithmeticOpInt64(int64(o.value), op, rhsIn.AsInt64())
	if err != nil {
		return nil, err
	}
	if ret == int64(o.value) {
		return o, nil
	}
	return o.GateKeeper().NewChar(frame, rune(ret)), nil
}

// Copy creates and returns a new instance of the Char object with the same values.
func (o *Char) Copy(frame int, _ int) IObject {
	return o.GateKeeper().NewChar(frame, o.value)
}

// Falsy checks whether the Char object represents a falsy state, returning true if the underlying values is 0.
func (o *Char) Falsy() bool {
	return o.value == 0
}

// Equals checks if the current Char object is equal to another IObject. Returns true if both objects are equal.
func (o *Char) Equals(x IObject) bool {
	t, ok := x.(*Char)
	if !ok {
		return false
	}
	return o.value == t.value
}
