package objects

// UserFunction represents a user-defined callable function object.
// It embeds ObjectImpl to inherit default behavior and implements a callable interface.
// Name is the identifier for the function.
// Value holds the callable logic of the function as a CallableFunc.
// EncodingID is an identifier used for serialization or external encoding purposes.
type UserFunction struct {
	ObjectImpl
	name       string
	value      CallableFunc
	encodingID string
}

// NewUserFunction creates a new UserFunction instance with the specified ID and callable function.
func NewUserFunction(id string, fn CallableFunc) *UserFunction {
	return &UserFunction{
		name:  id,
		value: fn,
	}
}

// TypeName returns the type name of the UserFunction as a string, prefixed with "user-function:".
func (o *UserFunction) TypeName() string {
	return "user-function:" + o.name
}

// String returns the string representation of a UserFunction object, which is always "<user-function>".
func (o *UserFunction) String() string {
	return "<user-function>"
}

// Copy creates and returns a new UserFunction instance with the same Value field as the original object.
func (o *UserFunction) Copy() IObject {
	return &UserFunction{value: o.value}
}

// Equals checks whether the current UserFunction is equal to another object of type IObject. Always returns false.
func (o *UserFunction) Equals(_ IObject) bool {
	return false
}

// Call invokes the function encapsulated within the UserFunction with the provided arguments and returns the result or an error.
func (o *UserFunction) Call(args ...IObject) (IObject, error) {
	return o.value(args...)
}

// CanCall checks whether the UserFunction instance can be invoked as a callable function. Always returns true.
func (o *UserFunction) CanCall() bool {
	return true
}
