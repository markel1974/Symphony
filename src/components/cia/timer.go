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
	crBitStart = uint8(0x1) //bit 0
	//1 = TIMER A/B output appears on PB6/PB7 - 0 = PB6/PB7 normal operation.
	crBitPBOn = uint8(0x2) //bit 1
	//1 = TOGGLE - 0 = PULSE
	crBitOutMode = uint8(0x4) //bit 2
	//1 = ONE-SHOT - 0 = CONTINUOUS
	crBitRunMode = uint8(0x8) //bit 3
	//1 = FORCE LOAD (this is a STROBE input, there is no data storage, bit 4 will always read back a zero and writing a zero has no effect).
	crBitForceLoad = uint8(0x10) //bit 4
	//1 = TIMER A counts positive CNT transitions. - 0 = TIMER A counts phi2 pulses.
	crBitInMode = uint8(0x20) //bit 5
	//1 = SERIAL PORT output (CNT sources shift clock) - 0 = SERIAL PORT input (external shift clock required)
	crBitSPMode = uint8(0x40) //bit 6
	//1 = 50Hz clock required on TOD pin for accurate time - 0 = 60Hz clock required on TOD pin for accurate time
	crBitTODIn = uint8(0x80) //bit 7
)

const (
	crBitStartUnset     = ^crBitStart
	crBitForceLoadUnset = ^crBitForceLoad
)

const (
	countModeTick              = 0
	countModeCNT               = 1
	countModeTimerUnderflow    = 2
	countModeTimerUnderflowCNT = 3
)

const defaultTimerInit = 0xffff

//The timer latch is loaded into the timer on any timer underflow.
//The timer latch is loaded into the timer on a force load.
//The timer latch is loaded into the timer after a write to the high byte of the prescaler while the timer is stopped.
//If the timer is running, a write to the high byte will load the timer latch, but not reload the counter

type Timer struct {
	id           string
	cr           uint8
	crNew        uint8      // New values for cr
	crNewPending bool       // New value for crNew pending
	timer        uint16     // Timer
	timerLatch   uint16     // Timer latch
	timerState   TimerState // Timer states
	// 0 = clock; 1 = positive CNT (Serial Port) transition; 2 = timerA underflow; 3 = timerA underflow while CNT (Serial Port) is high
	countMode     uint8
	count         func(bool) bool
	toggleMode    bool
	timerLatchLow uint16
	cnt           bool
}

