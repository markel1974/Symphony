package prg

import (
	"bytes"
	"fmt"
	"github.com/markel1974/c64emu/src/c64/banks"
	"github.com/markel1974/c64emu/src/c64/inputs"
	"log"
	"os"
)

type PRG struct {
	data      []byte
	startAddr uint16
	observer  *banks.Observer
	keys      *inputs.Keyboard
}

func NewPRG(observer *banks.Observer, keys *inputs.Keyboard) *PRG {
	return &PRG{
		observer:  observer,
		keys:      keys,
		data:      nil,
		startAddr: 0,
	}
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
	b.startAddr = ((uint16(src[1])) << 8) | uint16(src[0])
	size := uint32(len(b.data))
	// check range
	if end := uint32(b.startAddr) + (size - 1); end > 0xffff {
		return fmt.Errorf("invalid PRG size")
	}
	return nil
}

func (b *PRG) Inject(buffer []byte) bool {
	if !bytes.Contains(buffer, []byte("READY")) {
		return false
	}
	err := b.observer.Inject(false, b.startAddr, b.data)
	if err != nil {
		log.Printf("can't load prg: %s", err.Error())
	} else {
		b.keys.SetCommand("RUN\n")
	}
	return true
}
