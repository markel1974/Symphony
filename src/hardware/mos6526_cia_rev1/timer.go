package mos6526

import (
	"log"

	"github.com/markel1974/c64emu/src/common/signals"
	"github.com/markel1974/c64emu/src/kernel/component"
	"github.com/markel1974/c64emu/src/references"
)

// timerStop represents a state where the timer is stopped.
// timerWaitThenCount represents a state where the timer waits before starting to count.
// timerLoadThenStop represents a state where the timer loads and then stops.
// timerLoadThenCount represents a state where the timer loads and then starts counting.
// timerLoadThenWaitThenCount represents a state where the timer loads, waits, and then starts counting.
// timerCount represents a state where the timer is actively counting.
// timerCountThenStop represents a state where the timer counts and then stops.
const (
	timerStop = iota
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
	crBitStart     = uint8(0x1)  //bit 0
	crBitPBOn      = uint8(0x2)  //bit 1
	crBitOutMode   = uint8(0x4)  //bit 2
	crBitRunMode   = uint8(0x8)  //bit 3
	crBitForceLoad = uint8(0x10) //bit 4
	crBitInMode    = uint8(0x20) //bit 5
	crBitSPMode    = uint8(0x40) //bit 6
	crBitTODIn     = uint8(0x80) //bit 7
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
	countModeTick              = 0 // clock
	countModeCNT               = 1 // positive CNT (Serial Port) transition
	countModeTimerUnderflow    = 2 // timerA underflow
	countModeTimerUnderflowCNT = 3 // timerA underflow while CNT (Serial Port) is high
)

// defaultTimerInit is the initial value assigned to a timer, typically representing a fully loaded state (0xffff).
const defaultTimerInit = 0xffff

// Timer represents a configurable timer with settings for counting modes, latches, and state management.
// It includes mechanisms for updating and handling pending configurations.
// The timer's behavior is determined by the countMode and the corresponding count function.
// It supports various operational states applicable for different use cases.
type Timer struct {
	*component.BaseComponent
	reflect            *TimerReflect
	underflowOutSignal *signals.Signal
	countFn            func(bool) bool
	cr                 uint8  // symphony:export cr represents a configuration register or a control register as an unsigned 8-bit integer.
	crNew              uint8  // symphony:export crNew represents a new configuration register with an 8-bit unsigned integer value.
	crNewPending       bool   // symphony:export crNewPending indicates whether a new change request is in a pending state.
	timer              uint16 // symphony:export timer represents a 16-bit unsigned integer used for timing or count-related operations.
	timerLatch         uint16 // symphony:export timerLatch is a 16-bit value used to store the latch for a timer mechanism.
	timerState         uint8  // symphony:export timerState represents the current state of a timer, encoded as an unsigned 8-bit integer.
	countMode          uint8  // symphony:export countMode represents the mode or operation type for counting, defined as an unsigned 8-bit integer.
	toggleMode         bool   // symphony:export toggleMode indicates whether the current mode is toggled on or off.
	timerLatchLow      uint16 // symphony:export timerLatchLow stores the lower 16 bits of the timer latch value used for timing or counting operations.
	cntPulse           bool   // symphony:export cntPulse indicates whether the pulse counter is enabled or active.
	cntLevel           bool   // symphony:export cntLevel indicates whether the current level is active or enabled as a boolean value.
	underflowIn        bool   // symphony:export underflowIn indicates whether the timer underflow input signal is active or not.
	underflowOut       bool   // symphony:export underflowOut indicates the current state of the underflow output signal for the Timer instance.
}

