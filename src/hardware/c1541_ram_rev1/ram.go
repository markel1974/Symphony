package c1541_ram_rev1

import (
	"github.com/markel1974/c64emu/src/common/filler"
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/references"
)

const c1541RamSize = 0x0800

// Ram represents a memory component containing raw byte storage and write triggers, extending BaseComponent.
type Ram struct {
	*component.BaseComponent
	ram    []byte
	filler *filler.Filler
}

// NewRam initializes and returns a new Ram instance, registering it to the given parent component and factory setup.
func NewRam(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *Ram {
	rl := &Ram{
		BaseComponent: component.NewBaseComponent(),
		ram:           make([]byte, c1541RamSize),
		filler:        filler.New(255, 128, 0, 0, 0, 0, 0, 0),
	}
	rl.BaseComponent.Register(factory, parent, Identifier(), rl, references.IdIC1541Ram(rl, label, instance))
	return rl
}

// Setup initializes the necessary configurations for the Ram instance. It returns an error if the setup fails.
func (r *Ram) Setup() error {
	r.filler.InitWithPattern(r.ram, uint(len(r.ram)))
	return nil
}

// Bind integrates the current RAM instance with the given IC64RamSocket interface and initializes memory with a pattern.
func (r *Ram) Bind(_ references.IC1541RamSocket) error {
	return nil
}

// Connect establishes necessary connections or bindings for the Ram component. Returns an error if the operation fails.
func (r *Ram) Connect() error {
	return nil
}

// Reset clears or reinitializes the internal state of the RAM component to its default or initial configuration.
func (r *Ram) Reset() {
}

// EmulationRequired checks whether emulation is necessary for the RAM component and always returns false.
func (r *Ram) EmulationRequired() bool {
	return false
}

// Emulate runs the emulation logic for the RAM component during each cycle of the emulator.
func (r *Ram) Emulate() {
}

// Internal determines whether the RAM component is configured as an internal component.
func (r *Ram) Internal() bool {
	return false
}

// Read retrieves a byte from the specified address in the RAM.
func (r *Ram) Read(addr uint16) uint8 {
	return r.ram[addr]
}

// Write writes the provided 8-bit data to the specified 16-bit memory address in RAM and triggers any associated actions.
func (r *Ram) Write(addr uint16, data uint8) {
	r.ram[addr] = data
}

func (r *Ram) Size() int {
	return len(r.ram)
}
