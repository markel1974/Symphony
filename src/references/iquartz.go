package references

// IQuartz defines an interface for time-related operations within the emulation environment.
// Cycle retrieves the current clock cycle of the system.
// ToUSec converts a given clock cycle count into microseconds.
type IQuartz interface {
	Cycle() uint64

	AddCycle()

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
