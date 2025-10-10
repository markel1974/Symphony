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
	vTable  map[string]IObject
}

// newAny creates a new instance of Any, encapsulating the provided value and associating it with the specified allocator.
func newAny(allocator IAllocator, value interface{}) IObject {
	a := &Any{
		IAllocator: allocator,
		data:       value,
		vTable:     make(map[string]IObject),
	}
	return a
}

// setAllocator assigns the specified IAllocator to the Any instance for memory management and object lifecycle control.
func (o *Any) setAllocator(allocator IAllocator) {
	o.IAllocator = allocator
}

// Setup initializes the object with a frame ID and a value, setting up its metadata, methods, and virtual table.
func (o *Any) Setup(frameId int, value interface{}) {
	o.setFrame(frameId)
	o.data = value
	o.valueOf = reflect.ValueOf(value)
	o.kind = o.valueOf.Type()
	o.vTable = make(map[string]IObject)
	for x := 0; x < o.valueOf.NumMethod(); x++ {
		m := o.valueOf.Type().Method(x)
		o.vTable[m.Name] = o.GateKeeper().NewFuncImport(frameId, m.Name, -1, func(gk IGateKeeper, f int, args ...IObject) (uint, IObject, error) {
			return o.call(frameId, m.Func, args)
		})
	}
}

// TypeName returns the name of the underlying type as a string.
func (o *Any) TypeName() string {
	return o.kind.String()
}

// AsInterface converts the object into a generic interface{} type and returns the underlying data.
func (o *Any) AsInterface() interface{} {
	return o.data
}

// AsValue attempts to cast the underlying value of Any to the specified target type and checks if it is assignable.
func (o *Any) AsValue(target reflect.Type) (reflect.Value, bool) {
	ret := o.valueOf.Type().AssignableTo(target)
	return o.valueOf, ret
}

// AsString returns the string representation of the wrapped data. If the data is nil, it returns "<nil>".
func (o *Any) AsString() string {
	return fmt.Sprintf("%v", o.data)
}

// AsBool returns the boolean representation of the Any object by negating its Falsy status.
func (o *Any) AsBool() bool {
	return o.bool()
}

// AsInt64 converts the underlying value to an int64 if possible, otherwise returns 0.
func (o *Any) AsInt64() int64 {
	if v, ok := o.int64(); ok {
		return v
	}
	if o.bool() {
		return 1
	}
	return 0
}

// AsFloat64 returns the value of `Any` as a float64 if convertible, otherwise returns 0.
func (o *Any) AsFloat64() float64 {
	v, _ := o.float64()
	return v
}

// AsBytes converts the object elements into a single concatenated slice of bytes by calling AsBytes on each element.
func (o *Any) AsBytes() []byte {
	return nil
}

// Falsy checks if the underlying data is nil or if the value is considered zero, returning true in those cases.
func (o *Any) Falsy() bool {
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
			return o.call(f, method, args)
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
	return o.valueOf.IsNil() || o.valueOf.IsZero()
}

// LogicalOp performs a logical operation on the current object with the provided operator and right-hand side object.
// It returns the result of the operation if the operator is valid, otherwise it returns an ErrInvalidOperator error.
func (o *Any) LogicalOp(_ int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		v, err := logicalOpNil(op)
		return o.GateKeeper().FromBoolError(v, err)
	}
	v, err := logicalOpInt64(o.AsInt64(), op, rhsIn.AsInt64())
	return o.GateKeeper().FromBoolError(v, err)
}

// ArithmeticOp performs an arithmetic operation on the current object and the provided IObject using the specified operator.
// Returns the result as an IObject or an error if the operation is invalid.
func (o *Any) ArithmeticOp(frame int, op ArithmeticOperator, rhsIn IObject) (IObject, error) {
	val := o.AsInt64()
	ret, err := arithmeticOpInt64(val, op, rhsIn.AsInt64())
	if err != nil {
		return nil, err
	}
	return o.GateKeeper().NewInt(frame, ret), nil
}

