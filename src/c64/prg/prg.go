package prg

import (
	"fmt"
	"os"
)

type PRG struct {
	data      []byte
	startAddr uint16
}

func NewPRG() *PRG {
	return &PRG{
		data:      nil,
		startAddr: 0,
	}
}

func (b *PRG) GetData() []byte {
	return b.data
}

func (b *PRG) GetStartAddress() uint16 {
	return b.startAddr
}

func (b *PRG) Load(prgFile string) error {
	src, err := os.ReadFile(prgFile)
	if err != nil {
		return err
	}
	if len(src) < 3 {
		return fmt.Errorf("invalid PRG file len")
	}
	b.data = src[2:]
	// get load addr
	b.startAddr = uint16(src[1])<<8 | uint16(src[0])
	size := uint32(len(b.data))
	// check range
	if end := uint32(b.startAddr) + (size - 1); end > 0xffff {
		return fmt.Errorf("invalid PRG size")
	}
	return nil
}
