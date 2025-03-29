package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// ThrottleSocket wraps an IThrottle implementation for managing execution rate with a specific frame interval setting.
type ThrottleSocket struct {
	references.IThrottle
	label         string
	parent        references.IComponent
	component     references.IComponent
	frameInterval int64
	hwId          string
}

// NewThrottleSocket initializes and returns a new instance of ThrottleSocket with the specified frame interval.
func NewThrottleSocket(parent references.IComponent, label string, frameInterval int64) *ThrottleSocket {
	s := &ThrottleSocket{
		IThrottle:     nil,
		parent:        parent,
		label:         label,
		frameInterval: frameInterval,
	}
	s.hwId = references.IdIThrottle(s.IThrottle, s.label, 0)
	return s
}

func (s *ThrottleSocket) HardwareId() string {
	return s.hwId
}

// Mount initializes the ThrottleSocket by configuring its IThrottle implementation using the given components and configuration.
func (s *ThrottleSocket) Mount() error {
	var err error
	s.component = s.parent.GetChildByHardwareId(s.HardwareId())
	if s.IThrottle, err = references.ComponentToIThrottle(s.component); err != nil {
		return err
	}
	if err = s.IThrottle.Bind(s, s.frameInterval); err != nil {
		return err
	}
	return nil
}
