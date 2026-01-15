package terminal

import (
	"github.com/markel1974/symphony/src/kernel/interfaces"
	"github.com/markel1974/symphony/src/kernel/terminal/vt100"
)

// EquipmentFactory is responsible for creating terminal equipment instances by leveraging provided input-output interfaces.
type EquipmentFactory struct {
}

// NewEquipmentFactory creates and returns a new instance of EquipmentFactory.
func NewEquipmentFactory() *EquipmentFactory {
	return &EquipmentFactory{}
}

// Create initializes a new VT100 terminal instance using the provided input/output and enter key.
func (f *EquipmentFactory) Create(_ string, enterKey rune) interfaces.ITerminal {
	return vt100.NewVt100(enterKey)
}
