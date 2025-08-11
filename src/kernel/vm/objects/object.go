package objects

import (
	"time"

	"github.com/markel1974/c64emu/src/kernel/vm/errors"
	"github.com/markel1974/c64emu/src/kernel/vm/tokens"
)

type Object interface {
	// TypeName should return the name of the type.
	TypeName() string

	// String should return a string representation of the type's value.
	String() string

	// BinaryOp should return another object that is the result of a given
	// binary operator and a right-hand side object. If BinaryOp returns an
	// error, the VM will treat it as a run-time error.
	BinaryOp(op tokens.Token, rhs Object) (Object, error)

	// IsFalsy should return true if the value of the type should be considered
	// as falsy.
	IsFalsy() bool

	// Equals should return true if the value of the type should be considered
	// as equal to the value of another object.
	Equals(another Object) bool

	// Copy should return a copy of the type (and its value). Copy function
	// will be used for copy() builtin function which is expected to deep-copy
	// the values generally.
	Copy() Object

	// IndexGet should take an index Object and return a result Object or an
	// error for indexable objects. Indexable is an object that can take an
	// index and return an object. If error is returned, the runtime will treat
	// it as a run-time error and ignore returned value. If Object is not
	// indexable, ErrNotIndexable should be returned as error. If nil is
	// returned as value, it will be converted to UndefinedToken value by the
	// runtime.
	IndexGet(index Object) (value Object, err error)

	// IndexSet should take an index Object and a value Object for index
	// assignable objects. Index assignable is an object that can take an index
	// and a value on the left-hand side of the assignment statement. If Object
	// is not index assignable, ErrNotIndexAssignable should be returned as
	// error. If an error is returned, it will be treated as a run-time error.
	IndexSet(index, value Object) error

	// Iterate should return an Iterator for the type.
	Iterate() Iterator

	// CanIterate should return whether the Object can be Iterated.
	CanIterate() bool

	// Call should take an arbitrary number of arguments and returns a return
	// value and/or an error, which the VM will consider as a run-time error.
	Call(args ...Object) (ret Object, err error)

	// CanCall should return whether the Object can be Called.
	CanCall() bool
}

// Object represents an object in the VM.

// ObjectImpl represents a default Object Implementation. To defined a new
// value type, one can embed ObjectImpl in their type declarations to avoid
// implementing all non-significant methods. TypeName() and String() methods
// still need to be implemented.
type ObjectImpl struct {
}

// TypeName returns the name of the type.
func (o *ObjectImpl) TypeName() string {
	panic(errors.ErrNotImplemented)
}

func (o *ObjectImpl) String() string {
	panic(errors.ErrNotImplemented)
}

// BinaryOp returns another object that is the result of a given binary
// operator and a right-hand side object.
func (o *ObjectImpl) BinaryOp(_ tokens.Token, _ Object) (Object, error) {
	return nil, errors.ErrInvalidOperator
}

// Copy returns a copy of the type.
func (o *ObjectImpl) Copy() Object {
	return nil
}

// IsFalsy returns true if the value of the type is falsy.
func (o *ObjectImpl) IsFalsy() bool {
	return false
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *ObjectImpl) Equals(x Object) bool {
	return o == x
}

// IndexGet returns an element at a given index.
func (o *ObjectImpl) IndexGet(_ Object) (res Object, err error) {
	return nil, errors.ErrNotIndexable
}

// IndexSet sets an element at a given index.
func (o *ObjectImpl) IndexSet(_, _ Object) (err error) {
	return errors.ErrNotIndexAssignable
}

// Iterate returns an iterator.
func (o *ObjectImpl) Iterate() Iterator {
	return nil
}

// CanIterate returns whether the Object can be Iterated.
func (o *ObjectImpl) CanIterate() bool {
	return false
}

// Call takes an arbitrary number of arguments and returns a return value and/or an error.
func (o *ObjectImpl) Call(_ ...Object) (ret Object, err error) {
	return nil, nil
}

// CanCall returns whether the Object can be Called.
func (o *ObjectImpl) CanCall() bool {
	return false
}

// ToInterface attempts to convert an object o to an interface{} value
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

func FromMap(v map[string]interface{}) map[string]Object {
	kv := make(map[string]Object)
	for key, val := range v {
		kv[key] = FromInterface(val)
	}
	return kv
}

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
	//return nil, fmt.Errorf("cannot convert to object: %T", v)
}

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
