package mos6581

type DigitalRenderVoice struct {
	wave    WaveFormType        // Selected waveform
	egState EGState             // Current state of EG
	modBy   *DigitalRenderVoice // Voice that modulates this one
	modTo   *DigitalRenderVoice // Voice that is modulated by this one
	count   uint32              // Counter for waveform generator, 8.16 fixed
	add     uint32              // Added to counter in every frame
	freq    uint16              // SID frequency value
	pw      uint16              // SID pulse-width value
	aAdd    uint32              // EG parameters
	dSub    uint32
	sLevel  uint32
	rSub    uint32
	egLevel uint32 // Current EG level, 8.16 fixed
	noise   uint32 // Last noise generator output value
	gate    bool   // EG gate bit
	ring    bool   // Ring modulation bit
	test    bool   // Test bit
	filter  bool   // Flag: Voice filtered
	sync    bool   // The following bit is set for the modulating voice, not for the modulated one (as the SID bits)
}

func NewDigitalRenderVoice() *DigitalRenderVoice {
	return &DigitalRenderVoice{
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
	}
}
