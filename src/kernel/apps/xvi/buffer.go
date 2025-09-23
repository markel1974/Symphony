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
	lines := strings.Split(content, "\n")
	b := &Buffer{
		lines:    lines, // Direct assignment is fine here
		cursorX:  0,
		cursorY:  0,
		filePath: filePath,
	}
	// Ensure buffer is never truly empty, always has at least one (empty) line.
	if len(b.lines) == 0 {
		b.lines = []string{""}
	}
	return b
}

// GetLine retrieves the content of the line at the specified vertical index y from the buffer.
func (b *Buffer) GetLine(y int) string {
	if y >= 0 && y < len(b.lines) {
		return b.lines[y]
	}
	return ""
}

// DeleteLine removes a line at the specified index y and returns the deleted content.
func (b *Buffer) DeleteLine(y int) string {
	if len(b.lines) == 0 || y < 0 || y >= len(b.lines) {
		return ""
	}

	deletedLine := b.lines[y]
	b.lines = append(b.lines[:y], b.lines[y+1:]...)

	// If we deleted the last line, add an empty one back to prevent an empty buffer.
	if len(b.lines) == 0 {
		b.lines = []string{""}
		b.cursorY = 0
	} else if b.cursorY >= len(b.lines) {
		// Adjust cursor if it's now out of bounds.
		b.cursorY = len(b.lines) - 1
	}
	b.clampCursorX()
	return deletedLine
}

// InsertLineBelow inserts a new line with provided content after the specified index y.
func (b *Buffer) InsertLineBelow(y int, content string) {
	if y < 0 || y >= len(b.lines) {
		// If index is invalid but buffer is not empty, append at the end.
		if len(b.lines) > 0 {
			b.lines = append(b.lines, content)
		}
		return
	}
	head := b.lines[:y+1]
	tail := b.lines[y+1:]
	b.lines = append(append(head, content), tail...)
}

// LineCount returns the total number of lines currently in the buffer.
func (b *Buffer) LineCount() int {
	return len(b.lines)
}

// Cursor returns the current cursor position as (x, y) coordinates in the buffer.
func (b *Buffer) Cursor() (int, int) {
	return b.cursorX, b.cursorY
}

// clampCursorX ensures the horizontal cursor position is valid for the current line.
func (b *Buffer) clampCursorX() {
	if b.cursorY >= 0 && b.cursorY < len(b.lines) {
		lineLen := len(b.lines[b.cursorY])
		if b.cursorX > lineLen {
			b.cursorX = lineLen
		}
	} else {
		b.cursorX = 0
	}
}

// MoveCursor adjusts the cursor position, clamping it within valid boundaries.
func (b *Buffer) MoveCursor(dx, dy int) {
	// Clamp Y
	b.cursorY += dy
	if b.cursorY < 0 {
		b.cursorY = 0
	}
	if b.cursorY >= len(b.lines) {
		b.cursorY = len(b.lines) - 1
	}

	// Clamp X
	b.cursorX += dx
	if b.cursorX < 0 {
		b.cursorX = 0
	}
	b.clampCursorX()
}

// InsertRow splits the current line at the cursor, creating a new line.
func (b *Buffer) InsertRow() {
	if b.cursorY < 0 || b.cursorY >= len(b.lines) {
		return
	}

	line := b.lines[b.cursorY]
	b.clampCursorX() // Ensure cursorX is valid before slicing

	beforeCursor := line[:b.cursorX]
	afterCursor := line[b.cursorX:]

	b.lines[b.cursorY] = beforeCursor
	b.cursorY++
	b.cursorX = 0

	// Insert the new line into the buffer
	head := b.lines[:b.cursorY]
	tail := b.lines[b.cursorY:]
	b.lines = append(append(head, afterCursor), tail...)
}

// InsertChar inserts a character at the current cursor position.
func (b *Buffer) InsertChar(char rune) {
	if b.cursorY < 0 || b.cursorY >= len(b.lines) {
		return
	}
	b.clampCursorX()

	line := b.lines[b.cursorY]
	if b.cursorX >= len(line) {
		line += string(char)
	} else {
		line = line[:b.cursorX] + string(char) + line[b.cursorX:]
	}
	b.lines[b.cursorY] = line
	b.cursorX++
}

// DeleteChar removes the character to the left of the cursor.
func (b *Buffer) DeleteChar() {
	if b.cursorY < 0 || b.cursorY >= len(b.lines) {
		return
	}
	b.clampCursorX()

	if b.cursorX > 0 {
		line := b.lines[b.cursorY]
		line = line[:b.cursorX-1] + line[b.cursorX:]
		b.lines[b.cursorY] = line
		b.cursorX--
	}
}

// Clear removes all lines from the buffer and resets the cursor.
func (b *Buffer) Clear() {
	b.lines = []string{""} // Replace content with a single empty line
	b.cursorX = 0
	b.cursorY = 0
}

// GetContent joins all lines into a single string, representing the file's content.
func (b *Buffer) GetContent() string {
	return strings.Join(b.lines, "\n")
}
