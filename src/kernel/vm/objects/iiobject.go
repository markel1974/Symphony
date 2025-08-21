package objects

// IObject represents a generic interface for any object in the system.
// TypeName returns the name of the type of the object.
// String returns a string representation of the object.
// Frame retrieves the current execution frame associated with the object.
// BinaryOp performs a binary operation using the specified operator and operands, returning the result or an error.
// Boolean evaluates and returns the boolean value of the object.
// Equals returns true if the object is equal to another given object.
// Copy creates and returns a deep copy of the object within the given frame.
// IndexGet retrieves the value at a given index from the object, returning the value or an error.
// IndexSet updates the value at a given index within the object, returning an error if the operation fails.
// Iterate returns an iterator for the object, enabling traversal over its elements.
// CanIterate checks whether the object supports iteration.
// Call invokes the object as a callable function with the provided arguments, returning the result or an error.
// CanCall checks whether the object supports being called as a function.
// Length returns the length of the object, if applicable.
type IObject interface {
	TypeName() string

	String() string

	Frame() int

	BinaryOp(frame int, op Operator, rightHandSide IObject) (IObject, error)

	Boolean() bool

	Equals(other IObject) bool

	Copy(frame int) IObject

	IndexGet(frame int, index IObject) (value IObject, err error)

	IndexSet(index, value IObject) error

	Iterate(frame int) IIterator

	CanIterate() bool

	Call(frame int, args ...IObject) (ret IObject, err error)

	CanCall() bool

	Length() int
}
