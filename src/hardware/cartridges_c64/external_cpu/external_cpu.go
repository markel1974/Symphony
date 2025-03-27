package external_cpu

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/hardware/6510"
	"github.com/markel1974/c64emu/src/hardware/pic_6510"
	"github.com/markel1974/c64emu/src/hardware/quartz"
	"github.com/markel1974/c64emu/src/references"
)

// Id is a constant representing the identifier "SCPU" for a specific hardware type or component registration.
const Id = "SCPU"

// CartridgeExternalCPU represents an external CPU module with its associated components and connections for system integration.
// It includes an ID, expansion board, CPU socket, programmable interrupt controller, CPU, and quartz clock source.
type CartridgeExternalCPU struct {
	*component.BaseComponent
	loaderId  string
	board     references.IExpansionC64
	cpuSocket *CPUSocket
	pic       *pic_6510.Pic
	cpu       *mos6510.CPU
	quartz    *quartz.Quartz
}

// NewExternalCPU returns a new instance of the CartridgeExternalCPU struct implementing the ICartridgeC64 interface.
func NewExternalCPU(parent references.IComponent, factory references.IComponentFactory, instance int) *CartridgeExternalCPU {
	r := &CartridgeExternalCPU{
		BaseComponent: component.NewBaseComponent(),
		board:         nil,
		cpuSocket:     nil,
		pic:           nil,
		cpu:           nil,
		quartz:        nil,
	}
	r.BaseComponent.Register(factory, parent, Identifier(), r, references.IdICartridgeC64(r, instance))
	return r
}

func New(parent references.IComponent, factory references.IComponentFactory, instance int) references.ICartridgeC64 {
	return NewExternalCPU(parent, factory, instance)
}

// EmulationRequired determines if CPU emulation is needed for the associated external CPU instance. Always returns true.
func (s *CartridgeExternalCPU) EmulationRequired() bool {
	return true
}

// Setup initializes the CartridgeExternalCPU with the provided expansion board and CRT loader, configuring its internal components.
func (s *CartridgeExternalCPU) Setup(board references.IExpansionC64, ldr references.ICartridgeLoaderC64, cfg *config.Config) error {
	s.board = board
	s.loaderId = ldr.GetId()
	s.board.SetDMALow(true)
	s.quartz = quartz.NewQuartz(s, s.GetFactory(), 0)
	s.pic = pic_6510.NewPIC(s, s.GetFactory(), 0)
	if err := s.pic.Setup(s, cfg, s.quartz); err != nil {
		return err
	}
	if err := s.pic.Connect(); err != nil {
		return err
	}
	s.cpuSocket = NewCPUSocket()
	s.cpu = mos6510.NewCPU(s, s.GetFactory(), 0)
	if err := s.cpuSocket.Setup(s); err != nil {
		return err
	}
	if err := s.cpu.Setup(s.cpuSocket, cfg); err != nil {
		return err
	}
	s.board.IRQTriggerBind(s.pic.TriggerIRQ)
	s.board.IRQClearBind(s.pic.ClearIRQ)
	return nil
}

func (s *CartridgeExternalCPU) GetLoaderId() string {
	return s.loaderId
}

// Reset reinitializes the CPU to its default state by invoking its internal Reset method.
func (s *CartridgeExternalCPU) Reset() {
	s.cpu.Reset()
}

func (s *CartridgeExternalCPU) Connect() error {
	return nil
}

func (s *CartridgeExternalCPU) Internal() bool {
	return false
}

// Emulate simulates the operation of the external CPU at a fixed frequency by processing cycles and managing bus signals.
func (s *CartridgeExternalCPU) Emulate() {
	// TEST SUPER CPU.....
	const mhz = 20
	// TODO TRIGGER BALOW-AECLOW HAS SIGNAL
	ba := s.board.BusAvailable()
	aec := s.board.AECAvailable()
	s.cpu.SetRDYLow(!ba)
	s.cpu.SetAECLow(!aec)
	for x := 0; x < mhz; x++ {
		s.cpu.Emulate()
		s.quartz.Emulate()
	}
}

// GetExRom returns the external ROM configuration signal as a uint8 value. It often indicates the cartridge's EXROM state.
func (s *CartridgeExternalCPU) GetExRom() uint8 {
	return 1
}

// GetGame retrieves the current game signal state for the CartridgeExternalCPU, returning a uint8 representation of its value.
func (s *CartridgeExternalCPU) GetGame() uint8 {
	return 1
}

// IORead reads data from the specified memory address and returns the value along with a success status flag.
func (s *CartridgeExternalCPU) IORead(addr uint16) (uint8, bool) {
	return 0, false
}

// IOWrite writes a byte of data to the specified address in the I/O space.
// Returns true if the write operation is handled; otherwise, false.
func (s *CartridgeExternalCPU) IOWrite(addr uint16, data uint8) bool {
	return false
}

// Write stores an 8-bit data value to a specified address within a given ROM interval and returns success status as a boolean.
func (s *CartridgeExternalCPU) Write(i references.RomInterval, addr uint16, data uint8) bool {
	return false
}

// Read fetches a byte and status from the specified address within the ROM interval.
func (s *CartridgeExternalCPU) Read(i references.RomInterval, addr uint16) (uint8, bool) {
	return 0, false
}

// Detach releases the resources and internal bindings associated with the CartridgeExternalCPU, ensuring proper cleanup.
func (s *CartridgeExternalCPU) Detach() error {
	//TODO IMPLEMENT
	return nil
}
