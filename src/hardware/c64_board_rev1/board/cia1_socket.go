package board

import (
	"github.com/markel1974/c64emu/src/common/bits"
	"github.com/markel1974/c64emu/src/references"
)

// defaultLPState represents the default state value for Launchpad.
// defaultJoyState represents the default joystick state value.
// defaultKeyState represents the default keyboard state value.
const (
	defaultLPState  = 0x10
	defaultJoyState = 0xff
	defaultKeyState = 0xff
)

// CIA1SocketConnection is an interface for handling external interactions like IRQ triggers for CIA-1 components.
// IRQTrigger simulates an Interrupt Request (IRQ) activation for the given event identifier.
// IRQClearTrigger simulates clearing an IRQ activation for the provided event identifier.
type CIA1SocketConnection interface {
	IRQTrigger(d uint32)
	IRQClearTrigger(d uint32)
}

// CIA1Socket represents a connector for CIA-1 chip emulation, managing keyboard, joystick, and external connections.
// IMos6526 interface provides the core logic for the connected CIA chip functionality.
// IC64Keyboard interface is used to handle input and state changes from a keyboard.
// IC64Joystick interface specifies methods to interact with joystick devices (one or more).
// ICIA1SocketConnections interface provides external event interactions like IRQ and light pen triggers.
// intrId is used to track the IRQ identification bit for events.
// prevLPState represents the prior state of the light pen (LP) input signal.
// keyMatrix is a direct lookup table representing the state of the keyboard matrix.
// revMatrix is the reversed lookup table to assist with keyboard matrix computation.
// joy1State and joy2State track the current states of two connected joysticks.
type CIA1Socket struct {
	references.IMos6526
	label       string
	parent      references.IComponent
	component   references.IComponent
	connection  CIA1SocketConnection
	vic         references.IMos6569
	keys        references.IC64Keyboard
	joy1        references.IC64Joystick
	joy2        references.IC64Joystick
	intrId      uint32  //
	prevLPState uint8   // Previous state of LP line (bit 4)
	keyMatrix   []uint8 // keyboard matrix [0: down, 1: up]
	revMatrix   []uint8 // Reversed keyboard matrix
	joy1State   uint8   // Joystick 1
	joy2State   uint8   // Joystick 2
	hwId        string
}

// NewCIA1Socket creates and initializes a new instance of CIA1Socket with default state and properties.
func NewCIA1Socket(parent references.IComponent, label string, connection CIA1SocketConnection) *CIA1Socket {
	c := &CIA1Socket{
		parent:      parent,
		label:       label,
		connection:  connection,
		IMos6526:    nil,
		vic:         nil,
		keys:        nil,
		joy1:        nil,
		joy2:        nil,
		intrId:      intrIrqCia1Bit,
		prevLPState: defaultLPState,
		keyMatrix:   make([]uint8, 8),
		revMatrix:   make([]uint8, 8),
		joy1State:   defaultJoyState,
		joy2State:   defaultJoyState,
	}
	c.hwId = references.IdIMos6526(c.IMos6526, c.label, 0)
	return c
}

func (w *CIA1Socket) HardwareId() string {
	return w.hwId
}

// Mount initializes the CIA1Socket with the provided CIA instance, connections, keyboard, and joystick references.
// It sets up the CIA via the Setup method and returns any errors encountered during initialization.
func (w *CIA1Socket) Mount() error {
	var err error
	w.component = w.parent.GetChildByHardwareId(w.HardwareId())
	if w.IMos6526, err = references.ComponentToIMos6526(w.component); err != nil {
		return err
	}
	idVIC := references.IdIMos6569(w.vic, w.label, 0)
	if w.vic, err = references.ComponentToIMos6569(w.parent.GetChildByHardwareId(idVIC)); err != nil {
		return err
	}
	idKeys := references.IdIC64Keyboard(w.keys, w.label, 0)
	if w.keys, err = references.ComponentToIC64Keyboard(w.parent.GetChildByHardwareId(idKeys)); err != nil {
		return err
	}
	idJoy1 := references.IdIC64Joystick(w.joy1, w.label, 0)
	if w.joy1, err = references.ComponentToIC64Joystick(w.parent.GetChildByHardwareId(idJoy1)); err != nil {
		return err
	}
	idJoy2 := references.IdIC64Joystick(w.joy1, w.label, 1)
	if w.joy2, err = references.ComponentToIC64Joystick(w.parent.GetChildByHardwareId(idJoy2)); err != nil {
		return err
	}
	if err = w.IMos6526.Bind(w); err != nil {
		return err
	}
	return nil
}

