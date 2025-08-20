package objects

const (
	FuncBuiltinDef = "func_builtin"
	FuncPackageDef = "func_package"
)

// FuncCallable is a type alias for a function that takes a variadic list of IObject arguments and returns an IObject and an error.
type FuncCallable = func(args ...IObject) (ret IObject, err error)

// FuncPackage is a callable object type that encapsulates a function and provides execution context information.
type FuncPackage struct {
	*Object
	kind  string
	name  string
	value FuncCallable
}

// NewFuncPackage creates a new FuncPackage instance with the specified ID and callable function.
func _newFuncPackage(factory *Factory, kind string, name string, fn FuncCallable) *FuncPackage {
	return &FuncPackage{
		Object: factory.NewObject(),
		kind:   kind,
		name:   name,
		value:  fn,
	}
}

// Name returns the name of the FuncPackage as a string.
func (o *FuncPackage) Name() string {
	return o.name
}

// TypeName returns the type name of the FuncPackage as a string.
func (o *FuncPackage) TypeName() string {
	return o.kind + ":" + o.name
}

// String returns the string representation of a FuncPackage object.
func (o *FuncPackage) String() string {
	return "<" + o.kind + ">"
}

// Copy creates and returns a new FuncPackage instance with the same Value field as the original object.
func (o *FuncPackage) Copy() IObject {
	return o.Factory().NewFuncPackage(o.kind, o.name, o.value)
}

// Equals checks whether the current FuncPackage is equal to another object of type IObject. Always returns false.
func (o *FuncPackage) Equals(_ IObject) bool {
	return false
}

// Call invokes the function encapsulated within the FuncPackage with the provided arguments and returns the result or an error.
func (o *FuncPackage) Call(args ...IObject) (IObject, error) {
	return o.value(args...)
}

// CanCall checks whether the FuncPackage instance can be invoked as a callable function. Always returns true.
func (o *FuncPackage) CanCall() bool {
	return true
}
