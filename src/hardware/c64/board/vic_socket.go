package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

type IVICSocketConnection interface {
	LastCycleTrigger()

	VBlankTrigger()

	RDYLowTrigger(v bool)

	AECLowTrigger(v bool)
}

// VICSocket represents a virtual interface connector socket with a reference to a board and an interrupt identifier.
type VICSocket struct {
	references.IVIC
	connections IVICSocketConnection
	db          references.IDisplayBuffer
	pic         references.IPIC6510
	pla         references.IPlaC64
	quartz      references.IQuartz
	intrId      uint32
}

// NewVICSocket creates and returns a new instance of the VICSocket struct.
func NewVICSocket(connections IVICSocketConnection) *VICSocket {
	return &VICSocket{
		IVIC:        nil,
		connections: connections,
		db:          nil,
		pic:         nil,
		pla:         nil,
		quartz:      nil,
		intrId:      intrIrqVicBit,
	}
}

func (v *VICSocket) SetDisplayBuffer(db references.IDisplayBuffer) {
	v.db = db
}

// GetDisplayBuffer returns the IDisplayBuffer instance associated with the VICSocket's board.
func (v *VICSocket) GetDisplayBuffer() references.IDisplayBuffer {
	return v.db
}

// Setup initializes the VICSocket with the given board and interrupt ID.
func (v *VICSocket) Setup(c map[string]references.IComponent, cfg *config.Config) error {
	var err error
	if v.IVIC, err = references.ComponentsToIVIC(c, 0); err != nil {
		return err
	}
	if v.pic, err = references.ComponentsToIPIC6510(c, 0); err != nil {
		return err
	}
	if v.pla, err = references.ComponentsToIPLAc64(c, 0); err != nil {
		return err
	}
	if v.quartz, err = references.ComponentsToIQuartz(c, 0); err != nil {
		return err
	}
	if err = v.IVIC.Setup(v, cfg); err != nil {
		return err
	}
	return nil
}

func (v *VICSocket) Connect() error {
	return nil
}

// Cycle retrieves the current cycle count from the associated Quartz scheduler.
func (v *VICSocket) Cycle() uint64 {
	return v.quartz.Cycle()
}

// GetBanks returns an implementation of the mos6569.IVICBanks interface, which provides access to memory handling operations.
func (v *VICSocket) GetBanks() references.IVICBanks {
	return v.pla
}

// IRQTrigger signals an interrupt request by invoking the IRQ trigger mechanism on the associated board slot.
func (v *VICSocket) IRQTrigger() {
	v.pic.TriggerIRQ(v.intrId)
}

// IRQClear clears the interrupt request for the associated slot identified by the intrId of the VICSocket.
func (v *VICSocket) IRQClear() {
	v.pic.ClearIRQ(v.intrId)
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
