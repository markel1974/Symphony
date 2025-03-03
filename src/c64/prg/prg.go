package prg

import (
	"bytes"
	"fmt"
	"github.com/markel1974/c64emu/src/c64/inputs"
	"os"
)

// PRG represents a program loader, managing data injection, search strings, and command execution.
type PRG struct {
	data      []byte
	startAddr uint16
	observer  *Observer
	keys      *inputs.Keyboard
	search    []byte
	command   string
}

// NewPRG creates a new PRG instance, initializing its observer, keyboard, and default properties.
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

// SetSearch updates the search property by converting the provided string into a byte slice and assigning it.
func (b *PRG) SetSearch(search string) {
	b.search = []byte(search)
}

// SetCommand sets the command string for the PRG instance, which is used to configure the keyboard during execution.
func (b *PRG) SetCommand(cmd string) {
	b.command = cmd
}

// Load loads a PRG file into memory, validates its size, and calculates the start address from the file's header.
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

// Inject checks if the search pattern exists in the provided buffer and, if found, injects data into memory.
// It sets the command string and returns true if injection was successful, otherwise returns false.
func (b *PRG) Inject(buffer []byte) bool {
	if !bytes.Contains(buffer, b.search) {
		return false
	}
	b.observer.Inject(false, b.startAddr, b.data)
	b.keys.SetCommand(b.command)
	return true
}
