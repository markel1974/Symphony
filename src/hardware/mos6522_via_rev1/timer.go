package mos6522

import (
	"github.com/markel1974/symphony/src/kernel/component"
	"github.com/markel1974/symphony/src/references"
)

// defaultViaTimeout is a constant representing the default timeout value used to initialize or reset a Timer's counter.
const defaultViaTimeout = 0xffff

// Timer represents a 16-bit timer component used for counting cycles and managing underflow behavior.
type Timer struct {
	*component.BaseComponent
	counter    uint16 // symphony:export counter is a 16-bit register used for counting cycles or steps in the Timer's operation.
	latch      uint16 // symphony:export latch is a 16-bit register used to store a predefined value for reloading the counter during Timer operations.
	clockPulse bool   // symphony:export clockPulse indicates whether the timer's clock signal is active during the current emulation cycle.
	loadSignal bool   // symphony:export loadSignal indicates a flag to signal reloading the counter from the latch value in the Timer.
	underflow  bool   // symphony:export underflow indicates whether the timer has reached zero and triggered an underflow condition during operation.
	reflect    *TimerReflect
}

// NewTimer initializes a new Timer instance, registers it with the specified factory, and sets default counter values.
func NewTimer(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *Timer {
	t := &Timer{
		BaseComponent: component.NewBaseComponent(),
		counter:       defaultViaTimeout,
		latch:         defaultViaTimeout,
		clockPulse:    true,
	}
	t.reflect = NewTimerReflect(t, factory, parent, "timer", instance, references.IdInternalComponent(label, instance, "Timer"))
	return t
}

// Setup initializes the Timer with necessary configurations and ensures it is ready for operation. Returns an error if setup fails.
func (t *Timer) Setup() error {
	return nil
}

// Connect establishes connections or dependencies required by the Timer instance and ensures it is ready for use.
func (t *Timer) Connect() error {
	return nil
}

// EmulationRequired determines whether the timer requires emulation for its functionality. Returns false if not required.
func (t *Timer) EmulationRequired() bool {
	return false
}

// Internal returns a boolean indicating whether the timer is configured for internal operation or signaling.
func (t *Timer) Internal() bool {
	return true
}

// Reset initializes the Timer to its default state, resetting all internal fields to predefined default values.
func (t *Timer) Reset() {
	t.counter = defaultViaTimeout
	t.latch = defaultViaTimeout
	t.clockPulse = true
	t.loadSignal = false
	t.underflow = false
}

// Load sets the loadSignal flag to true, indicating a request to reload the counter from the latch value.
func (t *Timer) Load() {
	t.loadSignal = true
}

// CounterLow returns the lower byte (8 bits) of the 16-bit counter value.
func (t *Timer) CounterLow() uint8 {
	return uint8(t.counter)
}

// CounterHigh returns the high 8 bits of the 16-bit counter value as a uint8.
func (t *Timer) CounterHigh() uint8 {
	return uint8(t.counter >> 8)
}

// LatchLow returns the low byte of the 16-bit `latch` value of the Timer as an unsigned 8-bit integer.
func (t *Timer) LatchLow() uint8 {
	return uint8(t.latch)
}

// LatchHigh returns the high byte (bits 8-15) of the 16-bit latch value as an unsigned 8-bit integer.
func (t *Timer) LatchHigh() uint8 {
	return uint8(t.latch >> 8)
}

// SetClockPulse updates the clock pulse signal for the current emulation cycle.
func (t *Timer) SetClockPulse(clockPulse bool) {
	t.clockPulse = clockPulse
}

// SetLatchHigh sets the high byte of the latch value with the provided data, retaining the low byte unchanged.
func (t *Timer) SetLatchHigh(data uint8) {
	t.latch = (t.latch & 0x00FF) | (uint16(data) << 8)
}

// SetLatchLow sets the low byte of the latch register without altering the high byte.
// symphony:export
func (t *Timer) SetLatchLow(data uint8) {
	t.latch = (t.latch & 0xFF00) | uint16(data)
}

// Emulate performs a single emulation cycle for the timer, handling countdown, reload, and underflow detection.
func (t *Timer) Emulate() {
	t.underflow = false
	if t.loadSignal {
		t.counter = t.latch
		t.loadSignal = false // strobe
	}
	if !t.clockPulse {
		return
	}
	t.counter--
	if t.counter == defaultViaTimeout {
		t.underflow = true
	}
}

// Underflow returns true if the timer has experienced an underflow during the current emulation cycle.
func (t *Timer) Underflow() bool {
	return t.underflow
}
