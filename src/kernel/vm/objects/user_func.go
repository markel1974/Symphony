package objects

// UserFunction represents a user function.
type UserFunction struct {
	ObjectImpl
	Name       string
	Value      CallableFunc
	EncodingID string
}

func NewUserFunction(id string, fn CallableFunc) *UserFunction {
	return &UserFunction{
		Name:  "re_match",
		Value: fn,
	}
}

// TypeName returns the name of the type.
func (o *UserFunction) TypeName() string {
	return "user-function:" + o.Name
}

func (o *UserFunction) String() string {
	return "<user-function>"
}

// Copy returns a copy of the type.
func (o *UserFunction) Copy() IObject {
	return &UserFunction{Value: o.Value}
}

// Equals returns true if the values of the type is equal to the values of
// another object.
func (o *UserFunction) Equals(_ IObject) bool {
	return false
}

// Call invokes a user function.
func (o *UserFunction) Call(args ...IObject) (IObject, error) {
	return o.Value(args...)
}

// CanCall returns whether the IObject can be Called.
func (o *UserFunction) CanCall() bool {
	return true
}
