package component

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// GetSegmentKeys retrieves all the keys from a map if the provided interface is a map[string]interface{}.
// Returns an error if the input is nil or not a valid map.
func GetSegmentKeys(s interface{}) ([]string, error) {
	if s == nil {
		return nil, fmt.Errorf("error getting segment keys: %s", "nil interface")
	}
	data, ok := s.(map[string]interface{})
	if !ok || data == nil {
		return nil, fmt.Errorf("error getting segment keys: %s", "invalid object")
	}
	var out []string
	for k := range data {
		out = append(out, k)
	}
	return out, nil
}

// GetSegment retrieves a map segment identified by the given id from the provided interface.
// Returns an error if the interface is nil, invalid, or the segment does not exist.
func GetSegment(id string, s interface{}) (map[string]interface{}, error) {
	if s == nil {
		return nil, fmt.Errorf("error getting segment %s: %s", id, "nil interface")
	}
	data, ok := s.(map[string]interface{})
	if !ok || data == nil {
		return nil, fmt.Errorf("error getting segment %s: %s", id, "invalid interface")
	}
	segmentI, ok := data[id]
	if !ok {
		return nil, fmt.Errorf("error getting segment %s: %s", id, "missing property")
	}
	segment, ok := segmentI.(map[string]interface{})
	if !ok || segment == nil {
		return nil, fmt.Errorf("error getting segment %s: %s", id, "invalid segment interface")
	}
	return segment, nil
}

// ComponentData parses a colon-separated string into label, id, instance, and validates the input format.
func ComponentData(data string) (string, string, int, error) {
	p := strings.Split(data, ":")
	if len(p) < 3 {
		return "", "", 0, fmt.Errorf("error restoring component %s: %s", data, "invalid component id")
	}
	label := p[0]
	id := p[1]
	instance, err := strconv.Atoi(p[2])
	if err != nil {
		return "", "", 0, fmt.Errorf("error restoring component %s: %s", id, err.Error())
	}
	return label, id, instance, nil
}

func ConvertArgument(arg interface{}, targetType reflect.Type) (reflect.Value, bool) {
	argValue := reflect.ValueOf(arg)
	if argValue.Type().AssignableTo(targetType) {
		return argValue, true
	}
	if argValue.Type().ConvertibleTo(targetType) {
		return argValue.Convert(targetType), true
	}
	switch argValue.Kind() {
	case reflect.String:
		if v, err := ConvertStringArgument(argValue.Interface().(string), targetType); err == nil {
			return v, true
		}
	default:
	}
	return reflect.Value{}, false
}

// ConvertStringArgument converts a string argument to a reflect.Value of the specified targetType or returns an error.
func ConvertStringArgument(arg string, targetType reflect.Type) (reflect.Value, error) {
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
