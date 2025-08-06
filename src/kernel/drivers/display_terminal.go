package drivers

import (
	"io"

	"github.com/markel1974/c64emu/src/kernel/interfaces"
)

// DisplayTerminal is a struct that wraps an ITerminal interface and an io.Writer for managing terminal outputs.
type DisplayTerminal struct {
	interfaces.ITerminal
	writer io.Writer
}

// NewDisplayTerminal creates and returns a new DisplayTerminal instance with the provided writer and ITerminal implementation.
func NewDisplayTerminal(writer io.Writer, terminal interfaces.ITerminal) *DisplayTerminal {
	return &DisplayTerminal{
		writer:    writer,
		ITerminal: terminal,
	}
}

// Write writes the provided byte slice to the underlying writer and returns the number of bytes written and any errors encountered.
func (v *DisplayTerminal) Write(p []byte) (n int, err error) {
	return v.writer.Write(p)
}
