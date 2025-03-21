package dynamic_throttle

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/references"
	"time"
)

// DynamicThrottle dynamically regulates task execution intervals to maintain a desired frame rate or time spacing.
type DynamicThrottle struct {
	*component.BaseComponent
	frameInterval int64
	tuning        int64
	prev          int64
	counter       uint64
}

// NewDynamicThrottle creates a new instance of DynamicThrottling with the specified frameInterval in milliseconds.
func NewDynamicThrottle(parent references.IComponent, factory references.IComponentFactory, instance int) *DynamicThrottle {
	d := &DynamicThrottle{
		BaseComponent: component.NewBaseComponent(),
		prev:          time.Now().UnixMilli(),
		frameInterval: 0,
		tuning:        0,
		counter:       0,
	}
	d.BaseComponent.Register(factory, parent, Identifier(), instance, d, references.IdIThrottle(d))
	return d
}

// Setup initializes the DynamicThrottle by setting the desired frame interval in milliseconds.
func (s *DynamicThrottle) Setup(frameInterval int64) error {
	s.frameInterval = frameInterval
	return nil
}

// Reset resets the internal state of the DynamicThrottle, including counters and timestamps, to their initial values.
func (s *DynamicThrottle) Reset() {
}

// Throttle regulates code execution to maintain a consistent time interval between consecutive invocations.
// It calculates the time difference from the previous execution and sleeps if necessary to enforce the interval.
// Adjusts a tuning parameter dynamically to compensate for deviations in interval accuracy.
// Updates the internal state, including the previous execution timestamp and invocation counter.
func (s *DynamicThrottle) Throttle() {
	//https://codereview.stackexchange.com/questions/40473/portable-periodic-one-shot-timer-implementation?noredirect=1&lq=1
	//https://docs.microsoft.com/en-us/windows/win32/api/synchapi/nf-synchapi-signalobjectandwait
	now := time.Now().UnixMilli()
	diff := now - s.prev
	if interval := s.frameInterval - diff; interval < s.frameInterval {
		duration := now + interval
		if sleep := interval - s.tuning; sleep > 1 {
			time.Sleep(time.Duration(sleep) * time.Millisecond)
			now = time.Now().UnixMilli()
		}
		if s.tuning = now - duration; s.tuning < 0 {
			s.tuning = 0
		}
	}
	s.prev = now
	s.counter++
}

// Counter returns the current value of the counter field, which represents the number of throttling operations performed.
func (s *DynamicThrottle) Counter() uint64 {
	return s.counter
}
