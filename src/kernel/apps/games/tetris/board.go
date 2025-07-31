package tetris

import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
)

// Board represents a grid-based game board with a specific width and height.
// Each cell in the grid contains a color defined by interfaces.ColorDef.
// The dimensions of the board are specified by the fields w (width) and h (height).
type Board struct {
	colors [][]interfaces.ColorDef
	w      int
	h      int
}

// NewBoard creates a new Board with the specified width and height, initializing all cells to the default blank color.
func NewBoard(w int, h int) *Board {
	b := &Board{
		w: w,
		h: h,
	}
	b.colors = make([][]interfaces.ColorDef, w)
	for r := range b.colors {
		b.colors[r] = make([]interfaces.ColorDef, h)
		for c := range b.colors[r] {
			b.colors[r][c] = blankColor
		}
	}
	return b
}

// deleteLine removes a specific horizontal line at the given `y` index and shifts all lines above it down by one.
func (b *Board) deleteLine(y int) {
	for i := 0; i < b.w; i++ {
		b.colors[i][y] = blankColor
	}
	for j := y; j > 0; j-- {
		for i := 0; i < b.w; i++ {
			b.colors[i][j] = b.colors[i][j-1]
		}
	}
	for i := 0; i < b.w; i++ {
		b.colors[i][0] = blankColor
	}
}

// fullLines returns a slice of integers representing the indices of fully filled lines on the board.
func (b *Board) fullLines() []int {
	var fullLines []int
	for j := 0; j < b.h; j++ {
		if b.isFullLine(j) {
			fullLines = append(fullLines, j)
		}
	}
	return fullLines
}

// isFullLine checks if the row at the given y-coordinate is completely filled with non-blank colors.
func (b *Board) isFullLine(y int) bool {
	hasBlank := false
	for i := 0; i < b.w; i++ {
		if b.colors[i][y] == blankColor {
			hasBlank = true
			break
		}
	}
	return !hasBlank
}

// hasFullLine checks if the board contains at least one full line of cells.
// Returns true if a full line is found; otherwise, returns false.
func (b *Board) hasFullLine() bool {
	for j := 0; j < b.h; j++ {
		if b.isFullLine(j) {
			return true
		}
	}
	return false
}

/*
func (b *Board) text() string {
	text := ""
	for j := 0; j < b.h; j++ {
		for i := 0; i < b.w; i++ {
			text = fmt.Sprintf("%s%c", text, charByColor(b.colors[i][j]))
		}
		text = fmt.Sprintf("%s\n", text)
	}
	return text
}
*/

// setCell sets the color of the specified cell on the board using its x and y coordinates.
func (b *Board) setCell(cell *Cell) {
	b.colors[cell.x][cell.y] = cell.color
}

// setCells sets the colors of the specified cells on the board by updating their positions with the provided values.
func (b *Board) setCells(cells []*Cell) {
	for _, cell := range cells {
		b.setCell(cell)
	}
}

// isOnBoard checks if the given coordinates (x, y) are within the boundaries of the board.
func (b *Board) isOnBoard(x, y int) bool {
	return (0 <= x && x < b.w) && (0 <= y && y < b.h)
}