// UnaryOp performs a unary operation using the specified UnaryOperator. Returns a new object or an error.
func (o *Any) UnaryOp(frame int, op UnaryOperator) (IObject, error) {
	if val, ok := o.float64(); ok {
		r, err := unaryOpFloat64(op, val)
		if err != nil {
			return nil, err
		}
		return o.GateKeeper().NewFloat(frame, r), nil
	}
	val := o.AsInt64()
	r, err := unaryOpInt64(op, val)
	if err != nil {
		return nil, err
	}
	return o.GateKeeper().NewInt(frame, r), nil
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
		return o.call(frame, o.valueOf, args)
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
	return o.call(frame, method, methodArgs)
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

// Method retrieves an object from the iTable by name, returning the object and a boolean indicating success or failure.
func (o *Any) Method(name string) (IObject, bool) {
	m, ok := o.vTable[name]
	if !ok || m == nil {
		return o.GateKeeper().UndefinedValue(), false
	}
	return m, ok
}

// GobDecode implements the gob.GobDecoder interface, decoding the Allocator state from a byte slice representation.
func (o *Any) GobDecode(_ []byte) (err error) {
	return nil
}

// GobEncode encodes the Allocator instance into a byte slice for use with the gob package and returns it along with any error.
func (o *Any) GobEncode() ([]byte, error) {
	return nil, nil
}

// int64 attempts to convert the underlying value of Any to int64 and returns the result along with a success flag.
func (o *Any) int64() (int64, bool) {
	if o.valueOf.CanInt() {
		return o.valueOf.Int(), true
	} else if o.valueOf.CanUint() {
		return int64(o.valueOf.Uint()), true
	}
	return 0, false
}

// float64 attempts to convert the internal value of Any to a float64 and returns the value along with a success flag.
func (o *Any) float64() (float64, bool) {
	if o.valueOf.CanFloat() {
		return o.valueOf.Float(), true
	}
	return 0, false
}

// bool returns a boolean value indicating whether the data in the Any object is nil or represents a zero value.
func (o *Any) bool() bool {
	if o.valueOf.IsZero() || o.valueOf.IsNil() {
		return false
	}
	return true
}

// call invokes the provided function with specified arguments and ensures type compatibility for the invocation.
// It validates argument count and type, supports variadic functions, and handles return values, including errors.
// Returns a status code, the result of the call, or an error in case of a failure.
func (o *Any) call(frameId int, mFunc reflect.Value, args []IObject) (uint, IObject, error) {
	if mFunc.Type().NumIn() != len(args) && !mFunc.Type().IsVariadic() {
		return 0, nil, fmt.Errorf("wrong number of arguments: want %d, got %d", mFunc.Type().NumIn(), len(args))
	}
	var in []reflect.Value
	if len(args) > 0 {
		for idx, arg := range args {
			target := mFunc.Type().In(idx)
			val, ok := arg.AsValue(target)
			if !ok {
				return 0, nil, fmt.Errorf("wrong type for argument %d: want %s, got %s", idx, target, arg.TypeName())
			}
			if !val.Type().AssignableTo(target) {
				return 0, nil, fmt.Errorf("wrong type for argument %d: want %s, got %s", idx, target, val.Type())
			}
			in = append(in, val)
		}
	}
	//kk := m.Func.Type().In(0)
	//m.Func.Type().Out(0)
	results := mFunc.Call(in)
	switch len(results) {
	case 0:
		return 0, o.GateKeeper().UndefinedValue(), nil
	case 1:
		if err, isErr := results[0].Interface().(error); isErr {
			if err != nil {
				return 1, o.GateKeeper().NewError(frameId, err.Error()), nil
			}
			return 1, o.GateKeeper().UndefinedValue(), nil
		}
		return 1, o.GateKeeper().FromInterface(frameId, results[0].Interface()), nil
	default:
		// Handle multiple return values by packing them into an Array
		retArray := make([]IObject, len(results))
		for i, res := range results {
			retArray[i] = o.GateKeeper().FromInterface(frameId, res.Interface())
		}
		return 1, o.GateKeeper().NewArray(frameId, retArray), nil
	}
}
