// file: kernel/apps/xvi/buffer.go
package xvi

import "strings"

// Buffer represents a text editor's internal representation of a file or document with support for line-by-line operations.
type Buffer struct {
	lines    []string
	cursorX  int
	cursorY  int
	filePath string
}

// NewBuffer creates and returns a new Buffer instance initialized with the given file path and content.
func NewBuffer(filePath string, content string) *Buffer {
	return &Buffer{
		lines:    strings.Split(content, "\n"),
		cursorX:  0,
		cursorY:  0,
		filePath: filePath,
	}
}

// GetLine retrieves the content of the line at the specified vertical index y from the buffer.
// Returns an empty string if the index y is out of range.
func (b *Buffer) GetLine(y int) string {
	if y >= 0 && y < len(b.lines) {
		return b.lines[y]
	}
	return ""
}

// LineCount returns the total number of lines currently in the buffer.
func (b *Buffer) LineCount() int {
	return len(b.lines)
}

// Cursor returns the current cursor position as (x, y) coordinates in the buffer.
func (b *Buffer) Cursor() (int, int) {
	return b.cursorX, b.cursorY
}

// MoveCursor adjusts the cursor position by the specified deltas (dx, dy), clamping it within valid boundaries.
func (b *Buffer) MoveCursor(dx, dy int) {
	b.cursorY += dy
	if b.cursorY < 0 {
		b.cursorY = 0
	}
	if b.cursorY >= len(b.lines) {
		b.cursorY = len(b.lines) - 1
	}

	b.cursorX += dx
	if b.cursorX < 0 {
		b.cursorX = 0
	}
	// Clamp X to the end of the current line
	lineLen := len(b.lines[b.cursorY])
	if b.cursorX > lineLen {
		b.cursorX = lineLen
	}
}

// InsertChar inserts a given character at the current cursor position in the buffer and moves the cursor to the right.
func (b *Buffer) InsertChar(char rune) {
	line := b.lines[b.cursorY]
	if b.cursorX >= len(line) {
		line += string(char)
	} else {
		line = line[:b.cursorX] + string(char) + line[b.cursorX:]
	}
	b.lines[b.cursorY] = line
	b.cursorX++
}

// DeleteChar removes the character to the left of the cursor in the current line and updates the cursor's position.
func (b *Buffer) DeleteChar() {
	if b.cursorX > 0 {
		line := b.lines[b.cursorY]
		line = line[:b.cursorX-1] + line[b.cursorX:]
		b.lines[b.cursorY] = line
		b.cursorX--
	}
}
