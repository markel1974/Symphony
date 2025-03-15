package loader

import (
	"bytes"
	"fmt"
	"github.com/markel1974/c64emu/src/references"
	"path"
	"strings"
)

// https://sourceforge.net/p/vice-emu/code/HEAD/tree/trunk/vice/src/c64/cart/ocean.c#l240

// crtHeaderLen defines the length of the CRT header in bytes.
// chipHeaderSize specifies the size of the chip header in bytes.
// chipHeaderDef represents the default string identifier for a chip header.
const (
	crtHeaderLen   = 0x40
	chipHeaderSize = 0x10
	chipHeaderDef  = "CHIP"
)

// CRTLoader represents a loader for cartridge files with associated metadata and functionality for parsing CRT headers.
type CRTLoader struct {
	id           string
	rowCartridge []byte
	cursor       int
	mc           MachineType
	Version      uint16 /* version */
	Kind         uint16 /* type of cartridge */
	SubType      uint8  /* subtype/hardware revision of cartridge */
	exRom        int    /* exRom line status */
	game         int    /* game line status */
	name         string /* name of cartridge */
	kind         Type
}

func (cl *CRTLoader) Game() int {
	return cl.game
}

func (cl *CRTLoader) ExRom() int {
	return cl.exRom
}

func (cl *CRTLoader) Name() string {
	//TODO implement me
	return cl.name
}

// NewLoader initializes and returns a new CRTLoader instance with the given id and MachineType.
func NewLoader(id string, mc MachineType) *CRTLoader {
	return &CRTLoader{
		id:           id,
		rowCartridge: nil,
		cursor:       0,
		mc:           mc,
		kind:         TypeBin,
	}
}

// Setup initializes the CRTLoader with the provided id and data, determining the cartridge type and handling .crt files.
func (cl *CRTLoader) Setup(id string, data []byte) error {
	cl.kind = TypeBin
	cl.rowCartridge = data
	lp := strings.ToLower(strings.TrimSpace(id))
	if ext := path.Ext(lp); ext == ".crt" {
		cl.kind = TypeCrt
		if err := cl.open(); err != nil {
			return err
		}
	}
	return nil
}

// GetId returns the unique identifier associated with the CRTLoader instance.
func (cl *CRTLoader) GetId() string {
	return cl.id
}

// GetType returns the cartridge type of the CRTLoader instance as a Type.
func (cl *CRTLoader) GetType() int {
	return int(cl.kind)
}

// GetData retrieves the raw cartridge data as a byte slice.
func (cl *CRTLoader) GetData() []byte {
	return cl.rowCartridge
}

// open initializes the CRTLoader by validating and parsing the cartridge header based on its type and data.
// Returns an error if the header is invalid, improperly formatted, or any mandatory fields are inconsistent.
func (cl *CRTLoader) open() error {
	if cl.kind == TypeBin {
		return nil
	}
	cl.cursor = 0
	var skip uint32
	if crtHeaderLen > len(cl.rowCartridge) {
		return fmt.Errorf("invalid CRT header")
	}
	crtHeader := cl.rowCartridge[:crtHeaderLen]
	cl.cursor = crtHeaderLen
	//cl.Machine = -1
	err := cl.validateMachine(cl.mc, string(crtHeader[0:16]))
	if err != nil {
		return err
	}
	if skip, err = cl.buf2uint32(crtHeader[0x10:0x14]); err != nil {
		return err
	}
	if skip < uint32(len(crtHeader)) {
		return fmt.Errorf("crt header size is wrong (is 0x%02x, expected 0x%02x)", skip, len(crtHeader))
	}
	skip -= uint32(len(crtHeader))
	if cl.Version, err = cl.buf2uint16(crtHeader[0x14:0x16]); err != nil {
		return err
	}
	if cl.Kind, err = cl.buf2uint16(crtHeader[0x16:0x18]); err != nil {
		return err
	}
	cl.SubType = crtHeader[0x1a]
	if exRom := int(crtHeader[0x18]); exRom != 0 {
		cl.exRom = 0
	} else {
		cl.exRom = 1
	}
	if game := int(crtHeader[0x19]); game != 0 {
		cl.game = 0
	} else {
		cl.game = 1
	}
	cl.name = string(crtHeader[0x20:])
	return nil
}

