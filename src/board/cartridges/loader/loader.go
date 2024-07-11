package loader

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"
)

const (
	CARTRIDGE_FILETYPE_BIN = 0
	CARTRIDGE_FILETYPE_CRT = 1
)

// https://sourceforge.net/p/vice-emu/code/HEAD/tree/trunk/vice/src/c64/cart/ocean.c#l240

type Mode int

const (
	ModeBin = iota
	ModeCrt
)

const (
	MachineC64 = iota
	MachineC128
	MachineVic20
	MachinePet
	MachineCbm5x0
	MachineCbm6x0
	MachinePlus4
	MachineC64DTV
	MachineC64SC
	MachineVSid
	MachineSCpu64
)

const (
	CrtHeaderLen = 0x40
)

var ChipHeader = []byte("CHIP")

var CrtHeaderC64 = []byte("C64 CARTRIDGE   ")
var CrtHeaderC128 = []byte("C128 CARTRIDGE  ")
var CrtHeaderCbm2 = []byte("CBM2 CARTRIDGE  ")
var CrtHeaderVic20 = []byte("VIC20 CARTRIDGE ")
var CrtHeaderPlus4 = []byte("PLUS4 CARTRIDGE ")

type CRTLoader struct {
	rowCartridge []byte
	cursor       int
	mc           int
	Version      uint16 /* version */
	Kind         uint16 /* type of cartridge */
	SubType      uint8  /* subtype/hardware revision of cartridge */
	ExRom        int    /* exRom line status */
	Game         int    /* game line status */
	Name         string /* name of cartridge */
	Machine      int    /* detected machine for this crt file */
	mode         Mode
}

func NewLoader(mc int) *CRTLoader {
	return &CRTLoader{
		rowCartridge: nil,
		cursor:       0,
		mc:           mc,
		mode:         ModeBin,
	}
}

func (cl *CRTLoader) Setup(p string) error {
	data, err := os.ReadFile(p)
	if err != nil {
		return err
	}
	cl.mode = ModeBin
	cl.rowCartridge = data
	lp := strings.ToLower(strings.TrimSpace(p))
	if ext := path.Ext(lp); ext == ".crt" {
		cl.mode = ModeCrt
		if err = cl.open(); err != nil {
			return err
		}
	}
	return nil
}

func (cl *CRTLoader) GetMode() Mode {
	return cl.mode
}

func (cl *CRTLoader) GetData() []byte {
	return cl.rowCartridge
}

func (cl *CRTLoader) open() error {
	if cl.mode == ModeBin {
		return nil
	}
	cl.cursor = 0
	var skip uint32
	if CrtHeaderLen > len(cl.rowCartridge) {
		return fmt.Errorf("invalid CRT header")
	}
	crtHeader := cl.rowCartridge[:CrtHeaderLen]
	cl.cursor = CrtHeaderLen
	for {
		cl.Machine = -1
		if bytes.Compare(crtHeader[:16], CrtHeaderC64) == 0 {
			cl.Machine = MachineC64
			if !(cl.mc == MachineC64 || cl.mc == MachineC64SC || cl.mc == MachineC128 || cl.mc == MachineSCpu64) {
				return fmt.Errorf("invalid crt header")
			}
		} else if bytes.Compare(crtHeader[:16], CrtHeaderC128) == 0 {
			cl.Machine = MachineC128
			if !(cl.mc == MachineC128) {
				return fmt.Errorf("invalid crt header")
			}
		} else if bytes.Compare(crtHeader[:16], CrtHeaderCbm2) == 0 {
			cl.Machine = MachineCbm6x0
			if !(cl.mc == MachineCbm5x0) || (cl.mc == MachineCbm6x0) {
				return fmt.Errorf("invalid crt header")
			}
		} else if bytes.Compare(crtHeader[:16], CrtHeaderVic20) == 0 {
			fmt.Printf("Found header: '%s'\n", CrtHeaderVic20)
			cl.Machine = MachineVic20
			if !(cl.mc == MachineVic20) {
				return fmt.Errorf("invalid crt header")
			}
		} else if bytes.Compare(crtHeader[:16], CrtHeaderPlus4) == 0 {
			fmt.Printf("Found header: '%s'\n", CrtHeaderPlus4)
			cl.Machine = MachinePlus4
			if !(cl.mc == MachinePlus4) {
				return fmt.Errorf("invalid crt header")
			}
		} else {
			return fmt.Errorf("no crt header found")
		}
		var err error
		if skip, err = buf2dword(crtHeader[0x10:]); err != nil {
			return err
		}
		if skip < uint32(len(crtHeader)) {
			return fmt.Errorf("crt header size is wrong (is 0x%02x, expected 0x%02x)", skip, len(crtHeader))
		}
		skip -= uint32(len(crtHeader))
		if cl.Version, err = buf2word(crtHeader[0x14:]); err != nil {
			return err
		}
		if cl.Kind, err = buf2word(crtHeader[0x16:]); err != nil {
			return err
		}
		cl.SubType = crtHeader[0x1a]
		if exRom := int(crtHeader[0x18]); exRom != 0 {
			cl.ExRom = 0
		} else {
			cl.ExRom = 1
		}
		if game := int(crtHeader[0x19]); game != 0 {
			cl.Game = 0
		} else {
			cl.Game = 1
		}
		cl.Name = string(crtHeader[0x20:])
		return nil
	}
}

// ReadChipHeader
// crtReadChipHeader
func (cl *CRTLoader) ReadChipHeader() (*CrtChipHeader, error) {
	if cl.mode == ModeBin {
		return nil, nil
	}
	const chipHeaderSize = 0x10
	if cl.cursor+chipHeaderSize > len(cl.rowCartridge) {
		return nil, nil
	}
	rowCartridge := cl.rowCartridge[cl.cursor:]
	header := NewChipHeader()
	chipHeader := rowCartridge[:chipHeaderSize]
	if bytes.Compare(chipHeader[:4], ChipHeader) != 0 {
		return nil, fmt.Errorf("invalid header signature")
	}
	var err error
	if header.Skip, err = buf2dword(chipHeader[4:]); err != nil {
		return nil, err
	}
	if int(header.Skip) < chipHeaderSize {
		return nil, fmt.Errorf("invalid packet size")
	}
	header.Skip -= uint32(chipHeaderSize)
	if header.Size, err = buf2word(chipHeader[14:]); err != nil {
		return nil, err
	}
	if uint32(header.Size) > header.Skip {
		return nil, fmt.Errorf("rom bigger then total size")
	}
	header.Skip -= uint32(header.Size)
	if header.Kind, err = buf2word(chipHeader[8:]); err != nil {
		return nil, err
	}
	if header.Bank, err = buf2word(chipHeader[10:]); err != nil {
		return nil, err
	}
	if header.Start, err = buf2word(chipHeader[12:]); err != nil {
		return nil, err
	}
	if int(header.Start)+int(header.Size) > 0x10000 {
		return nil, fmt.Errorf("rom crossing the 64k boundary")
	}
	cl.cursor += chipHeaderSize
	if cl.cursor+int(header.Size) > len(cl.rowCartridge) {
		return nil, fmt.Errorf("corrupted data")
	}
	header.Data = cl.rowCartridge[cl.cursor : cl.cursor+int(header.Size)]
	cl.cursor += int(header.Size)
	if header.Skip > 0 {
		//TODO verify
		cl.cursor += int(header.Skip)
	}
	return header, nil
}
