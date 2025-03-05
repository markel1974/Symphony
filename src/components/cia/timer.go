package mos6526

import (
	"fmt"
	"github.com/markel1974/c64emu/src/components/board"
	"log"
)

// TimerState represents the states of a timer, defined as an unsigned 8-bit integer type.
type TimerState uint8

// timerStop represents the state where the timer is stopped.
// timerWaitThenCount represents the state where the timer waits before starting to count.
// timerLoadThenStop represents the state where the timer loads a value and then stops immediately.
// timerLoadThenCount represents the state where the timer loads a value and then starts counting.
// timerLoadThenWaitThenCount represents the state where the timer loads a value, waits, and then starts counting.
// timerCount represents the state where the timer is actively counting.
// timerCountThenStop represents the state where the timer counts and then stops.
const (
	timerStop = TimerState(iota)
	timerWaitThenCount
	timerLoadThenStop
	timerLoadThenCount
	timerLoadThenWaitThenCount
	timerCount
	timerCountThenStop
)

// crBitStart represents the control bit to start (1) or stop (0) TIMER A, automatically resets in one-shot mode on underflow.
// crBitPBOn represents the control bit to enable (1) TIMER A/B output on PB6/PB7 or retain (0) normal operation for PB6/PB7.
// crBitOutMode represents the output mode control bit, toggle (1) or pulse (0).
// crBitRunMode represents the run mode control bit, one-shot (1) or continuous (0).
// crBitForceLoad represents the force load control bit, acts as a strobe input; always reads 0 and ignores zero writes.
// crBitInMode represents the input source control bit; TIMER A counts CNT transitions (1) or phi2 pulses (0).
// crBitSPMode represents the serial port mode control bit; output (1) with CNT as shift clock or input (0) requiring external clock.
// crBitTODIn represents the TOD (Time of Day) clock source control bit; 50Hz required (1) or 60Hz required (0) on TOD pin.
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

// crBitStartUnset represents the negation of the crBitStart constant, used to unset the START TIMER A bit.
// crBitForceLoadUnset represents the negation of the crBitForceLoad constant, used to unset the FORCE LOAD bit.
const (
	crBitStartUnset     = ^crBitStart
	crBitForceLoadUnset = ^crBitForceLoad
)

// countModeTick represents the counting mode set to tick-based.
const (
	countModeTick              = 0
	countModeCNT               = 1
	countModeTimerUnderflow    = 2
	countModeTimerUnderflowCNT = 3
)

// defaultTimerInit is the default initial value for the timer and timer latch constants, set to 0xffff.
const defaultTimerInit = 0xffff

//The timer latch is loaded into the timer on any timer underflow.
//The timer latch is loaded into the timer on a force load.
//The timer latch is loaded into the timer after a write to the high byte of the prescaler while the timer is stopped.
//If the timer is running, a write to the high byte will load the timer latch, but not reload the counter

