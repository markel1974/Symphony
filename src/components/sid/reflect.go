package mos6581

import (
	"fmt"
	"github.com/markel1974/c64emu/src/components/board"
)

// Costanti per gli indici dei registri del SID (all'interno della slice 'registers').
const (
	freqLO1 = 0
	freqHI1 = 1
	pwLO1   = 2
	pwHI1   = 3
	cr1     = 4 // Control Register - Voice 1
	ad1     = 5 // Attack/Decay - Voice 1
	sr1     = 6 // Sustain/Release - Voice 1
	freqLO2 = 7
	freqHI2 = 8
	pwLO2   = 9
	pwHI2   = 10
	cr2     = 11
	ad2     = 12
	sr2     = 13
	freqLO3 = 14
	freqHI3 = 15
	pwLO3   = 16
	pwHI3   = 17
	cr3     = 18
	ad3     = 19
	sr3     = 20
	fcLO    = 21 // Filter Cutoff Low
	fcHI    = 22
	resFilt = 23 // Resonance and Filter Control
	modeVol = 24 // Mode and Volume
	potX    = 25 // Potentiometer X (read-only)
	potY    = 26 // Potentiometer Y (read-only)
	osc3    = 27 // Oscillator 3 / Random Number (read-only)
	env3    = 28 // Envelope 3 (read-only)
)

type Reflect struct {
	props *board.Properties
	sid   *SID
}

