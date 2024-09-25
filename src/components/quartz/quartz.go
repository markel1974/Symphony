package quartz

import (
	"container/list"
	"fmt"
)

type Quartz struct {
	cycle           uint64
	alarmsContainer map[*Alarm]*Alarm
	alarms          *list.List
}

func NewQuartz() *Quartz {
	return &Quartz{
		cycle:           0,
		alarmsContainer: make(map[*Alarm]*Alarm),
		alarms:          list.New(),
	}
}

func (s *Quartz) AddCycle() {
	s.cycle++
	s.alarmsCheck(s.cycle)
}

func (s *Quartz) Cycle() uint64 {
	return s.cycle
}

func (s *Quartz) NewAlarm(name string, callback AlarmCallback) *Alarm {
	a := NewAlarm(s, name, callback)
	s.alarmsContainer[a] = a
	return a
}

func (s *Quartz) alarmDestroy(alarm *Alarm) {
	delete(s.alarmsContainer, alarm)
}

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

func (s *Quartz) alarmUnset(alarm *Alarm) error {
	if alarm.element == nil {
		return fmt.Errorf("alarm not setted")
	}
	s.alarms.Remove(alarm.element)
	alarm.element = nil
	return nil
}

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
