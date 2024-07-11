package quartz

import (
	"container/list"
	"fmt"
)

type AlarmCallback func(mainCpuClk uint64, offset uint64)

type Alarm struct {
	quartz    *Quartz
	name      string
	callback  AlarmCallback
	destroyed bool
	cycle     uint64
	element   *list.Element
}

func NewAlarm(quartz *Quartz, name string, callback AlarmCallback) *Alarm {
	return &Alarm{quartz: quartz, name: name, callback: callback, destroyed: false, cycle: 0}
}

func (a *Alarm) Set(cycle uint64) error {
	if a.destroyed {
		return fmt.Errorf("alarm already destroyed")
	}
	a.cycle = cycle
	if err := a.quartz.alarmSet(a); err != nil {
		return err
	}
	return nil
}

func (a *Alarm) Unset() error {
	if a.destroyed {
		return fmt.Errorf("alarm already destroyed")
	}
	if err := a.quartz.alarmUnset(a); err != nil {
		return err
	}
	return nil
}

func (a *Alarm) Destroy() {
	a.destroyed = true
	a.quartz.alarmDestroy(a)
}

func (a *Alarm) exec(cycle uint64, offset uint64) {
	if a.destroyed {
		return
	}
	a.callback(cycle, offset)
	a.element = nil
}
