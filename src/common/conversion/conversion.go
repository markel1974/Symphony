package conversion

// BoolToUint8 converts a boolean value to its uint8 representation; returns 1 for true and 0 for false.
func BoolToUint8(v bool) uint8 {
	if v {
		return 1
	}
	return 0
}

// Uint8ToBool converts a uint8 value to a boolean, returning true if the value is non-zero, and false otherwise.
func Uint8ToBool(v uint8) bool {
	return v != 0
}
