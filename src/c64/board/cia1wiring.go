package board

import (
	"github.com/markel1974/c64emu/src/c64/keyboard"
	"github.com/markel1974/c64emu/src/signals"
)

type CIA1Wiring struct {
	signalLightPenTrigger *signals.Signal
	prevLPState           uint8    // Previous state of LP line (bit 4
	keyMatrix             [8]uint8 // C64 keyboard matrix, 1 bit/key (0: key down, 1: key up)
	revMatrix             [8]uint8 // Reversed keyboard matrix
	joy1                  uint8    // Joystick 1 AND value
	joy2                  uint8    // Joystick 2 AND value
	keys                  *keyboard.Keyboard
}

func NewCIA1AWiring() *CIA1Wiring {
	c := &CIA1Wiring{
		signalLightPenTrigger: signals.NewSignal(),
	}
	return c
}

func (w *CIA1Wiring) Setup(keys *keyboard.Keyboard, fn func()) {
	w.keys = keys
	w.signalLightPenTrigger.Bind(fn)
}

func (w *CIA1Wiring) Reset() {
	for i := 0; i < 8; i++ {
		w.keyMatrix[i] = 0xff
		w.revMatrix[i] = 0xff
	}
	w.joy1 = 0xff
	w.joy2 = 0xff
	w.prevLPState = 0x10
}

func (w *CIA1Wiring) ReadPortA(prA uint8, ddrA uint8, prB uint8, ddrB uint8) uint8 {
	w.updateInputs()
	//Joy port 2
	ret := prA | ^ddrA
	tst := (prB | ^ddrB) & w.joy1
	if (tst & 0x01) == 0 {
		ret &= w.revMatrix[0]
	}
	if (tst & 0x02) == 0 {
		ret &= w.revMatrix[1]
	}
	if (tst & 0x04) == 0 {
		ret &= w.revMatrix[2]
	}
	if (tst & 0x08) == 0 {
		ret &= w.revMatrix[3]
	}
	if (tst & 0x10) == 0 {
		ret &= w.revMatrix[4]
	}
	if (tst & 0x20) == 0 {
		ret &= w.revMatrix[5]
	}
	if (tst & 0x40) == 0 {
		ret &= w.revMatrix[6]
	}
	if (tst & 0x80) == 0 {
		ret &= w.revMatrix[7]
	}
	return ret & w.joy2
}

func (w *CIA1Wiring) ReadPortB(prA uint8, ddrA uint8, prB uint8, ddrB uint8) uint8 {
	w.updateInputs()
	//joy port 1
	ret := ^ddrB
	tst := (prA | ^ddrA) & w.joy2
	if (tst & 0x01) == 0 {
		ret &= w.keyMatrix[0]
	}
	if (tst & 0x02) == 0 {
		ret &= w.keyMatrix[1]
	}
	if (tst & 0x04) == 0 {
		ret &= w.keyMatrix[2]
	}
	if (tst & 0x08) == 0 {
		ret &= w.keyMatrix[3]
	}
	if (tst & 0x10) == 0 {
		ret &= w.keyMatrix[4]
	}
	if (tst & 0x20) == 0 {
		ret &= w.keyMatrix[5]
	}
	if (tst & 0x40) == 0 {
		ret &= w.keyMatrix[6]
	}
	if (tst & 0x80) == 0 {
		ret &= w.keyMatrix[7]
	}
	return (ret | (prB & ddrB)) & w.joy1
}

func (w *CIA1Wiring) WritePortA(_ uint8, _ uint8, _ uint8, _ uint8) {
}

func (w *CIA1Wiring) WriteDdrA(_ uint8, _ uint8, _ uint8, _ uint8) {
}

func (w *CIA1Wiring) WritePortB(_ uint8, _ uint8, prB uint8, ddrB uint8) {
	w.updateLightPen(prB, ddrB)
}

func (w *CIA1Wiring) WriteDdrB(_ uint8, _ uint8, prB uint8, ddrB uint8) {
	w.updateLightPen(prB, ddrB)
}

func (w *CIA1Wiring) updateLightPen(prB uint8, ddrB uint8) {
	if ((prB | ^ddrB) & 0x10) != w.prevLPState {
		w.signalLightPenTrigger.Emit()
	}
	w.prevLPState = (prB | ^ddrB) & 0x10
}

func (w *CIA1Wiring) updateInputs() {
	if joy1, ok := w.keys.PollJoystick1(); ok {
		w.joy1 = joy1
	}
	if joy2, ok := w.keys.PollJoystick2(); ok {
		w.joy2 = joy2
	}
	if keyM, revM, pressed, shifted, ok := w.keys.PollKeyboard(); ok {
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
