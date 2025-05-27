package mos6581

import (
	"github.com/markel1974/c64emu/src/references"
)

// LatencyMin represents the minimum threshold for network latency in milliseconds.
// LatencyMax represents the maximum threshold for network latency in milliseconds.
// LatencyAvg represents the average threshold for network latency in milliseconds.
const (
	LatencyMin = 80
	LatencyMax = 120
	LatencyAvg = 280
)

// _audioRegisters defines the array of SID audio-related register indices used for audio processing in the SID chip.
//var _audioRegisters = []uint8{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 22, 23, 24}

var _audioRegisters = []uint8{
	0, 1, 2, 3, 4, 5, 6, // Voce 1
	7, 8, 9, 10, 11, 12, 13, // Voce 2
	14, 15, 16, 17, 18, 19, 20, // Voce 3
	21, 22, 23, 24, // Filtro & Globali
}

// EGState represents the state of an envelope generator in a sound synthesis context.
type EGState int

// AudioBuilder is a structure used for constructing and managing audio processing components and state.
// player is the audio render interface for handling playback operations.
// fragSize specifies the size of audio fragments in samples.
// bufferFrags determines the number of fragments in the buffer.
// bufferSize represents the total buffer size in bytes for storing audio data.
// volume controls the master volume level of the audio system.
// voices holds the configuration and data for each individual voice generator.
// sampleBuf is a buffer for storing sampled voice data.
// sampleBufIdx is the current position in sampleBuf where new samples are written.
// soundBuffer is an array for holding audio output data before rendering.
// toOutput tracks the amount of data left to output in the soundBuffer.
// sbPos keeps the current position within the soundBuffer.
// divisor defines the current divider value used for audio processing timing.
// registerToVoice maps registers to corresponding voice instances.
// filters represents audio filters applied during playback.
// divisorTable holds precomputed divisor values for efficient audio processing.
type AudioBuilder struct {
	player          references.IAudioRender
	fragSize        int      // samples, not bytes
	bufferFrags     int      // frags the in buffer
	bufferSize      int      // bytes, not samples
	volume          uint8    // Master volume
	voices          []*Voice // Data for 3 voices
	sampleBuf       []uint8  // Buffer for sampled voices
	sampleBufIdx    int      // Index in sample_buf for writing
	soundBuffer     []uint32
	toOutput        int
	sbPos           int
	divisor         int
	registerToVoice []*Voice
	filters         *Filters
	divisorTable    *DivisorTable
	lead            *Lead
}

// NewAudioBuilder initializes and returns a new instance of AudioBuilder with the specified parameters.
func NewAudioBuilder(player references.IAudioRender, useFilters bool, fragFreq int, rasters int) *AudioBuilder {
	// one frag per frame
	fragSize := SampleFreq / fragFreq     // samples, not bytes
	bufferSize := 2 * fragSize * fragFreq // bytes, not samples

	d := &AudioBuilder{
		player:          player,
		sampleBuf:       make([]uint8, SampleBufSize),
		divisorTable:    NewDivisorTable(rasters, fragFreq),
		voices:          nil,
		fragSize:        fragSize,
		bufferFrags:     fragFreq,
		bufferSize:      bufferSize,
		registerToVoice: make([]*Voice, RegisterCount),
		filters:         NewFilters(),
		soundBuffer:     make([]uint32, 2*fragSize),
		lead:            NewLead(fragFreq),
	}

	voice0 := NewVoice(0)
	voice1 := NewVoice(1)
	voice2 := NewVoice(2)

	voice0.Setup(voice2, voice1)
	voice1.Setup(voice0, voice2)
	voice2.Setup(voice1, voice0)

	d.voices = append(d.voices, voice0, voice1, voice2)

	for x := range d.registerToVoice {
		vIdx := (x / 7) % len(d.voices)
		d.registerToVoice[x] = d.voices[vIdx]
	}

	d.Reset()
	return d
}

// Reset reinitializes the state of the AudioBuilder, clearing buffers, resetting filters, and setting default values.
func (dr *AudioBuilder) Reset() {
	dr.volume = 0
	for _, voice := range dr.voices {
		voice.Reset()
	}
	dr.filters.Reset()
	dr.sampleBufIdx = 0
	for x := range dr.sampleBuf {
		dr.sampleBuf[x] = 0
	}
	dr.toOutput = 0
	dr.divisor = 0
	for x := range dr.soundBuffer {
		dr.soundBuffer[x] = 0
	}
	dr.lead.Reset()
	dr.sbPos = 0
}

// LoadRegister processes audio register data to update voice properties, filter settings, and master volume, generating sound.
func (dr *AudioBuilder) LoadRegister(registers []uint8) {
	for _, reg := range _audioRegisters {
		data := registers[reg]
		switch reg {
		case 0, 7, 14:
			voice := dr.registerToVoice[reg]
			voice.UpdateFreqA(data)
		case 1, 8, 15:
			voice := dr.registerToVoice[reg]
			voice.UpdateFreqB(data)
		case 2, 9, 16:
			voice := dr.registerToVoice[reg]
			voice.UpdatePulseWidthA(data)
		case 3, 10, 17:
			voice := dr.registerToVoice[reg]
			voice.UpdatePulseWidthB(data)
		case 4, 11, 18:
			voice := dr.registerToVoice[reg]
			voice.UpdateWaveForm(data)
		case 5, 12, 19:
			voice := dr.registerToVoice[reg]
			voice.UpdateEnvelopeGenerators(data)
		case 6, 13, 20:
			voice := dr.registerToVoice[reg]
			voice.UpdateSustainLevel(data)
		case 21:
			dr.updateFilterFreqLow(data)
		case 22:
			// Il registro 22 ($D416) usa solo i 3 bit più bassi per la frequenza
			dr.updateFilterFreqHigh(data & 0x07)
		case 23:
			dr.updateVoiceFilters(data)
		case 24:
			dr.updateVolume(data)
		}
	}
}