// Reset resets the internal state of the CIA1Socket and its connected components to their default values.
func (w *CIA1Socket) Reset() {
	w.keys.Reset()
	w.joy1.Reset()
	w.joy2.Reset()
	w.IMos6526.Reset()
	for idx := range w.keyMatrix {
		w.keyMatrix[idx] = defaultKeyState
	}
	for idx := range w.revMatrix {
		w.revMatrix[idx] = defaultKeyState
	}
	w.joy1State = defaultJoyState
	w.joy2State = defaultJoyState
	w.prevLPState = defaultLPState
}

// Update polls the state of connected joysticks and keyboard, and updates the key and reverse matrix based on input events.
func (w *CIA1Socket) Update() {
	if joy1State, ok := w.joy1.Poll(); ok {
		w.joy1State = joy1State
	}
	if joy2State, ok := w.joy2.Poll(); ok {
		w.joy2State = joy2State
	}
	if v, ok := w.keys.Poll(); ok {
		pressed := (v & 0x20000) != 0
		shifted := (v & 0x10000) != 0
		keyM := uint8(v & 0xff)
		revM := uint8((v >> 8) & 0xff)
		if pressed {
			if shifted {
				w.keyMatrix[6] &= 0xef
				w.revMatrix[4] &= 0xbf
			}
			w.keyMatrix[keyM] &= ^(1 << revM)
			w.revMatrix[revM] &= ^(1 << keyM)
		} else {
			if shifted {
				w.keyMatrix[6] |= 0x10
				w.revMatrix[4] |= 0x40
			}
			w.keyMatrix[keyM] |= 1 << revM
			w.revMatrix[revM] |= 1 << keyM
		}
	}
	w.IMos6526.Update()
}

// ReadPortA reads data from port A by combining the given parameters with internal joystick states and matrices.
func (w *CIA1Socket) ReadPortA(prA uint8, ddrA uint8, prB uint8, ddrB uint8) uint8 {
	ret := prA | ^ddrA
	tst := (prB | ^ddrB) & w.joy1State
	for idx, bit := range bits.Uint8s {
		if (tst & bit) == 0 {
			ret &= w.revMatrix[idx]
		}
	}
	return ret & w.joy2State
}

// ReadPortB reads the state of Port B by combining the active bits of DDRB and PRB with the joystick and key matrix states.
func (w *CIA1Socket) ReadPortB(prA uint8, ddrA uint8, prB uint8, ddrB uint8) uint8 {
	ret := ^ddrB
	tst := (prA | ^ddrA) & w.joy2State
	for idx, bit := range bits.Uint8s {
		if (tst & bit) == 0 {
			ret &= w.keyMatrix[idx]
		}
	}
	return (ret | (prB & ddrB)) & w.joy1State
}

// WritePortA handles data writing to port A based on the current state and provided parameters.
func (w *CIA1Socket) WritePortA(_ uint8, _ uint8, _ uint8, _ uint8) {
}

// WriteDdrA updates the Data Direction Register A (DDRA) with the provided parameters without performing any operation.
func (w *CIA1Socket) WriteDdrA(_ uint8, _ uint8, _ uint8, _ uint8) {
}

// WritePortB handles writing to port B by updating the light pen state based on the given port and data direction registers.
func (w *CIA1Socket) WritePortB(_ uint8, _ uint8, prB uint8, ddrB uint8) {
	w.updateLightPen(prB, ddrB)
}

// WriteDdrB handles updates to the DDRB register and triggers updates to the light pen based on the PRB and DDRB values.
func (w *CIA1Socket) WriteDdrB(_ uint8, _ uint8, prB uint8, ddrB uint8) {
	w.updateLightPen(prB, ddrB)
}

func (w *CIA1Socket) ReadSP() bool {
	//TODO ATTACH
	return false
}

func (w *CIA1Socket) WriteSP(level bool) {
	//TODO ATTACH
}

// IRQTrigger triggers an interrupt request (IRQ) using the associated interrupt ID managed by the connection interface.
func (w *CIA1Socket) IRQTrigger() {
	w.connection.IRQTrigger(w.intrId)
}

// IRQClearTrigger clears the interrupt request associated with the socket by invoking the IRQClearTrigger method on connections.
func (w *CIA1Socket) IRQClearTrigger() {
	w.connection.IRQClearTrigger(w.intrId)
}

// updateLightPen checks the state of the light pen line and triggers an action if the state has changed.
func (w *CIA1Socket) updateLightPen(prB uint8, ddrB uint8) {
	if ((prB | ^ddrB) & 0x10) != w.prevLPState {
		w.vic.LightPenTrigger()
	}
	w.prevLPState = (prB | ^ddrB) & 0x10
}
