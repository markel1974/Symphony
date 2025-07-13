package mos6581

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

const (
	freqLO1 = 0
	freqHI1 = 1
	pwLO1   = 2
	pwHI1   = 3
	cr1     = 4 // Control Register - Voice 1
	ad1     = 5 // Attack/Decay - Voice 1
	sr1     = 6 // Sustain/Release - Voice 1
	freqLO2 = 7
	freqHI2 = 8
	pwLO2   = 9
	pwHI2   = 10
	cr2     = 11
	ad2     = 12
	sr2     = 13
	freqLO3 = 14
	freqHI3 = 15
	pwLO3   = 16
	pwHI3   = 17
	cr3     = 18
	ad3     = 19
	sr3     = 20
	fcLO    = 21 // Filter Cutoff Low
	fcHI    = 22
	resFilt = 23 // Resonance and Filter Control
	modeVol = 24 // Mode and Volume
	potX    = 25 // Potentiometer X (read-only)
	potY    = 26 // Potentiometer Y (read-only)
	osc3    = 27 // Oscillator 3 / Random Number (read-only)
	env3    = 28 // Envelope 3 (read-only)
)

const (
	sidVolumeMax       = 15.0    // 0-15
	normalizedIntValue = 32767.0 // interval -1.0, 1.0
	scalingFactor      = 1024.0  // scaling factor (eq: 1 >> 10)
	divisor            = normalizedIntValue * scalingFactor
	inverseDivisor     = 1.0 / divisor
)

// WriteFn defines a function type that processes an 8-bit unsigned integer as input.
type WriteFn func(reg uint8, data uint8)

// ReadFn defines a function type used to read an 8-bit value from a given 16-bit address.
type ReadFn func(reg uint8) uint8

// SID represents a Sound Interface Device, a component used to generate and handle audio synthesis in the system.
type SID struct {
	*component.BaseComponent
	registers                 []uint8
	cfg                       *config.Config
	reflect                   *SidReflect
	player                    references.IAudioRender
	voices                    *Voices
	filters                   *Filters
	writes                    [RegisterCount]WriteFn
	reads                     [RegisterCount]ReadFn
	volume                    uint8   // Master volume
	sampleBuf                 []uint8 // Buffer for sampled voices
	sampleBufIdx              int     // Index in sample_buf for writing
	soundBuffer               []float32
	audioSamplesPerVolumeStep float64
}

