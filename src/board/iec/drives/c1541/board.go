package c1541

import (
	"github.com/markel1974/c64emu/src/board/iec/drives/c1541/cpu"
	"github.com/markel1974/c64emu/src/board/iec/drives/c1541/mechanics"
	"github.com/markel1974/c64emu/src/board/iec/drives/c1541/ram"
	"github.com/markel1974/c64emu/src/board/iec/drives/c1541/via"
	"github.com/markel1974/c64emu/src/board/iec/virtualdrive"
	"github.com/markel1974/c64emu/src/board/quartz"
	"github.com/markel1974/c64emu/src/config"
)

type Board struct {
	pin          *Pin
	iec          virtualdrive.IIec
	quartz       *quartz.Quartz
	cpu          *cpu.MOS6502
	via1         *via.Via1
	via2         *via.Via2
	ram          *ram.Ram
	deviceNumber uint8
	cfg          *config.Config
}

func New(quartz *quartz.Quartz, iec virtualdrive.IIec, deviceNumber uint8) *Board {
	return &Board{
		iec:          iec,
		quartz:       quartz,
		deviceNumber: deviceNumber,
		pin:          nil,
		via1:         nil,
		via2:         nil,
		cpu:          nil,
		ram:          nil,
		cfg:          nil,
	}
}

func (m *Board) Setup(cfg *config.Config) {
	m.cfg = cfg
	m.pin = NewPin()
	m.cfg.Bind(m.configChanged)
	m.ram = ram.New(0xffff)
	m.cpu = cpu.NewMOS6502()
	job := mechanics.NewJob(m.ram, m.deviceNumber)
	m.via1 = via.NewVia1(m.iec, m.pin, m.deviceNumber)
	m.via2 = via.NewVia2(m.iec, m.pin, job)
	m.ram.Setup()
	m.cpu.Setup(m.pin, m.ram, cfg)
	m.via1.Setup()
	m.via2.Setup()
}

func (m *Board) Reset() {
	//TODO
}

func (m *Board) Emulate() {
	//TODO
}

func (m *Board) Ready() bool {
	//TODO
	return true
}

func (m *Board) GetDeviceNumber() uint8 {
	return m.deviceNumber
}

func (m *Board) AtnStateChanged(b bool, b2 bool) {
	//TODO implement me
	panic("implement me")
}

func (m *Board) BusStateChanged(u uint8) {
	//TODO implement me
	panic("implement me")
}

func (m *Board) configChanged() {

}