// ReadChipHeader parses the next chip header from the cartridge data and returns the CrtChipHeader or an error if invalid.
func (cl *CRTLoader) ReadChipHeader() (references.ICartridgeChipHeaderC64, error) {
	if cl.kind == TypeBin {
		return nil, nil
	}
	if cl.cursor == len(cl.rowCartridge) {
		return nil, nil
	}
	if cl.cursor+chipHeaderSize > len(cl.rowCartridge) {
		return nil, nil
	}
	chipHeader := cl.rowCartridge[cl.cursor : cl.cursor+chipHeaderSize]
	cl.cursor += chipHeaderSize
	if bytes.Compare(chipHeader[0:4], []byte(chipHeaderDef)) != 0 {
		return nil, fmt.Errorf("invalid header signature")
	}
	var err error
	header := NewChipHeader()
	if header.skip, err = cl.buf2uint32(chipHeader[4:8]); err != nil {
		return nil, err
	}
	if int(header.skip) < chipHeaderSize {
		return nil, fmt.Errorf("invalid packet size")
	}
	header.skip -= uint32(chipHeaderSize)
	if header.size, err = cl.buf2uint16(chipHeader[14:18]); err != nil {
		return nil, err
	}
	if uint32(header.size) > header.skip {
		return nil, fmt.Errorf("rom is bigger then total size")
	}
	header.skip -= uint32(header.size)
	if header.kind, err = cl.buf2uint16(chipHeader[8:12]); err != nil {
		return nil, err
	}
	if header.bank, err = cl.buf2uint16(chipHeader[10:14]); err != nil {
		return nil, err
	}
	if header.start, err = cl.buf2uint16(chipHeader[12:16]); err != nil {
		return nil, err
	}
	if int(header.start)+int(header.size) > 0x10000 {
		return nil, fmt.Errorf("rom crossing the 64k boundary")
	}
	if cl.cursor+int(header.size) > len(cl.rowCartridge) {
		return nil, fmt.Errorf("corrupted data")
	}
	header.data = make([]uint8, int(header.size))
	copy(header.data, cl.rowCartridge[cl.cursor:cl.cursor+int(header.size)])
	//header.Data = cl.rowCartridge[cl.cursor : cl.cursor+int(header.Size)]
	cl.cursor += int(header.size)
	if header.skip > 0 {
		//TODO verify
		cl.cursor += int(header.skip)
	}
	return header, nil
}

// buf2uint32 converts a 4-byte slice into a uint32 value in big-endian order. Returns an error if the slice is too short.
func (cl *CRTLoader) buf2uint32(buf []byte) (uint32, error) {
	if len(buf) < 4 {
		return 0, fmt.Errorf("invalid length")
	}
	t := buf[:4]
	data := uint32(t[3]) | (uint32(t[2]) << 8) | (uint32(t[1]) << 16) | (uint32(t[0]) << 24)
	return data, nil
}

// buf2uint16 converts a 2-byte buffer into a uint16 while ensuring the buffer has sufficient length. Returns an error if invalid.
func (cl *CRTLoader) buf2uint16(buf []byte) (uint16, error) {
	if len(buf) < 2 {
		return 0, fmt.Errorf("invalid length")
	}
	t := buf[:2]
	data := uint16(t[1]) | (uint16(t[0]) << 8)
	return data, nil
}

// validateMachine checks if the given MachineType and id combination is valid based on predefined mappings. Returns error if invalid.
func (cl *CRTLoader) validateMachine(m MachineType, id string) error {
	supported, ok := _machineContainer[id]
	if !ok {
		return fmt.Errorf("invalid crt header")
	}
	_, ok = supported[m]
	if !ok {
		return fmt.Errorf("invalid crt header")
	}
	return nil
}
