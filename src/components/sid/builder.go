package mos6581

import (
	"github.com/markel1974/c64emu/src/flag"
	"log"
	"math"
)

const (
	LatencyMin = 80
	LatencyMax = 120
	LatencyAvg = 280
	MPi        = math.Pi
)

type EGState int

const (
	EgIdle = EGState(iota)
	EgAttack
	EgDecay
	EgRelease
)

const (
	voicesNumber = 3
)

type DivisorTableData struct {
	divisor int32
	toOut   int32
}

type AudioBuilder struct {
	player             IPlayer
	rasters            int
	fragFreq           int // one frag per frame
	fragSize           int // samples, not bytes
	fragInterval       int // in milliseconds
	bufferFrags        int // frags the in buffer
	bufferSize         int // bytes, not samples
	maxLeadAvg         int // lead average count
	useFilters         bool
	triTable           [0x1000 * 2]uint16
	divisorTableData   []*DivisorTableData
	volume             uint8        // Master volume
	v3Mute             bool         // Voice 3 muted
	voices             []*Voice     // Data for 3 voices
	filterType         FilterType   // Filter type
	filterFreq         uint8        // SID filter frequency (upper 8 bits)
	filterRes          uint8        // Filter resonance (0..15)
	filterAmpl         float64      // IIR filter input attenuation;
	d1, d2, g1, g2     float64      // IIR filter coefficients
	xn1, xn2, yn1, yn2 float64      // IIR filter previous input/output signal
	resonanceLP        [256]float64 // shortcut for calc_filter
	resonanceHP        [256]float64
	sampleBuf          [SampleBufSize]uint8 // Buffer for sampled voices
	sampleInPtr        int                  // Index in sample_buf for writing
	soundBuffer        []uint32
	toOutput           int
	sbPos              int
	divisor            int
	lead               []int
	leadPos            int
	seed               uint32
	registerToVoice    []*Voice
}

func NewAudioBuilder(sp IPlayer, useFilters bool, fragFreq int, rasters int) *AudioBuilder {
	bufferFrags := fragFreq                  // one frag per frame
	fragSize := SampleFreq / fragFreq        // samples, not bytes
	fragInterval := 1000 / fragFreq          // in milliseconds
	bufferSize := 2 * fragSize * bufferFrags // bytes, not samples

	d := &AudioBuilder{
		player:           sp,
		useFilters:       useFilters,
		divisorTableData: nil,
		seed:             1,
		voices:           make([]*Voice, voicesNumber),
		rasters:          rasters,
		fragFreq:         fragFreq,
		fragSize:         fragSize,
		fragInterval:     fragInterval,
		bufferFrags:      fragFreq,
		bufferSize:       bufferSize,
		maxLeadAvg:       bufferFrags,
		registerToVoice:  make([]*Voice, RegisterCount),
	}

	for idx := range d.voices {
		d.voices[idx] = NewVoice(uint8(idx))
	}

	for x := range d.registerToVoice {
		vIdx := (x / 7) % voicesNumber
		d.registerToVoice[x] = d.voices[vIdx]
	}

	// Link voices together
	d.voices[0].modBy = d.voices[2]
	d.voices[1].modBy = d.voices[0]
	d.voices[2].modBy = d.voices[1]
	d.voices[0].modTo = d.voices[1]
	d.voices[1].modTo = d.voices[2]
	d.voices[2].modTo = d.voices[0]
	// Calculate triangle table
	for i := uint16(0); i < 0x1000; i++ {
		d.triTable[i] = (i << 4) | (i >> 8)
		d.triTable[0x1fff-i] = (i << 4) | (i >> 8)
	}
	for i := 0; i < 256; i++ {
		d.resonanceLP[i] = calcResonanceLp(float64(i))
		d.resonanceHP[i] = calcResonanceHp(float64(i))
	}
	d.Reset()
	d.createDivisorTable()
	return d
}

