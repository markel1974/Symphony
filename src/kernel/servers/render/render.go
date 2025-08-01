package render

import (
	"github.com/markel1974/c64emu/src/kernel/adaptiveticker"
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/messages"
)

// eol represents the end-of-line marker used for denoting line breaks in the output, set to "\r\n".
const eol = "\r\n"

// Render represents a rendering engine responsible for managing terminal dimensions, repainting logic, and paint tasks.
type Render struct {
	terminal  interfaces.ITerminal
	dirty     bool
	width     int
	height    int
	fullPaint bool
	ticker    *adaptiveticker.AdaptiveTicker
	timerChan chan *adaptiveticker.TimerHandler
}

// NewRender creates and initializes a new Render instance with the provided terminal implementation.
// Returns a pointer to the newly created Render object.
func NewRender(ticker *adaptiveticker.AdaptiveTicker, timerChan chan *adaptiveticker.TimerHandler, terminal interfaces.ITerminal) *Render {
	return &Render{
		ticker:    ticker,
		timerChan: timerChan,
		terminal:  terminal,
		dirty:     false,
		width:     80,
		height:    24,
		fullPaint: true,
	}
}

func (c *Render) Read(data []byte) (int, error) {
	return c.terminal.Read(data)
}

// GetScreenSize returns the current screen width and height of the Render instance.
func (c *Render) GetScreenSize() (int, int) {
	return c.width, c.height
}

// SetScreenSize updates the screen's width and height, marks the screen for a full repaint, and sets the terminal size.
func (c *Render) SetScreenSize(width int, height int) {
	c.width = width
	c.height = height
	c.fullPaint = true
	c.terminal.SetSize(width, height)
}

// IsDirty checks if the render state is marked as dirty, indicating that a repaint is needed. It returns true if dirty.
func (c *Render) IsDirty() bool {
	return c.dirty
}

// ExecPaint performs the rendering process by painting background tasks and a foreground task onto the terminal surface.
func (c *Render) ExecPaint(fgTask interfaces.IProcess, tasks []interfaces.IProcess) bool {
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

// PaintRequest marks the rendering system as requiring a paint and optionally marks it for a full repaint.
// Returns true if the state was not already marked as dirty.
func (c *Render) PaintRequest(full bool) bool {
	if full {
		c.fullPaint = true
	}
	ret := false
	if !c.dirty {
		c.dirty = true
		c.ticker.Create(c.timerChan, messages.NewMessagePaint(), -1, -1, 1)
		ret = true
	}
	return ret
}

// Write sends the given string data to the terminal's output stream.
func (c *Render) Write(data string) {
	_, _ = c.terminal.Write(data)
}

// WriteLn writes the provided string to the terminal followed by an end-of-line character.
func (c *Render) WriteLn(data string) {
	c.Write(data)
	c.Write(eol)
}

// WriteColor writes the given data string to the terminal with specified foreground and background colors, and color mode.
func (c *Render) WriteColor(data string, fg interfaces.ColorDef, bg interfaces.ColorDef, mode interfaces.ColorMode) {
	_, _ = c.terminal.WriteColor(data, fg, bg, mode)
}

// WriteColorLn writes the given text with specified foreground and background colors and mode, followed by a line break.
func (c *Render) WriteColorLn(data string, fg interfaces.ColorDef, bg interfaces.ColorDef, mode interfaces.ColorMode) {
	c.WriteColor(data, fg, bg, mode)
	c.Write(eol)
}

// ClearScreen clears the terminal screen using the underlying ITerminal implementation.
func (c *Render) ClearScreen() {
	_, _ = c.terminal.ClearScreen()
}

// Scan processes the provided byte data using the underlying terminal's Scan method.
func (c *Render) Scan(data []byte) {
	c.terminal.Scan(data)
}

// ClearLine clears the specified line from the terminal screen using the terminal implementation of the associated Render object.
func (c *Render) ClearLine(line string) {
	_, _ = c.terminal.ClearLine(line)
}

// MoveCursorLeft moves the terminal cursor one position to the left using the underlying terminal implementation.
func (c *Render) MoveCursorLeft() {
	_, _ = c.terminal.MoveCursorLeft()
}

// MoveCursorRight moves the cursor one position to the right in the terminal.
func (c *Render) MoveCursorRight() {
	_, _ = c.terminal.MoveCursorRight()
}

// SaveCursor saves the current cursor position in the terminal for future restoration.
func (c *Render) SaveCursor() {
	_, _ = c.terminal.SaveCursor()
}

// RestoreCursor restores the saved cursor position in the terminal using the associated ITerminal implementation.
func (c *Render) RestoreCursor() {
	_, _ = c.terminal.RestoreCursor()
}

// EOL returns the end-of-line marker used by the terminal.
func (c *Render) EOL() string {
	return eol
}

func (c *Render) WritePromptLine(prompt string, line string) {
	c.ClearLine(line)
	c.WriteColor(prompt, interfaces.ColorGreenDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
	c.WriteColor(line, interfaces.ColorNoneDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
}

func (c *Render) WritePromptEOL(prompt string, eol bool) {
	if eol {
		c.WriteColor(c.EOL(), interfaces.ColorNoneDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
	}
	c.WriteColor(prompt, interfaces.ColorGreenDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
}

func (c *Render) WriteCritical(line string) {
	c.WriteColor(line, interfaces.ColorRedDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
}

func (c *Render) WriteNormal(line string) {
	c.WriteColor(line, interfaces.ColorNoneDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
}

func (c *Render) WriteHighlight(line string) {
	c.WriteColor(line, interfaces.ColorBlueDef, interfaces.ColorRedDef, interfaces.ModeNormal)
}
