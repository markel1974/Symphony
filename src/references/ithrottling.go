package references

// IThrottling defines an interface for controlling the execution rate and tracking the count of operations executed.
// Throttle adjusts execution speed to meet a specified rate, ensuring operations run within defined limits.
// Counter returns the current count of throttled operations executed thus far.
type IThrottling interface {
	Throttle()

	Counter() uint64
}
