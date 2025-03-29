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
}

// NewThrottleSocket initializes and returns a new instance of ThrottleSocket with the specified frame interval.
func NewThrottleSocket(parent references.IComponent, label string, frameInterval int64) *ThrottleSocket {
	return &ThrottleSocket{
		IThrottle:     nil,
		parent:        parent,
		label:         label,
		frameInterval: frameInterval,
	}
}

// Mount initializes the ThrottleSocket by configuring its IThrottle implementation using the given components and configuration.
func (s *ThrottleSocket) Mount() error {
	var err error
	idThrottle := references.IdIThrottle(s.IThrottle, s.label, 0)
	if s.IThrottle, err = references.ComponentToIThrottle(s.parent.GetChildByHardwareId(idThrottle)); err != nil {
		return err
	}
	if err = s.IThrottle.Bind(s, s.frameInterval); err != nil {
		return err
	}
	return nil
}
