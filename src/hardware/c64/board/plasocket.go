package board

import (
	"fmt"
	"github.com/markel1974/c64emu/src/component"
	pla2 "github.com/markel1974/c64emu/src/hardware/pla_c64"
	"github.com/markel1974/c64emu/src/hardware/roms_c64"
	"github.com/markel1974/c64emu/src/references"
)

// PLASocket represents a socket encapsulating a PLA and its associated board.
type PLASocket struct {
	board *Board
	pla   *pla2.PLA
}

// NewPLASocket creates a new instance of PLASocket with default uninitialized board and PLA components.
func NewPLASocket() *PLASocket {
	c := &PLASocket{
		board: nil,
		pla:   nil,
	}
	return c
}

// Setup initializes the PLASocket with the provided board, PLA, and associated components like VIC, SID, CIAs, and cartridge manager.
func (w *PLASocket) Setup(board *Board, p *pla2.PLA, vic references.IVic, sid references.ISid, c1 component.IComponent, c2 component.IComponent, cartMan references.IExpansionSocketC64) error {
	cia1, ok := c1.(references.ICia)
	if !ok {
		return fmt.Errorf("unknown cia1 interface")
	}
	cia2, ok := c2.(references.ICia)
	if !ok {
		return fmt.Errorf("unknown cia2 interface")
	}
	w.board = board
	w.pla = p
	rl := roms_c64.NewRomLoader(w.board.cfg)
	w.pla.Setup(vic, sid, cia1, cia2, cartMan, rl, w.board.cfg)
	return nil
}

// Reset reinitializes the internal state of the PLA to its default configuration.
func (w *PLASocket) Reset() {
	w.pla.Reset()
}

// Emulate runs the core emulation cycle for this PLASocket instance, handling interactions and operations as defined.
func (w *PLASocket) Emulate() {
	// Emulate is empty because the PLA does not have its own emulation cycle,
	// it works as a memory controller.
}

// GetMemoryConfig retrieves the current memory configuration as a slice of uint8 values.
func (w *PLASocket) GetMemoryConfig() []uint8 {
	return w.pla.GetMemoryConfig()
}

// SetMemoryEntry updates the current memory entry to the specified configuration value.
func (w *PLASocket) SetMemoryEntry(m uint8) {
	w.pla.SetMemoryEntry(m)
}

// SetMemoryConfig updates the memory configuration of the associated PLA using the provided configuration slice.
func (w *PLASocket) SetMemoryConfig(m []uint8) {
	w.pla.SetMemoryConfig(m)
}

// RebuildMemoryConfig triggers a reconstruction of the memory configuration in the associated PLA.
func (w *PLASocket) RebuildMemoryConfig() {
	w.pla.RebuildMemoryConfig()
}

// Write writes a byte of data to the specified memory address using the PLA's write functionality.
func (w *PLASocket) Write(addr uint16, data uint8) {
	w.pla.Write(addr, data)
}

// Read reads a byte of data from the specified memory address using the underlying PLA logic of the PLASocket.
func (w *PLASocket) Read(addr uint16) uint8 {
	return w.pla.Read(addr)
}

// ReadCharRom reads a character ROM value at the specified address and returns the corresponding byte.
func (w *PLASocket) ReadCharRom(addr uint16) uint8 {
	return w.pla.ReadCharRom(addr)
}

// ReadDirect reads a byte directly from the specified address without applying any memory management or configuration rules.
func (w *PLASocket) ReadDirect(addr uint16) uint8 {
	return w.pla.ReadDirect(addr)
}

// ReadColor reads a color value from the specified memory address in the PLA socket and returns it as a uint8.
func (w *PLASocket) ReadColor(addr uint16) uint8 {
	return w.pla.ReadColor(addr)
}

// SetWriteTrigger sets a write trigger at the specified memory address and associates it with the provided callback function.
func (w *PLASocket) SetWriteTrigger(addr uint16, fn func(uint16, uint8)) int {
	return w.pla.SetWriteTrigger(addr, fn)
}

// RemoveRamTrigger removes the RAM trigger with the specified id at the given memory address.
func (w *PLASocket) RemoveRamTrigger(addr uint16, id int) {
	w.pla.RemoveRamTrigger(addr, id)
}
