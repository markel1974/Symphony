package catalog

import (
	"bytes"
	"fmt"
	"io"
)

// chipHeaderDef defines the header string used for CHIP format identifiers.
// chipHeaderSize specifies the size of the CHIP header in bytes.
const (
	chipHeaderDef  = "CHIP"
	chipHeaderSize = 0x10
)

// CrtChipHeader represents a structured header for a cartridge chip including metadata and associated data.
// Designed for handling cartridge chip data with specific attributes such as skip, kind, bank, start, size, and data.
type CrtChipHeader struct {
	skip  uint32
	kind  uint16
	bank  uint16
	start uint16
	size  uint16
	data  []byte
}

// NewChipHeader creates and returns a new instance of CrtChipHeader with its fields initialized to default zero values.
func NewChipHeader() *CrtChipHeader {
	ch := &CrtChipHeader{}
	return ch
}

// Setup initializes the CrtChipHeader by parsing the given byte slice and updating the cursor position. Returns updated cursor and error.
func (h *CrtChipHeader) Setup(rowCartridge []byte, cursor int) (int, error) {
	if cursor+chipHeaderSize > len(rowCartridge) {
		return cursor, nil
	}
	chipHeader := rowCartridge[cursor : cursor+chipHeaderSize]
	cursor += chipHeaderSize
	if bytes.Compare(chipHeader[0:4], []byte(chipHeaderDef)) != 0 {
		return cursor, fmt.Errorf("invalid header signature")
	}
	var err error
	if h.skip, err = Buf2uint32(chipHeader[4:8]); err != nil {
		return cursor, err
	}
	if int(h.skip) < chipHeaderSize {
		return cursor, fmt.Errorf("invalid packet size")
	}
	h.skip -= uint32(chipHeaderSize)
	if h.size, err = Buf2uint16(chipHeader[14:18]); err != nil {
		return cursor, err
	}
	if uint32(h.size) > h.skip {
		return cursor, fmt.Errorf("rom is bigger then total size")
	}
	h.skip -= uint32(h.size)
	if h.kind, err = Buf2uint16(chipHeader[8:12]); err != nil {
		return cursor, err
	}
	if h.bank, err = Buf2uint16(chipHeader[10:14]); err != nil {
		return cursor, err
	}
	if h.start, err = Buf2uint16(chipHeader[12:16]); err != nil {
		return 0, err
	}
	if int(h.start)+int(h.size) > 0x10000 {
		return 0, fmt.Errorf("rom crossing the 64k boundary")
	}
	if cursor+int(h.size) > len(rowCartridge) {
		return 0, fmt.Errorf("corrupted data")
	}

	h.data = make([]uint8, int(h.size))
	copy(h.data, rowCartridge[cursor:cursor+int(h.size)])
	cursor += int(h.size)

	if h.skip > 0 {
		//TODO verify
		cursor += int(h.skip)
	}
	return cursor, nil
}

// Skip returns the `skip` field from the CrtChipHeader, representing the remaining size after the chip header and data.
func (h *CrtChipHeader) Skip() uint32 {
	return h.skip
}

// Kind returns the kind field of the CrtChipHeader as a uint16.
func (h *CrtChipHeader) Kind() uint16 {
	return h.kind
}

// Bank returns the bank identifier as a 16-bit unsigned integer from the CrtChipHeader.
func (h *CrtChipHeader) Bank() uint16 {
	return h.bank
}

// Start retrieves the starting address of the chip data within the CrtChipHeader. It returns the value as a uint16.
func (h *CrtChipHeader) Start() uint16 {
	return h.start
}

// Size returns the size of the chip header as a 16-bit unsigned integer.
func (h *CrtChipHeader) Size() uint16 {
	return h.size
}

// Data retrieves the data slice associated with the CrtChipHeader instance.
func (h *CrtChipHeader) Data() []byte {
	return h.data
}

// Write writes the CrtChipHeader's contents to the provided io.Writer and returns an error if the write operation fails.
func (h *CrtChipHeader) Write(w io.Writer) error {
	//TODO IMPLEMENT
	return fmt.Errorf("unimplemented")
}
