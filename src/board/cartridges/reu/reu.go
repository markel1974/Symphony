package reu

import "github.com/markel1974/c64emu/src/board/cartridges/icartridge"

type REU struct {
}

func NewREU() *REU {
	return &REU{}
}

func (reu *REU) Setup(game uint8, exRom uint8) {
	//TODO IMPLEMENT
}

func (reu *REU) ReadRegister(addr uint16) uint8 {
	if (addr & 0xfff0) == 0xdf00 {
		//TODO IMPLEMENT
		//return s.reu.ReadRegister(addr & 0x0f)
	}
	return 0
}

func (reu *REU) WriteRegister(addr uint16, data uint8) {
	if (addr & 0xfff0) == 0xdf00 {
		//TODO IMPLEMENT
		//s.reu.WriteRegister(addr&0x0f, data)
		return
	}
}

func (reu *REU) Write(i icartridge.RomInterval, addr uint16, data uint8) bool {
	if i == icartridge.ROM_HI_2 && addr == 0xff00 {
		reu.ff00Trigger()
	}
	return false
}

func (reu *REU) Read(i icartridge.RomInterval, addr uint16) (uint8, bool) {
	return 0, false
}

func (reu *REU) Detach() {
	//TODO IMPLEMENT
}

func (reu *REU) ff00Trigger() {
	//TODO IMPLEMENT
}
