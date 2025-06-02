package mos6581

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// potXRegisterIndex represents the register index for the X-axis potentiometer.
// potYRegisterIndex represents the register index for the Y-axis potentiometer.
const (
	potXRegisterIndex = 25
	potYRegisterIndex = 26
)

// RegMax defines the maximum number of SID chip registers plus one, typically used as a boundary for register operations.
//const RegMax = 0x1f + 1

// WriteFn defines a function type that processes an 8-bit unsigned integer as input.
type WriteFn func(data uint8)

// ReadFn defines a function type used to read an 8-bit value from a given 16-bit address.
type ReadFn func(addr uint16) uint8

// SID represents a Sound Interface Device, a component used to generate and handle audio synthesis in the system.
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
	writes       [RegisterCount]WriteFn
	reads        [RegisterCount]ReadFn
	//divisorTable    *DivisorTable
	//divisor         int
}

// NewSID initializes and returns a new SID component instance with the given parent, factory, label, and instance number.
// It sets up the base component, registers, and reflection interface for the sound interface device.
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

// Setup initializes the SID instance by configuring settings, binding configuration changes, and setting up registers.
func (sid *SID) Setup() error {
	sid.cfg = sid.GetFactory().GetConfig()
	sid.cfg.Bind(sid.onConfigChanged)
	return nil
}

// Bind initializes the SID instance with the given socket, fragment frequency, and raster count, returning an error if any.
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

	sid.writes = sid.createWriteRegister()
	sid.reads = sid.createReadRegister()
	sid.Reset()

	return nil
}

// Connect establishes a connection using the SID configuration and returns an error if it fails.
func (sid *SID) Connect() error {
	return nil
}

// Internal determines if the SID instance is internal, returning true if it is, otherwise false.
func (sid *SID) Internal() bool {
	return false
}

// Emulate executes the emulation process for the SID object, simulating its behavior according to the defined parameters.
func (sid *SID) Emulate() {

}

// EmulationRequired checks whether emulation is required for the specified SID and returns a boolean result.
func (sid *SID) EmulationRequired() bool {
	return false
}

// onConfigChanged handles updates or adjustments to the configuration settings dynamically during runtime.
func (sid *SID) onConfigChanged() {
	//TODO
}

// Reset reinitializes the SID object by clearing its registers, buffers, and resetting all internal components to defaults.
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

// Prepare updates the sample buffer with the current volume and increments the buffer index in a circular manner.
func (sid *SID) Prepare() {
	// sid.volume viene aggiornato da WriteRegister quando $D418 è scritto.
	// Qui salviamo quel valore in sampleBuf.
	sid.sampleBuf[sid.sampleBufIdx] = sid.volume
	sid.sampleBufIdx = (sid.sampleBufIdx + 1) % SampleBufSize
}

// Update processes the sound buffer and writes updated sound data to the audio player.
func (sid *SID) Update() {
	sid.calcSoundBuffer()

	//TODO RIMUOVERE 2 * sid.fragSize e pos
	soundBufferSamples := 2 * sid.fragSize
	sid.player.Write(sid.soundBuffer, 0, soundBufferSamples)
}

// ReadRegister reads the value of a specified SID register identified by the provided address.
func (sid *SID) ReadRegister(addr uint16) uint8 {
	reg := addr & 0x1f
	return sid.reads[reg](addr)
}

// WriteRegister writes a value to a specific register within the range of the SID's addressable space.
func (sid *SID) WriteRegister(addr uint16, data uint8) {
	reg := addr & 0x1f
	sid.registers[reg] = data
	sid.writes[reg](data)
}

// SetPotX sets the value of the POTX (paddle) register to control the position of the X-axis paddle input.
func (sid *SID) SetPotX(pot uint8) {
	// PX7 PX6 PX5 PX4 PX3 PX2 PX1 PX0
	sid.registers[potXRegisterIndex] = pot
}

// SetPotY sets the pot value for the Y-axis register of the SID. This value updates the corresponding SID register directly.
func (sid *SID) SetPotY(pot uint8) {
	//PY7 PY6 PY5 PY4 PY3 PY2 PY1 PY0
	sid.registers[potYRegisterIndex] = pot
}

// GetLastByte retrieves the last byte of the SID structure and returns it as an unsigned 8-bit integer.
func (sid *SID) GetLastByte() uint8 {
	return 0
}

// calcSoundBuffer generates audio samples for the current block, applying volume changes and filtering for each voice.
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

// createReadRegister initializes and returns an array of ReadFn functions mapped to SID register addresses.
// It sets a default read function for most registers and customized functions for specific addresses like 27 and 28.
func (sid *SID) createReadRegister() [RegisterCount]ReadFn {
	var reads [RegisterCount]ReadFn
	defaultFn := func(addr uint16) uint8 {
		reg := addr & 0x1f
		return sid.registers[reg]
	}
	for idx := range reads {
		reads[idx] = defaultFn
	}
	reads[27] = sid.readVoice2Waveform
	reads[28] = sid.readVoice2EgLevel
	return reads
}

