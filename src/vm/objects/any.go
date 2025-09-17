package objects

import (
	"fmt"
	"reflect"
)

// Any represents a flexible polymorphic type that wraps and operates on arbitrary data using reflection and allocation.
type Any struct {
	IAllocator
	data interface{}
	v    reflect.Value
	t    reflect.Type
}

// newAny creates a new instance of Any, encapsulating the provided value and associating it with the specified allocator.
func newAny(allocator IAllocator, value interface{}) IObject {
	valOf := reflect.ValueOf(value)
	return &Any{
		IAllocator: allocator,
		data:       value,
		v:          valOf,
		t:          valOf.Type(),
	}
}

// setAllocator assigns the specified IAllocator to the Any instance for memory management and object lifecycle control.
func (o *Any) setAllocator(allocator IAllocator) {
	o.IAllocator = allocator
}

// TypeName returns the name of the underlying type as a string.
func (o *Any) TypeName() string {
	return o.t.String()
}

// AsString returns the string representation of the wrapped data. If the data is nil, it returns "<nil>".
func (o *Any) AsString() string {
	if o.data == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%v", o.data)
}

// AsBool returns the boolean representation of the Any object by negating its Falsy status.
func (o *Any) AsBool() bool {
	return !o.Falsy()
}

// AsInt64 converts the underlying value to an int64 if possible, otherwise returns 0.
func (o *Any) AsInt64() int64 {
	if o.v.CanInt() {
		return o.v.Int()
	}
	return 0
}

// AsFloat64 returns the value of `Any` as a float64 if convertible, otherwise returns 0.
func (o *Any) AsFloat64() float64 {
	if o.v.CanFloat() {
		return o.v.Float()
	}
	return 0
}

// Falsy checks if the underlying data is nil or if the value is considered zero, returning true in those cases.
func (o *Any) Falsy() bool {
	if o.data == nil {
		return true
	}
	return o.v.IsZero()
}

// Equals compares the current object with another IObject and returns true if they are considered equal, otherwise false.
func (o *Any) Equals(other IObject) bool {
	if otherAny, ok := other.(*Any); ok {
		return o.data == otherAny.data
	}
	return false
}

// Copy creates a new instance of Any with the same data as the current object and associates it with the given frame.
func (o *Any) Copy(frame int, _ int) IObject {
	return o.GateKeeper().NewAny(frame, o.data)
}

// IndexGet retrieves a field or method of a struct or pointer to a struct by its name, provided as the index parameter.
// Returns an error if the index is not a valid key or if the field/method cannot be found.
// If the index refers to a method, the returned IObject will be callable.
func (o *Any) IndexGet(frame int, index IObject) (IObject, error) {
	key := index.AsString()
	v := o.v
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() == reflect.Struct {
		field := v.FieldByName(key)
		if field.IsValid() {
			return o.GateKeeper().FromInterface(frame, field.Interface()), nil
		}
	}
	if fn, ok := o.buildFunc(frame, key); ok {
		return fn, nil
	}
	return nil, fmt.Errorf("field or method '%s' not found on type '%s'", key, o.TypeName())
}

// IndexSet updates the value of a struct field, specified by the index, with the given value if the field is assignable.
func (o *Any) IndexSet(index IObject, value IObject) error {
	if o.v.Kind() != reflect.Ptr || o.v.Elem().Kind() != reflect.Struct {
		return ErrNotAssignable
	}
	key := index.AsString()
	field := o.v.Elem().FieldByName(key)
	if !field.IsValid() || !field.CanSet() {
		return fmt.Errorf("field '%s' not found or cannot be set", key)
	}
	valToSet := reflect.ValueOf(o.GateKeeper().ToInterface(value))
	if valToSet.Type().AssignableTo(field.Type()) {
		field.Set(valToSet)
		return nil
	}
	return fmt.Errorf("cannot assign type %s to field %s (type %s)", valToSet.Type(), key, field.Type())
}

// AssignValue assigns a value from another IObject but always returns ErrNotAssignable indicating the operation is unsupported.
func (o *Any) AssignValue(_ IObject) error {
	return ErrNotAssignable
}

// Nil checks if the `data` field of the `Any` object is nil and returns true if it is, otherwise false.
func (o *Any) Nil() bool {
	return o.data == nil
}

