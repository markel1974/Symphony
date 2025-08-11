package objects

import (
	"time"

	"github.com/markel1974/c64emu/src/kernel/vm/errors"
)

// Object represents the fundamental interface for all types
// TypeName returns the name of the object's type.
// String returns a string representation of the object.
// BinaryOp performs a binary operation using the given operator and another object, returning a result or error.
// Falsy determines if the object is considered falsy in a boolean context.
// Equals compares the current object with another, checking for value equality.
// Copy creates and returns a copy of the object.
// IndexGet retrieves a value at the given index from the object, or an error if indexing is not supported.
// IndexSet assigns a value at the given index within the object, or returns an error if not assignable.
// Iterate returns an iterator for the object or nil if the object is not iterable.
// CanIterate checks whether the object supports iteration.
// Call invokes the object as a callable with arguments, returning the result or error if not callable.
// CanCall checks whether the object can be invoked as a callable.
type Object interface {
	TypeName() string

	String() string

	BinaryOp(op Operator, rhs Object) (Object, error)

	Falsy() bool

	Equals(another Object) bool

	Copy() Object

	IndexGet(index Object) (value Object, err error)

	IndexSet(index, value Object) error

	Iterate() Iterator

	CanIterate() bool

	Call(args ...Object) (ret Object, err error)

	CanCall() bool
}

// Object represents an object in the VM.

// ObjectImpl is a base implementation of the Object interface, providing default behaviors for common methods.
type ObjectImpl struct {
}

// TypeName returns the name of the type as a string. Currently, it raises an error indicating the method is not implemented.
func (o *ObjectImpl) TypeName() string {
	panic(errors.ErrNotImplemented)
}

// String returns a string representation of the ObjectImpl. It currently panics as it is not implemented.
func (o *ObjectImpl) String() string {
	panic(errors.ErrNotImplemented)
}

// BinaryOp performs a binary operation with the given operator and object, returning an error for unsupported operations.
func (o *ObjectImpl) BinaryOp(_ Operator, _ Object) (Object, error) {
	return nil, errors.ErrInvalidOperator
}

// Copy creates and returns a duplicate of the current object.
func (o *ObjectImpl) Copy() Object {
	return nil
}

// IsFalsy returns true if the object should be considered a falsy value, otherwise false.
func (o *ObjectImpl) Falsy() bool {
	return false
}

// Equals checks whether the given Object is equal to the current ObjectImpl instance.
func (o *ObjectImpl) Equals(x Object) bool {
	return o == x
}

// IndexGet attempts to retrieve a value for the provided index but returns an ErrNotIndexable since ObjectImpl is not indexable.
func (o *ObjectImpl) IndexGet(_ Object) (res Object, err error) {
	return nil, errors.ErrNotIndexable
}

// IndexSet sets a value at the specified index but returns an error indicating the object is not index-assignable.
func (o *ObjectImpl) IndexSet(_, _ Object) (err error) {
	return errors.ErrNotIndexAssignable
}

// Iterate returns an Iterator for the object, or nil if the object is not iterable.
func (o *ObjectImpl) Iterate() Iterator {
	return nil
}

// CanIterate determines if the object supports iteration and returns false for ObjectImpl.
func (o *ObjectImpl) CanIterate() bool {
	return false
}

// Call executes the ObjectImpl instance as a callable, passing the provided arguments, and returns the result or an error.
func (o *ObjectImpl) Call(_ ...Object) (ret Object, err error) {
	return nil, nil
}

// CanCall checks if the object can be invoked as a callable. Returns false for non-callable objects.
func (o *ObjectImpl) CanCall() bool {
	return false
}