func NewReflect(s *SID) *Reflect {
	r := &Reflect{
		props: nil,
		sid:   s,
	}
	r.props = board.NewProperties(r.RunCommand)
	kind := r.sid.registers[0] // Assuming all registers are uint8
	// Voice 1
	r.props.Add(board.CreatePropertyInfo("freqLO1", kind, "Voice 1: Frequency Low Byte (0xD400)", false, r.getFreqLO1, r.setFreqLO1))
	r.props.Add(board.CreatePropertyInfo("freqHI1", kind, "Voice 1: Frequency High Byte (0xD401)", false, r.getFreqHI1, r.setFreqHI1))
	r.props.Add(board.CreatePropertyInfo("pwLO1", kind, "Voice 1: Pulse Width Low Byte (0xD402)", false, r.getPwLO1, r.setPwLO1))
	r.props.Add(board.CreatePropertyInfo("pwHI1", kind, "Voice 1: Pulse Width High Byte (0xD403)", false, r.getPwHI1, r.setPwHI1))
	r.props.Add(board.CreatePropertyInfo("cr1", kind, "Voice 1: Control Register (0xD404) - Gate, Sync, Ring, Test, Waveform", false, r.getCr1, r.setCr1))
	r.props.Add(board.CreatePropertyInfo("ad1", kind, "Voice 1: Attack/Decay (0xD405) - Attack Rate, Decay Rate", false, r.getAd1, r.setAd1))
	r.props.Add(board.CreatePropertyInfo("sr1", kind, "Voice 1: Sustain/Release (0xD406) - Sustain Level, Release Rate", false, r.getSr1, r.setSr1))
	// Voice 2
	r.props.Add(board.CreatePropertyInfo("freqLO2", kind, "Voice 2: Frequency Low Byte (0xD407)", false, r.getFreqLO2, r.setFreqLO2))
	r.props.Add(board.CreatePropertyInfo("freqHI2", kind, "Voice 2: Frequency High Byte (0xD408)", false, r.getFreqHI2, r.setFreqHI2))
	r.props.Add(board.CreatePropertyInfo("pwLO2", kind, "Voice 2: Pulse Width Low Byte (0xD409)", false, r.getPwLO2, r.setPwLO2))
	r.props.Add(board.CreatePropertyInfo("pwHI2", kind, "Voice 2: Pulse Width High Byte (0xD40A)", false, r.getPwHI2, r.setPwHI2))
	r.props.Add(board.CreatePropertyInfo("cr2", kind, "Voice 2: Control Register (0xD40B) - Gate, Sync, Ring, Test, Waveform", false, r.getCr2, r.setCr2))
	r.props.Add(board.CreatePropertyInfo("ad2", kind, "Voice 2: Attack/Decay (0xD40C) - Attack Rate, Decay Rate", false, r.getAd2, r.setAd2))
	r.props.Add(board.CreatePropertyInfo("sr2", kind, "Voice 2: Sustain/Release (0xD40D) - Sustain Level, Release Rate", false, r.getSr2, r.setSr2))
	// Voice 3
	r.props.Add(board.CreatePropertyInfo("freqLO3", kind, "Voice 3: Frequency Low Byte (0xD40E)", false, r.getFreqLO3, r.setFreqLO3))
	r.props.Add(board.CreatePropertyInfo("freqHI3", kind, "Voice 3: Frequency High Byte (0xD40F)", false, r.getFreqHI3, r.setFreqHI3))
	r.props.Add(board.CreatePropertyInfo("pwLO3", kind, "Voice 3: Pulse Width Low Byte (0xD410)", false, r.getPwLO3, r.setPwLO3))
	r.props.Add(board.CreatePropertyInfo("pwHI3", kind, "Voice 3: Pulse Width High Byte (0xD411)", false, r.getPwHI3, r.setPwHI3))
	r.props.Add(board.CreatePropertyInfo("cr3", kind, "Voice 3: Control Register (0xD412) - Gate, Sync, Ring, Test, Waveform", false, r.getCr3, r.setCr3))
	r.props.Add(board.CreatePropertyInfo("ad3", kind, "Voice 3: Attack/Decay (0xD413) - Attack Rate, Decay Rate", false, r.getAd3, r.setAd3))
	r.props.Add(board.CreatePropertyInfo("sr3", kind, "Voice 3: Sustain/Release (0xD414) - Sustain Level, Release Rate", false, r.getSr3, r.setSr3))
	// Filtro
	r.props.Add(board.CreatePropertyInfo("fcLO", kind, "Filter Cutoff Frequency Low Byte (0xD415)", false, r.getFcLO, r.setFcLO))
	r.props.Add(board.CreatePropertyInfo("fcHI", kind, "Filter Cutoff Frequency High Byte (0xD416)", false, r.getFcHI, r.setFcHI))
	r.props.Add(board.CreatePropertyInfo("resFilt", kind, "Filter Resonance and Control Register (0xD417) - Resonance, Filter Mode, External Input", false, r.getResFilt, r.setResFilt))
	r.props.Add(board.CreatePropertyInfo("modeVol", kind, "SID Mode and Volume Register (0xD418) - Filter Mode, Voice 3 Mute, Volume", false, r.getModeVol, r.setModeVol))
	// POTX, POTY, OSC3, ENV3 (sola lettura)
	r.props.Add(board.CreatePropertyInfo("potX", kind, "POTX Register (0xD419) - Paddle X Input", false, r.getPotX, r.setPotX))
	r.props.Add(board.CreatePropertyInfo("potY", kind, "POTY Register (0xD41A) - Paddle Y Input", false, r.getPotY, r.setPotY))
	r.props.Add(board.CreatePropertyInfo("osc3", kind, "OSC3 Register (0xD41B) - Oscillator 3 Value", false, r.getOsc3, r.setOsc3))
	r.props.Add(board.CreatePropertyInfo("env3", kind, "ENV3 Register (0xD41C) - Envelope 3 Value", false, r.getEnv3, r.setEnv3))
	return r
}

// GetProperties restituisce la mappa delle proprietà del SID.
func (r *Reflect) GetProperties() *board.Properties {
	return r.props
}

