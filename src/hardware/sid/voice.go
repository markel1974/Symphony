package mos6581

// WaveFormType represents different waveform types used in a sound synthesis context.
type WaveFormType int

// WaveNone represents the absence of a waveform.
// WaveTri represents a triangular waveform.
// WaveSaw represents a sawtooth waveform.
// WaveTriSaw represents a combined triangular and sawtooth waveform.
// WaveRect represents a rectangular waveform.
// WaveTriRect represents a combined triangular and rectangular waveform.
// WaveSawRect represents a combined sawtooth and rectangular waveform.
// WaveTriSawRect represents a combined triangular, sawtooth, and rectangular waveform.
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

// EGState represents the state of an envelope generator (EG) in sound synthesis, defining phases like attack, decay, sustain, and release.
type EGState int

// EgIdle represents the idle state in the EGState enumeration.
// EgAttack represents the attack state in the EGState enumeration.
// EgDecay represents the decay state in the EGState enumeration.
// EgRelease represents the release state in the EGState enumeration.
const (
	EgIdle = EGState(iota)
	EgAttack
	EgDecay
	EgRelease
)

// Voice represents a sound generator within a synthesizer.
type Voice struct {
	number  uint8        // number represents the numerical identifier of the voice in the synthesizer.
	wave    WaveFormType // Selected waveform
	egState EGState      // Current state of EG
	modBy   *Voice       // Voice that modulates this one
	modTo   *Voice       // Voice that is modulated by this one
	count   uint32       // Counter for waveform generator, 8.16 fixed
	add     uint32       // Added to the counter in every frame
	freq    uint16       // SID frequency value
	pw      uint16       // SID pulse-width value
	aAdd    uint32       // EG parameters
	dSub    uint32       // dSub is the decrement value for the decay phase of the envelope generator.
	sLevel  uint32       // sLevel represents the sustain level of the envelope generator.
	rSub    uint32       // rSub is the decrement value for the release phase of the envelope generator.
	egLevel uint32       // Current EG level, 8.16 fixed
	noise   uint32       // Last noise generator output value
	gate    uint8        // EG gate bit
	ring    uint8        // Ring modulation bit
	test    uint8        // Test bit
	filter  uint8        // Flag: Voice filtered
	sync    uint8        // The following bit is set for the modulating voices, not for the modulated one (as the SID bits)
	//seed    uint32       // seed represents the current state of the random number generator for noise waveform generation.
	noiseLFSR uint32
	mute      bool
}

// NewVoice creates a new Voice instance with provided voice number and initializes its properties to default values.
func NewVoice(number uint8) *Voice {
	return &Voice{
		number:  number,
		wave:    0,
		egState: 0,
		modBy:   nil,
		modTo:   nil,
		count:   0,
		add:     0,
		freq:    0,
		pw:      0,
		aAdd:    0,
		dSub:    0,
		sLevel:  0,
		rSub:    0,
		egLevel: 0,
		noise:   0,
		gate:    0,
		ring:    0,
		test:    0,
		filter:  0,
		sync:    0,
		//seed:    1,
		noiseLFSR: 0x7FFFFF,
		mute:      false,
	}
}

// Setup initializes modulation relationships for the voice by setting the modulating and modulated voices.
func (v *Voice) Setup(modBy *Voice, modTo *Voice) {
	v.modBy = modBy
	v.modTo = modTo
}

// Reset reinitializes the Voice instance to its default state, resetting all parameters and flags to their initial values.
func (v *Voice) Reset() {
	v.wave = WaveNone
	v.egState = EgIdle
	v.add = 0
	v.count = 0
	v.pw = 0
	v.freq = 0
	v.sLevel = 0
	v.egLevel = 0
	v.rSub = _eGTable[0]
	v.dSub = _eGTable[0]
	v.aAdd = _eGTable[0]
	v.test = 0
	v.ring = 0
	v.gate = 0
	v.sync = 0
	v.filter = 0
	v.noiseLFSR = 0x7FFFFF
	v.mute = false
}

func (v *Voice) IsMuted() bool {
	return v.mute
}

func (v *Voice) EgLevel() uint32 {
	return v.egLevel
}

// SetFilter updates the filter flag for the voice to indicate whether it should be filtered, using the given value.
func (v *Voice) SetFilter(f uint8) {
	v.filter = f
}

// Filter returns the filter flag for the voice, indicating whether the voice is filtered.
func (v *Voice) Filter() uint8 {
	return v.filter
}

// UpdateFreqA updates the lower 8 bits of the frequency register and recalculates the increment value for the waveform generator.
func (v *Voice) UpdateFreqA(data uint8) {
	v.freq = (v.freq & 0xff00) | uint16(data)
	v.add = uint32(float64(v.freq) * Frequency / SampleFreq)
}

// UpdateFreqB updates the high byte of the frequency value and recalculates the corresponding increment value for the counter.
func (v *Voice) UpdateFreqB(data uint8) {
	v.freq = (v.freq & 0xff) | (uint16(data) << 8)
	v.add = uint32(float64(v.freq) * Frequency / SampleFreq)
}

// UpdatePulseWidthA updates the lower 8 bits of the pulse width value while preserving the upper 8 bits.
func (v *Voice) UpdatePulseWidthA(data uint8) {
	v.pw = (v.pw & 0x0f00) | uint16(data)
}

// UpdatePulseWidthB updates the high 4 bits of the pulse-width value by masking and shifting the input data.
func (v *Voice) UpdatePulseWidthB(data uint8) {
	v.pw = (v.pw & 0xff) | ((uint16(data) & 0xf) << 8)
}

