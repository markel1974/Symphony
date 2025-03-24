package board

import (
	"fmt"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

type IIECSocketConnection interface {
	LedTrigger(uint33 uint32)
}

type IECSocket struct {
	references.IIec
	connection IIECSocketConnection
	quartzC    references.IComponent
}

func NewIECSocket(connection IIECSocketConnection) *IECSocket {
	return &IECSocket{
		IIec:       nil,
		connection: connection,
	}
}

func (s *IECSocket) Setup(cc map[string]references.IComponent, cfg *config.Config) error {
	var err error
	s.IIec, err = references.ComponentsToIEC(cc, 0)
	if err != nil {
		return err
	}
	var ok bool
	if s.quartzC, ok = cc[references.IdIQuartz(nil, 0)]; !ok {
		return fmt.Errorf("nil quartz")
	}
	if err = s.IIec.Setup(s, cfg); err != nil {
		return err
	}
	s.IIec.LEDSignal().Bind(s.connection.LedTrigger)
	return nil
}

func (s *IECSocket) Connect() error {
	if err := s.IIec.Connect(s.quartzC); err != nil {
		return err
	}
	return nil
}
