package mos6581

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// potXRegisterIndex is the index used to reference the X-axis potentiometer register.
// potYRegisterIndex is the index used to reference the Y-axis potentiometer register.
const (
	potXRegisterIndex = 25
	potYRegisterIndex = 26
)

// SID represents a sound interface device component with registers, configuration, and audio rendering capabilities.
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

// NewSID initializes a new SID instance, registers it with its factory and parent, and sets the default state.
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

// Setup initializes the SID instance by binding configuration changes to a callback and obtaining configuration settings.
func (sid *SID) Setup() error {
	sid.cfg = sid.GetFactory().GetConfig()
	sid.cfg.Bind(sid.onConfigChanged)
	return nil
}

// Bind initializes the SID instance by setting up voices, sound buffers, filtering, and prepares it for audio rendering.
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

// Connect establishes a connection or initializes components required by the SID device.
func (sid *SID) Connect() error {
	return nil
}

// Internal determines if the SID instance operates in internal mode. Always returns false.
func (sid *SID) Internal() bool {
	return false
}

// Emulate processes and generates audio data for the SID chip emulation, updating internal states and sound buffers.
func (sid *SID) Emulate() {

}

// EmulationRequired checks if emulation is required for the current SID instance and returns a boolean result.
func (sid *SID) EmulationRequired() bool {
	return false
}

// SetPotX sets the X-coordinate potentiometer value in the SID chip by writing to the potX register.
func (sid *SID) SetPotX(pot uint8) {
	// PX7 PX6 PX5 PX4 PX3 PX2 PX1 PX0
	sid.registers[potXRegisterIndex] = pot
}

// SetPotY sets the value of the potentiometer Y register to the given 8-bit value.
func (sid *SID) SetPotY(pot uint8) {
	//PY7 PY6 PY5 PY4 PY3 PY2 PY1 PY0
	sid.registers[potYRegisterIndex] = pot
}

// onConfigChanged is triggered when the configuration settings bound to the SID instance are updated or modified.
func (sid *SID) onConfigChanged() {
	//TODO
}

// Reset resets the internal state of the SID, clearing registers, buffers, and resetting voices and filters.
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

// ReadRegister reads the value of a SID register at the specified address.
// Returns the register data or computed value depending on the address and SID voice state.
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

// WriteRegister writes a value to a specified SID register address and updates the corresponding voice or filter settings.
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

// GetLastByte retrieves the last byte stored in the sample buffer.
func (sid *SID) GetLastByte() uint8 {
	return 0
}

// Prepare updates the sample buffer with the current volume and increments the buffer index cyclically.
func (sid *SID) Prepare() {
	// sid.volume viene aggiornato da WriteRegister quando $D418 è scritto.
	// Qui salviamo quel valore in sampleBuf.
	sid.sampleBuf[sid.sampleBufIdx] = sid.volume
	sid.sampleBufIdx = (sid.sampleBufIdx + 1) % SampleBufSize
}

// Update processes and sends the audio buffer to the audio render system for playback.
func (sid *SID) Update() {
	sid.calcSoundBuffer()

	//TODO RIMUOVERE 2 * sid.fragSize e pos
	soundBufferSamples := 2 * sid.fragSize
	sid.player.Write(sid.soundBuffer, 0, soundBufferSamples)
}

