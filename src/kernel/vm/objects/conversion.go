package objects

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

// ToInterface converts an IObject to its corresponding native Go representation, such as int, string, float64, bool, etc.
func toInterface(factory *Factory, in IObject) (res interface{}) {
	switch o := in.(type) {
	case *Int:
		res = o.value
	case *String:
		res = o.value
	case *Float:
		res = o.value
	case *Bool:
		res = o == factory.TrueValue()
	case *Char:
		res = o.value
	case *Bytes:
		res = o.values
	case *Array:
		res = make([]interface{}, len(o.Values()))
		for i, val := range o.Values() {
			res.([]interface{})[i] = toInterface(factory, val)
		}
	case *ArrayImmutable:
		res = make([]interface{}, o.Length())
		for i, val := range o.Values() {
			res.([]interface{})[i] = toInterface(factory, val)
		}
	case *Map:
		res = make(map[string]interface{})
		for key, v := range o.values {
			res.(map[string]interface{})[key] = toInterface(factory, v)
		}
	case *MapImmutable:
		res = make(map[string]interface{})
		for key, v := range o.Values() {
			res.(map[string]interface{})[key] = toInterface(factory, v)
		}
	case *Time:
		res = o.value
	case *Error:
		res = errors.New(o.String())
	case *Undefined:
		res = nil
	case IObject:
		return o
	}
	return
}

// FromInterface converts a native Go value of various types into a corresponding IObject implementation.
func fromInterface(factory *Factory, frame int, in interface{}) IObject {
	switch v := in.(type) {
	case nil:
		return factory.UndefinedValue()
	case string:
		if len(v) > MaxStringLen {
			return factory.NewStringNoSize(frame, v[0:MaxStringLen])
		}
		return factory.NewStringNoSize(frame, v)
	case int64:
		return factory.NewInt(frame, v)
	case int:
		return factory.NewInt(frame, int64(v))
	case bool:
		if v {
			return factory.TrueValue()
		}
		return factory.FalseValue()
	case rune:
		return factory.NewChar(frame, v)
	case byte:
		return factory.NewChar(frame, rune(v))
	case float64:
		return factory.NewFloat(frame, v)
	case []byte:
		if len(v) > MaxBytesLen {
			return factory.NewBytes(frame, v[0:MaxBytesLen])
		}
		return factory.NewBytes(frame, v)
	case error:
		return factory.NewError(frame, v.Error())
	case map[string]IObject:
		return factory.NewMap(frame, v)
	case map[string]interface{}:
		kv := factory.FromMap(frame, v)
		return factory.NewMap(frame, kv)
	case []bool:
		arr := make([]IObject, len(v))
		for i, e := range v {
			if e {
				arr[i] = factory.TrueValue()
			} else {
				arr[i] = factory.FalseValue()
			}
		}
		return factory.NewArray(frame, arr)
	case []int:
		arr := make([]IObject, len(v))
		for i, e := range v {
			arr[i] = factory.NewInt(frame, int64(e))
		}
		return factory.NewArray(frame, arr)
	case []map[string]interface{}:
		arr := make([]IObject, len(v))
		for i, e := range v {
			kv := factory.FromMap(frame, e)
			vo := fromInterface(factory, frame, kv)
			arr[i] = vo
		}
		return factory.NewArray(frame, arr)
	case []IObject:
		return factory.NewArray(frame, v)
	case []interface{}:
		arr := make([]IObject, len(v))
		for i, e := range v {
			arr[i] = fromInterface(factory, frame, e)
		}
		return factory.NewArray(frame, arr)
	case time.Time:
		return factory.NewTime(frame, v)
	case IObject:
		return v
	case FuncCallable:
		return factory.NewFuncPackage(FuncPackageDef, "FuncCallable", v)
	}
	return factory.UndefinedValue()
}

// ToMap converts an IObject to a map[string]interface{} if the object is a *Map, recursively applying ToInterface.
func toMap(factory *Factory, o IObject) (res map[string]interface{}) {
	switch o := o.(type) {
	case *Map:
		res = make(map[string]interface{})
		for key, v := range o.values {
			res[key] = toInterface(factory, v)
		}
	}
	return
}

// FromMap converts a map with string keys and interface{} values into a map with string keys and IObject values.
func fromMap(factory *Factory, frame int, v map[string]interface{}) map[string]IObject {
	kv := make(map[string]IObject)
	for key, val := range v {
		kv[key] = fromInterface(factory, frame, val)
	}
	return kv
}

