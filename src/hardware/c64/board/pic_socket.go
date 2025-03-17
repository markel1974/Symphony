package board

import (
	"github.com/markel1974/c64emu/src/references"
)

type PICSocket struct {
	references.IPIC6510 // Incorpora l'interfaccia
}

func NewPICSocket() *PICSocket {
	return &PICSocket{
		IPIC6510: nil,
	}
}

func (s *PICSocket) Connect(pic references.IPIC6510, quart references.IQuartz) error {
	s.IPIC6510 = pic
	if err := s.IPIC6510.Setup(quart); err != nil {
		return err
	}
	return nil
}
