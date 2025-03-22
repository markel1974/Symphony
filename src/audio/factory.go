package audio

import (
	"github.com/markel1974/c64emu/src/audio/oto_player"
	"github.com/markel1974/c64emu/src/references"
	"strings"
)

// Factory is a type responsible for creating instances of IAudioRender implementations based on provided identifiers.
type Factory struct {
}

// NewFactory initializes and returns a new instance of the Factory type.
func NewFactory() *Factory {
	f := &Factory{}
	return f
}

// Create generates an IAudioRender instance based on the provided id. Defaults to "oto" if the id is unrecognized.
func (f *Factory) Create(id string) references.IAudioRender {
	switch strings.ToLower(strings.TrimSpace(id)) {
	case "oto":
		return oto_player.NewAudio()
	default:
		return oto_player.NewAudio()
	}
}
