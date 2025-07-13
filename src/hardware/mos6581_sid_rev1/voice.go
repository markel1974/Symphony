package mos6581

// WaveFormType represents the type of waveform used in a synthesizer voice.
type WaveFormType int

// WaveNone represents no waveform.
// WaveTri represents a triangle waveform.
// WaveSaw represents a sawtooth waveform.
// WaveTriSaw represents a combination of triangle and sawtooth waveforms.
// WaveRect represents a rectangular waveform.
// WaveTriRect represents a combination of triangle and rectangular waveforms.
// WaveSawRect represents a combination of sawtooth and rectangular waveforms.
// WaveTriSawRect represents a combination of triangle, sawtooth, and rectangular waveforms.
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

// DefaultNoiseLFSR represents the default value for a 24-bit Linear Feedback Shift Register (LFSR).
// Max24BitValue is the maximum value that can be represented in 24 bits.
const (
	DefaultNoiseLFSR = 0x7FFFFF
	Max24BitValue    = 0xFFFFFF
)

// EGState represents the state of an envelope generator in a synthesized voice module.
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

// Voice represents a synthesizer voice with properties for waveform, modulation, frequency, and envelope generation.
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

// NewVoice creates and initializes a new Voice instance with the specified number and default parameters.
// It sets up envelope generation, waveform generation, and test waveforms for the voice.
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

// Setup initializes the Voice instance by linking it with modulating Voices modBy and modTo.
func (v *Voice) Setup(modBy *Voice, modTo *Voice) {
	v.modBy = modBy
	v.modTo = modTo
}

// Reset reinitializes all properties of the Voice object to their default states.
func (v *Voice) Reset() {
	v.wave = WaveNone
	v.egState = EgIdle
	v.add = 0
	v.count = 0
	v.pw = 0
	v.freq = 0
	v.sLevel = 0
	v.egLevel = 0
	v.rSub = egLut(0)
	v.dSub = egLut(0)
	v.aAdd = egLut(0)
	v.test = 0
	v.ring = 0
	v.gate = 0
	v.sync = 0
	v.filter = 0
	v.noiseLFSR = DefaultNoiseLFSR
	v.mute = false
}

func (v *Voice) SetMute(m bool) {
	v.mute = m
}

// IsMuted checks if the voice is currently in a muted state and returns true if muted, otherwise false.
func (v *Voice) IsMuted() bool {
	return v.mute
}

// EgLevel returns the current envelope generator level of the Voice as a uint32.
func (v *Voice) EgLevel() uint32 {
	return v.egLevel
}

// SetFilter sets the filter value for the Voice instance.
func (v *Voice) SetFilter(f uint8) {
	v.filter = f
}

// Filter returns the filter value associated with the Voice instance as an unsigned 8-bit integer.
func (v *Voice) Filter() uint8 {
	return v.filter
}

// UpdateFreqA updates the lower 8 bits of the frequency value and recalculates the additive frequency increment.
func (v *Voice) UpdateFreqA(data uint8) {
	v.freq = (v.freq & 0xff00) | uint16(data)
	v.add = uint32(float64(v.freq) * Frequency / SampleFreq)
}

// UpdateFreqB updates the high byte of the frequency value and recalculates the add value based on updated frequency.
func (v *Voice) UpdateFreqB(data uint8) {
	v.freq = (v.freq & 0xff) | (uint16(data) << 8)
	v.add = uint32(float64(v.freq) * Frequency / SampleFreq)
}

// UpdatePulseWidthA updates the lower 8 bits of the pulse width (pw) with the provided 8-bit data value.
func (v *Voice) UpdatePulseWidthA(data uint8) {
	v.pw = (v.pw & 0x0f00) | uint16(data)
}

// UpdatePulseWidthB updates the high nibble of the pulse width value using the provided 8-bit data.
func (v *Voice) UpdatePulseWidthB(data uint8) {
	v.pw = (v.pw & 0xff) | ((uint16(data) & 0xf) << 8)
}

// UpdateWaveForm updates the waveform type and state flags based on the provided data byte.
// It adjusts gate, sync, ring, and test flags and manages the envelope generator's state accordingly.
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

// UpdateEnvelopeGenerators updates the envelope generator parameters based on the provided data.
// The method computes attack and decay values using the egLut function and updates the relevant fields.
func (v *Voice) UpdateEnvelopeGenerators(data uint8) {
	v.aAdd = egLut(data >> 4)
	v.dSub = egLut(data)
}

// UpdateSustainLevel adjusts the sustain level and release sublevel of the voice based on the provided data value.
func (v *Voice) UpdateSustainLevel(data uint8) {
	v.sLevel = (uint32(data) >> 4) * 0x111111
	v.rSub = egLut(data)
}

// UpdateCount updates the count value for a Voice instance based on specific conditions related to test, sync, and max value.
func (v *Voice) UpdateCount() {
	if v.test == 0 {
		v.count += v.add
	}
	if (v.sync != 0) && (v.count > 0x1000000) {
		v.modTo.count = 0
	}
	v.count &= Max24BitValue
}

// ComputeEnvelopeGenerators evaluates and progresses the envelope generators based on the current state or freezes them if TEST is active.
func (v *Voice) ComputeEnvelopeGenerators() {
	if v.test != 0 {
		// If the TEST bit is active, the envelopes are frozen/disabled.
		// No attack, decay or release progression is executed.
		// The current v.egLevel is maintained.
		return
	}
	v.eg[v.egState]()
}

// ComputeWaveForm generates and returns a waveform value based on the current wave index and test mode of the voice.
func (v *Voice) ComputeWaveForm() uint16 {
	idx := v.wave
	if v.test != 0 {
		return v.waveFormTest[idx]()
	}
	return v.waveForm[idx]()
}

