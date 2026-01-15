package tetris

import (
	"github.com/markel1974/symphony/src/kernel/interfaces"
)

// Cell represents a unit on a grid with two-dimensional coordinates (x, y) and an associated color.
type Cell struct {
	x, y  int
	color interfaces.ColorDef
}

// NewCell creates a new Cell instance with the specified x, y coordinates and assigns a color based on the input rune.
func NewCell(x, y int, ch rune) *Cell {
	return &Cell{x: x, y: y, color: colorMapping[ch]}
}

// conflicts checks if the current cell either overlaps another cell or is out of bounds on the given board.
func (c *Cell) conflicts(board *Board) bool {
	return c.isOnWall(board) || c.isOverlapped(board)
}

// isOverlapped checks if the cell overlaps with a non-blank cell on the board, returning true if it does.
func (c *Cell) isOverlapped(board *Board) bool {
	if !board.isOnBoard(c.x, c.y) {
		return false
	}
	return board.colors[c.x][c.y] != blankColor
}

// isOnWall checks if the cell is outside the board's boundaries on any side.
func (c *Cell) isOnWall(board *Board) bool {
	return c.x < 0 || board.w <= c.x || board.h <= c.y
}
