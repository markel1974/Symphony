package render

import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"strconv"
)

type Component struct {
	*interfaces.ProcessDescription
	*WindowOptions
}

func NewComponent(desc *interfaces.ProcessDescription) *Component {
	caption := strconv.Itoa(desc.PID())
	if len(desc.Name()) > 0 {
		caption += " - " + desc.Name()
	}
	return &Component{
		ProcessDescription: desc,
		WindowOptions:      NewWindowOptions(caption, 0, 0, 1.0),
	}
}
