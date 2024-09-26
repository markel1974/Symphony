package conversion

// BoolToUint8 converts a bool to a uint8.
// Returns 1 if true, 0 if false.
func BoolToUint8(v bool) uint8 {
	if v {
		return 1
	}
	return 0
}

// Uint8ToBool converts a uint8 to a bool.
// Returns false only if input is 0.
func Uint8ToBool(v uint8) bool {
	return v != 0
}
