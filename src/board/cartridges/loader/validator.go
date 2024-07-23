package loader

import (
	"errors"
)

func ValidateCartridge(data []byte) error {
	i4 := data[0x4]
	i5 := data[0x5]
	i6 := data[0x6]
	i7 := data[0x7]
	i8 := data[0x8]
	if i4 == 0xc3 && i5 == 0xc2 && i6 == 0xcd && i7 == 0x38 && i8 == 0x30 {
		return nil
		//valid cartridge
	}
	return errors.New("invalid cartridge")
}
