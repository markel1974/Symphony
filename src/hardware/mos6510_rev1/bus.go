package mos6510_rev1

import (
	"github.com/markel1974/symphony/src/kernel/component"
	"github.com/markel1974/symphony/src/references"
)

type Bus struct {
	*component.BaseComponent
	reflect     *BusReflect
	status      uint8
	realRead    func(uint16) uint8
	realWrite   func(uint16, uint8)
	setModeHalt func()
	Read        func(uint16) (uint8, bool) // Read is a function that reads a byte from a specified 16-bit memory address in the CPU's memory bank.
	Write       func(uint16, uint8)        // Write is a function that writes a byte to a specified 16-bit memory address in the CPU's memory bank.
}

func NewBus(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *Bus {
	p := &Bus{
		status:        0,
		BaseComponent: component.NewBaseComponent(),
	}
	p.reflect = NewBusReflect(p, factory, parent, "bus", instance, references.IdInternalComponent(label, instance, "Bus"))
	return p
}

// Setup initializes the Interrupts instance with the provided quartz reference and other dependencies. Returns an error on failure.
func (i *Bus) Setup() error {
	return nil
}

func (i *Bus) Bind(setModeHalt func(), banks references.IMos6510Banks) error {
	i.setModeHalt = setModeHalt
	i.realRead = banks.Read
	i.realWrite = banks.Write
	i.SetNormalMode()
	return nil
}

// Connect establishes necessary bindings or configurations between the PIC and other components, returning an error if it fails.
func (i *Bus) Connect() error {
	return nil
}

// Internal indicates if the `Interrupts` is set as an internal device. Always returns false in this implementation.
func (i *Bus) Internal() bool {
	return true
}

// Emulate processes the internal state of the Interrupts instance during an emulation cycle.
func (i *Bus) Emulate() {
}

// EmulationRequired determines if emulation is required for the current state of the Interrupts component. Always returns false.
func (i *Bus) EmulationRequired() bool {
	return false
}

// Reset reinitializes the Interrupts by clearing all internal state and resetting signals to their default values.
func (i *Bus) Reset() {
	i.SetNormalMode()
}

func (i *Bus) ReadDirect(addr uint16) uint8 {
	return i.realRead(addr)
}

func (i *Bus) SetAECLowMode() {
	// disconnecting the bus...
	i.Read = i.readAECLow
	i.Write = i.writeAECLow
}

func (i *Bus) SetRDYLowMode() {
	i.Read = i.readBALow
	i.Write = i.realWrite
}

func (i *Bus) SetNormalMode() {
	//normal mode
	i.Read = i.readNormal
	i.Write = i.realWrite
}

// readNormal performs a normal read operation from the specified address and always returns a successful status.
//
//go:nosplit
func (i *Bus) readNormal(addr uint16) (uint8, bool) {
	data := i.realRead(addr)
	return data, true
}

// readBALow halts the CPU by transitioning to halt mode and always returns 0 and false.
//
//go:nosplit
func (i *Bus) readBALow(_ uint16) (uint8, bool) {
	i.setModeHalt()
	return 0, false
}

// readAec performs a read operation while the CPU bus is disconnected, always returning a default value of 0.
//
//go:nosplit
func (i *Bus) readAECLow(_ uint16) (uint8, bool) {
	return 0, false
}

// writeAec is a placeholder method called when an illegal write operation is attempted on a disconnected bus.
//
//go:nosplit
func (i *Bus) writeAECLow(_ uint16, _ uint8) {
}
