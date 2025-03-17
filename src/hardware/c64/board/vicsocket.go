package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// VICSocket represents a virtual interface connector socket with a reference to a board and an interrupt identifier.
type VICSocket struct {
	references.IVIC
	board  *Board
	intrId uint32
}

// NewVICSocket creates and returns a new instance of the VICSocket struct.
func NewVICSocket() *VICSocket {
	return &VICSocket{
		IVIC:   nil,
		board:  nil,
		intrId: intrIrqVicBit,
	}
}

// Connect initializes the VICSocket with the given board and interrupt ID.
func (v *VICSocket) Connect(board *Board, vic references.IVIC) error {
	v.board = board
	v.IVIC = vic
	v.IVIC.Setup(v, v.board.cfg)
	return nil
}

// Cycle retrieves the current cycle count from the associated Quartz scheduler.
func (v *VICSocket) Cycle() uint64 {
	return v.board.quartzSocket.Cycle()
}

// GetDisplayBuffer returns the IDisplayBuffer instance associated with the VICSocket's board.
func (v *VICSocket) GetDisplayBuffer() references.IDisplayBuffer {
	return v.board.db
}

// GetBanks returns an implementation of the mos6569.IVICBanks interface, which provides access to memory handling operations.
func (v *VICSocket) GetBanks() references.IVICBanks {
	return v.board.plaSocket
}

// IRQTrigger signals an interrupt request by invoking the IRQ trigger mechanism on the associated board slot.
func (v *VICSocket) IRQTrigger() {
	v.board.irqTriggerSlot(v.intrId)
}

// IRQClear clears the interrupt request for the associated slot identified by the intrId of the VICSocket.
func (v *VICSocket) IRQClear() {
	v.board.irqClearSlot(v.intrId)
}

// BALow sets the BA (Bus Available) line low or high based on the provided boolean value.
func (v *VICSocket) BALow(d bool) {
	v.board.rdyLowSlot(d)
}

// AECLow controls the state of the Address Enable Control (AEC) line. It sets the AEC signal to low if d is true.
func (v *VICSocket) AECLow(d bool) {
	v.board.aecLowSlot(d)
}

// LastCycle triggers the last cycle slot operation on the VIC through the connected board.
func (v *VICSocket) LastCycle() {
	v.board.vicLastCycleSLot()
}

// VBlank handles the vertical blanking phase by triggering the corresponding slot on the associated board instance.
func (v *VICSocket) VBlank() {
	v.board.vicVBlankSlot()
}
