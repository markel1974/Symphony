package c64_ram_rev1

import (
	"github.com/markel1974/symphony/src/common/filler"
	"github.com/markel1974/symphony/src/kernel/component"
	"github.com/markel1974/symphony/src/references"
)

// Ram represents a memory component with associated data buffers and fillers for color and data management.
type Ram struct {
	*component.BaseComponent
	reflect *RamReflect
	ram     []uint8 // symphony:export ram is a byte slice representing the memory storage of the RAM component.
	filler  *filler.Filler
}

// NewRam creates and initializes a new Ram instance with a parent component, factory, label, and instance number.
func NewRam(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *Ram {
	rl := &Ram{
		BaseComponent: component.NewBaseComponent(),
		ram:           make([]uint8, 0x10000),
		filler:        filler.New(255, 128, 0, 0, 0, 0, 0, 0),
	}
	rl.BaseComponent.Register(factory, parent, Identifier(), instance, rl, references.IdIC64Ram(rl, label, instance))
	rl.reflect = NewRamReflect(rl)
	return rl
}

// Setup initializes the RAM and color memory with specific patterns using the associated Filler objects.
func (r *Ram) Setup() error {
	r.filler.InitWithPattern(r.ram, uint(len(r.ram)))
	return nil
}

// Bind connects the Ram component to the provided IC64RamSocket instance and initializes any required references or states.
func (r *Ram) Bind(_ references.IC64RamSocket) error {
	return nil
}

// Connect establishes necessary connections for the Ram component to function properly. Returns an error if it fails.
func (r *Ram) Connect() error {
	return nil
}

// Reset clears and reinitializes the RAM to its default state as defined during setup.
func (r *Ram) Reset() {
	//r.filler.InitWithPattern(r.ram, uint(len(r.ram)))
}

// EmulationRequired checks whether emulation is necessary for the current RAM implementation. Always returns false.
func (r *Ram) EmulationRequired() bool {
	return false
}

// Emulate performs the main emulation logic for the Ram component during a simulation cycle.
func (r *Ram) Emulate() {
}

// Internal determines whether the RAM component operates internally within the system. Always returns false.
func (r *Ram) Internal() bool {
	return false
}

// Read retrieves a byte from the RAM at the specified memory address.
func (r *Ram) Read(addr uint16) uint8 {
	return r.ram[addr]
}

// Write stores the provided data byte at the specified memory address in the RAM.
func (r *Ram) Write(addr uint16, data uint8) {
	r.ram[addr] = data
}

// Size returns the length of the `ram` slice, representing the total size of the RAM in bytes.
func (r *Ram) Size() int {
	return len(r.ram)
}
