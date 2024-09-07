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

type WaveFormType int

const (
	WaveNone = WaveFormType(iota)
	WaveTri
	WaveSaw
	WaveTriSaw
	WaveRect
	WaveTriRect
	WaveSawRect
	WaveTriSawRect
	WaveNoise
)

type EGState int

const (
	EgIdle = EGState(iota)
	EgAttack
	EgDecay
	EgRelease
)

type FilterType int

const (
	FilterNone = FilterType(iota)
	FilterLp
	FilterBp
	FilterLpBp
	FilterHp
	FilterNotch
	FilterHpBp
	FilterAll
)

func calcResonanceLp(X float64) float64 {
	v := 227.755 - 1.7635*X - 0.0176385*X*X + 0.00333484*X*X*X - 9.05683e-6*X*X*X*X
	return v
}

func calcResonanceHp(X float64) float64 {
	v := 366.374 - 14.0052*X + 0.603212*X*X - 0.000880196*X*X*X
	return v
}

type DivisorTableData struct {
	divisor int32
	toOut   int32
}

type Builder struct {
	player       IPlayer
	rasters      int
	fragFreq     int // one frag per frame
	fragSize     int // samples, not bytes
	fragInterval int // in milliseconds
	bufferFrags  int // frags the in buffer
	bufferSize   int // bytes, not samples
	maxLeadAvg   int // lead average count
	//regsHistory        [][]uint8
	//regsHistoryIndex   uint32
	useFilters         bool
	triTable           [0x1000 * 2]uint16
	divisorTableData   []*DivisorTableData
	volume             uint8        // Master volume
	v3Mute             bool         // Voice 3 muted
	voice              [3]*Voice    // Data for 3 voices
	filterType         FilterType   // Filter type
	filterFreq         uint8        // SID filter frequency (upper 8 bits)
	filterRes          uint8        // Filter resonance (0..15)
	filterAmpl         float64      // IIR filter input attenuation;
	d1, d2, g1, g2     float64      // IIR filter coefficients
	xn1, xn2, yn1, yn2 float64      // IIR filter previous input/output signal
	resonanceLP        [256]float64 // shortcut for calc_filter
	resonanceHP        [256]float64
	sampleBuf          [SampleBufSize]uint8 // Buffer for sampled voice
	sampleInPtr        int                  // Index in sample_buf for writing
	soundBuffer        []uint32
	toOutput           int
	sbPos              int
	divisor            int
	lead               []int
	leadPos            int
	seed               uint32
}

