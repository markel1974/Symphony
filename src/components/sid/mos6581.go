package mos6581

import (
	"github.com/markel1974/c64emu/src/config"
)

type SID struct {
	regs             []uint8
	regsHistory      [][]uint8
	regsHistoryIndex uint32
	cfg              *config.Config
}

func NewSID() *SID {
	s := &SID{
		regs:             make([]uint8, RegisterCount),
		regsHistory:      make([][]uint8, RegisterHistory),
		regsHistoryIndex: 0,
		cfg:              nil,
	}
	for x := range s.regsHistory {
		s.regsHistory[x] = make([]uint8, RegisterCount)
	}
	return s
}

func (sid *SID) Setup(cfg *config.Config) {
	sid.cfg = cfg
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
	for x := range sid.regsHistory {
		for y := range sid.regsHistory[x] {
			sid.regsHistory[x][y] = 0
		}
	}
	sid.regsHistoryIndex = 0
}

func (sid *SID) ReadRegister(addr uint16) uint8 {
	addr = addr & 0x1f
	return sid.regs[addr]
}

func (sid *SID) WriteRegister(addr uint16, data uint8) {
	addr = addr & 0x1f
	sid.regs[addr] = data
}

func (sid *SID) Emulate() {
	if sid.regsHistoryIndex < RegisterHistory {
		copy(sid.regsHistory[sid.regsHistoryIndex], sid.regs)
		sid.regsHistoryIndex++
	}
}

func (sid *SID) GetRegsHistory() [][]uint8 {
	return sid.regsHistory
}

func (sid *SID) ResetHistory() uint32 {
	cycle := sid.regsHistoryIndex
	sid.regsHistoryIndex = 0
	return cycle
}
