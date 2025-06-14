package references

import (
	"fmt"
)

// IThrottleSocket represents a socket interface used for binding and communication in throttling mechanisms.
type IThrottleSocket interface {
}

// IThrottle is an interface for managing and controlling throttled execution of operations.
// Setup initializes the throttle instance, preparing it for usage.
// Bind associates a throttle with a socket and sets a frame interval for execution control.
// Connect establishes connections or dependencies required for throttle operation.
// Counter retrieves the current count of throttle cycles or executions.
// Update performs an update cycle for the throttle, managing execution intervals or changes.
type IThrottle interface {
	Setup() error

	Bind(socket IThrottleSocket, frameInterval int64) error

	Connect() error

	Counter() uint64

	Update()
}

// IdIThrottle constructs a unique identifier string for an IThrottle instance using label, instance, and interface name.
func IdIThrottle(v IThrottle, label string, instance int) string {
	return IdInternalComponent(label, instance, InterfaceName(&v))
}

// ComponentToIThrottle converts an IComponent to an IThrottle, returning an error if conversion is not possible or input is nil.
func ComponentToIThrottle(component IComponent) (IThrottle, error) {
	if component == nil {
		return nil, fmt.Errorf("component is IThrottle nil")
	}
	v, ok := component.(IThrottle)
	if !ok {
		return nil, fmt.Errorf("component is not a IThrottle")
	}
	return v, nil
}
