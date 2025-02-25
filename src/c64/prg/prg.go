package prg

import (
	"bytes"
	"fmt"
	"github.com/markel1974/c64emu/src/c64/inputs"
	"os"
)

type PRG struct {
	data      []byte
	startAddr uint16
	observer  *Observer
	keys      *inputs.Keyboard
	search    []byte
	command   string
}

func NewPRG(b IAdapter, keys *inputs.Keyboard) *PRG {
	return &PRG{
		observer:  NewObserver(b),
		keys:      keys,
		data:      nil,
		startAddr: 0,
		search:    []byte("READY"),
		command:   "RUN\n",
	}
}

func (b *PRG) SetSearch(search string) {
	b.search = []byte(search)
}

func (b *PRG) SetCommand(cmd string) {
	b.command = cmd
}

func (b *PRG) Load(prgFile string) error {
	src, err := os.ReadFile(prgFile)
	if err != nil {
		return err
	}
	if len(src) < 3 {
		return fmt.Errorf("invalid prg file len")
	}
	b.data = src[2:]
	b.startAddr = ((uint16(src[1])) << 8) | uint16(src[0])
	size := uint32(len(b.data))
	if end := uint32(b.startAddr) + (size - 1); end > 0xffff {
		return fmt.Errorf("invalid prg size")
	}
	return nil
}

func (b *PRG) Inject(buffer []byte) bool {
	if !bytes.Contains(buffer, b.search) {
		return false
	}
	b.observer.Inject(false, b.startAddr, b.data)
	b.keys.SetCommand(b.command)
	return true
}
