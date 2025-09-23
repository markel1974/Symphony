package xvi

import (
	"fmt"
	"strings"

	"github.com/markel1974/c64emu/src/kernel/interfaces"
)

// Tui represents a text-based user interface for interacting with a Buffer and managing viewport rendering and modes.
type Tui struct {
	buffer        *Buffer
	mode          string
	status        string
	rows          int
	columns       int
	offsetY       int
	offsetX       int
	commandBuffer string
}

// NewTui creates and initializes a new Tui instance with the provided Buffer and default settings.
func NewTui(buffer *Buffer) *Tui {
	t := &Tui{
		buffer:  buffer,
		mode:    "normal",
		status:  "",
		offsetY: 0,
		offsetX: 0,
	}
	return t
}

// UpdateCommand updates the command buffer of the Tui instance with the provided command string.
func (t *Tui) UpdateCommand(cmd string) {
	t.commandBuffer = cmd
}

// SetMode sets the mode of the Tui instance to the specified string value.
func (t *Tui) SetMode(mode string) {
	t.mode = mode
}

// GetMode returns the current mode of the Tui instance as a string.
func (t *Tui) GetMode() string {
	return t.mode
}

// SetError sets the status of the Tui instance to the specified error value.
func (t *Tui) SetError(err error) {
	if err == nil {
		t.status = ""
	} else {
		t.status = err.Error()
	}
}

// Draw renders the content of the buffer onto the surface based on the viewport and cursor position.
func (t *Tui) Draw(process interfaces.IUserProcess, surface interfaces.ISurface) {
	rows, columns := surface.GetWindowSize()
	if rows != t.rows || columns != t.columns {
		t.rows = rows
		t.columns = columns
		t.offsetY = 0
		t.offsetX = 0
		surface.MoveCursor(0, 0)
		process.PaintRequest()
		return
	}
	cx, cy := t.buffer.Cursor()
	textAreaHeight := t.rows - 4
	textAreaWidth := t.columns
	if cy < t.offsetY {
		t.offsetY = cy
	}
	if cy >= t.offsetY+textAreaHeight {
		t.offsetY = cy - textAreaHeight + 1
	}
	if cx < t.offsetX {
		t.offsetX = cx
	}
	if cx >= t.offsetX+textAreaWidth {
		t.offsetX = cx - textAreaWidth + 1
	}
	for y := 0; y < textAreaHeight; y++ {
		bufferLineIndex := y + t.offsetY
		if bufferLineIndex < t.buffer.LineCount() {
			line := t.buffer.GetLine(bufferLineIndex)
			if t.offsetX < len(line) {
				visibleLine := line[t.offsetX:]
				if len(visibleLine) > textAreaWidth {
					visibleLine = visibleLine[:textAreaWidth]
				}
				surface.DrawTextColor(y, 0, visibleLine, interfaces.ColorWhiteDef, interfaces.ColorBlackDef, 0)
			}
		} else {
			surface.DrawTextColor(y, 0, "~", interfaces.ColorGreenDef, interfaces.ColorBlackDef, 0)
		}
	}

	// Status Bar
	var statusText string
	switch t.mode {
	case "command":
		statusText = ":"
	default:
		statusText = fmt.Sprintf("-- %s -- %s   %d, %d", strings.ToUpper(t.mode), t.buffer.filePath, cy+1, cx+1)
	}

	statusLine := t.status + strings.Repeat(" ", t.columns-len(t.status))
	surface.DrawTextColor(t.rows-3, 0, statusLine, interfaces.ColorBlackDef, interfaces.ColorWhiteDef, 0)
	surface.DrawTextColor(t.rows-4, 0, statusText, interfaces.ColorBlackDef, interfaces.ColorWhiteDef, 0)

	// Sposta il cursore fisico nella posizione corretta RELATIVA alla viewport
	surface.MoveCursor(cy-t.offsetY, cx-t.offsetX)

	/*
		// Status Bar
		var statusText string
		switch t.mode {
		case "command":
			statusText = ":" + t.commandBuffer
		default:
			cx, cy := t.buffer.Cursor()
			statusText = fmt.Sprintf("-- %s -- %s   %d, %d", strings.ToUpper(t.mode), t.buffer.filePath, cy+1, cx+1)
		}
		statusLine := t.status + strings.Repeat(" ", t.columns-len(t.status))
		surface.DrawTextColor(t.rows-2, 0, statusLine, interfaces.ColorBlackDef, interfaces.ColorWhiteDef, 0)
		commandLine := statusText + strings.Repeat(" ", t.columns-len(statusText))
		surface.DrawTextColor(t.rows-1, 0, commandLine, interfaces.ColorBlackDef, interfaces.ColorWhiteDef, 0)
		if t.mode == "command" {
			surface.MoveCursor(t.rows-1, len(statusText))
		} else {
			cx, cy := t.buffer.Cursor()
			surface.MoveCursor(cy-t.offsetY, cx-t.offsetX)
		}
	*/
}
