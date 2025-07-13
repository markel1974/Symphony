package mos6581

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

const (
	sidVolumeMax       = 15.0    // 0-15
	normalizedIntValue = 32767.0 // interval -1.0, 1.0
	scalingFactor      = 1024.0  // scaling factor (eq: 1 >> 10)
	divisor            = normalizedIntValue * scalingFactor
	inverseDivisor     = 1.0 / divisor
)

// potXRegisterIndex represents the register index for the X-axis potentiometer.
// potYRegisterIndex represents the register index for the Y-axis potentiometer.
const (
	potXRegisterIndex = 25
	potYRegisterIndex = 26
)

// WriteFn defines a function type that processes an 8-bit unsigned integer as input.
type WriteFn func(data uint8)

// ReadFn defines a function type used to read an 8-bit value from a given 16-bit address.
type ReadFn func(addr uint16) uint8

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
	sid.player = sid.GetFactory().GetIAudioRender()
	sid.sampleBuf = make([]uint8, SampleBufSize)
	sid.voices = nil
	fragSize := SampleFreq / fragFreq // samples, not bytes
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
	reg := addr & 0x1f
	return sid.reads[reg](addr)
}

// WriteRegister writes a value to a specific register within the range of the SID's addressable space.
func (sid *SID) WriteRegister(addr uint16, data uint8) {
	reg := addr & 0x1f
	sid.registers[reg] = data
	sid.writes[reg](data)
}

// SetPotX sets the value of the POTX (paddle) register to control the position of the X-axis paddle input.
func (sid *SID) SetPotX(pot uint8) {
	// PX7 PX6 PX5 PX4 PX3 PX2 PX1 PX0
	sid.registers[potXRegisterIndex] = pot
}

// SetPotY sets the pot value for the Y-axis register of the SID. This value updates the corresponding SID register directly.
func (sid *SID) SetPotY(pot uint8) {
	//PY7 PY6 PY5 PY4 PY3 PY2 PY1 PY0
	sid.registers[potYRegisterIndex] = pot
}

// GetVASignal retrieves the last byte of the SID structure and returns it as an unsigned 8-bit integer.
func (sid *SID) GetVASignal() uint8 {
	return 0
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

// createReadRegister initializes and returns an array of ReadFn functions mapped to SID register addresses.
// It sets a default read function for most registers and customized functions for specific addresses like 27 and 28.
func (sid *SID) createReadRegister() [RegisterCount]ReadFn {
	var reads [RegisterCount]ReadFn
	defaultFn := func(addr uint16) uint8 {
		reg := addr & 0x1f
		return sid.registers[reg]
	}
	for idx := range reads {
		reads[idx] = defaultFn
	}
	reads[27] = sid.voices.ReadVoice2Waveform
	reads[28] = sid.voices.ReadVoice2EgLevel
	return reads
}

// createWriteRegister initializes and returns an array of WriteFn mapped to SID register write operations.
func (sid *SID) createWriteRegister() [RegisterCount]WriteFn {
	var writes [RegisterCount]WriteFn
	defaultFn := func(data uint8) {}
	for idx := range writes {
		writes[idx] = defaultFn
	}
	writes[0] = sid.voices.WriteVoice0UpdateFreqA
	writes[1] = sid.voices.WriteVoice0UpdateFreqB
	writes[2] = sid.voices.WriteVoice0UpdatePulseWidthA
	writes[3] = sid.voices.WriteVoice0UpdatePulseWidthB
	writes[4] = sid.voices.WriteVoice0UpdateWaveForm
	writes[5] = sid.voices.writeVoice0UpdateEnvelopeGenerators
	writes[6] = sid.voices.WriteVoice0UpdateSustainLevel
	writes[7] = sid.voices.WriteVoice1UpdateFreqA
	writes[8] = sid.voices.WriteVoice1UpdateFreqB
	writes[9] = sid.voices.WriteVoice1UpdatePulseWidthA
	writes[10] = sid.voices.WriteVoice1UpdatePulseWidthB
	writes[11] = sid.voices.WriteVoice1UpdateWaveForm
	writes[12] = sid.voices.WriteVoice1UpdateEnvelopeGenerators
	writes[13] = sid.voices.WriteVoice1UpdateSustainLevel
	writes[14] = sid.voices.WriteVoice2UpdateFreqA
	writes[15] = sid.voices.WriteVoice2UpdateFreqB
	writes[16] = sid.voices.WriteVoice2UpdatePulseWidthA
	writes[17] = sid.voices.WriteVoice2UpdatePulseWidthB
	writes[18] = sid.voices.WriteVoice2UpdateWaveForm
	writes[19] = sid.voices.WriteVoice2UpdateEnvelopeGenerators
	writes[20] = sid.voices.WriteVoice2UpdateSustainLevel
	writes[21] = sid.filters.UpdateFreqLow
	writes[22] = sid.filters.UpdateFreqHigh
	writes[23] = sid.writeFiltersRegister
	writes[24] = sid.writeMasterVolumeAndFilterType
	return writes
}

// writeFiltersRegister configures filter settings for the SID voices based on the provided data value.
func (sid *SID) writeFiltersRegister(data uint8) {
	sid.voices.SetFilters(data)
	sid.filters.UpdateRes(data)
}

// writeMasterVolumeAndFilterType updates master volume, filter type, and mute state based on the given input data.
func (sid *SID) writeMasterVolumeAndFilterType(data uint8) {
	sid.volume = data & 0xf
	sid.voices.SetMuteVoice2(data)
	sid.filters.UpdateType(data)
}
