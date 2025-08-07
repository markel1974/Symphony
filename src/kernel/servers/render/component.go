package render

import (
	"strconv"

	"github.com/markel1974/c64emu/src/kernel/interfaces"
)

// Component represents a UI or system entity combining process details with associated rendering surfaces.
type Component struct {
	*interfaces.ProcessDescription
	descriptiveSurface *DescriptiveSurface
	surface            *Surface
}

// NewComponent initializes and returns a new Component instance using the provided process description and terminal interface.
func NewComponent(desc *interfaces.ProcessDescription, terminal interfaces.ITerminal, rows int, columns int) *Component {
	caption := strconv.Itoa(desc.PID())
	if len(desc.Name()) > 0 {
		caption += " - " + desc.Name()
	}
	c := &Component{
		ProcessDescription: desc,
		surface:            NewSurface(terminal, rows, columns, caption),
	}
	return c
}

// Surface returns the Surface instance associated with the Component.
func (c *Component) Surface() *Surface {
	return c.surface
}

// SetDescriptiveSurface sets the descriptive surface for the Component.
func (c *Component) SetDescriptiveSurface(surface interfaces.ISurface) {
	s, _ := surface.(*DescriptiveSurface)
	if s != nil {
		c.descriptiveSurface = s
	}
}