func (r *Reflect) RunCommand(cmd string, args []string) (map[string]interface{}, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (s *Reflect) getFreqLO1() uint8  { return s.sid.registers[freqLO1] }
func (s *Reflect) setFreqLO1(v uint8) { s.sid.registers[freqLO1] = v }

func (s *Reflect) getFreqHI1() uint8  { return s.sid.registers[freqHI1] }
func (s *Reflect) setFreqHI1(v uint8) { s.sid.registers[freqHI1] = v }

func (s *Reflect) getPwLO1() uint8  { return s.sid.registers[pwLO1] }
func (s *Reflect) setPwLO1(v uint8) { s.sid.registers[pwLO1] = v }

func (s *Reflect) getPwHI1() uint8  { return s.sid.registers[pwHI1] }
func (s *Reflect) setPwHI1(v uint8) { s.sid.registers[pwHI1] = v }

func (s *Reflect) getCr1() uint8  { return s.sid.registers[cr1] }
func (s *Reflect) setCr1(v uint8) { s.sid.registers[cr1] = v }

func (s *Reflect) getAd1() uint8  { return s.sid.registers[ad1] }
func (s *Reflect) setAd1(v uint8) { s.sid.registers[ad1] = v }

func (s *Reflect) getSr1() uint8  { return s.sid.registers[sr1] }
func (s *Reflect) setSr1(v uint8) { s.sid.registers[sr1] = v }

func (s *Reflect) getFreqLO2() uint8  { return s.sid.registers[freqLO2] }
func (s *Reflect) setFreqLO2(v uint8) { s.sid.registers[freqLO2] = v }

func (s *Reflect) getFreqHI2() uint8  { return s.sid.registers[freqHI2] }
func (s *Reflect) setFreqHI2(v uint8) { s.sid.registers[freqHI2] = v }

func (s *Reflect) getPwLO2() uint8  { return s.sid.registers[pwLO2] }
func (s *Reflect) setPwLO2(v uint8) { s.sid.registers[pwLO2] = v }

func (s *Reflect) getPwHI2() uint8  { return s.sid.registers[pwHI2] }
func (s *Reflect) setPwHI2(v uint8) { s.sid.registers[pwHI2] = v }

func (s *Reflect) getCr2() uint8  { return s.sid.registers[cr2] }
func (s *Reflect) setCr2(v uint8) { s.sid.registers[cr2] = v }

func (s *Reflect) getAd2() uint8  { return s.sid.registers[ad2] }
func (s *Reflect) setAd2(v uint8) { s.sid.registers[ad2] = v }

func (s *Reflect) getSr2() uint8  { return s.sid.registers[sr2] }
func (s *Reflect) setSr2(v uint8) { s.sid.registers[sr2] = v }

func (s *Reflect) getFreqLO3() uint8  { return s.sid.registers[freqLO3] }
func (s *Reflect) setFreqLO3(v uint8) { s.sid.registers[freqLO3] = v }

func (s *Reflect) getFreqHI3() uint8  { return s.sid.registers[freqHI3] }
func (s *Reflect) setFreqHI3(v uint8) { s.sid.registers[freqHI3] = v }

func (s *Reflect) getPwLO3() uint8  { return s.sid.registers[pwLO3] }
func (s *Reflect) setPwLO3(v uint8) { s.sid.registers[pwLO3] = v }

func (s *Reflect) getPwHI3() uint8  { return s.sid.registers[pwHI3] }
func (s *Reflect) setPwHI3(v uint8) { s.sid.registers[pwHI3] = v }

func (s *Reflect) getCr3() uint8  { return s.sid.registers[cr3] }
func (s *Reflect) setCr3(v uint8) { s.sid.registers[cr3] = v }

func (s *Reflect) getAd3() uint8  { return s.sid.registers[ad3] }
func (s *Reflect) setAd3(v uint8) { s.sid.registers[ad3] = v }

func (s *Reflect) getSr3() uint8  { return s.sid.registers[sr3] }
func (s *Reflect) setSr3(v uint8) { s.sid.registers[sr3] = v }

func (s *Reflect) getFcLO() uint8  { return s.sid.registers[fcLO] }
func (s *Reflect) setFcLO(v uint8) { s.sid.registers[fcLO] = v }

func (s *Reflect) getFcHI() uint8  { return s.sid.registers[fcHI] }
func (s *Reflect) setFcHI(v uint8) { s.sid.registers[fcHI] = v }

func (s *Reflect) getResFilt() uint8  { return s.sid.registers[resFilt] }
func (s *Reflect) setResFilt(v uint8) { s.sid.registers[resFilt] = v }

func (s *Reflect) getModeVol() uint8  { return s.sid.registers[modeVol] }
func (s *Reflect) setModeVol(v uint8) { s.sid.registers[modeVol] = v }

func (s *Reflect) getPotX() uint8  { return s.sid.registers[potX] }
func (s *Reflect) setPotX(v uint8) { s.sid.registers[potX] = v }

func (s *Reflect) getPotY() uint8  { return s.sid.registers[potY] }
func (s *Reflect) setPotY(v uint8) { s.sid.registers[potY] = v }

func (s *Reflect) getOsc3() uint8  { return s.sid.registers[osc3] }
func (s *Reflect) setOsc3(v uint8) { s.sid.registers[osc3] = v }

func (s *Reflect) getEnv3() uint8  { return s.sid.registers[env3] }
func (s *Reflect) setEnv3(v uint8) { s.sid.registers[env3] = v }
