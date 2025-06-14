package mos6581

// WaveFormType defines the type of waveform to be generated in a synthesizer, such as Triangle, Saw, or Noise.
type WaveFormType int

// WaveNone represents a waveform type with no signal.
// WaveTri represents a triangular waveform.
// WaveSaw represents a sawtooth waveform.
// WaveTriSaw represents a combination of triangular and sawtooth waveforms.
// WaveRect represents a rectangular waveform.
// WaveTriRect represents a combination of triangular and rectangular waveforms.
// WaveSawRect represents a combination of sawtooth and rectangular waveforms.
// WaveTriSawRect represents a combination of triangular, sawtooth, and rectangular waveforms.
// WaveNoise represents a noise waveform.
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

const (
	DefaultNoiseLFSR = 0x7FFFFF
	Max24BitValue    = 0xFFFFFF
)

// EGState defines the states of an envelope generator in a synthesizer, such as idle, attack, decay, and release.
type EGState int

// EgIdle represents the idle state in the envelope generator.
// EgAttack represents the attack state in the envelope generator.
// EgDecay represents the decay state in the envelope generator.
// EgRelease represents the release state in the envelope generator.
const (
	EgIdle = EGState(iota)
	EgAttack
	EgDecay
	EgRelease
)

// Voice represents a synthesizer voice with state and parameters for waveform generation and modulation.
type Voice struct {
	number       uint8        // number represents the numerical identifier of the voice in the synthesizer.
	wave         WaveFormType // Selected waveform
	egState      EGState      // Current state of EG
	modBy        *Voice       // Voice that modulates this one
	modTo        *Voice       // Voice that is modulated by this one
	count        uint32       // Counter for waveform generator, 8.16 fixed
	add          uint32       // Added to the counter in every frame
	freq         uint16       // SID frequency value
	pw           uint16       // SID pulse-width value
	aAdd         uint32       // EG parameters
	dSub         uint32       // dSub is the decrement value for the decay phase of the envelope generator.
	sLevel       uint32       // sLevel represents the sustain level of the envelope generator.
	rSub         uint32       // rSub is the decrement value for the release phase of the envelope generator.
	egLevel      uint32       // Current EG level, 8.16 fixed
	gate         uint8        // EG gate bit
	ring         uint8        // Ring modulation bit
	test         uint8        // Test bit
	filter       uint8        // Flag: Voice filtered
	sync         uint8        // The following bit is set for the modulating voices, not for the modulated one (as the SID bits)
	noiseLFSR    uint32
	mute         bool
	waveForm     []func() uint16
	waveFormTest []func() uint16
	eg           []func()
}

// NewVoice creates and initializes a new Voice instance with the specified number. Returns a pointer to the Voice object.
func NewVoice(number uint8) *Voice {
	v := &Voice{
		number:    number,
		wave:      0,
		egState:   0,
		modBy:     nil,
		modTo:     nil,
		count:     0,
		add:       0,
		freq:      0,
		pw:        0,
		aAdd:      0,
		dSub:      0,
		sLevel:    0,
		rSub:      0,
		egLevel:   0,
		gate:      0,
		ring:      0,
		test:      0,
		filter:    0,
		sync:      0,
		noiseLFSR: DefaultNoiseLFSR,
		mute:      false,
	}
	v.eg = v.buildEnvelopeGenerator()
	v.waveForm = v.buildWaveForm()
	v.waveFormTest = v.buildWaveFormTest()
	return v
}

// Setup initializes the voice's modulation relationships by assigning modBy and modTo properties.
func (v *Voice) Setup(modBy *Voice, modTo *Voice) {
	v.modBy = modBy
	v.modTo = modTo
}

// Reset reinitializes the Voice instance by setting its state and attributes to their default values.
func (v *Voice) Reset() {
	v.wave = WaveNone
	v.egState = EgIdle
	v.add = 0
	v.count = 0
	v.pw = 0
	v.freq = 0
	v.sLevel = 0
	v.egLevel = 0
	v.rSub = egTable(0)
	v.dSub = egTable(0)
	v.aAdd = egTable(0)
	v.test = 0
	v.ring = 0
	v.gate = 0
	v.sync = 0
	v.filter = 0
	v.noiseLFSR = DefaultNoiseLFSR
	v.mute = false
}

// IsMuted checks if the Voice is currently muted and returns true if it is, otherwise false.
func (v *Voice) IsMuted() bool {
	return v.mute
}

// EgLevel returns the current envelope generator level of the voice as a uint32.
func (v *Voice) EgLevel() uint32 {
	return v.egLevel
}

// SetFilter sets the filter value for the Voice instance, influencing its sound properties.
func (v *Voice) SetFilter(f uint8) {
	v.filter = f
}

