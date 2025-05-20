package prg

import (
	"bytes"
	"fmt"
)

// Keys represents an interface defining a method for setting a command string in an implementing struct.
type Keys interface {
	SetCommand(cmd string)
}

// PRG represents a memory-loaded program with a start address, search text, command, and associated observer and keys.
type PRG struct {
	data      []byte
	startAddr uint16
	observer  *Observer
	keys      Keys
	search    []byte
	command   string
}

// NewPRG creates and initializes a new PRG object with the provided IAdapter and Keys instances.
func NewPRG(b IAdapter, keys Keys) *PRG {
	return &PRG{
		observer:  NewObserver(b),
		keys:      keys,
		data:      nil,
		startAddr: 0,
		search:    []byte("READY"),
		command:   "RUN\n",
	}
}

// SetSearch configures the search value by converting the provided string into a byte slice and storing it in the PRG instance.
func (b *PRG) SetSearch(search string) {
	b.search = []byte(search)
}

// SetCommand assigns the given command string to the `command` field of the PRG instance.
func (b *PRG) SetCommand(cmd string) {
	b.command = cmd
}

// Load reads a PRG file from the specified path and loads its contents into memory, returning an error if it fails.
func (b *PRG) Load(src []byte) error {
	return b.LoadData(src)
}

// LoadData loads PRG data from the provided byte slice and sets the start address. Returns an error on invalid data.
func (b *PRG) LoadData(src []byte) error {
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

// Inject checks if the given buffer contains the search byte sequence.
// If found, it loads data into memory via the observer and sets a command in the keys object.
// Returns true if the operation is successful, otherwise false.
func (b *PRG) Inject(buffer []byte) bool {
	if !bytes.Contains(buffer, b.search) {
		return false
	}
	b.observer.Inject(false, b.startAddr, b.data)
	b.keys.SetCommand(b.command)
	return true
}
