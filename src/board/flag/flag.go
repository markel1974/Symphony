package flag

func BoolToUint8(v bool) uint8 {
	if v {
		return 1
	}
	return 0
}

func Uint8ToBool(v uint8) bool {
	if v == 0 {
		return false
	}
	return true
}
