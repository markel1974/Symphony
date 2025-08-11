package objects

// CallableFunc is a type alias for a function that takes a variadic list of IObject arguments and returns an IObject and an error.
type CallableFunc = func(args ...IObject) (ret IObject, err error)

// BuiltinFunction represents a callable built-in function with a name and specific implementation.
type BuiltinFunction struct {
	ObjectImpl
	name  string
	value CallableFunc
}

// NewBuiltinFunction creates a new BuiltinFunction instance with the given name and CallableFunc.
func NewBuiltinFunction(name string, fn CallableFunc) *BuiltinFunction {
	return &BuiltinFunction{name: name, value: fn}
}

// Name returns the name of the BuiltinFunction as a string.
func (o *BuiltinFunction) Name() string {
	return o.name
}

// Value returns the CallableFunc associated with the BuiltinFunction.
func (o *BuiltinFunction) Value() CallableFunc {
	return o.value
}

// TypeName returns the type name of the BuiltinFunction in the format "builtin-function:<name>".
func (o *BuiltinFunction) TypeName() string {
	return "builtin-function:" + o.name
}

// String returns the string representation of the BuiltinFunction object.
func (o *BuiltinFunction) String() string {
	return "<builtin-function>"
}

// Copy creates and returns a new instance of BuiltinFunction with the same value field as the original object.
func (o *BuiltinFunction) Copy() IObject {
	return &BuiltinFunction{value: o.value}
}

// Equals determines whether this BuiltinFunction object is equal to another IObject. Always returns false.
func (o *BuiltinFunction) Equals(_ IObject) bool {
	return false
}

// Call invokes the built-in function with the provided arguments and returns the result or an error if any occurs.
func (o *BuiltinFunction) Call(args ...IObject) (IObject, error) {
	return o.value(args...)
}

// CanCall returns true indicating that the BuiltinFunction object is callable.
func (o *BuiltinFunction) CanCall() bool {
	return true
}