// createWriteRegister initializes and returns an array of WriteFn mapped to SID register write operations.
func (sid *SID) createWriteRegister() [RegisterCount]WriteFn {
	var writes [RegisterCount]WriteFn
	emptyFn := func(data uint8) {}
	for idx := range writes {
		writes[idx] = emptyFn
	}
	writes[0] = sid.writeVoice0UpdateFreqA
	writes[1] = sid.writeVoice0UpdateFreqB
	writes[2] = sid.writeVoice0UpdatePulseWidthA
	writes[3] = sid.writeVoice0UpdatePulseWidthB
	writes[4] = sid.writeVoice0UpdateWaveForm
	writes[5] = sid.writeVoice0UpdateEnvelopeGenerators
	writes[6] = sid.writeVoice0UpdateSustainLevel
	writes[7] = sid.writeVoice1UpdateFreqA
	writes[8] = sid.writeVoice1UpdateFreqB
	writes[9] = sid.writeVoice1UpdatePulseWidthA
	writes[10] = sid.writeVoice1UpdatePulseWidthB
	writes[11] = sid.writeVoice1UpdateWaveForm
	writes[12] = sid.writeVoice1UpdateEnvelopeGenerators
	writes[13] = sid.writeVoice1UpdateSustainLevel
	writes[14] = sid.writeVoice2UpdateFreqA
	writes[15] = sid.writeVoice2UpdateFreqB
	writes[16] = sid.writeVoice2UpdatePulseWidthA
	writes[17] = sid.writeVoice2UpdatePulseWidthB
	writes[18] = sid.writeVoice2UpdateWaveForm
	writes[19] = sid.writeVoice2UpdateEnvelopeGenerators
	writes[20] = sid.writeVoice2UpdateSustainLevel
	writes[21] = sid.writeFiltersUpdateFreqLow
	writes[22] = sid.writeFiltersUpdateFreqHigh
	writes[23] = sid.writeFiltersRegister
	writes[24] = sid.writeMasterVolumeAndFilterType
	return writes
}

// readVoice2Waveform retrieves the most significant byte (MSB) of the current level of oscillator for voice 2.
// It uses the ComputeWaveForm() function from voice.go, which returns a uint16 value.
func (sid *SID) readVoice2Waveform(_ uint16) uint8 {
	// OSC3 - Oscillator 3 Value ($D41B)
	// Restituisce il byte più significativo (MSB) dell'output corrente
	// dell'oscillatore (waveform) della voce 2.
	// La funzione ComputeWaveForm() in voice.go restituisce un uint16.
	return uint8(sid.voices[2].ComputeWaveForm() >> 8)
}

// readVoice2EgLevel retrieves the most significant byte (MSB) of the current envelope output for voice 2.
func (sid *SID) readVoice2EgLevel(_ uint16) uint8 {
	// ENV3 - Envelope 3 Value ($D41C)
	// Restituisce il byte più significativo (MSB) del livello corrente
	// dell'inviluppo (Envelope Generator) della voce 2.
	// La funzione EgLevel() in voice.go restituisce un uint32 (valore a 24 bit).
	return uint8(sid.voices[2].EgLevel() >> 16)
}

// Voice 0

// writeVoice0UpdateFreqA updates frequency parameter A for voice 0 with the specified data value.
func (sid *SID) writeVoice0UpdateFreqA(data uint8) {
	sid.voices[0].UpdateFreqA(data)
}

// writeVoice0UpdateFreqB writes the frequency update value to voice 0 using the provided data byte.
func (sid *SID) writeVoice0UpdateFreqB(data uint8) {
	sid.voices[0].UpdateFreqB(data)
}

// writeVoice0UpdatePulseWidthA updates the pulse width modulation of Voice 0 through the given data value.
func (sid *SID) writeVoice0UpdatePulseWidthA(data uint8) {
	sid.voices[0].UpdatePulseWidthA(data)
}

// writeVoice0UpdatePulseWidthB updates the pulse width value for voice 0 using the provided 8-bit data.
func (sid *SID) writeVoice0UpdatePulseWidthB(data uint8) {
	sid.voices[0].UpdatePulseWidthB(data)
}

// writeVoice0UpdateWaveForm updates the waveform for voice 0 using the provided data value.
func (sid *SID) writeVoice0UpdateWaveForm(data uint8) {
	sid.voices[0].UpdateWaveForm(data)
}

// writeVoice0UpdateEnvelopeGenerators updates the envelope generators for voice 0 with the specified data.
func (sid *SID) writeVoice0UpdateEnvelopeGenerators(data uint8) {
	sid.voices[0].UpdateEnvelopeGenerators(data)
}

