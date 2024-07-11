package loader

import "fmt"

func buf2dword(buf []byte) (uint32, error) {
	if len(buf) < 4 {
		return 0, fmt.Errorf("invalid length")
	}
	t := buf[:4]
	data := uint32(t[3]) | (uint32(t[2]) << 8) | (uint32(t[1]) << 16) | (uint32(t[0]) << 24)
	return data, nil
}

func buf2word(buf []byte) (uint16, error) {
	if len(buf) < 2 {
		return 0, fmt.Errorf("invalid length")
	}
	t := buf[:2]
	data := uint16(t[1]) | (uint16(t[0]) << 8)
	return data, nil
}
