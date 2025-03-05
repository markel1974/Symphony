package board

import "log"

// IDumpable defines methods to serialize an object's state into a map and restore it from a map representation.
type IDumpable interface {
	Dump(d map[string]interface{}) error
	Restore(d map[string]interface{}) error
}

// DumpSetUint8 sets a uint8 value in the map by converting it to float64 and delegating to DumpSetNumber.
func DumpSetUint8(d map[string]interface{}, id string, val uint8) {
	DumpSetNumber(d, id, float64(val))
}

// DumpSetUint16 sets a uint16 value into a map with the given key, converting it to a float64 internally.
func DumpSetUint16(d map[string]interface{}, id string, val uint16) {
	DumpSetNumber(d, id, float64(val))
}

// DumpSetUint sets a uint value in the given map as a float64, using the provided key.
func DumpSetUint(d map[string]interface{}, id string, val uint) {
	DumpSetNumber(d, id, float64(val))
}

// DumpSetInt sets an integer value in the given map with the specified key by converting it to a float64.
func DumpSetInt(d map[string]interface{}, id string, val int) {
	DumpSetNumber(d, id, float64(val))
}

// DumpSetNumber sets a float64 value in the provided map using the specified key.
func DumpSetNumber(d map[string]interface{}, id string, val float64) {
	d[id] = val
}

// DumpSetBool sets a boolean value in the provided map using the specified key.
func DumpSetBool(d map[string]interface{}, id string, val bool) {
	d[id] = val
}

// DumpSetString sets a string value in the provided map using the given id as the key.
func DumpSetString(d map[string]interface{}, id string, val string) {
	d[id] = val
}

// DumpSetByteArray copies the provided byte slice into a new slice and sets it in the map using the given key.
func DumpSetByteArray(d map[string]interface{}, id string, val []byte) {
	c := make([]uint8, len(val))
	copy(c, val)
	d[id] = c
}

// dumpGet retrieves a value from a map using a key and checks its existence and non-nil status.
// Returns the value and a boolean indicating whether the key was found and the value is non-nil.
func dumpGet(d map[string]interface{}, id string) (interface{}, bool) {
	if d == nil {
		return 0, false
	}
	val, found := d[id]
	if !found || val == nil {
		return nil, false
	}
	return val, true
}

// DumpGetFloat64 retrieves a float64 value from a map by ID, returning true on success and assigning it to the provided pointer.
func DumpGetFloat64(d map[string]interface{}, id string, val *float64) bool {
	data, ok := dumpGet(d, id)
	if !ok {
		return false
	}
	t, ok := data.(float64)
	if !ok {
		log.Printf("DumpGetFloat64: invalid type for key '%s': expected float64, got %T", id, t)
		return false
	}
	*val = t
	return true
}

// DumpGetInt retrieves an integer value from a map by a given key and assigns it to the provided pointer if successful.
// Returns true if the value exists, is of type int, and the assignment is made; otherwise, returns false.
func DumpGetInt(d map[string]interface{}, id string, val *int) bool {
	data, ok := dumpGet(d, id)
	if !ok {
		return false
	}
	t, ok := data.(int)
	if !ok {
		log.Printf("DumpGetInt: invalid type for key '%s': expected int, got %T", id, t)
		return false
	}
	*val = t
	return true
}

// DumpGetUint retrieves a uint8 value from a map using a given key and stores it in the provided pointer if successful.
// Returns true if the key is found and the value can be converted to uint8, otherwise returns false.
func DumpGetUint(d map[string]interface{}, id string, val *uint) bool {
	data, ok := dumpGet(d, id)
	if !ok {
		return false
	}
	t, ok := data.(uint)
	if !ok {
		log.Printf("DumpGetUint8: invalid type for key '%s': expected uint, got %T", id, t)
		return false
	}
	*val = t
	return true
}

// DumpGetUint8 retrieves a uint8 value from a map using a given key and stores it in the provided pointer if successful.
// Returns true if the key is found and the value can be converted to uint8, otherwise returns false.
func DumpGetUint8(d map[string]interface{}, id string, val *uint8) bool {
	data, ok := dumpGet(d, id)
	if !ok {
		return false
	}
	t, ok := data.(uint8)
	if !ok {
		log.Printf("DumpGetUint8: invalid type for key '%s': expected uint8, got %T", id, t)
		return false
	}
	*val = t
	return true
}

// DumpGetUint16 retrieves a uint16 value from a map using the specified key and assigns it to the provided pointer.
// Returns true if the value is successfully retrieved and cast to uint16, otherwise returns false.
func DumpGetUint16(d map[string]interface{}, id string, val *uint16) bool {
	data, ok := dumpGet(d, id)
	if !ok {
		return false
	}
	t, ok := data.(uint16)
	if !ok {
		log.Printf("DumpGetUint16: invalid type for key '%s': expected uint16, got %T", id, t)
		return false
	}
	*val = t
	return true
}

// DumpGetBool retrieves a boolean value from the given map using the specified key and assigns it to the provided pointer.
// Returns true if the key exists, the value is of type bool, and the assignment is successful; otherwise, returns false.
func DumpGetBool(d map[string]interface{}, id string, val *bool) bool {
	t, ok := dumpGet(d, id)
	if !ok {
		return false
	}
	r, ok := t.(bool)
	if !ok {
		log.Printf("DumpGetBool: invalid type for key '%s': expected bool, got %T", id, t)
		return false
	}
	*val = r
	return true
}

// DumpGetString retrieves a string value from a map for a given key and stores it in the provided pointer if successful.
// Returns true if the key exists and the value is a string, otherwise false.
func DumpGetString(d map[string]interface{}, id string, val *string) bool {
	t, ok := dumpGet(d, id)
	if !ok {
		return false
	}
	r, ok := t.(string)
	if !ok {
		log.Printf("DumpGetString: invalid type for key '%s': expected string, got %T", id, t)
		return false
	}
	*val = r
	return true
}

// DumpGetByteArray retrieves a byte array from the specified map identified by `id` and validates its size and content.
func DumpGetByteArray(d map[string]interface{}, id string, res *[]byte, count int) bool {
	val, ok := dumpGet(d, id)
	if !ok {
		return false
	}
	r, ok := val.([]uint8)
	if !ok {
		log.Printf("DumpGetByteArray: invalid type for key '%s': expected []uint8, got %T", id, val)
		return false
	}
	if count > -1 {
		if len(r) != count {
			return false
		}
	}
	if len(r) == 0 {
		return false
	}
	*res = make([]byte, len(r))
	copy(*res, r)
	return true
}
