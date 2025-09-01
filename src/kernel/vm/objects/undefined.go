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
	gk    IGateKeeper
	frame int
}

func newUndefined(factory IGateKeeper, frame int) IObject {
	return &Undefined{
		gk:    factory,
		frame: frame,
	}
}

// GateKeeper returns a reference to the GateKeeper associated with the Object.
func (o *Undefined) GateKeeper() IGateKeeper {
	return o.gk
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

// SetStatic sets the frame to FrameStatic, marking it with a static execution context.
func (o *Undefined) SetStatic() {
	o.frame = FrameStatic
}

// Frame returns the current frame value of the Object.
func (o *Undefined) Frame() int {
	return o.frame
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

// IndexSet attempts to assign a value to an index in the object but always returns ErrUnsupportedIndex,
// as this operation is unsupported.
func (o *Undefined) IndexSet(_, _ IObject) (err error) {
	return ErrUnsupportedIndex
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *Undefined) Call(_ int, _ ...IObject) (ret IObject, err error) {
	return nil, nil
}

// CanCall determines if the object can be invoked as a callable. Returns false for non-callable objects.
func (o *Undefined) CanCall() bool {
	return false
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
	return UndefinedLabel
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

// CanIterate returns whether the IObject can be Iterated.
func (o *Undefined) CanIterate() bool {
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
