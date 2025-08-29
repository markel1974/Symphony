package objects

import "encoding/gob"

// InterfaceType defines a string constant representing the word "interface".
// InterfaceLabel defines a string constant combining "<" with InterfaceType and ">" to form a label.
const (
	InterfaceType  = "interface"
	InterfaceLabel = "<" + InterfaceType + ">"
)

// init registers the Interface type with the gob package for encoding and decoding operations.
func init() {
	gob.Register(&Interface{})
}

// Interface represents an object with contextual execution information and dynamic properties managed within a frame.
type Interface struct {
	factory IGateKeeper
	frame   int
	value   IObject
	iTable  map[string]IObject
}

// newInterface creates a new instance of Interface with the provided factory, frame ID, value, and interface table.
func newInterface(factory IGateKeeper, frame int, value IObject, itable map[string]IObject) IObject {
	return &Interface{
		factory: factory,
		frame:   frame,
		value:   value,
		iTable:  itable,
	}
}

// GateKeeper retrieves the IGateKeeper instance associated with the Interface.
func (o *Interface) GateKeeper() IGateKeeper {
	return o.factory
}

// Frame returns the current execution frame associated with the Interface instance.
func (o *Interface) Frame() int {
	return o.frame
}

// TypeName returns the type name of the underlying IObject.
func (o *Interface) TypeName() string {
	return o.value.TypeName()
}

// String returns the string representation of the Interface instance by delegating to the underlying IObject value.
func (o *Interface) String() string {
	return o.value.String()
}

// Falsy determines whether the underlying value of the object evaluates to a falsy value, returning true if it does.
func (o *Interface) Falsy() bool {
	return o.value.Falsy()
}

// Equals compares the current object with another IObject and returns true if they are equal.
func (o *Interface) Equals(other IObject) bool {
	return o.value.Equals(other)
}

// Copy creates and returns a deep copy of the current object with the specified frame and depth.
func (o *Interface) Copy(frame int, depth int) IObject {
	return o.value.Copy(frame, depth)
}

// BinaryOp performs a binary operation with the current object using the provided operator and right-hand side operand.
func (o *Interface) BinaryOp(frame int, op Operator, rightHandSide IObject) (IObject, error) {
	return o.value.BinaryOp(frame, op, rightHandSide)
}

// IndexGet retrieves the value at the specified index from the object, using the provided execution frame and index.
func (o *Interface) IndexGet(frame int, index IObject) (value IObject, err error) {
	return o.value.IndexGet(frame, index)
}

// IndexSet sets a value at the specified index in the IObject, returning an error if the operation fails.
func (o *Interface) IndexSet(index, value IObject) error {
	return o.value.IndexSet(index, value)
}

// Iterate returns an iterator for traversing over the elements of the value associated with the interface.
func (o *Interface) Iterate(frame int) IIterator {
	return o.value.Iterate(frame)
}

// CanIterate determines whether the wrapped IObject supports iteration, returning true if it does, or false otherwise.
func (o *Interface) CanIterate() bool {
	return o.value.CanIterate()
}

// Call invokes the object with the specified frame and arguments, returning the result or an error if unsupported.
func (o *Interface) Call(frame int, args ...IObject) (ret IObject, err error) {
	return nil, ErrUnsupportedOperation
}

// CanCall determines whether the object can be invoked as a callable function. Always returns false for this implementation.
func (o *Interface) CanCall() bool {
	return false
}

// Length returns the length of the underlying value represented by the IObject interface.
func (o *Interface) Length() int {
	return o.value.Length()
}

// Value returns the underlying IObject associated with the Interface instance.
func (o *Interface) Value() IObject {
	return o.value
}

// ITable returns the internal map representation `iTable` of the Interface instance.
func (o *Interface) ITable() map[string]IObject {
	return o.iTable
}
