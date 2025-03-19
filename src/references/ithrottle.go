package references

const IdIThrottle = "IThrottle"

// IThrottle defines an interface for controlling the execution rate and tracking the count of operations executed.
// Throttle adjusts execution speed to meet a specified rate, ensuring operations run within defined limits.
// Counter returns the current count of throttled operations executed thus far.
type IThrottle interface {
	Setup(int64) error

	Throttle()

	Counter() uint64
}
