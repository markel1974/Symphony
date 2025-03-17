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

// CIA1Socket represents a connector for CIA-1 chip emulation, managing keyboard, joystick, and external connections.
// ICIA interface provides the core logic for the connected CIA chip functionality.
// IKeyboard interface is used to handle input and state changes from a keyboard.
// IJoystick interface specifies methods to interact with joystick devices (one or more).
// ICIA1SocketConnections interface provides external event interactions like IRQ and light pen triggers.
// intrId is used to track the IRQ identification bit for events.
// prevLPState represents the prior state of the light pen (LP) input signal.
// keyMatrix is a direct lookup table representing the state of the keyboard matrix.
// revMatrix is the reversed lookup table to assist with keyboard matrix computation.
// joy1State and joy2State track the current states of two connected joysticks.
type CIA1Socket struct {
	references.ICIA
	pic         references.IPIC6510
	vic         references.IVIC
	keys        references.IKeyboard
	joy1        references.IJoystick
	joy2        references.IJoystick
	intrId      uint32  //
	prevLPState uint8   // Previous state of LP line (bit 4)
	keyMatrix   []uint8 // keyboard matrix [0: down, 1: up]
	revMatrix   []uint8 // Reversed keyboard matrix
	joy1State   uint8   // Joystick 1
	joy2State   uint8   // Joystick 2
}

// NewCIA1Socket creates and initializes a new instance of CIA1Socket with default state and properties.
func NewCIA1Socket() *CIA1Socket {
	c := &CIA1Socket{
		ICIA:        nil,
		pic:         nil,
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
	return c
}

// Connect initializes the CIA1Socket with the provided CIA instance, connections, keyboard, and joystick references.
// It sets up the CIA via the Setup method and returns any errors encountered during initialization.
func (w *CIA1Socket) Connect(cia1 references.ICIA, pic references.IPIC6510, vic references.IVIC, keys references.IKeyboard, joy1 references.IJoystick, joy2 references.IJoystick) error {
	w.ICIA = cia1
	w.pic = pic
	w.vic = vic
	w.keys = keys
	w.joy1 = joy1
	w.joy2 = joy2
	if err := w.ICIA.Setup(w); err != nil {
		return err
	}
	return nil
}

// Reset resets the internal state of the CIA1Socket and its connected components to their default values.
func (w *CIA1Socket) Reset() {
	w.keys.Reset()
	w.joy1.Reset()
	w.joy2.Reset()
	w.ICIA.Reset()
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
	w.ICIA.Update()
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

// updateLightPen checks the state of the light pen line and triggers an action if the state has changed.
func (w *CIA1Socket) updateLightPen(prB uint8, ddrB uint8) {
	if ((prB | ^ddrB) & 0x10) != w.prevLPState {
		w.vic.LightPenTrigger()
	}
	w.prevLPState = (prB | ^ddrB) & 0x10
}

// IRQTrigger triggers an interrupt request (IRQ) using the associated interrupt ID managed by the connection interface.
func (w *CIA1Socket) IRQTrigger() {
	w.pic.TriggerIRQ(w.intrId)
}

// IRQClear clears the interrupt request associated with the socket by invoking the IRQClear method on connections.
func (w *CIA1Socket) IRQClear() {
	w.pic.ClearIRQ(w.intrId)
}
