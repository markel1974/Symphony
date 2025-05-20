//go:build !js || !wasm

package wasm_render

import (
	"fmt"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

type RenderStub struct {
	error
}

func New() *RenderStub {
	return &RenderStub{
		error: fmt.Errorf("WASM renderer is not supported on this platform"),
	}
}

func (g *RenderStub) Setup(board references.IBoard, cfg *config.Config) error {
	return g.error
}

func (g *RenderStub) CreateDisplayBuffer(w int, h int) (references.IDisplayBuffer, error) {
	return nil, g.error
}

func (g *RenderStub) Start() error {
	return g.error
}

func (g *RenderStub) VBlank() {
}
