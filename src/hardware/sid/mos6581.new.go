package mos6581

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// potXRegisterIndex represents the index of the potentiometer X register.
// potYRegisterIndex represents the index of the potentiometer Y register.
const (
	potXRegisterIndex = 25
	potYRegisterIndex = 26
)

// SID represents the core structure of a Sound Interface Device that processes and renders audio signals.
type SID struct {
	*component.BaseComponent
	registers    []uint8
	cfg          *config.Config
	reflect      *SidReflect
	player       references.IAudioRender
	fragSize     int      // samples, not bytes
	bufferFrags  int      // frags the in buffer
	volume       uint8    // Master volume
	voices       []*Voice // Data for 3 voices
	sampleBuf    []uint8  // Buffer for sampled voices
	sampleBufIdx int      // Index in sample_buf for writing
	soundBuffer  []uint32
	filters      *Filters
	//divisorTable    *DivisorTable
	//divisor         int
}

// NewSID creates and initializes a new SID instance with specified parent component, factory, label, and instance number.
func NewSID(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *SID {
	s := &SID{
		BaseComponent: component.NewBaseComponent(),
		registers:     make([]uint8, RegisterCount),
		player:        nil,
		cfg:           nil,
	}
	s.BaseComponent.Register(factory, parent, Identifier(), s, references.IdISID(s, label, instance))
	s.reflect = NewSidReflect(s)
	return s
}

// Setup initializes the SID configuration and binds it to configuration change events.
func (sid *SID) Setup() error {
	sid.cfg = sid.GetFactory().GetConfig()
	sid.cfg.Bind(sid.onConfigChanged)
	return nil
}

// Bind initializes the SID player, voice structures, and buffer configurations for audio emulation.
func (sid *SID) Bind(_ references.ISIDSocket, fragFreq int, rasters int) error {
	fragSize := SampleFreq / fragFreq // samples, not bytes

	sid.player = sid.GetFactory().GetIAudioRender()
	sid.sampleBuf = make([]uint8, SampleBufSize)
	//sid.divisorTable = NewDivisorTable(rasters, fragFreq)
	sid.voices = nil
	sid.fragSize = fragSize
	sid.bufferFrags = fragFreq
	sid.filters = NewFilters()
	sid.soundBuffer = make([]uint32, 2*fragSize)

	voice0 := NewVoice(0)
	voice1 := NewVoice(1)
	voice2 := NewVoice(2)

	voice0.Setup(voice2, voice1)
	voice1.Setup(voice0, voice2)
	voice2.Setup(voice1, voice0)
	sid.voices = append(sid.voices, voice0, voice1, voice2)

	sid.Reset()

	return nil
}

// Connect establishes a connection or initializes necessary resources for the SID component. Returns an error if unsuccessful.
func (sid *SID) Connect() error {
	return nil
}

// Internal indicates whether the SID instance is operating in an internal mode. Always returns false.
func (sid *SID) Internal() bool {
	return false
}

// Emulate processes the SID chip emulation logic, generating audio output based on the internal state and registers.
func (sid *SID) Emulate() {

}

// EmulationRequired determines if emulation is required for the SID component. Returns false if not required.
func (sid *SID) EmulationRequired() bool {
	return false
}

// SetPotX updates the X-axis potentiometer value by assigning it to the corresponding register in the SID chip.
func (sid *SID) SetPotX(pot uint8) {
	// PX7 PX6 PX5 PX4 PX3 PX2 PX1 PX0
	sid.registers[potXRegisterIndex] = pot
}

// SetPotY sets the potentiometer value for the Y-axis in the SID registers at the potYRegisterIndex.
func (sid *SID) SetPotY(pot uint8) {
	//PY7 PY6 PY5 PY4 PY3 PY2 PY1 PY0
	sid.registers[potYRegisterIndex] = pot
}

// onConfigChanged is triggered when the configuration is modified, allowing the SID instance to handle configuration updates.
func (sid *SID) onConfigChanged() {
	//TODO
}

// Reset initializes all registers, voices, and buffers of the SID instance to default states.
func (sid *SID) Reset() {
	for x := range sid.registers {
		sid.registers[x] = 0
	}
	sid.SetPotX(0xff)
	sid.SetPotY(0xff)

	sid.volume = 0
	for _, voice := range sid.voices {
		voice.Reset()
	}
	sid.filters.Reset()
	sid.sampleBufIdx = 0
	for x := range sid.sampleBuf {
		sid.sampleBuf[x] = 0
	}
	//sid.fragCurrent = 0
	//sid.divisor = 0
	for x := range sid.soundBuffer {
		sid.soundBuffer[x] = 0
	}
	//sid.lead.Reset()
	//sid.sbPos = 0

	//PADDLE TEST
	//sid.WriteRegister(0xDC00, 0x40) //Control-Port 1 selected
	//sid.WriteRegister(0xD419, 0x7F) //Paddle X value
	//sid.WriteRegister(0xD419, 0x7F) //Paddle Y value
}

// ReadRegister reads the value of the SID register at the given address. The address is masked to a valid register range.
//
//	func (sid *SID) ReadRegister(addr uint16) uint8 {
//		reg := addr & 0x1f
//		v := sid.registers[reg]
//		return v
//	}
func (sid *SID) ReadRegister(addr uint16) uint8 {
	reg := addr & 0x1f

	switch reg {
	case 27: // OSC3 - Oscillator 3 Value ($D41B)
		// Restituisce il byte più significativo (MSB) dell'output corrente
		// dell'oscillatore della voce 2.
		// La funzione ComputeWaveForm() in voice.go restituisce un uint16.
		if len(sid.voices) > 2 && sid.voices[2] != nil {
			return uint8(sid.voices[2].ComputeWaveForm() >> 8) //
		}
		return 0 // Fallback se la voce non è inizializzata
	case 28: // ENV3 - Envelope 3 Value ($D41C)
		// Restituisce il byte più significativo (MSB) del livello corrente
		// dell'inviluppo della voce 2.
		// La funzione EgLevel() in voice.go restituisce un uint32 (valore a 24 bit).
		if len(sid.voices) > 2 && sid.voices[2] != nil {
			return uint8(sid.voices[2].EgLevel() >> 16) //
		}
		return 0 // Fallback se la voce non è inizializzata
	default:
		// Per tutti gli altri registri, restituisce il valore memorizzato.
		return sid.registers[reg]
	}
}

// WriteRegister updates a specific SID register at `addr` with the given `data` and triggers related updates for voices or filters.
func (sid *SID) WriteRegister(addr uint16, data uint8) {
	reg := addr & 0x1f
	sid.registers[reg] = data

	switch reg {
	case 0:
		sid.voices[0].UpdateFreqA(data)
	case 1:
		sid.voices[0].UpdateFreqB(data)
	case 2:
		sid.voices[0].UpdatePulseWidthA(data)
	case 3:
		sid.voices[0].UpdatePulseWidthB(data)
	case 4:
		sid.voices[0].UpdateWaveForm(data)
	case 5:
		sid.voices[0].UpdateEnvelopeGenerators(data)
	case 6:
		sid.voices[0].UpdateSustainLevel(data)
	case 7:
		sid.voices[1].UpdateFreqA(data)
	case 8:
		sid.voices[1].UpdateFreqB(data)
	case 9:
		sid.voices[1].UpdatePulseWidthA(data)
	case 10:
		sid.voices[1].UpdatePulseWidthB(data)
	case 11:
		sid.voices[1].UpdateWaveForm(data)
	case 12:
		sid.voices[1].UpdateEnvelopeGenerators(data)
	case 13:
		sid.voices[1].UpdateSustainLevel(data)
	case 14:
		sid.voices[2].UpdateFreqA(data)
	case 15:
		sid.voices[2].UpdateFreqB(data)
	case 16:
		sid.voices[2].UpdatePulseWidthA(data)
	case 17:
		sid.voices[2].UpdatePulseWidthB(data)
	case 18:
		sid.voices[2].UpdateWaveForm(data)
	case 19:
		sid.voices[2].UpdateEnvelopeGenerators(data)
	case 20:
		sid.voices[2].UpdateSustainLevel(data)
	case 21:
		sid.filters.UpdateFreqLow(data)
	case 22:
		// Il registro 22 ($D416) usa solo i 3 bit più bassi per la frequenza
		sid.filters.UpdateFreqHigh(data & 0x07)
	case 23:
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
		sid.voices[0].SetFilter(f1)
		sid.voices[1].SetFilter(f2)
		sid.voices[2].SetFilter(f3)
		sid.filters.UpdateRes(data)
	case 24:
		mute := false //uint8(0)
		if (data & 0x80) != 0 {
			mute = true
		}
		sid.volume = data & 0xf
		sid.voices[2].mute = mute
		sid.filters.UpdateType(data)
	}
}

// Prepare updates the sample buffer with the current volume, increments the buffer index, and calculates the divisor value.
func (sid *SID) Prepare() {
	sid.sampleBuf[sid.sampleBufIdx] = sid.volume
	sid.sampleBufIdx = (sid.sampleBufIdx + 1) % SampleBufSize
	//sid.divisor += SampleFreq
	//sid.divisor = int(sid.divisorTable.GetDivisor(sid.divisor))
}

// Update processes and writes sound data to the audio player buffer using the current state of the SID.
func (sid *SID) Update() {
	sid.calcSoundBuffer()

	//TODO RIMUOVERE 2 * sid.fragSize e pos
	soundBufferSamples := 2 * sid.fragSize
	sid.player.Write(sid.soundBuffer, 0, soundBufferSamples)
}

// GetLastByte retrieves the last byte value processed or stored in the SID instance.
func (sid *SID) GetLastByte() uint8 {
	return 0
}

// calcBuffer generates an audio buffer by combining waveform data, filters, and envelope generators for SID voices.
func (sid *SID) calcSoundBuffer() {
	const halfBufSize = SampleBufSize / 2
	const samples = ((0x138 * 50) << 16) / SampleFreq
	sampleCount := uint32((sid.sampleBufIdx + halfBufSize) << 16)
	count := len(sid.soundBuffer)
	count >>= 1 // 16 bit mono output, count is in bytes
	//count >>= 2; // 16 bit stereo output, count is in bytes
	idx := 0
	for ; count >= 0; count, idx = count-1, idx+1 {

		// Get current master volume from sample buffer, calculate sampled voices
		masterVolume := sid.sampleBuf[(sampleCount>>16)%SampleBufSize]
		sampleCount += samples
		sumOutputFilter := int32(0)

		sumOutput := _sampleTable[masterVolume] << 8
		for _, voice := range sid.voices {
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
		sumOutputFilter = sid.filters.Compute(sumOutputFilter)
		sid.soundBuffer[idx] = uint32((sumOutput + sumOutputFilter) >> 10)
	}
}
