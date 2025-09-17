package objects

import (
	"fmt"
	"reflect"
)

// Any represents a flexible polymorphic type that wraps and operates on arbitrary data using reflection and allocation.
type Any struct {
	IAllocator
	data    interface{}
	valueOf reflect.Value
	kind    reflect.Type
}

// newAny creates a new instance of Any, encapsulating the provided value and associating it with the specified allocator.
func newAny(allocator IAllocator, value interface{}) IObject {
	valueOf := reflect.ValueOf(value)
	return &Any{
		IAllocator: allocator,
		data:       value,
		valueOf:    valueOf,
		kind:       valueOf.Type(),
	}
}

// setAllocator assigns the specified IAllocator to the Any instance for memory management and object lifecycle control.
func (o *Any) setAllocator(allocator IAllocator) {
	o.IAllocator = allocator
}

// TypeName returns the name of the underlying type as a string.
func (o *Any) TypeName() string {
	return o.kind.String()
}

// AsInterface converts the object into a generic interface{} type and returns the underlying data.
func (o *Any) AsInterface() interface{} {
	return o.data
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
	if o.valueOf.CanInt() {
		return o.valueOf.Int()
	}
	return 0
}

// AsFloat64 returns the value of `Any` as a float64 if convertible, otherwise returns 0.
func (o *Any) AsFloat64() float64 {
	if o.valueOf.CanFloat() {
		return o.valueOf.Float()
	}
	return 0
}

// AsBytes converts the object elements into a single concatenated slice of bytes by calling AsBytes on each element.
func (o *Any) AsBytes() []byte {
	return nil
}

// Falsy checks if the underlying data is nil or if the value is considered zero, returning true in those cases.
func (o *Any) Falsy() bool {
	if o.data == nil {
		return true
	}
	return o.valueOf.IsZero()
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
	v := o.valueOf
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, ErrIndexNotIndexable
	}
	key := index.AsString()
	if field := v.FieldByName(key); field.IsValid() {
		return o.GateKeeper().FromInterface(frame, field.Interface()), nil
	}
	if method := o.valueOf.MethodByName(key); method.IsValid() {
		return o.GateKeeper().NewFuncImport(frame, key, -1, func(gk IGateKeeper, f int, args ...IObject) (uint, IObject, error) {
			if method.Type().NumIn() != len(args) && !method.Type().IsVariadic() {
				return 0, nil, fmt.Errorf("wrong number of arguments for method %s: want %d, got %d", key, method.Type().NumIn(), len(args))
			}
			return o.call(gk, f, method, args)
		}), nil
	}
	return nil, fmt.Errorf("field or method '%s' not found on type '%s'", key, o.TypeName())
}

// IndexSet updates the value of a struct field, specified by the index, with the given value if the field is assignable.
func (o *Any) IndexSet(index IObject, value IObject) error {
	if o.valueOf.Kind() != reflect.Ptr || o.valueOf.Elem().Kind() != reflect.Struct {
		return ErrNotAssignable
	}
	key := index.AsString()
	field := o.valueOf.Elem().FieldByName(key)
	if !field.IsValid() || !field.CanSet() {
		return fmt.Errorf("field '%s' not found or cannot be set", key)
	}
	valToSet := reflect.ValueOf(value.AsInterface())
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
	if o.valueOf.Kind() == reflect.Func {
		if o.kind.NumIn() != len(args) && !o.kind.IsVariadic() {
			return 0, nil, fmt.Errorf("wrong number of arguments for %s: want %d, got %d", o.TypeName(), o.kind.NumIn(), len(args))
		}
		return o.call(o.GateKeeper(), frame, o.valueOf, args)
	}
	s1, err := o.GateKeeper().ToStringArg(0, args)
	if err != nil {
		return 0, nil, ErrInvalidArgumentsNumber
	}
	method := o.valueOf.MethodByName(s1)
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
	return o.call(o.GateKeeper(), frame, method, methodArgs)
}

// Length returns the length of the underlying value if it is an array, channel, map, slice, or string; otherwise, 0.
func (o *Any) Length() int {
	switch o.valueOf.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return o.valueOf.Len()
	default:
		return 0
	}
}

// Count returns a fixed integer value of 1, representing a single object.
func (o *Any) Count() int {
	return 1
}

// call invokes a reflect.Value method with arguments, handles results, and converts to/from IObject representation.
func (o *Any) call(gk IGateKeeper, frame int, method reflect.Value, args []IObject) (uint, IObject, error) {
	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg.AsInterface())
	}
	results := method.Call(in)
	switch len(results) {
	case 0:
		return 0, gk.UndefinedValue(), nil
	case 1:
		if err, isErr := results[0].Interface().(error); isErr {
			if err != nil {
				return 1, gk.NewError(frame, err.Error()), nil
			}
			return 1, gk.UndefinedValue(), nil
		}
		return 1, gk.FromInterface(frame, results[0].Interface()), nil
	default:
		// Handle multiple return values by packing them into an Array
		retArray := make([]IObject, len(results))
		for i, res := range results {
			retArray[i] = gk.FromInterface(frame, res.Interface())
		}
		return 1, gk.NewArray(frame, retArray), nil
	}
}
