package references

import (
	"fmt"
	"github.com/markel1974/c64emu/src/config"
)

func IdIQuartz(_ IQuartz, instance int) string {
	return IdInternalComponent("IQuartz", instance)
}

type IQuartzSocket interface {
}

// IQuartz defines an interface for managing clock cycles, alarms, and time conversion in an emulation environment.
// Setup initializes the quartz instance and prepares it for usage.
// Cycle retrieves the current clock cycle count.
// Emulate increments the internal cycle count by one.
// ToUSec converts a given clock cycle count to microseconds.
// NewAlarm creates a new alarm instance with a specified name and callback function.
type IQuartz interface {
	Setup(socket IQuartzSocket, cfg *config.Config) error

	Connect() error

	Cycle() uint64

	Emulate()

	ToUSec(uint64) uint64

	NewAlarm(string, QuartzAlarmCallback) IQuartzAlarm
}

// QuartzAlarmCallback defines a function type executed by a Quartz instance, receiving main CPU clock and cycle offset.
type QuartzAlarmCallback func(mainCpuClk uint64, offset uint64)

// IQuartzAlarm defines an interface for managing alarms based on clock cycles or time intervals.
//
// Set sets the alarm to trigger after a specific distance in clock cycles.
// Unset removes the alarm from being triggered.
// Destroy cleans up any internal resources associated with the alarm.
type IQuartzAlarm interface {
	Set(dist uint64) error

	Unset() error

	Destroy()
}

func ComponentToIQuartz(component IComponent) (IQuartz, error) {
	if component == nil {
		return nil, fmt.Errorf("component is nil")
	}
	v, ok := component.(IQuartz)
	if !ok {
		return nil, fmt.Errorf("component is not a %s", IdIQuartz(v, 0))
	}
	return v, nil
}

func ComponentsToIQuartz(cc map[string]IComponent, instance int) (IQuartz, error) {
	quartzId := IdIQuartz(nil, instance)
	c, err := ComponentToIQuartz(cc[quartzId])
	if err != nil {
		return nil, err
	}
	return c, nil
}
