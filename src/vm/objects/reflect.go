package objects

import (
	"reflect"
)

// _reflectMap converts a map of string keys and IObject values into a Go reflect.Value of the specified map type.
// It supports mapping to various types of map values such as bool, int64, float64, string, interface{}, etc.
// Returns the constructed reflect.Value and a boolean indicating success or failure based on target type compatibility.
func _reflectMap(data map[string]IObject, target reflect.Type) (reflect.Value, bool) {
	key := target.Key().Kind()
	val := target.Elem().Kind()
	switch key {
	case reflect.String:
		switch val {
		case reflect.Bool:
			out := make(map[string]bool)
			for k, v := range data {
				out[k] = v.AsBool()
			}
			return reflect.ValueOf(out), true
		case reflect.Int64:
			out := make(map[string]int64)
			for k, v := range data {
				out[k] = v.AsInt64()
			}
			return reflect.ValueOf(out), true
		case reflect.Float64:
			out := make(map[string]float64)
			for k, v := range data {
				out[k] = v.AsFloat64()
			}
			return reflect.ValueOf(out), true
		case reflect.String:
			out := make(map[string]string)
			for k, v := range data {
				out[k] = v.AsString()
			}
			return reflect.ValueOf(out), true
		case reflect.Interface:
			out := make(map[string]interface{})
			for k, v := range data {
				out[k] = v.AsInterface()
			}
			return reflect.ValueOf(out), true
		default:
			return reflect.Value{}, false
		}
	default:
		return reflect.Value{}, false
	}
}

// _reflectArray converts a slice of IObject to a slice of the specified reflect.Type and returns it as reflect.Value.
// It supports various basic types (e.g., bool, int, float64, string), failing for unsupported types with a false return.
func _reflectArray(data []IObject, target reflect.Type) (reflect.Value, bool) {
	elemType := target.Elem()
	//name := elemType.Name()
	kind := elemType.Kind()
	switch kind {
	case reflect.Bool:
		out := make([]bool, len(data))
		for i, val := range data {
			out[i] = val.AsBool()
		}
		return reflect.ValueOf(out), true
	case reflect.Int:
		out := make([]int, len(data))
		for i, val := range data {
			out[i] = int(val.AsInt64())
		}
		return reflect.ValueOf(out), true
	case reflect.Int8:
		out := make([]int8, len(data))
		for i, val := range data {
			out[i] = int8(val.AsInt64())
		}
		return reflect.ValueOf(out), true
	case reflect.Int32:
		out := make([]int32, len(data))
		for i, val := range data {
			out[i] = int32(val.AsInt64())
		}
		return reflect.ValueOf(out), true
	case reflect.Int64:
		out := make([]int64, len(data))
		for i, val := range data {
			out[i] = val.AsInt64()
		}
		return reflect.ValueOf(out), true
	case reflect.Uint:
		out := make([]uint, len(data))
		for i, val := range data {
			out[i] = uint(val.AsInt64())
		}
		return reflect.ValueOf(out), true
	case reflect.Uint8:
		out := make([]uint8, len(data))
		for i, val := range data {
			out[i] = uint8(val.AsInt64())
		}
		return reflect.ValueOf(out), true
	case reflect.Uint32:
		out := make([]uint32, len(data))
		for i, val := range data {
			out[i] = uint32(val.AsInt64())
		}
		return reflect.ValueOf(out), true
	case reflect.Uint64:
		out := make([]uint64, len(data))
		for i, val := range data {
			out[i] = uint64(val.AsInt64())
		}
		return reflect.ValueOf(out), true
	case reflect.Float32:
		out := make([]float32, len(data))
		for i, val := range data {
			out[i] = float32(val.AsFloat64())
		}
		return reflect.ValueOf(out), true
	case reflect.Float64:
		out := make([]float64, len(data))
		for i, val := range data {
			out[i] = val.AsFloat64()
		}
		return reflect.ValueOf(out), true
	case reflect.String:
		out := make([]string, len(data))
		for i, val := range data {
			out[i] = val.AsString()
		}
		return reflect.ValueOf(out), true
	case reflect.Interface:
		out := make([]interface{}, len(data))
		for i, val := range data {
			out[i] = val.AsInterface()
		}
		return reflect.ValueOf(out), true
	default:
		return reflect.Value{}, false
	}
}

// _reflect attempts to convert an IObject to a reflect.Value of the specified target type.
// It returns the converted reflect.Value and a boolean indicating success or failure.
func _reflect(data IObject, target reflect.Type) (reflect.Value, bool) {
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
	case reflect.Array, reflect.Slice:
		t, ok := data.(*Array)
		if !ok {
			return reflect.Value{}, false
		}
		return _reflectArray(t.data, target)
	case reflect.Map:
		t, ok := data.(*Map)
		if !ok {
			return reflect.Value{}, false
		}
		return _reflectMap(t.data, target)
	default:
		return reflect.Value{}, false
	}
}
