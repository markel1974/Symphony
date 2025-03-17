package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

type IVICSocketConnection interface {
	IRQTrigger(v uint32)

	IRQClear(v uint32)

	RDYLowTrigger(d bool)

	AECLowTrigger(d bool)

	LastCycleTrigger()

	VBlankTrigger()
}

// VICSocket represents a virtual interface connector socket with a reference to a board and an interrupt identifier.
type VICSocket struct {
	references.IVIC
	connections IVICSocketConnection
	db          references.IDisplayBuffer
	pla         references.IPlaC64
	quartz      references.IQuartz
	intrId      uint32
}

// NewVICSocket creates and returns a new instance of the VICSocket struct.
func NewVICSocket() *VICSocket {
	return &VICSocket{
		IVIC:        nil,
		connections: nil,
		db:          nil,
		quartz:      nil,
		intrId:      intrIrqVicBit,
	}
}

// Connect initializes the VICSocket with the given board and interrupt ID.
func (v *VICSocket) Connect(vic references.IVIC, connections IVICSocketConnection, db references.IDisplayBuffer, pla references.IPlaC64, quartz references.IQuartz, cfg *config.Config) error {
	v.IVIC = vic
	v.connections = connections
	v.db = db
	v.pla = pla
	v.quartz = quartz
	v.IVIC.Setup(v, cfg)
	return nil
}

// Cycle retrieves the current cycle count from the associated Quartz scheduler.
func (v *VICSocket) Cycle() uint64 {
	return v.quartz.Cycle()
}

// GetDisplayBuffer returns the IDisplayBuffer instance associated with the VICSocket's board.
func (v *VICSocket) GetDisplayBuffer() references.IDisplayBuffer {
	return v.db
}

// GetBanks returns an implementation of the mos6569.IVICBanks interface, which provides access to memory handling operations.
func (v *VICSocket) GetBanks() references.IVICBanks {
	return v.pla
}

// IRQTrigger signals an interrupt request by invoking the IRQ trigger mechanism on the associated board slot.
func (v *VICSocket) IRQTrigger() {
	v.connections.IRQTrigger(v.intrId)
}

// IRQClear clears the interrupt request for the associated slot identified by the intrId of the VICSocket.
func (v *VICSocket) IRQClear() {
	v.connections.IRQClear(v.intrId)
}

// BALow sets the BA (Bus Available) line low or high based on the provided boolean value.
func (v *VICSocket) BALow(d bool) {
	v.connections.RDYLowTrigger(d)
}

// AECLow controls the state of the Address Enable Control (AEC) line. It sets the AEC signal to low if d is true.
func (v *VICSocket) AECLow(d bool) {
	v.connections.AECLowTrigger(d)
}

// LastCycle triggers the last cycle slot operation on the VIC through the connected board.
func (v *VICSocket) LastCycle() {
	v.connections.LastCycleTrigger()
}

// VBlank handles the vertical blanking phase by triggering the corresponding slot on the associated board instance.
func (v *VICSocket) VBlank() {
	v.connections.VBlankTrigger()
}
