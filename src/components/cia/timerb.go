package cia

import (
	"fmt"
	"github.com/markel1974/c64emu/src/signals"
)

type TimerB struct {
	crB             uint8
	hasNewCrB       bool   // Flag: New value for CRB pending
	newCrB          uint8  // New values for CRB
	timerB          uint16 // Timer B
	timerBState     uint8  // Timer B states
	latchB          uint16 // Timer latch B
	timerBCntPhi2   bool   // Flag: Timer B is counting Phi 2
	timerBCntTimerA bool   // Flag: Timer B is counting underflow's of Timer A

	signalTimerBUnderflow *signals.Signal
}

func NewTimerB() *TimerB {
	m := &TimerB{
		signalTimerBUnderflow: signals.NewSignal(),
	}
	m.Reset()
	return m
}

func (m *TimerB) Reset() {
	m.hasNewCrB = false
	m.timerB = 0xffff
	m.timerBCntPhi2 = false
	m.timerBCntTimerA = false
	m.timerBState = timerStop
	m.latchB = 1
}

func (m *TimerB) SignalTimerBUnderflowBind(fn func()) {
	m.signalTimerBUnderflow.Bind(fn)
}

func (m *TimerB) Emulate(timerAUnderflow bool) {
	tbUseInterrupt := false

	// Timer B state machine
	switch m.timerBState {
	case timerWaitThenCount:
		// fall through
		m.timerBState = timerCount
		goto tbIdle

	case timerStop:
		goto tbIdle

	case timerLoadThenStop:
		m.timerBState = timerStop
		// Reload timer
		m.timerB = m.latchB
		goto tbIdle

	case timerLoadThenCount:
		m.timerBState = timerCount
		// Reload timer
		m.timerB = m.latchB
		goto tbIdle

	case timerLoadThenWaitThenCount:
		m.timerBState = timerWaitThenCount
		if m.timerB == 1 {
			// Interrupt if timer == 1
			tbUseInterrupt = true
			goto tbCount
		} else {
			// Reload timer
			m.timerB = m.latchB
			goto tbIdle
		}

	case timerCount:
		goto tbCount

	case timerCountThenStop:
		m.timerBState = timerStop
		goto tbCount
	}

	// Count timer B
tbCount:
	if !tbUseInterrupt {
		if m.timerBCntPhi2 || (m.timerBCntTimerA && timerAUnderflow) {
			// Decrement timer, underflow?
			tb := m.timerB
			m.timerB--
			if (tb == 0) || (m.timerB == 0) {
				if m.timerBState != timerStop {
					// tb_interrupt
					tbUseInterrupt = true
				}
			}
		}
	}

	if tbUseInterrupt {
		// Reload timer
		m.timerB = m.latchB
		m.signalTimerBUnderflow.Emit()
		// One-shot?
		if (m.crB & 8) != 0 {
			// Yes, stop timer
			m.crB &= 0xfe
			m.newCrB &= 0xfe
			// Reload in next cycle
			m.timerBState = timerLoadThenStop
		} else {
			// No, delay one cycle (and reload)
			m.timerBState = timerLoadThenCount
		}
	}

	// Delayed write to CRB?
tbIdle:
	if m.hasNewCrB {
		switch m.timerBState {
		case timerStop, timerLoadThenStop:
			// Timer started, wasn't running
			if (m.newCrB & 1) != 0 {
				if (m.newCrB & 0x10) != 0 {
					// Force load
					m.timerBState = timerLoadThenWaitThenCount
				} else {
					// No force load
					m.timerBState = timerWaitThenCount
				}
			} else {
				// Timer stopped, was already stopped
				if (m.newCrB & 0x10) != 0 {
					// Force load
					m.timerBState = timerLoadThenStop
				}
			}

		case timerCount:
			if (m.newCrB & 1) != 0 {
				// Timer started, was already running
				if (m.newCrB & 0x10) != 0 {
					// Force load
					m.timerBState = timerLoadThenWaitThenCount
				}
			} else {
				// Timer stopped, was running
				if (m.newCrB & 0x10) != 0 {
					// Force load
					m.timerBState = timerLoadThenStop
				} else {
					// No force load
					m.timerBState = timerCountThenStop
				}
			}

		case timerLoadThenCount, timerWaitThenCount:
			if (m.newCrB & 1) != 0 {
				// One-shot?
				if (m.newCrB & 8) != 0 {
					// Yes, stop timer
					m.newCrB &= 0xfe
					m.timerBState = timerStop
				} else if (m.newCrB & 0x10) != 0 {
					// Force load
					m.timerBState = timerLoadThenWaitThenCount
				}
			} else {
				m.timerBState = timerStop
			}
		default:
			fmt.Println("TIMER B - UNDEFINED", m.timerBState)
		}
		m.crB = m.newCrB & 0xef
		m.hasNewCrB = false
	}
}
