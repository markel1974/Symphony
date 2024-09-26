package mos6581

import (
	"github.com/markel1974/c64emu/src/config"
)

const (
	potXRegisterIndex = 25
	potYRegisterIndex = 26
)

type SID struct {
	id           string
	socket       ISocket
	registers    []uint8
	cfg          *config.Config
	audioBuilder *AudioBuilder
}

func NewSID(id string) *SID {
	s := &SID{
		id:           id,
		socket:       nil,
		registers:    make([]uint8, RegisterCount),
		cfg:          nil,
		audioBuilder: nil,
	}
	return s
}

func (sid *SID) Setup(socket ISocket, cfg *config.Config, fragFreq int, rasters int) {
	sid.socket = socket
	sid.audioBuilder = NewAudioBuilder(sid.socket.GetPlayer(), true, fragFreq, rasters)
	sid.cfg = cfg
	sid.cfg.Bind(sid.onConfigChanged)
}

func (sid *SID) SetPotX(pot uint8) {
	// PX7 PX6 PX5 PX4 PX3 PX2 PX1 PX0
	sid.registers[potXRegisterIndex] = pot
}

func (sid *SID) SetPotY(pot uint8) {
	//PY7 PY6 PY5 PY4 PY3 PY2 PY1 PY0
	sid.registers[potYRegisterIndex] = pot
}

func (sid *SID) onConfigChanged() {
	//TODO
}

func (sid *SID) Reset() {
	for x := range sid.registers {
		sid.registers[x] = 0
	}
	sid.SetPotX(0xff)
	sid.SetPotY(0xff)

	sid.audioBuilder.Reset()

	//PADDLE TEST
	//sid.WriteRegister(0xDC00, 0x40) //Control-Port 1 selected
	//sid.WriteRegister(0xD419, 0x7F) //Paddle X value
	//sid.WriteRegister(0xD419, 0x7F) //Paddle Y value
}

func (sid *SID) ReadRegister(addr uint16) uint8 {
	reg := addr & 0x1f
	v := sid.registers[reg]
	return v
}

func (sid *SID) WriteRegister(addr uint16, data uint8) {
	reg := addr & 0x1f
	sid.registers[reg] = data
}

func (sid *SID) Prepare() {
	for _, x := range _audioRegisters {
		sid.audioBuilder.LoadRegister(x, sid.registers[x])
	}
}

func (sid *SID) Update() {
	sid.audioBuilder.Update()
}

func (sid *SID) GetLastByte() uint8 {
	return 0
}
