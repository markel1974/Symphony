package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// IIECSocketConnection defines an interface for managing socket communication and triggering LED signals.
type IIECSocketConnection interface {
	LedTrigger(uint33 uint32)
}

// IECSocket represents a socket interface for managing communication and connections within an emulated environment.
type IECSocket struct {
	label     string
	parent    references.IComponent
	component references.IComponent
	references.IIec
	connection IIECSocketConnection
	hwId       string
}

// NewIECSocket creates and returns a new IECSocket instance using the provided IIECSocketConnection.
func NewIECSocket(parent references.IComponent, label string, connection IIECSocketConnection) *IECSocket {
	s := &IECSocket{
		IIec:       nil,
		parent:     parent,
		label:      label,
		connection: connection,
	}
	s.hwId = references.IdIIec(s.IIec, s.label, 0)
	return s
}

func (s *IECSocket) HardwareId() string {
	return s.hwId
}

// Mount initializes the IECSocket instance, configuring components, quartz, and binding LED signals. Returns an error if any issue occurs.
func (s *IECSocket) Mount() error {
	var err error
	s.component = s.parent.GetChildByHardwareId(s.HardwareId())
	if s.IIec, err = references.ComponentToIEC(s.component); err != nil {
		return err
	}
	if err = s.IIec.Bind(s); err != nil {
		return err
	}
	s.IIec.LEDSignal().Bind(s.connection.LedTrigger)
	return nil
}