// UpdateWaveForm updates the waveform type and other voice properties such as gate, sync, ring modulation, and test mode.
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
	//v.test = test
	//if v.test != 0 {
	//	v.count = 0
	//}
	if test != v.test { // Se lo stato del test bit cambia
		v.test = test
		if v.test != 0 {
			v.count = 0            // Resetta l'accumulatore di fase principale
			v.noiseLFSR = 0x7FFFFF // Resetta l'LFSR del rumore al suo stato iniziale
			// Qui andrebbero gestiti anche altri effetti del test bit sugli inviluppi, etc.
		}
	}
}

// UpdateEnvelopeGenerators updates the attack increment and decay decrement rates of the envelope generator using the provided data.
func (v *Voice) UpdateEnvelopeGenerators(data uint8) {
	v.aAdd = _eGTable[(data >> 4)]
	v.dSub = _eGTable[(data & 0xf)]
}

// UpdateSustainLevel adjusts the sustain level and release decrement based on the given data value for the envelope generator.
func (v *Voice) UpdateSustainLevel(data uint8) {
	v.sLevel = (uint32(data) >> 4) * 0x111111
	v.rSub = _eGTable[data&0xf]
}

func (v *Voice) UpdateCount() {
	if v.test == 0 {
		v.count += v.add
	}
	if (v.sync != 0) && (v.count > 0x1000000) {
		v.modTo.count = 0
	}
	v.count &= 0xffffff
}

// ComputeEnvelopeGenerators updates the envelope generator levels and transitions between states based on current values.
// File: voice.go

// ComputeEnvelopeGenerators updates the envelope generator levels and transitions between states based on current values.
func (v *Voice) ComputeEnvelopeGenerators() {
	if v.test != 0 {
		// Se il bit TEST è attivo, gli inviluppi sono congelati/disabilitati.
		// Non viene eseguita alcuna progressione di Attack, Decay o Release.
		// Il livello corrente v.egLevel viene mantenuto.
		return
	}

	switch v.egState {
	case EgAttack:
		v.egLevel += v.aAdd
		if v.egLevel > 0xffffff {
			v.egLevel = 0xffffff
			v.egState = EgDecay
		}
	case EgDecay:
		if v.egLevel <= v.sLevel || v.egLevel > 0xffffff { // La condizione > 0xffffff gestisce l'underflow
			v.egLevel = v.sLevel
		} else {
			v.egLevel -= v.dSub >> _eGDRShiftTable[v.egLevel>>16]
			if v.egLevel <= v.sLevel || v.egLevel > 0xffffff {
				v.egLevel = v.sLevel
			}
		}
	case EgRelease:
		v.egLevel -= v.rSub >> _eGDRShiftTable[v.egLevel>>16]
		if v.egLevel > 0xffffff { // Underflow (diventato > 0xffffff dopo la sottrazione)
			v.egLevel = 0
			v.egState = EgIdle
		}
	case EgIdle:
		v.egLevel = 0
	}
}

// ComputeWaveForm generates and returns the waveform output for the voice based on its current waveform type and settings.
func (v *Voice) ComputeWaveForm() uint16 {
	if v.test != 0 {
		// Comportamento fisico reale del SID quando TEST=1
		switch v.wave {
		case WaveTri:
			// Forza l'onda triangolare a valore massimo
			return 0xFFFF
		case WaveSaw:
			// Congela il dente di sega al valore corrente
			frozen := v.count >> 8
			return uint16(frozen | (frozen << 8))
		case WaveRect:
			// Modalità "ring modulation forzata"
			if v.modBy.count&0x800000 != 0 {
				return 0xFFFF
			}
			return 0x0000
		case WaveNoise:
			// Modalità deterministic output
			lfsr := v.noiseLFSR | 0x400000
			return uint16(((lfsr>>12)&0xFF)<<8 | (lfsr & 0xFF))
		default:
			// Per waveform combinate: usa maschera bitwise
			return _triTable[v.count>>11] & _sawRectTable[v.count>>16]
		}
	}

	// Logica normale di generazione della forma d'onda (se v.test == 0)
	output := uint16(0)
	switch v.wave {
	case WaveTri:
		if v.ring != 0 {
			output = _triTable[(v.count^(v.modBy.count&0x800000))>>11]
		} else {
			output = _triTable[v.count>>11]
		}
	case WaveSaw:
		output = uint16(v.count >> 8)
	case WaveRect:
		// La soglia pw è a 12 bit, l'accumulatore count è a 24 bit.
		// Il confronto è (count_24bit > pw_12bit_shl_12)
		if v.count > (uint32(v.pw) << 12) {
			output = 0xffff
		} else {
			output = 0
		}
	case WaveTriSaw:
		output = _triSawTable[v.count>>16]
	case WaveTriRect:
		if v.count > (uint32(v.pw) << 12) {
			output = _triRectTable[v.count>>16]
		} else {
			output = 0 // O _triRectTable_low[v.count>>16] se esistesse una parte bassa
		}
	case WaveSawRect:
		if v.count > (uint32(v.pw) << 12) {
			output = _sawRectTable[v.count>>16]
		} else {
			output = 0 // O _sawRectTable_low[v.count>>16]
		}
	case WaveTriSawRect:
		if v.count > (uint32(v.pw) << 12) {
			output = _triSawRectTable[v.count>>16]
		} else {
			output = 0 // O _triSawRectTable_low[v.count>>16]
		}
	case WaveNoise:
		// Avanza l'LFSR a 23 bit
		msb := (v.noiseLFSR >> 22) & 1
		tapBit := (v.noiseLFSR >> 17) & 1
		feedback := msb ^ tapBit
		v.noiseLFSR = ((v.noiseLFSR << 1) | feedback) & 0x7FFFFF

		output = uint16(((v.noiseLFSR >> 15) & 0xFF) << 8)
	default:
		output = 0x8000
	}
	return output
}
