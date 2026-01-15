package benchmark

import (
	"github.com/markel1974/symphony/src/config"
	"github.com/markel1974/symphony/src/hardware"
	"github.com/markel1974/symphony/src/hardware/mos6569_vic_rev1"
	"github.com/markel1974/symphony/src/renderers/audio/null_audio_render"
	"github.com/markel1974/symphony/src/renderers/graphics/null_graphic_render"
	"log"
	"os"
)

type VicSocket struct {
}

func (v *VicSocket) Cycle() uint64 {
	return 0
}
func (v *VicSocket) BALow(d bool) {
}
func (v *VicSocket) AECLow(d bool) {
}
func (v *VicSocket) VBlank() {
}
func (v *VicSocket) LastCycle() {
}
func (v *VicSocket) ScreenFreq() int {
	return 50
}
func (v *VicSocket) TotalRaster() int {
	return 312
}
func (v *VicSocket) ReadRam(addr uint16) uint8 {
	return 0xff
}
func (v *VicSocket) ReadColorRam(addr uint16) uint8 {
	return 0xff
}
func (v *VicSocket) ReadCharRom(addr uint16) uint8 {
	return 0xff
}
func (v *VicSocket) IRQTrigger() {
}
func (v *VicSocket) IRQClearTrigger() {
}
func (v *VicSocket) RDYLowTrigger(bool) {
}
func (v *VicSocket) AECLowTrigger(bool) {
}
func (v *VicSocket) LastCycleTrigger() {
}
func (v *VicSocket) VBlankTrigger() {
}

func VIC(hz int, interval int, batch int, components int) {
	cfg := config.New()
	audioRender := null_audio_render.NewAudio()
	displayRender := null_graphic_render.NewDisplayBuffer()
	hwFactory := hardware.NewFactory(displayRender, audioRender, cfg)
	var emulate []func()
	for x := 0; x < components; x++ {
		vic := mos6569.NewVIC(nil, hwFactory, "test", 0)
		if err := vic.Setup(); err != nil {
			log.Fatal(err)
		}
		socket := &VicSocket{}
		if err := vic.Bind(socket); err != nil {
			log.Fatal(err)
		}
		emulate = append(emulate, vic.Emulate)
	}
	ticker := NewTicker(hz, interval, batch, func() {
		for _, v := range emulate {
			v()
		}
	})
	ticker.Start()
	os.Exit(1)
}
