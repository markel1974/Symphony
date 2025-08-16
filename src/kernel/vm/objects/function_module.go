package objects

const (
	FunctionBuiltinDef = "function_builtin"
	FunctionModuleDef  = "function_module"
)

// CallableFunc is a type alias for a function that takes a variadic list of IObject arguments and returns an IObject and an error.
type CallableFunc = func(args ...IObject) (ret IObject, err error)

// FunctionModule is a callable object type that encapsulates a function and provides execution context information.
type FunctionModule struct {
	Object
	kind  string
	name  string
	value CallableFunc
}

// NewFunctionModule creates a new FunctionModule instance with the specified ID and callable function.
func NewFunctionModule(kind string, name string, fn CallableFunc) *FunctionModule {
	return &FunctionModule{
		kind:  kind,
		name:  name,
		value: fn,
	}
}

// Name returns the name of the FunctionModule as a string.
func (o *FunctionModule) Name() string {
	return o.name
}

// TypeName returns the type name of the FunctionModule as a string.
func (o *FunctionModule) TypeName() string {
	return o.kind + ":" + o.name
}

// String returns the string representation of a FunctionModule object.
func (o *FunctionModule) String() string {
	return "<" + o.kind + ">"
}

// Copy creates and returns a new FunctionModule instance with the same Value field as the original object.
func (o *FunctionModule) Copy() IObject {
	return &FunctionModule{value: o.value}
}

// Equals checks whether the current FunctionModule is equal to another object of type IObject. Always returns false.
func (o *FunctionModule) Equals(_ IObject) bool {
	return false
}

// Call invokes the function encapsulated within the FunctionModule with the provided arguments and returns the result or an error.
func (o *FunctionModule) Call(args ...IObject) (IObject, error) {
	return o.value(args...)
}

// CanCall checks whether the FunctionModule instance can be invoked as a callable function. Always returns true.
func (o *FunctionModule) CanCall() bool {
	return true
}
