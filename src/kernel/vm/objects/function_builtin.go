package objects

// CallableFunc is a type alias for a function that takes a variadic list of IObject arguments and returns an IObject and an error.
type CallableFunc = func(args ...IObject) (ret IObject, err error)

// FunctionBuiltin represents a callable built-in function with a name and specific implementation.
type FunctionBuiltin struct {
	ObjectImpl
	name  string
	value CallableFunc
}

// NewFunctionBuiltin creates a new FunctionBuiltin instance with the given name and CallableFunc.
func NewFunctionBuiltin(name string, fn CallableFunc) *FunctionBuiltin {
	return &FunctionBuiltin{name: name, value: fn}
}

// Name returns the name of the FunctionBuiltin as a string.
func (o *FunctionBuiltin) Name() string {
	return o.name
}

// Value returns the CallableFunc associated with the FunctionBuiltin.
func (o *FunctionBuiltin) Value() CallableFunc {
	return o.value
}

// TypeName returns the type name of the FunctionBuiltin in the format "builtin-function:<name>".
func (o *FunctionBuiltin) TypeName() string {
	return "builtin-function:" + o.name
}

// String returns the string representation of the FunctionBuiltin object.
func (o *FunctionBuiltin) String() string {
	return "<builtin-function>"
}

// Copy creates and returns a new instance of FunctionBuiltin with the same values field as the original object.
func (o *FunctionBuiltin) Copy() IObject {
	return &FunctionBuiltin{value: o.value}
}

// Equals determines whether this FunctionBuiltin object is equal to another IObject. Always returns false.
func (o *FunctionBuiltin) Equals(_ IObject) bool {
	return false
}

// Call invokes the built-in function with the provided arguments and returns the result or an error if any occurs.
func (o *FunctionBuiltin) Call(args ...IObject) (IObject, error) {
	return o.value(args...)
}

// CanCall returns true indicating that the FunctionBuiltin object is callable.
func (o *FunctionBuiltin) CanCall() bool {
	return true
}
