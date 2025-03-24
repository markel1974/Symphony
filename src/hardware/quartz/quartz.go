package quartz

import (
	"container/list"
	"fmt"
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// Quartz represents a scheduler that manages alarms based on cycles, enabling timed execution of associated callbacks.
// cycle tracks the current cycle count in the Quartz instance.
// alarmsContainer provides a mapping of active alarms for quick accessibility and management.
// alarms stores active alarms in a doubly linked list, sorted by their scheduled cycle execution times.
type Quartz struct {
	*component.BaseComponent
	cycle           uint64
	alarmsContainer map[*Alarm]*Alarm
	alarms          *list.List
}

// NewQuartz creates and returns a new instance of Quartz, initializing its cycle counter, alarms container, and alarms list.
func NewQuartz(parent references.IComponent, factory references.IComponentFactory, instance int) *Quartz {
	q := &Quartz{
		BaseComponent:   component.NewBaseComponent(),
		cycle:           0,
		alarmsContainer: make(map[*Alarm]*Alarm),
		alarms:          list.New(),
	}
	q.BaseComponent.Register(factory, parent, Identifier(), q, references.IdIQuartz(q, instance))
	//q.BaseComponent.Register(parent, q)
	return q
}

func (s *Quartz) Setup(_ references.IQuartzSocket, _ *config.Config) error {
	return nil
}

func (s *Quartz) Connect() error {
	return nil
}

// Emulate increments the internal cycle counter and checks scheduled alarms against the updated cycle value.
func (s *Quartz) Emulate() {
	s.cycle++
	if s.alarms.Len() > 0 {
		s.alarmsCheck(s.cycle)
	}
}

func (m *Quartz) EmulationRequired() bool {
	return true
}

func (s *Quartz) Reset() {
	//
}

// Cycle returns the current cycle count of the Quartz instance.
func (s *Quartz) Cycle() uint64 {
	return s.cycle
}

// ToUSec converts a time value in microseconds (us) to clock cycles based on a specific conversion factor.
func (s *Quartz) ToUSec(v uint64) uint64 {
	return v
}

// NewAlarm creates a new alarm with the given name and callback, associates it with the Quartz instance, and returns it.
func (s *Quartz) NewAlarm(name string, callback references.QuartzAlarmCallback) references.IQuartzAlarm {
	a := NewAlarm(s, name, callback)
	s.alarmsContainer[a] = a
	return a
}

// alarmDestroy removes the specified alarm from the alarmsContainer within the Quartz instance.
func (s *Quartz) alarmDestroy(alarm *Alarm) {
	delete(s.alarmsContainer, alarm)
}

// alarmSet schedules an alarm to trigger after a specified number of cycles. Returns an error if the alarm is already set.
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

// alarmUnset removes the specified alarm from the scheduled list and clears its element reference. Returns an error if the alarm is not set.
func (s *Quartz) alarmUnset(alarm *Alarm) error {
	if alarm.element == nil {
		return fmt.Errorf("alarm not setted")
	}
	s.alarms.Remove(alarm.element)
	alarm.element = nil
	return nil
}

// alarmsCheck executes all alarms whose scheduled cycle is less than or equal to the current cycle and removes them from the list.
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
