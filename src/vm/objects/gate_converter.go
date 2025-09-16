package objects

import (
	"errors"
	"strconv"
	"time"
)

// GateConverter provides methods to convert between various data types and IObject representations, leveraging GateKeeper.
type GateConverter struct {
	gk *GateKeeper
}

// NewGateConverter creates and returns a new instance of GateConverter with the provided GateKeeper for conversions.
func NewGateConverter(gk *GateKeeper) *GateConverter {
	return &GateConverter{gk: gk}
}

// AssignBool assigns a boolean Code to the destination object by converting it to an integer (1 for true, 0 for false).
func (gc *GateConverter) AssignBool(val bool, dstObj IObject) error {
	v := int64(0)
	if val {
		v = 1
	}
	return gc.AssignInt(v, dstObj)
}

// AssignInt assigns an int64 Code to the destination object (dstObj) if it is a supported type.
// It handles multiple destination types including Bool, Int, Float, Char, and String.
// Returns an error if the destination type is not supported.
func (gc *GateConverter) AssignInt(val int64, dstObj IObject) error {
	switch out := dstObj.(type) {
	case *Bool:
		if val == 0 {
			out.SetValue(false)
		} else {
			out.SetValue(true)
		}
		return nil
	case *Int:
		out.SetValue(val)
		return nil
	case *Float:
		out.SetValue(float64(val))
		return nil
	case *Char:
		out.SetValue(rune(val))
		return nil
	case *String:
		return out.SetValue(strconv.FormatInt(val, 10))
	default:
		return ErrNotAssignable
	}
}

// FromStringArray converts an input slice of strings into an array of IObject and wraps it in a new IObject array.
func (gc *GateConverter) FromStringArray(frame int, in []string) (IObject, error) {
	var data []IObject
	if len(in) > 0 {
		data = make([]IObject, len(in))
		for idx, v := range in {
			r := gc.gk.NewString(frame, v)
			data[idx] = r
		}
	}
	return gc.gk.NewArray(frame, data), nil
}

// FromInterface converts an interface Code into an IObject instance based on its underlying type using the provided frame.
func (gc *GateConverter) FromInterface(frame int, in interface{}) IObject {
	switch v := in.(type) {
	case nil:
		return gc.gk.UndefinedValue()
	case string:
		return gc.gk.NewString(frame, v)
	case uint:
		return gc.gk.NewInt(frame, int64(v))
	case uint16:
		return gc.gk.NewInt(frame, int64(v))
	case uint32:
		return gc.gk.NewInt(frame, int64(v))
	case uint64:
		return gc.gk.NewInt(frame, int64(v))
	case int:
		return gc.gk.NewInt(frame, int64(v))
	case int16:
		return gc.gk.NewInt(frame, int64(v))
	case int64:
		return gc.gk.NewInt(frame, v)
	case bool:
		if v {
			return gc.gk.TrueValue()
		}
		return gc.gk.FalseValue()
	case rune:
		return gc.gk.NewChar(frame, v)
	case byte:
		return gc.gk.NewChar(frame, rune(v))
	case float64:
		return gc.gk.NewFloat(frame, v)
	case float32:
		return gc.gk.NewFloat(frame, float64(v))
	case []byte:
		return gc.gk.NewBytes(frame, v)
	case error:
		return gc.gk.NewError(frame, v.Error())
	case map[string]IObject:
		return gc.gk.NewMap(frame, v)
	case map[string]interface{}:
		kv := gc.FromMap(frame, v)
		return gc.gk.NewMap(frame, kv)
	case []bool:
		arr := make([]IObject, len(v))
		for i, e := range v {
			if e {
				arr[i] = gc.gk.TrueValue()
			} else {
				arr[i] = gc.gk.FalseValue()
			}
		}
		return gc.gk.NewArray(frame, arr)
	case []int:
		arr := make([]IObject, len(v))
		for i, e := range v {
			arr[i] = gc.gk.NewInt(frame, int64(e))
		}
		return gc.gk.NewArray(frame, arr)
	case []map[string]interface{}:
		arr := make([]IObject, len(v))
		for i, e := range v {
			kv := gc.FromMap(frame, e)
			vo := gc.FromInterface(frame, kv)
			arr[i] = vo
		}
		return gc.gk.NewArray(frame, arr)
	case []IObject:
		return gc.gk.NewArray(frame, v)
	case []interface{}:
		arr := make([]IObject, len(v))
		for i, e := range v {
			arr[i] = gc.FromInterface(frame, e)
		}
		return gc.gk.NewArray(frame, arr)
	case time.Time:
		return gc.gk.NewTime(frame, v)
	case IObject:
		return v
	case FuncCallable:
		return gc.gk.NewFuncImport(FrameStatic, "FuncCallable", 0, v)
	}
	return gc.gk.UndefinedValue()
}

