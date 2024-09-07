package mos6581

import (
	"github.com/markel1974/c64emu/src/config"
)

type SID struct {
	id               string
	socket           ISocket
	regs             []uint8
	regsHistory      [][]uint8
	regsHistoryIndex uint32
	cfg              *config.Config
	player           *DigitalRender
}

func NewSID(id string) *SID {
	s := &SID{
		id:               id,
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

func (sid *SID) Setup(socket ISocket, cfg *config.Config) {
	sid.socket = socket
	sid.cfg = cfg
	sid.player = NewDigitalRenderer(socket.GetPlayer(), false)
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
	sid.player.Reset()

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
		//sidCounter := s.sid.ResetHistory()
		sid.player.Render(sid.regsHistory[sid.regsHistoryIndex])
		_ = sid.ResetHistory()
	}
	if lastVicCycle {
		if sid.regsHistoryIndex < RegisterHistory {
			copy(sid.regsHistory[sid.regsHistoryIndex], sid.regs)
			sid.regsHistoryIndex++
		}
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

func (sid *SID) GetLastByte() uint8 {
	return 0
}
