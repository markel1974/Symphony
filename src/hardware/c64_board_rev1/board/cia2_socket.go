package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// CIA2SocketConnection defines an interface for managing trigger and clear operations for CIA2 socket connections.
// NMITrigger triggers a Non-Maskable Interrupt (NMI) signal in the CIA2 socket.
// NMIClearTrigger clears the Non-Maskable Interrupt (NMI) signal in the CIA2 socket.
type CIA2SocketConnection interface {
	NMITrigger()
	NMIClearTrigger()
}

// CIA2Socket represents a specialized implementation of IMos6526, handling interactions with CIA2 and emulation components.
// The connections field manages CIA2-specific port and communication events.
// The iec field facilitates communication over the IEC serial bus.
// The intrId field defines the unique interrupt identifier for CIA2 within the system.
type CIA2Socket struct {
	references.IMos6526
	label                      string
	parent                     references.IComponent
	component                  references.IComponent
	vicRef                     references.IMos6569
	iecRef                     references.IIec
	intrId                     uint32
	hwId                       string
	selfReadDDRA               func() uint8
	selfReadPRA                func() uint8
	connectionsNMITrigger      func()
	connectionsNMIClearTrigger func()
	vicChangedVA               func(uint8)
	iecCpuRead                 func() uint8
	iecCpuWrite                func(uint8)
}

// NewCIA2Socket creates and returns a pointer to a new instance of CIA2Socket with default uninitialized fields.
func NewCIA2Socket(parent references.IComponent, label string, connections CIA2SocketConnection) *CIA2Socket {
	c := &CIA2Socket{
		parent:                     parent,
		label:                      label,
		connectionsNMITrigger:      connections.NMITrigger,
		connectionsNMIClearTrigger: connections.NMIClearTrigger,
		IMos6526:                   nil,
		vicRef:                     nil,
		iecRef:                     nil,
		intrId:                     intrIrqCia2Bit,
	}
	c.hwId = references.IdIMos6526(c.IMos6526, c.label, 1)
	return c
}

func (w *CIA2Socket) HardwareId() string {
	return w.hwId
}

// Wire initializes the CIA2Socket instance with the provided CIA, connections, and IEC interface, and sets up the CIA.
func (w *CIA2Socket) Wire() error {
	var err error
	w.component = w.parent.GetChildByHardwareId(w.HardwareId())
	if w.IMos6526, err = references.ComponentToIMos6526(w.component); err != nil {
		return err
	}
	idVIC := references.IdIMos6569(w.vicRef, w.label, 0)
	if w.vicRef, err = references.ComponentToIMos6569(w.parent.GetChildByHardwareId(idVIC)); err != nil {
		return err
	}
	idIEC := references.IdIIec(w.iecRef, w.label, 0)
	if w.iecRef, err = references.ComponentToIEC(w.parent.GetChildByHardwareId(idIEC)); err != nil {
		return err
	}
	if err = w.IMos6526.Bind(w); err != nil {
		return err
	}
	w.selfReadDDRA = w.ReadDDRA
	w.selfReadPRA = w.ReadPRA
	w.vicChangedVA = w.vicRef.ChangedVA
	w.iecCpuRead = w.iecRef.CpuRead
	w.iecCpuWrite = w.iecRef.CpuWrite
	return nil
}

// ReadPortA reads the value of port A by combining peripheral data, data direction bits, and the IEC CpuRead result.
func (w *CIA2Socket) ReadPortA(prA uint8 /* prb */, _ uint8, ddrA uint8 /*ddrB */, _ uint8) uint8 {
	data := w.iecCpuRead()
	ret := ((prA | (^ddrA)) & 0x3f) | data
	return ret
}

// ReadPortB reads the current state of Port B by combining the peripheral register and inverted direction register.
func (w *CIA2Socket) ReadPortB( /* prA */ _ uint8, prB uint8 /* ddrA */, _ uint8, ddrB uint8) uint8 {
	ret := prB | (^ddrB)
	return ret
}

// SignalPRA writes the state of Port A using the given peripheral and direction register values.
func (w *CIA2Socket) SignalPRA(prA uint8) {
	ddrA := w.selfReadDDRA()
	w.updateVA(prA, ddrA)
	w.iecCpuWrite(prA)
}

// SignalPRB updates the state of port B using the provided peripheral and direction registers.
func (w *CIA2Socket) SignalPRB(_ uint8) {
}

// SignalDDRA updates the data direction register for port A and triggers actions based on the updated state.
func (w *CIA2Socket) SignalDDRA(ddrA uint8) {
	prA := w.selfReadPRA()
	w.updateVA(prA, ddrA)
}

// SignalDDRB updates the data direction register for port B with the provided parameters.
func (w *CIA2Socket) SignalDDRB(_ uint8) {
}

// ReadSP reads the state of the SP (Serial Port) line and returns its current boolean value.
func (w *CIA2Socket) ReadSP() bool {
	return false
}

// SignalSP sets the state of the SP (Serial Port) line to the specified level, controlling serial communication output.
func (w *CIA2Socket) SignalSP( /*level*/ _ bool) {
}

// Update triggers internal state updates and recalculations based on the latest system inputs or interactions.
//func (w *CIA2Socket) Update() {
//}

// updateVA updates the VIC-memory bank based on the current states of prA and ddrA and triggers a corresponding VA change event.
func (w *CIA2Socket) updateVA(prA uint8, ddrA uint8) {
	//Bit 0..1: Select the position of the VIC-memory
	//Bit 2: RS-232: TXD Output, userPort: Data PA 2 (pin M)
	//Bit 3..5: serial bus Output (0=High/Inactive, 1=Low/Active)
	//Bit 6..7: serial bus Input (0=Low/Active, 1=High/Inactive)

	//%00, 0: Bank 3: $C000-$FFFF, 49152-65535
	//%01, 1: Bank 2: $8000-$BFFF, 32768-49151
	//%10, 2: Bank 1: $4000-$7FFF, 16384-32767
	//%11, 3: Bank 0: $0000-$3FFF, 0-16383 (standard)
	va := (^(prA | (^ddrA))) & 3
	w.vicChangedVA(va)
}

// IRQTrigger triggers a Non-Maskable Interrupt (NMI) by invoking the NMITrigger method on the connected socket.
func (w *CIA2Socket) IRQTrigger() {
	w.connectionsNMITrigger()
}

// IRQClearTrigger clears the NMI (Non-Maskable Interrupt) request by invoking the NMIClear method on the connections object.
func (w *CIA2Socket) IRQClearTrigger() {
	w.connectionsNMIClearTrigger()
}
