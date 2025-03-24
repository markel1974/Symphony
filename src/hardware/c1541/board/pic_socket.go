package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// PICSocket is a type that embeds the IPIC6510 interface for interfacing with a programmable interrupt controller.
type PICSocket struct {
	references.IPIC6510
	quartz references.IQuartz
}

// NewPICSocket creates and initializes a new instance of PICSocket with a nil IPIC6510 interface.
func NewPICSocket() *PICSocket {
	return &PICSocket{
		IPIC6510: nil,
	}
}

// Setup sets up the programmable interrupt controller (PIC) with the provided quartz instance and assigns it to the socket.
// Returns an error if the setup process fails.
func (s *PICSocket) Setup(c map[string]references.IComponent, _ *config.Config) error {
	var err error
	s.IPIC6510, err = references.ComponentsToIPIC6510(c, 0)
	if err != nil {
		return err
	}
	s.quartz, err = references.ComponentsToIQuartz(c, 0)
	if err != nil {
		return err
	}
	if err = s.IPIC6510.Setup(s.quartz); err != nil {
		return err
	}
	return nil
}

func (s *PICSocket) Connect() error {
	return nil
}
