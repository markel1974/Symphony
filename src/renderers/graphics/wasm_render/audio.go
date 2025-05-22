//go:build js && wasm

package wasm_render

import "github.com/markel1974/c64emu/src/config"

type Audio struct {
	cfg *config.Config
	pos int
}

func NewAudio() *Audio {
	return &Audio{
		pos: 0,
	}
}

func (a *Audio) Setup(cfg *config.Config) error {
	a.cfg = cfg
	return nil
}

func (a *Audio) GetCurrentPosition() int {
	return a.pos
}

func (a *Audio) Write(_ []uint32, pos int, samples int) {
	//TODO
	a.pos = pos + samples
	//fmt.Println("AUDIO STREAM ", b, pos, samples)
}

func (a *Audio) Play() {
	//TODO
}

func (a *Audio) Pause() {

}

func (a *Audio) Resume() {

}