// Filter returns the current filter value for the Voice instance.
func (v *Voice) Filter() uint8 {
	return v.filter
}

// UpdateFreqA updates the low byte of the frequency register and recalculates the phase increment for waveform generation.
func (v *Voice) UpdateFreqA(data uint8) {
	v.freq = (v.freq & 0xff00) | uint16(data)
	v.add = uint32(float64(v.freq) * Frequency / SampleFreq)
}

// UpdateFreqB updates the high byte of the frequency and recalculates the phase increment based on `Frequency` and `SampleFreq`.
func (v *Voice) UpdateFreqB(data uint8) {
	v.freq = (v.freq & 0xff) | (uint16(data) << 8)
	v.add = uint32(float64(v.freq) * Frequency / SampleFreq)
}

// UpdatePulseWidthA updates the lower 8 bits of the pulse width while preserving the upper 4 bits.
func (v *Voice) UpdatePulseWidthA(data uint8) {
	v.pw = (v.pw & 0x0f00) | uint16(data)
}

// UpdatePulseWidthB updates the high nibble of the pulse width value with the provided 8-bit data.
func (v *Voice) UpdatePulseWidthB(data uint8) {
	v.pw = (v.pw & 0xff) | ((uint16(data) & 0xf) << 8)
}

// UpdateWaveForm updates the waveform type and controls, including gate, sync, ring, and test flags, based on the given data.
func (v *Voice) UpdateWaveForm(data uint8) {
	v.wave = WaveFormType(data>>4) & 0xf
	gate := uint8(0)
	ring := uint8(0)
	test := uint8(0)
	sync := uint8(0)
	if (data & 1) != 0 {
		gate = 1
	}
	if (data & 2) != 0 {
		sync = 1
	}
	if (data & 4) != 0 {
		ring = 1
	}
	if (data & 8) != 0 {
		test = 1
	}
	if gate != v.gate {
		if gate != 0 {
			v.egState = EgAttack
		} else {
			if v.egState != EgIdle {
				v.egState = EgRelease
			}
		}
	}
	v.gate = gate
	v.modBy.sync = sync
	v.ring = ring
	if test != v.test {
		v.test = test
		if v.test != 0 {
			v.count = 0
			v.noiseLFSR = DefaultNoiseLFSR
		}
	}
}

// UpdateEnvelopeGenerators updates the attack increment and decay decrement values using a lookup table based on input data.
func (v *Voice) UpdateEnvelopeGenerators(data uint8) {
	v.aAdd = egTable(data >> 4)
	v.dSub = egTable(data)
}

// UpdateSustainLevel sets the sustain level and release decrement rate based on the provided data.
func (v *Voice) UpdateSustainLevel(data uint8) {
	v.sLevel = (uint32(data) >> 4) * 0x111111
	v.rSub = egTable(data)
}

// UpdateCount updates the voice's internal counter based on the current add value and test flag.
// Resets the modulator target counter if sync is enabled and count exceeds a threshold.
// Ensures the count remains within a 24-bit range.
func (v *Voice) UpdateCount() {
	if v.test == 0 {
		v.count += v.add
	}
	if (v.sync != 0) && (v.count > 0x1000000) {
		v.modTo.count = 0
	}
	v.count &= Max24BitValue
}

// ComputeEnvelopeGenerators updates the envelope generator state based on the current state and settings of the Voice.
// If the TEST bit is active, the envelope progression is halted, and the current level is maintained.
func (v *Voice) ComputeEnvelopeGenerators() {
	if v.test != 0 {
		// If the TEST bit is active, the envelopes are frozen/disabled.
		// No attack, decay or release progression is executed.
		// The current v.egLevel is maintained.
		return
	}
	v.eg[v.egState]()
}

// ComputeWaveForm generates the current waveform value based on the voice's wave type and test status.
func (v *Voice) ComputeWaveForm() uint16 {
	idx := v.wave
	if v.test != 0 {
		return v.waveFormTest[idx]()
	}
	return v.waveForm[idx]()
}

// buildEnvelopeGenerator initializes and returns envelope generator functions representing different envelope states.
func (v *Voice) buildEnvelopeGenerator() []func() {
	eg := make([]func(), 0xf)
	defaultFn := func() {}
	for x := range eg {
		eg[x] = defaultFn
	}
	eg[EgAttack] = func() {
		v.egLevel += v.aAdd
		if v.egLevel > Max24BitValue {
			v.egLevel = Max24BitValue
			v.egState = EgDecay
		}
	}
	eg[EgDecay] = func() {
		if v.egLevel <= v.sLevel || v.egLevel > Max24BitValue {
			// The condition > egLevelMax handles underflow
			v.egLevel = v.sLevel
		} else {
			level := v.dSub >> eGDRShiftTable(v.egLevel)
			v.egLevel -= level
			if v.egLevel <= v.sLevel || v.egLevel > Max24BitValue {
				v.egLevel = v.sLevel
			}
		}
	}
	eg[EgRelease] = func() {
		level := v.rSub >> eGDRShiftTable(v.egLevel)
		v.egLevel -= level
		if v.egLevel > Max24BitValue {
			// Underflow (become > 0egLevelMax after subtraction)
			v.egLevel = 0
			v.egState = EgIdle
		}
	}
	eg[EgIdle] = func() {
		v.egLevel = 0
	}
	return eg
}

