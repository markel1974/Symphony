package mos6581

import "github.com/markel1974/c64emu/src/flag"

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

type Voice struct {
	number  uint8
	wave    WaveFormType // Selected waveform
	egState EGState      // Current state of EG
	modBy   *Voice       // Voice that modulates this one
	modTo   *Voice       // Voice that is modulated by this one
	count   uint32       // Counter for waveform generator, 8.16 fixed
	add     uint32       // Added to counter in every frame
	freq    uint16       // SID frequency value
	pw      uint16       // SID pulse-width value
	aAdd    uint32       // EG parameters
	dSub    uint32       // dSub is the decrement value for the decay phase of the envelope generator.
	sLevel  uint32       // sLevel represents the sustain level of the envelope generator.
	rSub    uint32       // rSub is the decrement value for the release phase of the envelope generator.
	egLevel uint32       // Current EG level, 8.16 fixed
	noise   uint32       // Last noise generator output value
	gate    bool         // EG gate bit
	ring    bool         // Ring modulation bit
	test    bool         // Test bit
	filter  bool         // Flag: Voice filtered
	sync    bool         // The following bit is set for the modulating voices, not for the modulated one (as the SID bits)
	seed    uint32
}

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
		gate:    false,
		ring:    false,
		test:    false,
		filter:  false,
		sync:    false,
		seed:    1,
	}
}

func (v *Voice) Setup(modBy *Voice, modTo *Voice) {
	v.modBy = modBy
	v.modTo = modTo
}

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
	v.test = false
	v.ring = false
	v.gate = false
	v.sync = false
	v.filter = false
}

func (v *Voice) UpdateFreqA(regIdx uint16) {
	v.freq = (v.freq & 0xff00) | regIdx
	v.add = uint32(float64(v.freq) * Frequency / SampleFreq)
}

func (v *Voice) UpdateFreqB(data uint8) {
	v.freq = (v.freq & 0xff) | (uint16(data) << 8)
	v.add = uint32(float64(v.freq) * Frequency / SampleFreq)
}

func (v *Voice) UpdatePulseWidthA(data uint8) {
	v.pw = (v.pw & 0x0f00) | uint16(data)
}

func (v *Voice) UpdatePulseWidthB(data uint8) {
	v.pw = (v.pw & 0xff) | ((uint16(data) & 0xf) << 8)
}

func (v *Voice) UpdateWaveForm(data uint8) {
	v.wave = WaveFormType(data>>4) & 0xf
	if flag.Uint8ToBool(data&1) != v.gate {
		if (data & 1) != 0 {
			// Gate turned on
			v.egState = EgAttack
		} else {
			// Gate turned off
			if v.egState != EgIdle {
				v.egState = EgRelease
			}
		}
	}
	v.gate = flag.Uint8ToBool(data & 1)
	v.modBy.sync = flag.Uint8ToBool(data & 2)
	v.ring = flag.Uint8ToBool(data & 4)
	v.test = flag.Uint8ToBool(data & 8)
	if v.test {
		v.count = 0
	}
}

func (v *Voice) UpdateEnvelopeGenerators(data uint8) {
	v.aAdd = _eGTable[data>>4]
	v.dSub = _eGTable[data&0xf]
}

func (v *Voice) UpdateSustainLevel(data uint8) {
	v.sLevel = (uint32(data) >> 4) * 0x111111
	v.rSub = _eGTable[data&0xf]
}

func (v *Voice) ComputeEnvelopeGenerators() {
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
}

func (v *Voice) ComputeWaveForm() uint16 {
	output := uint16(0)
	switch v.wave {
	case WaveTri:
		if v.ring {
			output = _triTable[(v.count^(v.modBy.count&0x800000))>>11]
		} else {
			output = _triTable[v.count>>11]
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
			v.seed = (v.seed * 1103515245) + 12345
			noise := v.seed >> 16
			v.noise = noise << 8
			output = uint16(v.noise)
			v.count &= 0xfffff
		} else {
			output = uint16(v.noise)
		}
	default:
		output = 0x8000
	}
	return output
}
