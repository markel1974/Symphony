package mos6526

import (
	"github.com/markel1974/c64emu/src/signals"
	"log"
)

type TimerState uint8

const (
	timerStop = TimerState(iota)
	timerWaitThenCount
	timerLoadThenStop
	timerLoadThenCount
	timerLoadThenWaitThenCount
	timerCount
	timerCountThenStop
)

const (
	//1 = START TIMER A - 0 = STOP TIMER A s (This bit is automatically reset when underflow occurs during one-shot mode).
	crBitStart = 0x1 //bit 0
	//1 = TIMER A output appears on PB6 - 0 = PB6 normal operation.
	crBitPBOn = 0x2 //bit 1
	//1 = TOGGLE - 0 = PULSE
	crBitOutMode = 0x4 //bit 2
	//1 = ONE-SHOT - 0 = CONTINUOUS
	crBitRunMode = 0x8 //bit 3
	//1 = FORCE LOAD (this is a STROBE input, there is no data storage, bit 4 will always read back a zero and writing a zero has no effect).
	crBitForceLoad = 0x10 //bit 4
	//1 = TIMER A counts positive CNT transitions. - 0 = TIMER A counts phi2 pulses.
	crBitInMode = 0x20 //bit 5
	//1 = SERIAL PORT output (CNT sources shift clock) - 0 = SERIAL PORT input (external shift clock required)
	crBitSPMode = 0x40 //bit 6
	//1 = 50Hz clock required on TOD pin for accurate time - 0 = 60Hz clock required on TOD pin for accurate time
	crBitTODIn = 0x80 //bit 7
)

const defaultTimerInit = 0xffff

type Timer struct {
	id              string
	signalUnderflow *signals.Signal
	cr              uint8
	crNew           uint8      // New values for cr
	crPending       bool       // New value for crNew pending
	timer           uint16     // Timer
	timerLatch      uint16     // Timer latch
	timerState      TimerState // Timer states
	countMode       uint8
}

func NewTimer(id string) *Timer {
	m := &Timer{
		id:              id,
		signalUnderflow: signals.NewSignal(),
		cr:              0,
		crNew:           0,
		crPending:       false,
		timer:           defaultTimerInit,
		timerLatch:      defaultTimerInit,
		timerState:      timerStop,
		countMode:       0,
	}
	m.Reset()
	return m
}

func (m *Timer) Reset() {
	m.cr = 0
	m.crNew = 0
	m.crPending = false
	m.timer = defaultTimerInit
	m.timerLatch = defaultTimerInit
	m.timerState = timerStop
	m.countMode = 0
}

func (m *Timer) SignalUnderflowBind(fn func()) {
	m.signalUnderflow.Bind(fn)
}

func (m *Timer) GetRTC() bool {
	return (m.cr & crBitTODIn) != 0
}

func (m *Timer) GetCR() uint8 {
	return m.cr
}

func (m *Timer) GetTimerLow() uint8 {
	return uint8(m.timer)
}

func (m *Timer) GetTimerHigh() uint8 {
	return uint8(m.timer >> 8)
}

func (m *Timer) SetTimerLow(data uint8) {
	timerLow := uint16(data)
	timerHigh := m.timerLatch & 0xff00
	m.timerLatch = timerLow | timerHigh
	if (m.cr & crBitForceLoad) != 0 {
		m.timer = m.timerLatch // Reload
	}
}

func (m *Timer) SetTimerHigh(data uint8) {
	timerLow := m.timerLatch & 0xff
	timerHigh := uint16(data) << 8
	m.timerLatch = timerLow | timerHigh
	if (m.cr&crBitStart) == 0 || (m.cr&crBitForceLoad) != 0 {
		m.timer = m.timerLatch // Reload
	}

}

func (m *Timer) SetControlRegister(data uint8, countMode uint8) {
	//m.printTimerControlData(data)
	m.crPending = true
	m.crNew = data
	m.countMode = countMode

	if (m.crNew & crBitPBOn) != 0 {
		log.Printf("[SetControlRegister] %s Unimplemented TIMER A on PB6", m.id)
	}
	if (m.crNew & crBitOutMode) != 0 {
		log.Printf("[SetControlRegister] %s Unimplemented OUT MODE", m.id)
	}
	if (m.crNew & crBitSPMode) != 0 {
		log.Printf("[SetControlRegister] %s Unimplemented SERIAL PORT output", m.id)
	}
	//count: 0 clock - 1 positive CNT (Serial Port) transition; 2 - timerA underflow pulse - 3 timerA underflow pulse while CNT (Serial Port) is high
}

