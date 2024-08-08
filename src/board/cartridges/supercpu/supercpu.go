package supercpu

import (
	"github.com/markel1974/c64emu/src/board/cartridges/icartridge"
	"github.com/markel1974/c64emu/src/board/cartridges/loader"
	"github.com/markel1974/c64emu/src/board/cpu"
)

const Id = "SCPU"

type SuperCPU struct {
	id    string
	board icartridge.IExpansion
	pic   *cpu.Pic
	cpu   *cpu.MOS6510
}

func New() icartridge.ICartridge {
	r := &SuperCPU{
		id:    "",
		board: nil,
		pic:   nil,
		cpu:   nil,
	}
	return r
}

func (s *SuperCPU) GetId() string {
	return s.id
}

func (s *SuperCPU) EmulationRequired() bool {
	return true
}

func (s *SuperCPU) Emulate() {
	// TODO CLEAR-TRIGGER RST/NMI/IRQ BALOW-AECLOW HAS SIGNAL
	irqLine := s.board.IRQLine()
	s.pic.SetIRQLine(irqLine)

	aecLow := s.board.AECAvailable()
	s.cpu.SetRDYLow(aecLow)

	s.cpu.Emulate()
}

func (s *SuperCPU) Reset() {
	s.cpu.Reset()
}

func (s *SuperCPU) Setup(board icartridge.IExpansion, ldr *loader.CRTLoader) error {
	s.board = board
	s.id = ldr.GetId()
	s.board.SetDMALow(true)

	s.pic = cpu.NewPic()
	s.pic.Setup(board.GetQuartz())

	s.cpu = cpu.NewMOS6510()
	s.cpu.Setup(s.pic, board, nil)

	//TODO MISSING IRQ/NMI/RESET
	//m.via1.SignalTriggerIRQBind(m.pic.TriggerIRQ)
	//m.via1.SignalClearIRQBind(m.pic.ClearIRQ)

	return nil
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
