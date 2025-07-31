package mos6581

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/references"
)

const (
	voiceNumber = 3
)

// Voices represents a collection of multiple Voice instances.
type Voices struct {
	*component.BaseComponent
	reflect     *VoicesReflect
	voiceNumber int
	voices      []*Voice // Data for 3 voices
}

var _voicesJoin = [][]uint8{
	{2, 1},
	{0, 2},
	{1, 0},
}

// NewVoices creates and initializes a Voices instance containing three interconnected Voice objects.
func NewVoices(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *Voices {
	vs := &Voices{
		BaseComponent: component.NewBaseComponent(),
		voiceNumber:   voiceNumber,
	}
	vs.reflect = NewVoicesReflect(vs, factory, parent, "voices", instance, references.IdInternalComponent(label, instance, "Voices"))
	vs.voices = make([]*Voice, voiceNumber)
	for idx := range vs.voices {
		v := NewVoice(vs, factory, label, idx)
		vs.voices[idx] = v
	}
	for idx, v := range vs.voices {
		if idx < len(_voicesJoin) {
			join := _voicesJoin[idx]
			modBy := vs.voices[join[0]]
			modTo := vs.voices[join[1]]
			v.JoinVoice(modBy, modTo)
		}
	}
	return vs
}

// Setup initializes the Voices instance, preparing it for use, and returns an error if the initialization fails.
func (v *Voices) Setup() error {
	return nil
}

// Connect establishes a connection to the specified voices service and returns an error if the connection fails.
func (v *Voices) Connect() error {
	return nil
}

// EmulationRequired checks if voice emulation is necessary, returning false if not required.
func (v *Voices) EmulationRequired() bool {
	return false
}

// Emulate generates a simulated behavior or operation based on the current state of the Voices instance.
func (v *Voices) Emulate() {
}

// Internal checks and returns whether the Voices instance operates in an internal mode.
func (v *Voices) Internal() bool {
	return true
}

// Reset resets the state of all voices in the collection by calling their individual Reset methods.
func (v *Voices) Reset() {
	for _, voice := range v.voices {
		voice.Reset()
	}
}

// Compute calculates and returns the summed output of filtered and non-filtered voices in int32 format.
func (v *Voices) Compute() (int32, int32) {
	sumOutputNonFiltered := int32(0)
	sumOutputFiltered := int32(0)

	for _, voice := range v.voices {
		voice.ComputeEnvelopeGenerators()
		effectiveEnvelope := voice.EgLevel() >> 16 // 8-bit envelope
		if voice.IsMuted() {
			continue
		}
		voice.UpdateCount()
		waveOutput := voice.ComputeWaveForm()
		signedWaveOutput := int32(int16(waveOutput ^ 0x8000))
		voiceContribution := signedWaveOutput * int32(effectiveEnvelope)
		if voice.Filter() != 0 {
			sumOutputFiltered += voiceContribution
		} else {
			sumOutputNonFiltered += voiceContribution
		}
	}
	return sumOutputFiltered, sumOutputNonFiltered
}

// Voice 0

// WriteVoice0UpdateFreqA updates the frequency parameter A for voice 0 using the provided 8-bit data value.
func (v *Voices) WriteVoice0UpdateFreqA(_ uint8, data uint8) {
	v.voices[0].UpdateFreqA(data)
}

// WriteVoice0UpdateFreqB updates the frequency parameter B for voice 0 using the provided data value.
func (v *Voices) WriteVoice0UpdateFreqB(_ uint8, data uint8) {
	v.voices[0].UpdateFreqB(data)
}

// WriteVoice0UpdatePulseWidthA updates the pulse width A parameter for voice 0 using the provided 8-bit data.
func (v *Voices) WriteVoice0UpdatePulseWidthA(_ uint8, data uint8) {
	v.voices[0].UpdatePulseWidthA(data)
}

// WriteVoice0UpdatePulseWidthB updates the pulse width B parameter for voice 0 with the provided data.
func (v *Voices) WriteVoice0UpdatePulseWidthB(_ uint8, data uint8) {
	v.voices[0].UpdatePulseWidthB(data)
}

// WriteVoice0UpdateWaveForm updates the waveform data for voice 0 by delegating the operation to its UpdateWaveForm method.
func (v *Voices) WriteVoice0UpdateWaveForm(_ uint8, data uint8) {
	v.voices[0].UpdateWaveForm(data)
}

// writeVoice0UpdateEnvelopeGenerators updates the envelope generators for voice 0 using the provided data.
func (v *Voices) writeVoice0UpdateEnvelopeGenerators(_ uint8, data uint8) {
	v.voices[0].UpdateEnvelopeGenerators(data)
}

// WriteVoice0UpdateSustainLevel updates the sustain level of the first voice (voice 0) with the given data value.
func (v *Voices) WriteVoice0UpdateSustainLevel(_ uint8, data uint8) {
	v.voices[0].UpdateSustainLevel(data)
}

// Voice 1

// WriteVoice1UpdateFreqA updates the frequency A of voice 1 using the provided 8-bit data value.
func (v *Voices) WriteVoice1UpdateFreqA(_ uint8, data uint8) {
	v.voices[1].UpdateFreqA(data)
}

// WriteVoice1UpdateFreqB updates the frequency parameter B for voice 1 using the provided data value.
func (v *Voices) WriteVoice1UpdateFreqB(_ uint8, data uint8) {
	v.voices[1].UpdateFreqB(data)
}

// WriteVoice1UpdatePulseWidthA updates the pulse width parameter A for voice 1 with the provided data.
func (v *Voices) WriteVoice1UpdatePulseWidthA(_ uint8, data uint8) {
	v.voices[1].UpdatePulseWidthA(data)
}

// WriteVoice1UpdatePulseWidthB updates the pulse width parameter B of voice 1 with the specified data value.
func (v *Voices) WriteVoice1UpdatePulseWidthB(_ uint8, data uint8) {
	v.voices[1].UpdatePulseWidthB(data)
}

// WriteVoice1UpdateWaveForm updates the waveform configuration for voice 1 using the provided data value.
func (v *Voices) WriteVoice1UpdateWaveForm(_ uint8, data uint8) {
	v.voices[1].UpdateWaveForm(data)
}

// WriteVoice1UpdateEnvelopeGenerators updates the envelope generators for voice 1 using the provided data.
func (v *Voices) WriteVoice1UpdateEnvelopeGenerators(_ uint8, data uint8) {
	v.voices[1].UpdateEnvelopeGenerators(data)
}

// WriteVoice1UpdateSustainLevel updates the sustain level of voice 1 using the provided data.
func (v *Voices) WriteVoice1UpdateSustainLevel(_ uint8, data uint8) {
	v.voices[1].UpdateSustainLevel(data)
}

// Voice 2

// WriteVoice2UpdateFreqA updates the frequency parameter A for voice 2 with the provided data.
func (v *Voices) WriteVoice2UpdateFreqA(_ uint8, data uint8) {
	v.voices[2].UpdateFreqA(data)
}

// WriteVoice2UpdateFreqB updates the frequency parameter B for the third voice using the provided data byte.
func (v *Voices) WriteVoice2UpdateFreqB(_ uint8, data uint8) {
	v.voices[2].UpdateFreqB(data)
}

// WriteVoice2UpdatePulseWidthA updates the pulse width A register of voice 2 with the specified data.
func (v *Voices) WriteVoice2UpdatePulseWidthA(_ uint8, data uint8) {
	v.voices[2].UpdatePulseWidthA(data)
}

// WriteVoice2UpdatePulseWidthB updates the pulse width parameter B for voice 2 with the provided data.
func (v *Voices) WriteVoice2UpdatePulseWidthB(_ uint8, data uint8) {
	v.voices[2].UpdatePulseWidthB(data)
}

// WriteVoice2UpdateWaveForm updates the waveform of the second voice using the provided data value.
func (v *Voices) WriteVoice2UpdateWaveForm(_ uint8, data uint8) {
	v.voices[2].UpdateWaveForm(data)
}

// WriteVoice2UpdateEnvelopeGenerators updates the envelope generator state for the third voice with the provided data.
func (v *Voices) WriteVoice2UpdateEnvelopeGenerators(_ uint8, data uint8) {
	v.voices[2].UpdateEnvelopeGenerators(data)
}

// WriteVoice2UpdateSustainLevel updates the sustain level parameter for voice 2 using the provided data.
func (v *Voices) WriteVoice2UpdateSustainLevel(_ uint8, data uint8) {
	v.voices[2].UpdateSustainLevel(data)
}

// ReadVoice2Waveform retrieves the MSB of the current oscillator output (waveform) for voice 2, derived from ComputeWaveForm.
func (v *Voices) ReadVoice2Waveform(_ uint8) uint8 {
	// OSC3 - Oscillator 3 Value ($D41B)
	// Returns the most significant byte (MSB) of the current output
	// of the oscillator (waveform) for voice 2.
	// The ComputeWaveForm() function in voice.go returns an uint16.
	return uint8(v.voices[2].ComputeWaveForm() >> 8)
}

// ReadVoice2EgLevel reads the most significant byte (MSB) of the current envelope generator level for voice 2.
func (v *Voices) ReadVoice2EgLevel(_ uint8) uint8 {
	// ENV3 - Envelope 3 Value ($D41C)
	// Returns the most significant byte (MSB) of the current level
	// of the envelope (Envelope Generator) for voice 2.
	// The EgLevel() function in voice.go returns an uint32 (24-bit value).
	return uint8(v.voices[2].EgLevel() >> 16)
}

// SetFilters sets the filter values for all three voices using the provided parameters f1, f2, and f3 respectively.
func (v *Voices) SetFilters(data uint8) {
	var f1, f2, f3 uint8 = 0, 0, 0
	if (data & 1) != 0 {
		f1 = 1
	}
	if (data & 2) != 0 {
		f2 = 1
	}
	if (data & 4) != 0 {
		f3 = 1
	}
	v.voices[0].SetFilter(f1)
	v.voices[1].SetFilter(f2)
	v.voices[2].SetFilter(f3)
}

// SetFilterVoice0 modifies the filter setting for voice 0 using the specified filter value f0.
func (v *Voices) SetFilterVoice0(f0 uint8) {
	v.voices[0].SetFilter(f0)
}

// SetFilterVoice1 configures the filter state for voice 1 using the provided filter value f1.
func (v *Voices) SetFilterVoice1(f1 uint8) {
	v.voices[1].SetFilter(f1)
}

// SetFilterVoice2 sets the filter state for voice 2 using the specified filter value.
func (v *Voices) SetFilterVoice2(f2 uint8) {
	v.voices[2].SetFilter(f2)
}

// SetMuteVoice2 sets the mute state of voice 2 based on the boolean parameter m.
func (v *Voices) SetMuteVoice2(data uint8) {
	mute := false
	if (data & 0x80) != 0 {
		mute = true
	}
	v.voices[2].SetMute(mute)
}
