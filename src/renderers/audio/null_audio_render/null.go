package null_audio_render

import "github.com/markel1974/symphony/src/config"

// NullAudio is a no-op audio implementation that fulfills audio-related interfaces without performing any real operations.
type NullAudio struct {
}

// NewAudio creates and returns a new instance of NullAudio with no-op implementations of audio functionalities.
func NewAudio() *NullAudio {
	return &NullAudio{}
}

// Setup initializes the NullAudio instance with the provided configuration and returns an error if initialization fails.
func (d *NullAudio) Setup(_ *config.Config) error {
	return nil
}

// Write sends audio samples to the NullAudio device for processing or discarding without playback.
func (d *NullAudio) Write(_ *[]float32, _ int) {
}

// Play starts audio playback for the NullAudio instance without producing any audible output.
func (d *NullAudio) Play() {
}

// Pause stops audio playback temporarily without releasing resources.
func (d *NullAudio) Pause() {
}

// Resume resumes audio playback in the NullAudio implementation.
func (d *NullAudio) Resume() {
}
