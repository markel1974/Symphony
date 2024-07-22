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

const chipHeaderSize = 0x10

var ChipHeader = []byte("CHIP")

var machineContainer = map[string]map[int]bool{
	"C64 CARTRIDGE   ": {MachineC64: true, MachineC64SC: true, MachineC128: true, MachineSCpu64: true},
	"C128 CARTRIDGE  ": {MachineC128: true},
	"CBM2 CARTRIDGE  ": {MachineCbm6x0: true, MachineCbm5x0: true},
	"VIC20 CARTRIDGE ": {MachineVic20: true},
	"PLUS4 CARTRIDGE ": {MachinePlus4: true},
}

type CRTLoader struct {
	id           string
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

func NewLoader(id string, mc int) *CRTLoader {
	return &CRTLoader{
		id:           id,
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

func (cl *CRTLoader) GetId() string {
	return cl.id
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
	cl.Machine = -1
	supported, ok := machineContainer[string(crtHeader[0:16])]
	if !ok {
		return fmt.Errorf("invalid crt header")
	}
	_, ok = supported[cl.mc]
	if !ok {
		return fmt.Errorf("invalid crt header")
	}
	var err error
	if skip, err = buf2dword(crtHeader[0x10:0x14]); err != nil {
		return err
	}
	if skip < uint32(len(crtHeader)) {
		return fmt.Errorf("crt header size is wrong (is 0x%02x, expected 0x%02x)", skip, len(crtHeader))
	}
	skip -= uint32(len(crtHeader))
	if cl.Version, err = buf2word(crtHeader[0x14:0x16]); err != nil {
		return err
	}
	if cl.Kind, err = buf2word(crtHeader[0x16:0x18]); err != nil {
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

func (cl *CRTLoader) ReadChipHeader() (*CrtChipHeader, error) {
	if cl.mode == ModeBin {
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
	if bytes.Compare(chipHeader[0:4], ChipHeader) != 0 {
		return nil, fmt.Errorf("invalid header signature")
	}
	var err error
	header := NewChipHeader()
	if header.Skip, err = buf2dword(chipHeader[4:8]); err != nil {
		return nil, err
	}
	if int(header.Skip) < chipHeaderSize {
		return nil, fmt.Errorf("invalid packet size")
	}
	header.Skip -= uint32(chipHeaderSize)
	if header.Size, err = buf2word(chipHeader[14:18]); err != nil {
		return nil, err
	}
	if uint32(header.Size) > header.Skip {
		return nil, fmt.Errorf("rom is bigger then total size")
	}
	header.Skip -= uint32(header.Size)
	if header.Kind, err = buf2word(chipHeader[8:12]); err != nil {
		return nil, err
	}
	if header.Bank, err = buf2word(chipHeader[10:14]); err != nil {
		return nil, err
	}
	if header.Start, err = buf2word(chipHeader[12:16]); err != nil {
		return nil, err
	}
	if int(header.Start)+int(header.Size) > 0x10000 {
		return nil, fmt.Errorf("rom crossing the 64k boundary")
	}
	if cl.cursor+int(header.Size) > len(cl.rowCartridge) {
		return nil, fmt.Errorf("corrupted data")
	}
	header.Data = make([]uint8, int(header.Size))
	copy(header.Data, cl.rowCartridge[cl.cursor:cl.cursor+int(header.Size)])
	//header.Data = cl.rowCartridge[cl.cursor : cl.cursor+int(header.Size)]
	cl.cursor += int(header.Size)
	if header.Skip > 0 {
		//TODO verify
		cl.cursor += int(header.Skip)
	}
	return header, nil
}
