package supercpu

import (
	"github.com/markel1974/c64emu/src/board/cartridges/icartridge"
	"github.com/markel1974/c64emu/src/board/cartridges/loader"
	cpu2 "github.com/markel1974/c64emu/src/components/6510"
	"github.com/markel1974/c64emu/src/components/quartz"
)

const Id = "SCPU"

type SuperCPU struct {
	id     string
	board  icartridge.IExpansion
	pic    *cpu2.Pic
	cpu    *cpu2.MOS6510
	quartz *quartz.Quartz
}

func New() icartridge.ICartridge {
	r := &SuperCPU{
		id:     "",
		board:  nil,
		pic:    nil,
		cpu:    nil,
		quartz: nil,
	}
	return r
}

func (s *SuperCPU) GetId() string {
	return s.id
}

func (s *SuperCPU) EmulationRequired() bool {
	return true
}

func (s *SuperCPU) Setup(board icartridge.IExpansion, ldr *loader.CRTLoader) error {
	s.board = board
	s.id = ldr.GetId()
	s.board.SetDMALow(true)

	s.quartz = quartz.NewQuartz()
	s.pic = cpu2.NewPic()

	s.pic.Setup(s.quartz)

	s.cpu = cpu2.NewMOS6510("superCpu")
	s.cpu.Setup(s.pic, board, nil)

	s.board.IRQTriggerBind(s.pic.TriggerIRQ)
	s.board.IRQClearBind(s.pic.ClearIRQ)

	return nil
}

func (s *SuperCPU) Reset() {
	s.cpu.Reset()
}

func (s *SuperCPU) Emulate() {
	// TEST SUPER CPU.....
	const mhz = 20
	// TODO TRIGGER BALOW-AECLOW HAS SIGNAL
	aecLow := s.board.AECAvailable()
	s.cpu.SetRDYLow(aecLow)
	for x := 0; x < mhz; x++ {
		s.cpu.Emulate()
		s.quartz.AddCycle()
	}
}

func (s *SuperCPU) GetExRom() uint8 {
	return 1
}

func (s *SuperCPU) GetGame() uint8 {
	return 1
}

func (s *SuperCPU) IORead(addr uint16) (uint8, bool) {
	return 0, false
}

func (s *SuperCPU) IOWrite(addr uint16, data uint8) bool {
	return false
}

func (s *SuperCPU) Write(i icartridge.RomInterval, addr uint16, data uint8) bool {
	return false
}

func (s *SuperCPU) Read(i icartridge.RomInterval, addr uint16) (uint8, bool) {
	return 0, false
}

func (s *SuperCPU) Detach() error {
	//TODO IMPLEMENT
	return nil
}
