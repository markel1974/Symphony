package objects

import (
	"fmt"
	"strconv"
	"time"
)

// IObject defines a generic interface for objects that can perform various operations and support multiple behaviors.
// TypeName returns the type name of the object.
// String returns the string representation of the object.
// BinaryOp performs a binary operation between the object and a right-hand side operand.
// Boolean checks if the object represents a falsy value.
// Equals checks whether the object is equal to another object.
// Copy creates and returns a copy of the object.
// IndexGet retrieves the values at the specified index from the object.
// IndexSet assigns a value to the specified index within the object.
// Iterate returns an iterator for the object, enabling iteration.
// CanIterate checks if the object can be iterated over.
// Call invokes the object as a callable function with provided arguments.
// CanCall checks if the object can be called as a function.
// Length returns the length of the object.
type IObject interface {
	TypeName() string

	String() string

	BinaryOp(op Operator, rightHandSide IObject) (IObject, error)

	Boolean() bool

	Equals(other IObject) bool

	Copy() IObject

	IndexGet(index IObject) (value IObject, err error)

	IndexSet(index, value IObject) error

	Iterate() IIterator

	CanIterate() bool

	Call(args ...IObject) (ret IObject, err error)

	CanCall() bool

	Length() int
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
	case *ArrayImmutable:
		res = make([]interface{}, o.Length())
		for i, val := range o.Values() {
			res.([]interface{})[i] = ToInterface(val)
		}
	case *Map:
		res = make(map[string]interface{})
		for key, v := range o.values {
			res.(map[string]interface{})[key] = ToInterface(v)
		}
	case *MapImmutable:
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

// ToMap converts an IObject to a map[string]interface{} if the object is a *Map, recursively applying ToInterface.
func ToMap(o IObject) (res map[string]interface{}) {
	switch o := o.(type) {
	case *Map:
		res = make(map[string]interface{})
		for key, v := range o.values {
			res[key] = ToInterface(v)
		}
	}
	return
}

// FromInterface converts a native Go value of various types into a corresponding IObject implementation.
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
		return NewFunctionUser("CallableFunc", v)
	}
	return UndefinedValue
}

// FromMap converts a map with string keys and interface{} values into a map with string keys and IObject values.
func FromMap(v map[string]interface{}) map[string]IObject {
	kv := make(map[string]IObject)
	for key, val := range v {
		kv[key] = FromInterface(val)
	}
	return kv
}

// ToInt64 attempts to convert the given IObject to an int64 value.
// It returns the converted value and a boolean indicating success or failure.
func ToInt64(o IObject) (int64, bool) {
	switch o := o.(type) {
	case *Int:
		return o.value, true
	case *Float:
		return int64(o.value), true
	case *Char:
		return int64(o.value), true
	case *Bool:
		if o == TrueValue {
			return 1, true
		}
		return 0, true
	case *String:
		c, err := strconv.ParseInt(o.value, 10, 64)
		if err == nil {
			return c, true
		}
		return 0, false
	default:
		return 0, false
	}
}

// ToInt64Arg converts an IObject to an int64, returning an error if the conversion is not possible or the type is invalid.
func ToInt64Arg(index int, o IObject) (int64, error) {
	v, ok := ToInt64(o)
	if !ok {
		return 0, NewInvalidArgumentError(index, "int(compatible)", o.TypeName())
	}
	return v, nil
}

// ToRune attempts to convert an IObject to a rune if it is of type *Int or *Char, returning the rune and a boolean success flag.
func ToRune(o IObject) (v rune, ok bool) {
	switch o := o.(type) {
	case *Int:
		v = rune(o.value)
		ok = true
	case *Char:
		v = o.value
		ok = true
	}
	return
}

// ToString converts an IObject to its string representation and determines whether the conversion is valid.
func ToString(o IObject) (string, bool) {
	if o == nil {
		return "", false
	}
	if o == UndefinedValue {
		return "", false
	}
	if str, isStr := o.(*String); isStr {
		return str.value, true
	}
	return o.String(), true
}

// ToStringArg attempts to convert an IObject to a string. Returns an error if conversion fails or type is incompatible.
func ToStringArg(index int, o IObject) (string, error) {
	v, ok := ToString(o)
	if !ok {
		return "", NewInvalidArgumentError(index, "string(compatible)", o.TypeName())
	}
	return v, nil
}

