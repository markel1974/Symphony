//go:build !js || !wasm

package wasm_render

import (
	"fmt"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

type RenderStub struct {
	err error
}

func New() *RenderStub {
	return &RenderStub{
		err: fmt.Errorf("WASM renderer is not supported on this platform"),
	}
}

func (g *RenderStub) Setup(board references.IBoard, cfg *config.Config) error {
	return g.err
}

func (g *RenderStub) CreateDisplayBuffer(w int, h int) (references.IDisplayBuffer, error) {
	return nil, g.err
}

func (g *RenderStub) Start() error {
	return g.err
}

func (g *RenderStub) VBlank() {
}

type AudioStub struct {
	err error
}

func NewAudio() *AudioStub {
	return &AudioStub{
		err: fmt.Errorf("WASM renderer is not supported on this platform"),
	}
}

func (a *AudioStub) Setup(cfg *config.Config) error {
	return a.err
}

func (a *AudioStub) GetCurrentPosition() int {
	return 0
}

func (a *AudioStub) Write(_ []uint32, pos int, samples int) {
}

func (a *AudioStub) Play() {
}

func (a *AudioStub) Pause() {
}

func (a *AudioStub) Resume() {
}
