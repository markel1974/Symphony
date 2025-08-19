package xvi

import (
	"fmt"
	"strings"

	"github.com/markel1974/c64emu/src/kernel/interfaces"
)

type Tui struct {
	buffer  *Buffer
	mode    string
	rows    int
	columns int
}

func NewTui(buffer *Buffer) *Tui {
	t := &Tui{
		buffer: buffer,
		mode:   "normal",
	}
	return t
}

func (t *Tui) SetMode(mode string) {
	t.mode = mode
}

func (t *Tui) GetMode() string {
	return t.mode
}

func (t *Tui) Draw(process interfaces.IUserProcess, surface interfaces.ISurface) {
	rows, columns := surface.GetSize()
	fmt.Println(rows, columns)
	if rows != t.rows || columns != t.columns {
		t.rows = rows
		t.columns = columns
		process.PaintRequest()
		return
	}
	for y := 0; y < t.buffer.LineCount(); y++ {
		line := t.buffer.GetLine(y)
		for idx, l := range line {
			surface.DrawColor(y, idx, l, interfaces.ColorWhiteDef, interfaces.ColorBlackDef, 0)
		}
	}
	cx, cy := t.buffer.Cursor()
	statusText := fmt.Sprintf("-- %s -- %s   %d, %d", strings.ToUpper(t.mode), t.buffer.filePath, cy+1, cx+1)
	for idx, l := range statusText {
		surface.DrawColor(t.rows-3, idx, l, interfaces.ColorBlackDef, interfaces.ColorWhiteDef, 0)
	}
	surface.MoveCursor(cy, cx)
}
