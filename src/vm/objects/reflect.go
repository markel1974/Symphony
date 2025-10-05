package objects

import (
	"reflect"
)

func _reflectMap(data map[string]IObject, target reflect.Type) (reflect.Value, bool) {
	//key := target.Key().Kind()
	//val := target.Elem().Kind()
	//fmt.Println(key, val)
	//TODO
	return reflect.Value{}, false
}

func _reflectArray(data []IObject, target reflect.Type) (reflect.Value, bool) {
	elemType := target.Elem()
	name := elemType.Name()
	//z := elemType.Kind().String()
	//fmt.Println(z)
	switch name {
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