// NewTimer initializes and returns a new Timer instance with the given parentId and suffix.
// The Timer is set to its default state and its Reset method is called to ensure initialization.
func NewTimer(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *Timer {
	m := &Timer{
		BaseComponent:      component.NewBaseComponent(),
		cr:                 0,
		crNew:              0,
		crNewPending:       false,
		timer:              defaultTimerInit,
		timerLatch:         defaultTimerInit,
		timerState:         timerStop,
		countMode:          countModeTick,
		toggleMode:         false,
		timerLatchLow:      0,
		countFn:            nil,
		cntLevel:           false,
		cntPulse:           false,
		reflect:            nil,
		underflowOutSignal: signals.NewSignal(),
		underflowIn:        false,
		underflowOut:       false,
	}
	//m.BaseComponent.Register(factory, parent, "timer", instance, m, references.IdInternalComponent(label, instance, "Timer"))
	m.reflect = NewTimerReflect(m, factory, parent, "timer", instance, references.IdInternalComponent(label, instance, "Timer"))
	m.Reset()
	return m
}

func (m *Timer) Setup() error {
	return nil
}

func (m *Timer) Connect() error {
	return nil
}

func (m *Timer) Internal() bool {
	return true
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
	m.countFn = m.countTick
	m.cntLevel = false
	m.cntPulse = false
}

// UnderflowSignal returns the signal triggered on timer underflow for further handling or binding.
func (m *Timer) UnderflowSignal() *signals.Signal {
	return m.underflowOutSignal
}

// GetUnderflowOut retrieves the current state of the underflow output signal from the timer instance.
func (m *Timer) GetUnderflowOut() bool {
	return m.underflowOut
}

// SetUnderflowIn sets the underflow input signal state for the timer instance to the specified boolean value.
func (m *Timer) SetUnderflowIn(u bool) {
	m.underflowIn = u
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

// ReadCR retrieves the current value of the control register (CR) for the Timer instance.
func (m *Timer) ReadCR() uint8 {
	return m.cr
}

// ReadTimerLow retrieves the lower 8 bits of the timer's current value and returns it as an unsigned 8-bit integer.
func (m *Timer) ReadTimerLow() uint8 {
	return uint8(m.timer)
}

// ReadTimerHigh returns the high byte (upper 8 bits) of the timer value by shifting the timer's value 8 bits to the right.
func (m *Timer) ReadTimerHigh() uint8 {
	return uint8(m.timer >> 8)
}

// WriteTimerLow sets the lower 8 bits of the timer latch value by converting the given prescaler to a uint16.
func (m *Timer) WriteTimerLow(prescaler uint8) {
	m.timerLatchLow = uint16(prescaler)
}

// WriteTimerHigh configures the high byte of the timer latch using the given prescaler and updates the timer if not started.
func (m *Timer) WriteTimerHigh(prescaler uint8) {
	timerLatchHigh := uint16(prescaler) << 8
	m.timerLatch = m.timerLatchLow | timerLatchHigh
	if (m.cr & crBitStart) == 0 {
		m.timer = m.timerLatch
	}
}

// SetControlRegister updates the control register with new data and sets the count mode for the timer.
func (m *Timer) SetControlRegister(data uint8, countMode uint8) {
	//if m.crNewPending {
	//	fmt.Printf("TIMER %s has cr pending\n", m.GetId())
	//}
	m.crNewPending = true
	m.crNew = data
	m.setCountMode(countMode)
}

// updateCountMode updates the counting mode for the timer based on the provided countMode parameter.
// It sets the appropriate counting function (m.count) and logs unsupported or partially supported modes.
func (m *Timer) setCountMode(countMode uint8) {
	m.countMode = countMode
	switch m.countMode {
	case countModeTick:
		m.countFn = m.countTick
	case countModeCNT:
		m.countFn = m.countCNT
	case countModeTimerUnderflow:
		m.countFn = m.countTimerUnderflow
	case countModeTimerUnderflowCNT:
		m.countFn = m.countTimerUnderflowCNT
	default:
		log.Printf("[timerCount] %s UNSUPPORTED Count Mode %d", m.GetId(), m.countMode)
		m.countFn = m.countTick
	}
}

// Emulate processes the current timer state and performs actions such as counting or toggling based on the timer configuration.
func (m *Timer) Emulate() {
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
		if m.countFn(m.underflowIn) {
			m.toggleMode = !m.toggleMode // Toggle PB6/PB7 output
			if (m.cr & crBitRunMode) != 0 {
				m.timer = m.timerLatch
				m.timerState = timerStop
				m.cr &= crBitStartUnset
			} else {
				m.timer = m.timerLatch
				m.timerState = timerLoadThenCount
			}
			m.underflowOut = true
			m.underflowOutSignal.Emit()
			return
		}
		m.cr &= crBitStartUnset //0xfe
		m.timerState = timerStop
	case timerCount:
		if m.countFn(m.underflowIn) {
			m.toggleMode = !m.toggleMode // Toggle PB6/PB7 output
			if (m.cr & crBitRunMode) != 0 {
				m.timer = m.timerLatch
				m.timerState = timerStop
				m.cr &= crBitStartUnset
			} else {
				m.timer = m.timerLatch
				m.timerState = timerLoadThenCount
			}
			m.underflowOut = true
			m.underflowOutSignal.Emit()
			return
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
	m.underflowOut = false
	return
}

// EmulationRequired checks if emulation is required for the Timer instance and returns true if so.
func (m *Timer) EmulationRequired() bool {
	return true
}

// SetCNTLevel sets the CNT (counter) level for the Timer instance based on the provided boolean value.
func (m *Timer) SetCNTLevel(level bool) {
	m.cntLevel = level
}

// SetCNTPulse enables the CNT pulse by setting the cntPulse field to true.
func (m *Timer) SetCNTPulse() {
	m.cntPulse = true
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
	if m.cntPulse {
		m.cntPulse = false
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
	if underflowX && m.cntLevel {
		if m.timer <= 1 {
			m.timer = 0
			return true
		}
		m.timer--
	}
	return false
}
