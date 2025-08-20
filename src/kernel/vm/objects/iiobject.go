package objects

// IObject defines a generic interface for objects that can perform various operations and support multiple behaviors.
// TypeName returns the type name of the object.
// String returns the string representation of the object.
// BinaryOp performs a binary operation between the object and a right-hand side operand.
// Boolean checks if the object represents a falsy value.
// Equals checks whether the object is equal to another object.
// Copy creates and returns a copy of the object.
// IndexGet retrieves the values at the specified index from the object.
// IndexSet assigns a value to the specified index within the object.
// Iterate returns an iterator for the object, enabling iteration.
// CanIterate checks if the object can be iterated over.
// Call invokes the object as a callable function with provided arguments.
// CanCall checks if the object can be called as a function.
// Length returns the length of the object.
type IObject interface {
	TypeName() string

	String() string

	BinaryOp(frame int, op Operator, rightHandSide IObject) (IObject, error)

	Boolean() bool

	Equals(other IObject) bool

	Copy(frame int) IObject

	IndexGet(frame int, index IObject) (value IObject, err error)

	IndexSet(index, value IObject) error

	Iterate(frame int) IIterator

	CanIterate() bool

	Call(args ...IObject) (ret IObject, err error)

	CanCall() bool

	Length() int
}