// ToMap converts an IObject to a map[string]interface{}, recursively transforming nested elements.
func (gc *GateConverter) ToMap(in IObject) (res map[string]interface{}) {
	switch o := in.(type) {
	case *Map:
		res = make(map[string]interface{})
		for key, v := range o.data {
			res[key] = gc.ToInterface(v)
		}
	}
	return
}

// ToInterface converts an IObject to its native Go interface representation, based on the underlying type of the object.
func (gc *GateConverter) ToInterface(in IObject) (res interface{}) {
	switch o := in.(type) {
	case *Int:
		res = o.data
	case *String:
		res = o.data
	case *Float:
		res = o.data
	case *Bool:
		res = o == gc.gk.TrueValue()
	case *Char:
		res = o.data
	case *Bytes:
		res = o.data
	case *Array:
		res = make([]interface{}, len(o.Values()))
		for i, val := range o.Values() {
			res.([]interface{})[i] = gc.ToInterface(val)
		}
	case *Map:
		res = make(map[string]interface{})
		for key, v := range o.data {
			res.(map[string]interface{})[key] = gc.ToInterface(v)
		}
	case *Time:
		res = o.data
	case *Error:
		res = errors.New(o.AsString())
	case *Undefined:
		res = nil
	case IObject:
		return o
	}
	return
}

// FromMap converts a map of string keys and interface Code into a map of string keys and IObject Code.
func (gc *GateConverter) FromMap(frame int, v map[string]interface{}) map[string]IObject {
	kv := make(map[string]IObject)
	for key, val := range v {
		kv[key] = gc.FromInterface(frame, val)
	}
	return kv
}

// ToInt64 converts an IObject to an int64 and returns the Code along with a boolean indicating success or failure.
func (gc *GateConverter) ToInt64(in IObject) (int64, bool) {
	switch o := in.(type) {
	case *Int:
		return o.data, true
	case *Float:
		return int64(o.data), true
	case *Char:
		return int64(o.data), true
	case *Bool:
		if o == gc.gk.TrueValue() {
			return 1, true
		}
		return 0, true
	case *String:
		c, err := strconv.ParseInt(o.data, 10, 64)
		if err == nil {
			return c, true
		}
		return 0, false
	default:
		return 0, false
	}
}

// ToInt64Arg extracts an int64 Code from the input slice at the specified index or returns an error if invalid.
func (gc *GateConverter) ToInt64Arg(index int, in []IObject) (int64, error) {
	if index < 0 || index >= len(in) {
		return 0, ErrInvalidArgumentsNumber
	}
	o := in[index]
	v, ok := gc.ToInt64(o)
	if !ok {
		return 0, NewInvalidArgumentError(index, "int", o.TypeName())
	}
	return v, nil
}

// ToRune converts the given IObject to a rune if possible, returning the result and a success flag.
func (gc *GateConverter) ToRune(in IObject) (v rune, ok bool) {
	switch o := in.(type) {
	case *Int:
		v = rune(o.data)
		ok = true
	case *Char:
		v = o.data
		ok = true
	}
	return
}

// ToString converts the given IObject to its string representation. It returns the string Code and a boolean for success.
func (gc *GateConverter) ToString(in IObject) (string, bool) {
	if in == nil {
		return "", false
	}
	return in.AsString(), true
}

// ToStringArg converts an IObject at the given index in the slice to a string. Returns an error if conversion fails.
func (gc *GateConverter) ToStringArg(index int, in []IObject) (string, error) {
	if index < 0 || index >= len(in) {
		return "", ErrInvalidArgumentsNumber
	}
	o := in[index]
	v, ok := gc.ToString(o)
	if !ok {
		return "", NewInvalidArgumentError(index, "string", o.TypeName())
	}
	return v, nil
}

