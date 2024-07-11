package loader

import (
	"fmt"
	"io"
)

type CrtChipHeader struct {
	Skip  uint32 /* bytes to skip after ROM */
	Kind  uint16 /* chip type */
	Bank  uint16 /* bank number */
	Start uint16 /* start address of ROM */
	Size  uint16 /* size of ROM in bytes */
	Data  []byte
}

func NewChipHeader() *CrtChipHeader {
	return &CrtChipHeader{}
}

func (c *CrtChipHeader) Write(w io.Writer) error {
	//TODO IMPLEMENT
	return fmt.Errorf("unimplemented")
}