func (dr *AudioBuilder) Reset() {
	dr.volume = 0
	dr.v3Mute = false
	for _, voice := range dr.voices {
		voice.Reset()
	}
	dr.filterType = FilterNone
	dr.filterFreq = 0
	dr.filterRes = 0
	dr.filterAmpl = 1.0
	dr.d1 = 0.0
	dr.d2 = 0.0
	dr.g1 = 0.0
	dr.g2 = 0.0
	dr.xn1 = 0.0
	dr.xn2 = 0.0
	dr.yn1 = 0.0
	dr.yn2 = 0.0
	dr.sampleInPtr = 0

	for x := range dr.sampleBuf {
		dr.sampleBuf[x] = 0
	}
	dr.toOutput = 0
	dr.divisor = 0
	dr.soundBuffer = make([]uint32, 2*dr.fragSize)
	dr.lead = make([]int, dr.maxLeadAvg)
	for lIdx := 0; lIdx < dr.maxLeadAvg; lIdx++ {
		dr.lead[lIdx] = 0
	}
	dr.sbPos = 0
	dr.leadPos = 0
}

func (dr *AudioBuilder) loadRegister(regIdx uint16, data uint8) {
	switch regIdx {
	case 0, 7, 14:
		voice := dr.registerToVoice[regIdx]
		voice.UpdateFreqA(regIdx)
	case 1, 8, 15:
		voice := dr.registerToVoice[regIdx]
		voice.UpdateFreqB(data)
	case 2, 9, 16:
		voice := dr.registerToVoice[regIdx]
		voice.UpdatePulseWidthA(data)
	case 3, 10, 17:
		voice := dr.registerToVoice[regIdx]
		voice.UpdatePulseWidthB(data)
	case 4, 11, 18:
		voice := dr.registerToVoice[regIdx]
		voice.UpdateWaveForm(data)
	case 5, 12, 19:
		voice := dr.registerToVoice[regIdx]
		voice.UpdateEG(data)
	case 6, 13, 20:
		voice := dr.registerToVoice[regIdx]
		voice.UpdateSustainLevel(data)
	case 22:
		dr.updateFilterFreq(data)
	case 23:
		dr.updateVoiceFilters(data)
	case 24:
		dr.updateVolume(data)
	}
}

func (dr *AudioBuilder) updateFilterFreq(data uint8) {
	if data != dr.filterFreq {
		dr.filterFreq = data
		if dr.useFilters {
			dr.calcFilter()
		}
	}
}

func (dr *AudioBuilder) updateVoiceFilters(data uint8) {
	dr.voices[0].filter = flag.Uint8ToBool(data & 1)
	dr.voices[1].filter = flag.Uint8ToBool(data & 2)
	dr.voices[2].filter = flag.Uint8ToBool(data & 4)
	if (data >> 4) != dr.filterRes {
		dr.filterRes = data >> 4
		if dr.useFilters {
			dr.calcFilter()
		}
	}
}

func (dr *AudioBuilder) updateVolume(data uint8) {
	dr.volume = data & 0xf
	dr.v3Mute = flag.Uint8ToBool(data & 0x80)
	if FilterType((data>>4)&7) != dr.filterType {
		dr.filterType = FilterType((data >> 4) & 7)
		dr.xn1 = 0.0
		dr.xn2 = 0.0
		dr.yn1 = 0.0
		dr.yn2 = 0.0
		if dr.useFilters {
			dr.calcFilter()
		}
	}
}

