package render

import (
	"bytes"
	"math"

	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/servers/render/plotter"
)

// Surface represents a two-dimensional grid-based rendering surface for text-based terminal output.
type Surface struct {
	interpreted    *InterpretedSurface
	terminal       interfaces.ITerminal
	surface        [][]string
	columnsCleaner []string
	caption        string
	rows           int
	columns        int
	rowMax         int
	scale          float64
	offsetX        int
	offsetY        int
	border         int
	user           bool
	selection      bool
	iRows          int
	iColumns       int
	cursorRow      int
	cursorColumn   int
}

// NewSurface initializes a new Surface object with the provided terminal, row, and column dimensions.
func NewSurface(terminal interfaces.ITerminal, rows int, columns int, caption string) *Surface {
	s := &Surface{
		terminal:     terminal,
		caption:      caption,
		rows:         -1,
		columns:      -1,
		scale:        1.0,
		offsetX:      0,
		offsetY:      0,
		rowMax:       0,
		border:       1,
		interpreted:  nil,
		cursorRow:    -1,
		cursorColumn: -1,
	}
	s.Prepare(rows, columns)
	return s
}

// GetWindowSize returns the calculated width and height of the window based on scale and user-defined adjustments.
func (s *Surface) GetWindowSize() (int, int) {
	rows := s.rows
	columns := s.columns
	if s.scale > 0 && s.scale < 1 {
		rows = int(math.Round(float64(rows) * s.scale))
		columns = int(math.Round(float64(columns) * s.scale))
	}
	if s.user {
		rows -= s.border * 2
		columns -= s.border * 2
	}
	return rows, columns
}

// GetScreenSize returns the number of rows and columns representing the screen size of the surface.
func (s *Surface) GetScreenSize() (int, int) {
	return s.rows, s.columns
}

// Assign copies the properties of the given Surface to the current one.
func (s *Surface) Assign(surface *Surface) {
	s.caption = surface.caption
	s.scale = surface.scale
	s.offsetX = surface.offsetX
	s.offsetY = surface.offsetY
	s.border = surface.border
	s.selection = surface.selection
}

// SetInterpretedSurface sets the descriptive surface for the Surface.
func (s *Surface) SetInterpretedSurface(surface *InterpretedSurface) {
	s.interpreted = surface
}

// GetInterpretedSurface retrieves the descriptive surface for the Surface.
func (s *Surface) GetInterpretedSurface() *InterpretedSurface {
	return s.interpreted
}

// SetOption updates the task's X, Y offsets or Scale based on the given option ('x', 'y', or 'z') and value.
func (s *Surface) SetOption(option rune, value float64) {
	switch option {
	case 'y':
		s.offsetY = s.offsetY + int(value)
	case 'x':
		s.offsetX = s.offsetX + int(value)
	case 'z':
		if scale := s.scale + value; scale >= 0.2 && scale <= 1 {
			s.scale = scale
		}
	}
}

// SetRowMax sets the maximum allowed rows for the Surface instance to the specified value.
func (s *Surface) SetRowMax(rowMax int) {
	s.rowMax = rowMax
}

// RowMax retrieves the maximum row index of the Surface.
func (s *Surface) RowMax() int {
	return s.rowMax
}

// Prepare adjusts the dimensions of the Surface to the specified rows and columns, resetting or reallocating data as needed.
func (s *Surface) Prepare(rows int, columns int) {
	s.rowMax = 0
	if rows == s.rows && columns == s.columns {
		for r := range s.surface {
			copy(s.surface[r], s.columnsCleaner)
		}
		return
	}
	empty := s.terminal.CreateEmpty()
	s.columnsCleaner = make([]string, columns)
	for c := range s.columnsCleaner {
		s.columnsCleaner[c] = empty
	}
	s.surface = make([][]string, rows)
	for r := range s.surface {
		s.surface[r] = make([]string, columns)
		copy(s.surface[r], s.columnsCleaner)
	}
	s.rows = rows
	s.columns = columns
	return
}

// Merge merges the given Surface with the current one, overriding non-empty cells in the target Surface.
func (s *Surface) Merge(surface *Surface) {
	empty := s.terminal.CreateEmpty()
	for i, row := range surface.surface {
		if i >= len(s.surface) {
			continue
		}
		for j, val := range row {
			if j >= len(s.surface[i]) {
				continue
			}
			if val != empty {
				s.surface[i][j] = val
			}
		}
	}
}

// Begin initializes the Surface for user-defined modifications and updates its internal dimensions.
func (s *Surface) Begin() {
	s.user = true
	s.iRows, s.iColumns = s.GetWindowSize()
	s.cursorRow = -1
	s.cursorColumn = -1
}

// End terminates the user interaction mode and re-renders the window.
func (s *Surface) End() {
	s.user = false
	s.drawWindow()
}

// SetSelectionMode sets the selection mode for the Surface. If true, the surface will render in selection mode.
func (s *Surface) SetSelectionMode(selection bool) {
	s.selection = selection
}

