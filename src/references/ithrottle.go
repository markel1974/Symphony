package references

import (
	"fmt"
	"github.com/markel1974/c64emu/src/config"
)

func IdIThrottle(_ IThrottle, instance int) string {
	return IdInternalComponent("IThrottle", instance)
}

type IThrottleSocket interface {
}

// IThrottle defines an interface for controlling the execution rate and tracking the count of operations executed.
// Throttle adjusts execution speed to meet a specified rate, ensuring operations run within defined limits.
// Counter returns the current count of throttled operations executed thus far.
type IThrottle interface {
	Setup(socket IThrottleSocket, cfg *config.Config, interval int64) error

	Connect() error

	Counter() uint64

	Update()
}

func ComponentToIThrottle(component IComponent) (IThrottle, error) {
	if component == nil {
		return nil, fmt.Errorf("component is nil")
	}
	v, ok := component.(IThrottle)
	if !ok {
		return nil, fmt.Errorf("component is not a %s", IdIThrottle(v, 0))
	}
	return v, nil
}

func ComponentsToIThrottle(cc map[string]IComponent, instance int) (IThrottle, error) {
	id := IdIThrottle(nil, instance)
	c, err := ComponentToIThrottle(cc[id])
	if err != nil {
		return nil, err
	}
	return c, nil
}