func (dr *AudioBuilder) calcFilter() {
	var fr float64
	// Check for some trivial cases
	if dr.filterType == FilterAll {
		dr.d1 = 0.0
		dr.d2 = 0.0
		dr.g1 = 0.0
		dr.g2 = 0.0
		dr.filterAmpl = 1.0
		return
	} else if dr.filterType == FilterNone {
		dr.d1 = 0.0
		dr.d2 = 0.0
		dr.g1 = 0.0
		dr.g2 = 0.0
		dr.filterAmpl = 0.0
		return
	}

	// Calculate resonance frequency
	if dr.filterType == FilterLp || dr.filterType == FilterLpBp {
		fr = dr.resonanceLP[dr.filterFreq]
	} else {
		fr = dr.resonanceHP[dr.filterFreq]
	}
	// Limit to <1/2 sample frequency, avoid div by 0 in case FilterBp below
	arg := fr / float64(SampleFreq>>1)
	if arg > 0.99 {
		arg = 0.99
	}
	if arg < 0.01 {
		arg = 0.01
	}
	// Calculate poles (resonance frequency and resonance)
	dr.g2 = 0.55 + 1.2*arg*arg - 1.2*arg + float64(dr.filterRes)*0.0133333333
	dr.g1 = -2.0 * math.Sqrt(dr.g2) * math.Cos(MPi*arg)
	// Increase resonance if LP/HP combined with BP
	if dr.filterType == FilterLpBp || dr.filterType == FilterHpBp {
		dr.g2 += 0.1
	}
	// Stabilize filter
	if math.Abs(dr.g1) >= dr.g2+1.0 {
		if dr.g1 > 0.0 {
			dr.g1 = dr.g2 + 0.99
		} else {
			dr.g1 = -(dr.g2 + 0.99)
		}
	}
	// Calculate roots (filter characteristic) and input attenuation
	switch dr.filterType {
	case FilterLpBp, FilterLp:
		dr.d1 = 2.0
		dr.d2 = 1.0
		dr.filterAmpl = 0.25 * (1.0 + dr.g1 + dr.g2)
	case FilterHpBp, FilterHp:
		dr.d1 = -2.0
		dr.d2 = 1.0
		dr.filterAmpl = 0.25 * (1.0 - dr.g1 + dr.g2)
	case FilterBp:
		dr.d1 = 0.0
		dr.d2 = -1.0
		dr.filterAmpl = 0.25 * (1.0 + dr.g1 + dr.g2) * (1.0 + math.Cos(MPi*arg)) / math.Sin(MPi*arg)
	case FilterNotch:
		dr.d1 = -2.0 * math.Cos(MPi*arg)
		dr.d2 = 1.0
		dr.filterAmpl = 0.25 * (1.0 + dr.g1 + dr.g2) * (1.0 + math.Cos(MPi*arg)) / math.Sin(MPi*arg)
	default:
		log.Printf("SID FILTER NOT IMPLEMENTED %d\n", dr.filterType)
	}
}

func (dr *AudioBuilder) noiseRandom() uint8 {
	dr.seed = dr.seed*1103515245 + 12345
	return (uint8)(dr.seed >> 16)
}

func (dr *AudioBuilder) calcBuffer(buf []uint32) {
	// Get filter coefficients, so the emulator won't change them in the middle of calculations
	cfAmpl := dr.filterAmpl
	cd1 := dr.d1
	cd2 := dr.d2
	cg1 := dr.g1
	cg2 := dr.g2
	// Index in sample_buf for reading, 16.16 fixed
	sampleCount := (dr.sampleInPtr + SampleBufSize/2) << 16
	// count >>= 2; // 16 bit stereo output, count is in bytes
	count := len(buf)
	count >>= 1 // 16 bit mono output, count is in bytes
	idx := 0
	for ; count >= 0; count, idx = count-1, idx+1 {
		var sumOutputFilter int32 = 0
		// Get current master volume from sample buffer, calculate sampled voices
		masterVolume := dr.sampleBuf[(sampleCount>>16)%SampleBufSize]
		sampleCount += ((0x138 * 50) << 16) / SampleFreq
		sumOutput := _sampleTable[masterVolume] << 8
		// Loop for all three voices
		for j := 0; j < 3; j++ {
			v := dr.voices[j]
			// Envelope generators
			var envelope uint16
			switch v.egState {
			case EgAttack:
				v.egLevel += v.aAdd
				if v.egLevel > 0xffffff {
					v.egLevel = 0xffffff
					v.egState = EgDecay
				}
			case EgDecay:
				if v.egLevel <= v.sLevel || v.egLevel > 0xffffff {
					v.egLevel = v.sLevel
				} else {
					v.egLevel -= v.dSub >> _eGDRShiftTable[v.egLevel>>16]
					if v.egLevel <= v.sLevel || v.egLevel > 0xffffff {
						v.egLevel = v.sLevel
					}
				}
			case EgRelease:
				v.egLevel -= v.rSub >> _eGDRShiftTable[v.egLevel>>16]
				if v.egLevel > 0xffffff {
					v.egLevel = 0
					v.egState = EgIdle
				}
			case EgIdle:
				v.egLevel = 0
			}
			envelope = uint16((v.egLevel * uint32(masterVolume)) >> 20)
			// Waveform generator
			var output uint16
			if !v.test {
				v.count += v.add
			}
			if v.sync && (v.count > 0x1000000) {
				v.modTo.count = 0
			}
			v.count &= 0xffffff
			switch v.wave {
			case WaveTri:
				if v.ring {
					output = dr.triTable[(v.count^(v.modBy.count&0x800000))>>11]
				} else {
					output = dr.triTable[v.count>>11]
				}
			case WaveSaw:
				output = uint16(v.count >> 8)
			case WaveRect:
				if v.count > uint32(v.pw<<12) {
					output = 0xffff
				} else {
					output = 0
				}
			case WaveTriSaw:
				output = _triSawTable[v.count>>16]
			case WaveTriRect:
				if v.count > uint32(v.pw<<12) {
					output = _triRectTable[v.count>>16]
				} else {
					output = 0
				}
			case WaveSawRect:
				if v.count > uint32(v.pw<<12) {
					output = _sawRectTable[v.count>>16]
				} else {
					output = 0
				}
			case WaveTriSawRect:
				if v.count > uint32(v.pw<<12) {
					output = _triSawRectTable[v.count>>16]
				} else {
					output = 0
				}
			case WaveNoise:
				if v.count > 0x100000 {
					v.noise = uint32(dr.noiseRandom()) << 8
					output = uint16(v.noise)
					v.count &= 0xfffff
				} else {
					output = uint16(v.noise)
				}
			default:
				output = 0x8000
			}
			if v.filter {
				sumOutputFilter += int32((output ^ 0x8000) * envelope)
			} else {
				sumOutput += int32((output ^ 0x8000) * envelope)
			}
		}
		// Filter
		if dr.useFilters {
			xn := float64(sumOutputFilter) * cfAmpl
			yn := xn + cd1*dr.xn1 + cd2*dr.xn2 - cg1*dr.yn1 - cg2*dr.yn2
			dr.yn2 = dr.yn1
			dr.yn1 = yn
			dr.xn2 = dr.xn1
			dr.xn1 = xn
			sumOutputFilter = int32(yn)
		}
		// Write to buffer
		buf[idx] = uint32((sumOutput + sumOutputFilter) >> 10)
	}
}

