package cia

import (
	"github.com/markel1974/c64emu/src/flag"
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

const (
	IRQUnderflowTimerA = 0x1
	IRQUnderflowTimerB = 0x2
	IRQTODAlarmEqual   = 0x4
	IRQSDRFullOtEmpty  = 0x8
	IRQFlagPin         = 0x10
	IRQOccurred        = 0x80
)

// Timer states
const (
	T_STOP = iota
	T_WAIT_THEN_COUNT
	T_LOAD_THEN_STOP
	T_LOAD_THEN_COUNT
	T_LOAD_THEN_WAIT_THEN_COUNT
	T_COUNT
	T_COUNT_THEN_STOP
)

type Timers struct {
	prA                uint8
	ddrB               uint8
	crA                uint8
	newCrA             uint8  // New values for CRA
	hasNewCrA          bool   // Flag: New value for CRA pending
	timerA             uint16 // Timer A
	timerAState        uint8  // Timer A states
	latchA             uint16 // Timer latch A
	timerACntPhi2      bool   // Flag: Timer A is counting Phi 2
	timerAIrqNextCycle bool   // Flag: Trigger TA IRQ in next cycle

	ddrA               uint8
	prB                uint8
	crB                uint8
	hasNewCrB          bool   // Flag: New value for CRB pending
	newCrB             uint8  // New values for CRB
	timerB             uint16 // Timer B
	timerBState        uint8  // Timer B states
	latchB             uint16 // Timer latch B
	timerBCntPhi2      bool   // Flag: Timer B is counting Phi 2
	timerBCntTimerA    bool   // Flag: Timer B is counting underflow's of Timer A
	timerBIrqNextCycle bool   // Flag: Trigger Timer B IRQ in next cycle

	sdr uint8
	icr uint8 // Pending interrupts

	intMask         uint8 // Enabled interrupts
	interruptSignal *signals.SignalByte
}

func NewMOS6526() *Timers {
	m := &Timers{
		interruptSignal: signals.NewSignalByte(),
	}
	m.Reset()
	return m
}

func (m *Timers) Reset() {
	m.hasNewCrA = false
	m.hasNewCrB = false
	m.timerA = 0xffff
	m.timerACntPhi2 = false
	m.timerAIrqNextCycle = false
	m.timerAState = T_STOP
	m.timerB = 0xffff
	m.timerBCntPhi2 = false
	m.timerBCntTimerA = false
	m.timerBIrqNextCycle = false
	m.timerBState = T_STOP
	m.latchA = 1
	m.latchB = 1
}

func (m *Timers) SignalInterruptBind(fn func(uint8)) {
	m.interruptSignal.Bind(fn)
}

func (m *Timers) CheckIRQs() {
	// Trigger pending interrupts
	if m.timerAIrqNextCycle {
		m.timerAIrqNextCycle = false
		m.interruptSignal.Emit(IRQUnderflowTimerA)
	}
	if m.timerBIrqNextCycle {
		m.timerBIrqNextCycle = false
		m.interruptSignal.Emit(IRQUnderflowTimerB)
	}
}

func (m *Timers) Emulate() {
	taUnderflow := m.emulateTimerA()
	m.emulateTimerB(taUnderflow)
}

func (m *Timers) emulateTimerA() bool {
	taUnderflow := false
	taUseInterrupt := false

	// Timer A state machine
	switch m.timerAState {
	case T_WAIT_THEN_COUNT:
		// fall through
		m.timerAState = T_COUNT
		goto ta_idle

	case T_STOP:
		goto ta_idle

	case T_LOAD_THEN_STOP:
		m.timerAState = T_STOP
		// Reload timer
		m.timerA = m.latchA
		goto ta_idle

	case T_LOAD_THEN_COUNT:
		m.timerAState = T_COUNT
		// Reload timer
		m.timerA = m.latchA
		goto ta_idle

	case T_LOAD_THEN_WAIT_THEN_COUNT:
		m.timerAState = T_WAIT_THEN_COUNT
		if m.timerA == 1 {
			// Interrupt if timer == 1
			taUseInterrupt = true
			goto ta_count
			//goto ta_interrupt
		} else {
			m.timerA = m.latchA // Reload timer
			goto ta_idle
		}

	case T_COUNT:
		goto ta_count

	case T_COUNT_THEN_STOP:
		m.timerAState = T_STOP
		goto ta_count
	}

	// Count timer A
ta_count:
	if !taUseInterrupt {
		if m.timerACntPhi2 {
			ta := m.timerA
			m.timerA--
			if (ta == 0) || (m.timerA == 0) {
				// Decrement timer, underflow?
				if m.timerAState != T_STOP {
					//ta_interrupt:
					taUseInterrupt = true
				}
				taUnderflow = true
			}
		}
	}

	if taUseInterrupt {
		// Reload timer
		m.timerA = m.latchA
		// Trigger interrupt in next cycle
		m.timerAIrqNextCycle = true
		// But set ICR bit now
		m.icr |= 1

		// One-shot?
		if flag.Uint8ToBool(m.crA & 8) {
			// Yes, stop timer
			m.crA &= 0xfe
			m.newCrA &= 0xfe
			// Reload in next cycle
			m.timerAState = T_LOAD_THEN_STOP
		} else {
			// No, delay one cycle (and reload)
			m.timerAState = T_LOAD_THEN_COUNT
		}
		//TODO VERIFY!!!!!
		taUnderflow = true
	}

	// Delayed write to CRA?
ta_idle:
	if m.hasNewCrA {
		switch m.timerAState {
		case T_STOP, T_LOAD_THEN_STOP:
			if flag.Uint8ToBool(m.newCrA & 1) {
				// Timer started, wasn't running
				if flag.Uint8ToBool(m.newCrA & 0x10) {
					// Force load
					m.timerAState = T_LOAD_THEN_WAIT_THEN_COUNT
				} else {
					// No force load
					m.timerAState = T_WAIT_THEN_COUNT
				}
			} else {
				// Timer stopped, was already stopped
				if flag.Uint8ToBool(m.newCrA & 0x10) {
					// Force load
					m.timerAState = T_LOAD_THEN_STOP
				}
			}

		case T_COUNT:
			if flag.Uint8ToBool(m.newCrA & 1) {
				// Timer started, was already running
				if flag.Uint8ToBool(m.newCrA & 0x10) {
					// Force load
					m.timerAState = T_LOAD_THEN_WAIT_THEN_COUNT
				}
			} else {
				// Timer stopped, was running
				if flag.Uint8ToBool(m.newCrA & 0x10) {
					// Force load
					m.timerAState = T_LOAD_THEN_STOP
				} else {
					// No force load
					m.timerAState = T_COUNT_THEN_STOP
				}
			}

		case T_LOAD_THEN_COUNT, T_WAIT_THEN_COUNT:
			if flag.Uint8ToBool(m.newCrA & 1) {
				// One-shot?
				if flag.Uint8ToBool(m.newCrA & 8) {
					// Yes, stop timer
					m.newCrA &= 0xfe
					m.timerAState = T_STOP
				} else if flag.Uint8ToBool(m.newCrA & 0x10) {
					// Force load
					m.timerAState = T_LOAD_THEN_WAIT_THEN_COUNT
				}
			} else {
				m.timerAState = T_STOP
			}
		}
		m.crA = m.newCrA & 0xef
		m.hasNewCrA = false
	}
	return taUnderflow
}

func (m *Timers) emulateTimerB(timerAUnderflow bool) {
	tbUseInterrupt := false

	// Timer B state machine
	switch m.timerBState {
	case T_WAIT_THEN_COUNT:
		// fall through
		m.timerBState = T_COUNT
		goto tb_idle

	case T_STOP:
		goto tb_idle

	case T_LOAD_THEN_STOP:
		m.timerBState = T_STOP
		// Reload timer
		m.timerB = m.latchB
		goto tb_idle

	case T_LOAD_THEN_COUNT:
		m.timerBState = T_COUNT
		// Reload timer
		m.timerB = m.latchB
		goto tb_idle

	case T_LOAD_THEN_WAIT_THEN_COUNT:
		m.timerBState = T_WAIT_THEN_COUNT
		if m.timerB == 1 {
			// Interrupt if timer == 1
			tbUseInterrupt = true
			goto tb_count
		} else {
			// Reload timer
			m.timerB = m.latchB
			goto tb_idle
		}

	case T_COUNT:
		goto tb_count

	case T_COUNT_THEN_STOP:
		m.timerBState = T_STOP
		goto tb_count
	}

	// Count timer B
tb_count:
	if !tbUseInterrupt {
		if m.timerBCntPhi2 || (m.timerBCntTimerA && timerAUnderflow) {
			// Decrement timer, underflow?
			tb := m.timerB
			m.timerB--
			if (tb == 0) || (m.timerB == 0) {
				if m.timerBState != T_STOP {
					// tb_interrupt
					tbUseInterrupt = true
				}
			}
		}
	}

	if tbUseInterrupt {
		// Reload timer
		m.timerB = m.latchB
		// Trigger interrupt in next cycle
		m.timerBIrqNextCycle = true
		// But set ICR bit now
		m.icr |= 2
		// One-shot?
		if flag.Uint8ToBool(m.crB & 8) {
			// Yes, stop timer
			m.crB &= 0xfe
			m.newCrB &= 0xfe
			// Reload in next cycle
			m.timerBState = T_LOAD_THEN_STOP
		} else {
			// No, delay one cycle (and reload)
			m.timerBState = T_LOAD_THEN_COUNT
		}
	}

	// Delayed write to CRB?
tb_idle:
	if m.hasNewCrB {
		switch m.timerBState {
		case T_STOP, T_LOAD_THEN_STOP:
			// Timer started, wasn't running
			if flag.Uint8ToBool(m.newCrB & 1) {
				if flag.Uint8ToBool(m.newCrB & 0x10) {
					// Force load
					m.timerBState = T_LOAD_THEN_WAIT_THEN_COUNT
				} else {
					// No force load
					m.timerBState = T_WAIT_THEN_COUNT
				}
			} else {
				// Timer stopped, was already stopped
				if flag.Uint8ToBool(m.newCrB & 0x10) {
					// Force load
					m.timerBState = T_LOAD_THEN_STOP
				}
			}

		case T_COUNT:
			if flag.Uint8ToBool(m.newCrB & 1) {
				// Timer started, was already running
				if flag.Uint8ToBool(m.newCrB & 0x10) {
					// Force load
					m.timerBState = T_LOAD_THEN_WAIT_THEN_COUNT
				}
			} else {
				// Timer stopped, was running
				if flag.Uint8ToBool(m.newCrB & 0x10) {
					// Force load
					m.timerBState = T_LOAD_THEN_STOP
				} else {
					// No force load
					m.timerBState = T_COUNT_THEN_STOP
				}
			}

		case T_LOAD_THEN_COUNT, T_WAIT_THEN_COUNT:
			if flag.Uint8ToBool(m.newCrB & 1) {
				// One-shot?
				if flag.Uint8ToBool(m.newCrB & 8) {
					// Yes, stop timer
					m.newCrB &= 0xfe
					m.timerBState = T_STOP
				} else if flag.Uint8ToBool(m.newCrB & 0x10) {
					// Force load
					m.timerBState = T_LOAD_THEN_WAIT_THEN_COUNT
				}
			} else {
				m.timerBState = T_STOP
			}
		}
		m.crB = m.newCrB & 0xef
		m.hasNewCrB = false
	}
}
