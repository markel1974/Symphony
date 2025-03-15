package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// VicSocket represents a virtual interface connector socket with a reference to a board and an interrupt identifier.
type VicSocket struct {
	vic    references.IVic
	board  *Board
	intrId uint32
}

// NewVicSocket creates and returns a new instance of the VicSocket struct.
func NewVicSocket() *VicSocket {
	return &VicSocket{
		vic:    nil,
		board:  nil,
		intrId: intrIrqVicBit,
	}
}

// Setup initializes the VicSocket with the given board and interrupt ID.
func (v *VicSocket) Setup(board *Board, vic references.IVic) error {
	v.board = board
	v.vic = vic
	v.vic.Setup(v, v.board.cfg)
	return nil
}

// Reset resets the VIC component of the associated board by invoking its Reset method.
func (v *VicSocket) Reset() {
	v.vic.Reset()
}

func (v *VicSocket) Emulate() {
	v.vic.Emulate()
}

func (v *VicSocket) GetText() []byte {
	return v.vic.GetText()
}

func (v *VicSocket) GetBALow() bool {
	return v.vic.GetBALow()
}

func (v *VicSocket) GetAECLow() bool {
	return v.vic.GetAECLow()
}

func (v *VicSocket) LightPenTrigger() {
	v.vic.LightPenTrigger()
}

func (v *VicSocket) ChangedVA(newVA uint8) {
	v.vic.ChangedVA(newVA)
}

// Cycle retrieves the current cycle count from the associated Quartz scheduler.
func (v *VicSocket) Cycle() uint64 {
	return v.board.quartz.Cycle()
}

// GetDisplayBuffer returns the IDisplayBuffer instance associated with the VicSocket's board.
func (v *VicSocket) GetDisplayBuffer() references.IDisplayBuffer {
	return v.board.db
}

// GetBanks returns an implementation of the mos6569.IVicBanks interface, which provides access to memory handling operations.
func (v *VicSocket) GetBanks() references.IVicBanks {
	return v.board.plaSocket
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
