package objects

import "encoding/gob"

const (
	UndefinedType  = "undefined"
	UndefinedLabel = "<" + UndefinedType + ">"
)

func init() {
	gob.Register(&Undefined{})
}

// Undefined represents an undefined values.
type Undefined struct {
	IAllocator
}

func newUndefined(allocator IAllocator) IObject {
	return &Undefined{
		IAllocator: allocator,
	}
}

// setAllocator sets the allocator for the instance, defining its memory management and lifecycle behavior.
func (o *Undefined) setAllocator(allocator IAllocator) {
	o.IAllocator = allocator
}

// AsBool returns the boolean representation of the Undefined object, which is always false.
func (o *Undefined) AsBool() bool {
	return false
}

// AsInt64 returns the length of the array as an int64 value.
func (o *Undefined) AsInt64() int64 {
	return 0
}

// AsFloat64 returns the length of the array as an int64 value.
func (o *Undefined) AsFloat64() float64 {
	return 0
}

// AssignValue sets the current object to the provided IObject, returning ErrNotAssignable if the operation is not supported.
func (o *Undefined) AssignValue(_ IObject) error {
	return ErrNotAssignable
}

// Nil checks if the object is nil and always returns false.
func (o *Undefined) Nil() bool {
	return true
}

// LogicalOp performs a logical operation using the specified operator and right-hand side operand. Always returns an error.
func (o *Undefined) LogicalOp(_ int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		switch op {
		case OperatorLogicalEq:
			return o.GateKeeper().TrueValue(), nil
		case OperatorLogicalNotEq:
			return o.GateKeeper().FalseValue(), nil
		default:
			return o.GateKeeper().UndefinedValue(), ErrInvalidOperator
		}
	} else {
		switch op {
		case OperatorLogicalEq:
			return o.GateKeeper().FalseValue(), nil
		case OperatorLogicalNotEq:
			return o.GateKeeper().TrueValue(), nil
		default:
			return o.GateKeeper().UndefinedValue(), ErrInvalidOperator
		}
	}
}

// ArithmeticOp performs an arithmetic operation on the object with the specified operator and operand, returning an error for invalid operators.
func (o *Undefined) ArithmeticOp(_ int, _ ArithmeticOperator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// IndexSet attempts to assign a value to an index in the object but always returns ErrIndexUnsupported,
// as this operation is unsupported.
func (o *Undefined) IndexSet(_, _ IObject) (err error) {
	return ErrIndexUnsupported
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *Undefined) Call(_ int, _ ...IObject) (retCount uint, ret IObject, err error) {
	return 0, nil, nil
}

// Length returns the length of the Int object.
func (o *Undefined) Length() int {
	return 0
}

// TypeName returns the name of the type.
func (o *Undefined) TypeName() string {
	return UndefinedType
}

func (o *Undefined) AsString() string {
	return ""
}

// Copy returns a copy of the type.
func (o *Undefined) Copy(_ int, _ int) IObject {
	return o
}

// Falsy returns true.
func (o *Undefined) Falsy() bool {
	return true
}

// Equals returns true if the values of the type are equal to the values of
// another object.
func (o *Undefined) Equals(x IObject) bool {
	return o == x
}

// IndexGet returns an element at a given index.
func (o *Undefined) IndexGet(_ int, _ IObject) (IObject, error) {
	return o.GateKeeper().UndefinedValue(), nil
}

// Iterate creates a map iterator.
func (o *Undefined) Iterate(_ int) IIterator {
	return o
}

// Iterable returns whether the IObject can be Iterated.
func (o *Undefined) Iterable() bool {
	return true
}

// Next returns true if there are more elements to iterate.
func (o *Undefined) Next() bool {
	return false
}

// Key returns the key or index values of the current element.
func (o *Undefined) Key(_ int) IObject {
	return o
}

// Value returns the values of the current element.
func (o *Undefined) Value(_ int) IObject {
	return o
}

// Count returns the total number of elements in the instance and its sub-elements.
func (o *Undefined) Count() int {
	return 1
}
