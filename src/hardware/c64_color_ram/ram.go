package ram_c64

import (
	"github.com/markel1974/c64emu/src/common/filler"
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/references"
)

type ColorRam struct {
	*component.BaseComponent
	color       []byte //SRAM 2114 (1K x 4 bit)
	colorFiller *filler.Filler
}

// NewColorRam creates and initializes a new Ram instance with a parent component, factory, label, and instance number.
func NewColorRam(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *ColorRam {
	rl := &ColorRam{
		BaseComponent: component.NewBaseComponent(),
		color:         make([]byte, 0x0400),
		colorFiller:   filler.New(255, 128, 0, 0, 0, 0, 0, filler.InitRandomChanceHalf),
	}
	rl.BaseComponent.Register(factory, parent, Identifier(), rl, references.IdIC64ColorRam(rl, label, instance))
	return rl
}

// Setup initializes the RAM and color memory with specific patterns using the associated Filler objects.
func (r *ColorRam) Setup() error {
	r.colorFiller.InitWithPattern(r.color, uint(len(r.color)))
	return nil
}

// Bind connects the Ram component to the provided IC64RamSocket instance and initializes any required references or states.
func (r *ColorRam) Bind(_ references.IColorRamC64Socket) error {
	return nil
}

// Connect establishes necessary connections for the Ram component to function properly. Returns an error if it fails.
func (r *ColorRam) Connect() error {
	return nil
}

// Reset clears and reinitializes the RAM to its default state as defined during setup.
func (r *ColorRam) Reset() {
}

// EmulationRequired checks whether emulation is necessary for the current RAM implementation. Always returns false.
func (r *ColorRam) EmulationRequired() bool {
	return false
}

// Emulate performs the main emulation logic for the Ram component during a simulation cycle.
func (r *ColorRam) Emulate() {
}

// Internal determines whether the RAM component operates internally within the system. Always returns false.
func (r *ColorRam) Internal() bool {
	return false
}

// ReadColor retrieves a byte from the color memory at the specified address.
func (r *ColorRam) Read(addr uint16) uint8 {
	return r.color[addr&0x03ff]
}

// WriteColor writes an 8-bit value to the specified address in the color memory segment of the RAM.
func (r *ColorRam) Write(addr uint16, data uint8) {
	r.color[addr&0x03ff] = data & 0x0f
}

// Size returns the length of the `ram` slice, representing the total size of the RAM in bytes.
func (r *ColorRam) Size() int {
	return len(r.color)
}