// calcSoundBuffer generates the audio samples for the current block, applying mixing, filtering, and volume adjustments.
func (sid *SID) calcSoundBuffer() {
	const halfSampleBufSize = SampleBufSize / 2 // 624 / 2 = 312
	// Numero di campioni audio da generare in questa chiamata (per 20ms a 44.1kHz)
	numAudioSamplesInBlock := sid.fragSize // Assumiamo sid.fragSize = 882
	// Numero di aggiornamenti del volume da sampleBuf attesi per questo blocco
	// (corrispondenti alle chiamate a Prepare() in un frame PAL)
	numVolumeUpdatesInBlock := halfSampleBufSize
	// Calcola l'indice di partenza da cui leggere in sampleBuf per questo blocco.
	// sid.sampleBufIdx è il puntatore di SCRITTURA (dove Prepare scriverà il prossimo valore).
	// Dobbiamo leggere i 312 valori scritti NEL FRAME PRECEDENTE.
	// Quindi, l'ultimo valore scritto per il blocco precedente è a (sid.sampleBufIdx - 1 + SampleBufSize) % SampleBufSize
	// Il primo valore per il blocco precedente (di 312 valori) è a
	// (sid.sampleBufIdx - numVolumeUpdatesInBlock + SampleBufSize) % SampleBufSize
	currentVolumeBufferReadIdx := (sid.sampleBufIdx - numVolumeUpdatesInBlock + SampleBufSize) % SampleBufSize
	// Legge il primo valore di volume che sarà usato per i primi campioni audio.
	currentVolumeValue := sid.sampleBuf[currentVolumeBufferReadIdx]
	// Rapporto per determinare quando passare al successivo valore di volume da sampleBuf
	audioSamplesPerVolumeStep := float64(numAudioSamplesInBlock) / float64(numVolumeUpdatesInBlock) // Circa 882 / 312 = 2.8269...
	// `nextChangeAtAudioSampleIdx` è l'indice del campione audio `idx` (a partire da 0.0)
	// al quale dovremmo passare al prossimo valore di volume da `sampleBuf`.
	// Il primo cambio avverrà dopo `audioSamplesPerVolumeStep` campioni.
	nextChangeAtAudioSampleIdx := audioSamplesPerVolumeStep
	// `volumeStepsTaken` conta quanti dei 312 valori di volume abbiamo già utilizzato.
	// Inizia da 0 perché il primo valore (indice 0 dei 312) è già in currentVolumeValue.
	volumeStepsTaken := 0
	for idx := 0; idx < numAudioSamplesInBlock; idx++ { // Loop per 882 campioni audio
		// Controlla se è il momento di aggiornare il currentVolumeValue leggendo il prossimo
		// valore da sampleBuf.
		// Questo avviene quando l'indice del campione audio corrente (idx) supera o eguaglia
		// la soglia calcolata (nextChangeAtAudioSampleIdx).
		// Ci assicuriamo anche di non superare il numero di aggiornamenti del volume disponibili.
		if float64(idx) >= nextChangeAtAudioSampleIdx && volumeStepsTaken < numVolumeUpdatesInBlock-1 {
			volumeStepsTaken++
			// Avanza all'indice successivo nel ring buffer sampleBuf
			currentVolumeBufferReadIdx = (currentVolumeBufferReadIdx + 1) % SampleBufSize
			currentVolumeValue = sid.sampleBuf[currentVolumeBufferReadIdx]
			// Calcola la soglia per il prossimo cambio di volume
			nextChangeAtAudioSampleIdx += audioSamplesPerVolumeStep
		}

		// Voice Mixing
		sumOutputNonFiltered := int32(0)
		sumOutputFiltered := int32(0)

		for _, voice := range sid.voices {
			voice.ComputeEnvelopeGenerators()
			effectiveEnvelope := voice.EgLevel() >> 16 // 8-bit envelope
			if voice.IsMuted() {
				continue
			}
			voice.UpdateCount()
			waveOutput := voice.ComputeWaveForm() // La gestione del bit TEST è interna
			signedWaveOutput := int32(int16(waveOutput ^ 0x8000))
			voiceContribution := signedWaveOutput * int32(effectiveEnvelope)
			if voice.Filter() != 0 {
				sumOutputFiltered += voiceContribution
			} else {
				sumOutputNonFiltered += voiceContribution
			}
		}
		sumOutputFiltered = sid.filters.Compute(sumOutputFiltered)
		mixedSignal := sumOutputNonFiltered + sumOutputFiltered

		// Applica il volume letto da sampleBuf (0-15)
		// Questo currentVolumeValue cambia circa ogni 2-3 campioni audio.
		volumeAppliedSignal := (mixedSignal * int32(currentVolumeValue)) / 15

		finalSampleValue := volumeAppliedSignal >> 10 // Scala verso il basso

		sid.soundBuffer[idx] = uint32(finalSampleValue)
	}
}
