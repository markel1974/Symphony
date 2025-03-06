package board

import (
	"github.com/markel1974/c64emu/src/common/bits"
	mos6526 "github.com/markel1974/c64emu/src/components/cia"
)

// defaultLPState represents the default state for LP (Launchpad) set to 0x10.
// defaultJoyState represents the default state for a joystick set to 0xff.
// defaultKeyState represents the default state for a keyboard set to 0xff.
const (
	defaultLPState  = 0x10
	defaultJoyState = 0xff
	defaultKeyState = 0xff
)

// CIA1Socket represents the state and behavior of the CIA1 socket in a computing system.
// board refers to the system board connected to the CIA1 socket.
// intrId represents the ID of the interrupt associated with CIA1.
// prevLPState tracks the previous state of the light pen signal (bit 4).
// keyMatrix holds the current state of the keyboard matrix (0: down, 1: up).
// revMatrix represents the reversed keyboard matrix for compatibility or emulation purposes.
// joy1 denotes the state of joystick 1 connected to the system.
// joy2 denotes the state of joystick 2 connected to the system.
type CIA1Socket struct {
	cia1        *mos6526.CIA
	board       *Board  //
	intrId      uint32  //
	prevLPState uint8   // Previous state of LP line (bit 4)
	keyMatrix   []uint8 // keyboard matrix [0: down, 1: up]
	revMatrix   []uint8 // Reversed keyboard matrix
	joy1        uint8   // Joystick 1
	joy2        uint8   // Joystick 2
}

// NewCIA1Socket creates and initializes a new instance of CIA1Socket with default values for its fields.
func NewCIA1Socket() *CIA1Socket {
	c := &CIA1Socket{
		cia1:        nil,
		board:       nil,
		intrId:      intrIrqCia1Bit,
		prevLPState: defaultLPState,
		keyMatrix:   make([]uint8, 8),
		revMatrix:   make([]uint8, 8),
		joy1:        defaultJoyState,
		joy2:        defaultJoyState,
	}
	return c
}

// Setup initializes the CIA1Socket with the provided Board reference and interrupt ID.
func (w *CIA1Socket) Setup(board *Board, cia1 *mos6526.CIA) {
	w.board = board
	w.cia1 = cia1
	w.cia1.Setup(w)
}

// Reset reinitializes the CIA1Socket by resetting its board components, key matrices, joystick states, and light pen state.
func (w *CIA1Socket) Reset() {
	w.board.keys.Reset()
	w.board.joy1.Reset()
	w.board.joy2.Reset()
	w.cia1.Reset()
	for idx := range w.keyMatrix {
		w.keyMatrix[idx] = defaultKeyState
	}
	for idx := range w.revMatrix {
		w.revMatrix[idx] = defaultKeyState
	}
	w.joy1 = defaultJoyState
	w.joy2 = defaultJoyState
	w.prevLPState = defaultLPState
}

func (w *CIA1Socket) Emulate() {
	w.cia1.Emulate()
}

// Update synchronizes the joystick and keyboard states by polling their current inputs and updates the key matrices accordingly.
func (w *CIA1Socket) Update() {
	if joy1State, ok := w.board.joy1.Poll(); ok {
		w.joy1 = joy1State
	}
	if joy2State, ok := w.board.joy2.Poll(); ok {
		w.joy2 = joy2State
	}
	if v, ok := w.board.keys.Poll(); ok {
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
	w.cia1.Update()
}

// ReadPortA reads data from Port A taking into account the provided port registers and joystick states.
// The method integrates Port B and reads the reverse keyboard matrix when joystick inputs are active.
// The result is masked with joystick 2's state before being returned.
func (w *CIA1Socket) ReadPortA(prA uint8, ddrA uint8, prB uint8, ddrB uint8) uint8 {
	ret := prA | ^ddrA
	tst := (prB | ^ddrB) & w.joy1
	for idx, bit := range bits.Uint8s {
		if (tst & bit) == 0 {
			ret &= w.revMatrix[idx]
		}
	}
	return ret & w.joy2
}

// ReadPortB reads from Port B based on the provided port registers and data direction registers.
// It considers joystick inputs and keyboard matrix states to compute the resulting value.
// The method applies bitwise operations to determine the state of the port and filters results accordingly.
func (w *CIA1Socket) ReadPortB(prA uint8, ddrA uint8, prB uint8, ddrB uint8) uint8 {
	ret := ^ddrB
	tst := (prA | ^ddrA) & w.joy2
	for idx, bit := range bits.Uint8s {
		if (tst & bit) == 0 {
			ret &= w.keyMatrix[idx]
		}
	}
	return (ret | (prB & ddrB)) & w.joy1
}

// WritePortA writes data to Port A while considering the specified parameters but does not implement any functionality yet.
func (w *CIA1Socket) WritePortA(_ uint8, _ uint8, _ uint8, _ uint8) {
}

// WriteDdrA handles writing to Data Direction Register A of CIA1 and updates relevant internal states accordingly.
func (w *CIA1Socket) WriteDdrA(_ uint8, _ uint8, _ uint8, _ uint8) {
}

// WritePortB updates the light pen state based on the provided port B latch and data direction register values.
func (w *CIA1Socket) WritePortB(_ uint8, _ uint8, prB uint8, ddrB uint8) {
	w.updateLightPen(prB, ddrB)
}

// WriteDdrB updates the light pen state based on the provided DDRB and PRB values.
func (w *CIA1Socket) WriteDdrB(_ uint8, _ uint8, prB uint8, ddrB uint8) {
	w.updateLightPen(prB, ddrB)
}

// updateLightPen updates the light pen state based on Port B and DDR B values and triggers the light pen event if state changes.
func (w *CIA1Socket) updateLightPen(prB uint8, ddrB uint8) {
	if ((prB | ^ddrB) & 0x10) != w.prevLPState {
		w.board.vicSocket.LightPenTrigger()
	}
	w.prevLPState = (prB | ^ddrB) & 0x10
}

// IRQTrigger triggers an interrupt request on the board's IRQ slot using the associated interrupt ID.
func (w *CIA1Socket) IRQTrigger() {
	w.board.irqTriggerSlot(w.intrId)
}

// IRQClear clears the interrupt request for the associated hardware by invoking the board's irqClearSlot method.
func (w *CIA1Socket) IRQClear() {
	w.board.irqClearSlot(w.intrId)
}