func NewTimer(id string) *Timer {
	m := &Timer{
		id:            id,
		cr:            0,
		crNew:         0,
		crNewPending:  false,
		timer:         defaultTimerInit,
		timerLatch:    defaultTimerInit,
		timerState:    timerStop,
		countMode:     countModeTick,
		toggleMode:    false,
		timerLatchLow: 0,
		count:         nil,
		cnt:           false,
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
	m.countMode = countModeTick
	m.toggleMode = false
	m.count = m.countTick
	m.cnt = false
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

func (m *Timer) SetTimerLow(prescaler uint8) {
	m.timerLatchLow = uint16(prescaler)
}

func (m *Timer) SetTimerHigh(prescaler uint8) {
	timerLatchHigh := uint16(prescaler) << 8
	m.timerLatch = m.timerLatchLow | timerLatchHigh
	if (m.cr & crBitStart) == 0 {
		m.timer = m.timerLatch
	}
}

func (m *Timer) SetControlRegister(data uint8, countMode uint8) {
	if m.crNewPending {
		fmt.Printf("TIMER %s has cr pending\n", m.id)
	}
	m.crNewPending = true
	m.crNew = data
	m.countMode = countMode
	switch m.countMode {
	case countModeTick:
		m.count = m.countTick
	case countModeCNT:
		log.Printf("[timerCount] %s TODO Count Mode countModeCNT", m.id)
		m.count = m.countCNT
	case countModeTimerUnderflow:
		m.count = m.countTimerUnderflow
	case countModeTimerUnderflowCNT:
		log.Printf("[timerCount] %s TODO Count Mode countModeTimerUnderflowCNT", m.id)
		m.count = m.countTimerUnderflowCNT
	default:
		log.Printf("[timerCount] %s UNSUPPORTED Count Mode %d", m.id, m.countMode)
		m.count = m.countTick
	}
}

func (m *Timer) Emulate(underflowX bool) bool {
	if m.crNewPending {
		m.crNewPending = false
		m.pendingVerify()
	}

	switch m.timerState {
	case timerStop:
		//nothing to do
	case timerLoadThenStop:
		m.timer = m.timerLatch
		m.timerState = timerStop
		m.cr &= crBitStartUnset
	case timerCountThenStop:
		if m.count(underflowX) {
			m.toggleMode = !m.toggleMode // Toggle PB6/PB7 output
			if (m.cr & crBitRunMode) != 0 {
				m.timer = m.timerLatch
				m.timerState = timerStop
				m.cr &= crBitStartUnset
			} else {
				m.timer = m.timerLatch
				m.timerState = timerLoadThenCount
			}
			return true
		}
		m.cr &= crBitStartUnset //0xfe
		m.timerState = timerStop
	case timerCount:
		if m.count(underflowX) {
			m.toggleMode = !m.toggleMode // Toggle PB6/PB7 output
			if (m.cr & crBitRunMode) != 0 {
				m.timer = m.timerLatch
				m.timerState = timerStop
				m.cr &= crBitStartUnset
			} else {
				m.timer = m.timerLatch
				m.timerState = timerLoadThenCount
			}
			return true
		}
	case timerWaitThenCount:
		m.timerState = timerCount
	case timerLoadThenCount:
		m.timer = m.timerLatch
		m.timerState = timerCount
	case timerLoadThenWaitThenCount:
		m.timer = m.timerLatch
		m.timerState = timerWaitThenCount
	}
	return false
}

func (m *Timer) pendingVerify() {
	m.cr = m.crNew & crBitForceLoadUnset //no force load set (strobe) (0xef)

	if (m.crNew & crBitStart) != 0 {
		if m.timerState == timerStop || m.timerState == timerLoadThenStop {
			m.toggleMode = true
		}
		if (m.crNew & crBitForceLoad) != 0 {
			m.timerState = timerLoadThenWaitThenCount
			return
		}
		switch m.timerState {
		case timerStop:
			m.timerState = timerWaitThenCount
		case timerLoadThenStop:
			m.timerState = timerLoadThenWaitThenCount
		case timerCountThenStop:
			m.timerState = timerWaitThenCount
		case timerCount:
		case timerLoadThenWaitThenCount:
		case timerLoadThenCount:
		case timerWaitThenCount:
		}
	} else {
		if (m.crNew & crBitForceLoad) != 0 {
			m.timerState = timerLoadThenStop
			return
		}
		switch m.timerState {
		case timerStop:
		case timerLoadThenStop:
		case timerCountThenStop:
		case timerLoadThenCount:
			m.timerState = timerLoadThenStop
		case timerLoadThenWaitThenCount:
			m.timerState = timerLoadThenStop
		case timerCount:
			m.timerState = timerCountThenStop
		case timerWaitThenCount:
			m.timerState = timerCountThenStop
		}
	}
}

func (m *Timer) countTick(_ bool) bool {
	if m.timer <= 1 {
		m.timer = 0
		return true
	}
	m.timer--
	return false
}

func (m *Timer) countCNT(_ bool) bool {
	if m.cnt {
		if m.timer <= 1 {
			m.timer = 0
			return true
		}
		m.timer--
	}
	return false
}

func (m *Timer) countTimerUnderflow(underflowX bool) bool {
	if underflowX {
		if m.timer <= 1 {
			m.timer = 0
			return true
		}
		m.timer--
	}
	return false
}

func (m *Timer) countTimerUnderflowCNT(underflowX bool) bool {
	if underflowX && m.cnt {
		if m.timer <= 1 {
			m.timer = 0
			return true
		}
		m.timer--
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