// NewSID initializes and returns a new SID component instance with the given parent, factory, label, and instance number.
// It sets up the base component, registers, and reflection interface for the sound interface device.
func NewSID(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *SID {
	s := &SID{
		BaseComponent: component.NewBaseComponent(),
		registers:     make([]uint8, RegisterCount),
		player:        nil,
		cfg:           nil,
	}
	s.BaseComponent.Register(factory, parent, Identifier(), s, references.IdIMos6581(s, label, instance))
	s.reflect = NewSidReflect(s)
	return s
}

// Setup initializes the SID instance by configuring settings, binding configuration changes, and setting up registers.
func (sid *SID) Setup() error {
	sid.cfg = sid.GetFactory().GetConfig()
	sid.cfg.Bind(sid.onConfigChanged)
	return nil
}

// Bind initializes the SID instance with the given socket, fragment frequency, and raster count, returning an error if any.
func (sid *SID) Bind(_ references.IMos6581Socket, fragFreq int /* rasters */, _ int) error {
	fragSize := SampleFreq / fragFreq
	sid.player = sid.GetFactory().GetIAudioRender()
	sid.sampleBuf = make([]uint8, SampleBufSize)
	sid.filters = NewFilters()
	sid.audioSamplesPerVolumeStep = float64(fragSize) / float64(SampleBufHalfSize)
	sid.soundBuffer = make([]float32, fragSize)
	sid.voices = NewVoices()
	sid.writes = sid.createWriteRegister()
	sid.reads = sid.createReadRegister()
	sid.Reset()
	return nil
}

// Connect establishes a connection using the SID configuration and returns an error if it fails.
func (sid *SID) Connect() error {
	return nil
}

// Internal determines if the SID instance is internal, returning true if it is, otherwise false.
func (sid *SID) Internal() bool {
	return false
}

// Emulate executes the emulation process for the SID object, simulating its behavior according to the defined parameters.
func (sid *SID) Emulate() {

}

// EmulationRequired checks whether emulation is required for the specified SID and returns a boolean result.
func (sid *SID) EmulationRequired() bool {
	return false
}

// onConfigChanged handles updates or adjustments to the configuration settings dynamically during runtime.
func (sid *SID) onConfigChanged() {
	//TODO
}

// Reset reinitializes the SID object by clearing its registers, buffers, and resetting all internal components to defaults.
func (sid *SID) Reset() {
	for x := range sid.registers {
		sid.registers[x] = 0
	}
	sid.SetPotX(0xff)
	sid.SetPotY(0xff)

	sid.volume = 0
	sid.voices.Reset()
	sid.filters.Reset()
	sid.sampleBufIdx = 0
	for x := range sid.sampleBuf {
		sid.sampleBuf[x] = 0
	}
	for x := range sid.soundBuffer {
		sid.soundBuffer[x] = 0
	}
	//PADDLE TEST
	//sid.WriteRegister(0xDC00, 0x40) //Control-Port 1 selected
	//sid.WriteRegister(0xD419, 0x7F) //Paddle X value
	//sid.WriteRegister(0xD419, 0x7F) //Paddle Y value
}

// Prepare updates the sample buffer with the current volume and increments the buffer index circularly.
func (sid *SID) Prepare() {
	// sid.volume is updated by WriteRegister when $D418 is written.
	sid.sampleBuf[sid.sampleBufIdx] = sid.volume
	sid.sampleBufIdx = (sid.sampleBufIdx + 1) % SampleBufSize
}

// Update processes the sound buffer and writes updated sound data to the audio player.
func (sid *SID) Update() {
	sid.calcSoundBuffer()
	sid.player.Write(&sid.soundBuffer, len(sid.soundBuffer))
}

// ReadRegister reads the value of a specified SID register identified by the provided address.
func (sid *SID) ReadRegister(addr uint16) uint8 {
	reg := uint8(addr & RegisterSize)
	return sid.reads[reg](reg)
}

// WriteRegister writes a value to a specific register within the range of the SID's addressable space.
func (sid *SID) WriteRegister(addr uint16, data uint8) {
	reg := uint8(addr & RegisterSize)
	sid.registers[reg] = data
	sid.writes[reg](reg, data)
}

// SetPotX sets the value of the POTX (paddle) register to control the position of the X-axis paddle input.
func (sid *SID) SetPotX(pot uint8) {
	// PX7 PX6 PX5 PX4 PX3 PX2 PX1 PX0
	sid.registers[potX] = pot
}

// SetPotY sets the pot value for the Y-axis register of the SID. This value updates the corresponding SID register directly.
func (sid *SID) SetPotY(pot uint8) {
	//PY7 PY6 PY5 PY4 PY3 PY2 PY1 PY0
	sid.registers[potY] = pot
}

// calcSoundBuffer generates audio samples for the current block, applying volume changes and filtering for each voice.
func (sid *SID) calcSoundBuffer() {
	// Calculate the starting index to read from sampleBuf for this block
	currentVolumeBufferReadIdx := (sid.sampleBufIdx - SampleBufHalfSize + SampleBufSize) % SampleBufSize
	// Read the first volume value to be used for initial audio samples
	currentVolumeValue := sid.sampleBuf[currentVolumeBufferReadIdx]
	// Ratio to determine when to move to the next volume value from sampleBuf
	// nextChangeAtAudioSampleIdx is the audio sample index at which we should
	// move to the next volume value from sampleBuf
	nextChangeAtAudioSampleIdx := sid.audioSamplesPerVolumeStep
	// Count how many of the 312 volume values we've already used
	volumeSteps := 0
	for idx := range sid.soundBuffer {
		// Check if it's time to update currentVolumeValue by reading next value from sampleBuf.
		// This happens when the current audio sample index (idx)
		// exceeds or equals a calculated threshold (nextChangeAtAudioSampleIdx).
		// Also ensure we don't exceed available volume updates.
		if float64(idx) >= nextChangeAtAudioSampleIdx && volumeSteps < SampleBufHalfSize-1 {
			volumeSteps++
			// Advance to the next index in ring buffer sampleBuf
			currentVolumeBufferReadIdx = (currentVolumeBufferReadIdx + 1) % SampleBufSize
			currentVolumeValue = sid.sampleBuf[currentVolumeBufferReadIdx]
			// Calculate a threshold for the next volume change
			nextChangeAtAudioSampleIdx += sid.audioSamplesPerVolumeStep
		}
		// Voice Mixing
		sumOutputFiltered, sumOutputNonFiltered := sid.voices.Compute()
		// Filters
		computedFilter := float32(sid.filters.Compute(sumOutputFiltered))
		mixedSignal := float32(sumOutputNonFiltered) + computedFilter
		volumeFactor := float32(currentVolumeValue) / sidVolumeMax
		sid.soundBuffer[idx] = (mixedSignal * volumeFactor) * inverseDivisor
	}
}

// readRegisterDefault reads the value from the specified SID register address, normalized within the valid register range.
func (sid *SID) readDefault(reg uint8) uint8 {
	return sid.registers[reg]
}

// writeDefault is a default write handler for SID registers that performs no operation when invoked.
func (sid *SID) writeDefault(_ uint8, _ uint8) {
}

// createReadRegister initializes and returns an array of ReadFn functions mapped to SID register addresses.
// It sets a default read function for most registers and customized functions for specific addresses like 27 and 28.
func (sid *SID) createReadRegister() [RegisterCount]ReadFn {
	var reads [RegisterCount]ReadFn
	for idx := range reads {
		reads[idx] = sid.readDefault
	}
	reads[osc3] = sid.voices.ReadVoice2Waveform
	reads[env3] = sid.voices.ReadVoice2EgLevel
	return reads
}

// createWriteRegister initializes and returns an array of WriteFn mapped to SID register write operations.
func (sid *SID) createWriteRegister() [RegisterCount]WriteFn {
	var writes [RegisterCount]WriteFn
	for idx := range writes {
		writes[idx] = sid.writeDefault
	}
	writes[freqLO1] = sid.voices.WriteVoice0UpdateFreqA
	writes[freqHI1] = sid.voices.WriteVoice0UpdateFreqB
	writes[pwLO1] = sid.voices.WriteVoice0UpdatePulseWidthA
	writes[pwHI1] = sid.voices.WriteVoice0UpdatePulseWidthB
	writes[cr1] = sid.voices.WriteVoice0UpdateWaveForm
	writes[ad1] = sid.voices.writeVoice0UpdateEnvelopeGenerators
	writes[sr1] = sid.voices.WriteVoice0UpdateSustainLevel
	writes[freqLO2] = sid.voices.WriteVoice1UpdateFreqA
	writes[freqHI2] = sid.voices.WriteVoice1UpdateFreqB
	writes[pwLO2] = sid.voices.WriteVoice1UpdatePulseWidthA
	writes[pwHI2] = sid.voices.WriteVoice1UpdatePulseWidthB
	writes[cr2] = sid.voices.WriteVoice1UpdateWaveForm
	writes[ad2] = sid.voices.WriteVoice1UpdateEnvelopeGenerators
	writes[sr2] = sid.voices.WriteVoice1UpdateSustainLevel
	writes[freqLO3] = sid.voices.WriteVoice2UpdateFreqA
	writes[freqHI3] = sid.voices.WriteVoice2UpdateFreqB
	writes[pwLO3] = sid.voices.WriteVoice2UpdatePulseWidthA
	writes[pwHI3] = sid.voices.WriteVoice2UpdatePulseWidthB
	writes[cr3] = sid.voices.WriteVoice2UpdateWaveForm
	writes[ad3] = sid.voices.WriteVoice2UpdateEnvelopeGenerators
	writes[sr3] = sid.voices.WriteVoice2UpdateSustainLevel
	writes[fcLO] = sid.filters.UpdateFreqLow
	writes[fcHI] = sid.filters.UpdateFreqHigh
	writes[resFilt] = sid.writeFiltersRegister
	writes[modeVol] = sid.writeMasterVolumeAndFilterType
	return writes
}

// writeFiltersRegister configures filter settings for the SID voices based on the provided data value.
func (sid *SID) writeFiltersRegister(_ uint8, data uint8) {
	sid.voices.SetFilters(data)
	sid.filters.UpdateRes(data)
}

// writeMasterVolumeAndFilterType updates master volume, filter type, and mute state based on the given input data.
func (sid *SID) writeMasterVolumeAndFilterType(_ uint8, data uint8) {
	sid.volume = data & 0xf
	sid.voices.SetMuteVoice2(data)
	sid.filters.UpdateType(data)
}
