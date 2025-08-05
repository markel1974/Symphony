package render

import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"strconv"
)

type Component struct {
	*interfaces.ProcessDescription
	*WindowOptions
	surface *Surface
}

func NewComponent(desc *interfaces.ProcessDescription, terminal interfaces.ITerminal) *Component {
	caption := strconv.Itoa(desc.PID())
	if len(desc.Name()) > 0 {
		caption += " - " + desc.Name()
	}
	return &Component{
		ProcessDescription: desc,
		WindowOptions:      NewWindowOptions(caption, 0, 0, 1.0),
		surface:            NewSurface(terminal, 24, 80),
	}
}
