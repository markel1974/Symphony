package loader

import (
	"fmt"
	"io"
)

// CrtChipHeader represents the header structure for a chip in the cartridge, containing metadata and ROM data.
type CrtChipHeader struct {
	skip  uint32 /* bytes to skip after ROM */
	kind  uint16 /* chip type */
	bank  uint16 /* bank number */
	start uint16 /* start address of ROM */
	size  uint16 /* size of ROM in bytes */
	data  []byte
}

// NewChipHeader creates and returns a new instance of CrtChipHeader with default zero values initialized.
func NewChipHeader() *CrtChipHeader {
	return &CrtChipHeader{}
}

func (c *CrtChipHeader) Skip() uint32 {
	return c.skip
}

func (c *CrtChipHeader) Kind() uint16 {
	return c.kind
}

func (c *CrtChipHeader) Bank() uint16 {
	return c.bank
}

func (c *CrtChipHeader) Start() uint16 {
	return c.start
}

func (c *CrtChipHeader) Size() uint16 {
	return c.size
}

func (c *CrtChipHeader) Data() []byte {
	return c.data
}

// Write writes the CrtChipHeader data to the provided io.Writer and returns an error if the operation fails.
func (c *CrtChipHeader) Write(w io.Writer) error {
	//TODO IMPLEMENT
	return fmt.Errorf("unimplemented")
}
