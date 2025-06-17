package c64_color_ram_rev1

import (
	"github.com/markel1974/c64emu/src/common/filler"
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/references"
)

// ColorRam represents a color memory component that interfaces with a BaseComponent and includes an SRAM buffer.
type ColorRam struct {
	*component.BaseComponent
	color       []byte //SRAM 2114 (1K x 4 bit)
	colorFiller *filler.Filler
}

// NewColorRam creates and initializes a new ColorRam instance, registers it with the factory, and sets up essential components.
func NewColorRam(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *ColorRam {
	rl := &ColorRam{
		BaseComponent: component.NewBaseComponent(),
		color:         make([]byte, 0x0400),
		colorFiller:   filler.New(255, 128, 0, 0, 0, 0, 0, filler.InitRandomChanceHalf),
	}
	rl.BaseComponent.Register(factory, parent, Identifier(), rl, references.IdIC64ColorRam(rl, label, instance))
	return rl
}

// Setup initializes the `ColorRam` by setting up the color filler with a pattern applied to the `color` memory slice.
func (r *ColorRam) Setup() error {
	r.colorFiller.InitWithPattern(r.color, uint(len(r.color)))
	return nil
}

// Bind establishes a connection between the ColorRam component and the provided IColorRamC64Socket interface.
func (r *ColorRam) Bind(_ references.IColorRamC64Socket) error {
	return nil
}

// Connect initializes and establishes connections for the ColorRam component, preparing it for proper operation.
func (r *ColorRam) Connect() error {
	return nil
}

// Reset clears and reinitializes the internal color RAM to its default state.
func (r *ColorRam) Reset() {
}

// EmulationRequired determines if emulation logic is required for the ColorRam component. Always returns false.
func (r *ColorRam) EmulationRequired() bool {
	return false
}

// Emulate performs the emulation process for the ColorRam component during the system's operational cycle.
func (r *ColorRam) Emulate() {
}

// Internal determines if the ColorRam is configured for internal use. Always returns false.
func (r *ColorRam) Internal() bool {
	return false
}

// Read retrieves a 4-bit color value from the color RAM at the specified 10-bit address.
func (r *ColorRam) Read(addr uint16) uint8 {
	return r.color[addr&0x03ff]
}

// Write stores the lower 4 bits of the provided data at the address masked to fit the 1K memory range.
func (r *ColorRam) Write(addr uint16, data uint8) {
	r.color[addr&0x03ff] = data & 0x0f
}

// ReadRegister reads the 8-bit value from the color SRAM at the specified address, masking the address to 10 bits.
func (r *ColorRam) ReadRegister(addr uint16) uint8 {
	return r.color[addr&0x03ff]
}

// WriteRegister writes the provided 4-bit data to the specified address in the color RAM, masking the address to 10 bits.
func (r *ColorRam) WriteRegister(addr uint16, data uint8) {
	r.color[addr&0x03ff] = data & 0x0f
}

// Size returns the current size of the color memory array.
func (r *ColorRam) Size() int {
	return len(r.color)
}
