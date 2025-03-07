package mos6526

import (
	"fmt"
	"github.com/markel1974/c64emu/src/components/board"
	"log"
)

// crId is the identifier for 'cr'.
const (
	crId            = "cr"
	crNewId         = "crNew"
	crNewPendingId  = "crNewPending"
	timerId         = "timer"
	timerLatchId    = "timerLatch"
	toggleModeId    = "toggleMode"
	timerLatchLowId = "timerLatchLow"
	cntId           = "cnt"
	timerStateId    = "timerState"
	countModeId     = "countMode"
)

// TimerState represents the state of a timer, typically used to indicate different phases or statuses of a timing mechanism.
type TimerState uint8

// timerStop represents a state where the timer is stopped.
// timerWaitThenCount represents a state where the timer waits before starting to count.
// timerLoadThenStop represents a state where the timer loads and then stops.
// timerLoadThenCount represents a state where the timer loads and then starts counting.
// timerLoadThenWaitThenCount represents a state where the timer loads, waits, and then starts counting.
// timerCount represents a state where the timer is actively counting.
// timerCountThenStop represents a state where the timer counts and then stops.
const (
	timerStop = TimerState(iota)
	timerWaitThenCount
	timerLoadThenStop
	timerLoadThenCount
	timerLoadThenWaitThenCount
	timerCount
	timerCountThenStop
)

// crBitStart represents START/STOP TIMER A control bit: 1 = START TIMER A, 0 = STOP TIMER A (auto reset on one-shot mode underflow).
// crBitPBOn enables TIMER A/B output on PB6/PB7: 1 = output on PB6/PB7, 0 = PB6/PB7 normal operation.
// crBitOutMode configures output mode: 1 = TOGGLE mode, 0 = PULSE mode.
// crBitRunMode sets run mode: 1 = ONE-SHOT mode, 0 = CONTINUOUS mode.
// crBitForceLoad forces a load operation: 1 = FORCE LOAD (strobe input, always reads back as zero, writing zero has no effect).
// crBitInMode selects TIMER A counting mode: 1 = counts positive CNT transitions, 0 = counts phi2 pulses.
// crBitSPMode configures serial port mode: 1 = SERIAL PORT output (CNT sources shift clock), 0 = SERIAL PORT input (external shift clock).
// crBitTODIn specifies TOD pin time accuracy clock: 1 = 50Hz clock required, 0 = 60Hz clock required.
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

// crBitStartUnset defines the inverted value for stopping Timer A operation in one-shot mode when `crBitStart` is unset.
// crBitForceLoadUnset defines the inverted value representing no effect state when `crBitForceLoad` is unset.
const (
	crBitStartUnset     = ^crBitStart
	crBitForceLoadUnset = ^crBitForceLoad
)

// countModeTick represents the counting mode that increments on each tick.
// countModeCNT represents the counting mode that increments based on the CNT signal.
// countModeTimerUnderflow represents the counting mode that increments on timer underflow.
// countModeTimerUnderflowCNT represents the counting mode that increments on both timer underflow and CNT signal.
const (
	countModeTick              = 0
	countModeCNT               = 1
	countModeTimerUnderflow    = 2
	countModeTimerUnderflowCNT = 3
)

// defaultTimerInit is the initial value assigned to a timer, typically representing a fully loaded state (0xffff).
const defaultTimerInit = 0xffff

//The timer latch is loaded into the timer on any timer underflow.
//The timer latch is loaded into the timer on a force load.
//The timer latch is loaded into the timer after a write to the high byte of the prescaler while the timer is stopped.
//If the timer is running, a write to the high byte will load the timer latch, but not reload the counter

// Timer represents a configurable timer with settings for counting modes, latches, and state management.
// It includes mechanisms for updating and handling pending configurations.
// The timer's behavior is determined by the countMode and the corresponding count function.
// It supports various operational states applicable for different use cases.
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

// NewTimer initializes and returns a new Timer instance with the given parentId and suffix.
// The Timer is set to its default state and its Reset method is called to ensure initialization.
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

// Reset reinitializes the Timer's internal state to its default values and clears any pending or current configurations.
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

// GetId retrieves the unique identifier of the Timer instance.
func (m *Timer) GetId() string {
	return m.id
}

// GetParentId retrieves the identifier of the parent associated with the timer instance.
func (m *Timer) GetParentId() string {
	return m.parentId
}

// GetProperties returns a map representing the properties of the timer and their respective metadata, or an error if any occurs.
func (t *Timer) GetProperties() *board.Properties {
	//TODO IMPLEMENT
	return nil
	/*
		return map[string]*board.PropertyInfo{
			crId:            board.MustCreatePropertyInfo(t.cr, "Control Register (CR) of the timer.", false),
			crNewId:         board.MustCreatePropertyInfo(t.crNew, "New value for the Control Register (CR).", false),
			crNewPendingId:  board.MustCreatePropertyInfo(t.crNewPending, "Flag indicating if a new value for CR is pending.", false),
			timerId:         board.MustCreatePropertyInfo(t.timer, "Current value of the timer.", false),
			timerLatchId:    board.MustCreatePropertyInfo(t.timerLatch, "Latch value for the timer.", false),
			toggleModeId:    board.MustCreatePropertyInfo(t.toggleMode, "Flag indicating if the timer is in toggle mode.", false),
			timerLatchLowId: board.MustCreatePropertyInfo(t.timerLatchLow, "Low byte of the timer latch value.", false),
			cntId:           board.MustCreatePropertyInfo(t.cnt, "CNT flag (specific to the CIA).", false),
			timerStateId:    board.MustCreatePropertyInfo(t.timerState, "Current state of the timer.", false),
			countModeId:     board.MustCreatePropertyInfo(t.countMode, "Current count mode of the timer.", false),
		}, nil

	*/
}

