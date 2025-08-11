package objects

// CallableFunc is a function signature for the callable functions.
type CallableFunc = func(args ...Object) (ret Object, err error)

// BuiltinFunction represents a builtin function.
type BuiltinFunction struct {
	ObjectImpl
	Name  string
	Value CallableFunc
}

// TypeName returns the name of the type.
func (o *BuiltinFunction) TypeName() string {
	return "builtin-function:" + o.Name
}

func (o *BuiltinFunction) String() string {
	return "<builtin-function>"
}

// Copy returns a copy of the type.
func (o *BuiltinFunction) Copy() Object {
	return &BuiltinFunction{Value: o.Value}
}

// Equals returns true if the value of the type is equal to the value of another object.
func (o *BuiltinFunction) Equals(_ Object) bool {
	return false
}

// Call executes a builtin function.
func (o *BuiltinFunction) Call(args ...Object) (Object, error) {
	return o.Value(args...)
}

// CanCall returns whether the Object can be Called.
func (o *BuiltinFunction) CanCall() bool {
	return true
}
