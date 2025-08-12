package objects

// FunctionUser represents a user-defined callable function object.
// It embeds ObjectImpl to inherit default behavior and implements a callable interface.
// Name is the identifier for the function.
// Value holds the callable logic of the function as a CallableFunc.
// EncodingID is an identifier used for serialization or external encoding purposes.
type FunctionUser struct {
	ObjectImpl
	name       string
	value      CallableFunc
	encodingID string
}

// NewFunctionUser creates a new FunctionUser instance with the specified ID and callable function.
func NewFunctionUser(id string, fn CallableFunc) *FunctionUser {
	return &FunctionUser{
		name:  id,
		value: fn,
	}
}

// TypeName returns the type name of the FunctionUser as a string, prefixed with "user-function:".
func (o *FunctionUser) TypeName() string {
	return "user-function:" + o.name
}

// String returns the string representation of a FunctionUser object, which is always "<user-function>".
func (o *FunctionUser) String() string {
	return "<user-function>"
}

// Copy creates and returns a new FunctionUser instance with the same Value field as the original object.
func (o *FunctionUser) Copy() IObject {
	return &FunctionUser{value: o.value}
}

// Equals checks whether the current FunctionUser is equal to another object of type IObject. Always returns false.
func (o *FunctionUser) Equals(_ IObject) bool {
	return false
}

// Call invokes the function encapsulated within the FunctionUser with the provided arguments and returns the result or an error.
func (o *FunctionUser) Call(args ...IObject) (IObject, error) {
	return o.value(args...)
}

// CanCall checks whether the FunctionUser instance can be invoked as a callable function. Always returns true.
func (o *FunctionUser) CanCall() bool {
	return true
}
