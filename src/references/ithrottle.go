package references

import "fmt"

func IdIThrottle(_ IThrottle) string {
	return "IThrottle"
}

// IThrottle defines an interface for controlling the execution rate and tracking the count of operations executed.
// Throttle adjusts execution speed to meet a specified rate, ensuring operations run within defined limits.
// Counter returns the current count of throttled operations executed thus far.
type IThrottle interface {
	Setup(int64) error

	Counter() uint64

	Throttle()
}

func ComponentToIThrottle(component IComponent, err error) (IThrottle, error) {
	if err = ComponentValidate(component, err); err != nil {
		return nil, err
	}
	v, ok := component.(IThrottle)
	if !ok {
		return nil, fmt.Errorf("component is not a %s", IdIThrottle(v))
	}
	return v, nil
}
