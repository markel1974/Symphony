package context

import (
	"github.com/markel1974/c64emu/src/shell/interfaces"
)

type Render struct {
	terminal  interfaces.ITerminal
	dirty     bool
	width     int
	height    int
	fullPaint bool
}

func NewRender(terminal interfaces.ITerminal) *Render {
	return &Render{
		terminal:  terminal,
		dirty:     false,
		width:     80,
		height:    24,
		fullPaint: true,
	}
}

// GetScreenSize retrieves the current screen width and height as two integer values.
func (c *Render) GetScreenSize() (int, int) {
	return c.width, c.height
}

func (c *Render) SetScreenSize(width int, height int) {
	c.width = width
	c.height = height
	c.fullPaint = true
	c.terminal.SetSize(width, height)
}

func (c *Render) IsDirty() bool {
	return c.dirty
}

// ExecPaint renders the current task manager state onto the terminal by painting tasks and handling selection logic.
func (c *Render) ExecPaint(fgTask interfaces.ITask, tasks []interfaces.ITask) bool {
	w, h := c.GetScreenSize()
	surface := newSurface(c.terminal, h, w)
	if c.fullPaint {
		surface.SetCompletePaint()
		c.fullPaint = false
	}
	//zOrder
	for _, task := range tasks {
		surface.SetSelectionMode(false)
		task.Paint(surface)
	}
	//zOrder
	if fgTask != nil {
		surface.SetSelectionMode(true)
		fgTask.Paint(surface)
	}
	surface.Render()
	c.dirty = false
	return true
}

// PaintRequest marks the Kernel as dirty and initiates a paint request if not already pending.
func (c *Render) PaintRequest(full bool) bool {
	if full {
		c.fullPaint = true
	}
	ret := false
	if !c.dirty {
		c.dirty = true
		ret = true
	}
	return ret
}

func (c *Render) Write(data string) {
	_, _ = c.terminal.Write(data)
}

func (c *Render) WriteLn(data string) {
	c.Write(data)
	c.Write(eol)
}

func (c *Render) WriteColor(data string, fg interfaces.ColorDef, bg interfaces.ColorDef, mode interfaces.ColorMode) {
	_, _ = c.terminal.WriteColor(data, fg, bg, mode)
}

func (c *Render) WriteColorLn(data string, fg interfaces.ColorDef, bg interfaces.ColorDef, mode interfaces.ColorMode) {
	c.WriteColor(data, fg, bg, mode)
	c.Write(eol)
}

func (c *Render) ClearScreen() {
	_, _ = c.terminal.ClearScreen()
}

func (c *Render) Scan(data []byte) {
	c.terminal.Scan(data)
}
