package objects

import (
	"time"
)

// IObject defines a generic interface for objects that can perform various operations and support multiple behaviors.
// TypeName returns the type name of the object.
// String returns the string representation of the object.
// BinaryOp performs a binary operation between the object and a right-hand side operand.
// Falsy checks if the object represents a falsy values.
// Equals checks whether the object is equal to another object.
// Copy creates and returns a copy of the object.
// IndexGet retrieves the values at the specified index from the object.
// IndexSet assigns a values to the specified index within the object.
// Iterate returns an iterator for the object, enabling iteration.
// CanIterate checks if the object can be iterated over.
// Call invokes the object as a callable function with provided arguments.
// CanCall checks if the object can be called as a function.
type IObject interface {
	TypeName() string

	String() string

	BinaryOp(op Operator, rhs IObject) (IObject, error)

	Falsy() bool

	Equals(another IObject) bool

	Copy() IObject

	IndexGet(index IObject) (value IObject, err error)

	IndexSet(index, value IObject) error

	Iterate() IIterator

	CanIterate() bool

	Call(args ...IObject) (ret IObject, err error)

	CanCall() bool
}

// ObjectImpl is a default implementation of the IObject interface with unimplemented or default behavior for methods.
type ObjectImpl struct {
}

// TypeName returns the name of the type as a string. This method must be implemented by objects inheriting ObjectImpl.
func (o *ObjectImpl) TypeName() string {
	panic(ErrNotImplemented)
}

// String returns the string representation of the ObjectImpl. Currently, it is not implemented and will panic.
func (o *ObjectImpl) String() string {
	panic(ErrNotImplemented)
}

// BinaryOp performs a binary operation on the current object and another object using the specified operator.
// Returns the result of the operation or an error if the operation is not supported.
func (o *ObjectImpl) BinaryOp(_ Operator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// Copy creates and returns a new instance of the object, duplicating its state.
func (o *ObjectImpl) Copy() IObject {
	return nil
}

// Falsy returns false, indicating the object is not considered falsy in a boolean context.
func (o *ObjectImpl) Falsy() bool {
	return false
}

// Equals checks whether the current object is equal to another object of type IObject.
func (o *ObjectImpl) Equals(x IObject) bool {
	return o == x
}

// IndexGet attempts to retrieve a values at the given index and returns an error if the object is not indexable.
func (o *ObjectImpl) IndexGet(_ IObject) (res IObject, err error) {
	return nil, ErrNotIndexable
}

// IndexSet attempts to assign a values to an index in the object but always returns ErrNotIndexAssignable, as this operation is unsupported.
func (o *ObjectImpl) IndexSet(_, _ IObject) (err error) {
	return ErrNotIndexAssignable
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (o *ObjectImpl) Iterate() IIterator {
	return nil
}

// CanIterate determines if the object can be iterated over and returns false for this implementation.
func (o *ObjectImpl) CanIterate() bool {
	return false
}

// Call invokes the ObjectImpl with the provided arguments, returning a result object and an error, if any.
func (o *ObjectImpl) Call(_ ...IObject) (ret IObject, err error) {
	return nil, nil
}

// CanCall determines if the object can be invoked as a callable. Returns false for non-callable objects.
func (o *ObjectImpl) CanCall() bool {
	return false
}

// ToInterface converts an IObject to its corresponding native Go representation, such as int, string, float64, bool, etc.
func ToInterface(in IObject) (res interface{}) {
	switch o := in.(type) {
	case *Int:
		res = o.value
	case *String:
		res = o.value
	case *Float:
		res = o.value
	case *Bool:
		res = o == TrueValue
	case *Char:
		res = o.value
	case *Bytes:
		res = o.values
	case *Array:
		res = make([]interface{}, len(o.Values()))
		for i, val := range o.Values() {
			res.([]interface{})[i] = ToInterface(val)
		}
	case *ImmutableArray:
		res = make([]interface{}, o.Length())
		for i, val := range o.Values() {
			res.([]interface{})[i] = ToInterface(val)
		}
	case *Map:
		res = make(map[string]interface{})
		for key, v := range o.values {
			res.(map[string]interface{})[key] = ToInterface(v)
		}
	case *ImmutableMap:
		res = make(map[string]interface{})
		for key, v := range o.Values() {
			res.(map[string]interface{})[key] = ToInterface(v)
		}
	case *Time:
		res = o.value
	case *Error:
		res = New(o.String())
	case *Undefined:
		res = nil
	case IObject:
		return o
	}
	return
}

// FromMap converts a map with string keys and interface{} values into a map with string keys and IObject values.
func FromMap(v map[string]interface{}) map[string]IObject {
	kv := make(map[string]IObject)
	for key, val := range v {
		kv[key] = FromInterface(val)
	}
	return kv
}

// FromInterface converts a native Go values of various types into a corresponding IObject implementation.
func FromInterface(in interface{}) IObject {
	switch v := in.(type) {
	case nil:
		return UndefinedValue
	case string:
		if len(v) > MaxStringLen {
			return &String{value: v[0:MaxStringLen]}
		}
		return &String{value: v}
	case int64:
		return &Int{value: v}
	case int:
		return &Int{value: int64(v)}
	case bool:
		if v {
			return TrueValue
		}
		return FalseValue
	case rune:
		return &Char{value: v}
	case byte:
		return &Char{value: rune(v)}
	case float64:
		return &Float{value: v}
	case []byte:
		if len(v) > MaxBytesLen {
			return &Bytes{values: v[0:MaxBytesLen]}
		}
		return &Bytes{values: v}
	case error:
		return &Error{value: &String{value: v.Error()}}
	case map[string]IObject:
		return NewMap(v)
	case map[string]interface{}:
		kv := FromMap(v)
		return NewMap(kv)
	case []bool:
		arr := make([]IObject, len(v))
		for i, e := range v {
			if e {
				arr[i] = TrueValue
			} else {
				arr[i] = FalseValue
			}
		}
		return NewArray(arr)
	case []int:
		arr := make([]IObject, len(v))
		for i, e := range v {
			arr[i] = &Int{value: int64(e)}
		}
		return NewArray(arr)
	case []map[string]interface{}:
		arr := make([]IObject, len(v))
		for i, e := range v {
			kv := FromMap(e)
			vo := FromInterface(kv)
			arr[i] = vo
		}
		return NewArray(arr)
	case []IObject:
		return NewArray(v)
	case []interface{}:
		arr := make([]IObject, len(v))
		for i, e := range v {
			arr[i] = FromInterface(e)
		}
		return NewArray(arr)
	case time.Time:
		return &Time{value: v}
	case IObject:
		return v
	case CallableFunc:
		return NewUserFunction("CallableFunc", v)
	}
	return UndefinedValue
}

// CountObjects recursively counts the total number of objects contained in the given IObject, including nested structures.
func CountObjects(in IObject) (c int) {
	c = 1
	switch o := in.(type) {
	case *Array:
		for _, v := range o.Values() {
			c += CountObjects(v)
		}
	case *ImmutableArray:
		for _, v := range o.Values() {
			c += CountObjects(v)
		}
	case *Map:
		for _, v := range o.values {
			c += CountObjects(v)
		}
	case *ImmutableMap:
		for _, v := range o.Values() {
			c += CountObjects(v)
		}
	case *Error:
		c += CountObjects(o.value)
	}
	return
}
