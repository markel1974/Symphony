package render

import (
	"strconv"

	"github.com/markel1974/c64emu/src/kernel/interfaces"
)

// Component represents a UI or system entity combining process details with associated rendering surfaces.
type Component struct {
	*interfaces.ProcessDescription
	surface *Surface
}

// NewComponent initializes and returns a new Component instance using the provided process description and terminal interface.
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

// Surface returns the Surface instance associated with the Component.
func (c *Component) Surface() *Surface {
	return c.surface
}

// CloneSurface returns a clone of the Surface instance associated with the Component.
func (c *Component) CloneSurface() *Surface {
	s := c.surface.Clone()
	return s
}

// SetSurface assigns a provided ISurface implementation to the Component if it can be cast to the specific Surface type.
func (c *Component) SetSurface(surface interfaces.ISurface) {
	s, _ := surface.(*Surface)
	if s != nil {
		c.surface = s
	}
}