func (dr *AudioBuilder) Flush() {
	dr.sampleBuf[dr.sampleBufIdx] = dr.volume
	dr.sampleBufIdx = (dr.sampleBufIdx + 1) % SampleBufSize
	dr.divisor += SampleFreq
	dr.toOutput += int(dr.divisorTable.GetOut(dr.divisor))
	dr.divisor = int(dr.divisorTable.GetDivisor(dr.divisor))

	// Calculate the sound data only when we have enough to fill the buffer entirely.
	if dr.toOutput >= dr.fragSize {
		dr.toOutput -= dr.fragSize
		dr.write()
	}
}

// updateFilterFreq updates the filter frequency in the Filters instance using the given data.
func (dr *AudioBuilder) updateFilterFreqLow(data uint8) {
	dr.filters.UpdateFreqLow(data)
}

// updateFilterFreq updates the filter frequency in the Filters instance using the given data.
func (dr *AudioBuilder) updateFilterFreqHigh(data uint8) {
	dr.filters.UpdateFreqHigh(data)
}

// updateVoiceFilters updates the individual filter settings for each voice and the overall filter resonance settings.
func (dr *AudioBuilder) updateVoiceFilters(data uint8) {
	f1 := uint8(0)
	f2 := uint8(0)
	f3 := uint8(0)
	if (data & 1) != 0 {
		f1 = 1
	}
	if (data & 2) != 0 {
		f2 = 1
	}
	if (data & 4) != 0 {
		f3 = 1
	}
	dr.voices[0].SetFilter(f1)
	dr.voices[1].SetFilter(f2)
	dr.voices[2].SetFilter(f3)
	dr.filters.UpdateRes(data)
}

// updateVolume adjusts the master volume, toggles voice 3 mute status, and updates the filter type based on input data.
func (dr *AudioBuilder) updateVolume(data uint8) {
	mute := false
	if (data & 0x80) != 0 {
		mute = true
	}
	dr.volume = data & 0xf
	dr.voices[2].mute = mute
	dr.filters.UpdateType(data)
}

// calcBuffer processes audio data for the given buffer by applying waveform calculations and filters to generate output.
// It iterates through the buffer, calculating the mixed output for each audio voice and summing the results.
// The method uses sample data, master volume, envelopes, and other voice parameters for precise audio mixing.
// Filtered and unfiltered outputs are computed and combined, then written into the provided buffer.
func (dr *AudioBuilder) calcBuffer(buf []uint32, sampleBufIdx int) {
	const halfBufSize = SampleBufSize / 2
	const samples = ((0x138 * 50) << 16) / SampleFreq
	sampleCount := uint32((sampleBufIdx + halfBufSize) << 16)
	count := len(buf)
	count >>= 1 // 16 bit mono output, count is in bytes
	//count >>= 2; // 16 bit stereo output, count is in bytes
	idx := 0
	for ; count >= 0; count, idx = count-1, idx+1 {
		sumOutputFilter := int32(0)
		// Get current master volume from sample buffer, calculate sampled voices
		masterVolume := dr.sampleBuf[(sampleCount>>16)%SampleBufSize]
		sampleCount += samples
		sumOutput := _sampleTable[masterVolume] << 8
		for _, voice := range dr.voices {
			voice.ComputeEnvelopeGenerators()
			envelope := uint16((voice.EgLevel() * uint32(masterVolume)) >> 20)
			if voice.IsMuted() {
				continue
			}
			voice.UpdateCount()
			output := voice.ComputeWaveForm()
			if voice.Filter() != 0 {
				sumOutputFilter += int32(int16(output^0x8000)) * int32(envelope)
			} else {
				sumOutput += int32(int16(output^0x8000)) * int32(envelope)
			}
		}
		sumOutputFilter = dr.filters.Compute(sumOutputFilter)
		buf[idx] = uint32((sumOutput + sumOutputFilter) >> 10)
	}
}

// write processes audio data and manages buffer positioning, adjusting for audio playback timing and synchronization.
func (dr *AudioBuilder) write() {
	// Compute the current lead in frags.
	currentPosition := dr.player.GetCurrentPosition()
	if currentPosition == -1 {
		return
	}
	leadInBytes := (dr.sbPos - currentPosition + dr.bufferSize) % dr.bufferSize
	if leadInBytes >= dr.bufferSize/2 {
		leadInBytes -= dr.bufferSize
	}
	leadInFrags := leadInBytes / 2 * dr.fragSize
	avgLead, ok := dr.lead.Average(leadInFrags)
	if !ok {
		return
	}
	// Calculate one frag
	nSamples := dr.fragSize
	dr.calcBuffer(dr.soundBuffer, dr.sampleBufIdx)
	// If we're getting too far behind the audio add an extra frag.
	if avgLead < dr.lead.GetLoWater() {
		dr.lead.Update()
		//fmt.Printf("Adding an extra frag...\n");
		dr.calcBuffer(dr.soundBuffer[dr.fragSize:], dr.sampleBufIdx)
		nSamples += dr.fragSize
	}
	// Write the frags to the player and update out write position.
	currPos := dr.sbPos
	samples := 2 * nSamples
	dr.sbPos = (dr.sbPos + samples) % dr.bufferSize
	dr.player.Write(dr.soundBuffer, currPos, samples)
}
