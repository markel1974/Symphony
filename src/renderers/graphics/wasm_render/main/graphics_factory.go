//go:build js && wasm

package main

import (
	"github.com/markel1974/symphony/src/references"
	"github.com/markel1974/symphony/src/renderers/graphics/wasm_render"
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
