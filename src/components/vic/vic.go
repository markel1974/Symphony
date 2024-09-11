package mos6569

import (
	"github.com/markel1974/c64emu/src/config"
)

//https://dustlayer.com/c64-architecture

type cycleData struct {
	fn    func(vic *VIC)
	next  *cycleData
	cycle uint8
}

type VIC struct {
	id              string
	core            *Core
	cfg             *config.Config
	collisions      *Collisions
	sprites         *Sprites
	graphics        *Graphics
	borders         *Borders
	lineStart       int
	drawLine        bool
	vBlankNextCycle bool
	curr            *cycleData
}

func NewVIC(id string) *VIC {
	vic := &VIC{
		id: id,
	}
	return vic
}

func (vic *VIC) Setup(socket ISocket, cfg *config.Config) {
	vic.cfg = cfg
	db := socket.GetDisplayBuffer()
	vic.core = NewCore(socket)
	vic.collisions = NewCollisions(vic.core)
	vic.graphics = NewGraphics(vic.core, vic.collisions, db)
	vic.sprites = NewSprites(vic.core, vic.collisions, db)
	vic.borders = NewBorder(vic.core, db, 13)
	vic.vBlankNextCycle = false
	vic.drawLine = false
	vic.cfg.Bind(vic.configChanged)
	vic.core.Setup()
	vic.graphics.Setup()
	vic.sprites.Setup()
	vic.curr = _pal
}

func (vic *VIC) Reset() {
	vic.core.ready = false
}

func (vic *VIC) GetText() []byte {
	return vic.graphics.GetText()
}

func (vic *VIC) GetLastByte() uint8 {
	return vic.core.lastByte
}

func (vic *VIC) GetBALow() bool {
	return vic.core.GetBALow()
}

func (vic *VIC) GetAECLow() bool {
	return vic.core.GetAECLow()
}

func (vic *VIC) configChanged() {
	//vic.skipFrames = vic.cfg.SkipFrames()
}

func (vic *VIC) ReadRegister(addr uint16) uint8 {
	return vic.core.ReadRegister(addr)
}

func (vic *VIC) WriteRegister(addr uint16, data uint8) {
	vic.core.WriteRegister(addr, data)
}

func (vic *VIC) ChangedVA(va uint8) {
	vic.core.ChangedVA(va)
}

func (vic *VIC) LightPenTrigger() {
	vic.core.LightPenTrigger()
}

func (vic *VIC) Emulate() {
	vic.core.TryAcquireAEC()
	vic.curr.fn(vic)
	vic.curr = vic.curr.next
	vic.core.UpdateRasterX()
}
