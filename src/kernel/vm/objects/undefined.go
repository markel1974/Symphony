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
	factory IGateKeeper
	frame   int
}

func newUndefined(factory IGateKeeper, frame int) IObject {
	return &Undefined{
		factory: factory,
		frame:   frame,
	}
}

// GateKeeper returns a reference to the GateKeeper associated with the Object.
func (o *Undefined) GateKeeper() IGateKeeper {
	return o.factory
}

// Frame returns the current frame value of the Object.
func (o *Undefined) Frame() int {
	return o.frame
}

// BinaryOp performs a binary operation on the current object and another object using the specified operator.
// Returns the result of the operation or an error if the operation is not supported.
func (o *Undefined) BinaryOp(_ int, _ Operator, _ IObject) (IObject, error) {
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

func (o *Undefined) String() string {
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
