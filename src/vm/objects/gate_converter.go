package objects

import (
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
	case *Time:
		out.SetValue(time.Unix(val, 0))
		return nil
	default:
		return ErrNotAssignable
	}
}

// ToAny converts a given interface{} into an IObject using the specified frame as context.
func (gc *GateConverter) ToAny(frame int, i interface{}) IObject {
	return gc.gk.NewAny(frame, i)
}

// FromAny converts an IObject into its underlying data if it is of type *Any. Returns an error if not assignable.
func (gc *GateConverter) FromAny(obj IObject) (interface{}, error) {
	i, ok := obj.(*Any)
	if !ok {
		return nil, ErrNotAssignable
	}
	return i.data, nil
}

// StructFromMap converts a map of key-value pairs to an IObject representation of a structured data entity.
func (gc *GateConverter) StructFromMap(frame int, name string, v map[string]interface{}) IObject {
	base := gc.createObjectMap(frame, v)
	out := gc.gk.NewStruct(frame, name, base)
	return out
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
	case float32:
		return gc.gk.NewFloat(frame, float64(v))
	case float64:
		return gc.gk.NewFloat(frame, v)
	case error:
		return gc.gk.NewError(frame, v.Error())
	case map[string]IObject:
		return gc.gk.NewMap(frame, v)
	case map[string]interface{}:
		kv := gc.createObjectMap(frame, v)
		return gc.gk.NewMap(frame, kv)
	case []map[string]interface{}:
		arr := make([]IObject, len(v))
		for i, e := range v {
			kv := gc.createObjectMap(frame, e)
			vo := gc.FromInterface(frame, kv)
			arr[i] = vo
		}
		return gc.gk.NewArray(frame, arr)
	case []byte:
		return gc.gk.NewBytes(frame, v)
	case []string:
		arr := make([]IObject, len(v))
		for idx, v := range v {
			r := gc.gk.NewString(frame, v)
			arr[idx] = r
		}
		return gc.gk.NewArray(frame, arr)
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
	case []int8:
		arr := make([]IObject, len(v))
		for i, e := range v {
			arr[i] = gc.gk.NewInt(frame, int64(e))
		}
		return gc.gk.NewArray(frame, arr)
	case []int16:
		arr := make([]IObject, len(v))
		for i, e := range v {
			arr[i] = gc.gk.NewInt(frame, int64(e))
		}
		return gc.gk.NewArray(frame, arr)
	case []int32:
		arr := make([]IObject, len(v))
		for i, e := range v {
			arr[i] = gc.gk.NewInt(frame, int64(e))
		}
		return gc.gk.NewArray(frame, arr)
	case []int64:
		arr := make([]IObject, len(v))
		for i, e := range v {
			arr[i] = gc.gk.NewInt(frame, e)
		}
		return gc.gk.NewArray(frame, arr)
	case []uint:
		arr := make([]IObject, len(v))
		for i, e := range v {
			arr[i] = gc.gk.NewInt(frame, int64(e))
		}
		return gc.gk.NewArray(frame, arr)
	case []uint16:
		arr := make([]IObject, len(v))
		for i, e := range v {
			arr[i] = gc.gk.NewInt(frame, int64(e))
		}
		return gc.gk.NewArray(frame, arr)
	case []uint32:
		arr := make([]IObject, len(v))
		for i, e := range v {
			arr[i] = gc.gk.NewInt(frame, int64(e))
		}
		return gc.gk.NewArray(frame, arr)
	case []uint64:
		arr := make([]IObject, len(v))
		for i, e := range v {
			arr[i] = gc.gk.NewInt(frame, int64(e))
		}
		return gc.gk.NewArray(frame, arr)
	case []float32:
		arr := make([]IObject, len(v))
		for i, e := range v {
			arr[i] = gc.gk.NewFloat(frame, float64(e))
		}
		return gc.gk.NewArray(frame, arr)
	case []float64:
		arr := make([]IObject, len(v))
		for i, e := range v {
			arr[i] = gc.gk.NewFloat(frame, e)
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
	case Invocable:
		return gc.gk.NewFuncImport(FrameStatic, "Invocable", 0, v)
	default:
		return gc.gk.NewAny(frame, v)
	}
}

// CreateObjectMap converts a map of string-interface pairs into a map of string-IObject pairs using the provided frame.
func (gc *GateConverter) createObjectMap(frame int, v map[string]interface{}) map[string]IObject {
	kv := make(map[string]IObject)
	for key, val := range v {
		kv[key] = gc.FromInterface(frame, val)
	}
	return kv
}

// ToInt64Arg extracts an int64 Code from the input slice at the specified index or returns an error if invalid.
func (gc *GateConverter) ToInt64Arg(index int, in []IObject) (int64, error) {
	if index < 0 || index >= len(in) {
		return 0, ErrInvalidArgumentsNumber
	}
	o := in[index]
	v := o.AsInt64()
	return v, nil
}

// ToStringArg converts an IObject at the given index in the slice to a string. Returns an error if conversion fails.
func (gc *GateConverter) ToStringArg(index int, in []IObject) (string, error) {
	if index < 0 || index >= len(in) {
		return "", ErrInvalidArgumentsNumber
	}
	o := in[index]
	return o.AsString(), nil
}

// ToBytesArg converts the IObject at the specified index in the input slice to a byte array or returns an error if invalid.
func (gc *GateConverter) ToBytesArg(index int, in []IObject) ([]byte, error) {
	if index < 0 || index >= len(in) {
		return nil, ErrInvalidArgumentsNumber
	}
	o := in[index]
	b := o.AsBytes()
	return b, nil
}

// ToFloat64Arg converts an argument at the specified index from an []IObject to a float64 or returns an error if conversion fails.
func (gc *GateConverter) ToFloat64Arg(index int, in []IObject) (float64, error) {
	if index < 0 || index >= len(in) {
		return 0, ErrInvalidArgumentsNumber
	}
	o := in[index]
	v := o.AsFloat64()
	return v, nil
}

// ToTime converts an IObject to a time.Time Code if the conversion is possible. Returns the time and a bool indicating success.
func (gc *GateConverter) toTime(in IObject) (time.Time, bool) {
	switch o := in.(type) {
	case *Char:
		return time.Unix(int64(o.data), 0), true
	case *Int:
		return time.Unix(o.data, 0), true
	case *Float:
		return time.Unix(int64(o.data), 0), true
	case *Time:
		return o.data, true
	case *Bytes:
		if v, err := strconv.ParseInt(string(o.data), 10, 64); err == nil {
			return time.Unix(v, 0), true
		}
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
	v, ok := gc.toTime(o)
	if !ok {
		return time.Time{}, NewInvalidArgumentError(index, "time", o.TypeName())
	}
	return v, nil
}

// ToBoolArg extracts a boolean Code from the IObject array at the specified index or returns an error if invalid.
func (gc *GateConverter) ToBoolArg(index int, in []IObject) (bool, error) {
	if index < 0 || index >= len(in) {
		return false, ErrInvalidArgumentsNumber
	}
	o := in[index]
	b1 := o.AsBool()
	return b1, nil
}
