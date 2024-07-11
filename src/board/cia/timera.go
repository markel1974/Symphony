package cia

type TimerA struct {
	newCr        uint8 // New values for CRA
	hasNewCr     bool  // Flag: New value for CRA pending
	pendingCrIrq *uint8 // Pending interrupts //TODO NON DEVE ESSERE CONDIVISO!!!!
	cr           uint8
	timer        uint16 // Timer A
	state        uint8  // Timer A states
	latch        uint16 // Timer latch A
	cntPhi2      bool   // Flag: Timer A is counting Phi 2
	irqNextCycle bool   // Flag: Trigger TA IRQ in next cycle
	underflow    bool
	useInterrupt bool
	statesFn     [7]func()
}

func NewTimerA() *TimerA {
	t := &TimerA{
		//pendingCrIrq: 0,
		timer:        0,
		state:        0,
		latch:        0,
		cntPhi2:      false,
		irqNextCycle: false,
		underflow:    false,
		useInterrupt: false,
		newCr:        0,
	}
	t.statesFn[T_STOP] = t.stateStop
	t.statesFn[T_WAIT_THEN_COUNT] = t.stateWaitThenCount
	t.statesFn[T_LOAD_THEN_STOP] = t.stateLoadThenStop
	t.statesFn[T_LOAD_THEN_COUNT] = t.stateLoadThenCount
	t.statesFn[T_LOAD_THEN_WAIT_THEN_COUNT] = t.stateLoadThenWaitThenCount
	t.statesFn[T_COUNT] = t.stateCount
	t.statesFn[T_COUNT_THEN_STOP] = t.stateCountThenStop
	t.Reset()
	return t
}

func (m *TimerA) Reset() {
	//m.pendingCrIrq = 0
	m.newCr = 0
	m.timer = 0xffff
	m.cntPhi2 = false
	m.irqNextCycle = false
	m.state = T_STOP
	m.latch = 1
	m.underflow = false
	m.useInterrupt = false
}

func (m * TimerA) CheckIrq() bool {
	return m.irqNextCycle
}

func (m * TimerA) SetCr(v uint8) {
	m.cr = v
}

func (m * TimerA) GetCr() uint8 {
	return m.cr
}

func (m * TimerA) SetCountPhi2(v bool) {
	m.cntPhi2 = v
}

func (m *TimerA) EmulateTimer() bool {
	m.underflow = false
	m.useInterrupt = false
	m.statesFn[m.state]()
	return m.underflow
}

func (m *TimerA) stateStop() {
	m.idle()
}

func (m *TimerA) stateWaitThenCount() {
	// fall through
	m.state = T_COUNT
	m.idle()
}

func (m *TimerA) stateLoadThenStop() {
	m.state = T_STOP
	// Reload timer
	m.timer = m.latch
	m.idle()
}

func (m *TimerA) stateLoadThenCount() {
	m.state = T_COUNT
	// Reload timer
	m.timer = m.latch
	m.idle()
}

func (m *TimerA) stateLoadThenWaitThenCount() {
	m.state = T_WAIT_THEN_COUNT
	if m.timer == 1 {
		m.useInterrupt = true
		m.count()
		m.idle()
		return
	}
	// Reload timer
	m.timer = m.latch
	m.idle()
}

func (m *TimerA) stateCount() {
	m.count()
	m.idle()
}

func (m *TimerA) stateCountThenStop() {
	m.state = T_STOP
	m.count()
	m.idle()
}

func (m *TimerA) count() {
	if !m.useInterrupt {
		if m.cntPhi2 {
			ta := m.timer
			m.timer--
			if (ta == 0) || (m.timer == 0) {
				// Decrement timer, underflow?
				if m.state != T_STOP {
					m.useInterrupt = true
				}
				m.underflow = true
			}
		}
	}

	if m.useInterrupt {
		// Reload timer
		m.timer = m.latch
		// Trigger interrupt in next cycle
		m.irqNextCycle = true
		// But set ICR bit now
		*m.pendingCrIrq |= 1

		// One-shot?
		if (m.cr & 8) != 0 {
			// Yes, stop timer
			m.cr &= 0xfe
			m.newCr &= 0xfe
			// Reload in next cycle
			m.state = T_LOAD_THEN_STOP
		} else {
			// No, delay one cycle (and reload)
			m.state = T_LOAD_THEN_COUNT
		}
		//TODO VERIFY!!!!!
		m.underflow = true
	}
}

func (m *TimerA) idle() {
	if !m.hasNewCr {
		return
	}
	if m.state == T_STOP || m.state == T_LOAD_THEN_STOP {
		if (m.newCr & 1) != 0 {
			// Timer started, wasn't running
			if (m.newCr & 0x10) != 0 {
				// Force load
				m.state = T_LOAD_THEN_WAIT_THEN_COUNT
			} else {
				// No force load
				m.state = T_WAIT_THEN_COUNT
			}
		} else {
			// Timer stopped, was already stopped
			if (m.newCr & 0x10) != 0 {
				// Force load
				m.state = T_LOAD_THEN_STOP
			}
		}
	} else if m.state == T_COUNT {
		if (m.newCr & 1) != 0 {
			// Timer started, was already running
			if (m.newCr & 0x10) != 0 {
				// Force load
				m.state = T_LOAD_THEN_WAIT_THEN_COUNT
			}
		} else {
			// Timer stopped, was running
			if (m.newCr & 0x10) != 0 {
				// Force load
				m.state = T_LOAD_THEN_STOP
			} else {
				// No force load
				m.state = T_COUNT_THEN_STOP
			}
		}
	} else if m.state == T_LOAD_THEN_COUNT || m.state == T_WAIT_THEN_COUNT {
		if (m.newCr & 1) != 0 {
			// One-shot?
			if (m.newCr & 8) != 0 {
				m.newCr &= 0xfe
				m.state = T_STOP
			} else if (m.newCr & 0x10) != 0 {
				m.state = T_LOAD_THEN_WAIT_THEN_COUNT
			}
		} else {
			m.state = T_STOP
		}
	}
	m.cr = m.newCr & 0xef
	m.hasNewCr = false
}