// ToBytes converts the provided IObject to a byte slice. Returns the byte slice and true if successful, otherwise nil and false.
func (gc *GateConverter) ToBytes(in IObject) ([]byte, bool) {
	switch o := in.(type) {
	case *Int:
		return []byte{byte(o.data)}, true
	case *Float:
		return []byte{byte(o.data)}, true
	case *Char:
		return []byte{byte(o.data)}, true
	case *Bool:
		if o == gc.gk.TrueValue() {
			return []byte{1}, true
		} else {
			return []byte{0}, true
		}
	case *Bytes:
		return o.data, true
	case *String:
		return []byte(o.data), true
	default:
		return nil, false
	}
}

// ToBytesArg converts the IObject at the specified index in the input slice to a byte array or returns an error if invalid.
func (gc *GateConverter) ToBytesArg(index int, in []IObject) ([]byte, error) {
	if index < 0 || index >= len(in) {
		return nil, ErrInvalidArgumentsNumber
	}
	o := in[index]
	b, ok := gc.ToBytes(o)
	if !ok {
		return nil, NewInvalidArgumentError(index, "bytes", o.TypeName())
	}
	return b, nil
}

// ToFloat64 converts an IObject to a float64, returning the Code and a boolean indicating success or failure.
func (gc *GateConverter) ToFloat64(in IObject) (float64, bool) {
	switch o := in.(type) {
	case *Int:
		return float64(o.data), true
	case *Float:
		return o.data, true
	case *Char:
		return float64(o.data), true
	case *Bool:
		if o == gc.gk.TrueValue() {
			return 1, true
		}
		return 0, true
	case *String:
		c, err := strconv.ParseFloat(o.data, 64)
		if err == nil {
			return c, true
		}
		return 0, false
	default:
		return 0, false
	}
}

// ToFloat64Arg converts an argument at the specified index from an []IObject to a float64 or returns an error if conversion fails.
func (gc *GateConverter) ToFloat64Arg(index int, in []IObject) (float64, error) {
	if index < 0 || index >= len(in) {
		return 0, ErrInvalidArgumentsNumber
	}
	o := in[index]
	v, ok := gc.ToFloat64(o)
	if !ok {
		return 0, NewInvalidArgumentError(index, "float64", o.TypeName())
	}
	return v, nil
}

// ToTime converts an IObject to a time.Time Code if the conversion is possible. Returns the time and a bool indicating success.
func (gc *GateConverter) ToTime(in IObject) (time.Time, bool) {
	switch o := in.(type) {
	case *Time:
		return o.data, true
	case *Int:
		return time.Unix(o.data, 0), true
	case *Float:
		return time.Unix(int64(o.data), 0), true
	case *Char:
		return time.Unix(int64(o.data), 0), true
	case *String:
		if v, err := strconv.ParseInt(o.data, 10, 64); err == nil {
			return time.Unix(v, 0), true
		}
	}
	return time.Time{}, false
}

// ToTimeArg extracts a time.Time Code from the IObject slice at the specified index, returning an error if invalid.
func (gc *GateConverter) ToTimeArg(index int, in []IObject) (time.Time, error) {
	if index < 0 || index >= len(in) {
		return time.Time{}, ErrInvalidArgumentsNumber
	}
	o := in[index]
	v, ok := gc.ToTime(o)
	if !ok {
		return time.Time{}, NewInvalidArgumentError(index, "time", o.TypeName())
	}
	return v, nil
}

// ToBool converts the provided IObject to a boolean Code. Returns the computed boolean and a success flag.
func (gc *GateConverter) ToBool(o IObject) (v bool, ok bool) {
	ok = true
	v = !o.Falsy()
	return
}

// FromBool converts a boolean Code into an IObject representation based on the GateKeeper's true and false Code.
func (gc *GateConverter) FromBool(v bool) IObject {
	if v {
		return gc.gk.TrueValue()
	}
	return gc.gk.FalseValue()
}

// ToBoolArg extracts a boolean Code from the IObject array at the specified index or returns an error if invalid.
func (gc *GateConverter) ToBoolArg(index int, in []IObject) (bool, error) {
	if index < 0 || index >= len(in) {
		return false, ErrInvalidArgumentsNumber
	}
	o := in[index]
	b1, ok := o.(*Bool)
	if !ok {
		return false, NewInvalidArgumentError(index, "bool", o.TypeName())
	}
	return b1.data, nil
}
