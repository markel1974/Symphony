package mos6526

import (
	"fmt"
	"log"
)

// TimerState represents the state of a timer, defined as an unsigned 8-bit integer.
type TimerState uint8

// timerStop represents a state where the timer is stopped.
// timerWaitThenCount represents a state where the timer waits before starting to count.
// timerLoadThenStop represents a state where the timer loads a value and then transitions to a stopped state.
// timerLoadThenCount represents a state where the timer loads a value and starts counting.
// timerLoadThenWaitThenCount represents a state where the timer loads a value, waits, and then starts counting.
// timerCount represents a state where the timer is actively counting.
// timerCountThenStop represents a state where the timer counts and then transitions to a stopped state.
const (
	timerStop = TimerState(iota)
	timerWaitThenCount
	timerLoadThenStop
	timerLoadThenCount
	timerLoadThenWaitThenCount
	timerCount
	timerCountThenStop
)

// crBitStart controls TIMER A start/stop: 1 = start, 0 = stop. Auto-resets on underflow in one-shot mode.
// crBitPBOn enables TIMER A/B output on PB6/PB7: 1 = output on, 0 = normal PB6/PB7 operation.
// crBitOutMode sets output mode: 1 = toggle, 0 = pulse.
// crBitRunMode sets run mode: 1 = one-shot, 0 = continuous.
// crBitForceLoad forces load: 1 = strobe input, always reads back as 0, writing 0 has no effect.
// crBitInMode sets TIMER A count mode: 1 = counts positive CNT transitions, 0 = counts phi2 pulses.
// crBitSPMode sets SERIAL PORT mode: 1 = output with CNT as shift clock, 0 = input requiring external shift clock.
// crBitTODIn configures TOD pin clock requirement: 1 = 50Hz clock, 0 = 60Hz clock.
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

// crBitStartUnset is the inverted value of crBitStart, representing a stop or reset state for the START TIMER A bit.
// crBitForceLoadUnset is the inverted value of crBitForceLoad, indicating a reset or no effect state for the FORCE LOAD bit.
const (
	crBitStartUnset     = ^crBitStart
	crBitForceLoadUnset = ^crBitForceLoad
)

// Constants defining various count modes for a timer system or counter logic.
const (
	countModeTick              = 0
	countModeCNT               = 1
	countModeTimerUnderflow    = 2
	countModeTimerUnderflowCNT = 3
)

// defaultTimerInit is the initial value for timers, set to the maximum value (0xffff).
const defaultTimerInit = 0xffff

//The timer latch is loaded into the timer on any timer underflow.
//The timer latch is loaded into the timer on a force load.
//The timer latch is loaded into the timer after a write to the high byte of the prescaler while the timer is stopped.
//If the timer is running, a write to the high byte will load the timer latch, but not reload the counter

// Timer represents a configurable timer capable of operating in various modes and states with latch functionality.
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

// NewTimer initializes and returns a new Timer instance with the given ID and resets its internal state.
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

// Reset resets the Timer's internal state, control registers, and modes to their default values.
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

// HasPBOn checks if the TIMER A/B output is enabled to appear on PB6/PB7, based on the control register configuration.
func (m *Timer) HasPBOn() bool {
	return (m.cr & crBitPBOn) != 0
}

// ToggleModeApply checks the toggle mode of the timer and applies it if the output mode bit is set, returning the result.
func (m *Timer) ToggleModeApply(d bool) bool {
	if (m.cr & crBitOutMode) != 0 {
		d = m.toggleMode
	}
	return d
}

// GetRTC checks if the RTC (Real-Time Clock) mode bit is set in the control register and returns true if active.
func (m *Timer) GetRTC() bool {
	return (m.cr & crBitTODIn) != 0
}

// GetCR retrieves the current value of the control register (cr) for the Timer instance as an unsigned 8-bit integer.
func (m *Timer) GetCR() uint8 {
	return m.cr
}

// GetTimerLow returns the lower 8 bits of the timer value as a uint8.
func (m *Timer) GetTimerLow() uint8 {
	return uint8(m.timer)
}

// GetTimerHigh returns the high byte of the 16-bit timer value as an unsigned 8-bit integer.
func (m *Timer) GetTimerHigh() uint8 {
	return uint8(m.timer >> 8)
}

// SetTimerLow updates the low byte of the timer latch with the given prescaler value.
func (m *Timer) SetTimerLow(prescaler uint8) {
	m.timerLatchLow = uint16(prescaler)
}

// SetTimerHigh sets the high byte of the timer latch using the given prescaler and updates the timer if it's not started.
func (m *Timer) SetTimerHigh(prescaler uint8) {
	timerLatchHigh := uint16(prescaler) << 8
	m.timerLatch = m.timerLatchLow | timerLatchHigh
	if (m.cr & crBitStart) == 0 {
		m.timer = m.timerLatch
	}
}

// SetControlRegister updates the control register and sets the counting mode of the Timer. Pending changes are marked accordingly.
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

// Emulate performs state-based emulation of the timer based on its current configuration and underflow triggers.
// Returns true if the timer underflow condition is met during emulation.
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

// pendingVerify updates the internal timer state and control register based on the current configuration and new control bits.
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

// countTick decrements the timer value. If the timer reaches zero or below, it resets to zero and returns true.
func (m *Timer) countTick(_ bool) bool {
	if m.timer <= 1 {
		m.timer = 0
		return true
	}
	m.timer--
	return false
}

// countCNT decreases the timer by 1 if the cnt field is true and the timer is above 1. Returns true if timer reaches 0.
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

// countTimerUnderflow decrements the timer when underflowX is true, resets it to 0 if it reaches 1, and returns true if underflow occurred.
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

// countTimerUnderflowCNT decrements the timer if both underflowX and CNT are true, resetting it to 0 if it reaches 1, returning true.
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

// printTimerControlData outputs the Timer's control register bit states to the console for debugging purposes.
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