// ToStringArrayArg attempts to convert an array of IObjects to a slice of strings.
func ToStringArrayArg(index int, arr []IObject) ([]string, error) {
	var sArr []string
	for idx, elem := range arr {
		str, ok := ToString(elem)
		if !ok {
			return nil, NewInvalidArgumentError(index, fmt.Sprintf("%d - string array(compatible)", idx), elem.TypeName())
		}
		sArr = append(sArr, str)
	}
	return sArr, nil
}

// ToByteSlice converts an IObject to a byte slice if the object is of type *Bytes or *String.
// It returns the converted byte slice and a boolean indicating success.
func ToByteSlice(o IObject) ([]byte, bool) {
	switch o := o.(type) {
	case *Bytes:
		return o.values, true
	case *String:
		return []byte(o.value), true
	default:
		return nil, false
	}
}

// ToByteSliceArg attempts to convert an IObject to a byte slice. Returns an error if the conversion fails or the type is incompatible.
func ToByteSliceArg(index int, o IObject) ([]byte, error) {
	b, ok := ToByteSlice(o)
	if !ok {
		return nil, NewInvalidArgumentError(index, "byte slice(compatible)", o.TypeName())
	}
	return b, nil
}

// ToFloat64 attempts to convert an IObject to a float64 and returns the values along with a success flag.
func ToFloat64(o IObject) (float64, bool) {
	switch o := o.(type) {
	case *Int:
		return float64(o.value), true
	case *Float:
		return o.value, true
	case *Char:
		return float64(o.value), true
	case *Bool:
		if o == TrueValue {
			return 1, true
		}
		return 0, true
	case *String:
		c, err := strconv.ParseFloat(o.value, 64)
		if err == nil {
			return c, true
		}
		return 0, false
	default:
		return 0, false
	}
}

// ToFloat64Arg converts an IObject to a float64 and returns an error if the conversion fails or the type is incompatible.
func ToFloat64Arg(index int, o IObject) (float64, error) {
	v, ok := ToFloat64(o)
	if !ok {
		return 0, NewInvalidArgumentError(index, "float64(compatible)", o.TypeName())
	}
	return v, nil
}

// ToTime converts an IObject into a time.Time if it is time-compatible (e.g., *Time or *Int). Returns the time and a boolean.
func ToTime(o IObject) (time.Time, bool) {
	switch o := o.(type) {
	case *Time:
		return o.value, true
	case *Int:
		return time.Unix(o.value, 0), true
	}
	return time.Time{}, false
}

func ToTimeArg(index int, o IObject) (time.Time, error) {
	v, ok := ToTime(o)
	if !ok {
		return time.Time{}, NewInvalidArgumentError(index, "time(compatible)", o.TypeName())
	}
	return v, nil
}

// CountObjects recursively counts the total number of objects contained in the given IObject, including nested structures.
func CountObjects(in IObject) int {
	c := 1
	switch o := in.(type) {
	case *Array:
		for _, v := range o.Values() {
			c += CountObjects(v)
		}
	case *ArrayImmutable:
		for _, v := range o.Values() {
			c += CountObjects(v)
		}
	case *Map:
		for _, v := range o.values {
			c += CountObjects(v)
		}
	case *MapImmutable:
		for _, v := range o.Values() {
			c += CountObjects(v)
		}
	case *Error:
		c += CountObjects(o.value)
	}
	return c
}

// IndexAssign assigns a value to a nested structure, using selectors to determine the target location.
// It navigates through the provided selectors and performs an assignment on the target object at the final index.
// Returns an error if any selector is invalid, the object is not indexable, or the assignment fails.
func IndexAssign(dst IObject, src IObject, selectors []IObject) error {
	numSel := len(selectors)
	for sIdx := numSel - 1; sIdx > 0; sIdx-- {
		next, err := dst.IndexGet(selectors[sIdx])
		if err != nil {
			if Is(err, ErrNotIndexable) {
				return fmt.Errorf("not indexable: %s", dst.TypeName())
			}
			if Is(err, ErrInvalidIndexType) {
				return fmt.Errorf("invalid index type: %s",
					selectors[sIdx].TypeName())
			}
			return err
		}
		dst = next
	}
	if err := dst.IndexSet(selectors[0], src); err != nil {
		if Is(err, ErrNotIndexAssignable) {
			return fmt.Errorf("not index-assignable: %s", dst.TypeName())
		}
		if Is(err, ErrInvalidIndexValueType) {
			return fmt.Errorf("invaid index values type: %s", src.TypeName())
		}
		return err
	}
	return nil
}
