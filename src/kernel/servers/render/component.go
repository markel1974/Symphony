package render

import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"strconv"
)

type Component struct {
	*interfaces.ProcessDescription
	surface *Surface
}

func NewComponent(desc *interfaces.ProcessDescription, terminal interfaces.ITerminal) *Component {
	caption := strconv.Itoa(desc.PID())
	if len(desc.Name()) > 0 {
		caption += " - " + desc.Name()
	}
	c := &Component{
		ProcessDescription: desc,
		surface:            NewSurface(terminal, 24, 80, caption),
	}
	return c
}
