package references

import (
	"fmt"
)

func IdIThrottle(_ IThrottle, label string, instance int) string {
	return IdInternalComponent(label, instance, "IThrottle")
}

type IThrottleSocket interface {
}

// IThrottle defines an interface for controlling the execution rate and tracking the count of operations executed.
// Throttle adjusts execution speed to meet a specified rate, ensuring operations run within defined limits.
// Counter returns the current count of throttled operations executed thus far.
type IThrottle interface {
	Setup() error

	Bind(socket IThrottleSocket, frameInterval int64) error

	Connect() error

	Counter() uint64

	Update()
}

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
