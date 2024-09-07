package mos6581

import (
	"github.com/markel1974/c64emu/src/config"
)

type SID struct {
	id      string
	socket  ISocket
	regs    []uint8
	cfg     *config.Config
	builder *Builder
}

func NewSID(id string) *SID {
	s := &SID{
		id:   id,
		regs: make([]uint8, RegisterCount),
		cfg:  nil,
	}
	return s
}

func (sid *SID) Setup(socket ISocket, cfg *config.Config, fragFreq int, rasters int) {
	sid.socket = socket
	sid.cfg = cfg
	sid.builder = NewBuilder(socket.GetPlayer(), true, fragFreq, rasters)
	sid.cfg.Bind(sid.configChanged)
}

func (sid *SID) SetPotXSlot(pot uint8) {
	// PX7 PX6 PX5 PX4 PX3 PX2 PX1 PX0
	sid.regs[25] = pot
}

func (sid *SID) SetPotYSlot(pot uint8) {
	//PY7 PY6 PY5 PY4 PY3 PY2 PY1 PY0
	sid.regs[26] = pot
}

func (sid *SID) configChanged() {
	//TODO
}

func (sid *SID) Reset() {
	for x := range sid.regs {
		sid.regs[x] = 0
	}
	sid.builder.Reset()

	//PADDLE TEST
	//sid.WriteRegister(0xDC00, 0x40) //Control-Port 1 selected
	//sid.WriteRegister(0xD419, 0x7F) //Paddle X value
	//sid.WriteRegister(0xD419, 0x7F) //Paddle Y value
}

func (sid *SID) ReadRegister(addr uint16) uint8 {
	reg := addr & 0x1f
	v := sid.regs[reg]
	//fmt.Printf("[%s][ReadRegister] addr %X [%x] -> %d\n", sid.id, addr, reg, v)
	return v
}

func (sid *SID) WriteRegister(addr uint16, data uint8) {
	reg := addr & 0x1f
	//fmt.Printf("[%s][WriteRegister] addr %X [%x] -> %d\n", sid.id, addr, reg, data)
	sid.regs[reg] = data
}

func (sid *SID) Emulate(vBlank bool, lastVicCycle bool) {
	if vBlank {
		sid.builder.Render()
	}
	if lastVicCycle {
		sid.builder.AddToHistory(sid.regs)
	}
}

func (sid *SID) GetLastByte() uint8 {
	return 0
}
