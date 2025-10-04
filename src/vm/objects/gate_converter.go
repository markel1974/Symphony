package objects

import (
	"fmt"
	"reflect"
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

// FromBoolError converts a boolean value and error into an IObject representation, returning an error if provided.
func (gc *GateConverter) FromBoolError(val bool, err error) (IObject, error) {
	if err != nil {
		return gc.gk.UndefinedValue(), err
	}
	if val {
		return gc.gk.TrueValue(), nil
	}
	return gc.gk.FalseValue(), nil
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

// FromArrayInterfaces converts a slice of interface{} into a slice of IObject using the specified frame for context.
func (gc *GateConverter) FromArrayInterfaces(frame int, in []interface{}) []IObject {
	if len(in) == 0 {
		return nil
	}
	argsObj := make([]IObject, len(in))
	for idx, arg := range in {
		argObj := gc.gk.FromInterface(frame, arg)
		argsObj[idx] = argObj
	}
	return argsObj

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

func (gc *GateConverter) ReflectMap(data map[string]IObject, target reflect.Type) (reflect.Value, bool) {
	//TODO
	return reflect.Value{}, false
}

func (gc *GateConverter) ReflectArray(data []IObject, target reflect.Type) (reflect.Value, bool) {
	elemType := target.Elem()
	switch elemType.Name() {
	case "bool":
		out := make([]bool, len(data))
		for i, val := range data {
			out[i] = val.AsBool()
		}
		return reflect.ValueOf(out), true
	case "byte":
		out := make([]byte, len(data))
		for i, val := range data {
			out[i] = byte(val.AsInt64())
		}
		return reflect.ValueOf(out), true
	case "rune":
		out := make([]rune, len(data))
		for i, val := range data {
			out[i] = rune(val.AsInt64())
		}
		return reflect.ValueOf(out), true
	case "int":
		out := make([]int, len(data))
		for i, val := range data {
			out[i] = int(val.AsInt64())
		}
		return reflect.ValueOf(out), true
	case "int8":
		out := make([]int8, len(data))
		for i, val := range data {
			out[i] = int8(val.AsInt64())
		}
		return reflect.ValueOf(out), true
	case "int32":
		out := make([]int32, len(data))
		for i, val := range data {
			out[i] = int32(val.AsInt64())
		}
		return reflect.ValueOf(out), true
	case "int64":
		out := make([]int64, len(data))
		for i, val := range data {
			out[i] = val.AsInt64()
		}
		return reflect.ValueOf(out), true
	case "uint":
		out := make([]uint, len(data))
		for i, val := range data {
			out[i] = uint(val.AsInt64())
		}
		return reflect.ValueOf(out), true
	case "uint8":
		out := make([]uint8, len(data))
		for i, val := range data {
			out[i] = uint8(val.AsInt64())
		}
		return reflect.ValueOf(out), true
	case "uint32":
		out := make([]uint32, len(data))
		for i, val := range data {
			out[i] = uint32(val.AsInt64())
		}
		return reflect.ValueOf(out), true
	case "uint64":
		out := make([]uint64, len(data))
		for i, val := range data {
			out[i] = uint64(val.AsInt64())
		}
		return reflect.ValueOf(out), true
	case "float32":
		out := make([]float32, len(data))
		for i, val := range data {
			out[i] = float32(val.AsFloat64())
		}
		return reflect.ValueOf(out), true
	case "float64":
		out := make([]float64, len(data))
		for i, val := range data {
			out[i] = val.AsFloat64()
		}
		return reflect.ValueOf(out), true
	case "string":
		out := make([]string, len(data))
		for i, val := range data {
			out[i] = val.AsString()
		}
		return reflect.ValueOf(out), true
	case "interface", "interface{}", "any":
		out := make([]interface{}, len(data))
		for i, val := range data {
			out[i] = val.AsInterface()
		}
		return reflect.ValueOf(out), true
	default:
		return reflect.Value{}, false
	}
}

func (gc *GateConverter) Reflect(data IObject, target reflect.Type) (reflect.Value, bool) {
	//elemType := target.Elem()
	//es elemType.Name => int
	switch target.Kind() {
	case reflect.Bool:
		return reflect.ValueOf(data.AsBool()), true
	case reflect.Int:
		return reflect.ValueOf(int(data.AsInt64())), true
	case reflect.Int8:
		return reflect.ValueOf(int8(data.AsInt64())), true
	case reflect.Int32:
		return reflect.ValueOf(int32(data.AsInt64())), true
	case reflect.Int64:
		return reflect.ValueOf(data.AsInt64()), true
	case reflect.Uint:
		return reflect.ValueOf(uint(data.AsInt64())), true
	case reflect.Uint8:
		return reflect.ValueOf(uint8(data.AsInt64())), true
	case reflect.Uint32:
		return reflect.ValueOf(uint32(data.AsInt64())), true
	case reflect.Uint64:
		return reflect.ValueOf(uint64(data.AsInt64())), true
	case reflect.Float32:
		return reflect.ValueOf(float32(data.AsFloat64())), true
	case reflect.Float64:
		return reflect.ValueOf(data.AsFloat64()), true
	case reflect.String:
		return reflect.ValueOf(data.AsString()), true
	case reflect.Interface:
		return reflect.ValueOf(data.AsInterface()), true
	default:
		return reflect.Value{}, false
	}
}

// ConvertArgument attempts to convert a given argument to the specified target type, returning the converted value and success status.
func (gc *GateConverter) ConvertArgument(arg interface{}, targetType reflect.Type) (reflect.Value, bool) {
	argValue := reflect.ValueOf(arg)
	if argValue.Type().AssignableTo(targetType) {
		return argValue, true
	}
	if argValue.Type().ConvertibleTo(targetType) {
		return argValue.Convert(targetType), true
	}
	switch argValue.Kind() {
	case reflect.String:
		if v, err := gc.convertStringArgument(argValue.Interface().(string), targetType); err == nil {
			return v, true
		}
	default:
	}
	return reflect.Value{}, false
}

// convertStringArgument converts a string argument to the specified target type using reflection.
// It supports basic types like string, integers, unsigned integers, floats, and booleans.
// Returns the converted value or an error if conversion fails or the type is unsupported.
func (gc *GateConverter) convertStringArgument(arg string, targetType reflect.Type) (reflect.Value, error) {
	switch targetType.Kind() {
	case reflect.String:
		return reflect.ValueOf(arg).Convert(targetType), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(arg, 10, 64)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("not a valid integer")
		}
		return reflect.ValueOf(i).Convert(targetType), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := strconv.ParseUint(arg, 10, 64)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("not a valid unsigned integer")
		}
		return reflect.ValueOf(u).Convert(targetType), nil
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("not a valid float")
		}
		return reflect.ValueOf(f).Convert(targetType), nil
	case reflect.Bool:
		b, err := strconv.ParseBool(arg)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("not a valid boolean")
		}
		return reflect.ValueOf(b).Convert(targetType), nil
	default:
		return reflect.Value{}, fmt.Errorf("unsupported argument type: %s", targetType.Kind())
	}
}
