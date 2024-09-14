package board

type CIA1Socket struct {
	board       *Board
	intrId      uint32
	prevLPState uint8    // Previous state of LP line (bit 4
	keyMatrix   [8]uint8 // C64 keyboard matrix, 1 bit/key (0: key down, 1: key up)
	revMatrix   [8]uint8 // Reversed keyboard matrix
	joy1        uint8    // Joystick 1 AND value
	joy2        uint8    // Joystick 2 AND value
}

func NewCIA1Socket() *CIA1Socket {
	c := &CIA1Socket{}
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

	for i := 0; i < 8; i++ {
		w.keyMatrix[i] = 0xff
		w.revMatrix[i] = 0xff
	}
	w.joy1 = 0xff
	w.joy2 = 0xff
	w.prevLPState = 0x10
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

func (w *CIA1Socket) ReadPortB(prA uint8, ddrA uint8, prB uint8, ddrB uint8) uint8 {
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
		w.board.vic.LightPenTrigger() // signalLightPenTrigger.Emit()
	}
	w.prevLPState = (prB | ^ddrB) & 0x10
}

func (w *CIA1Socket) IRQTrigger() {
	w.board.irqTriggerSlot(w.intrId)
}

func (w *CIA1Socket) IRQClear() {
	w.board.irqClearSlot(w.intrId)
}
