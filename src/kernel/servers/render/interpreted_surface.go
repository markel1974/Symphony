package render

import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/servers/render/plotter"
)

// maxCommands defines the maximum number of commands that can be stored in the command buffer before further additions are ignored.
const maxCommands = 16384

// InterpretedCommand represents a command with configurable properties for rendering surfaces and visual output.
// It includes text placement, color definitions, and series data for graphical representation.
type InterpretedCommand struct {
	Type       InterpretedCommandType
	Rows       int
	Column     int
	Text       string
	Fg         interfaces.ColorDef
	Bg         interfaces.ColorDef
	Mode       interfaces.ColorMode
	SeriesData []float64
	SeriesW    int
	SeriesH    int
	SeriesMin  float64
	SeriesMax  float64
}

// InterpretedCommandType defines the type of a parsed or processed command in the application.
type InterpretedCommandType int

// InterpretedCommandDrawColor represents a command type for drawing colors.
// InterpretedCommandDrawTextColor represents a command type for drawing text colors.
// InterpretedCommandDrawSeries represents a command type for drawing a series.
const (
	InterpretedCommandDrawColor InterpretedCommandType = iota
	InterpretedCommandDrawTextColor
	InterpretedCommandDrawSeries
)

// InterpretedSurface represents a 2D surface capable of handling interpreted drawing commands like text, colors, and series.
// It enables operations such as drawing, moving cursors, and rendering data through sequences of commands.
// The type maintains the dimensions of the surface, available rows/columns, and a list of accumulated commands for rendering.
type InterpretedSurface struct {
	screenRows    int
	screenColumns int
	windowRows    int
	windowColumns int
	commands      []InterpretedCommand
	moveRows      int
	moveCols      int
}

// NewInterpretedSurface initializes and returns a pointer to a new InterpretedSurface with the specified dimensions.
func NewInterpretedSurface(screenRows int, screenColumns int, windowRows int, windowColumns int) *InterpretedSurface {
	return &InterpretedSurface{
		screenRows:    screenRows,
		screenColumns: screenColumns,
		windowRows:    windowRows,
		windowColumns: windowColumns,
		moveRows:      -1,
		moveCols:      -1,
	}
}

// GetScreenSize returns the number of rows and columns that define the screen size of the surface.
func (s *InterpretedSurface) GetScreenSize() (int, int) {
	return s.screenRows, s.screenColumns
}

// GetWindowSize returns the dimensions of the current window as the number of rows and columns.
func (s *InterpretedSurface) GetWindowSize() (int, int) {
	return s.windowRows, s.windowColumns
}

// DrawColor appends a command to draw a colored rune at the specified position with foreground and background colors, and a mode setting.
func (s *InterpretedSurface) DrawColor(rows int, columns int, text rune, fg interfaces.ColorDef, bg interfaces.ColorDef, mode interfaces.ColorMode) {
	if len(s.commands) > maxCommands {
		return
	}
	s.commands = append(s.commands, InterpretedCommand{
		Type:   InterpretedCommandDrawColor,
		Rows:   rows,
		Column: columns,
		Text:   string(text),
		Fg:     fg,
		Bg:     bg,
		Mode:   mode,
	})
}

// Begin initializes a new sequence of operations or commands on the surface.
func (s *InterpretedSurface) Begin() {
}

// End finalizes the current sequence of drawing or surface commands on the InterpretedSurface.
func (s *InterpretedSurface) End() {
}

// MoveCursor updates the cursor position by setting the specified row and column offsets.
func (s *InterpretedSurface) MoveCursor(rows int, column int) {
	s.moveRows = rows
	s.moveCols = column
}

// GetMoveCursor returns the current cursor position as a tuple of rows and columns.
func (s *InterpretedSurface) GetMoveCursor() (int, int) {
	return s.moveRows, s.moveCols
}

// DrawTextColor queues a command to draw text at a specific location with foreground, background colors, and a color mode.
func (s *InterpretedSurface) DrawTextColor(rows int, column int, text string, fg interfaces.ColorDef, bg interfaces.ColorDef, mode interfaces.ColorMode) {
	if len(s.commands) > maxCommands {
		return
	}
	s.commands = append(s.commands, InterpretedCommand{
		Type:   InterpretedCommandDrawTextColor,
		Rows:   rows,
		Column: column,
		Text:   text,
		Fg:     fg,
		Bg:     bg,
		Mode:   mode,
	})
}

// DrawSeries adds a command to draw a data series with specified dimensions and value range to the command queue.
func (s *InterpretedSurface) DrawSeries(data []float64, w int, h int, min float64, max float64) {
	if len(s.commands) > maxCommands {
		return
	}
	s.commands = append(s.commands, InterpretedCommand{
		Type:       InterpretedCommandDrawSeries,
		SeriesData: data,
		SeriesW:    w,
		SeriesH:    h,
		SeriesMin:  min,
		SeriesMax:  max,
	})
}

// Draw renders a single character at the specified row and column using default colors and normal mode.
func (s *InterpretedSurface) Draw(rows int, column int, c rune) {
	s.DrawColor(rows, column, c, interfaces.ColorNoneDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
}

// DrawText draws a string `c` at the specified row and column positions using default colors and normal display mode.
func (s *InterpretedSurface) DrawText(rows int, column int, c string) {
	s.DrawTextColor(rows, column, c, interfaces.ColorNoneDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
}

// Appy applies a series of commands stored in the surface to the provided ISurface target, executing drawing operations.
func (s *InterpretedSurface) Appy(target interfaces.ISurface) {
	for _, cmd := range s.commands {
		switch cmd.Type {
		case InterpretedCommandDrawColor:
			if len(cmd.Text) > 0 {
				target.DrawColor(cmd.Rows, cmd.Column, []rune(cmd.Text)[0], cmd.Fg, cmd.Bg, cmd.Mode)
			}
		case InterpretedCommandDrawTextColor:
			target.DrawTextColor(cmd.Rows, cmd.Column, cmd.Text, cmd.Fg, cmd.Bg, cmd.Mode)
		case InterpretedCommandDrawSeries:
			rows, columns := target.GetWindowSize()
			w, h := cmd.SeriesW, cmd.SeriesH
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
			g.Setup(cmd.SeriesData, cmd.SeriesMin, cmd.SeriesMax)
			g.Draw(target)
		}
	}
}

// ApplyMoveCursor adjusts the target surface's cursor position to match the moveRows and moveCols values of the current surface.
func (s *InterpretedSurface) ApplyMoveCursor(target interfaces.ISurface) {
	if s.moveRows < 0 || s.moveCols < 0 {
		return
	}
	target.MoveCursor(s.moveRows, s.moveCols)
}
