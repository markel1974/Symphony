package external_cpu

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/hardware/mos6510_pic_rev1"
	"github.com/markel1974/c64emu/src/hardware/mos6510_rev1"
	"github.com/markel1974/c64emu/src/hardware/quartz_rev1"
	"github.com/markel1974/c64emu/src/references"
)

// Id is a constant representing the identifier "SCPU" for a specific hardware type or component registration.
const Id = "SCPU"

// CartridgeExternalCPU represents an external CPU module with its associated components and connections for system integration.
// It includes an ID, expansion board, CPU socket, programmable interrupt controller, CPU, and quartz clock source.
type CartridgeExternalCPU struct {
	*component.BaseComponent
	loaderId string
	board    references.IC64Expansion
	pic      references.IMos6510Pic
	cpu      references.IMos6510
	quartz   references.IQuartz
	cfg      *config.Config
	spec     *references.C64CartridgeSpec
}

// NewExternalCPU returns a new instance of the CartridgeExternalCPU struct implementing the IC64Cartridge interface.
func NewExternalCPU(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *CartridgeExternalCPU {
	r := &CartridgeExternalCPU{
		BaseComponent: component.NewBaseComponent(),
		board:         nil,
		pic:           nil,
		cpu:           nil,
		quartz:        nil,
		cfg:           nil,
		spec:          references.C64CartridgeSpecOff,
	}
	r.BaseComponent.Register(factory, parent, Identifier(), r, references.IdIC64Cartridge(r, label, instance))
	return r
}

func New(parent references.IComponent, factory references.IComponentFactory, label string, instance int) references.IC64Cartridge {
	return NewExternalCPU(parent, factory, label, instance)
}

// EmulationRequired determines if CPU emulation is needed for the associated external CPU instance. Always returns true.
func (s *CartridgeExternalCPU) EmulationRequired() bool {
	return true
}

func (s *CartridgeExternalCPU) Setup() error {
	s.cfg = s.GetFactory().GetConfig()
	return nil
}

func (s *CartridgeExternalCPU) Bind(board references.IC64Expansion, ldr references.IC64CartridgeLoader) error {
	//TODO REWRITE!!!
	const instance = 1000
	s.board = board
	s.loaderId = ldr.Id()
	s.board.SetDMALow(true)

	cpu, err := s.GetFactory().Create(s, Identifier(), mos6510_rev1.Identifier(), instance)
	if err != nil {
		return err
	}
	if s.cpu, err = references.ComponentToIMos6510(cpu); err != nil {
		return err
	}

	p, err := s.GetFactory().Create(s, Identifier(), mos6510_pic_rev1.Identifier(), instance)
	if err != nil {
		return err
	}
	if s.pic, err = references.ComponentToIMos6510Pic(p); err != nil {
		return err
	}
	q, err := s.GetFactory().Create(s, Identifier(), quartz_rev1.Identifier(), instance)
	if err != nil {
		return err
	}
	if s.quartz, err = references.ComponentToIQuartz(q); err != nil {
		return err
	}

	cc := make(map[string]references.IComponent)
	cc[cpu.HardwareId()] = cpu
	cc[p.HardwareId()] = p
	cc[q.HardwareId()] = q
	for _, v := range cc {
		if err = v.Setup(); err != nil {
			return err
		}
	}
	if err = s.cpu.Bind(s, s.pic, s.board); err != nil {
		return err
	}
	for _, v := range cc {
		if err = v.Connect(); err != nil {
			return err
		}
	}
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

// Config returns the Game line status, ExROM line state, and a boolean indicating successful configuration retrieval.
func (s *CartridgeExternalCPU) Config() (uint8, uint8, bool) {
	return s.spec.Game, s.spec.ExRom, true
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

// IRQ handles the Interrupt Request (IRQ) signal for the CartridgeGeneric, enabling appropriate cartridge-specific behavior.
func (s *CartridgeExternalCPU) IRQ(d uint32) {
	s.pic.TriggerIRQ(d)
}

// IRQCLear clears the state of any active Interrupt Requests (IRQ) for the CartridgeGeneric.
func (s *CartridgeExternalCPU) IRQCLear(d uint32) {
	s.pic.ClearIRQ(d)
}

// HardwareButton handles the system response to a physical button press event, updating cartridge state as necessary.
func (s *CartridgeExternalCPU) HardwareButton(pressed bool, value uint8) {
}

// Read fetches a byte and status from the specified address within the ROM interval.
func (s *CartridgeExternalCPU) Read(addr uint16) uint8 {
	return 0
}

// Detach releases the resources and internal bindings associated with the CartridgeExternalCPU, ensuring proper cleanup.
func (s *CartridgeExternalCPU) Detach() error {
	//TODO IMPLEMENT
	return nil
}
