package catalog

import (
	"fmt"
	"github.com/markel1974/c64emu/src/references"
	"path"
	"strings"
)

// https://sourceforge.net/p/vice-emu/code/HEAD/tree/trunk/vice/src/c64/cart/ocean.c#l240

// crtHeaderLen defines the fixed length in bytes of the CRT header used for identifying and parsing cartridge data.
const (
	crtHeaderLen = 0x40
)

// Type defines an integer-based enumeration typically used for specifying cartridge or loader types in the system.
type Type int

// TypeBin represents a type enumeration with a binary identifier.
// TypeCrt represents a type enumeration with a certificate identifier.
const (
	TypeBin = Type(iota)
	TypeCrt
)

// Loader represents a structure for handling a cartridge file with relevant metadata and methods for processing.
type Loader struct {
	id           string
	rowCartridge []byte
	cursor       int
	mc           MachineType
	Version      uint16 // version
	Kind         uint16 // type of cartridge
	SubType      uint8  // subtype/hardware revision of cartridge
	exRom        int    // exRom line status
	game         int    // game line status
	name         string // name of cartridge
	kind         Type
}

// NewLoader initializes and returns a new CRTLoader instance with the specified ID and machine type.
func NewLoader(id string, mc MachineType) *Loader {
	return &Loader{
		id:           id,
		rowCartridge: nil,
		cursor:       0,
		mc:           mc,
		kind:         TypeBin,
	}
}

// Setup initializes the CRTLoader with the given identifier and data, determining its type and processing accordingly.
func (cl *Loader) Setup(id string, data []byte) error {
	cl.kind = TypeBin
	cl.rowCartridge = data
	lp := strings.ToLower(strings.TrimSpace(id))
	if ext := path.Ext(lp); ext == ".crt" {
		cl.kind = TypeCrt
		if err := cl.read(); err != nil {
			return err
		}
	}
	return nil
}

// GetId returns the unique identifier of the CRTLoader instance as a string.
func (cl *Loader) GetId() string {
	return cl.id
}

// GetType returns the type of the cartridge as an integer by converting the internal kind field to an int.
func (cl *Loader) GetType() int {
	return int(cl.kind)
}

// GetData returns the raw cartridge data as a byte slice.
func (cl *Loader) GetData() []byte {
	return cl.rowCartridge
}

// Game returns the current game line status of the CRTLoader instance.
func (cl *Loader) Game() int {
	return cl.game
}

// ExRom retrieves the status of the ExRom line from the CRTLoader instance. Returns 0 or 1 based on the ExRom state.
func (cl *Loader) ExRom() int {
	return cl.exRom
}

// Name retrieves the name of the cartridge associated with the CRTLoader instance.
func (cl *Loader) Name() string {
	return cl.name
}

// read initializes the CRTLoader by parsing the CRT header and validating its format. Returns an error if validation fails.
func (cl *Loader) read() error {
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
	err := ValidateMachine(cl.mc, string(crtHeader[0:16]))
	if err != nil {
		return err
	}
	if skip, err = Buf2uint32(crtHeader[0x10:0x14]); err != nil {
		return err
	}
	if skip < uint32(len(crtHeader)) {
		return fmt.Errorf("crt header size is wrong (is 0x%02x, expected 0x%02x)", skip, len(crtHeader))
	}
	skip -= uint32(len(crtHeader))
	if cl.Version, err = Buf2uint16(crtHeader[0x14:0x16]); err != nil {
		return err
	}
	if cl.Kind, err = Buf2uint16(crtHeader[0x16:0x18]); err != nil {
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
	name := crtHeader[0x20:]
	beginName := 0
	endName := 0
	for x := 0; x < len(name); x++ {
		if name[x] == 0 {
			endName = x
			break
		}
	}
	cl.name = string(name[beginName:endName])
	return nil
}

// ReadChipHeader reads and initializes a chip header from the cartridge data at the current cursor position.
// Returns the parsed ICartridgeChipHeaderC64 instance or an error if the operation fails.
func (cl *Loader) ReadChipHeader() (references.ICartridgeChipHeaderC64, error) {
	if cl.kind == TypeBin {
		return nil, nil
	}
	if cl.cursor == len(cl.rowCartridge) {
		return nil, nil
	}
	header := NewChipHeader()
	cursor, err := header.Setup(cl.rowCartridge, cl.cursor)
	if err != nil {
		return nil, err
	}
	cl.cursor = cursor
	if header == nil {
		return nil, nil
	}
	return header, nil
}
