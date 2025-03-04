package board

import "fmt"

// IDumpable defines methods for serializing and deserializing an object's state to and from a map representation.
type IDumpable interface {
	Dump() (map[string]interface{}, error)
	Restore(d map[string]interface{}) error
}

// DumpSetUint8 sets a uint8 value in the provided map as a float64 for the specified key.
func DumpSetUint8(d map[string]interface{}, id string, val uint8) {
	DumpSetNumber(d, id, float64(val))
}

// DumpSetUint16 sets a uint16 value in the given map, converting it to float64 for compatibility.
func DumpSetUint16(d map[string]interface{}, id string, val uint16) {
	DumpSetNumber(d, id, float64(val))
}

// DumpSetUint sets a uint value in the given map by converting it to float64 and associating it with the specified id key.
func DumpSetUint(d map[string]interface{}, id string, val uint) {
	DumpSetNumber(d, id, float64(val))
}

// DumpSetInt sets an integer value in the provided map using a given key by converting the integer to a float64.
func DumpSetInt(d map[string]interface{}, id string, val int) {
	DumpSetNumber(d, id, float64(val))
}

// DumpSetNumber sets the value of a given key in the map to the specified float64 value.
func DumpSetNumber(d map[string]interface{}, id string, val float64) {
	d[id] = val
}

// DumpSetBool sets a boolean value in the provided map with the specified key.
func DumpSetBool(d map[string]interface{}, id string, val bool) {
	d[id] = val
}

// DumpSetString sets a string value in a map with the provided key and value.
func DumpSetString(d map[string]interface{}, id string, val string) {
	d[id] = val
}

// DumpSetByteArray copies a byte slice into a map under a specified key as a new slice of uint8.
func DumpSetByteArray(d map[string]interface{}, id string, val []byte) {
	c := make([]uint8, len(val))
	copy(c, val)
	d[id] = c
}

// dumpGet retrieves a value from a map using a given key and returns an error if the value is not found or is nil.
func dumpGet(d map[string]interface{}, id string) (interface{}, error) {
	if d == nil {
		return 0, fmt.Errorf("%s: nil dump in", id)
	}
	val, found := d[id]
	if !found {
		return 0, fmt.Errorf("%s: missing dump", id)
	}
	if val == nil {
		return 0, fmt.Errorf("%s: nil interface", id)
	}
	return val, nil
}

// dumpGetNumber retrieves the value of a specified key from a map as a float64 or returns an error if type conversion fails.
func dumpGetNumber(d map[string]interface{}, id string) (float64, error) {
	val, err := dumpGet(d, id)
	if err != nil {
		return 0, err
	}
	r, ok := val.(float64)
	if !ok {
		return 0, fmt.Errorf("%s: isn't an number", id)
	}
	return r, nil
}

// DumpGetFloat64 retrieves a float64 value from a map by key and returns an error if the value is not found or not a float64.
func DumpGetFloat64(d map[string]interface{}, id string) (float64, error) {
	return dumpGetNumber(d, id)
}

// DumpMustGetFloat64 retrieves and returns the float64 value associated with the given key from the provided map.
// Panics if the conversion to float64 fails or the key does not exist.
func DumpMustGetFloat64(d map[string]interface{}, id string) float64 {
	v, _ := dumpGetNumber(d, id)
	return v
}

// DumpGetInt retrieves an integer value from the map by the given key and returns an error if the value is not a valid number.
func DumpGetInt(d map[string]interface{}, id string) (int, error) {
	val, err := dumpGetNumber(d, id)
	if err != nil {
		return 0, err
	}
	return int(val), nil
}

// DumpMustGetInt retrieves the value associated with the given key from the map, converts it to an int, and returns it.
// Panics if the key does not exist or the value is not a number.
func DumpMustGetInt(d map[string]interface{}, id string) int {
	val, _ := dumpGetNumber(d, id)
	return int(val)
}

// DumpGetUInt8 retrieves a uint8 value associated with the given id from the provided map.
// It returns an error if the id is missing, the value is nil, or cannot be converted to uint8.
func DumpGetUInt8(d map[string]interface{}, id string) (uint8, error) {
	val, err := dumpGetNumber(d, id)
	if err != nil {
		return 0, err
	}
	return uint8(val), nil
}

// DumpMustGetUInt8 retrieves a uint8 value from a map by key and ignores any errors that might occur during retrieval.
func DumpMustGetUInt8(d map[string]interface{}, id string) uint8 {
	val, _ := DumpGetUInt8(d, id)
	return val
}

// DumpGetUInt16 retrieves a uint16 value corresponding to the specified key from the provided map.
// It returns an error if the key is not found, the value is not a number, or the type conversion fails.
func DumpGetUInt16(d map[string]interface{}, id string) (uint16, error) {
	val, err := dumpGetNumber(d, id)
	if err != nil {
		return 0, err
	}
	return uint16(val), nil
}

// DumpMustGetUInt16 retrieves a uint16 value from the map using the given key or returns a zero value if an error occurs.
func DumpMustGetUInt16(d map[string]interface{}, id string) uint16 {
	val, _ := DumpGetUInt16(d, id)
	return val
}

// DumpGetBool retrieves a boolean value from the provided map using the given id and returns an error if not found or invalid.
func DumpGetBool(d map[string]interface{}, id string) (bool, error) {
	val, err := dumpGet(d, id)
	if err != nil {
		return false, err
	}
	r, ok := val.(bool)
	if !ok {
		return false, fmt.Errorf("%s: isn't an bool", id)
	}
	return r, nil
}

// DumpMustGetBool retrieves a boolean value from the provided map using the given key and defaults to false on error.
func DumpMustGetBool(d map[string]interface{}, id string) bool {
	val, _ := DumpGetBool(d, id)
	return val
}

// DumpGetString retrieves a string value associated with the given id from the map or returns an error if not found or invalid.
func DumpGetString(d map[string]interface{}, id string) (string, error) {
	val, err := dumpGet(d, id)
	if err != nil {
		return "", err
	}
	r, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("%s: isn't an string", id)
	}
	return r, nil
}

// DumpMustGetString retrieves the string value associated with the given id from the map or returns an empty string on error.
func DumpMustGetString(d map[string]interface{}, id string) string {
	val, _ := DumpGetString(d, id)
	return val
}

// DumpGetByteArray retrieves a byte array from a map by key, validates its length if specified, and returns a copy.
// Returns an error if the key is not found, the value is not a byte array, or the length is mismatched.
func DumpGetByteArray(d map[string]interface{}, id string, count int) ([]byte, error) {
	val, err := dumpGet(d, id)
	if err != nil {
		return nil, err
	}
	r, ok := val.([]uint8)
	if !ok {
		return nil, fmt.Errorf("%s: isn't an byte array", id)
	}
	if count > -1 {
		if len(r) != count {
			return nil, fmt.Errorf("%s: array length mismatch", id)
		}
	}
	if len(r) == 0 {
		return []byte{}, nil
	}
	out := make([]byte, len(r))
	copy(out, r)
	return out, nil
}

// DumpMustGetByteArray retrieves a byte array from the map and panics if an error occurs during retrieval or casting.
func DumpMustGetByteArray(d map[string]interface{}, id string, count int) []byte {
	val, _ := DumpGetByteArray(d, id, count)
	return val
}