// Dump serializes the Timer's internal state into the provided map using predefined keys.
func (m *Timer) Dump(d map[string]interface{}) error {
	//board.DumpSetUint8(d, crId, m.cr)
	//board.DumpSetUint8(d, crNewId, m.crNew)
	//board.DumpSetBool(d, crNewPendingId, m.crNewPending)
	//board.DumpSetUint16(d, timerId, m.timer)
	//board.DumpSetUint16(d, timerLatchId, m.timerLatch)
	//board.DumpSetBool(d, toggleModeId, m.toggleMode)
	//board.DumpSetUint16(d, timerLatchLowId, m.timerLatchLow)
	//board.DumpSetBool(d, cntId, m.cnt)
	//board.DumpSetUint8(d, timerStateId, uint8(m.timerState))
	//board.DumpSetUint8(d, countModeId, m.countMode)
	return nil
}

// Restore populates the Timer's fields using the provided map and handles type-specific conversions. Returns an error if any conversion fails.
func (m *Timer) Restore(d map[string]interface{}) error {
	/*
		for k, v := range d {
			var err error
			switch k {
			case crId:
				err = board.DumpGetUint8(v, &m.cr)
			case crNewId:
				err = board.DumpGetUint8(v, &m.crNew)
			case crNewPendingId:
				err = board.DumpGetBool(v, &m.crNewPending)
			case timerId:
				err = board.DumpGetUint16(v, &m.timer)
			case timerLatchId:
				err = board.DumpGetUint16(v, &m.timerLatch)
			case toggleModeId:
				err = board.DumpGetBool(v, &m.toggleMode)
			case timerLatchLowId:
				err = board.DumpGetUint16(v, &m.timerLatchLow)
			case cntId:
				err = board.DumpGetBool(v, &m.cnt)
			case timerStateId:
				timerState := uint8(m.timerState)
				if err = board.DumpGetUint8(v, &timerState); err == nil {
					m.timerState = TimerState(timerState)
				}
			case countModeId:
				countMode := m.countMode
				if err = board.DumpGetUint8(v, &countMode); err == nil {
					m.updateCountMode(countMode)
				}
			}
			if err != nil {
				return err
			}
		}
	*/
	return nil
}

// HasPBOn checks whether bit `crBitPBOn` is set in the control register `cr`. Returns true if set, otherwise false.
func (m *Timer) HasPBOn() bool {
	return (m.cr & crBitPBOn) != 0
}

// ToggleModeApply applies the toggle mode state based on the control register flags and returns the updated boolean value.
func (m *Timer) ToggleModeApply(d bool) bool {
	if (m.cr & crBitOutMode) != 0 {
		d = m.toggleMode
	}
	return d
}

// GetRTC checks if the RTC (Real-Time Clock) bit is set in the control register and returns true if it is.
func (m *Timer) GetRTC() bool {
	return (m.cr & crBitTODIn) != 0
}

// GetCR retrieves the current value of the control register (CR) for the Timer instance.
func (m *Timer) GetCR() uint8 {
	return m.cr
}

// GetTimerLow retrieves the lower 8 bits of the timer's current value and returns it as an unsigned 8-bit integer.
func (m *Timer) GetTimerLow() uint8 {
	return uint8(m.timer)
}

// GetTimerHigh returns the high byte (upper 8 bits) of the timer value by shifting the timer's value 8 bits to the right.
func (m *Timer) GetTimerHigh() uint8 {
	return uint8(m.timer >> 8)
}

// SetTimerLow sets the lower 8 bits of the timer latch value by converting the given prescaler to a uint16.
func (m *Timer) SetTimerLow(prescaler uint8) {
	m.timerLatchLow = uint16(prescaler)
}

// SetTimerHigh configures the high byte of the timer latch using the given prescaler and updates the timer if not started.
func (m *Timer) SetTimerHigh(prescaler uint8) {
	timerLatchHigh := uint16(prescaler) << 8
	m.timerLatch = m.timerLatchLow | timerLatchHigh
	if (m.cr & crBitStart) == 0 {
		m.timer = m.timerLatch
	}
}

// SetControlRegister updates the control register with new data and sets the count mode for the timer.
func (m *Timer) SetControlRegister(data uint8, countMode uint8) {
	if m.crNewPending {
		fmt.Printf("TIMER %s has cr pending\n", m.id)
	}
	m.crNewPending = true
	m.crNew = data
	m.updateCountMode(countMode)
}

// updateCountMode updates the counting mode for the timer based on the provided countMode parameter.
// It sets the appropriate counting function (m.count) and logs unsupported or partially supported modes.
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

// Emulate processes the current timer state and performs actions such as counting or toggling based on the timer configuration.
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

// pendingVerify handles the state transitions of the timer based on the control register and force load conditions.
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

// countTick decrements the timer if greater than 1 and resets it to 0 when it reaches 1, returning true in that case.
func (m *Timer) countTick(_ bool) bool {
	if m.timer <= 1 {
		m.timer = 0
		return true
	}
	m.timer--
	return false
}

// countCNT checks if the `cnt` flag is true, then decrements the timer and returns true if it underflows, otherwise false.
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

// countTimerUnderflow verifies if the timer should underflow when underflowX is true and decrements the timer if applicable.
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

// countTimerUnderflowCNT checks for a combined condition of underflow and CNT signal to decrement the timer or reset it to 0.
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

// printTimerControlData displays detailed control register bit states for a Timer using the provided data byte.
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