// LogicalOp performs a logical operation on the current object with the provided operator and right-hand side object.
// It returns the result of the operation if the operator is valid, otherwise it returns an ErrInvalidOperator error.
func (o *Any) LogicalOp(_ int, _ LogicalOperator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// ArithmeticOp performs an arithmetic operation on the current object and the provided IObject using the specified operator.
// Returns the result as an IObject or an error if the operation is invalid.
func (o *Any) ArithmeticOp(_ int, _ ArithmeticOperator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// Iterate returns an iterator (IIterator) that traverses the object's data. Returns nil if the object is not iterable.
func (o *Any) Iterate(_ int) IIterator {
	return nil
}

// Iterable checks if the object is iterable and returns false, indicating it is not iterable by default.
func (o *Any) Iterable() bool {
	return false
}

// Call invokes a function or a method on the object with the provided arguments and returns the result or an error.
func (o *Any) Call(frame int, args ...IObject) (uint, IObject, error) {
	if o.v.Kind() == reflect.Func {
		if o.t.NumIn() != len(args) && !o.t.IsVariadic() {
			return 0, nil, fmt.Errorf("wrong number of arguments for %s: want %d, got %d", o.TypeName(), o.t.NumIn(), len(args))
		}
		return o.call(o.v, o.GateKeeper(), frame, args)
	}
	s1, err := o.GateKeeper().ToStringArg(0, args)
	if err != nil {
		return 0, nil, ErrInvalidArgumentsNumber
	}
	method := o.v.MethodByName(s1)
	if !method.IsValid() {
		return 0, nil, fmt.Errorf("function not found on type '%s'", o.TypeName())
	}
	var methodArgs []IObject
	if len(args) > 1 {
		methodArgs = args[1:]
	}
	if method.Type().NumIn() != len(methodArgs) && !method.Type().IsVariadic() {
		return 0, nil, fmt.Errorf("wrong number of arguments for method %s: want %d, got %d", s1, method.Type().NumIn(), len(methodArgs))
	}
	return o.call(o.v, o.GateKeeper(), frame, methodArgs)
}

// Length returns the length of the underlying value if it is an array, channel, map, slice, or string; otherwise, 0.
func (o *Any) Length() int {
	switch o.v.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return o.v.Len()
	default:
		return 0
	}
}

// Count returns a fixed integer value of 1, representing a single object.
func (o *Any) Count() int {
	return 1
}

// buildFunc defines and returns a method wrapper for the specified key if it exists, otherwise returns false.
func (o *Any) buildFunc(frame int, key string) (IObject, bool) {
	method := o.v.MethodByName(key)
	if !method.IsValid() {
		return nil, false
	}
	return o.GateKeeper().NewFuncImport(frame, key, -1, func(gk IGateKeeper, f int, args ...IObject) (uint, IObject, error) {
		if method.Type().NumIn() != len(args) && !method.Type().IsVariadic() {
			return 0, nil, fmt.Errorf("wrong number of arguments for method %s: want %d, got %d", key, method.Type().NumIn(), len(args))
		}
		return o.call(method, gk, f, args)
	}), true
}

// call invokes a reflect.Value method with arguments, handles results, and converts to/from IObject representation.
func (o *Any) call(method reflect.Value, gk IGateKeeper, frame int, args []IObject) (uint, IObject, error) {
	// Convert IObject arguments to reflect.Value
	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = reflect.ValueOf(gk.ToInterface(arg))
	}
	// Call the native Go method
	results := method.Call(in)
	// Convert results back to IObject
	if len(results) == 0 {
		return 0, gk.UndefinedValue(), nil
	}
	// Simplification: handle single return value, but could be extended
	if len(results) == 1 {
		// Special handling for 'error'
		if err, isErr := results[0].Interface().(error); isErr {
			if err != nil {
				return 1, gk.NewError(frame, err.Error()), nil
			}
			return 1, gk.UndefinedValue(), nil // nil error
		}
		return 1, gk.FromInterface(frame, results[0].Interface()), nil
	}
	// Handle multiple return values by packing them into an Array
	retArray := make([]IObject, len(results))
	for i, res := range results {
		retArray[i] = gk.FromInterface(frame, res.Interface())
	}
	return 1, gk.NewArray(frame, retArray), nil
}
