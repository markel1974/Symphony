//go:build js && wasm

package wasm_render

import "github.com/markel1974/c64emu/src/config"

// Audio represents a structure managing audio configurations and playback functionalities.
// Audio maintains an associated configuration and position tracker for processing audio streams.
type Audio struct {
	cfg *config.Config
	pos int
}

// NewAudio initializes and returns a new instance of the Audio struct with the default position set to 0.
func NewAudio() *Audio {
	return &Audio{
		pos: 0,
	}
}

// Setup initializes the Audio instance by assigning the provided configuration. Returns an error if initialization fails.
func (a *Audio) Setup(cfg *config.Config) error {
	a.cfg = cfg
	return nil
}

// GetCurrentPosition retrieves the current playback or stream position of the audio as an integer.
func (a *Audio) GetCurrentPosition() int {
	return a.pos
}

// Write updates the current position of the audio stream by adding the number of samples to the given position.
func (a *Audio) Write(_ []uint32, pos int, samples int) {
	//TODO
	a.pos = pos + samples
	//fmt.Println("AUDIO STREAM ", b, pos, samples)
}

// Play starts audio playback from the current position or last paused state.
func (a *Audio) Play() {
	//TODO
}

// Pause temporarily halts the playback of the audio stream.
func (a *Audio) Pause() {

}

// Resume resumes audio playback from the current position.
func (a *Audio) Resume() {

}