// ToInt64 attempts to convert the given IObject to an int64 value.
// It returns the converted value and a boolean indicating success or failure.
func toInt64(factory *Factory, o IObject) (int64, bool) {
	switch o := o.(type) {
	case *Int:
		return o.value, true
	case *Float:
		return int64(o.value), true
	case *Char:
		return int64(o.value), true
	case *Bool:
		if o == factory.TrueValue() {
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
func toInt64Arg(factory *Factory, index int, o IObject) (int64, error) {
	v, ok := toInt64(factory, o)
	if !ok {
		return 0, NewInvalidArgumentError(index, "int(compatible)", o.TypeName())
	}
	return v, nil
}

// ToRune attempts to convert an IObject to a rune if it is of type *Int or *Char, returning the rune and a boolean success flag.
func toRune(o IObject) (v rune, ok bool) {
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
func toString(factory *Factory, o IObject) (string, bool) {
	if o == nil {
		return "", false
	}
	if o == factory.UndefinedValue() {
		return "", false
	}
	if str, isStr := o.(*String); isStr {
		return str.value, true
	}
	return o.String(), true
}

// ToStringArg attempts to convert an IObject to a string. Returns an error if conversion fails or type is incompatible.
func toStringArg(factory *Factory, index int, o IObject) (string, error) {
	v, ok := toString(factory, o)
	if !ok {
		return "", NewInvalidArgumentError(index, "string(compatible)", o.TypeName())
	}
	return v, nil
}

// ToStringArrayArg attempts to convert an array of IObjects to a slice of strings.
func toStringArrayArg(factory *Factory, index int, arr []IObject) ([]string, error) {
	var sArr []string
	for idx, elem := range arr {
		str, ok := toString(factory, elem)
		if !ok {
			return nil, NewInvalidArgumentError(index, fmt.Sprintf("%d - string array(compatible)", idx), elem.TypeName())
		}
		sArr = append(sArr, str)
	}
	return sArr, nil
}

// ToByteSlice converts an IObject to a byte slice if the object is of type *Bytes or *String.
// It returns the converted byte slice and a boolean indicating success.
func toByteSlice(o IObject) ([]byte, bool) {
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
func toByteSliceArg(factory *Factory, index int, o IObject) ([]byte, error) {
	b, ok := toByteSlice(o)
	if !ok {
		return nil, NewInvalidArgumentError(index, "byte slice(compatible)", o.TypeName())
	}
	return b, nil
}

// ToFloat64 attempts to convert an IObject to a float64 and returns the values along with a success flag.
func toFloat64(factory *Factory, o IObject) (float64, bool) {
	switch o := o.(type) {
	case *Int:
		return float64(o.value), true
	case *Float:
		return o.value, true
	case *Char:
		return float64(o.value), true
	case *Bool:
		if o == factory.TrueValue() {
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
func toFloat64Arg(factory *Factory, index int, o IObject) (float64, error) {
	v, ok := toFloat64(factory, o)
	if !ok {
		return 0, NewInvalidArgumentError(index, "float64(compatible)", o.TypeName())
	}
	return v, nil
}

// ToTime converts an IObject into a time.Time if it is time-compatible (e.g., *Time or *Int). Returns the time and a boolean.
func toTime(o IObject) (time.Time, bool) {
	switch o := o.(type) {
	case *Time:
		return o.value, true
	case *Int:
		return time.Unix(o.value, 0), true
	}
	return time.Time{}, false
}

// ToTimeArg attempts to convert an IObject to a time.Time. Returns an error if the conversion fails or the type is incompatible.
func toTimeArg(index int, o IObject) (time.Time, error) {
	v, ok := toTime(o)
	if !ok {
		return time.Time{}, NewInvalidArgumentError(index, "time(compatible)", o.TypeName())
	}
	return v, nil
}

// ToBool converts the given IObject to a bool based on its Boolean() method and returns the result along with a success flag.
func toBool(o IObject) (v bool, ok bool) {
	ok = true
	v = !o.Boolean()
	return
}

// FromBool converts a boolean values into its corresponding IObject representation, returning TrueValue or FalseValue.
func fromBool(factory *Factory, v bool) IObject {
	if v {
		return factory.TrueValue()
	}
	return factory.FalseValue()
}

// ToBoolArg converts the given IObject to a boolean if possible or returns an error indicating an invalid argument type.
func toBoolArg(index int, o IObject) (bool, error) {
	b1, ok := o.(*Bool)
	if !ok {
		return false, NewInvalidArgumentError(index, "bool(compatible)", o.TypeName())
	}
	return b1.value, nil
}

// FromStringArray converts a slice of strings into an array of IObjects.
func fromStringArray(factory *Factory, frame int, in []string) (IObject, error) {
	var data []IObject
	if len(in) > 0 {
		data = make([]IObject, len(in))
		for idx, v := range in {
			r, err := factory.NewString(frame, v)
			if err != nil {
				return nil, err
			}
			data[idx] = r
		}
	}
	return factory.NewArray(frame, data), nil
}
