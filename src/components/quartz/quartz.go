package quartz

import (
	"container/list"
	"fmt"
)

// Quartz is a type that manages a clock cycle and provides scheduling functionality through alarms.
type Quartz struct {
	cycle           uint64
	alarmsContainer map[*Alarm]*Alarm
	alarms          *list.List
}

// NewQuartz creates and returns a new instance of Quartz with initialized cycle and empty alarm structures.
func NewQuartz() *Quartz {
	return &Quartz{
		cycle:           0,
		alarmsContainer: make(map[*Alarm]*Alarm),
		alarms:          list.New(),
	}
}

// AddCycle increments the internal cycle counter of the Quartz instance and triggers alarm checks for the current cycle.
func (s *Quartz) AddCycle() {
	s.cycle++
	s.alarmsCheck(s.cycle)
}

// Cycle returns the current cycle count maintained by the Quartz instance.
func (s *Quartz) Cycle() uint64 {
	return s.cycle
}

// NewAlarm creates a new alarm with the specified name and callback, associates it with the Quartz instance, and registers it.
func (s *Quartz) NewAlarm(name string, callback AlarmCallback) *Alarm {
	a := NewAlarm(s, name, callback)
	s.alarmsContainer[a] = a
	return a
}

// alarmDestroy removes the specified alarm from the Quartz instance's alarms container.
func (s *Quartz) alarmDestroy(alarm *Alarm) {
	delete(s.alarmsContainer, alarm)
}

// alarmSet schedules an alarm to trigger after a specified distance in cycles. Returns an error if the alarm is already set.
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

// alarmUnset removes an alarm from the scheduled list. Returns an error if the alarm is not currently set.
func (s *Quartz) alarmUnset(alarm *Alarm) error {
	if alarm.element == nil {
		return fmt.Errorf("alarm not setted")
	}
	s.alarms.Remove(alarm.element)
	alarm.element = nil
	return nil
}

// alarmsCheck processes all alarms that are scheduled to execute at or before the given cycle and removes them from the list.
func (s *Quartz) alarmsCheck(cycle uint64) {
	if s.alarms.Len() == 0 {
		return
	}
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