// buildEnvelopeGenerator constructs the envelope generator functions and initializes them for attack, decay, release, and idle states.
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
			level := v.dSub >> eGDRShiftLut(v.egLevel)
			v.egLevel -= level
			if v.egLevel <= v.sLevel || v.egLevel > Max24BitValue {
				v.egLevel = v.sLevel
			}
		}
	}
	eg[EgRelease] = func() {
		level := v.rSub >> eGDRShiftLut(v.egLevel)
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

// buildWaveFormTest initializes and returns a slice of functions used to generate different waveform test patterns for a voice.
func (v *Voice) buildWaveFormTest() []func() uint16 {
	waveFormTest := make([]func() uint16, 0xf+1)
	defaultFn := func() uint16 {
		return v.waveTriTest() & v.waveSawTest() & v.waveRectTest()
	}
	for x := range waveFormTest {
		waveFormTest[x] = defaultFn
	}
	waveFormTest[WaveTri] = v.waveTriTest
	waveFormTest[WaveSaw] = v.waveSawTest
	waveFormTest[WaveRect] = v.waveRectTest
	waveFormTest[WaveNoise] = v.waveNoiseTest
	return waveFormTest
}

// buildWaveForm initializes and returns a slice of waveform generator functions for different wave types.
func (v *Voice) buildWaveForm() []func() uint16 {
	waveForm := make([]func() uint16, 0xf+1)
	defaultFn := func() uint16 {
		return 0x8000
	}
	for x := range waveForm {
		waveForm[x] = defaultFn
	}
	waveForm[WaveTri] = v.waveTri
	waveForm[WaveSaw] = v.waveSaw
	waveForm[WaveRect] = v.waveRect
	waveForm[WaveTriSaw] = v.waveTriSaw
	waveForm[WaveTriRect] = v.waveTriRect
	waveForm[WaveSawRect] = v.waveSawRect
	waveForm[WaveTriSawRect] = v.waveTriSawRect
	waveForm[WaveNoise] = v.waveNoise
	return waveForm
}

// waveTri generates a triangular waveform based on the current count and modulation ring settings.
func (v *Voice) waveTri() uint16 {
	if v.ring != 0 {
		count := v.count ^ (v.modBy.count & 0x800000)
		return triLut(count)
	}
	return triLut(v.count)
}

// waveSaw generates a sawtooth waveform based on the current count value and returns its corresponding uint16 value.
func (v *Voice) waveSaw() uint16 {
	return sawLut(v.count)
}

// waveRect generates a rectangular waveform by comparing a 24-bit counter with a 12-bit pulse width threshold.
func (v *Voice) waveRect() uint16 {
	// The pw threshold is 12 bit, the count accumulator is 24 bit.
	if v.count > (uint32(v.pw) << 12) {
		return 0xffff
	}
	return 0
}

// waveTriSaw generates a combination waveform by bitwise ANDing triangle and sawtooth waveforms.
func (v *Voice) waveTriSaw() uint16 {
	return v.waveTri() & v.waveSaw()
}

// waveTriRect combines the outputs of waveTri and waveRect methods using a bitwise AND operation and returns the result.
func (v *Voice) waveTriRect() uint16 {
	return v.waveTri() & v.waveRect()
}

// waveSawRect combines waveSaw and waveRect results with a bitwise AND operation and returns the resulting value.
func (v *Voice) waveSawRect() uint16 {
	return v.waveSaw() & v.waveRect()
}

// waveTriSawRect generates a combined wave by applying a bitwise AND operation on tri, saw, and rect waveforms.
func (v *Voice) waveTriSawRect() uint16 {
	return v.waveTri() & v.waveSaw() & v.waveRect()
}

// waveNoise generates a 16-bit noise value using a 23-bit Linear Feedback Shift Register (LFSR) algorithm.
// It updates the internal LFSR state of the Voice and calculates noise based on XOR feedback of specific bits.
// Returns a 16-bit unsigned integer representing the generated noise level.
func (v *Voice) waveNoise() uint16 {
	// defines the LFSR polynomial
	msb := (v.noiseLFSR >> 22) & 1
	tapBit := (v.noiseLFSR >> 17) & 1
	feedback := msb ^ tapBit
	// performs shift and inserts new bit
	v.noiseLFSR = ((v.noiseLFSR << 1) | feedback) & DefaultNoiseLFSR
	// maps the output
	return uint16(((v.noiseLFSR >> 15) & 0xFF) << 8)
}

// waveTriTest generates a triangular waveform test value for the voice and returns it as a 16-bit unsigned integer.
func (v *Voice) waveTriTest() uint16 {
	return 0xFFF
}

// waveSawTest generates a waveform by shifting and scaling the `count` property and returns the computed value as uint16.
func (v *Voice) waveSawTest() uint16 {
	frozen := (v.count >> 12) << 4
	return uint16(frozen)
}

// waveRectTest calculates the waveform rectangle test value based on the modulation mode and returns a 16-bit result.
func (v *Voice) waveRectTest() uint16 {
	// Forced ring modulation mode
	if (v.modBy.count & 0x800000) != 0 {
		return 0xFFFF
	}
	return 0x0000
}

// waveNoiseTest generates a 16-bit output based on the linear feedback shift register (LFSR) state of the Voice instance.
func (v *Voice) waveNoiseTest() uint16 {
	// Deterministic output mode
	lfsr := v.noiseLFSR | 0x400000
	return uint16(((lfsr>>12)&0xff)<<8 | (lfsr & 0xff))
}
