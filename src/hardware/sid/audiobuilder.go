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

// EGState represents the state of an envelope generator in a sound synthesis context.
type EGState int

// AudioBuilder is a structure used for constructing and managing audio processing components and state.
// player is the audio render interface for handling playback operations.
// fragSize specifies the size of audio fragments in samples.
// bufferFrags determines the number of fragments in the buffer.
// bufferSize represents the total buffer size in bytes for storing audio data.
// volume controls the master volume level of the audio system.
// v3Mute indicates whether voice 3 is muted.
// voices holds the configuration and data for each individual voice generator.
// sampleBuf is a buffer for storing sampled voice data.
// sampleInPtr is the current position in sampleBuf where new samples are written.
// soundBuffer is an array for holding audio output data before rendering.
// toOutput tracks the amount of data left to output in the soundBuffer.
// sbPos keeps the current position within the soundBuffer.
// divisor defines the current divider value used for audio processing timing.
// lead is an array for managing lead data used in audio rendering.
// leadPos specifies the current position within the lead array.
// registerToVoice maps registers to corresponding voice instances.
// filters represents audio filters applied during playback.
// divisorTable holds precomputed divisor values for efficient audio processing.
// leadSmooth determines the smoothing level applied to the lead buffer.
// leadHiWater defines the high-water threshold for lead buffer management.
// leadLoWater sets the low-water threshold for lead buffer management.
type AudioBuilder struct {
	player          references.IAudioRender
	fragSize        int                  // samples, not bytes
	bufferFrags     int                  // frags the in buffer
	bufferSize      int                  // bytes, not samples
	volume          uint8                // Master volume
	v3Mute          uint8                // Voice 3 muted
	voices          []*Voice             // Data for 3 voices
	sampleBuf       [SampleBufSize]uint8 // Buffer for sampled voices
	sampleInPtr     int                  // Index in sample_buf for writing
	soundBuffer     []uint32
	toOutput        int
	sbPos           int
	divisor         int
	lead            []int
	leadPos         int
	registerToVoice []*Voice
	filters         *Filters
	divisorTable    *DivisorTable
	leadSmooth      int
	leadHiWater     int
	leadLoWater     int
}

