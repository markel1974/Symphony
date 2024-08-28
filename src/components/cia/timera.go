package cia

import (
	"fmt"
	"github.com/markel1974/c64emu/src/signals"
)

/*
 * Notes:
 * ------
 *
 *  - The Emulate() function is called for every emulated Phi2 clock cycle.
 *    It counts down the timers and triggers interrupts if necessary.
 *  - The TOD clocks are counted by TODUpdate() during the VBlank, so the input frequency is 50Hz
 *  - The fields keyMatrix and revMatrix contain one bit for each
 *    key on the C64 keyboard (0: key pressed, 1: key released).
 *    keyMatrix is used for normal keyboard polling (PRA->PRB),
 *    revMatrix for reversed polling (PRB->PRA).
 *
 * Incompatibilities:
 * ------------------
 *
 *  - The TOD clock should not be stopped on a read access, but be latched
 *  - The SDR interrupt is faked
 *  - Some small incompatibilities with the timers
 */

// Timer states
const (
	timerStop = iota
	timerWaitThenCount
	timerLoadThenStop
	timerLoadThenCount
	timerLoadThenWaitThenCount
	timerCount
	timerCountThenStop
)

type TimerA struct {
	crA                   uint8
	newCrA                uint8  // New values for CRA
	hasNewCrA             bool   // Flag: New value for CRA pending
	timerA                uint16 // Timer A
	timerAState           uint8  // Timer A states
	latchA                uint16 // Timer latch A
	timerACntPhi2         bool   // Flag: Timer A is counting Phi 2
	signalTimerAUnderflow *signals.Signal
}

func NewTimerA() *TimerA {
	m := &TimerA{
		signalTimerAUnderflow: signals.NewSignal(),
	}
	m.Reset()
	return m
}

func (m *TimerA) Reset() {
	m.hasNewCrA = false
	m.timerA = 0xffff
	m.timerACntPhi2 = false
	m.timerAState = timerStop
	m.latchA = 1
}

func (m *TimerA) SignalTimerAUnderflowBind(fn func()) {
	m.signalTimerAUnderflow.Bind(fn)
}

func (m *TimerA) Emulate() bool {
	taUnderflow := false
	taUseInterrupt := false

	// Timer A state machine
	switch m.timerAState {
	case timerWaitThenCount:
		// fall through
		m.timerAState = timerCount
		goto taIdle

	case timerStop:
		goto taIdle

	case timerLoadThenStop:
		m.timerAState = timerStop
		// Reload timer
		m.timerA = m.latchA
		goto taIdle

	case timerLoadThenCount:
		m.timerAState = timerCount
		// Reload timer
		m.timerA = m.latchA
		goto taIdle

	case timerLoadThenWaitThenCount:
		m.timerAState = timerWaitThenCount
		if m.timerA == 1 {
			// Interrupt if timer == 1
			taUseInterrupt = true
			goto taCount
			//goto ta_interrupt
		} else {
			m.timerA = m.latchA // Reload timer
			goto taIdle
		}

	case timerCount:
		goto taCount

	case timerCountThenStop:
		m.timerAState = timerStop
		goto taCount
	}

	// Count timer A
taCount:
	if !taUseInterrupt {
		if m.timerACntPhi2 {
			ta := m.timerA
			m.timerA--
			if (ta == 0) || (m.timerA == 0) {
				// Decrement timer, underflow?
				if m.timerAState != timerStop {
					taUseInterrupt = true
				}
				taUnderflow = true
			}
		}
	}

	if taUseInterrupt {
		// Reload timer
		m.timerA = m.latchA
		m.signalTimerAUnderflow.Emit()
		// One-shot?
		if (m.crA & 8) != 0 {
			// Yes, stop timer
			m.crA &= 0xfe
			m.newCrA &= 0xfe
			// Reload in next cycle
			m.timerAState = timerLoadThenStop
		} else {
			// No, delay one cycle (and reload)
			m.timerAState = timerLoadThenCount
		}
		//TODO VERIFY!!!!!
		taUnderflow = true
	}

	// Delayed write to CRA?
taIdle:
	if m.hasNewCrA {
		switch m.timerAState {
		case timerStop, timerLoadThenStop:
			if (m.newCrA & 1) != 0 {
				// Timer started, wasn't running
				if (m.newCrA & 0x10) != 0 {
					// Force load
					m.timerAState = timerLoadThenWaitThenCount
				} else {
					// No force load
					m.timerAState = timerWaitThenCount
				}
			} else {
				// Timer stopped, was already stopped
				if (m.newCrA & 0x10) != 0 {
					// Force load
					m.timerAState = timerLoadThenStop
				}
			}

		case timerCount:
			if (m.newCrA & 1) != 0 {
				// Timer started, was already running
				if (m.newCrA & 0x10) != 0 {
					// Force load
					m.timerAState = timerLoadThenWaitThenCount
				}
			} else {
				// Timer stopped, was running
				if (m.newCrA & 0x10) != 0 {
					// Force load
					m.timerAState = timerLoadThenStop
				} else {
					// No force load
					m.timerAState = timerCountThenStop
				}
			}

		case timerLoadThenCount, timerWaitThenCount:
			if (m.newCrA & 1) != 0 {
				// One-shot?
				if (m.newCrA & 8) != 0 {
					// Yes, stop timer
					m.newCrA &= 0xfe
					m.timerAState = timerStop
				} else if (m.newCrA & 0x10) != 0 {
					// Force load
					m.timerAState = timerLoadThenWaitThenCount
				}
			} else {
				m.timerAState = timerStop
			}
		default:
			fmt.Println("TIMER A - UNDEFINED", m.timerAState)
		}
		m.crA = m.newCrA & 0xef
		m.hasNewCrA = false
	}
	return taUnderflow
}
