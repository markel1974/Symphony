package board

import (
	"github.com/markel1974/c64emu/src/components/board"
	mos6569 "github.com/markel1974/c64emu/src/components/vic"
)

// VicSocket represents a virtual interface connector socket with a reference to a board and an interrupt identifier.
type VicSocket struct {
	board  *Board
	intrId uint32
}

// NewVicSocket creates and returns a new instance of the VicSocket struct.
func NewVicSocket() *VicSocket {
	return &VicSocket{}
}

// Setup initializes the VicSocket with the given board and interrupt ID.
func (v *VicSocket) Setup(board *Board, intrId uint32) {
	v.board = board
	v.intrId = intrId
}

// Reset resets the VIC component of the associated board by invoking its Reset method.
func (v *VicSocket) Reset() {
	v.board.vic.Reset()
}

// Cycle retrieves the current cycle count from the associated Quartz scheduler.
func (v *VicSocket) Cycle() uint64 {
	return v.board.quartz.Cycle()
}

// GetDisplayBuffer returns the IDisplayBuffer instance associated with the VicSocket's board.
func (v *VicSocket) GetDisplayBuffer() board.IDisplayBuffer {
	return v.board.db
}

// GetBanks returns an implementation of the mos6569.IBanks interface, which provides access to memory handling operations.
func (v *VicSocket) GetBanks() mos6569.IBanks {
	return v.board.pla
}

// IRQTrigger signals an interrupt request by invoking the IRQ trigger mechanism on the associated board slot.
func (v *VicSocket) IRQTrigger() {
	v.board.irqTriggerSlot(v.intrId)
}

// IRQClear clears the interrupt request for the associated slot identified by the intrId of the VicSocket.
func (v *VicSocket) IRQClear() {
	v.board.irqClearSlot(v.intrId)
}

// BALow sets the BA (Bus Available) line low or high based on the provided boolean value.
func (v *VicSocket) BALow(d bool) {
	v.board.rdyLowSlot(d)
}

// AECLow controls the state of the Address Enable Control (AEC) line. It sets the AEC signal to low if d is true.
func (v *VicSocket) AECLow(d bool) {
	v.board.aecLowSlot(d)
}

// LastCycle triggers the last cycle slot operation on the VIC through the connected board.
func (v *VicSocket) LastCycle() {
	v.board.vicLastCycleSLot()
}

// VBlank handles the vertical blanking phase by triggering the corresponding slot on the associated board instance.
func (v *VicSocket) VBlank() {
	v.board.vicVBlankSlot()
}