func NewBuilder(sp IPlayer, useFilters bool, fragFreq int, rasters int) *Builder {
	d := &Builder{
		player:           sp,
		useFilters:       useFilters,
		divisorTableData: nil,
		seed:             1,
		//regsHistory:      make([][]uint8, RegisterHistory),
		//regsHistoryIndex: 0,
	}
	//for x := range d.regsHistory {
	//	d.regsHistory[x] = make([]uint8, RegisterCount)
	//}
	d.rasters = rasters
	d.fragFreq = fragFreq                         // one frag per frame
	d.fragSize = SampleFreq / d.fragFreq          // samples, not bytes
	d.fragInterval = 1000 / d.fragFreq            // in milliseconds
	d.bufferFrags = d.fragFreq                    // frags the in buffer
	d.bufferSize = 2 * d.fragSize * d.bufferFrags // bytes, not samples
	d.maxLeadAvg = d.bufferFrags                  // lead average count

	d.voice[0] = NewVoice()
	d.voice[1] = NewVoice()
	d.voice[2] = NewVoice()
	// Link voices together
	d.voice[0].modBy = d.voice[2]
	d.voice[1].modBy = d.voice[0]
	d.voice[2].modBy = d.voice[1]
	d.voice[0].modTo = d.voice[1]
	d.voice[1].modTo = d.voice[2]
	d.voice[2].modTo = d.voice[0]
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

func (dr *Builder) Reset() {
	//for x := range dr.regsHistory {
	//	for y := range dr.regsHistory[x] {
	//		dr.regsHistory[x][y] = 0
	//	}
	//}
	//dr.regsHistoryIndex = 0

	dr.volume = 0
	dr.v3Mute = false
	for v := 0; v < 3; v++ {
		dr.voice[v].wave = WaveNone
		dr.voice[v].egState = EgIdle
		dr.voice[v].add = 0
		dr.voice[v].count = 0
		dr.voice[v].pw = 0
		dr.voice[v].freq = 0
		dr.voice[v].sLevel = 0
		dr.voice[v].egLevel = 0
		dr.voice[v].rSub = _eGTable[0]
		dr.voice[v].dSub = _eGTable[0]
		dr.voice[v].aAdd = _eGTable[0]
		dr.voice[v].test = false
		dr.voice[v].ring = false
		dr.voice[v].gate = false
		dr.voice[v].sync = false
		dr.voice[v].filter = false
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

func (dr *Builder) loadRegister(addr uint16, data uint8) {
	v := addr / 7 // Voice number
	switch addr {
	case 0, 7, 14:
		dr.voice[v].freq = (dr.voice[v].freq & 0xff00) | addr
		dr.voice[v].add = uint32(float64(dr.voice[v].freq) * Frequency / SampleFreq)
	case 1, 8, 15:
		dr.voice[v].freq = (dr.voice[v].freq & 0xff) | (uint16(data) << 8)
		dr.voice[v].add = uint32(float64(dr.voice[v].freq) * Frequency / SampleFreq)
	case 2, 9, 16:
		dr.voice[v].pw = (dr.voice[v].pw & 0x0f00) | uint16(data)
	case 3, 10, 17:
		dr.voice[v].pw = (dr.voice[v].pw & 0xff) | ((uint16(data) & 0xf) << 8)
	case 4, 11, 18:
		dr.voice[v].wave = WaveFormType(data>>4) & 0xf
		if flag.Uint8ToBool(data&1) != dr.voice[v].gate {
			if (data & 1) != 0 {
				// Gate turned on
				dr.voice[v].egState = EgAttack
			} else {
				// Gate turned off
				if dr.voice[v].egState != EgIdle {
					dr.voice[v].egState = EgRelease
				}
			}
		}
		dr.voice[v].gate = flag.Uint8ToBool(data & 1)
		dr.voice[v].modBy.sync = flag.Uint8ToBool(data & 2)
		dr.voice[v].ring = flag.Uint8ToBool(data & 4)
		dr.voice[v].test = flag.Uint8ToBool(data & 8)
		if dr.voice[v].test {
			dr.voice[v].count = 0
		}
	case 5, 12, 19:
		dr.voice[v].aAdd = _eGTable[data>>4]
		dr.voice[v].dSub = _eGTable[data&0xf]
	case 6, 13, 20:
		dr.voice[v].sLevel = (uint32(data) >> 4) * 0x111111
		dr.voice[v].rSub = _eGTable[data&0xf]
	case 22:
		if data != dr.filterFreq {
			dr.filterFreq = data
			if dr.useFilters {
				dr.calcFilter()
			}
		}
	case 23:
		dr.voice[0].filter = flag.Uint8ToBool(data & 1)
		dr.voice[1].filter = flag.Uint8ToBool(data & 2)
		dr.voice[2].filter = flag.Uint8ToBool(data & 4)
		if (data >> 4) != dr.filterRes {
			dr.filterRes = data >> 4
			if dr.useFilters {
				dr.calcFilter()
			}
		}
	case 24:
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
}

func (dr *Builder) calcFilter() {
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

func (dr *Builder) noiseRandom() uint8 {
	dr.seed = dr.seed*1103515245 + 12345
	return (uint8)(dr.seed >> 16)
}

func (dr *Builder) calcBuffer(buf []uint32) {
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
		// Get current master volume from sample buffer, calculate sampled voice
		masterVolume := dr.sampleBuf[(sampleCount>>16)%SampleBufSize]
		sampleCount += ((0x138 * 50) << 16) / SampleFreq
		sumOutput := _sampleTable[masterVolume] << 8
		// Loop for all three voices
		for j := 0; j < 3; j++ {
			v := dr.voice[j]
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

func (dr *Builder) Prepare(regs []uint8) {
	for y := uint16(0); y < uint16(len(regs)); y++ {
		dr.loadRegister(y, regs[y])
	}

	//if dr.regsHistoryIndex < RegisterHistory {
	//	copy(dr.regsHistory[dr.regsHistoryIndex], regs)
	//	dr.regsHistoryIndex++
	//}
}

func (dr *Builder) Render() {
	//for x := uint32(0); x < dr.regsHistoryIndex; x++ {
	//	for y := uint16(0); y < uint16(len(dr.regsHistory[x])); y++ {
	//		dr.loadRegister(y, dr.regsHistory[x][y])
	//	}
	//}
	//dr.regsHistoryIndex = 0

	dr.sampleBuf[dr.sampleInPtr] = dr.volume
	dr.sampleInPtr = (dr.sampleInPtr + 1) % SampleBufSize
	dr.divisor += SampleFreq
	dr.toOutput += int(dr.divisorTableData[dr.divisor].toOut)
	dr.divisor = int(dr.divisorTableData[dr.divisor].divisor)
	// Calculate the sound data only when we have enough to fill the buffer entirely.
	if dr.toOutput >= dr.fragSize {
		dr.toOutput -= dr.fragSize
		dr.playSound()
	}
}

func (dr *Builder) playSound() {
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
	dr.player.Write(dr.soundBuffer, dr.sbPos, 2*nSamples)
	dr.sbPos = (dr.sbPos + 2*nSamples) % dr.bufferSize
}

func (dr *Builder) Pause() {
	if dr.player == nil {
		return
	}
	dr.player.Pause()
}

func (dr *Builder) Resume() {
	if dr.player == nil {
		return
	}
	dr.player.Resume()
}

func (dr *Builder) createDivisorTable() {
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
