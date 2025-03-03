package board

import (
	"github.com/markel1974/c64emu/src/components/quartz"
	"github.com/markel1974/c64emu/src/signals"
)

type Expansion struct {
	board *Board
}

func NewExpansion(board *Board) *Expansion {
	return &Expansion{board: board}
}

func (s *Expansion) Reset() {
}

func (s *Expansion) Read(addr uint16) uint8 {
	return s.board.pla.Read(addr)
}

func (s *Expansion) Write(addr uint16, data uint8) {
	s.board.pla.Write(addr, data)
}

func (s *Expansion) GameExRomConfigChanged() {
	s.board.pla.RebuildMemoryConfig()
}

func (s *Expansion) NMITrigger() {
	s.board.pic.TriggerNMI()
}

func (s *Expansion) SetDMALow(v bool) {
	//s.board.dmaLowSlot(v)
}

func (s *Expansion) ResetTrigger() {
	s.board.pic.TriggerReset()
}

func (s *Expansion) IRQTrigger() {
	s.board.pic.TriggerIRQ(intrIrqExpansionBit)
}

func (s *Expansion) IRQClear() {
	s.board.pic.ClearIRQ(intrIrqExpansionBit)
}

func (s *Expansion) IRQTriggerBind(fn func(uint32)) {
	if s.board.expansionIrqTrigger == nil {
		s.board.expansionIrqTrigger = signals.NewSignalUint32()
	}
	s.board.expansionIrqTrigger.Bind(fn)
}

func (s *Expansion) IRQClearBind(fn func(uint32)) {
	if s.board.expansionIrqClear == nil {
		s.board.expansionIrqClear = signals.NewSignalUint32()
	}
	s.board.expansionIrqClear.Bind(fn)
}

func (s *Expansion) BusAvailable() bool {
	//return !s.board.vic.GetBALow()
	return true
}

func (s *Expansion) AECAvailable() bool {
	//return !s.board.vic.GetAECLow()
	return true
}

func (s *Expansion) Cycle() uint64 {
	//return s.board.quartz.Cycle()
	return 0
}

func (s *Expansion) CycleAlarm(id string, callback quartz.AlarmCallback) *quartz.Alarm {
	return nil
	//return s.board.quartz.NewAlarm(id, callback)
}

func (s *Expansion) RamSetWriteTrigger(addr uint16, fn func(uint16, uint8)) int {
	return s.board.pla.SetWriteTrigger(addr, fn)
}

func (s *Expansion) RamRemoveWriteTrigger(addr uint16, id int) {
	s.board.pla.RemoveRamTrigger(addr, id)
}

func (s *Expansion) RmwFlags() uint8 {
	//TODO IMPLEMENT cpu rmw flags
	return 0
}
