package board

import (
	"github.com/markel1974/symphony/src/common/bits"
	"github.com/markel1974/symphony/src/references"
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
	label                     string
	parent                    references.IComponent
	component                 references.IComponent
	vicRef                    references.IMos6569
	keysRef                   references.IC64Keyboard
	joy1Ref                   references.IC64Joystick
	joy2Ref                   references.IC64Joystick
	intrId                    uint32  //
	prevLPState               uint8   // Previous state of LP line (bit 4)
	keyMatrix                 []uint8 // keyboard matrix [0: down, 1: up]
	revMatrix                 []uint8 // Reversed keyboard matrix
	joy1State                 uint8   // Joystick 1
	joy2State                 uint8   // Joystick 2
	hwId                      string
	selfReadDDRB              func() uint8
	selfReadPRB               func() uint8
	connectionIRQTrigger      func(uint32)
	connectionIRQClearTrigger func(uint32)
	keysReset                 func()
	joy1Reset                 func()
	joy2Reset                 func()
	joy1Poll                  func() (uint8, bool)
	joy2Poll                  func() (uint8, bool)
	keysPoll                  func() (uint32, bool)
	vicLightPenTrigger        func()
}

// NewCIA1Socket creates and initializes a new instance of CIA1Socket with default state and properties.
func NewCIA1Socket(parent references.IComponent, label string, connection CIA1SocketConnection) *CIA1Socket {
	c := &CIA1Socket{
		parent:                    parent,
		label:                     label,
		connectionIRQTrigger:      connection.IRQTrigger,
		connectionIRQClearTrigger: connection.IRQClearTrigger,
		IMos6526:                  nil,
		vicRef:                    nil,
		keysRef:                   nil,
		joy1Ref:                   nil,
		joy2Ref:                   nil,
		intrId:                    intrIrqCia1Bit,
		prevLPState:               defaultLPState,
		keyMatrix:                 make([]uint8, 8),
		revMatrix:                 make([]uint8, 8),
		joy1State:                 defaultJoyState,
		joy2State:                 defaultJoyState,
	}
	c.hwId = references.IdIMos6526(c.IMos6526, c.label, 0)
	return c
}

func (w *CIA1Socket) HardwareId() string {
	return w.hwId
}

// Wire initializes the CIA1Socket with the provided CIA instance, connections, keyboard, and joystick references.
// It sets up the CIA via the Bind method and returns any errors encountered during initialization.
func (w *CIA1Socket) Wire() error {
	var err error
	w.component = w.parent.GetChildByHardwareId(w.HardwareId())
	if w.IMos6526, err = references.ComponentToIMos6526(w.component); err != nil {
		return err
	}
	idVIC := references.IdIMos6569(w.vicRef, w.label, 0)
	if w.vicRef, err = references.ComponentToIMos6569(w.parent.GetChildByHardwareId(idVIC)); err != nil {
		return err
	}
	idKeys := references.IdIC64Keyboard(w.keysRef, w.label, 0)
	if w.keysRef, err = references.ComponentToIC64Keyboard(w.parent.GetChildByHardwareId(idKeys)); err != nil {
		return err
	}
	idJoy1 := references.IdIC64Joystick(w.joy1Ref, w.label, 0)
	if w.joy1Ref, err = references.ComponentToIC64Joystick(w.parent.GetChildByHardwareId(idJoy1)); err != nil {
		return err
	}
	idJoy2 := references.IdIC64Joystick(w.joy1Ref, w.label, 1)
	if w.joy2Ref, err = references.ComponentToIC64Joystick(w.parent.GetChildByHardwareId(idJoy2)); err != nil {
		return err
	}
	if err = w.IMos6526.Bind(w); err != nil {
		return err
	}
	w.selfReadDDRB = w.ReadDDRB
	w.selfReadPRB = w.ReadPRB
	w.keysReset = w.keysRef.Reset
	w.joy1Reset = w.joy1Ref.Reset
	w.joy2Reset = w.joy2Ref.Reset
	w.joy1Poll = w.joy1Ref.Poll
	w.joy2Poll = w.joy2Ref.Poll
	w.keysPoll = w.keysRef.Poll
	w.vicLightPenTrigger = w.vicRef.LightPenTrigger
	return nil
}

