package board

import (
	"github.com/markel1974/c64emu/src/bits"
)

const (
	defaultLPState  = 0x10
	defaultJoyState = 0xff
	defaultKeyState = 0xff
)

type CIA1Socket struct {
	board       *Board  //
	intrId      uint32  //
	prevLPState uint8   // Previous state of LP line (bit 4)
	keyMatrix   []uint8 // keyboard matrix [0: down, 1: up]
	revMatrix   []uint8 // Reversed keyboard matrix
	joy1        uint8   // Joystick 1
	joy2        uint8   // Joystick 2
}

func NewCIA1Socket() *CIA1Socket {
	c := &CIA1Socket{
		board:       nil,
		intrId:      0,
		prevLPState: defaultLPState,
		keyMatrix:   make([]uint8, 8),
		revMatrix:   make([]uint8, 8),
		joy1:        defaultJoyState,
		joy2:        defaultJoyState,
	}
	return c
}

func (w *CIA1Socket) Setup(board *Board, intrId uint32) {
	w.board = board
	w.intrId = intrId
}

func (w *CIA1Socket) Reset() {
	w.board.keys.Reset()
	w.board.joy1.Reset()
	w.board.joy2.Reset()
	w.board.cia1.Reset()
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
	w.board.cia1.Update()
}

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

func (w *CIA1Socket) WritePortA(_ uint8, _ uint8, _ uint8, _ uint8) {
}

func (w *CIA1Socket) WriteDdrA(_ uint8, _ uint8, _ uint8, _ uint8) {
}

func (w *CIA1Socket) WritePortB(_ uint8, _ uint8, prB uint8, ddrB uint8) {
	w.updateLightPen(prB, ddrB)
}

func (w *CIA1Socket) WriteDdrB(_ uint8, _ uint8, prB uint8, ddrB uint8) {
	w.updateLightPen(prB, ddrB)
}

func (w *CIA1Socket) updateLightPen(prB uint8, ddrB uint8) {
	if ((prB | ^ddrB) & 0x10) != w.prevLPState {
		w.board.vic.LightPenTrigger()
	}
	w.prevLPState = (prB | ^ddrB) & 0x10
}

func (w *CIA1Socket) IRQTrigger() {
	w.board.irqTriggerSlot(w.intrId)
}

func (w *CIA1Socket) IRQClear() {
	w.board.irqClearSlot(w.intrId)
}
