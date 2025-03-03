package loader

import (
	"fmt"
	"io"
)

// CrtChipHeader represents the header structure for a chip in the cartridge, containing metadata and ROM data.
type CrtChipHeader struct {
	Skip  uint32 /* bytes to skip after ROM */
	Kind  uint16 /* chip type */
	Bank  uint16 /* bank number */
	Start uint16 /* start address of ROM */
	Size  uint16 /* size of ROM in bytes */
	Data  []byte
}

// NewChipHeader creates and returns a new instance of CrtChipHeader with default zero values initialized.
func NewChipHeader() *CrtChipHeader {
	return &CrtChipHeader{}
}

// Write writes the CrtChipHeader data to the provided io.Writer and returns an error if the operation fails.
func (c *CrtChipHeader) Write(w io.Writer) error {
	//TODO IMPLEMENT
	return fmt.Errorf("unimplemented")
}
