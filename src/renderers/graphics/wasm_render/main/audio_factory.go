//go:build js && wasm

package main

import (
	"github.com/markel1974/symphony/src/references"
	"github.com/markel1974/symphony/src/renderers/graphics/wasm_render"
)

type AudioFactory struct {
}

func NewAudioFactory() *AudioFactory {
	return &AudioFactory{}
}

func (a *AudioFactory) Create(_ string) references.IAudioRender {
	audio := wasm_render.NewAudio()
	return audio
}
