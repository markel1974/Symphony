package board

import (
	"github.com/markel1974/c64emu/src/common/signals"
	"github.com/markel1974/c64emu/src/components/board"
	"github.com/markel1974/c64emu/src/components/quartz"
)

// Expansion represents a logical structure that associates a parent ID, its own ID, and a reference to a Board.
type Expansion struct {
	*board.BaseComponent
	board *Board
}

// NewExpansion creates a new instance of the Expansion structure with the given parentId and suffix.
// The `id` property is initialized by concatenating "expansion" with the specified suffix.
// Returns a pointer to the newly created Expansion object.
func NewExpansion(parentNode *board.Node, suffix string) *Expansion {
	e := &Expansion{
		BaseComponent: board.NewBaseComponent("expansion", suffix, nil),
	}
	board.AssignNode(parentNode, e)
	return e
}

// Setup initializes the Expansion by associating it with the provided board instance.
func (s *Expansion) Setup(board *Board) {
	s.board = board
}

// Reset reinitializes the state of the Expansion to its default configuration.
func (s *Expansion) Reset() {
}

// Read reads a byte of data from the specified memory address using the associated board's PLA component.
func (s *Expansion) Read(addr uint16) uint8 {
	return s.board.plaSocket.Read(addr)
}

// Write writes a value to the specified address through the board's PLA module.
func (s *Expansion) Write(addr uint16, data uint8) {
	s.board.plaSocket.Write(addr, data)
}

// GameExRomConfigChanged triggers a rebuild of the memory configuration within the PLA board.
func (s *Expansion) GameExRomConfigChanged() {
	s.board.plaSocket.RebuildMemoryConfig()
}

// NMITrigger triggers a Non-Maskable Interrupt (NMI) by invoking the associated method on the programmable interrupt controller.
func (s *Expansion) NMITrigger() {
	s.board.pic.TriggerNMI()
}

// SetDMALow sets the DMA low signal state on the associated board.
func (s *Expansion) SetDMALow(v bool) {
	s.board.dmaLowSlot(v)
}

// ResetTrigger activates the reset sequence by invoking the TriggerReset function of the associated Pic instance.
func (s *Expansion) ResetTrigger() {
	s.board.pic.TriggerReset()
}

// IRQTrigger triggers the IRQ interrupt for the expansion module using the programmable interrupt controller (PIC).
func (s *Expansion) IRQTrigger() {
	s.board.pic.TriggerIRQ(intrIrqExpansionBit)
}

// IRQClear clears the IRQ signal for the expansion bit via the associated programmable interrupt controller (PIC).
func (s *Expansion) IRQClear() {
	s.board.pic.ClearIRQ(intrIrqExpansionBit)
}

// IRQTriggerBind binds the provided function to the IRQ trigger event, initializing the signal if not already done.
func (s *Expansion) IRQTriggerBind(fn func(uint32)) {
	if s.board.expansionIrqTrigger == nil {
		s.board.expansionIrqTrigger = signals.NewSignalUint32()
	}
	s.board.expansionIrqTrigger.Bind(fn)
}

// IRQClearBind binds a callback function to the clear interrupt request signal for the expansion.
func (s *Expansion) IRQClearBind(fn func(uint32)) {
	if s.board.expansionIrqClear == nil {
		s.board.expansionIrqClear = signals.NewSignalUint32()
	}
	s.board.expansionIrqClear.Bind(fn)
}

// BusAvailable checks if the bus is available by verifying if the BA (Bus Available) signal is not asserted low.
func (s *Expansion) BusAvailable() bool {
	return !s.board.vicSocket.GetBALow()
}

// AECAvailable determines if the Address Enable Control (AEC) signal is available by checking if the AEC line is not low.
func (s *Expansion) AECAvailable() bool {
	return !s.board.vicSocket.GetAECLow()
}

// Cycle returns the current cycle count as a uint64 from the underlying Quartz instance associated with the Expansion.
func (s *Expansion) Cycle() uint64 {
	return s.board.quartz.Cycle()
}

// CycleAlarm creates and registers a new alarm with a unique identifier and a callback function for execution upon triggering.
func (s *Expansion) CycleAlarm(id string, callback quartz.AlarmCallback) *quartz.Alarm {
	return s.board.quartz.NewAlarm(id, callback)
}

// RamSetWriteTrigger sets a write trigger for a specific RAM address and executes the provided callback on writes.
func (s *Expansion) RamSetWriteTrigger(addr uint16, fn func(uint16, uint8)) int {
	return s.board.plaSocket.SetWriteTrigger(addr, fn)
}

// RamRemoveWriteTrigger removes a write trigger from the specified RAM address using the given trigger id.
func (s *Expansion) RamRemoveWriteTrigger(addr uint16, id int) {
	s.board.plaSocket.RemoveRamTrigger(addr, id)
}

// RmwFlags computes and returns the read-modify-write flags for CPU operations.
func (s *Expansion) RmwFlags() uint8 {
	//TODO IMPLEMENT cpu rmw flags
	return 0
}
