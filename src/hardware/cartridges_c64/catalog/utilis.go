package catalog

import "fmt"

// Buf2uint32 converts a 4-byte slice into a uint32 value in big-endian order. Returns an error if the slice is too short.
func Buf2uint32(buf []byte) (uint32, error) {
	if len(buf) < 4 {
		return 0, fmt.Errorf("invalid length")
	}
	t := buf[:4]
	data := uint32(t[3]) | (uint32(t[2]) << 8) | (uint32(t[1]) << 16) | (uint32(t[0]) << 24)
	return data, nil
}

// Buf2uint16 converts a 2-byte buffer into a uint16 while ensuring the buffer has sufficient length. Returns an error if invalid.
func Buf2uint16(buf []byte) (uint16, error) {
	if len(buf) < 2 {
		return 0, fmt.Errorf("invalid length")
	}
	t := buf[:2]
	data := uint16(t[1]) | (uint16(t[0]) << 8)
	return data, nil
}