// writeVoice0UpdateSustainLevel updates the sustain level of voice 0 using the provided data value.
func (sid *SID) writeVoice0UpdateSustainLevel(data uint8) {
	sid.voices[0].UpdateSustainLevel(data)
}

// Voice 1

// writeVoice1UpdateFreqA writes the provided data to update the frequency parameter A for voice 1 of the SID.
func (sid *SID) writeVoice1UpdateFreqA(data uint8) {
	sid.voices[1].UpdateFreqA(data)
}

// writeVoice1UpdateFreqB updates the frequency B of voice 1 using the provided data value.
func (sid *SID) writeVoice1UpdateFreqB(data uint8) {
	sid.voices[1].UpdateFreqB(data)
}

// writeVoice1UpdatePulseWidthA updates the pulse width of Voice 1 using the provided 8-bit data.
func (sid *SID) writeVoice1UpdatePulseWidthA(data uint8) {
	sid.voices[1].UpdatePulseWidthA(data)
}

// writeVoice1UpdatePulseWidthB updates the pulse width modulation for voice 1 using the provided data value.
func (sid *SID) writeVoice1UpdatePulseWidthB(data uint8) {
	sid.voices[1].UpdatePulseWidthB(data)
}

// writeVoice1UpdateWaveForm updates the waveform of voice 1 in the SID chip using the provided data.
func (sid *SID) writeVoice1UpdateWaveForm(data uint8) {
	sid.voices[1].UpdateWaveForm(data)
}

// writeVoice1UpdateEnvelopeGenerators updates the envelope generators for voice 1 using the provided data.
func (sid *SID) writeVoice1UpdateEnvelopeGenerators(data uint8) {
	sid.voices[1].UpdateEnvelopeGenerators(data)
}

// writeVoice1UpdateSustainLevel updates the sustain level of voice 1 with the provided data value.
func (sid *SID) writeVoice1UpdateSustainLevel(data uint8) {
	sid.voices[1].UpdateSustainLevel(data)
}

// Voice 2

// writeVoice2UpdateFreqA writes a frequency update value to voice 2's frequency register A using the provided data.
func (sid *SID) writeVoice2UpdateFreqA(data uint8) {
	sid.voices[2].UpdateFreqA(data)
}

// writeVoice2UpdateFreqB updates the frequency B register of voice 2 with the provided 8-bit data.
func (sid *SID) writeVoice2UpdateFreqB(data uint8) {
	sid.voices[2].UpdateFreqB(data)
}

// writeVoice2UpdatePulseWidthA updates the pulse width modulation for voice 2 with the provided data value.
func (sid *SID) writeVoice2UpdatePulseWidthA(data uint8) {
	sid.voices[2].UpdatePulseWidthA(data)
}

// writeVoice2UpdatePulseWidthB updates the pulse width modulation parameter B for voice 2 using the given data.
func (sid *SID) writeVoice2UpdatePulseWidthB(data uint8) {
	sid.voices[2].UpdatePulseWidthB(data)
}

// writeVoice2UpdateWaveForm updates the waveform data for the second voice of the SID chip using the provided data value.
func (sid *SID) writeVoice2UpdateWaveForm(data uint8) {
	sid.voices[2].UpdateWaveForm(data)
}

// writeVoice2UpdateEnvelopeGenerators updates the envelope generators for voice 2 using the provided data.
func (sid *SID) writeVoice2UpdateEnvelopeGenerators(data uint8) {
	sid.voices[2].UpdateEnvelopeGenerators(data)
}

// writeVoice2UpdateSustainLevel updates the sustain level for voice 2 using the provided data value.
func (sid *SID) writeVoice2UpdateSustainLevel(data uint8) {
	sid.voices[2].UpdateSustainLevel(data)
}

// writeFiltersUpdateFreqLow updates the low byte of the filter frequency using the provided data value.
func (sid *SID) writeFiltersUpdateFreqLow(data uint8) {
	sid.filters.UpdateFreqLow(data)
}

// writeFiltersUpdateFreqHigh updates the high-frequency filter settings with the provided data masked to 3 bits.
func (sid *SID) writeFiltersUpdateFreqHigh(data uint8) {
	sid.filters.UpdateFreqHigh(data & 0x07)
}

// writeFiltersRegister configures filter settings for the SID voices based on the provided data value.
func (sid *SID) writeFiltersRegister(data uint8) {
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
}

// writeMasterVolumeAndFilterType updates the master volume, filter type, and mute state based on the given input data.
func (sid *SID) writeMasterVolumeAndFilterType(data uint8) {
	mute := false //uint8(0)
	if (data & 0x80) != 0 {
		mute = true
	}
	sid.volume = data & 0xf
	sid.voices[2].mute = mute
	sid.filters.UpdateType(data)
}
