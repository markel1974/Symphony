package loader

import (
	"bytes"
	"fmt"
	"path"
	"strings"
)

// https://sourceforge.net/p/vice-emu/code/HEAD/tree/trunk/vice/src/c64/cart/ocean.c#l240

const (
	crtHeaderLen   = 0x40
	chipHeaderSize = 0x10
	chipHeaderDef  = "CHIP"
)

type CRTLoader struct {
	id           string
	rowCartridge []byte
	cursor       int
	mc           MachineType
	Version      uint16 /* version */
	Kind         uint16 /* type of cartridge */
	SubType      uint8  /* subtype/hardware revision of cartridge */
	ExRom        int    /* exRom line status */
	Game         int    /* game line status */
	Name         string /* name of cartridge */
	kind         Type
}

func NewLoader(id string, mc MachineType) *CRTLoader {
	return &CRTLoader{
		id:           id,
		rowCartridge: nil,
		cursor:       0,
		mc:           mc,
		kind:         TypeBin,
	}
}

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

func (cl *CRTLoader) GetId() string {
	return cl.id
}

func (cl *CRTLoader) GetType() Type {
	return cl.kind
}

func (cl *CRTLoader) GetData() []byte {
	return cl.rowCartridge
}

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
	if header.Skip, err = cl.buf2uint32(chipHeader[4:8]); err != nil {
		return nil, err
	}
	if int(header.Skip) < chipHeaderSize {
		return nil, fmt.Errorf("invalid packet size")
	}
	header.Skip -= uint32(chipHeaderSize)
	if header.Size, err = cl.buf2uint16(chipHeader[14:18]); err != nil {
		return nil, err
	}
	if uint32(header.Size) > header.Skip {
		return nil, fmt.Errorf("rom is bigger then total size")
	}
	header.Skip -= uint32(header.Size)
	if header.Kind, err = cl.buf2uint16(chipHeader[8:12]); err != nil {
		return nil, err
	}
	if header.Bank, err = cl.buf2uint16(chipHeader[10:14]); err != nil {
		return nil, err
	}
	if header.Start, err = cl.buf2uint16(chipHeader[12:16]); err != nil {
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

func (cl *CRTLoader) buf2uint32(buf []byte) (uint32, error) {
	if len(buf) < 4 {
		return 0, fmt.Errorf("invalid length")
	}
	t := buf[:4]
	data := uint32(t[3]) | (uint32(t[2]) << 8) | (uint32(t[1]) << 16) | (uint32(t[0]) << 24)
	return data, nil
}

func (cl *CRTLoader) buf2uint16(buf []byte) (uint16, error) {
	if len(buf) < 2 {
		return 0, fmt.Errorf("invalid length")
	}
	t := buf[:2]
	data := uint16(t[1]) | (uint16(t[0]) << 8)
	return data, nil
}

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
