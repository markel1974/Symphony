package mos6526

import (
	"fmt"
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
	//1 = TIMER A/B output appears on PB6/PB7 - 0 = PB6/PB7 normal operation.
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
	id           string
	cr           uint8
	crNew        uint8      // New values for cr
	crNewPending bool       // New value for crNew pending
	timer        uint16     // Timer
	timerLatch   uint16     // Timer latch
	timerState   TimerState // Timer states
	// 0 = clock; 1 = positive CNT (Serial Port) transition; 2 = timerA underflow; 3 = timerA underflow while CNT (Serial Port) is high
	countMode  uint8
	toggleMode bool
}

func NewTimer(id string) *Timer {
	m := &Timer{
		id:           id,
		cr:           0,
		crNew:        0,
		crNewPending: false,
		timer:        defaultTimerInit,
		timerLatch:   defaultTimerInit,
		timerState:   timerStop,
		countMode:    0,
		toggleMode:   false,
	}
	m.Reset()
	return m
}

func (m *Timer) Reset() {
	m.cr = 0
	m.crNew = 0
	m.crNewPending = false
	m.timer = defaultTimerInit
	m.timerLatch = defaultTimerInit
	m.timerState = timerStop
	m.countMode = 0
	m.toggleMode = false
}

func (m *Timer) HasPBOn() bool {
	return (m.cr & crBitPBOn) != 0
}

func (m *Timer) ToggleModeApply(d bool) bool {
	if (m.cr & crBitOutMode) != 0 {
		d = m.toggleMode
	}
	return d
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
	//if m.timerLatch == 0 {
	//	fmt.Println("SetTimerLow: timerLatch => 0")
	//}
	if (m.cr & crBitForceLoad) != 0 {
		m.timer = m.timerLatch
	}
	//if m.id == "cia2_TIMER_B" {
	//	fmt.Printf("%s [SetTimerLow] %d\n", m.id, m.timerLatch)
	//}
}

func (m *Timer) SetTimerHigh(data uint8) {
	timerLow := m.timerLatch & 0xff
	timerHigh := uint16(data) << 8
	m.timerLatch = timerLow | timerHigh

	//if m.timerLatch == 0 {
	//	fmt.Println("SetTimerHigh: timerLatch => 0")
	//}
	if (m.cr&crBitStart) == 0 || (m.cr&crBitForceLoad) != 0 {
		m.timer = m.timerLatch
	}
	//if m.id == "cia2_TIMER_B" {
	//	fmt.Printf("%s [SetTimerHigh] %d\n", m.id, m.timerLatch)
	//}
}

func (m *Timer) SetControlRegister(data uint8, countMode uint8) {
	m.crNewPending = true
	m.crNew = data
	m.countMode = countMode
	//if m.id == "cia2_TIMER_B" {
	//	fmt.Printf("%s [SetControlRegister] %8b [CONTINUOS: %v]\n", m.id, m.crNew, (m.crNew&crBitRunMode) == 0)
	//}
}

func (m *Timer) Emulate(underflowX bool) bool {
	underflow := false
	switch m.timerState {
	case timerWaitThenCount:
		m.timerState = timerCount
	case timerStop:
		//nothing to do
	case timerLoadThenStop:
		m.timer = m.timerLatch
		m.timerState = timerStop
	case timerLoadThenCount:
		m.timer = m.timerLatch
		m.timerState = timerCount
	case timerLoadThenWaitThenCount:
		if m.timer == 1 {
			underflow = true
		} else {
			m.timer = m.timerLatch
		}
		m.timerState = timerWaitThenCount
	case timerCount:
		underflow = m.count(underflowX)
	case timerCountThenStop:
		underflow = m.count(underflowX)
		m.timerState = timerStop
	}

	if underflow {
		m.toggleMode = !m.toggleMode // Toggle PB6/PB7 output
		if (m.cr & crBitRunMode) != 0 {
			m.cr &= 0xfe                     // stop timer
			m.crNew &= 0xfe                  // stop timer
			m.timer = m.timerLatch           // Reload timer
			m.timerState = timerLoadThenStop // Reload in next cycle
		} else {
			m.timer = m.timerLatch            // Reload timer
			m.timerState = timerLoadThenCount // Reload in next cycle
		}
	}

	if m.crNewPending {
		switch m.timerState {
		case timerStop, timerLoadThenStop:
			// Timer started, wasn't running
			if (m.crNew & crBitStart) != 0 {
				m.toggleMode = true
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
			log.Printf("[Emulate] %s TIMER - UNDEFINED Timer %d", m.id, m.timerState)
		}
		//no force load set
		m.cr = m.crNew & 0xef
		m.crNewPending = false
	}

	return underflow
}

var _unsupportedPrinted = false

func (m *Timer) count(underflowX bool) bool {
	if m.countMode == 0 {
		if m.timer <= 1 {
			m.timer = 0
			return true
		}
		m.timer--
		return false
	}
	if m.countMode == 2 {
		if underflowX {
			if m.timer <= 1 {
				m.timer = 0
				return true
			}
			m.timer--
		}
		return false
	}
	// TODO UNSUPPORTED!!!!!
	if !_unsupportedPrinted {
		log.Printf("[timerCount] %s UNSUPPORTED Timer counts CNT %d", m.id, m.countMode)
		_unsupportedPrinted = true
	}
	return false
}

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