// Reset resets the internal state of the CIA1Socket and its connected components to their default values.
func (w *CIA1Socket) Reset() {
	w.keysReset()
	w.joy1Reset()
	w.joy2Reset()
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

// ReadPortA reads data from port A by combining the given parameters with internal joystick states and matrices.
func (w *CIA1Socket) ReadPortA(prA uint8, prB uint8, ddrA uint8, ddrB uint8) uint8 {
	w.pollJoy2()
	w.pollKeyboard()
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
func (w *CIA1Socket) ReadPortB(prA uint8, prB uint8, ddrA uint8, ddrB uint8) uint8 {
	w.pollJoy1()
	w.pollKeyboard()
	ret := ^ddrB
	tst := (prA | ^ddrA) & w.joy2State
	for idx, bit := range bits.Uint8s {
		if (tst & bit) == 0 {
			ret &= w.keyMatrix[idx]
		}
	}
	return (ret | (prB & ddrB)) & w.joy1State
}

// SignalPRA handles data writing to port A based on the current state and provided parameters.
func (w *CIA1Socket) SignalPRA(_ uint8) {
}

// SignalDDRA updates the Code Direction Register A (DDRA) with the provided parameters without performing any operation.
func (w *CIA1Socket) SignalDDRA(_ uint8) {
}

// SignalPRB handles writing to port B by updating the light pen state based on the given port and data direction registers.
func (w *CIA1Socket) SignalPRB(prB uint8) {
	ddrB := w.selfReadDDRB()
	w.updateLightPen(prB, ddrB)
}

// SignalDDRB handles updates to the DDRB register and triggers updates to the light pen based on the PRB and DDRB values.
func (w *CIA1Socket) SignalDDRB(ddrB uint8) {
	prB := w.selfReadPRB()
	w.updateLightPen(prB, ddrB)
}

// ReadSP reads the state of the SP (Serial Port) line and returns its boolean value.
func (w *CIA1Socket) ReadSP() bool {
	return false
}

// SignalSP sets the level of the SP (Serial Port) line, typically for signaling or data output purposes.
func (w *CIA1Socket) SignalSP( /*level*/ _ bool) {
}

// IRQTrigger triggers an interrupt request (IRQ) using the associated interrupt Id managed by the connection interface.
func (w *CIA1Socket) IRQTrigger() {
	w.connectionIRQTrigger(w.intrId)
}

// IRQClearTrigger clears the interrupt request associated with the socket by invoking the IRQClearTrigger method on connections.
func (w *CIA1Socket) IRQClearTrigger() {
	w.connectionIRQClearTrigger(w.intrId)
}

// updateLightPen checks the state of the light pen line and triggers an action if the state has changed.
func (w *CIA1Socket) updateLightPen(prB uint8, ddrB uint8) {
	if ((prB | ^ddrB) & 0x10) != w.prevLPState {
		w.vicLightPenTrigger()
	}
	w.prevLPState = (prB | ^ddrB) & 0x10
}

// pollJoy1 polls the state of the first joystick and updates its internal state if a valid response is received.
func (w *CIA1Socket) pollJoy1() {
	if joy1State, ok := w.joy1Poll(); ok {
		w.joy1State = joy1State
	}
}

// pollJoy2 polls the state of the second joystick and updates its internal state if a valid response is received.
func (w *CIA1Socket) pollJoy2() {
	if joy2State, ok := w.joy2Poll(); ok {
		w.joy2State = joy2State
	}
}

// pollKeyboard updates the keyMatrix and revMatrix based on the state of the keyboard, handling key presses and releases.
func (w *CIA1Socket) pollKeyboard() {
	if v, ok := w.keysPoll(); ok {
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
}
