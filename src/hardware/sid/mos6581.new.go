package mos6581

/*
import (
	"fmt"
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

// LatencyMin represents the minimum threshold for network latency in milliseconds.
// LatencyMax represents the maximum threshold for network latency in milliseconds.
// LatencyAvg represents the average threshold for network latency in milliseconds.
const (
	LatencyMin = 80
	LatencyMax = 120
	LatencyAvg = 280
)

type EGState int

// SID represents a chip emulation containing configurations, registers, and audio handling functionality.
type SID struct {
	*component.BaseComponent
	registers []uint8
	//history      [][]uint8
	cfg *config.Config
	//audioBuilder *AudioBuilder
	reflect *SidReflect
	//historyCount int
	player      references.IAudioRender
	fragSize    int   // samples, not bytes
	bufferFrags int   // frags the in buffer
	bufferSize  int   // bytes, not samples
	volume      uint8 // Master volume
	voices             []*Voice // Data for 3 voices
	sampleBuf          []uint8  // Buffer for sampled voices
	sampleBufIdx       int      // Index in sample_buf for writing
	soundBuffer        []uint32
	soundBufferSamples int
	toOutput           int
	sbPos              int
	sbCurrentPosition  int
	divisor            int
	registerToVoice    []*Voice
	filters            *Filters
	divisorTable       *DivisorTable
	lead               *Lead
}

// NewSID creates a new SID instance with a specified parent ID and suffix, initializing its registers and settings.
func NewSID(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *SID {
	s := &SID{
		BaseComponent: component.NewBaseComponent(),
		registers:     make([]uint8, RegisterCount),
		//history:       make([][]uint8, RegisterHistory),
		player: nil,
		cfg:    nil,
		//audioBuilder: nil,
		//historyCount:  0,
	}
	//for x := range s.history {
	//	s.history[x] = make([]uint8, RegisterCount)
	//}
	s.BaseComponent.Register(factory, parent, Identifier(), s, references.IdISID(s, label, instance))
	s.reflect = NewSidReflect(s)
	return s
}

func (sid *SID) Setup() error {
	sid.cfg = sid.GetFactory().GetConfig()
	sid.cfg.Bind(sid.onConfigChanged)
	return nil
}

func (sid *SID) Bind(_ references.ISIDSocket, fragFreq int, rasters int) error {
	//sid.audioBuilder = NewAudioBuilder(sid.GetFactory().GetIAudioRender(), true, fragFreq, rasters)

	fragSize := SampleFreq / fragFreq     // samples, not bytes
	bufferSize := 2 * fragSize * fragFreq // bytes, not samples

	sid.player = sid.GetFactory().GetIAudioRender()
	sid.sampleBuf = make([]uint8, SampleBufSize)
	sid.divisorTable = NewDivisorTable(rasters, fragFreq)
	sid.voices = nil
	sid.fragSize = fragSize
	sid.bufferFrags = fragFreq
	sid.bufferSize = bufferSize
	sid.registerToVoice = make([]*Voice, RegisterCount)
	sid.filters = NewFilters()
	sid.soundBuffer = make([]uint32, 2*fragSize)
	sid.lead = NewLead(fragFreq)
	sid.sbCurrentPosition = 0

	voice0 := NewVoice(0)
	voice1 := NewVoice(1)
	voice2 := NewVoice(2)

	voice0.Setup(voice2, voice1)
	voice1.Setup(voice0, voice2)
	voice2.Setup(voice1, voice0)

	sid.voices = append(sid.voices, voice0, voice1, voice2)

	for x := range sid.registerToVoice {
		vIdx := (x / 7) % len(sid.voices)
		sid.registerToVoice[x] = sid.voices[vIdx]
	}

	sid.Reset()

	return nil
}

func (sid *SID) Connect() error {
	return nil
}

func (sid *SID) Internal() bool {
	return false
}

// Emulate processes the main emulation logic for the SID component, handling internal updates and state changes.
func (sid *SID) Emulate() {

}

func (sid *SID) EmulationRequired() bool {
	return false
}

// SetPotX sets the value of the POT X register in the SID chip using the given 8-bit value.
func (sid *SID) SetPotX(pot uint8) {
	// PX7 PX6 PX5 PX4 PX3 PX2 PX1 PX0
	sid.registers[potXRegisterIndex] = pot
}

// SetPotY sets the value of the POT Y register to the specified 8-bit value in the SID chip.
func (sid *SID) SetPotY(pot uint8) {
	//PY7 PY6 PY5 PY4 PY3 PY2 PY1 PY0
	sid.registers[potYRegisterIndex] = pot
}

// onConfigChanged is triggered when the configuration bound to the SID instance changes.
func (sid *SID) onConfigChanged() {
	//TODO
}

// Reset initializes all SID registers to 0 and sets default values for PotX and PotY. It also resets the audio builder.
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
	sid.toOutput = 0
	sid.divisor = 0
	for x := range sid.soundBuffer {
		sid.soundBuffer[x] = 0
	}
	sid.lead.Reset()
	sid.sbPos = 0
	sid.sbCurrentPosition = 0

	//PADDLE TEST
	//sid.WriteRegister(0xDC00, 0x40) //Control-Port 1 selected
	//sid.WriteRegister(0xD419, 0x7F) //Paddle X value
	//sid.WriteRegister(0xD419, 0x7F) //Paddle Y value
}

// ReadRegister retrieves the value from the specified address within the SID's registers. Only the lower 5 bits of the address are used.
func (sid *SID) ReadRegister(addr uint16) uint8 {
	reg := addr & 0x1f
	v := sid.registers[reg]
	return v
}

// WriteRegister writes an 8-bit value to a specific register at the given address by mapping it within a 32-register range.
func (sid *SID) WriteRegister(addr uint16, data uint8) {
	reg := addr & 0x1f
	sid.registers[reg] = data

	switch reg {
	case 0, 7, 14:
		voice := sid.registerToVoice[reg]
		voice.UpdateFreqA(data)
	case 1, 8, 15:
		voice := sid.registerToVoice[reg]
		voice.UpdateFreqB(data)
	case 2, 9, 16:
		voice := sid.registerToVoice[reg]
		voice.UpdatePulseWidthA(data)
	case 3, 10, 17:
		voice := sid.registerToVoice[reg]
		voice.UpdatePulseWidthB(data)
	case 4, 11, 18:
		voice := sid.registerToVoice[reg]
		voice.UpdateWaveForm(data)
	case 5, 12, 19:
		voice := sid.registerToVoice[reg]
		voice.UpdateEnvelopeGenerators(data)
	case 6, 13, 20:
		voice := sid.registerToVoice[reg]
		voice.UpdateSustainLevel(data)
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

var _callFillBuffer int = 0

// Prepare loads necessary SID register values into the AudioBuilder for audio processing.
func (sid *SID) Prepare() {
	sid.sampleBuf[sid.sampleBufIdx] = sid.volume
	sid.sampleBufIdx = (sid.sampleBufIdx + 1) % SampleBufSize
	sid.divisor += SampleFreq
	sid.toOutput += int(sid.divisorTable.GetOut(sid.divisor))
	sid.divisor = int(sid.divisorTable.GetDivisor(sid.divisor))

	// Calculate the sound data only when we have enough to fill the buffer entirely.
	if sid.toOutput >= sid.fragSize {
		sid.toOutput -= sid.fragSize
		sid.fillBuffer()
		_callFillBuffer++
	}
}

// Update triggers the audioBuilder's internal Update method, updating audio sampling and processing within the SID.
func (sid *SID) Update() {
	if _callFillBuffer > 1 {
		fmt.Println("ERROR")
	}
	_callFillBuffer = 0
	sid.player.Write(sid.soundBuffer, 0, sid.soundBufferSamples)
}

// GetLastByte retrieves the last byte from the SID's internal state or configuration.
func (sid *SID) GetLastByte() uint8 {
	return 0
}

// calcBuffer processes audio data for the given buffer by applying waveform calculations and filters to generate output.
// It iterates through the buffer, calculating the mixed output for each audio voice and summing the results.
// The method uses sample data, master volume, envelopes, and other voice parameters for precise audio mixing.
// Filtered and unfiltered outputs are computed and combined, then written into the provided buffer.
func (sid *SID) calcBuffer(buf []uint32, sampleBufIdx int) {
	const halfBufSize = SampleBufSize / 2
	const samples = ((0x138 * 50) << 16) / SampleFreq
	sampleCount := uint32((sampleBufIdx + halfBufSize) << 16)
	count := len(buf)
	count >>= 1 // 16 bit mono output, count is in bytes
	//count >>= 2; // 16 bit stereo output, count is in bytes
	idx := 0
	for ; count >= 0; count, idx = count-1, idx+1 {
		sumOutputFilter := int32(0)
		// Get current master volume from sample buffer, calculate sampled voices
		masterVolume := sid.sampleBuf[(sampleCount>>16)%SampleBufSize]
		sampleCount += samples
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
		buf[idx] = uint32((sumOutput + sumOutputFilter) >> 10)
	}
}

// write processes audio data and manages buffer positioning, adjusting for audio playback timing and synchronization.
func (sid *SID) fillBuffer() {
	// Compute the current lead in frags.
	//currentPosition := dr.player.GetCurrentPosition()
	currentPosition := sid.sbCurrentPosition
	if currentPosition == -1 {
		return
	}
	leadInBytes := (sid.sbPos - currentPosition + sid.bufferSize) % sid.bufferSize
	if leadInBytes >= sid.bufferSize/2 {
		leadInBytes -= sid.bufferSize
	}
	leadInFrags := leadInBytes / 2 * sid.fragSize
	avgLead, ok := sid.lead.Average(leadInFrags)
	if !ok {
		return
	}
	// Calculate one frag
	nSamples := sid.fragSize
	sid.calcBuffer(sid.soundBuffer, sid.sampleBufIdx)
	// If we're getting too far behind the audio add an extra frag.
	if avgLead < sid.lead.GetLoWater() {
		sid.lead.Update()
		//fmt.Printf("Adding an extra frag...\n");
		sid.calcBuffer(sid.soundBuffer[sid.fragSize:], sid.sampleBufIdx)
		nSamples += sid.fragSize
	}
	// Write the frags to the player and update out write position.
	//currPos := dr.sbPos
	samples := 2 * nSamples
	sid.sbCurrentPosition += samples
	sid.sbPos = (sid.sbPos + samples) % sid.bufferSize
	sid.soundBufferSamples = samples
}
*/
