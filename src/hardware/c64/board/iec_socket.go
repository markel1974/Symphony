package board

import (
	"fmt"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// IIECSocketConnection defines an interface for managing socket communication and triggering LED signals.
type IIECSocketConnection interface {
	LedTrigger(uint33 uint32)
}

// IECSocket represents a socket interface for managing communication and connections within an emulated environment.
type IECSocket struct {
	references.IIec
	connection IIECSocketConnection
	quartz     references.IComponent
}

// NewIECSocket creates and returns a new IECSocket instance using the provided IIECSocketConnection.
func NewIECSocket(connection IIECSocketConnection) *IECSocket {
	return &IECSocket{
		IIec:       nil,
		connection: connection,
	}
}

// Mount initializes the IECSocket instance, configuring components, quartz, and binding LED signals. Returns an error if any issue occurs.
func (s *IECSocket) Mount(cc map[string]references.IComponent, _ *config.Config, label string) error {
	var err error
	s.IIec, err = references.ComponentsToIEC(cc, label, 0)
	if err != nil {
		return err
	}
	var ok bool
	if s.quartz, ok = cc[references.IdIQuartz(nil, label, 0)]; !ok {
		return fmt.Errorf("nil quartz")
	}
	if err = s.IIec.Bind(s, s.quartz); err != nil {
		return err
	}
	s.IIec.LEDSignal().Bind(s.connection.LedTrigger)
	return nil
}
