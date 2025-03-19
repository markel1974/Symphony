package external_cpu

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/hardware/6510"
	"github.com/markel1974/c64emu/src/hardware/pic_6510"
	"github.com/markel1974/c64emu/src/hardware/quartz"
	"github.com/markel1974/c64emu/src/references"
)

// Id is a constant representing the identifier "SCPU" for a specific hardware type or component registration.
const Id = "SCPU"

// ExternalCPU represents an external CPU module with its associated components and connections for system integration.
// It includes an ID, expansion board, CPU socket, programmable interrupt controller, CPU, and quartz clock source.
type ExternalCPU struct {
	*component.BaseComponent
	factory   references.IComponentFactory
	loaderId  string
	board     references.IExpansionC64
	cpuSocket *CPUSocket
	pic       *pic_6510.Pic
	cpu       *mos6510.CPU
	quartz    *quartz.Quartz
}

// New returns a new instance of the ExternalCPU struct implementing the ICartridgeC64 interface.
func New(parent references.IComponent, factory references.IComponentFactory, label int) references.ICartridgeC64 {
	r := &ExternalCPU{
		BaseComponent: component.NewBaseComponent("externalCpu", label, references.IdICartridgeC64),
		factory:       factory,
		board:         nil,
		cpuSocket:     nil,
		pic:           nil,
		cpu:           nil,
		quartz:        nil,
	}
	component.Register(parent, r)
	return r
}

// EmulationRequired determines if CPU emulation is needed for the associated external CPU instance. Always returns true.
func (s *ExternalCPU) EmulationRequired() bool {
	return true
}

// Setup initializes the ExternalCPU with the provided expansion board and CRT loader, configuring its internal components.
func (s *ExternalCPU) Setup(board references.IExpansionC64, ldr references.ICartridgeLoaderC64) error {
	s.board = board
	s.loaderId = ldr.GetId()
	s.board.SetDMALow(true)

	s.quartz = quartz.NewQuartz(s, s.factory, 0)
	s.pic = pic_6510.NewPIC(s, s.factory, 0)

	s.pic.Setup(s.quartz)

	s.cpuSocket = NewCPUSocket()
	s.cpu = mos6510.NewCPU(s, s.factory, 0)

	s.cpuSocket.Setup(s)
	s.cpu.Setup(s.cpuSocket)

	s.board.IRQTriggerBind(s.pic.TriggerIRQ)
	s.board.IRQClearBind(s.pic.ClearIRQ)

	return nil
}

func (s *ExternalCPU) GetLoaderId() string {
	return s.loaderId
}

// Reset reinitializes the CPU to its default state by invoking its internal Reset method.
func (s *ExternalCPU) Reset() {
	s.cpu.Reset()
}

// Emulate simulates the operation of the external CPU at a fixed frequency by processing cycles and managing bus signals.
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

// GetExRom returns the external ROM configuration signal as a uint8 value. It often indicates the cartridge's EXROM state.
func (s *ExternalCPU) GetExRom() uint8 {
	return 1
}

// GetGame retrieves the current game signal state for the ExternalCPU, returning a uint8 representation of its value.
func (s *ExternalCPU) GetGame() uint8 {
	return 1
}

// IORead reads data from the specified memory address and returns the value along with a success status flag.
func (s *ExternalCPU) IORead(addr uint16) (uint8, bool) {
	return 0, false
}

// IOWrite writes a byte of data to the specified address in the I/O space.
// Returns true if the write operation is handled; otherwise, false.
func (s *ExternalCPU) IOWrite(addr uint16, data uint8) bool {
	return false
}

// Write stores an 8-bit data value to a specified address within a given ROM interval and returns success status as a boolean.
func (s *ExternalCPU) Write(i references.RomInterval, addr uint16, data uint8) bool {
	return false
}

// Read fetches a byte and status from the specified address within the ROM interval.
func (s *ExternalCPU) Read(i references.RomInterval, addr uint16) (uint8, bool) {
	return 0, false
}

// Detach releases the resources and internal bindings associated with the ExternalCPU, ensuring proper cleanup.
func (s *ExternalCPU) Detach() error {
	//TODO IMPLEMENT
	return nil
}