// ToInterface converts an Object to its corresponding Go native interface representation.
func ToInterface(o Object) (res interface{}) {
	switch o := o.(type) {
	case *Int:
		res = o.Value
	case *String:
		res = o.Value
	case *Float:
		res = o.Value
	case *Bool:
		res = o == TrueValue
	case *Char:
		res = o.Value
	case *Bytes:
		res = o.Value
	case *Array:
		res = make([]interface{}, len(o.Value))
		for i, val := range o.Value {
			res.([]interface{})[i] = ToInterface(val)
		}
	case *ImmutableArray:
		res = make([]interface{}, len(o.Value))
		for i, val := range o.Value {
			res.([]interface{})[i] = ToInterface(val)
		}
	case *Map:
		res = make(map[string]interface{})
		for key, v := range o.Value {
			res.(map[string]interface{})[key] = ToInterface(v)
		}
	case *ImmutableMap:
		res = make(map[string]interface{})
		for key, v := range o.Value {
			res.(map[string]interface{})[key] = ToInterface(v)
		}
	case *Time:
		res = o.Value
	case *Error:
		res = errors.New(o.String())
	case *Undefined:
		res = nil
	case Object:
		return o
	}
	return
}

// FromMap converts a map[string]interface{} to a map[string]Object by recursively transforming its values to Object types.
func FromMap(v map[string]interface{}) map[string]Object {
	kv := make(map[string]Object)
	for key, val := range v {
		kv[key] = FromInterface(val)
	}
	return kv
}

// FromInterface converts a native Go interface{} into a corresponding Object type based on its underlying value.
func FromInterface(v interface{}) Object {
	switch v := v.(type) {
	case nil:
		return UndefinedValue
	case string:
		if len(v) > MaxStringLen {
			return &String{Value: v[0:MaxStringLen]}
		}
		return &String{Value: v}
	case int64:
		return &Int{Value: v}
	case int:
		return &Int{Value: int64(v)}
	case bool:
		if v {
			return TrueValue
		}
		return FalseValue
	case rune:
		return &Char{Value: v}
	case byte:
		return &Char{Value: rune(v)}
	case float64:
		return &Float{Value: v}
	case []byte:
		if len(v) > MaxBytesLen {
			return &Bytes{Value: v[0:MaxBytesLen]}
		}
		return &Bytes{Value: v}
	case error:
		return &Error{Value: &String{Value: v.Error()}}
	case map[string]Object:
		return &Map{Value: v}
	case map[string]interface{}:
		kv := FromMap(v)
		return &Map{Value: kv}
	case []bool:
		arr := make([]Object, len(v))
		for i, e := range v {
			if e {
				arr[i] = TrueValue
			} else {
				arr[i] = FalseValue
			}
		}
		return &Array{Value: arr}
	case []int:
		arr := make([]Object, len(v))
		for i, e := range v {
			arr[i] = &Int{Value: int64(e)}
		}
		return &Array{Value: arr}
	case []map[string]interface{}:
		arr := make([]Object, len(v))
		for i, e := range v {
			kv := FromMap(e)
			vo := FromInterface(kv)
			arr[i] = vo
		}
		return &Array{Value: arr}
	case []Object:
		return &Array{Value: v}
	case []interface{}:
		arr := make([]Object, len(v))
		for i, e := range v {
			arr[i] = FromInterface(e)
		}
		return &Array{Value: arr}
	case time.Time:
		return &Time{Value: v}
	case Object:
		return v
	case CallableFunc:
		return &UserFunction{Value: v}
	}
	return UndefinedValue
}

// CountObjects counts the total number of objects in the given object, including nested objects within supported types.
func CountObjects(o Object) (c int) {
	c = 1
	switch o := o.(type) {
	case *Array:
		for _, v := range o.Value {
			c += CountObjects(v)
		}
	case *ImmutableArray:
		for _, v := range o.Value {
			c += CountObjects(v)
		}
	case *Map:
		for _, v := range o.Value {
			c += CountObjects(v)
		}
	case *ImmutableMap:
		for _, v := range o.Value {
			c += CountObjects(v)
		}
	case *Error:
		c += CountObjects(o.Value)
	}
	return
}
