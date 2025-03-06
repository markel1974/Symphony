package board

import "fmt"

// dumpSetByteArray sets a byte array (slice of uint8) value in the container map at the specified key.
// It validates the input type and returns an error if the type does not match.
func dumpSetByteArray(container map[string]interface{}, key string, t interface{}) error {
	val, ok := t.([]uint8)
	if !ok {
		return fmt.Errorf("invalid type: got %T", t)
	}
	container[key] = val
	return nil
}

// dumpGetByteArray retrieves a value of type []uint8 from the container map, validates it, and copies it into the provided pointer.
// It expects the key for the target value and exactly two arguments: the expected length (int) and a *[]uint8 to store the result.
// Returns an error if the value is nil, arguments are invalid, types mismatch, or the byte array length is incorrect.
func dumpGetByteArray(container map[string]interface{}, key string, args []interface{}) error {
	t := container[key]
	if t == nil {
		return fmt.Errorf("nil value")
	}
	if len(args) != 2 {
		return fmt.Errorf("invalid arguments")
	}
	count, ok := args[0].(int)
	if !ok {
		return fmt.Errorf("invalid type: got %T", t)
	}
	val, ok := args[0].(*[]uint8)
	if !ok {
		return fmt.Errorf("invalid type: got %T", t)
	}
	r, ok := t.([]uint8)
	if !ok {
		return fmt.Errorf("invalid type: got %T", t)
	}
	if count > -1 {
		if len(r) != count {
			return fmt.Errorf("invalid length: got %d, expected %d", len(r), count)
		}
	}
	if len(r) == 0 {
		*val = []byte{}
	} else {
		*val = make([]byte, len(r))
		copy(*val, r)
	}
	return nil
}

// dumpGetUint8 retrieves a uint8 value from the container by key and assigns it to the first argument in args if provided.
// Returns an error if the key does not exist, the arguments are invalid, or the types do not match.
func dumpGetUint8(container map[string]interface{}, key string, args []interface{}) error {
	t := container[key]
	if t == nil {
		return fmt.Errorf("nil value")
	}
	if len(args) != 1 {
		return fmt.Errorf("invalid arguments")
	}
	val, ok := args[0].(*uint8)
	if !ok {
		return fmt.Errorf("invalid type: got %T", t)
	}
	r, ok := t.(uint8)
	if !ok {
		return fmt.Errorf("invalid type: got %T", t)
	}
	*val = r
	return nil
}

// dumpSetUint8 sets a uint8 value to the specified key in the container map.
// Returns an error if the value t is not of type uint8.
func dumpSetUint8(container map[string]interface{}, key string, t interface{}) error {
	val, ok := t.(uint8)
	if !ok {
		return fmt.Errorf("invalid type: got %T", t)
	}
	container[key] = val
	return nil
}

// propertySetFunction defines a function type that sets a value in a map container using a key and a provided value.
// The function returns an error if the operation fails.
type propertySetFunction func(container map[string]interface{}, key string, t interface{}) error

// propertyGetFunction defines a function type for retrieving a property value from a container with given arguments.
type propertyGetFunction func(container map[string]interface{}, key string, args []interface{}) error

// _typeUint8 is a placeholder variable representing a uint8 type with an initial value of 0.
var _typeUint8 = uint8(0)

// _typeUint8Array is a variable of type slice of uint8 used to represent a sequence of unsigned 8-bit integers.
var _typeUint8Array []uint8

// _supportedUInt8 represents the stringified type of a uint8 variable, used for type identification in property handling.
var _supportedUInt8 = fmt.Sprintf("%T", _typeUint8)

// _supportedUInt8Array holds the string representation of the type for a slice of uint8 values ([]uint8).
var _supportedUInt8Array = fmt.Sprintf("%T", _typeUint8Array)

// PropertyInfo represents metadata associated with a property, including its type, description, and access methods.
type PropertyInfo struct {
	kind        string
	description string
	readOnly    bool
	get         propertyGetFunction
	set         propertySetFunction
}

// MustCreatePropertyInfo creates a PropertyInfo and panics if an error occurs during its initialization.
func MustCreatePropertyInfo(kind interface{}, description string, readonly bool) *PropertyInfo {
	v, err := NewPropertyInfo(kind, description, readonly)
	if err != nil {
		panic(err)
	}
	return v
}

// NewPropertyInfo creates a new PropertyInfo instance with the specified type, description, and read-only configuration.
// It returns an error if the specified type is not supported.
func NewPropertyInfo(kind interface{}, desc string, ro bool) (*PropertyInfo, error) {
	t := fmt.Sprintf("%T", kind)
	switch t {
	case _supportedUInt8:
		return &PropertyInfo{kind: t, description: desc, readOnly: ro, get: dumpGetUint8, set: dumpSetUint8}, nil
	case _supportedUInt8Array:
		return &PropertyInfo{kind: t, description: desc, readOnly: ro, get: dumpGetByteArray, set: dumpSetByteArray}, nil
	}
	return nil, fmt.Errorf("unsupported type: got %T", t)
}