func (m *Timer) Emulate(underflowTimerX bool) bool {
	switch m.timerState {
	case timerWaitThenCount:
		m.timerState = timerCount
		m.checkPending()
		return false
	case timerStop:
		m.checkPending()
		return false
	case timerLoadThenStop:
		m.timerState = timerStop
		m.timer = m.timerLatch // Reload
		m.checkPending()
		return false
	case timerLoadThenCount:
		m.timerState = timerCount
		m.timer = m.timerLatch // Reload
		m.checkPending()
		return false
	case timerLoadThenWaitThenCount:
		m.timerState = timerWaitThenCount
		if m.timer == 1 {
			underflow := m.timerCount(true, underflowTimerX)
			m.checkPending()
			return underflow
		}
		m.timer = m.timerLatch // Reload
		m.checkPending()
		return false
	case timerCount:
		underflow := m.timerCount(false, underflowTimerX)
		m.checkPending()
		return underflow
	case timerCountThenStop:
		m.timerState = timerStop
		underflow := m.timerCount(false, underflowTimerX)
		m.checkPending()
		return underflow
	}
	//never happen
	return false
}

func (m *Timer) timerCount(signal bool, underflowTimerX bool) bool {
	underflow := false
	if !signal {
		count := false
		if m.countMode == 0 {
			count = true
		} else if m.countMode == 2 {
			if underflowTimerX {
				count = true
			}
		} else {
			log.Printf("[timerCount] %s UNSUPPORTED Timer counts CNT %d", m.id, m.countMode)
		}
		if count {
			timer := m.timer
			m.timer--
			if (timer == 0) || (m.timer == 0) {
				underflow = true
				if m.timerState != timerStop {
					signal = true
				}
			}
		}
	}
	if signal {
		underflow = true
		m.signalUnderflow.Emit()
		if (m.cr & crBitRunMode) != 0 {
			m.timer = m.timerLatch // Reload timer
			// stop timer
			m.cr &= 0xfe
			m.crNew &= 0xfe
			m.timerState = timerLoadThenStop // Reload in next cycle
		} else {
			m.timer = m.timerLatch            // Reload timer
			m.timerState = timerLoadThenCount // Delay one cycle (and reload)
		}
	}
	return underflow
}

func (m *Timer) checkPending() {
	// Delayed write to CR?
	if m.crPending {
		switch m.timerState {
		case timerStop, timerLoadThenStop:
			// Timer started, wasn't running
			if (m.crNew & crBitStart) != 0 {
				if (m.crNew & crBitForceLoad) != 0 {
					m.timerState = timerLoadThenWaitThenCount
				} else {
					m.timerState = timerWaitThenCount
				}
			} else {
				// Timer stopped, was already stopped
				if (m.crNew & crBitForceLoad) != 0 {
					m.timerState = timerLoadThenStop
				}
			}
		case timerCount:
			if (m.crNew & crBitStart) != 0 {
				// Timer started, was already running
				if (m.crNew & crBitForceLoad) != 0 {
					m.timerState = timerLoadThenWaitThenCount
				}
			} else {
				// Timer stopped, was running
				if (m.crNew & crBitForceLoad) != 0 {
					m.timerState = timerLoadThenStop
				} else {
					m.timerState = timerCountThenStop
				}
			}
		case timerLoadThenCount, timerWaitThenCount:
			if (m.crNew & crBitStart) != 0 {
				if (m.crNew & crBitRunMode) != 0 {
					// One-shot, stop timer
					m.crNew &= 0xfe
					m.timerState = timerStop
				} else if (m.crNew & crBitForceLoad) != 0 {
					// No One-shot, force load
					m.timerState = timerLoadThenWaitThenCount
				}
			} else {
				m.timerState = timerStop
			}
		default:
			log.Printf("[checkPending] %s TIMER - UNDEFINED Timer %d", m.id, m.timerState)
		}
		//no force load set
		m.cr = m.crNew & 0xef
		m.crPending = false
	}
}

/*
func (m *Timer) printTimerControlData(data uint8) {
	fmt.Printf("\n")
	fmt.Printf("%s Timer Control -> crBitStart: %v\n", m.id, data&crBitStart != 0)
	fmt.Printf("%s Timer Control -> crBitSignalNoUnderflow: %v\n", m.id, data&crBitPBOn != 0)
	fmt.Printf("%s Timer Control -> crBitSignalUnderflowInverted: %v\n", m.id, data&crBitOutMode != 0)
	fmt.Printf("%s Timer Control -> crBitRunMode: %v\n", m.id, data&crBitRunMode != 0)
	fmt.Printf("%s Timer Control -> crBitForceLoad: %v\n", m.id, data&crBitForceLoad != 0)
	fmt.Printf("%s Timer Control -> crBitInMode: %v\n", m.id, data&crBitInMode != 0)
	fmt.Printf("%s Timer Control -> crBitSPMode: %v\n", m.id, data&crBitSPMode != 0)
	fmt.Printf("%s Timer Control -> crBitTODIn: %v\n", m.id, data&crBitTODIn != 0)
	fmt.Printf("\n")
}
*/
