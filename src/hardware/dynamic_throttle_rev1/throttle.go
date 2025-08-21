package dynamic_throttle_rev1

import (
	"time"

	"github.com/markel1974/c64emu/src/kernel/component"
	"github.com/markel1974/c64emu/src/references"
)

const (
	lagThreshold   = 2 * 1_000_000 // 2ms
	adjustmentStep = 100_000       // 0.1ms
)

// DynamicThrottle dynamically regulates task execution intervals to maintain a desired frame rate or time spacing.
type DynamicThrottle struct {
	*component.BaseComponent
	frameInterval int64
	//tuning        int64
	prev                   int64
	counter                uint64
	lagAccumulator         int64
	idealFrameInterval     int64
	idealFrameIntervalHalf int64
	currentBracket         int
}

// NewDynamicThrottle creates a new instance of DynamicThrottling with the specified frameInterval in milliseconds.
func NewDynamicThrottle(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *DynamicThrottle {
	d := &DynamicThrottle{
		BaseComponent:          component.NewBaseComponent(),
		prev:                   time.Now().UnixNano(),
		frameInterval:          0,
		lagAccumulator:         0,
		idealFrameInterval:     0,
		counter:                0,
		currentBracket:         0,
		idealFrameIntervalHalf: 0,
		//tuning:        0,
	}
	d.BaseComponent.Register(factory, parent, Identifier(), instance, d, references.IdIThrottle(d, label, instance))
	return d
}

func (s *DynamicThrottle) Setup() error {
	return nil
}

func (s *DynamicThrottle) Bind(_ references.IThrottleSocket, frameDistanceMs int64) error {
	s.idealFrameInterval = frameDistanceMs * 1_000_000
	s.idealFrameIntervalHalf = s.idealFrameInterval / 2
	s.frameInterval = s.idealFrameInterval
	return nil
}

// Internal indicates if the `Pic` is set as an internal device. Always returns false in this implementation.
func (s *DynamicThrottle) Internal() bool {
	return false
}

// Connect establishes a connection or initializes any required resources for the DynamicThrottle component.
func (s *DynamicThrottle) Connect() error {
	return nil
}

// Emulate triggers the dynamic throttling process, ensuring consistent execution intervals based on configuration.
func (s *DynamicThrottle) Emulate() {

}

// EmulationRequired determines whether emulation is currently required by the DynamicThrottle instance.
func (s *DynamicThrottle) EmulationRequired() bool {
	return false
}

// Reset resets the internal state of the DynamicThrottle, including counters and timestamps, to their initial values.
func (s *DynamicThrottle) Reset() {
}

// Counter returns the current value of the counter field, which represents the number of throttling operations performed.
func (s *DynamicThrottle) Counter() uint64 {
	return s.counter
}

// Update dynamically adjusts the frame interval to minimize timing lags and manages performance brackets for throttling.
func (s *DynamicThrottle) Update() {
	targetWakeupTime := s.prev + s.frameInterval
	now := time.Now().UnixNano()
	currentLag := now - targetWakeupTime
	s.lagAccumulator += currentLag

	if s.lagAccumulator > lagThreshold {
		s.frameInterval += adjustmentStep
		s.lagAccumulator = 0
	} else if s.lagAccumulator < -lagThreshold {
		if s.frameInterval > s.idealFrameInterval {
			s.frameInterval -= adjustmentStep
		}
		s.lagAccumulator = 0
	}

	// Bracket Detection Ans Notification
	// Calculate which "gear" we are in now.
	// Adding idealFrameInterval/2 prevents oscillations at the boundary.
	newBracket := int((s.frameInterval + s.idealFrameIntervalHalf) / s.idealFrameInterval)
	if newBracket < 1 {
		newBracket = 1
	}
	if newBracket != s.currentBracket {
		s.currentBracket = newBracket
		// OnPerformanceBracketChanged.Emit(s.currentBracket) // Notifica il moltiplicatore (1, 2, 3...)
		//log.Printf("PERFORMANCE BRACKET CHANGED to %dx (Target: %.2fms)", s.currentBracket, float64(s.frameInterval)/1e6)
	}

	sleepDuration := targetWakeupTime - now
	if sleepDuration > 0 {
		time.Sleep(time.Duration(sleepDuration))
	}
	s.prev = targetWakeupTime
	s.counter++
}

/*
func (s *DynamicThrottle) Update() {
	//https://codereview.stackexchange.com/questions/40473/portable-periodic-one-shot-timer-implementation?noredirect=1&lq=1
	//https://docs.microsoft.com/en-us/windows/win32/api/synchapi/nf-synchapi-signalobjectandwait
	now := time.Now().UnixNano()
	diff := now - s.prev
	if interval := s.frameInterval - diff; interval < s.frameInterval {
		duration := now + interval
		if sleep := interval - s.tuning; sleep > 1 {
			time.sleep(time.Duration(sleep) * time.Nanosecond)
			now = time.Now().UnixNano()
		}
		s.tuning = now - duration
		//fmt.Println(s.tuning)
		//if s.tuning < 0 {
		//	s.tuning = 0
		//}
	}
	s.prev = now
	s.counter++
}
*/
