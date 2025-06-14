package references

import (
	"fmt"
	"time"
)

// IQuartz1Mhz defines a constant representing the number of clock cycles in 1 MHz based on time.Second and time.Microsecond.
const IQuartz1Mhz = uint64(time.Second / time.Microsecond)

// IQuartzSocket represents an interface used for binding and connecting Quartz components for emulation and timing purposes.
type IQuartzSocket interface {
}

// IQuartz defines an interface for managing clock-based operations in an emulation environment.
// Setup initializes the quartz object.
// Bind associates the quartz object with a socket and a specified frequency.
// Connect establishes necessary connections.
// Cycle retrieves the current cycle count.
// Emulate performs a single step emulation based on the quartz clock.
// USecToCycle converts microseconds to clock cycles as a float value.
// USecToCycleRounded converts microseconds to clock cycles with rounding to the nearest integer.
// NewAlarm creates a new alarm associated with the quartz object and callback.
type IQuartz interface {
	Setup() error

	Bind(socket IQuartzSocket, hz uint64) error

	Connect() error

	Cycle() uint64

	Emulate()

	USecToCycle(uint64) float64

	USecToCycleRounded(v uint64) uint64

	NewAlarm(string, QuartzAlarmCallback) IQuartzAlarm
}

// QuartzAlarmCallback defines a function type used as a callback for Quartz alarms, receiving main cycle and offset values.
type QuartzAlarmCallback func(mainCpuClk uint64, offset uint64)

// IQuartzAlarm defines an interface for managing quartz-based alarms with set, unset, and destroy capabilities.
type IQuartzAlarm interface {
	Set(dist uint64) error

	Unset() error

	Destroy()
}

// IdIQuartz generates a unique identifier string for an IQuartz implementation, combining a label, instance number, and interface name.
func IdIQuartz(v IQuartz, label string, instance int) string {
	return IdInternalComponent(label, instance, InterfaceName(&v))
}

// ComponentToIQuartz converts an IComponent to an IQuartz interface.
// Returns the IQuartz instance if conversion is successful; otherwise, returns an error.
func ComponentToIQuartz(component IComponent) (IQuartz, error) {
	if component == nil {
		return nil, fmt.Errorf("component IQuartz is nil")
	}
	v, ok := component.(IQuartz)
	if !ok {
		return nil, fmt.Errorf("component is not a IQuartz")
	}
	return v, nil
}