// NewAudioBuilder initializes and returns a new instance of AudioBuilder with the specified parameters.
func NewAudioBuilder(player references.IAudioRender, useFilters bool, fragFreq int, rasters int) *AudioBuilder {
	bufferFrags := fragFreq                  // one frag per frame
	fragSize := SampleFreq / fragFreq        // samples, not bytes
	fragInterval := 1000 / fragFreq          // in milliseconds
	bufferSize := 2 * fragSize * bufferFrags // bytes, not samples
	maxLeadAvg := bufferFrags
	d := &AudioBuilder{
		player:          player,
		divisorTable:    NewDivisorTable(rasters, fragFreq),
		voices:          nil,
		fragSize:        fragSize,
		bufferFrags:     fragFreq,
		bufferSize:      bufferSize,
		registerToVoice: make([]*Voice, RegisterCount),
		filters:         NewFilters(useFilters),
		lead:            make([]int, maxLeadAvg),
		leadSmooth:      LatencyAvg / fragInterval,
		leadHiWater:     LatencyMax / fragInterval,
		leadLoWater:     LatencyMin / fragInterval,
		soundBuffer:     make([]uint32, 2*fragSize),
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
	dr.v3Mute = 0
	for _, voice := range dr.voices {
		voice.Reset()
	}
	dr.filters.Reset()
	dr.sampleInPtr = 0
	for x := range dr.sampleBuf {
		dr.sampleBuf[x] = 0
	}
	dr.toOutput = 0
	dr.divisor = 0
	for x := range dr.soundBuffer {
		dr.soundBuffer[x] = 0
	}
	for x := range dr.lead {
		dr.lead[x] = 0
	}
	dr.leadPos = 0
	dr.sbPos = 0
}

// LoadRegister processes a given register and applies an update to the corresponding voice or system parameter based on the register value.
func (dr *AudioBuilder) LoadRegister(reg uint8, data uint8) {
	switch reg {
	case 0, 7, 14:
		voice := dr.registerToVoice[reg]
		voice.UpdateFreqA(reg)
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
	case 22:
		dr.updateFilterFreq(data)
	case 23:
		dr.updateVoiceFilters(data)
	case 24:
		dr.updateVolume(data)
	}
}

// updateFilterFreq updates the filter frequency in the Filters instance using the given data.
func (dr *AudioBuilder) updateFilterFreq(data uint8) {
	dr.filters.UpdateFreq(data)
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
	dr.voices[0].filter = f1
	dr.voices[1].filter = f2
	dr.voices[2].filter = f3
	dr.filters.UpdateRes(data)
}

// updateVolume adjusts the master volume, toggles voice 3 mute status, and updates the filter type based on input data.
func (dr *AudioBuilder) updateVolume(data uint8) {
	mute := uint8(0)
	if (data & 0x80) != 0 {
		mute = 1
	}
	dr.volume = data & 0xf
	dr.v3Mute = mute
	dr.filters.UpdateType(data)
}

// calcBuffer processes audio data for the given buffer by applying waveform calculations and filters to generate output.
// It iterates through the buffer, calculating the mixed output for each audio voice and summing the results.
// The method uses sample data, master volume, envelopes, and other voice parameters for precise audio mixing.
// Filtered and unfiltered outputs are computed and combined, then written into the provided buffer.
func (dr *AudioBuilder) calcBuffer(buf []uint32) {
	const halfBufSize = SampleBufSize / 2
	sampleCount := (dr.sampleInPtr + halfBufSize) << 16
	count := len(buf)
	count >>= 1 // 16 bit mono output, count is in bytes
	//count >>= 2; // 16 bit stereo output, count is in bytes
	idx := 0
	for ; count >= 0; count, idx = count-1, idx+1 {
		var sumOutputFilter int32 = 0
		// Get current master volume from sample buffer, calculate sampled voices
		masterVolume := dr.sampleBuf[(sampleCount>>16)%SampleBufSize]
		sampleCount += ((0x138 * 50) << 16) / SampleFreq
		sumOutput := _sampleTable[masterVolume] << 8
		for _, v := range dr.voices {
			v.ComputeEnvelopeGenerators()
			envelope := uint16((v.egLevel * uint32(masterVolume)) >> 20)
			if v.test == 0 {
				v.count += v.add
			}
			if (v.sync != 0) && (v.count > 0x1000000) {
				v.modTo.count = 0
			}
			v.count &= 0xffffff
			output := v.ComputeWaveForm()
			if v.filter != 0 {
				sumOutputFilter += int32((output ^ 0x8000) * envelope)
			} else {
				sumOutput += int32((output ^ 0x8000) * envelope)
			}
		}
		sumOutputFilter = dr.filters.Compute(sumOutputFilter)

		buf[idx] = uint32((sumOutput + sumOutputFilter) >> 10)
	}
}

// Update processes audio data for the buffer, updates sample values, and handles fragment writing when thresholds are met.
func (dr *AudioBuilder) Update() {
	dr.sampleBuf[dr.sampleInPtr] = dr.volume
	dr.sampleInPtr = (dr.sampleInPtr + 1) % SampleBufSize
	dr.divisor += SampleFreq
	dr.toOutput += int(dr.divisorTable.GetOut(dr.divisor))
	dr.divisor = int(dr.divisorTable.GetDivisor(dr.divisor))

	// Calculate the sound data only when we have enough to fill the buffer entirely.
	if dr.toOutput >= dr.fragSize {
		dr.toOutput -= dr.fragSize
		dr.write()
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
	dr.lead[dr.leadPos] = leadInFrags
	dr.leadPos++
	if dr.leadPos == dr.leadSmooth {
		dr.leadPos = 0
	}
	// Compute the average lead in frags.
	avgLead := 0
	for i := 0; i < dr.leadSmooth; i++ {
		avgLead += dr.lead[i]
	}
	avgLead /= dr.leadSmooth
	//fmt.Printf("lead = %d, avg = %d\n", leadInFrags, avgLead)
	//If we're getting too far ahead of the audio skip a frag.
	if avgLead > dr.leadHiWater {
		for i := 0; i < dr.leadSmooth; i++ {
			dr.lead[i]--
		}
		//fmt.Printf("Skipping a frag...\n")
		return
	}
	// Calculate one frag
	nSamples := dr.fragSize
	dr.calcBuffer(dr.soundBuffer)
	// If we're getting too far behind the audio add an extra frag.
	if avgLead < dr.leadLoWater {
		for i := 0; i < dr.leadSmooth; i++ {
			dr.lead[i]++
		}
		//fmt.Printf("Adding an extra frag...\n");
		dr.calcBuffer(dr.soundBuffer[dr.fragSize:])
		nSamples += dr.fragSize
	}
	// Write the frags to the player and update out write position.
	currPos := dr.sbPos
	samples := 2 * nSamples
	dr.sbPos = (dr.sbPos + samples) % dr.bufferSize
	dr.player.Write(dr.soundBuffer, currPos, samples)
}
