package quartz_rev1

import (
	"container/list"
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/component"
	"github.com/markel1974/c64emu/src/references"
)

// Quartz represents a timing system with configurable frequency, alarms, and a processing cycle mechanism.
type Quartz struct {
	*component.BaseComponent
	cycle           uint64
	alarmsContainer map[*Alarm]*Alarm
	alarms          *list.List
	hz              uint64
	factor          float64
}

// NewQuartz creates and initializes a new Quartz instance with the specified parent, factory, label, and instance number.
func NewQuartz(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *Quartz {
	q := &Quartz{
		BaseComponent:   component.NewBaseComponent(),
		cycle:           0,
		alarmsContainer: make(map[*Alarm]*Alarm),
		alarms:          list.New(),
		hz:              0,
		factor:          0,
	}
	q.BaseComponent.Register(factory, parent, Identifier(), instance, q, references.IdIQuartz(q, label, instance))
	return q
}

// Setup initializes the Quartz component and prepares it for operation.
func (s *Quartz) Setup() error {
	return nil
}

// Bind sets the clock frequency (hz) for the Quartz instance and calculates the scale factor based on 1 MHz reference.
func (s *Quartz) Bind(_ references.IQuartzSocket, hz uint64) error {
	s.hz = hz
	s.factor = float64(s.hz) / float64(references.IQuartz1Mhz)
	return nil
}

// Connect establishes the necessary connections or initializes the state for the Quartz component to operate properly.
func (s *Quartz) Connect() error {
	return nil
}

// Internal determines if the Quartz operates as an internal component. Returns false by default.
func (s *Quartz) Internal() bool {
	return false
}

// EmulationRequired returns true to indicate that emulation by the Quartz component is necessary.
func (s *Quartz) EmulationRequired() bool {
	return true
}

// Emulate processes a single emulation cycle by incrementing the cycle counter and checking for any active alarms.
func (s *Quartz) Emulate() {
	s.cycle++
	if s.alarms.Len() > 0 {
		s.alarmsCheck(s.cycle)
	}
}

// Reset clears the internal counters and reinitializes the alarms and container for the Quartz instance.
func (s *Quartz) Reset() {
	s.cycle = 0
	s.alarmsContainer = make(map[*Alarm]*Alarm)
	s.alarms = list.New()
}

// Cycle returns the current cycle count of the Quartz component.
func (s *Quartz) Cycle() uint64 {
	return s.cycle
}

// USecToCycle converts the given microsecond duration `v` into cycles based on the Quartz clock's frequency factor.
func (s *Quartz) USecToCycle(v uint64) float64 {
	return float64(v) * s.factor
}

// USecToCycleRounded converts a given duration in microseconds to the corresponding cycle count (rounded to the nearest integer).
func (s *Quartz) USecToCycleRounded(v uint64) uint64 {
	return uint64(float64(v)*s.factor + 0.5)
}

// NewAlarm creates a new alarm with the specified name and callback, adds it to the alarms container, and returns it.
func (s *Quartz) NewAlarm(name string, callback references.QuartzAlarmCallback) references.IQuartzAlarm {
	a := NewAlarm(s, name, callback)
	s.alarmsContainer[a] = a
	return a
}

// alarmDestroy removes the specified alarm from the alarmsContainer, effectively disabling it.
func (s *Quartz) alarmDestroy(alarm *Alarm) {
	delete(s.alarmsContainer, alarm)
}

// alarmSet schedules an alarm to trigger after a specified cycle distance and adjusts its placement in the alarm list.
func (s *Quartz) alarmSet(alarm *Alarm, dist uint64) error {
	if alarm.element != nil {
		return fmt.Errorf("alarm already set")
	}
	if dist == 0 {
		dist = 1
	}
	alarm.cycle = s.cycle + dist
	if s.alarms.Len() > 0 {
		for e := s.alarms.Front(); e != nil; e = e.Next() {
			curr := e.Value.(*Alarm)
			if alarm.cycle <= curr.cycle {
				alarm.element = s.alarms.InsertBefore(alarm, e)
				break
			}
		}
	}
	if alarm.element == nil {
		alarm.element = s.alarms.PushBack(alarm)
	}
	return nil
}

// alarmUnset removes the specified alarm from the scheduling list. Returns an error if the alarm is not currently set.
func (s *Quartz) alarmUnset(alarm *Alarm) error {
	if alarm.element == nil {
		return fmt.Errorf("alarm not setted")
	}
	s.alarms.Remove(alarm.element)
	alarm.element = nil
	return nil
}

// alarmsCheck processes the alarms linked to the Quartz instance, triggering and removing those due for execution.
func (s *Quartz) alarmsCheck(cycle uint64) {
	var next *list.Element
	for e := s.alarms.Front(); e != nil; e = next {
		next = e.Next()
		alarm := e.Value.(*Alarm)
		if cycle < alarm.cycle {
			break
		}
		alarm.exec(cycle, cycle-alarm.cycle)
		alarm.element = nil
		s.alarms.Remove(e)
	}
}
