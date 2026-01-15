package board

import (
	"github.com/markel1974/symphony/src/references"
)

// ICartridgeManagerConnections defines the interface for managing hardware connection signals for a cartridge system.
// BALow checks the state of the Bus Available (BA) signal, returning true if low.
// AECLow checks the state of the Address Enable Control (AEC) signal, returning true if low.
// NMITrigger triggers a Non-Maskable Interrupt (NMI) signal in the system.
// IRQTrigger sends an Interrupt Request (IRQ) signal with a specified delay value.
// IRQClearTrigger clears a triggered IRQ signal after the specified delay value.
// RSTTrigger triggers a system Reset (RST) signal in the hardware.
// DMALowTrigger controls the state of the Direct Memory Access (DMA) signal, accepting a boolean value.
type ICartridgeManagerConnections interface {
	BALow() bool

	AECLow() bool

	NMITrigger()

	IRQTrigger(d uint32)

	IRQClearTrigger(d uint32)

	RSTTrigger()

	DMALowTrigger(v bool)

	LedActivityTrigger(device uint8, led bool)
}

// CartridgeManagerSocket represents a connection point integrating various components like connections, PIC, PLA, VIC, and Quartz.
// It enables coordinated interaction and communication between connected emulation subsystems.
type CartridgeManagerSocket struct {
	references.IC64CartridgeManager
	connections ICartridgeManagerConnections
	label       string
	parent      references.IComponent
	component   references.IComponent
	pla         references.IC64Pla
	quartz      references.IQuartz
	hwId        string
}

// NewExpansionSocket initializes and returns a pointer to a new CartridgeManagerSocket instance with default nil values.
func NewExpansionSocket(parent references.IComponent, label string, connections ICartridgeManagerConnections) *CartridgeManagerSocket {
	e := &CartridgeManagerSocket{
		parent:      parent,
		label:       label,
		connections: connections,
		pla:         nil,
		quartz:      nil,
	}
	e.hwId = references.IdIC64CartridgeManager(e.IC64CartridgeManager, e.label, 0)
	return e
}

func (s *CartridgeManagerSocket) HardwareId() string {
	return s.hwId
}

// Wire initializes the CartridgeManagerSocket with its dependencies and sets up the required connections.
func (s *CartridgeManagerSocket) Wire() error {
	var err error
	s.component = s.parent.GetChildByHardwareId(s.HardwareId())
	if s.IC64CartridgeManager, err = references.ComponentToIC64CartridgeManager(s.component); err != nil {
		return err
	}
	if err = s.IC64CartridgeManager.Bind(s, s); err != nil {
		return err
	}
	idIPLA := references.IdIC64Pla(s.pla, s.label, 0)
	if s.pla, err = references.ComponentToIC64Pla(s.parent.GetChildByHardwareId(idIPLA)); err != nil {
		return err
	}
	idQuartz := references.IdIQuartz(s.quartz, s.label, 0)
	if s.quartz, err = references.ComponentToIQuartz(s.parent.GetChildByHardwareId(idQuartz)); err != nil {
		return err
	}
	return nil
}

func (s *CartridgeManagerSocket) Connect() error {
	return nil
}

// Read performs a read operation from the specified 16-bit memory address and returns the corresponding 8-bit value.
func (s *CartridgeManagerSocket) Read(addr uint16) uint8 {
	return s.pla.Read(addr)
}

// Write performs a write operation to the specified memory address with the given data value.
func (s *CartridgeManagerSocket) Write(addr uint16, data uint8) {
	s.pla.Write(addr, data)
}

// GameExRomConfigChanged updates the memory configuration of the PLA by triggering a rebuild process.
func (s *CartridgeManagerSocket) GameExRomConfigChanged() {
	s.pla.RebuildMemoryConfig()
}

// NMITrigger triggers a non-maskable interrupt (NMI) by invoking the corresponding method on the programmable interrupt controller.
func (s *CartridgeManagerSocket) NMITrigger() {
	s.connections.NMITrigger()
}

// SetDMALow triggers the DMA low state on the associated expansion socket connection based on the boolean value provided.
func (s *CartridgeManagerSocket) SetDMALow(v bool) {
	s.connections.DMALowTrigger(v)
}

// ResetTrigger invokes a reset trigger on the connected IMos6510Pic instance to reinitialize its state.
func (s *CartridgeManagerSocket) ResetTrigger() {
	s.connections.RSTTrigger()
}

// IRQTrigger triggers an interrupt request (IRQ) using the programmable interrupt controller (PIC).
func (s *CartridgeManagerSocket) IRQTrigger() {
	s.connections.IRQTrigger(intrIrqExpansionBit)
}

// IRQClearTrigger clears the IRQ signal associated with the expansion bit from the programmable interrupt controller (PIC).
func (s *CartridgeManagerSocket) IRQClearTrigger() {
	s.connections.IRQClearTrigger(intrIrqExpansionBit)
}

// BusAvailable checks whether the expansion bus is available by verifying the BA (Bus Available) line status of the VIC.
func (s *CartridgeManagerSocket) BusAvailable() bool {
	return !s.connections.BALow()
}

// AECAvailable checks if the AEC (Address Enable Control) line is available and returns true if it's not low.
func (s *CartridgeManagerSocket) AECAvailable() bool {
	return !s.connections.AECLow()
}

// Cycle retrieves the current clock cycle count from the associated quartz instance.
func (s *CartridgeManagerSocket) Cycle() uint64 {
	return s.quartz.Cycle()
}

// CycleAlarm creates a new quartz alarm with the given id and callback, and returns the associated IQuartzAlarm instance.
func (s *CartridgeManagerSocket) CycleAlarm(id string, callback references.QuartzAlarmCallback) references.IQuartzAlarm {
	return s.quartz.NewAlarm(id, callback)
}

// RamSetWriteTrigger sets a write trigger callback for a specified memory address, returning an identifier for the trigger.
func (s *CartridgeManagerSocket) RamSetWriteTrigger(addr uint16, fn func(uint16, uint8)) int {
	return s.pla.SetWriteTrigger(addr, fn)
}

// RamRemoveWriteTrigger removes a write trigger callback associated with the specified address and identifier.
func (s *CartridgeManagerSocket) RamRemoveWriteTrigger(addr uint16, id int) {
	s.pla.RemoveRamTrigger(addr, id)
}

func (s *CartridgeManagerSocket) LedActivity(deviceNumber uint8, led bool) {
	s.connections.LedActivityTrigger(deviceNumber, led)
}
