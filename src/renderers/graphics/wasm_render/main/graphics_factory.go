//go:build js && wasm

package main

import (
	"github.com/markel1974/c64emu/src/references"
	"github.com/markel1974/c64emu/src/renderers/graphics/wasm_render"
)

type GraphicsFactory struct {
}

func NewGraphicsFactory() *GraphicsFactory {
	return &GraphicsFactory{}
}

func (g *GraphicsFactory) Create(_ string) references.IDisplayRender {
	render := wasm_render.NewRender()
	return render
}
