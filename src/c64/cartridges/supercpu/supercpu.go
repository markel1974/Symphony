package supercpu

import (
	"github.com/markel1974/c64emu/src/c64/cartridges/icartridge"
	"github.com/markel1974/c64emu/src/c64/cartridges/loader"
	mos6510 "github.com/markel1974/c64emu/src/components/6510"
	"github.com/markel1974/c64emu/src/components/quartz"
)

const Id = "SCPU"

const (
	intrRstBit = 1
	intrNmiBit = 2
	intrIrqBit = 3
)

type SuperCPU struct {
	id     string
	board  icartridge.IExpansion
	pic    *mos6510.Pic
	cpu    *mos6510.CPU
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
	s.pic = mos6510.NewPic(mos6510.MinIrqCycleDistance, intrRstBit, intrNmiBit, intrIrqBit)

	s.pic.Setup(s.quartz)

	s.cpu = mos6510.NewCPU("superCpu")
	s.cpu.Setup(s.pic, board)

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
