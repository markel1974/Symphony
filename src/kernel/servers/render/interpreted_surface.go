package render

import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/servers/render/plotter"
)

// InterpretedCommand represents a drawing action to be processed, including text, colors, positions, and series data.
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

// InterpretedCommandType represents the type of command to be executed in a drawing operation or similar context.
type InterpretedCommandType int

// InterpretedCommandDrawColor represents a command for drawing a color.
// InterpretedCommandDrawTextColor represents a command for drawing text color.
// InterpretedCommandDrawSeries represents a command for drawing a series.
const (
	InterpretedCommandDrawColor InterpretedCommandType = iota
	InterpretedCommandDrawTextColor
	InterpretedCommandDrawSeries
)

// InterpretedSurface is a type used for storing and managing a sequence of drawing commands.
// It enables creation and manipulation of visual elements through DrawCommand instances.
// Commands can be executed on a target surface to reflect the desired graphical output.
type InterpretedSurface struct {
	rows     int
	columns  int
	commands []InterpretedCommand
}

func NewInterpretedSurface(rows int, columns int) *InterpretedSurface {
	return &InterpretedSurface{
		rows:    rows,
		columns: columns,
	}
}

// DrawColor adds a drawing command to the surface, specifying position, rune, foreground, background colors, and color mode.
func (s *InterpretedSurface) DrawColor(rows int, columns int, text rune, fg interfaces.ColorDef, bg interfaces.ColorDef, mode interfaces.ColorMode) {
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

func (s *InterpretedSurface) GetSize() (int, int) {
	return s.rows, s.columns
}

func (s *InterpretedSurface) Begin() {
}

func (s *InterpretedSurface) End() {
}

// DrawTextColor queues a command to render a string of text at the specified row and column with foreground and background colors.
func (s *InterpretedSurface) DrawTextColor(rows int, column int, text string, fg interfaces.ColorDef, bg interfaces.ColorDef, mode interfaces.ColorMode) {
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

// DrawSeries appends a DrawCommand for rendering a series to the command list with its data, dimensions, and range.
func (s *InterpretedSurface) DrawSeries(data []float64, w int, h int, min float64, max float64) {
	s.commands = append(s.commands, InterpretedCommand{
		Type:       InterpretedCommandDrawSeries,
		SeriesData: data,
		SeriesW:    w,
		SeriesH:    h,
		SeriesMin:  min,
		SeriesMax:  max,
	})
}

// Draw renders a single rune at the specified row and column using default color and mode settings.
func (s *InterpretedSurface) Draw(rows int, column int, c rune) {
	s.DrawColor(rows, column, c, interfaces.ColorNoneDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
}

// DrawText renders a string at the specified row and column without applying colors or modes.
func (s *InterpretedSurface) DrawText(rows int, column int, c string) {
	s.DrawTextColor(rows, column, c, interfaces.ColorNoneDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
}

// Appy applies a series of drawing commands to a target surface based on their type and provided parameters.
func (s *InterpretedSurface) Appy(target interfaces.ISurface) {
	for _, cmd := range s.commands {
		switch cmd.Type {
		case InterpretedCommandDrawColor:
			target.DrawColor(cmd.Rows, cmd.Column, []rune(cmd.Text)[0], cmd.Fg, cmd.Bg, cmd.Mode)
		case InterpretedCommandDrawTextColor:
			target.DrawTextColor(cmd.Rows, cmd.Column, cmd.Text, cmd.Fg, cmd.Bg, cmd.Mode)
		case InterpretedCommandDrawSeries:
			rows, columns := target.GetSize()
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
