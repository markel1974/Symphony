package render

import (
	"strconv"

	"github.com/markel1974/c64emu/src/kernel/interfaces"
)

// Component represents a UI or system entity combining process details with associated rendering surfaces.
type Component struct {
	surface   *Surface
	pid       int
	available bool
	zIndex    int
}

// NewComponent initializes and returns a new Component instance using the provided process description and terminal interface.
func NewComponent(pid int, description string, terminal interfaces.ITerminal, rows int, columns int) *Component {
	caption := strconv.Itoa(pid)
	if len(description) > 0 {
		caption += " - " + description
	}
	c := &Component{
		pid:       pid,
		surface:   NewSurface(terminal, rows, columns, caption),
		available: false,
		zIndex:    0,
	}
	return c
}

func (c *Component) ZIndex() int {
	return c.zIndex
}

func (c *Component) SetZIndex(zIndex int) {
	c.zIndex = zIndex
}

func (c *Component) Available() bool {
	return c.available
}

func (c *Component) SetAvailable() {
	c.available = true
}

// PID returns the process ID associated with the Component.
func (c *Component) PID() int {
	return c.pid
}

// Surface returns the Surface instance associated with the Component.
func (c *Component) Surface() *Surface {
	return c.surface
}

// RowMax returns the maximum row index for the Component's surface.'
func (c *Component) RowMax() int {
	return c.surface.RowMax()
}

/*
// Compile renders the Component's descriptive surface and applies it to the Component's surface.
func (c *Component) Compile(height int, width int) {
	c.surface.Prepare(height, width)
	//c.surface.SetZIndex(0)
	c.surface.SetSelectionMode(false)
	if c.surface.interpreted != nil {
		//if c.PID() == activePid {
		//	c.surface.SetZIndex(255)
		//	c.surface.SetSelectionMode(true)
		//}
		c.surface.Begin()
		c.surface.interpreted.Appy(c.surface)
		c.surface.End()
	}
}

*/

// SetInterpretedSurface sets the descriptive surface for the Component.
func (c *Component) SetInterpretedSurface(surface interfaces.ISurface) {
	s, _ := surface.(*InterpretedSurface)
	if s != nil {
		c.surface.SetInterpretedSurface(s)
	}
}
