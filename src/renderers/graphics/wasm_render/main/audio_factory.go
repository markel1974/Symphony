package main

import (
	"github.com/markel1974/c64emu/src/references"
	"github.com/markel1974/c64emu/src/renderers/graphics/wasm_render"
)

type AudioFactory struct {
}

func NewAudioFactory() *AudioFactory {
	return &AudioFactory{}
}

func (a *AudioFactory) Create(_ string) references.IAudioRender {
	return wasm_render.NewAudio()
}
