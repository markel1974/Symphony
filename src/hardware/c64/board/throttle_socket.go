package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

type ThrottleSocket struct {
	references.IThrottle
	frameInterval int64
}

func NewThrottleSocket(frameInterval int64) *ThrottleSocket {
	return &ThrottleSocket{
		IThrottle:     nil,
		frameInterval: frameInterval,
	}
}

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

func (s *ThrottleSocket) Connect() error {
	if err := s.IThrottle.Connect(); err != nil {
		return err
	}
	return nil
}