// Timer is a structure representing a timer with configurable control registers and modes of operation.
type Timer struct {
	parentId     string
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

// NewTimer initializes and returns a new instance of Timer with the provided id, setting default states and values.
func NewTimer(parentId string, suffix string) *Timer {
	m := &Timer{
		parentId:      parentId,
		id:            "timer" + suffix,
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

// Reset reinitializes the Timer fields to their default states, clearing all states and setting default configurations.
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

func (m *Timer) GetId() string {
	return m.id
}

func (m *Timer) GetParentId() string {
	return m.parentId
}

// Dump serializes the Timer's internal state into a map with string keys and interface{} values.
func (m *Timer) Dump(d map[string]interface{}) error {
	board.DumpSetUint8(d, []string{m.id, "cr"}, m.cr)
	board.DumpSetUint8(d, []string{m.id, "crNew"}, m.crNew)
	board.DumpSetBool(d, []string{m.id, "crNewPending"}, m.crNewPending)
	board.DumpSetUint16(d, []string{m.id, "timer"}, m.timer)
	board.DumpSetUint16(d, []string{m.id, "timerLatch"}, m.timerLatch)
	board.DumpSetBool(d, []string{m.id, "toggleMode"}, m.toggleMode)
	board.DumpSetUint16(d, []string{m.id, "timerLatchLow"}, m.timerLatchLow)
	board.DumpSetBool(d, []string{m.id, "cnt"}, m.cnt)
	board.DumpSetUint8(d, []string{m.id, "timerState"}, uint8(m.timerState))
	board.DumpSetUint8(d, []string{m.id, "countMode"}, m.countMode)
	return nil
}

// Restore restores the Timer's state from a provided data map and returns an error if any step fails.
func (m *Timer) Restore(d map[string]interface{}) error {
	_ = board.DumpGetUint8(d, []string{m.id, "cr"}, &m.cr)
	_ = board.DumpGetUint8(d, []string{m.id, "crNew"}, &m.crNew)
	_ = board.DumpGetBool(d, []string{m.id, "crNewPending"}, &m.crNewPending)
	_ = board.DumpGetUint16(d, []string{m.id, "timer"}, &m.timer)
	_ = board.DumpGetUint16(d, []string{m.id, "timerLatch"}, &m.timerLatch)
	_ = board.DumpGetBool(d, []string{m.id, "toggleMode"}, &m.toggleMode)
	_ = board.DumpGetUint16(d, []string{m.id, "timerLatchLow"}, &m.timerLatchLow)
	_ = board.DumpGetBool(d, []string{m.id, "cnt"}, &m.cnt)
	timerState := uint8(m.timerState)
	if ok := board.DumpGetUint8(d, []string{m.id, "timerState"}, &timerState); ok {
		m.timerState = TimerState(timerState)
	}
	countMode := m.countMode
	if ok := board.DumpGetUint8(d, []string{m.id, "countMode"}, &countMode); ok {
		m.updateCountMode(countMode)
	}
	return nil
}

// HasPBOn checks if the TIMER A/B output is set to appear on PB6/PB7, based on the state of the control register (cr).
func (m *Timer) HasPBOn() bool {
	return (m.cr & crBitPBOn) != 0
}

// ToggleModeApply modifies the value of `d` based on the current toggle mode and returns the resulting value.
func (m *Timer) ToggleModeApply(d bool) bool {
	if (m.cr & crBitOutMode) != 0 {
		d = m.toggleMode
	}
	return d
}

// GetRTC determines if a 50Hz clock is required on the TOD pin for accurate time by checking the crBitTODIn flag in the control register.
func (m *Timer) GetRTC() bool {
	return (m.cr & crBitTODIn) != 0
}

// GetCR retrieves the control register (CR) value of the Timer instance.
func (m *Timer) GetCR() uint8 {
	return m.cr
}

// GetTimerLow returns the lower 8 bits of the timer value as a uint8.
func (m *Timer) GetTimerLow() uint8 {
	return uint8(m.timer)
}

// GetTimerHigh retrieves the high byte of the current timer value by performing a bitwise right shift and casting to uint8.
func (m *Timer) GetTimerHigh() uint8 {
	return uint8(m.timer >> 8)
}

// SetTimerLow sets the lower 8 bits of the timer latch using the provided prescaler value.
func (m *Timer) SetTimerLow(prescaler uint8) {
	m.timerLatchLow = uint16(prescaler)
}

// SetTimerHigh sets the higher 8 bits of the timer latch using the provided prescaler value and updates the timer if not running.
func (m *Timer) SetTimerHigh(prescaler uint8) {
	timerLatchHigh := uint16(prescaler) << 8
	m.timerLatch = m.timerLatchLow | timerLatchHigh
	if (m.cr & crBitStart) == 0 {
		m.timer = m.timerLatch
	}
}

// SetControlRegister sets the control register to the specified data and updates the counting mode accordingly.
func (m *Timer) SetControlRegister(data uint8, countMode uint8) {
	if m.crNewPending {
		fmt.Printf("TIMER %s has cr pending\n", m.id)
	}
	m.crNewPending = true
	m.crNew = data
	m.updateCountMode(countMode)
}

// updateCountMode updates the count mode of the Timer and sets the corresponding count function based on the mode provided.
func (m *Timer) updateCountMode(countMode uint8) {
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

// Emulate simulates the timer's behavior based on the current state and control register, returning true for certain transitions.
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

// pendingVerify updates the timer's control register and determines the next timer state based on current settings.
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

// countTick decrements the timer and returns true if the timer reaches zero or below; otherwise, it returns false.
func (m *Timer) countTick(_ bool) bool {
	if m.timer <= 1 {
		m.timer = 0
		return true
	}
	m.timer--
	return false
}

// countCNT decrements the timer if the CNT flag is true and returns true if the timer reaches zero; otherwise, returns false.
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

// countTimerUnderflow decrements the timer if underflowX is true, and returns true if the timer reaches 0.
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

// countTimerUnderflowCNT counts down the timer when both underflowX is true and the CNT flag is set, returning true on underflow.
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

// printTimerControlData prints the control register bit state of the timer based on the provided `data` value.
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
