package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// IExpansionSocketConnections defines an interface for managing connections and triggers in an expansion socket system.
// DMALowTrigger sets the DMA low line state to the specified boolean value.
// IRQTriggerBind binds a callback function to handle IRQ triggers with a uint32 parameter.
// IRQClearBind binds a callback function to handle IRQ clearing with a uint32 parameter.
type IExpansionSocketConnections interface {
	DMALowTrigger(v bool)
}

// ExpansionSocket represents a connection point integrating various components like connections, PIC, PLA, VIC, and Quartz.
// It enables coordinated interaction and communication between connected emulation subsystems.
type ExpansionSocket struct {
	connections IExpansionSocketConnections
	pic         references.IPIC6510
	pla         references.IPlaC64
	vic         references.IVIC
	quartz      references.IQuartz
}

// NewExpansionSocket initializes and returns a pointer to a new ExpansionSocket instance with default nil values.
func NewExpansionSocket(connections IExpansionSocketConnections) *ExpansionSocket {
	e := &ExpansionSocket{
		connections: connections,
		pic:         nil,
		pla:         nil,
		vic:         nil,
		quartz:      nil,
	}
	return e
}

// Setup initializes the ExpansionSocket with its dependencies and sets up the required connections.
func (s *ExpansionSocket) Setup(c map[string]references.IComponent, _ *config.Config) error {
	var err error
	s.pic, err = references.ComponentsToIPIC6510(c, 0)
	if err != nil {
		return err
	}
	s.pla, err = references.ComponentsToIPLAc64(c, 0)
	if err != nil {
		return err
	}
	s.vic, err = references.ComponentsToIVIC(c, 0)
	if err != nil {
		return err
	}
	s.quartz, err = references.ComponentsToIQuartz(c, 0)
	if err != nil {
		return err
	}
	return nil
}

func (w *ExpansionSocket) Connect() error {
	return nil
}

// Reset reinitializes the state of the ExpansionSocket and its connected components to their default values.
func (s *ExpansionSocket) Reset() {
}

// Read performs a read operation from the specified 16-bit memory address and returns the corresponding 8-bit value.
func (s *ExpansionSocket) Read(addr uint16) uint8 {
	return s.pla.Read(addr)
}

// Write performs a write operation to the specified memory address with the given data value.
func (s *ExpansionSocket) Write(addr uint16, data uint8) {
	s.pla.Write(addr, data)
}

// GameExRomConfigChanged updates the memory configuration of the PLA by triggering a rebuild process.
func (s *ExpansionSocket) GameExRomConfigChanged() {
	s.pla.RebuildMemoryConfig()
}

// NMITrigger triggers a non-maskable interrupt (NMI) by invoking the corresponding method on the programmable interrupt controller.
func (s *ExpansionSocket) NMITrigger() {
	s.pic.TriggerNMI()
}

// SetDMALow triggers the DMA low state on the associated expansion socket connection based on the boolean value provided.
func (s *ExpansionSocket) SetDMALow(v bool) {
	s.connections.DMALowTrigger(v)
}

// ResetTrigger invokes a reset trigger on the connected IPIC6510 instance to reinitialize its state.
func (s *ExpansionSocket) ResetTrigger() {
	s.pic.TriggerReset()
}

// IRQTrigger triggers an interrupt request (IRQ) using the programmable interrupt controller (PIC).
func (s *ExpansionSocket) IRQTrigger() {
	s.pic.TriggerIRQ(intrIrqExpansionBit)
}

// IRQClear clears the IRQ signal associated with the expansion bit from the programmable interrupt controller (PIC).
func (s *ExpansionSocket) IRQClear() {
	s.pic.ClearIRQ(intrIrqExpansionBit)
}

// IRQTriggerBind connects a callback function to the IRQ trigger event, enabling custom handling of IRQ signals.
func (s *ExpansionSocket) IRQTriggerBind(fn func(uint32)) {
	s.pic.IRQTriggerBind(fn)
}

// IRQClearBind binds a callback function that is triggered when the IRQ clear event occurs.
func (s *ExpansionSocket) IRQClearBind(fn func(uint32)) {
	s.pic.IRQClearBind(fn)
}

// BusAvailable checks whether the expansion bus is available by verifying the BA (Bus Available) line status of the VIC.
func (s *ExpansionSocket) BusAvailable() bool {
	return !s.vic.GetBALow()
}

// AECAvailable checks if the AEC (Address Enable Control) line is available and returns true if it's not low.
func (s *ExpansionSocket) AECAvailable() bool {
	return !s.vic.GetAECLow()
}

// Cycle retrieves the current clock cycle count from the associated quartz instance.
func (s *ExpansionSocket) Cycle() uint64 {
	return s.quartz.Cycle()
}

// CycleAlarm creates a new quartz alarm with the given id and callback, and returns the associated IQuartzAlarm instance.
func (s *ExpansionSocket) CycleAlarm(id string, callback references.QuartzAlarmCallback) references.IQuartzAlarm {
	return s.quartz.NewAlarm(id, callback)
}

// RamSetWriteTrigger sets a write trigger callback for a specified memory address, returning an identifier for the trigger.
func (s *ExpansionSocket) RamSetWriteTrigger(addr uint16, fn func(uint16, uint8)) int {
	return s.pla.SetWriteTrigger(addr, fn)
}

// RamRemoveWriteTrigger removes a write trigger callback associated with the specified address and identifier.
func (s *ExpansionSocket) RamRemoveWriteTrigger(addr uint16, id int) {
	s.pla.RemoveRamTrigger(addr, id)
}

// RmwFlags retrieves the Read-Modify-Write flags for CPU operations. Currently, this method is a placeholder for implementation.
func (s *ExpansionSocket) RmwFlags() uint8 {
	//TODO IMPLEMENT cpu rmw flags
	return 0
}
