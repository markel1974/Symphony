package cia

import (
	"fmt"
	"github.com/markel1974/c64emu/src/signals"
)

type Timer struct {
	id                   string
	cr                   uint8
	hasNewCr             bool   // Flag: New value for CRB pending
	newCr                uint8  // New values for CRB
	timer                uint16 // Timer
	timerState           uint8  // Timer states
	latch                uint16 // Timer latch
	timerCntPhi2         bool   // Flag: Timer is counting Phi 2
	timerCntTimerX       bool   // Flag: Timer is counting underflow's of Timer X
	signalTimerUnderflow *signals.Signal
}

func NewTimer(id string) *Timer {
	m := &Timer{
		id:                   id,
		signalTimerUnderflow: signals.NewSignal(),
	}
	m.Reset()
	return m
}

func (m *Timer) Reset() {
	m.hasNewCr = false
	m.timer = 0xffff
	m.timerCntPhi2 = false
	m.timerCntTimerX = false
	m.timerState = timerStop
	m.latch = 1
}

func (m *Timer) GetTimerLow() uint8 {
	return uint8(m.timer)
}

func (m *Timer) GetTimerHigh() uint8 {
	return uint8(m.timer >> 8)
}

func (m *Timer) GetCR() uint8 {
	return m.cr
}

func (m *Timer) SetLatchLowByte(data uint8) {
	m.latch = (m.latch & 0xff00) | uint16(data)
}

func (m *Timer) SetLatchHighByte(data uint8) {
	m.latch = (m.latch & 0xff) | (uint16(data) << 8)
	if (m.cr & 1) == 0 {
		// Reload timer if stopped
		m.timer = m.latch
	}
}

func (m *Timer) TimerControl(data uint8, count bool) {
	m.hasNewCr = true
	m.newCr = data
	if count {
		m.timerCntPhi2 = (data & 0x60) == 0x00
		m.timerCntTimerX = (data & 0x60) == 0x40
	} else {
		m.timerCntPhi2 = (data & 0x20) == 0x00
	}
}

func (m *Timer) SignalTimerUnderflowBind(fn func()) {
	m.signalTimerUnderflow.Bind(fn)
}

func (m *Timer) Emulate(timerXUnderflow bool) bool {
	taUnderflow := false
	interrupt := false

	// Timer state machine
	switch m.timerState {
	case timerWaitThenCount:
		// fall through
		m.timerState = timerCount
		goto labelIdle
	case timerStop:
		goto labelIdle
	case timerLoadThenStop:
		m.timerState = timerStop
		// Reload timer
		m.timer = m.latch
		goto labelIdle
	case timerLoadThenCount:
		m.timerState = timerCount
		// Reload timer
		m.timer = m.latch
		goto labelIdle
	case timerLoadThenWaitThenCount:
		m.timerState = timerWaitThenCount
		if m.timer == 1 {
			// Interrupt if timer == 1
			interrupt = true
			goto labelCount
		} else {
			// Reload timer
			m.timer = m.latch
			goto labelIdle
		}
	case timerCount:
		goto labelCount
	case timerCountThenStop:
		m.timerState = timerStop
		goto labelCount
	}

	// Count timer
labelCount:
	if !interrupt {
		if m.timerCntPhi2 || (m.timerCntTimerX && timerXUnderflow) {
			// Decrement timer, underflow?
			timer := m.timer
			m.timer--
			if (timer == 0) || (m.timer == 0) {
				if m.timerState != timerStop {
					interrupt = true
				}
				taUnderflow = true
			}
		}
	}

	if interrupt {
		// Reload timer
		m.timer = m.latch
		//fmt.Println(m.id, "EMITTING CR", m.cr)
		m.signalTimerUnderflow.Emit()
		if (m.cr & 8) != 0 {
			// One-shot, stop timer
			m.cr &= 0xfe
			//fmt.Println(m.id, "ONE SHOT, STOPPING CR", m.cr)
			m.newCr &= 0xfe
			// Reload in next cycle
			m.timerState = timerLoadThenStop
		} else {
			// No One-shot, delay one cycle (and reload)
			m.timerState = timerLoadThenCount
		}
		//TODO VERIFY!!!!!
		taUnderflow = true
	}

	// Delayed write to CR?
labelIdle:
	if m.hasNewCr {
		switch m.timerState {
		case timerStop, timerLoadThenStop:
			// Timer started, wasn't running
			if (m.newCr & 1) != 0 {
				if (m.newCr & 0x10) != 0 {
					// Force load
					m.timerState = timerLoadThenWaitThenCount
				} else {
					// No force load
					m.timerState = timerWaitThenCount
				}
			} else {
				// Timer stopped, was already stopped
				if (m.newCr & 0x10) != 0 {
					// Force load
					m.timerState = timerLoadThenStop
				}
			}
		case timerCount:
			if (m.newCr & 1) != 0 {
				// Timer started, was already running
				if (m.newCr & 0x10) != 0 {
					// Force load
					m.timerState = timerLoadThenWaitThenCount
				}
			} else {
				// Timer stopped, was running
				if (m.newCr & 0x10) != 0 {
					// Force load
					m.timerState = timerLoadThenStop
				} else {
					// No force load
					m.timerState = timerCountThenStop
				}
			}
		case timerLoadThenCount, timerWaitThenCount:
			if (m.newCr & 1) != 0 {
				if (m.newCr & 8) != 0 {
					// One-shot, stop timer
					m.newCr &= 0xfe
					m.timerState = timerStop
				} else if (m.newCr & 0x10) != 0 {
					// No One-shot, force load
					m.timerState = timerLoadThenWaitThenCount
				}
			} else {
				m.timerState = timerStop
			}
		default:
			fmt.Println("TIMER - UNDEFINED", m.timerState)
		}
		m.cr = m.newCr & 0xef
		//fmt.Println(m.id, "CREATING NEW CR", m.cr)
		m.hasNewCr = false
	}
	return taUnderflow
}
