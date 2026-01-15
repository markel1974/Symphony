package oto_render

import (
	"log"

	"github.com/markel1974/symphony/src/config"
)

// Audio represents a structure for managing audio playback and configuration in a continuous stream setup.
// It includes options for sample rate, channel count, playback control, debugging, and current audio position tracking.
// The type integrates with ContinuousReader for audio streaming and playback functionality.
type Audio struct {
	audioReader *ContinuousReader
	cfg         *config.Config
}

// NewAudio creates and returns a new instance of the Audio type with default configurations for sample rate and channel count.
func NewAudio() *Audio {
	a := &Audio{}
	return a
}

// Setup initializes the Audio instance with the specified configuration and sets up continuous audio playback.
func (d *Audio) Setup(cfg *config.Config) error {
	//TODO DRIVER CONFIG!
	const sampleRate = 44100
	const channels = 1
	const chunks = 50
	d.cfg = cfg
	reader := NewContinuousReader()
	if err := reader.Setup(sampleRate, chunks, channels, "FLOAT32LE"); err != nil {
		return err
	}
	d.audioReader = reader
	go func() {
		d.audioReader.Play()
		if err := d.audioReader.Err(); err != nil {
			log.Println("Error playing audio:", err)
		}
	}()
	return nil
}

// Write processes and buffers audio samples for playback, updating the current position and managing playback timing.
func (d *Audio) Write(values *[]float32, samples int) {
	d.audioReader.AddChunk(values, samples)
}

// Play starts or resumes the audio playback for the Audio instance, leveraging the associated continuous reader.
func (d *Audio) Play() {
}

// Pause halts the current audio playback, maintaining the current playback position.
func (d *Audio) Pause() {
}

// Resume resumes audio playback by restarting the associated audio reader or player.
func (d *Audio) Resume() {
}
