package board

import (
	"github.com/markel1974/c64emu/src/common/signals"
	"github.com/markel1974/c64emu/src/components/quartz"
)

// Expansion represents an extension of the game board.
type Expansion struct {
	board *Board
}

// NewExpansion creates and returns a new Expansion instance associated with the provided Board.
func NewExpansion(board *Board) *Expansion {
	return &Expansion{board: board}
}

// Reset resets the Expansion object to its initial default state.
func (s *Expansion) Reset() {
}

// Read retrieves an 8-bit value from the specified 16-bit address using the PLA component in the board.
func (s *Expansion) Read(addr uint16) uint8 {
	return s.board.pla.Read(addr)
}

// Write writes the specified 8-bit data to the given 16-bit address via the PLA in the Expansion's board.
func (s *Expansion) Write(addr uint16, data uint8) {
	s.board.pla.Write(addr, data)
}

// GameExRomConfigChanged updates the memory configuration by invoking the RebuildMemoryConfig method on the pla component.
func (s *Expansion) GameExRomConfigChanged() {
	s.board.pla.RebuildMemoryConfig()
}

// NMITrigger triggers a Non-Maskable Interrupt (NMI) through the programmable interrupt controller of the board.
func (s *Expansion) NMITrigger() {
	s.board.pic.TriggerNMI()
}

// SetDMALow sets the DMA low signal to the provided boolean value, affecting the DMA operation of the board.
func (s *Expansion) SetDMALow(v bool) {
	s.board.dmaLowSlot(v)
}

// ResetTrigger triggers a system reset by invoking the reset mechanism in the programmable interrupt controller (PIC).
func (s *Expansion) ResetTrigger() {
	s.board.pic.TriggerReset()
}

// IRQTrigger triggers an IRQ interrupt for the expansion component by setting the corresponding interrupt bit in the PIC.
func (s *Expansion) IRQTrigger() {
	s.board.pic.TriggerIRQ(intrIrqExpansionBit)
}

// IRQClear clears the expansion IRQ by resetting the associated interrupt bit in the PIC.
func (s *Expansion) IRQClear() {
	s.board.pic.ClearIRQ(intrIrqExpansionBit)
}

// IRQTriggerBind binds a given function to the IRQ trigger signal for the expansion using a uint32 parameter.
func (s *Expansion) IRQTriggerBind(fn func(uint32)) {
	if s.board.expansionIrqTrigger == nil {
		s.board.expansionIrqTrigger = signals.NewSignalUint32()
	}
	s.board.expansionIrqTrigger.Bind(fn)
}

// IRQClearBind binds a callback function to be executed when the expansion IRQ clear signal is emitted.
func (s *Expansion) IRQClearBind(fn func(uint32)) {
	if s.board.expansionIrqClear == nil {
		s.board.expansionIrqClear = signals.NewSignalUint32()
	}
	s.board.expansionIrqClear.Bind(fn)
}

// BusAvailable checks if the bus is available by determining whether the VIC chip's `BALow` signal is inactive.
func (s *Expansion) BusAvailable() bool {
	return !s.board.vic.GetBALow()
}

// AECAvailable checks whether AEC (Address Enable Control) is available by verifying if AECLow signal is inactive.
func (s *Expansion) AECAvailable() bool {
	return !s.board.vic.GetAECLow()
}

// Cycle returns the current cycle count from the associated Quartz instance.
func (s *Expansion) Cycle() uint64 {
	return s.board.quartz.Cycle()
}

// CycleAlarm creates a new alarm with the specified ID and callback, using the Quartz instance tied to the expansion board.
func (s *Expansion) CycleAlarm(id string, callback quartz.AlarmCallback) *quartz.Alarm {
	return s.board.quartz.NewAlarm(id, callback)
}

// RamSetWriteTrigger sets a write trigger on a specific RAM address using the provided callback function.
// Returns an identifier for the trigger.
func (s *Expansion) RamSetWriteTrigger(addr uint16, fn func(uint16, uint8)) int {
	return s.board.pla.SetWriteTrigger(addr, fn)
}

// RamRemoveWriteTrigger removes a write trigger for a specific RAM address and identifier.
func (s *Expansion) RamRemoveWriteTrigger(addr uint16, id int) {
	s.board.pla.RemoveRamTrigger(addr, id)
}

// RmwFlags retrieves the read-modify-write flags for the CPU operation.
func (s *Expansion) RmwFlags() uint8 {
	//TODO IMPLEMENT cpu rmw flags
	return 0
}
