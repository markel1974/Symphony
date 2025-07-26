package null_audio_render

import "github.com/markel1974/c64emu/src/config"

type NullAudio struct {
}

func NewAudio() *NullAudio {
	return &NullAudio{}
}

func (d *NullAudio) Setup(_ *config.Config) error {
	return nil
}

// Write processes and buffers audio samples for playback, updating the current position and managing playback timing.
func (d *NullAudio) Write(values *[]float32, samples int) {
}

// Play starts or resumes the audio playback for the Audio instance, leveraging the associated continuous reader.
func (d *NullAudio) Play() {
}

// Pause halts the current audio playback, maintaining the current playback position.
func (d *NullAudio) Pause() {
}

// Resume resumes audio playback by restarting the associated audio reader or player.
func (d *NullAudio) Resume() {
}
