package quartz_rev1

import (
	"container/list"
	"fmt"
	"github.com/markel1974/c64emu/src/references"
)

// AlarmCallback defines a function that is triggered by an alarm and receives the main CPU clock and an offset as parameters.

// Alarm represents a scheduled action, associated with a Quartz instance, capable of being set, unset, and destroyed.
type Alarm struct {
	quartz    *Quartz
	name      string
	callback  references.QuartzAlarmCallback
	destroyed bool
	cycle     uint64
	element   *list.Element
}

// NewAlarm creates and returns a new Alarm instance with a specified Quartz, name, and callback. It initializes the Alarm.
func NewAlarm(quartz *Quartz, name string, callback references.QuartzAlarmCallback) *Alarm {
	return &Alarm{
		quartz:    quartz,
		name:      name,
		callback:  callback,
		destroyed: false,
		cycle:     0,
	}
}

// Set configures the alarm to trigger after a specified distance in cycles. Returns an error if the alarm is destroyed or already set.
func (a *Alarm) Set(dist uint64) error {
	if a.destroyed {
		return fmt.Errorf("alarm already destroyed")
	}
	if err := a.quartz.alarmSet(a, dist); err != nil {
		return err
	}
	return nil
}

// Unset removes the alarm from the scheduled list. Returns an error if the alarm is destroyed or not currently set.
func (a *Alarm) Unset() error {
	if a.destroyed {
		return fmt.Errorf("alarm already destroyed")
	}
	if err := a.quartz.alarmUnset(a); err != nil {
		return err
	}
	return nil
}

// Destroy marks the alarm as destroyed and removes it from the Quartz instance's alarms container.
func (a *Alarm) Destroy() {
	a.destroyed = true
	a.quartz.alarmDestroy(a)
}

// exec triggers the alarm callback if the alarm is not destroyed and clears its element reference.
func (a *Alarm) exec(cycle uint64, offset uint64) {
	if a.destroyed {
		return
	}
	a.callback(cycle, offset)
	a.element = nil
}
