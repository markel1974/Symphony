package board

import (
	"github.com/markel1974/symphony/src/references"
)

// IIECSocketConnection is an interface for managing LED activity in IEC socket connections.
// LedActivityTrigger toggles the LED state for a given device number with a boolean value.
type IIECSocketConnection interface {
	LedActivityTrigger(deviceNumber uint8, led bool)
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

// Wire initializes the IECSocket instance, configuring components, quartz, and binding LED signals. Returns an error if any issue occurs.
func (s *IECSocket) Wire() error {
	var err error
	s.component = s.parent.GetChildByHardwareId(s.HardwareId())
	if s.IIec, err = references.ComponentToIEC(s.component); err != nil {
		return err
	}
	if err = s.IIec.Bind(s); err != nil {
		return err
	}
	return nil
}

func (s *IECSocket) LedActivity(deviceNumber uint8, led bool) {
	s.connection.LedActivityTrigger(deviceNumber, led)
}
