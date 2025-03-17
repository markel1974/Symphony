package mos6581

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// potXRegisterIndex represents the register index for the X-axis potentiometer.
// potYRegisterIndex represents the register index for the Y-axis potentiometer.
const (
	potXRegisterIndex = 25
	potYRegisterIndex = 26
)

// SID represents a chip emulation containing configurations, registers, and audio handling functionality.
type SID struct {
	*component.BaseComponent
	factory      references.IComponentFactory
	socket       references.ISIDSocket
	registers    []uint8
	cfg          *config.Config
	audioBuilder *AudioBuilder
	reflect      *SidReflect
}

// NewSID creates a new SID instance with a specified parent ID and suffix, initializing its registers and settings.
func NewSID(parent references.IComponent, factory references.IComponentFactory, suffix string) *SID {
	s := &SID{
		BaseComponent: component.NewBaseComponent(componentId, suffix),
		factory:       factory,
		socket:        nil,
		registers:     make([]uint8, RegisterCount),
		cfg:           nil,
		audioBuilder:  nil,
	}
	component.Register(parent, s)
	s.reflect = NewSidReflect(s)
	return s
}

// Setup initializes the SID instance with the provided socket, configuration, fragment frequency, and raster count.
func (sid *SID) Setup(socket references.ISIDSocket, fragFreq int, rasters int, cfg *config.Config) error {
	sid.socket = socket
	sid.audioBuilder = NewAudioBuilder(sid.socket.GetPlayer(), true, fragFreq, rasters)
	sid.cfg = cfg
	sid.cfg.Bind(sid.onConfigChanged)

	return nil
}

// Emulate processes the main emulation logic for the SID component, handling internal updates and state changes.
func (sid *SID) Emulate() {

}

// SetPotX sets the value of the POT X register in the SID chip using the given 8-bit value.
func (sid *SID) SetPotX(pot uint8) {
	// PX7 PX6 PX5 PX4 PX3 PX2 PX1 PX0
	sid.registers[potXRegisterIndex] = pot
}

// SetPotY sets the value of the POT Y register to the specified 8-bit value in the SID chip.
func (sid *SID) SetPotY(pot uint8) {
	//PY7 PY6 PY5 PY4 PY3 PY2 PY1 PY0
	sid.registers[potYRegisterIndex] = pot
}

// onConfigChanged is triggered when the configuration bound to the SID instance changes.
func (sid *SID) onConfigChanged() {
	//TODO
}

// Reset initializes all SID registers to 0 and sets default values for PotX and PotY. It also resets the audio builder.
func (sid *SID) Reset() {
	for x := range sid.registers {
		sid.registers[x] = 0
	}
	sid.SetPotX(0xff)
	sid.SetPotY(0xff)

	sid.audioBuilder.Reset()

	//PADDLE TEST
	//sid.WriteRegister(0xDC00, 0x40) //Control-Port 1 selected
	//sid.WriteRegister(0xD419, 0x7F) //Paddle X value
	//sid.WriteRegister(0xD419, 0x7F) //Paddle Y value
}

// ReadRegister retrieves the value from the specified address within the SID's registers. Only the lower 5 bits of the address are used.
func (sid *SID) ReadRegister(addr uint16) uint8 {
	reg := addr & 0x1f
	v := sid.registers[reg]
	return v
}

// WriteRegister writes an 8-bit value to a specific register at the given address by mapping it within a 32-register range.
func (sid *SID) WriteRegister(addr uint16, data uint8) {
	reg := addr & 0x1f
	sid.registers[reg] = data
}

// Prepare loads necessary SID register values into the AudioBuilder for audio processing.
func (sid *SID) Prepare() {
	for _, x := range _audioRegisters {
		sid.audioBuilder.LoadRegister(x, sid.registers[x])
	}
}

// Update triggers the audioBuilder's internal Update method, updating audio sampling and processing within the SID.
func (sid *SID) Update() {
	sid.audioBuilder.Update()
}

// GetLastByte retrieves the last byte from the SID's internal state or configuration.
func (sid *SID) GetLastByte() uint8 {
	return 0
}
