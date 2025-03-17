package board

import (
	"github.com/markel1974/c64emu/src/references"
)

type ThrottleSocket struct {
	references.IThrottle
}

func NewThrottleSocket() *ThrottleSocket {
	return &ThrottleSocket{
		IThrottle: nil,
	}
}

func (s *ThrottleSocket) Connect(throttle references.IThrottle, frameInterval int64) error {
	s.IThrottle = throttle
	if err := s.IThrottle.Setup(frameInterval); err != nil {
		return err
	}
	return nil
}
