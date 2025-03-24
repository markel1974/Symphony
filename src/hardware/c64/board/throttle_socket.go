package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// ThrottleSocket wraps an IThrottle implementation for managing execution rate with a specific frame interval setting.
type ThrottleSocket struct {
	references.IThrottle
	frameInterval int64
}

// NewThrottleSocket initializes and returns a new instance of ThrottleSocket with the specified frame interval.
func NewThrottleSocket(frameInterval int64) *ThrottleSocket {
	return &ThrottleSocket{
		IThrottle:     nil,
		frameInterval: frameInterval,
	}
}

// Setup initializes the ThrottleSocket by configuring its IThrottle implementation using the given components and configuration.
func (s *ThrottleSocket) Setup(c map[string]references.IComponent, cfg *config.Config) error {
	var err error
	if s.IThrottle, err = references.ComponentsToIThrottle(c, 0); err != nil {
		return err
	}
	if err = s.IThrottle.Setup(s, cfg, s.frameInterval); err != nil {
		return err
	}
	return nil
}

// Connect establishes a connection using the underlying IThrottle interface. Returns an error if the connection fails.
func (s *ThrottleSocket) Connect() error {
	if err := s.IThrottle.Connect(); err != nil {
		return err
	}
	return nil
}
