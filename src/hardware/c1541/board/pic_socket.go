package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// PICSocket is a type that embeds the IPIC6510 interface for interfacing with a programmable interrupt controller.
type PICSocket struct {
	references.IPIC6510 // Incorpora l'interfaccia
}

// NewPICSocket creates and initializes a new instance of PICSocket with a nil IPIC6510 interface.
func NewPICSocket() *PICSocket {
	return &PICSocket{
		IPIC6510: nil,
	}
}

// Connect sets up the programmable interrupt controller (PIC) with the provided quartz instance and assigns it to the socket.
// Returns an error if the setup process fails.
func (s *PICSocket) Connect(pic references.IPIC6510, quart references.IQuartz) error {
	s.IPIC6510 = pic
	if err := s.IPIC6510.Setup(quart); err != nil {
		return err
	}
	return nil
}
