package references

import (
	"fmt"
	"github.com/markel1974/c64emu/src/config"
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
	Setup(cc map[string]IComponent, cfg *config.Config) error

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

func ComponentsToIThrottle(cc map[string]IComponent, label string, instance int) (IThrottle, error) {
	id := IdIThrottle(nil, label, instance)
	c, err := ComponentToIThrottle(cc[id])
	if err != nil {
		return nil, err
	}
	return c, nil
}