func (dr *AudioBuilder) Prepare(regs []uint8) {
	for y, r := range regs {
		dr.loadRegister(uint16(y), r)
	}
}

func (dr *AudioBuilder) Update() {
	dr.sampleBuf[dr.sampleInPtr] = dr.volume
	dr.sampleInPtr = (dr.sampleInPtr + 1) % SampleBufSize
	dr.divisor += SampleFreq
	dr.toOutput += int(dr.divisorTableData[dr.divisor].toOut)
	dr.divisor = int(dr.divisorTableData[dr.divisor].divisor)
	// Calculate the sound data only when we have enough to fill the buffer entirely.
	if dr.toOutput >= dr.fragSize {
		dr.toOutput -= dr.fragSize
		dr.write()
	}
}

func (dr *AudioBuilder) write() {
	if dr.player == nil {
		return
	}
	// Convert latency preferences from milliseconds to frags.
	leadSmooth := LatencyAvg / dr.fragInterval
	leadHiWater := LatencyMax / dr.fragInterval
	leadLoWater := LatencyMin / dr.fragInterval
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
	if dr.leadPos == leadSmooth {
		dr.leadPos = 0
	}
	// Compute the average lead in frags.
	avgLead := 0
	for i := 0; i < leadSmooth; i++ {
		avgLead += dr.lead[i]
	}
	avgLead /= leadSmooth
	//fmt.Printf("lead = %d, avg = %d\n", leadInFrags, avgLead)
	//If we're getting too far ahead of the audio skip a frag.
	if avgLead > leadHiWater {
		for i := 0; i < leadSmooth; i++ {
			dr.lead[i]--
		}
		//fmt.Printf("Skipping a frag...\n")
		return
	}
	// Calculate one frag
	nSamples := dr.fragSize
	dr.calcBuffer(dr.soundBuffer)
	// If we're getting too far behind the audio add an extra frag.
	if avgLead < leadLoWater {
		for i := 0; i < leadSmooth; i++ {
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

func (dr *AudioBuilder) createDivisorTable() {
	tmp := int32(dr.rasters * dr.fragFreq)
	dr.divisorTableData = make([]*DivisorTableData, SampleFreq+1)
	for x := int32(0); x <= SampleFreq; x++ {
		dtd := &DivisorTableData{}
		dtd.divisor = x
		dtd.toOut = 0
		for dtd.divisor >= 0 {
			dtd.divisor -= tmp
			dtd.toOut++
		}
		dr.divisorTableData[x] = dtd
	}
}
