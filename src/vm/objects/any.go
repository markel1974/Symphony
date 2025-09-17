package objects

import (
	"fmt"
	"reflect"
)

// Any is a flexible type wrapping a Go-native value, offering reflection-based type introspection and manipulation.
type Any struct {
	IAllocator
	data interface{}
	v    reflect.Value
	t    reflect.Type
}

// newAny creates a new instance of Any, initializing it with the provided allocator and value.
func newAny(allocator IAllocator, value interface{}) IObject {
	valOf := reflect.ValueOf(value)
	return &Any{
		IAllocator: allocator,
		data:       value,
		v:          valOf,
		t:          valOf.Type(),
	}
}

// setAllocator assigns a custom memory allocator for managing lifecycle and memory references of the object.
func (o *Any) setAllocator(allocator IAllocator) {
	o.IAllocator = allocator
}

// TypeName returns the string representation of the cached reflected type of the object.
func (o *Any) TypeName() string {
	return o.t.String()
}

// AsString converts the stored data in the Any object to a string representation. If the data is nil, it returns "<nil>".
func (o *Any) AsString() string {
	if o.data == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%v", o.data)
}

// AsBool returns the boolean representation of the underlying value, interpreted as !o.Falsy().
func (o *Any) AsBool() bool {
	return !o.Falsy()
}

// AsInt64 converts the value held by the Any object to an int64 if possible; otherwise, it returns 0.
func (o *Any) AsInt64() int64 {
	if o.v.CanInt() {
		return o.v.Int()
	}
	return 0
}

// AsFloat64 returns the value of the Any object as a float64 if it is convertible; otherwise, it returns 0.
func (o *Any) AsFloat64() float64 {
	if o.v.CanFloat() {
		return o.v.Float()
	}
	return 0
}

// Falsy checks if the contained value is considered "falsy." Returns true if the value is nil or has a zero value.
func (o *Any) Falsy() bool {
	if o.data == nil {
		return true
	}
	return o.v.IsZero()
}

// Equals compares the current instance with another IObject and returns true if they are equivalent, false otherwise.
func (o *Any) Equals(other IObject) bool {
	if otherAny, ok := other.(*Any); ok {
		return o.data == otherAny.data
	}
	return false
}

// Copy creates a new instance of the object with the provided frame and retains the original data.
func (o *Any) Copy(frame int, _ int) IObject {
	return o.GateKeeper().NewAny(frame, o.data)
}

// IndexGet retrieves the value of a field or method in the object based on the provided index as key.
func (o *Any) IndexGet(frame int, index IObject) (IObject, error) {
	key, ok := o.GateKeeper().ToString(index)
	if !ok {
		return nil, ErrIndexInvalidType
	}

	v := o.v
	if v.Kind() == reflect.Ptr {
		v = v.Elem() // Dereference pointer to access fields
	}
	// 1. Try to access a struct field
	if v.Kind() == reflect.Struct {
		field := v.FieldByName(key)
		if field.IsValid() {
			return o.GateKeeper().FromInterface(frame, field.Interface()), nil
		}
	}
	// 2. If not a field, try to find a method
	method := o.v.MethodByName(key)
	if method.IsValid() {
		// Wrap the Go method in a function that the VM can call
		return o.GateKeeper().NewFuncImport(frame, key, -1, func(gk IGateKeeper, f int, args ...IObject) (uint, IObject, error) {
			// Convert IObject arguments to reflect.Value
			if method.Type().NumIn() != len(args) {
				return 0, nil, fmt.Errorf("wrong number of arguments for method %s: want %d, got %d", key, method.Type().NumIn(), len(args))
			}
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
				// Special handling for 'error' type
				if err, isErr := results[0].Interface().(error); isErr {
					if err != nil {
						return 1, gk.NewError(f, err.Error()), nil
					}
					return 1, gk.UndefinedValue(), nil // nil error
				}
				return 1, gk.FromInterface(f, results[0].Interface()), nil
			}
			// Handle multiple return values by packing them into an Array
			retArray := make([]IObject, len(results))
			for i, res := range results {
				retArray[i] = gk.FromInterface(f, res.Interface())
			}
			return 1, gk.NewArray(f, retArray), nil
		}), nil
	}
	return nil, fmt.Errorf("field or method '%s' not found on type '%s'", key, o.TypeName())
}

// IndexSet sets a value for a specific index on the object, usually by modifying a field of a struct via reflection.
// It returns an error if the object is not a pointer to a struct, the index type is not valid, or the field cannot be set.
func (o *Any) IndexSet(index IObject, value IObject) error {
	if o.v.Kind() != reflect.Ptr || o.v.Elem().Kind() != reflect.Struct {
		return ErrNotAssignable // Puoi modificare solo i campi di un puntatore a struct
	}
	key, ok := o.GateKeeper().ToString(index)
	if !ok {
		return ErrIndexInvalidType
	}
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

// AssignValue assigns the value of the given IObject to the current object. Returns an error if assignment is not supported.
func (o *Any) AssignValue(v IObject) error {
	return ErrNotAssignable
}

// Nil returns true if the underlying data of the Any object is nil, otherwise false.
func (o *Any) Nil() bool {
	return o.data == nil
}

// LogicalOp performs a logical operation on the current object using the specified operator and right-hand operand.
// It returns the result of the operation or an error if the operation is not supported.
func (o *Any) LogicalOp(_ int, _ LogicalOperator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// ArithmeticOp performs an arithmetic operation using the given operator and right-hand-side operand.
// Returns an error if the operation is not supported or invalid.
func (o *Any) ArithmeticOp(_ int, _ ArithmeticOperator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// Iterate returns an IIterator instance for iterating over the underlying data. Always returns nil for non-iterable types.
func (o *Any) Iterate(_ int) IIterator {
	return nil
}

// Iterable determines if the object can be iterated over, returning false if iteration is not supported.
func (o *Any) Iterable() bool {
	return false
}

// Call invokes the object as a callable with the specified frame and arguments.
// Returns the execution status, result object, and an error if the object is not callable.
func (o *Any) Call(frame int, args ...IObject) (uint, IObject, error) {
	if o.v.Kind() == reflect.Func {
		v, _ := o.IndexGet(frame, o.GateKeeper().NewString(frame, ""))
		return 0, v, nil
	}
	return 0, nil, fmt.Errorf("object of type %s is not callable", o.TypeName())
}

// Length returns the length of the underlying value if it is an Array, Chan, Map, Slice, or String. Otherwise, it returns 0.
func (o *Any) Length() int {
	switch o.v.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return o.v.Len()
	default:
		return 0
	}
}

// Count returns the static value of 1 for an Any instance.
func (o *Any) Count() int {
	return 1
}