// buildWaveFormTest initializes and returns a slice of waveform functions for testing purposes, including default behavior.
// Each function determines the waveform output based on the voice's current state and attributes, such as count and noiseLFSR.
func (v *Voice) buildWaveFormTest() []func() uint16 {
	waveFormTest := make([]func() uint16, 0xf+1)
	defaultFn := func() uint16 {
		p1 := triTable(v.count)
		p2 := sawRectTable(v.count)
		return p1 & p2
	}
	for x := range waveFormTest {
		waveFormTest[x] = defaultFn
	}
	waveFormTest[WaveTri] = func() uint16 {
		return 0xFFF
	}
	waveFormTest[WaveSaw] = func() uint16 {
		//return uint16(v.count >> 8)
		frozen := (v.count >> 12) << 4
		return uint16(frozen)
	}
	waveFormTest[WaveRect] = func() uint16 {
		// Forced ring modulation mode
		if (v.modBy.count & 0x800000) != 0 {
			return 0xFFFF
		}
		return 0x0000
	}
	waveFormTest[WaveNoise] = func() uint16 {
		// Deterministic output mode
		lfsr := v.noiseLFSR | 0x400000
		return uint16(((lfsr>>12)&0xFF)<<8 | (lfsr & 0xFF))
	}
	return waveFormTest
}

// buildWaveForm initializes and returns an array of waveform generation functions for the Voice instance.
func (v *Voice) buildWaveForm() []func() uint16 {
	waveForm := make([]func() uint16, 0xf+1)
	defaultFn := func() uint16 {
		return 0x8000
	}
	for x := range waveForm {
		waveForm[x] = defaultFn
	}
	waveForm[WaveTri] = func() uint16 {
		if v.ring != 0 {
			count := v.count ^ (v.modBy.count & 0x800000)
			return triTable(count)
		}
		return triTable(v.count)
	}
	waveForm[WaveSaw] = func() uint16 {
		const scaleFactorRegular = 17
		const scaleFactorFaulty = 14
		accum := v.count >> 12
		// 6581 DAC's non-linearity:
		// The idea is that not all bits have the same "weight".
		// The contribution of the first 11 bits (the more "regular" part of the ramp).
		output := (accum & 0x7FF) * scaleFactorRegular
		// The most significant bit (MSB, value 0x800) is the "faulty" one
		// and contributes differently. We apply its contribution separately.
		if (accum & 0x800) != 0 {
			output += 0x800 * scaleFactorFaulty
		}
		return uint16(output >> 4)
	}
	waveForm[WaveRect] = func() uint16 {
		// The pw threshold is 12 bit, the count accumulator is 24 bit.
		// The comparison is (count_24bit > pw_12bit_shl_12)
		if v.count > (uint32(v.pw) << 12) {
			return 0xffff
		}
		return 0
	}
	waveForm[WaveTriSaw] = func() uint16 {
		return triSawTable(v.count)
	}
	waveForm[WaveTriRect] = func() uint16 {
		if v.count > (uint32(v.pw) << 12) {
			return triRectTable(v.count)
		} else {
			// _triRectTable_low[v.count>>16] if a low part exists
			return 0
		}
	}
	waveForm[WaveSawRect] = func() uint16 {
		if v.count > (uint32(v.pw) << 12) {
			return sawRectTable(v.count)
		}
		// _sawRectTable_low[v.count>>16] if a low part exists
		return 0
	}
	waveForm[WaveTriSawRect] = func() uint16 {
		if v.count > (uint32(v.pw) << 12) {
			return triSawRectTable(v.count)
		}
		// _triSawRectTable_low[v.count>>16] if a low part exists
		return 0
	}
	waveForm[WaveNoise] = func() uint16 {
		// Advance the 23-bit LFSR
		msb := (v.noiseLFSR >> 22) & 1
		tapBit := (v.noiseLFSR >> 17) & 1
		feedback := msb ^ tapBit
		v.noiseLFSR = ((v.noiseLFSR << 1) | feedback) & DefaultNoiseLFSR
		return uint16(((v.noiseLFSR >> 15) & 0xFF) << 8)
	}
	return waveForm
}
