package objects

import (
	"encoding/gob"
	"fmt"
)

func init() {
	gob.Register(&Builtin{})
}

// BuiltinType is a constant string identifying the type name of builtin objects.
const (
	BuiltinType = "builtin"
)

// Builtin defines a struct representing a built-in object with a name and an integer value.
type Builtin struct {
	gk    IGateKeeper
	frame int
	name  string
	index int
}

// NewBuiltin creates a new instance of Builtin with the specified name and value.
func newBuiltin(factory IGateKeeper, frame int, name string, index int) IObject {
	return &Builtin{
		gk:    factory,
		frame: frame,
		name:  name,
		index: index,
	}
}

// GateKeeper returns a reference to the GateKeeper associated with the Object.
func (o *Builtin) GateKeeper() IGateKeeper {
	return o.gk
}

// Frame returns the current frame value of the Object.
func (o *Builtin) Frame() int {
	return o.frame
}

// LogicalOp performs a logical operation on the Builtin object using the given operator and operand, returning the result or an error.
func (o *Builtin) LogicalOp(_ int, _ LogicalOperator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// ArithmeticOp performs an arithmetic operation on the Builtin object using the specified operator and operand.
// Always returns ErrInvalidOperator as arithmetic operations are not supported.
func (o *Builtin) ArithmeticOp(_ int, _ ArithmeticOperator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// IndexGet attempts to retrieve a value at the given index and returns an error if the object is not indexable.
func (o *Builtin) IndexGet(_ int, _ IObject) (res IObject, err error) {
	return nil, ErrNotIndexable
}

// IndexSet attempts to assign a value to an index in the object but always returns ErrUnsupportedIndex,
// as this operation is unsupported.
func (o *Builtin) IndexSet(_, _ IObject) (err error) {
	return ErrUnsupportedIndex
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (o *Builtin) Iterate(_ int) IIterator {
	return nil
}

// CanIterate determines if the object can be iterated over and returns false for this implementation.
func (o *Builtin) CanIterate() bool {
	return false
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *Builtin) Call(_ int, _ ...IObject) (ret IObject, err error) {
	return nil, nil
}

// CanCall determines if the object can be invoked as a callable. Returns false for non-callable objects.
func (o *Builtin) CanCall() bool {
	return false
}

// Length returns the length of the Int object.
func (o *Builtin) Length() int {
	return 0
}

// Name returns the name of the Builtin object.
func (o *Builtin) Name() string {
	return o.name
}

// Value returns the integer value associated with the Builtin object.
func (o *Builtin) Value() int {
	return o.index
}

// String returns a string representation of the Builtin object, formatted as "name: value".
func (o *Builtin) String() string {
	return fmt.Sprintf("%s: %d", o.name, o.index)
}

// TypeName returns the type name of the Builtin object as a string.
func (o *Builtin) TypeName() string {
	return BuiltinType
}

// Copy creates and returns a new Builtin object with the same name and value as the current instance.
func (o *Builtin) Copy(frame int, _ int) IObject {
	return o.GateKeeper().NewBuiltin(frame, o.name, o.index)
}

// Falsy returns false, indicating the boolean representation of the Builtin object is always false.
func (o *Builtin) Falsy() bool {
	return false
}

// Equals checks if the current Builtin object is equal to another IObject based on its name and value.
func (o *Builtin) Equals(x IObject) bool {
	t, ok := x.(*Builtin)
	if !ok {
		return false
	}
	return o.name == t.name && o.index == t.index
}
