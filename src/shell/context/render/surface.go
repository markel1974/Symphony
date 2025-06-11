/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package render

import (
	"bytes"
	"github.com/markel1974/c64emu/src/shell/context/render/plotter"
	"github.com/markel1974/c64emu/src/shell/interfaces"
	"math"
)

// Surface represents a two-dimensional grid-based rendering surface for text-based terminal output.
type Surface struct {
	terminal  interfaces.ITerminal
	rows      int
	columns   int
	surface   [][]string
	rMax      int
	scale     float64
	offsetX   int
	offsetY   int
	border    int
	user      bool
	caption   string
	selection bool
	full      bool
	iRows     int
	iColumns  int
}

// newSurface initializes a new Surface object with the provided terminal, row, and column dimensions.
func newSurface(terminal interfaces.ITerminal, rows int, columns int) *Surface {
	s := &Surface{
		terminal: terminal,
		rows:     rows,
		columns:  columns,
		scale:    1.0,
		offsetX:  0,
		offsetY:  0,
		rMax:     0,
		border:   1,
		full:     false,
	}

	s.surface = make([][]string, s.rows)
	for r := range s.surface {
		s.surface[r] = make([]string, s.columns)
		for c := range s.surface[r] {
			s.surface[r][c] = " "
		}
	}
	return s
}

// Begin initializes the Surface for user-defined modifications and updates its internal dimensions.
func (s *Surface) Begin() {
	s.user = true
	s.iRows, s.iColumns = s.GetSize()
}

// End terminates the user interaction mode and re-renders the window.
func (s *Surface) End() {
	s.user = false
	s.drawWindow()
}

// GetSize computes the dimensions of the surface, adjusting for scale and borders if necessary, and returns rows and columns.
func (s *Surface) GetSize() (int, int) {
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

// SetCompletePaint enables full rendering mode, marking the entire surface as requiring a complete repaint.
func (s *Surface) SetCompletePaint() {
	s.full = true
}

// SetSelectionMode sets the selection state of the Surface. When true, enables selection mode; when false, disables it.
func (s *Surface) SetSelectionMode(selection bool) {
	s.selection = selection
}

// SetCaption sets the caption text for the Surface. It updates the internal caption field with the provided string.
func (s *Surface) SetCaption(caption string) {
	s.caption = caption
}

// SetOffsetX updates the horizontal offset (offsetX) for the Surface, affecting the horizontal position of drawn elements.
func (s *Surface) SetOffsetX(offsetX int) {
	s.offsetX = offsetX
}

// SetOffsetY sets the vertical offset of the surface to the specified value.
func (s *Surface) SetOffsetY(offsetY int) {
	s.offsetY = offsetY
}

// SetScale updates the scaling factor for the Surface, affecting the computed size during rendering and operations.
func (s *Surface) SetScale(scale float64) {
	s.scale = scale
}

// Draw places a rune at the specified row and column on the surface, considering offsets and boundaries.
func (s *Surface) Draw(rs int, cs int, text rune) {
	rows, columns := s.compute(rs, cs)
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
			if rows > s.rMax {
				s.rMax = rows
			}
		}
	}
}

// DrawColor draws a colored character at the specified row and column with given foreground, background, and color mode.
func (s *Surface) DrawColor(rs int, cs int, text rune, fg interfaces.ColorDef, bg interfaces.ColorDef, mode interfaces.ColorMode) {
	rows, columns := s.compute(rs, cs)
	if rows < 0 {
		return
	}
	if columns < 0 {
		return
	}

	ars, acs := s.GetSize()
	if rs >= ars {
		return
	}
	if cs >= acs {
		return
	}

	if len(s.surface) > rows {
		if len(s.surface[rows]) > columns {
			colorized := s.terminal.Colorize(string(text), int(fg), int(bg), mode)
			s.surface[rows][columns] = colorized
			if rows > s.rMax {
				s.rMax = rows
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
	rows, columns := s.GetSize()
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
func (s *Surface) compute(r int, c int) (int, int) {
	rows := r + s.offsetY
	column := c + s.offsetX
	if s.user {
		rows += s.border
		column += s.border
	}
	return rows, column
}

// GetBuffer generates a byte slice representing the surface content, limited by the render boundary.
func (s *Surface) GetBuffer() []byte {
	var lines bytes.Buffer
	var max int

	if s.full {
		max = s.rows * s.columns
	} else {
		max = (s.rMax + 1) * s.columns
	}

	var counter = 0
	var halt = false

	for h, horizontal := range s.surface {
		if h != 0 {
			lines.WriteString("\r\n")
		}
		for _, v := range horizontal {
			if counter < max {
				lines.WriteString(v)
				counter++
			} else {
				halt = true
				break
			}
		}
		if halt {
			break
		}
	}

	return lines.Bytes()
}

// Render updates the terminal by rendering the current buffer to the screen and restoring the cursor position.
func (s *Surface) Render() {
	var buffer = string(s.GetBuffer())
	_, _ = s.terminal.SaveCursor()
	_, _ = s.terminal.MoveCursorTopLeft()
	_, _ = s.terminal.Write(buffer)
	_, _ = s.terminal.RestoreCursor()
}

// drawWindow draws a bordered window on the surface with optional caption and selection mode colors.
func (s *Surface) drawWindow() {
	rows, columns := s.GetSize()
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
