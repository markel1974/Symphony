package externalcpu

import (
	"github.com/markel1974/c64emu/src/c64/cartridges/icartridge"
	"github.com/markel1974/c64emu/src/c64/cartridges/loader"
	"github.com/markel1974/c64emu/src/components/6510"
	"github.com/markel1974/c64emu/src/components/quartz"
)

const Id = "SCPU"

type ExternalCPU struct {
	id        string
	board     icartridge.IExpansion
	cpuSocket *CPUSocket
	pic       *mos6510.Pic
	cpu       *mos6510.CPU
	quartz    *quartz.Quartz
}

func New() icartridge.ICartridge {
	r := &ExternalCPU{
		id:        "",
		board:     nil,
		cpuSocket: nil,
		pic:       nil,
		cpu:       nil,
		quartz:    nil,
	}
	return r
}

func (s *ExternalCPU) GetId() string {
	return s.id
}

func (s *ExternalCPU) EmulationRequired() bool {
	return true
}

func (s *ExternalCPU) Setup(board icartridge.IExpansion, ldr *loader.CRTLoader) error {
	s.board = board
	s.id = ldr.GetId()
	s.board.SetDMALow(true)

	s.quartz = quartz.NewQuartz()
	s.pic = mos6510.NewPic()

	s.pic.Setup(s.quartz)

	s.cpuSocket = NewCPUSocket()
	s.cpu = mos6510.NewCPU("superCpu")

	s.cpuSocket.Setup(s)
	s.cpu.Setup(s.cpuSocket)

	s.board.IRQTriggerBind(s.pic.TriggerIRQ)
	s.board.IRQClearBind(s.pic.ClearIRQ)

	return nil
}

func (s *ExternalCPU) Reset() {
	s.cpu.Reset()
}

func (s *ExternalCPU) Emulate() {
	// TEST SUPER CPU.....
	const mhz = 20
	// TODO TRIGGER BALOW-AECLOW HAS SIGNAL
	ba := s.board.BusAvailable()
	aec := s.board.AECAvailable()
	s.cpu.SetRDYLow(!ba)
	s.cpu.SetAECLow(!aec)
	for x := 0; x < mhz; x++ {
		s.cpu.Emulate()
		s.quartz.AddCycle()
	}
}

func (s *ExternalCPU) GetExRom() uint8 {
	return 1
}

func (s *ExternalCPU) GetGame() uint8 {
	return 1
}

func (s *ExternalCPU) IORead(addr uint16) (uint8, bool) {
	return 0, false
}

func (s *ExternalCPU) IOWrite(addr uint16, data uint8) bool {
	return false
}

func (s *ExternalCPU) Write(i icartridge.RomInterval, addr uint16, data uint8) bool {
	return false
}

func (s *ExternalCPU) Read(i icartridge.RomInterval, addr uint16) (uint8, bool) {
	return 0, false
}

func (s *ExternalCPU) Detach() error {
	//TODO IMPLEMENT
	return nil
}
