package drivers

import (
	"io"

	"github.com/markel1974/c64emu/src/kernel/interfaces"
)

type KeyboardTerminal struct {
	reader   io.Reader
	terminal interfaces.ITerminal
}

func NewKeyboardTerminal(reader io.Reader, terminal interfaces.ITerminal) *KeyboardTerminal {
	return &KeyboardTerminal{
		reader:   reader,
		terminal: terminal,
	}
}

// ScanKey processes the provided byte data using the underlying terminal's Scan method.
func (c *KeyboardTerminal) ScanKey(readBuffer []byte) (interfaces.KeyType, rune, error) {
	n, err := c.reader.Read(readBuffer)
	if err != nil {
		return interfaces.KeyTypeNone, 0, err
	}
	if n == 0 {
		return interfaces.KeyTypeNone, 0, nil
	}
	k, v := c.terminal.CreateScanKey(readBuffer)
	return k, v, nil
}