// MoveCursor moves the cursor to the specified row and column.
func (s *Surface) MoveCursor(rs int, cs int) {
	//if rs < 1 {
	//	rs = 1
	//}
	//if cs < 1 {
	//	cs = 1
	//}
	rows, columns := s.compute(rs, cs, true)
	rows++
	columns++
	//if rows < s.border {
	//	rows = s.border
	//}
	//if columns < s.border {
	//	columns = s.border
	//}
	maxRows := s.iRows + s.offsetY - s.border
	maxColumns := s.iColumns + s.offsetX - s.border
	if rows >= maxRows {
		rows = maxRows
	}
	if columns >= maxColumns {
		columns = maxColumns
	}
	s.cursorRow = rows
	s.cursorColumn = columns
	//fmt.Printf("Cursor: %d, %d, Received: %d, %d | %d | %d\n", s.cursorRow, s.cursorColumn, rs, cs, maxRows, maxColumns)
}

// Cursor returns the current position of the cursor as a row and column.
func (s *Surface) Cursor() (int, int) {
	return s.cursorRow, s.cursorColumn
}

// Draw places a rune at the specified row and column on the surface, considering offsets and boundaries.
func (s *Surface) Draw(rs int, cs int, text rune) {
	rows, columns := s.compute(rs, cs, s.user)
	if rows < 0 {
		return
	}
	if columns < 0 {
		return
	}
	if rs >= s.iRows {
		return
	}
	if cs >= s.iColumns {
		return
	}
	if len(s.surface) > rows {
		if len(s.surface[rows]) > columns {
			s.surface[rows][columns] = string(text)
			if rows > s.rowMax {
				s.rowMax = rows
			}
		}
	}
}

// DrawColor renders a colored character at the specified row and column on the surface using the given foreground color, background color, and color mode.
func (s *Surface) DrawColor(rs int, cs int, text rune, fg interfaces.ColorDef, bg interfaces.ColorDef, mode interfaces.ColorMode) {
	rows, columns := s.compute(rs, cs, s.user)
	if rows < 0 {
		return
	}
	if columns < 0 {
		return
	}
	ars, acs := s.GetWindowSize()
	if rs >= ars {
		return
	}
	if cs >= acs {
		return
	}

	if len(s.surface) > rows {
		if len(s.surface[rows]) > columns {
			colorized := s.terminal.CreateColorize(string(text), int(fg), int(bg), mode)
			s.surface[rows][columns] = colorized
			if rows > s.rowMax {
				s.rowMax = rows
			}
		}
	}
}

// DrawText writes a string `text` to the surface, starting at the specified `rows` and `column`.
func (s *Surface) DrawText(rows int, column int, text string) {
	for x, d := range text {
		s.Draw(rows, column+x, d)
	}
}

// DrawTextColor renders a colored string at the specified row and column positions on the surface.
// Each character in the string is drawn with the provided foreground color, background color, and color mode.
// The method uses internal character positioning and adjusts text placement based on the surface offset and scale.
func (s *Surface) DrawTextColor(rows int, column int, text string, fg interfaces.ColorDef, bg interfaces.ColorDef, mode interfaces.ColorMode) {
	for x, d := range text {
		s.DrawColor(rows, column+x, d, fg, bg, mode)
	}
}

// DrawSeries renders a data series onto the surface using specified width, height, and value range (min and max).
func (s *Surface) DrawSeries(data []float64, w int, h int, min float64, max float64) {
	rows, columns := s.GetWindowSize()
	if h <= 0 {
		h = rows
	}
	if w <= 0 {
		w = columns
	}
	if h >= rows {
		h = rows - 1
	}

	g := plotter.NewPlotter(w, h)
	g.Setup(data, min, max)
	g.Draw(s)
}

// compute adjusts the given row and column by applying surface offsets and border if in user mode. Returns new indices.
func (s *Surface) compute(r int, c int, user bool) (int, int) {
	rows := r + s.offsetY
	column := c + s.offsetX
	if user {
		rows += s.border
		column += s.border
	}
	return rows, column
}

// GetBuffer generates a byte slice representing the surface content, limited by the render boundary.
func (s *Surface) GetBuffer(lines *bytes.Buffer, full bool) {
	var maximum int
	if full {
		maximum = s.rows * s.columns
	} else {
		maximum = (s.rowMax + 1) * s.columns
	}
	var counter = 0
	for h, horizontal := range s.surface {
		if h != 0 {
			lines.WriteString("\r\n")
		}
		for _, v := range horizontal {
			if counter < maximum {
				lines.WriteString(v)
				counter++
			} else {
				return
			}
		}
	}
	//return lines.Bytes()
}

// drawWindow draws a bordered window on the surface with optional caption and selection mode colors.
func (s *Surface) drawWindow() {
	rows, columns := s.GetWindowSize()
	fg := interfaces.ColorWhiteDef
	bg := interfaces.ColorNoneDef
	mode := interfaces.ModeNormal

	if s.selection {
		fg = interfaces.ColorRedDef
	}

	for y := 0; y < rows; y++ {
		s.DrawColor(y, 0, '│', fg, bg, mode)
		s.DrawColor(y, columns-1, '│', fg, bg, mode)
	}

	for x := 0; x < columns; x++ {
		s.DrawColor(0, x, '─', fg, bg, mode)
		s.DrawColor(rows-1, x, '─', fg, bg, mode)
	}

	s.DrawColor(0, 0, '╭', fg, bg, mode)
	s.DrawColor(0, columns-1, '╮', fg, bg, mode)

	s.DrawColor(rows-1, 0, '╰', fg, bg, mode)
	s.DrawColor(rows-1, columns-1, '╯', fg, bg, mode)

	s.DrawTextColor(0, 2, s.caption, fg, bg, mode)
}
