package board

import mos6569 "github.com/markel1974/c64emu/src/components/vic"

type VicSocket struct {
	board  *Board
	intrId uint32
}

func NewVicSocket() *VicSocket {
	return &VicSocket{}
}

func (v *VicSocket) Setup(board *Board, intrId uint32) {
	v.board = board
	v.intrId = intrId
}

func (v *VicSocket) Reset() {
	v.board.vic.Reset()
}

func (v *VicSocket) Cycle() uint64 {
	return v.board.quartz.Cycle()
}

func (v *VicSocket) GetDisplayBuffer() mos6569.IDisplayBuffer {
	return v.board.db
}

func (v *VicSocket) GetBanks() mos6569.IBanks {
	return v.board.banks
}

func (v *VicSocket) IRQTrigger() {
	v.board.irqTriggerSlot(v.intrId)
}

func (v *VicSocket) IRQClear() {
	v.board.irqClearSlot(v.intrId)
}

func (v *VicSocket) BALow(d bool) {
	v.board.rdyLowSlot(d)
}

func (v *VicSocket) AECLow(d bool) {
	v.board.aecLowSlot(d)
}

func (v *VicSocket) LastCycle() {
	v.board.vicLastCycleSLot()
}

func (v *VicSocket) VBlank() {
	v.board.vicVBlankSlot()
}
